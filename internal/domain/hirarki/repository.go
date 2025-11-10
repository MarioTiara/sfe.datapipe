package hirarki

type Repository interface {
	Save(hirarki Hirarki) error
	SaveRange(hirarkilist []Hirarki) error
	HasThisMonthData() (bool, error)
	RemoveThisMonthData() (int, error)
}
