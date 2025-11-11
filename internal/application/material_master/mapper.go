package materialmaster

import (
	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"
)

type MaterialMasterMapper interface {
	MapRowsToMaterialMaster(rows <-chan []string) ([]*materialmaster.MaterialMaster, error)
}
