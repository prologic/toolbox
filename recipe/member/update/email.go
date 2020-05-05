package update

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_member"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_member"
	"github.com/watermint/toolbox/essentials/log/es_log"
	"github.com/watermint/toolbox/infra/api/api_parser"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/feed/fd_file"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
)

type MsgEmail struct {
	ProgressUpdate app_msg.Message
}

var (
	MEmail = app_msg.Apply(&MsgEmail{}).(*MsgEmail)
)

type EmailRow struct {
	FromEmail string `json:"from_email"`
	ToEmail   string `json:"to_email"`
}

type EmailWorker struct {
	transaction *EmailRow
	member      *mo_member.Member
	ctx         dbx_context.Context
	rep         rp_model.TransactionReport
	ctl         app_control.Control
}

func (z *EmailWorker) Exec() error {
	ui := z.ctl.UI()
	ui.Progress(MEmail.ProgressUpdate.With("EmailFrom", z.transaction.FromEmail).With("EmailTo", z.transaction.ToEmail))

	l := z.ctl.Log().With(es_log.Any("beforeMember", z.member))

	newEmail := &mo_member.Member{}
	if err := api_parser.ParseModelRaw(newEmail, z.member.Raw); err != nil {
		l.Debug("Unable to clone member data", es_log.Error(err))
		z.rep.Failure(err, z.transaction)
		return err
	}

	newEmail.Email = z.transaction.ToEmail
	newMember, err := sv_member.New(z.ctx).Update(newEmail)
	if err != nil {
		l.Debug("API returned an error", es_log.Error(err))
		z.rep.Failure(err, z.transaction)
		return err
	}

	z.rep.Success(z.transaction, newMember)
	return nil
}

type Email struct {
	Peer                dbx_conn.ConnBusinessMgmt
	File                fd_file.RowFeed
	UpdateUnverified    bool
	OperationLog        rp_model.TransactionReport
	SkipSameFromToEmail app_msg.Message
	SkipUnverifiedEmail app_msg.Message
}

func (z *Email) Preset() {
	z.File.SetModel(&EmailRow{})
	z.OperationLog.SetModel(
		&EmailRow{},
		&mo_member.Member{},
		rp_model.HiddenColumns(
			"result.team_member_id",
			"result.familiar_name",
			"result.abbreviated_name",
			"result.member_folder_id",
			"result.external_id",
			"result.account_id",
			"result.persistent_id",
		),
	)
}

func (z *Email) Exec(c app_control.Control) error {
	l := c.Log()
	ctx := z.Peer.Context()

	members, err := sv_member.New(ctx).List()
	if err != nil {
		return err
	}
	emailToMember := mo_member.MapByEmail(members)

	err = z.OperationLog.Open()
	if err != nil {
		return err
	}

	q := c.NewQueue()
	err = z.File.EachRow(func(m interface{}, rowIndex int) error {
		row := m.(*EmailRow)
		ll := l.With(es_log.Any("row", row))

		if row.FromEmail == row.ToEmail {
			ll.Debug("Skip")
			z.OperationLog.Skip(z.SkipSameFromToEmail, row)
			return nil
		}

		member, ok := emailToMember[row.FromEmail]
		if !ok {
			ll.Debug("Member not found for email")
			z.OperationLog.Failure(errors.New("member not found for email"), row)
			return nil
		}

		if !member.EmailVerified && !z.UpdateUnverified {
			ll.Debug("Do not update unverified email")
			z.OperationLog.Skip(z.SkipUnverifiedEmail, row)
			return nil
		}

		q.Enqueue(&EmailWorker{
			transaction: row,
			member:      member,
			ctx:         ctx,
			rep:         z.OperationLog,
			ctl:         c,
		})

		return nil
	})
	q.Wait()
	return err
}

