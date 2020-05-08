package sv_file

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_file"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestFilesImpl_ListChunked(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := NewFiles(ctx)
		err := sv.ListChunked(qtr_endtoend.NewTestDropboxFolderPath(), func(entry mo_file.Entry) {},
			Recursive(),
			IncludeMediaInfo(),
			IncludeDeleted(),
			IncludeHasExplicitSharedMembers(),
		)
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestFilesImpl_List(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := NewFiles(ctx)
		_, err := sv.List(qtr_endtoend.NewTestDropboxFolderPath(),
			Recursive(),
			IncludeMediaInfo(),
			IncludeDeleted(),
			IncludeHasExplicitSharedMembers(),
		)
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestFilesImpl_Poll(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := NewFiles(ctx)
		err := sv.Poll(qtr_endtoend.NewTestDropboxFolderPath(), func(entry mo_file.Entry) {
		})
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestFilesImpl_Remove(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := NewFiles(ctx)
		_, err := sv.Remove(qtr_endtoend.NewTestDropboxFolderPath(), RemoveRevision("test"))
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestFilesImpl_Resolve(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := NewFiles(ctx)
		_, err := sv.Resolve(qtr_endtoend.NewTestDropboxFolderPath())
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestFilesImpl_Search(t *testing.T) {
	qtr_endtoend.TestWithDbxContext(t, func(ctx dbx_context.Context) {
		sv := NewFiles(ctx)
		_, err := sv.Search("test",
			SearchPath(qtr_endtoend.NewTestDropboxFolderPath()),
			SearchMaxResults(100),
			SearchFileDeleted(),
			SearchFileNameOnly(),
			SearchFileExtension("test"),
			SearchCategories("pdf"),
			SearchIncludeHighlights(),
		)
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}
