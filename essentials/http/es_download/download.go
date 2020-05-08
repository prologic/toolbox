package es_download

import (
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/network/nw_bandwidth"
	"io"
	"net/http"
	"os"
)

func Download(l esl.Logger, url string, path string) error {
	l.Debug("Try download", esl.String("url", url))
	resp, err := http.Get(url)
	if err != nil {
		l.Debug("Unable to create download request")
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		l.Debug("Unable to create download file")
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, nw_bandwidth.WrapReader(resp.Body))
	if err != nil {
		l.Debug("Unable to copy from response", esl.Error(err))
		return err
	}
	l.Debug("Finished download", esl.String("url", url))
	return nil
}
