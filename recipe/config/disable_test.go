package config

import (
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestDisable_Exec(t *testing.T) {
	qtr_endtoend.TestRecipe(t, &Disable{})
}
