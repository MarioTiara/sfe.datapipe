package hirarki

import (
	"context"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
)

type HirarkiMapper interface {
	MapRowsToHirarki(ctx context.Context, rows <-chan []string) ([]*hirarki.Hirarki, error)
}
