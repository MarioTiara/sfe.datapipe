package excel

import materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"

type excelMaterialMasterMapper struct{}

func NewExcelMaterialMasterMapper() *excelMaterialMasterMapper {
	return &excelMaterialMasterMapper{}
}

func (m *excelMaterialMasterMapper) MapRowsToMaterialMaster(rows <-chan []string) (<-chan *materialmaster.MaterialMaster, <-chan error) {
	outCh := make(chan *materialmaster.MaterialMaster)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)

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

			outCh <- m
		}
	}()

	return outCh, errCh
}
