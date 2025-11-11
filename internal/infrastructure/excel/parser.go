package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelParser struct{}

func (e *ExcelParser) Parse(filePath string) (<-chan []string, <-chan error) {
	dataCh := make(chan []string)
	errCh := make(chan error, 1)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		f, err := excelize.OpenFile(filePath)
		if err != nil {
			errCh <- fmt.Errorf("open file: %w", err)
			return
		}

		// Get all sheet names
		sheetList := f.GetSheetList()
		if len(sheetList) == 0 {
			errCh <- fmt.Errorf("no sheets found")
		}

		// Use the first sheet
		firstSheet := sheetList[0]

		// Get all rows from the first sheet
		rows, err := f.GetRows(firstSheet)
		if err != nil {
			errCh <- fmt.Errorf("get rows: %w", err)
			return
		}

		for _, row := range rows {
			dataCh <- row
		}
	}()

	return dataCh, errCh
}
