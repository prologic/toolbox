package batch

import (
	"errors"
	"github.com/watermint/toolbox/domain/model/mo_teamfolder"
	"github.com/watermint/toolbox/domain/service/sv_teamfolder"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/feed/fd_file"
	"github.com/watermint/toolbox/infra/recipe/rc_conn"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/quality/infra/qt_endtoend"
	"strings"
)

type Permdelete struct {
	ErrTeamFolderNotFound                 app_msg.Message
	ErrUnableToDelete                     app_msg.Message
	ErrUnableToRetrieveCurrentTeamFolders app_msg.Message
	File                                  fd_file.RowFeed
	OperationLog                          rp_model.TransactionReport
	Peer                                  rc_conn.ConnBusinessFile
	ProgressDeleteFolder                  app_msg.Message
}

func (z *Permdelete) Exec(c app_control.Control) error {
	ui := c.UI()
	if err := z.OperationLog.Open(); err != nil {
		return err
	}

	folders, err := sv_teamfolder.New(z.Peer.Context()).List()
	if err != nil {
		ui.Error(z.ErrUnableToRetrieveCurrentTeamFolders.With("Error", err.Error()))
		return err
	}

	return z.File.EachRow(func(m interface{}, rowIndex int) error {
		r := m.(*TeamFolderName)
		ui.Info(z.ProgressDeleteFolder.With("Name", r.Name))

		var folder *mo_teamfolder.TeamFolder
		for _, tf := range folders {
			if strings.ToLower(r.Name) == strings.ToLower(tf.Name) {
				folder = tf
				break
			}
		}
		if folder == nil {
			ui.Error(z.ErrTeamFolderNotFound.With("Name", r.Name))
			z.OperationLog.Failure(errors.New("team folder not found"), r)
			return nil
		}

		err := sv_teamfolder.New(z.Peer.Context()).PermDelete(folder)
		if err != nil {
			ui.Error(z.ErrUnableToDelete.With("Name", r.Name).With("Error", err.Error()))
			z.OperationLog.Failure(err, r)
			return nil
		}
		z.OperationLog.Success(r, nil)
		return nil
	})
}

func (z *Permdelete) Test(c app_control.Control) error {
	return qt_endtoend.ImplementMe()
}

func (z *Permdelete) Preset() {
	z.File.SetModel(&TeamFolderName{})
	z.OperationLog.SetModel(&TeamFolderName{}, nil)
}
