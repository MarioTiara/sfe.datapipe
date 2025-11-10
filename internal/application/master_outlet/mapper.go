package masteroutlet

import (
	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"
)

type MasterOutletMapper interface {
	MapRowsToMasterOutlet(rows <-chan []string) (<-chan *masteroutlet.MasterOutlet, <-chan error)
}
