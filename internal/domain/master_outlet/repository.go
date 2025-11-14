package masteroutlet

import "context"

type Repository interface {
	Save(ctx context.Context, material *MasterOutlet) error
	SaveRange(ctx context.Context, materials []*MasterOutlet) error
	HasThisMonthData(ctx context.Context) (bool, error)
	RemoveThisMonthData(ctx context.Context) (int, error)
}
