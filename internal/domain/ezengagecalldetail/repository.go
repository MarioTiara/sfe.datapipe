package ezengagecalldetail

import "context"

type Repository interface {
	Save(ctx context.Context, calldetail *EZEngageCallDetail) error
	SaveRange(ctx context.Context, calldetail_list []*EZEngageCallDetail) error
	HasThisMonthData(ctx context.Context) (bool, error)
	RemoveThisMonthData(ctx context.Context) (int, error)
}
