package sv_file_revision

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestRevisionImpl_List(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := New(ctx)
		_, err := sv.List(qtr_endtoend.NewTestDropboxFolderPath())
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestRevisionImpl_ListById(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := New(ctx)
		_, err := sv.ListById(qtr_endtoend.NewTestDropboxFolderPath())
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}
