package customerfe

import (
	"context"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/customerfe"
)

type CustomerFEMapper interface {
	MapRowsToCustomerFE(ctx context.Context, rows <-chan []string) ([]*customerfe.CustomerFE, error)
}
