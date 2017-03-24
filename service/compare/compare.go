package compare

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/watermint/toolbox/infra"
	"github.com/watermint/toolbox/service/report"
	"io"
	"os"
	"sync"
)

const (
	BLOCK_SIZE     = 4194304
	HASH_FOR_EMPTY = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

// Calculate File content hash
func ContentHash(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		seelog.Warnf("Unable to acquire information about path [%s] error[%s]", path, err)
		return "", err
	}
	if info.Size() < 1 {
		return HASH_FOR_EMPTY, nil
	}

	f, err := os.Open(path)
	if err != nil {
		seelog.Warnf("Unable to open file [%s] error[%s]", path, err)
		return "", err
	}
	defer f.Close()

	var loadedBytes, totalBytes int64

	loadedBytes = 0
	totalBytes = info.Size()
	hashPerBlock := make([][32]byte, 0)

	for (totalBytes - loadedBytes) > 0 {
		r := io.LimitReader(f, BLOCK_SIZE)
		block := make([]byte, BLOCK_SIZE)
		readBytes, err := r.Read(block)
		if err == io.EOF {
			break
		}
		if err != nil {
			seelog.Warnf("Unable to load file [%s] error[%s]", path, err)
			return "", err
		}

		h := sha256.Sum256(block[:readBytes])
		hashPerBlock = append(hashPerBlock, h)
	}

	concatenated := make([]byte, 0)
	for _, h := range hashPerBlock {
		concatenated = append(concatenated, h[:]...)
	}
	h := sha256.Sum256(concatenated)
	return hex.EncodeToString(h[:]), nil
}

type CompareOpts struct {
	InfraOpts  *infra.InfraOpts
	ReportOpts *report.MultiReportOpts

	DropboxToken    string
	LocalBasePath   string
	DropboxBasePath string
}

// Returns (true, nil) when equal contents between local and dropbox
func Compare(opts *CompareOpts) (bool, error) {
	var err error
	trav := Traverse{
		DropboxToken:    opts.DropboxToken,
		DropboxBasePath: opts.DropboxBasePath,
		LocalBasePath:   opts.LocalBasePath,
		InfraOpts:       opts.InfraOpts,
	}
	err = trav.Prepare()
	if err != nil {
		return false, err
	}
	defer trav.Close()
	err = opts.ReportOpts.BeginMultiReport()
	if err != nil {
		seelog.Warnf("Unable to prepare report : error[%s]", err)
		return false, err
	}
	defer opts.ReportOpts.EndMultiReport()

	err = compareScan(&trav)
	if err != nil {
		seelog.Warnf("Unable to scan: error[%s]", err)
		return false, err
	}

	matchSum, _ := reportSummary(&trav, opts)
	matchDtl, _ := reportDropboxToLocal(&trav, opts)
	matchLtd, _ := reportLocalToDropbox(&trav, opts)
	matchSah, _ := reportSizeAndHash(&trav, opts)

	opts.ReportOpts.EndMultiReport()

	return matchSum && matchDtl && matchLtd && matchSah, nil
}

func compareScan(trav *Traverse) error {
	seelog.Info("Start scanning local files")
	err := trav.ScanLocal()
	if err != nil {
		seelog.Warnf("Unable to scan local files : error[%s]", err)
		return err
	}

	seelog.Info("Start scanning dropbox files")
	err = trav.ScanDropbox()
	if err != nil {
		seelog.Warnf("Unable to scan Dropbox files : error[%s]", err)
		return err
	}
	return nil
}

func reportSummary(trav *Traverse, opts *CompareOpts) (bool, error) {
	seelog.Debug("Start reporting Summary of comparison")
	dbxCount, dbxSize, err := trav.SummaryDropbox()
	if err != nil {
		seelog.Warnf("Unable to summarise results : error[%s]", err)
		return false, err
	}
	loCount, loSize, err := trav.SummaryLocal()
	if err != nil {
		seelog.Warnf("Unable to summarise results : error[%s]", err)
		return false, err
	}
	repo := make(chan report.ReportRow)
	wg := &sync.WaitGroup{}

	go opts.ReportOpts.Write("Summary", repo, wg)

	repo <- report.ReportHeader{
		Headers: []string{
			"Source",
			"File count",
			"Total size (bytes)",
		},
	}

	repo <- report.ReportData{
		Data: []interface{}{
			fmt.Sprintf("Local(%s)", opts.LocalBasePath),
			loCount,
			loSize,
		},
	}

	repo <- report.ReportData{
		Data: []interface{}{
			"Dropbox",
			dbxCount,
			dbxSize,
		},
	}

	repo <- nil

	wg.Wait()

	return dbxCount == loCount && dbxSize == loSize, nil
}

