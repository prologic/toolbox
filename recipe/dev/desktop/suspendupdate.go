package desktop

import (
	"github.com/watermint/toolbox/essentials/log/es_log"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/quality/infra/qt_endtoend"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"os"
	"path/filepath"
	"runtime"
)

type Suspendupdate struct {
	rc_recipe.RemarkSecret
	UpdaterPath string
	UpdaterName string
	Unsuspend   bool
}

func (z *Suspendupdate) Exec(c app_control.Control) error {
	l := c.Log()
	if runtime.GOOS != "windows" {
		l.Info("Skip: operation is not supported on this platform")
		return nil
	}
	mode := "Suspend"
	oldName := z.UpdaterName
	newName := "_" + z.UpdaterName
	if z.Unsuspend {
		mode = "Unsuspend"
		oldName = "_" + z.UpdaterName
		newName = z.UpdaterName
	}
	l = l.With(es_log.String("mode", mode))

	path := filepath.Join(z.UpdaterPath, oldName)
	ls, err := os.Lstat(path)
	if err != nil {
		l.Info("Unable to locate Updater", es_log.Error(err), es_log.String("path", path))
		return err
	}
	l.Debug("Updater", es_log.Any("lstat", ls))

	l.Info("Trying to rename Updater", es_log.String("path", path), es_log.String("newName", newName))
	if err = os.Rename(path, filepath.Join(z.UpdaterPath, newName)); err != nil {
		l.Error("Unable to rename Updater", es_log.Error(err))
		return err
	}
	l.Info("Updater")
	return nil
}

func (z *Suspendupdate) Test(c app_control.Control) error {
	if qt_endtoend.IsSkipEndToEndTest() {
		return nil
	}
	return qt_errors.ErrorNoTestRequired
}

func (z *Suspendupdate) Preset() {
	z.UpdaterPath = "C:/Program Files (x86)/Dropbox/Update"
	z.UpdaterName = "DropboxUpdate.exe"
}
