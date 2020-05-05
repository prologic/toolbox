package dbx_auth_attr

import (
	"github.com/watermint/toolbox/essentials/log/es_log"
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/ui/app_msg"
)

type MsgAttr struct {
	ErrorVerifyFailed app_msg.Message
	AttrTeamLicenses  app_msg.Message
}

var (
	MAttr = app_msg.Apply(&MsgAttr{}).(*MsgAttr)
)

type Attr struct {
	ctl  app_control.Control
	auth api_auth.Console
}

func (z *Attr) PeerName() string {
	return z.auth.PeerName()
}

func (z *Attr) Auth(scope string) (token api_auth.Context, err error) {
	l := z.ctl.Log().With(es_log.String("peerName", z.PeerName()), es_log.String("scope", scope))
	ui := z.ctl.UI()

	tc, err := z.auth.Auth(scope)
	if err != nil {
		return nil, err
	}

	l.Debug("Start verify token")

	tc, err = VerifyToken(tc, z.ctl)
	if err != nil {
		l.Debug("failed verify token", es_log.Error(err))
		ui.Error(MAttr.ErrorVerifyFailed.With("Error", err))
		return nil, err
	}
	return tc, nil
}
