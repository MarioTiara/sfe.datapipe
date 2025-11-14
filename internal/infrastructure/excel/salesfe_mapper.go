package excel

import (
	"context"
	"fmt"
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/salesfe"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type excelSalesFEMapper struct {
	logger logger.Logger
}

func NewExcelSalesFEMapper(logger logger.Logger) *excelSalesFEMapper {
	return &excelSalesFEMapper{logger: logger}
}

func (m *excelSalesFEMapper) MapRowsToSaleFE(ctx context.Context, rows <-chan []string) ([]*salesfe.SalesFE, error) {
	var result []*salesfe.SalesFE

	for row := range rows {
		if len(row) == 0 {
			continue // skip empty rows
		}

		// Parse InvoiceDate if exists and in valid format
		var invoiceDate time.Time
		if len(row) > 0 {
			if parsed, err := time.Parse("20060102", safeGet(row, 0)); err == nil {
				invoiceDate = parsed
			}
		}

		s := &salesfe.SalesFE{
			InvoiceDate:           invoiceDate,
			PONumber:              safeGet(row, 1),
			POType:                safeGet(row, 2),
			POTypeDesc:            safeGet(row, 3),
			CustomerCode:          safeGet(row, 4),
			CustomerName:          safeGet(row, 5),
			ChannelIC1:            safeGet(row, 6),
			ChannelIC1Description: safeGet(row, 7),
			ChannelIC4:            safeGet(row, 8),
			ChannelIC4Description: safeGet(row, 9),
			Plant:                 safeGet(row, 10),
			Branch:                safeGet(row, 11),
			Principal:             safeGet(row, 12),
			ProductGroup:          safeGet(row, 13),
			ItemCode:              safeGet(row, 14),
			ItemName:              safeGet(row, 15),
			Sales:                 parseFloat(safeGet(row, 16)),
			NetSales:              parseFloat(safeGet(row, 17)),
			SalesUnit:             parseFloat(safeGet(row, 18)),
			BonusUnit:             parseFloat(safeGet(row, 19)),
			BUN1:                  safeGet(row, 20),
		}

		result = append(result, s)
	}

	m.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(result)))
	return result, nil
}
