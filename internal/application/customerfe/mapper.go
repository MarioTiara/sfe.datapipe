package customerfe

import (
	"github.com/mariotiara/sfe-data-pipe/internal/domain/customerfe"
)

type CustomerFEMapper interface {
	MapRowsToCustomerFE(rows <-chan []string) (<-chan *customerfe.CustomerFE, <-chan error)
}
