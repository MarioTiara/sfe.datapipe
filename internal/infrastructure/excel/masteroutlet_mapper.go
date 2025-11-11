package excel

import masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"

type excelMasterOutletMapper struct{}

func NewExcelMasterOutletMapper() *excelMasterOutletMapper {
	return &excelMasterOutletMapper{}
}

func (m *excelMasterOutletMapper) MapRowsToMasterOutlet(rows <-chan []string) ([]*masteroutlet.MasterOutlet, error) {

	var result []*masteroutlet.MasterOutlet

	for row := range rows {
		if len(row) == 0 {
			continue // skip empty rows
		}

		m := &masteroutlet.MasterOutlet{
			CustomerCode:      safeGet(row, 0),
			CustomerName:      safeGet(row, 1),
			Channel:           safeGet(row, 2),
			Plant:             safeGet(row, 3),
			BranchName:        safeGet(row, 4),
			NIKSalesman:       safeGet(row, 5),
			NameSalesman:      safeGet(row, 6),
			RayonCode:         safeGet(row, 7),
			Rayon:             safeGet(row, 8),
			NewClass:          safeGet(row, 9),
			CallPlanFullMonth: safeGet(row, 10),
			TargetFreq:        safeGet(row, 11),
			TerrCode:          safeGet(row, 12),
			Username:          safeGet(row, 13),
		}

		result = append(result, m)
	}

	return result, nil
}
