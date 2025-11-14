package excel

import (
	"context"
	"fmt"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type excelHirarkiMapper struct {
	logger logger.Logger
}

func NewHirarkiExcelMapper(logger logger.Logger) *excelHirarkiMapper {
	return &excelHirarkiMapper{}
}

func (m *excelHirarkiMapper) MapRowsToHirarki(ctx context.Context, rows <-chan []string) ([]*hirarki.Hirarki, error) {
	var result []*hirarki.Hirarki
	m.logger.Info(ctx, "mapping streams to entities")
	for row := range rows {
		if len(row) == 0 {
			continue
		}

		h := &hirarki.Hirarki{
			RayonCode:              safeGet(row, 0),
			Plant:                  safeGet(row, 1),
			RayonType:              safeGet(row, 2),
			BUM:                    safeGet(row, 3),
			BUMName:                safeGet(row, 4),
			NSM:                    safeGet(row, 5),
			NSMName:                safeGet(row, 6),
			ASM:                    safeGet(row, 7),
			ASMName:                safeGet(row, 8),
			FSS:                    safeGet(row, 9),
			FSSName:                safeGet(row, 10),
			SLM:                    safeGet(row, 11),
			SLMName:                safeGet(row, 12),
			SalesmanCategoryUpdate: safeGet(row, 13),
			BranchName:             safeGet(row, 14),
			MLO:                    safeGet(row, 15),
			Remarks:                safeGet(row, 16),
			TerrCode:               safeGet(row, 17),
			Username:               safeGet(row, 18),
			Change:                 safeGet(row, 19),
			CategoryRayon:          safeGet(row, 20),
			RayonDetail:            safeGet(row, 21),
		}

		result = append(result, h)
	}

	m.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(result)))
	return result, nil
}
