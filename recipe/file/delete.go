package file

import (
	"github.com/watermint/toolbox/domain/model/mo_path"
	"github.com/watermint/toolbox/domain/service/sv_file"
	"github.com/watermint/toolbox/infra/api/api_util"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recpie/app_conn"
	"github.com/watermint/toolbox/infra/recpie/app_kitchen"
	"github.com/watermint/toolbox/infra/recpie/app_vo"
	"github.com/watermint/toolbox/infra/report/rp_spec"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
)

type DeleteVO struct {
	Peer app_conn.ConnUserFile
	Path string
}

type Delete struct {
}

func (z *Delete) Console() {
}

func (z *Delete) Requirement() app_vo.ValueObject {
	return &DeleteVO{}
}

func (z *Delete) Exec(k app_kitchen.Kitchen) error {
	vo := k.Value().(*DeleteVO)
	ui := k.UI()
	ctx, err := vo.Peer.Connect(k.Control())
	if err != nil {
		return err
	}

	var delete func(path mo_path.Path) error
	delete = func(path mo_path.Path) error {
		ui.Info("recipe.file.delete.progress.deleting", app_msg.P{
			"Path": path.Path(),
		})
		_, err = sv_file.NewFiles(ctx).Remove(path)
		if err == nil {
			return nil
		}
		switch api_util.ErrorSummary(err) {
		case "too_many_files":
			entries, err := sv_file.NewFiles(ctx).List(path)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				if f, ok := entry.File(); ok {
					delete(f.Path())
				}
				if f, ok := entry.Folder(); ok {
					delete(f.Path())
				}
			}
			return delete(path)

		default:
			return err
		}
	}

	return delete(mo_path.NewPath(vo.Path))
}

func (z *Delete) Test(c app_control.Control) error {
	return qt_recipe.ScenarioTest()
}

func (z *Delete) Reports() []rp_spec.ReportSpec {
	return []rp_spec.ReportSpec{}
}
