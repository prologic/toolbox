package sv_profile

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

// Mock tests

func TestProfileImpl_Current(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := NewProfile(ctx)
		_, err := sv.Current()
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestTeamImpl_Admin(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := NewTeam(ctx)
		_, err := sv.Admin()
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}
