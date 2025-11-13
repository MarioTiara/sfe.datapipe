package salesfe

import "context"

type Repository interface {
	Save(ctx context.Context, sale *SalesFE) error
	SaveRange(ctx context.Context, sales []*SalesFE) error
	HasThisMonthData(ctx context.Context) (bool, error)
	RemoveThisMonthData(ctx context.Context) (int, error)
}
