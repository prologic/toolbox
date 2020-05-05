package file

import (
	mo_path2 "github.com/watermint/toolbox/domain/common/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_file"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_file_content"
	"github.com/watermint/toolbox/essentials/log/es_log"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
	"os"
	"path/filepath"
)

type Download struct {
	rc_recipe.RemarkExperimental
	Peer         dbx_conn.ConnUserFile
	DropboxPath  mo_path.DropboxPath
	LocalPath    mo_path2.FileSystemPath
	OperationLog rp_model.RowReport
}

func (z *Download) Preset() {
	z.OperationLog.SetModel(
		&mo_file.ConcreteEntry{},
		rp_model.HiddenColumns(
			"id",
			"path_lower",
			"revision",
			"content_hash",
			"shared_folder_id",
			"parent_shared_folder_id",
		),
	)
}

func (z *Download) Exec(c app_control.Control) error {
	l := c.Log()
	ctx := z.Peer.Context()

	if err := z.OperationLog.Open(); err != nil {
		return err
	}

	entry, f, err := sv_file_content.NewDownload(ctx).Download(z.DropboxPath)
	if err != nil {
		return err
	}
	if err := os.Rename(f.Path(), filepath.Join(z.LocalPath.Path(), entry.Name())); err != nil {
		l.Debug("Unable to move file to specified path",
			es_log.Error(err),
			es_log.String("downloaded", f.Path()),
			es_log.String("destination", z.LocalPath.Path()),
		)
		return err
	}

	z.OperationLog.Row(entry.Concrete())
	return nil
}

func (z *Download) Test(c app_control.Control) error {
	return rc_exec.ExecMock(c, &Download{}, func(r rc_recipe.Recipe) {
		m := r.(*Download)
		m.LocalPath = qt_recipe.NewTestFileSystemFolderPath(c, "download")
		m.DropboxPath = qt_recipe.NewTestDropboxFolderPath("file-download")
	})
}
