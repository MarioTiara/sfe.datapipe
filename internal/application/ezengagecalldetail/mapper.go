package ezengagecalldetail

import (
	"context"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/ezengagecalldetail"
)

type CallDetailMapper interface {
	MapRowsToCallDetails(ctx context.Context, rows <-chan []string) ([]*ezengagecalldetail.EZEngageCallDetail, error)
}
