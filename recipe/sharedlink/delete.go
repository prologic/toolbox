package sharedlink

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_sharedlink"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_sharedlink"
	"github.com/watermint/toolbox/essentials/log/es_log"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
	"path/filepath"
	"strings"
)

type Delete struct {
	rc_recipe.RemarkIrreversible
	Peer                 dbx_conn.ConnUserFile
	Path                 mo_path.DropboxPath
	Recursive            bool
	SharedLink           rp_model.TransactionReport
	InfoNoLinksAtThePath app_msg.Message
	ProgressDelete       app_msg.Message
}

func (z *Delete) Preset() {
	z.SharedLink.SetModel(
		&mo_sharedlink.Metadata{},
		nil,
		rp_model.HiddenColumns(
			"input.id",
		),
	)
}

func (z *Delete) Exec(c app_control.Control) error {
	if err := z.SharedLink.Open(); err != nil {
		return err
	}

	if z.Recursive {
		return z.removeRecursive(c)
	} else {
		return z.removePathAt(c)
	}
}

func (z *Delete) removePathAt(c app_control.Control) error {
	ui := c.UI()
	l := c.Log()
	links, err := sv_sharedlink.New(z.Peer.Context()).ListByPath(z.Path)
	if err != nil {
		return err
	}
	if len(links) < 1 {
		ui.Info(z.InfoNoLinksAtThePath.With("Path", z.Path.Path()))
		return nil
	}

	var lastErr error
	for _, link := range links {
		ui.Progress(z.ProgressDelete.With("Url", link.LinkUrl()).With("Path", link.LinkPathLower()))
		err = sv_sharedlink.New(z.Peer.Context()).Remove(link)
		if err != nil {
			l.Debug("Unable to remove link", es_log.Error(err), es_log.Any("link", link))
			z.SharedLink.Failure(err, link)
			lastErr = err
		} else {
			z.SharedLink.Success(link, nil)
		}
	}
	return lastErr
}

func (z *Delete) removeRecursive(c app_control.Control) error {
	ui := c.UI()
	l := c.Log().With(es_log.String("path", z.Path.Path()))
	links, err := sv_sharedlink.New(z.Peer.Context()).List()
	if err != nil {
		return err
	}
	if len(links) < 1 {
		ui.Info(z.InfoNoLinksAtThePath.With("Path", z.Path.Path()))
		return nil
	}

	var lastErr error
	for _, link := range links {
		l = l.With(es_log.String("linkPath", link.LinkPathLower()))
		rel, err := filepath.Rel(strings.ToLower(z.Path.Path()), link.LinkPathLower())
		if err != nil {
			l.Debug("Skip due to path calc error", es_log.Error(err))
			continue
		}
		if strings.HasPrefix(rel, "..") {
			l.Debug("Skip due to non related path")
			continue
		}

		ui.Progress(z.ProgressDelete.With("Url", link.LinkUrl()).With("Path", link.LinkPathLower()))
		err = sv_sharedlink.New(z.Peer.Context()).Remove(link)
		if err != nil {
			l.Debug("Unable to remove link", es_log.Error(err), es_log.Any("link", link))
			z.SharedLink.Failure(err, link)
			lastErr = err
		} else {
			z.SharedLink.Success(link, nil)
		}
	}
	return lastErr
}

func (z *Delete) Test(c app_control.Control) error {
	// Non-recursive
	err := rc_exec.ExecMock(c, &Delete{}, func(r rc_recipe.Recipe) {
		m := r.(*Delete)
		m.Path = qt_recipe.NewTestDropboxFolderPath("sharedlink-delete")
		m.Recursive = false
	})
	if e, _ := qt_errors.ErrorsForTest(c.Log(), err); e != nil {
		return e
	}

	// Recursive
	err = rc_exec.ExecMock(c, &Delete{}, func(r rc_recipe.Recipe) {
		m := r.(*Delete)
		m.Path = qt_recipe.NewTestDropboxFolderPath("sharedlink-delete")
		m.Recursive = true
	})
	if e, _ := qt_errors.ErrorsForTest(c.Log(), err); e != nil {
		return e
	}
	return nil
}
