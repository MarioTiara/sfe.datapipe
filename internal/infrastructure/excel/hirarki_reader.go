package excel

import (
	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
	"github.com/xuri/excelize/v2"
)

func ReadHirarkiFromExcel(path string) ([]hirarki.Hirarki, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	var result []hirarki.Hirarki
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		h := hirarki.Hirarki{
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

	return result, nil
}
