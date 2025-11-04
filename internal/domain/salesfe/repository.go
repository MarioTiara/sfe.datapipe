package salesfe

type Repository interface {
	Save(sale SalesFE) error
	SaveRange(sales []SalesFE) error
}
