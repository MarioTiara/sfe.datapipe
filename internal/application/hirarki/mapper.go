package hirarki

import (
	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
)

type HirarkiMapper interface {
	MapRowsToCallDetails(rows <-chan []string) (<-chan *hirarki.Hirarki, <-chan error)
}
