package job

import (
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestLoop_Exec(t *testing.T) {
	qtr_endtoend.TestRecipe(t, &Loop{})
}
