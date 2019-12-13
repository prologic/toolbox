package sv_file_relocation

import (
	"fmt"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/infra/api/api_test"
	"testing"
	"time"
)

func TestCopy(t *testing.T) {
	api_test.DoTestTokenFull(func(ctx api_context.Context) {
		r := New(ctx)
		src := api_test.ToolboxTestSuiteFolder.ChildPath("copy/F0.jpg")
		name := fmt.Sprintf("copy-%x.jpg", time.Now().Unix())
		dest := api_test.ToolboxTestSuiteFolder.ChildPath("copy", name)

		entry, err := r.Copy(src, dest)
		if err != nil {
			t.Error(err)
		}
		if entry.Name() != name {
			t.Error("invalid")
		}
	})
}
