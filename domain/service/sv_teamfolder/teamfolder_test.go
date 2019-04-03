package sv_teamfolder

import (
	"github.com/watermint/toolbox/domain/infra/api_context"
	"github.com/watermint/toolbox/domain/infra/api_test"
	"testing"
)

func TestTeamFolderImpl_List(t *testing.T) {
	api_test.DoTestBusinessFile(func(ctx api_context.Context) {
		svc := New(ctx)
		list, err := svc.List()
		if err != nil {
			t.Error(err)
			return
		}

		for _, tf := range list {
			if tf.TeamFolderId == "" {
				t.Error("invalid")
			}
			r, err := svc.Resolve(tf.TeamFolderId)
			if err != nil {
				t.Error(err)
			}
			if r.TeamFolderId != tf.TeamFolderId || r.Name != tf.Name {
				t.Error("invalid")
			}
		}
	})
}
