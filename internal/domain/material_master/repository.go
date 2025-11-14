package materialmaster

import "context"

type Repository interface {
	Save(ctx context.Context, material *MaterialMaster) error
	SaveRange(ctx context.Context, materials []*MaterialMaster) error
	HasThisMonthData(ctx context.Context) (bool, error)
	RemoveThisMonthData(ctx context.Context) (int, error)
}
