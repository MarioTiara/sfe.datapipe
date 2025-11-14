package materialmaster

import (
	"context"

	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"
)

type MaterialMasterMapper interface {
	MapRowsToMaterialMaster(ctx context.Context, rows <-chan []string) ([]*materialmaster.MaterialMaster, error)
}
