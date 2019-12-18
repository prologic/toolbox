package file

import (
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recpie/rc_kitchen"
	"github.com/watermint/toolbox/infra/recpie/rc_vo"
	"github.com/watermint/toolbox/recipe/team/namespace/file"
)

type Size struct {
}

func (z *Size) Requirement() rc_vo.ValueObject {
	return &file.SizeVO{
		IncludeSharedFolder: false,
		IncludeTeamFolder:   true,
		Depth:               2,
	}
}

func (z *Size) Exec(k rc_kitchen.Kitchen) error {
	fs := file.Size{}
	return fs.Exec(k)
}

func (z *Size) Test(c app_control.Control) error {
	fs := file.Size{}
	return fs.Test(c)
}