func reportDropboxToLocal(trav *Traverse, opts *CompareOpts) (bool, error) {
	seelog.Debug("Start reporting Dropbox to Local")
	wg := &sync.WaitGroup{}
	cmpRows := make(chan *CompareRowDropboxToLocal)
	repRows := make(chan report.ReportRow)

	go trav.CompareDropboxToLocal(cmpRows, wg)
	go opts.ReportOpts.Write("NotFoundInLocal", repRows, wg)

	repRows <- report.ReportHeader{
		Headers: []string{
			"Path",
			"File Size (bytes)",
			"Content Hash",
			"Dropbox File ID",
			"Dropbox Revision",
		},
	}

	seelog.Debug("*** Record: files not found in Local")
	rowCount := 0
	for row := range cmpRows {
		if row == nil {
			break
		}
		rowCount++

		seelog.Debugf("Path[%s] (lower:%s) Size[%d] Hash[%s] DropboxFileId[%s] DropboxRev[%s]\n",
			row.Path,
			row.PathLower,
			row.Size,
			row.ContentHash,
			row.DropboxFileId,
			row.DropboxRevision,
		)

		repRows <- report.ReportData{
			Data: []interface{}{
				row.Path,
				row.Size,
				row.ContentHash,
				row.DropboxFileId,
				row.DropboxRevision,
			},
		}
	}
	repRows <- nil

	wg.Wait()

	return rowCount == 0, nil
}

func reportLocalToDropbox(trav *Traverse, opts *CompareOpts) (bool, error) {
	seelog.Debug("Start reporting Local to Dropbox comparison")
	wg := &sync.WaitGroup{}
	cmpRows := make(chan *CompareRowLocalToDropbox)
	repRows := make(chan report.ReportRow)

	go trav.CompareLocalToDropbox(cmpRows, wg)
	go opts.ReportOpts.Write("NotFoundInDropbox", repRows, wg)

	repRows <- report.ReportHeader{
		Headers: []string{
			"Path",
			"File Size (bytes)",
			"Content Hash",
		},
	}

	seelog.Debug("*** Record: files not found in Dropbox")
	rowCount := 0
	for row := range cmpRows {
		if row == nil {
			break
		}
		rowCount++
		seelog.Debugf("Path[%s] (lower:%s) Size[%d] Hash[%s]\n",
			row.Path,
			row.PathLower,
			row.Size,
			row.ContentHash,
		)

		repRows <- report.ReportData{
			Data: []interface{}{
				row.Path,
				row.Size,
				row.ContentHash,
			},
		}
	}
	repRows <- nil

	wg.Wait()

	return rowCount == 0, nil
}

func reportSizeAndHash(trav *Traverse, opts *CompareOpts) (bool, error) {
	seelog.Debug("Start reporting size and hash comparison")
	wg := &sync.WaitGroup{}
	cmpRows := make(chan *CompareRowSizeAndHash)
	repRows := make(chan report.ReportRow)

	go opts.ReportOpts.Write("DifferentContent", repRows, wg)
	go trav.CompareSizeAndHash(cmpRows, wg)

	repRows <- report.ReportHeader{
		Headers: []string{
			"Path",
			"Local File Size (bytes)",
			"Dropbox File Size (bytes)",
			"Local Content Hash",
			"Dropbox Content Hash",
		},
	}

	seelog.Debug("*** Record: files size and/or hash not mached")
	rowCount := 0
	for row := range cmpRows {
		if row == nil {
			break
		}
		rowCount++

		seelog.Debugf("Path[%s] (lower:%s) Size(Local:%d, Dropbox:%d), Hash(Local:%s, Dropbox:%s)\n",
			row.Path,
			row.PathLower,
			row.LocalSize,
			row.DropboxSize,
			row.LocalContentHash,
			row.DropboxContentHash,
		)

		repRows <- report.ReportData{
			Data: []interface{}{
				row.Path,
				row.LocalSize,
				row.DropboxSize,
				row.LocalContentHash,
				row.DropboxContentHash,
			},
		}
	}
	repRows <- nil

	wg.Wait()

	return rowCount == 0, nil
}
