package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelParser struct{}

func (e *ExcelParser) Parse(filePath string) (<-chan []string, <-chan error) {
	dataCh := make(chan []string)
	errCh := make(chan error, 1)

	fmt.Println("start parsing excel")
	go func() {
		defer close(dataCh)
		defer close(errCh)

		f, err := excelize.OpenFile(filePath)
		if err != nil {
			errCh <- err
			return
		}

		rows, err := f.GetRows("Sheet1")
		if err != nil {
			errCh <- err
			return
		}

		for _, row := range rows {
			dataCh <- row
		}
	}()

	fmt.Println("Parsing excel Done")
	return dataCh, errCh
}
