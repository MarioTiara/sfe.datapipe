package salesfe

type Repository interface {
	Save(sale SalesFE) error
	SaveRange(sales []SalesFE) error
	HasThisMonthData() (bool, error)
	RemoveThisMonthData() (int, error)
}
