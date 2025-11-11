package ezengagecalldetail

import (
	"github.com/mariotiara/sfe-data-pipe/internal/domain/ezengagecalldetail"
)

type CallDetailMapper interface {
	MapRowsToCallDetails(rows <-chan []string) ([]*ezengagecalldetail.EZEngageCallDetail, error)
}
