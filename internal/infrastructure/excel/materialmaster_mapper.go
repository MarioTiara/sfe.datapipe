package excel

import (
	"context"
	"fmt"

	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type excelMaterialMasterMapper struct {
	logger logger.Logger
}

func NewExcelMaterialMasterMapper(logger logger.Logger) *excelMaterialMasterMapper {
	return &excelMaterialMasterMapper{logger: logger}
}

func (m *excelMaterialMasterMapper) MapRowsToMaterialMaster(ctx context.Context, rows <-chan []string) ([]*materialmaster.MaterialMaster, error) {
	var result []*materialmaster.MaterialMaster
	for row := range rows {
		if len(row) == 0 {
			continue // skip empty rows
		}
		m := &materialmaster.MaterialMaster{
			PrincipalCode:       safeGet(row, 0),
			PrincipalName:       safeGet(row, 1),
			MaterialCode:        safeGet(row, 2),
			MaterialName:        safeGet(row, 3),
			Uom:                 safeGet(row, 4),
			Division:            safeGet(row, 5),
			DivisionDescription: safeGet(row, 6),
			ProductHierarchy:    safeGet(row, 7),
			Brand:               safeGet(row, 8),
			MATKL:               safeGet(row, 9),
			PH2Code:             safeGet(row, 10),
			PH2Name:             safeGet(row, 11),
			PH3Code:             safeGet(row, 12),
			PH3Name:             safeGet(row, 13),
			ATCCode:             safeGet(row, 14),
			ATCDescription:      safeGet(row, 15),
			DGIndicatorCode:     safeGet(row, 16),
			DGIndicatorName:     safeGet(row, 17),
			GeneralDGCode:       safeGet(row, 18),
			GeneralDGName:       safeGet(row, 19),
		}

		result = append(result, m)
	}
	m.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(result)))
	return result, nil
}
