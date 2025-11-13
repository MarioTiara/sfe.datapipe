package hirarki

import "context"

type Repository interface {
	Save(ctx context.Context, hirarki *Hirarki) error
	SaveRange(ctx context.Context, hirarkilist []*Hirarki) error
	HasThisMonthData(ctx context.Context) (bool, error)
	RemoveThisMonthData(ctx context.Context) (int, error)
}