func (z *Email) Test(c app_control.Control) error {
	return nil
	// TODO: Replace TestResource with new stuff
	//l := c.Log()
	//key := rc_recipe.Key(z)
	//res, found := c.TestResource(key)
	//if !found || !res.IsArray() {
	//	l.Debug("SKIP: Test resource not found")
	//	return qt_errors.ErrorNotEnoughResource
	//}
	//
	//pair := make(map[string]string)
	//noExist := make(map[string]bool)
	//
	//for _, row := range res.Array() {
	//	from := row.Get("from").String()
	//	to := row.Get("to").String()
	//	exists := row.Get("exists").Bool()
	//
	//	if !dbx_util.RegexEmail.MatchString(from) || !dbx_util.RegexEmail.MatchString(to) {
	//		l.Error("from or to email address unmatched to email address format", es_log.String("from", from), es_log.String("to", to))
	//		return errors.New("invalid input")
	//	}
	//	pair[from] = to
	//	noExist[from] = !exists
	//}
	//
	//createCsv := func(path string, reverse bool) error {
	//	l.Info("Create test file", es_log.String("path", path))
	//	f, err := os.Create(path)
	//	if err != nil {
	//		l.Debug("Unable to create test file", es_log.Error(err))
	//		return err
	//	}
	//	cw := csv.NewWriter(f)
	//	if err := cw.Write([]string{"from_email", "to_email"}); err != nil {
	//		return err
	//	}
	//
	//	for k, v := range pair {
	//		if reverse {
	//			if err := cw.Write([]string{v, k}); err != nil {
	//				return err
	//			}
	//		} else {
	//			if err := cw.Write([]string{k, v}); err != nil {
	//				return err
	//			}
	//		}
	//	}
	//	cw.Flush()
	//	return f.Close()
	//}
	//
	//pathForward := filepath.Join(c.Workspace().Test(), "testdata_forward.csv")
	//pathBackward := filepath.Join(c.Workspace().Test(), "testdata_backward.csv")
	//
	//if err := createCsv(pathForward, false); err != nil {
	//	l.Error("Unable to create test file", es_log.String("pathForward", pathForward), es_log.Error(err))
	//	return err
	//}
	//if err := createCsv(pathBackward, true); err != nil {
	//	l.Error("Unable to create test file", es_log.String("pathForward", pathForward), es_log.Error(err))
	//	return err
	//}
	//
	//var lastErr error
	//
	//preserveReport := func(suffix string) error {
	//	repPath := c.Workspace().Report() + suffix
	//	err := os.Rename(c.Workspace().Report(), repPath)
	//	if err != nil {
	//		l.Warn("Unable to preserve forward report", es_log.Error(err))
	//		repPath = c.Workspace().Report()
	//	}
	//
	//	// create alt report folder
	//	err = os.MkdirAll(c.Workspace().Report(), 0701)
	//	if err != nil {
	//		l.Error("Unable to create workspace path", es_log.Error(err))
	//		return err
	//	}
	//	return nil
	//}
	//
	//scanReport := func() {
	//	resultPath := filepath.Join(c.Workspace().Report(), "operation_log.json")
	//	resultFile, err := os.Open(resultPath)
	//	if err != nil {
	//		l.Warn("Unable to open", es_log.Error(err))
	//	} else {
	//		scanner := bufio.NewScanner(resultFile)
	//		for scanner.Scan() {
	//			row := gjson.Parse(scanner.Text())
	//
	//			status := row.Get("status_tag").String()
	//			reason := row.Get("reason").String()
	//			inputFrom := row.Get("input.from_email").String()
	//			inputTo := row.Get("input.to_email").String()
	//			resultEmail := row.Get("result.email").String()
	//
	//			ll := l.With(
	//				es_log.String("status", status),
	//				es_log.String("inputFrom", inputFrom),
	//				es_log.String("inputTo", inputTo),
	//				es_log.String("resultEmail", resultEmail),
	//				es_log.String("reason", reason),
	//			)
	//			isNonExistent := noExist[inputFrom] || noExist[inputTo]
	//
	//			ll.Info("Feed file row", es_log.Bool("isNonExist", isNonExistent))
	//
	//			switch {
	//			case status == rp_model.StatusTagFailure && isNonExistent:
	//				ll.Info("Successfully failed for non existent")
	//			case status == rp_model.StatusTagFailure:
	//				ll.Warn("Unexpected failure")
	//				lastErr = errors.New("unexpected failure")
	//			case status == rp_model.StatusTagSuccess && isNonExistent:
	//				ll.Warn("Unexpected failure")
	//				lastErr = errors.New("unexpected failure")
	//			case status == rp_model.StatusTagSuccess:
	//				if inputTo == resultEmail {
	//					ll.Info("Successfully changed for non existent")
	//				} else {
	//					ll.Warn("Email address unchanged")
	//					lastErr = errors.New("email address unchanged")
	//				}
	//			default:
	//				ll.Warn("Unexpected status")
	//				lastErr = errors.New("unexpected status")
	//			}
	//		}
	//	}
	//}
	//
	//// forward
	//{
	//	lastErr = rc_exec.Exec(c, &Email{}, func(r rc_recipe.Recipe) {
	//		rr := r.(*Email)
	//		rr.UpdateUnverified = true
	//		rr.File.SetFilePath(pathForward)
	//	})
	//	if lastErr != nil {
	//		l.Warn("Error in backward operation")
	//	}
	//	scanReport()
	//	if err := preserveReport("_forward"); err != nil {
	//		return err
	//	}
	//}
	//
	//// backward
	//{
	//	lastErr = rc_exec.Exec(c, &Email{}, func(r rc_recipe.Recipe) {
	//		rr := r.(*Email)
	//		rr.UpdateUnverified = true
	//		rr.File.SetFilePath(pathBackward)
	//	})
	//	if lastErr != nil {
	//		l.Warn("Error in backward operation")
	//	}
	//	scanReport()
	//	if err := preserveReport("_backward"); err != nil {
	//		return err
	//	}
	//}
	//
	//return lastErr
}
