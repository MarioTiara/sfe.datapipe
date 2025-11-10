package salesfe

import (
	"github.com/mariotiara/sfe-data-pipe/internal/domain/salesfe"
)

type SalesFEMapper interface {
	MapRowsToSaleFE(rows <-chan []string) (<-chan *salesfe.SalesFE, <-chan error)
}
