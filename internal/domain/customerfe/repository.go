package customerfe

import "context"

type Repository interface {
	Save(ctx context.Context, customer *CustomerFE) error
	SaveRange(ctx context.Context, customers []*CustomerFE) error
	HasThisMonthData(ctx context.Context) (bool, error)
	RemoveThisMonthData(ctx context.Context) (int, error)
}
