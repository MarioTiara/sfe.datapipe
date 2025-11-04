package excel

import (
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/salesfe"
	"github.com/xuri/excelize/v2"
)

func ReadSalesFromExcel(path string) ([]salesfe.SalesFE, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName((0)))
	if err != nil {
		return nil, err
	}

	var result []salesfe.SalesFE
	for i, row := range rows {
		if i == 0 {
			continue
		}

		date, _ := time.Parse("20060102", row[0])

		result = append(result, salesfe.SalesFE{
			InvoiceDate:           date,
			PONumber:              row[1],
			POType:                row[2],
			POTypeDesc:            row[3],
			CustomerCode:          row[4],
			CustomerName:          row[5],
			ChannelIC1:            row[6],
			ChannelIC1Description: row[7],
			ChannelIC4:            row[8],
			ChannelIC4Description: row[9],
			Plant:                 row[10],
			Branch:                row[11],
			Principal:             row[12],
			ProductGroup:          row[13],
			ItemCode:              row[14],
			ItemName:              row[15],
		})
	}

	return result, nil
}
