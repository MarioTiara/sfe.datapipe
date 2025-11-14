package salesfe

import (
	"context"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/salesfe"
)

type SalesFEMapper interface {
	MapRowsToSaleFE(ctx context.Context, rows <-chan []string) ([]*salesfe.SalesFE, error)
}
