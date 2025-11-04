package customerfe

type Repository interface {
	Save(customer CustomerFE) error
	SaveRange(customers []CustomerFE) error
}
