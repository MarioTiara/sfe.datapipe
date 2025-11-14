package masteroutlet

import (
	"context"

	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"
)

type MasterOutletMapper interface {
	MapRowsToMasterOutlet(ctx context.Context, rows <-chan []string) ([]*masteroutlet.MasterOutlet, error)
}
