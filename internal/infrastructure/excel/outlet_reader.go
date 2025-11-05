package excel

import (
	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"
	"github.com/xuri/excelize/v2"
)

func ReadMasterOutletFromExcel(path string) ([]masteroutlet.MasterOutlet, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	var result []masteroutlet.MasterOutlet
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		m := masteroutlet.MasterOutlet{
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
