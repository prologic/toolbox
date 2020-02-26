package rp_model_impl

import (
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_opt"
	"github.com/watermint/toolbox/infra/report/rp_model"
)

func newCascade(name string, ctl app_control.Control) Writer {
	fileWriters := make([]Writer, 0)
	consoleWriters := make([]Writer, 0)

	fileWriters = append(fileWriters, NewJsonWriter(name, ctl, false))
	if !ctl.IsLowMemory() {
		fileWriters = append(fileWriters, newCsvWriter(name, ctl))
		fileWriters = append(fileWriters, NewXlsxWriter(name, ctl))
	}
	if !ctl.IsQuiet() {
		if ctl.UIFormat() == app_opt.OutputJson {
			consoleWriters = append(consoleWriters, NewJsonWriter(name, ctl, true))
		} else {
			consoleWriters = append(consoleWriters, newUIWriter(name, ctl))
		}
	}

	return &cascadeWriter{
		ctl:            ctl,
		name:           name,
		fileWriters:    fileWriters,
		consoleWriters: consoleWriters,
	}
}

type cascadeWriter struct {
	ctl            app_control.Control
	name           string
	writers        []Writer
	fileWriters    []Writer
	consoleWriters []Writer
	isClosed       bool
}

func (z *cascadeWriter) Name() string {
	return z.name
}

func (z *cascadeWriter) Open(ctl app_control.Control, model interface{}, opts ...rp_model.ReportOpt) error {
	ro := &rp_model.ReportOpts{}
	for _, o := range opts {
		o(ro)
	}

	z.writers = make([]Writer, 0)
	z.writers = append(z.writers, z.fileWriters...)
	if !ro.NoConsoleOutput {
		z.writers = append(z.writers, z.consoleWriters...)
	}

	for _, w := range z.writers {
		if err := w.Open(ctl, model, opts...); err != nil {
			z.Close()
			return err
		}
	}
	return nil
}

func (z *cascadeWriter) Row(r interface{}) {
	if z.isClosed {
		return
	}

	for _, w := range z.writers {
		w.Row(r)
	}
}

func (z *cascadeWriter) Close() {
	for _, w := range z.writers {
		w.Close()
	}

	p := z.ctl.Workspace().Report()
	ui := z.ctl.UI()
	ui.OpenArtifact(p, z.ctl.IsAutoOpen())
	z.isClosed = true
}
