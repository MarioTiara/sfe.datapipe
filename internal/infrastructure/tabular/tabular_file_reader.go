package tabular

import (
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/mariotiara/sfe-data-pipe/internal/application/ports"
	"github.com/xuri/excelize/v2"
)

type TabularFileReader struct {
	FileSource ports.FileSource
}

func New(fs ports.FileSource) *TabularFileReader {
	return &TabularFileReader{FileSource: fs}
}

func (t *TabularFileReader) ReadRows(name string) (<-chan []string, error) {
	rc, err := t.FileSource.Open(name)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(name))

	switch ext {
	case ".csv":
		return readCSV(rc), nil
	case ".xlsx":
		return readExcel(rc)
	default:
		rc.Close()
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

func readCSV(r io.ReadCloser) <-chan []string {
	out := make(chan []string)

	go func() {
		defer close(out)
		defer r.Close()

		reader := csv.NewReader(r)
		for {
			row, err := reader.Read()
			if err == io.EOF {
				return
			}
			if err != nil {
				continue
			}
			out <- row
		}
	}()

	return out
}

func readExcel(r io.ReadCloser) (<-chan []string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		r.Close()
		return nil, err
	}

	sheet := f.GetSheetName(0)
	out := make(chan []string)

	go func() {
		defer close(out)
		defer r.Close()
		defer f.Close()

		rows, err := f.Rows(sheet)
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			cols, err := rows.Columns()
			if err != nil {
				continue
			}
			out <- cols
		}
	}()

	return out, nil
}
