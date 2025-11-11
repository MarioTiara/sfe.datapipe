package hirarki

import (
	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
)

type HirarkiMapper interface {
	MapRowsToHirarki(rows <-chan []string) ([]*hirarki.Hirarki, error)
}
