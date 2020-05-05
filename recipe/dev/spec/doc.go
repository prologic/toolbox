package spec

import (
	"encoding/json"
	"github.com/watermint/toolbox/domain/common/model/mo_string"
	"github.com/watermint/toolbox/essentials/io/es_stdout"
	"github.com/watermint/toolbox/essentials/log/es_log"
	"github.com/watermint/toolbox/infra/control/app_catalogue"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_doc"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/recipe/rc_spec"
	"io"
	"os"
)

type Doc struct {
	rc_recipe.RemarkSecret
	Lang     mo_string.OptionalString
	FilePath mo_string.OptionalString
}

func (z *Doc) Preset() {
}

func (z *Doc) traverseCatalogue(c app_control.Control) error {
	l := c.Log()
	sd := make(map[string]*rc_doc.Recipe)
	cat := app_catalogue.Current()

	for _, r := range cat.Recipes() {
		s := rc_spec.New(r)

		l.Debug("Generating", es_log.String("recipe", s.CliPath()))
		d := s.Doc(c.UI())
		sd[d.Path] = d
	}

	var w io.WriteCloser
	var err error
	shouldClose := false
	if !z.FilePath.IsExists() {
		w = es_stdout.NewDefaultOut(c.Feature().IsTest())
	} else {
		w, err = os.Create(z.FilePath.Value())
		if err != nil {
			l.Error("Unable to create spec file", es_log.Error(err), es_log.String("path", z.FilePath.Value()))
			return err
		}
		shouldClose = true
	}
	defer func() {
		if shouldClose {
			w.Close()
		}
	}()

	je := json.NewEncoder(w)
	je.SetIndent("", "  ")
	je.SetEscapeHTML(false)
	if err := je.Encode(sd); err != nil {
		l.Error("Unable to generate spec doc", es_log.Error(err))
		return err
	}
	return nil
}

func (z *Doc) Exec(c app_control.Control) error {
	if z.Lang.IsExists() {
		return z.traverseCatalogue(c.WithLang(z.Lang.Value()))
	} else {
		return z.traverseCatalogue(c)
	}
}

func (z *Doc) Test(c app_control.Control) error {
	return rc_exec.Exec(c, z, rc_recipe.NoCustomValues)
}
