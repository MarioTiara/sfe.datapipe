package materialmaster

type Repository interface {
	Save(material *MaterialMaster) error
	SaveRange(materials []*MaterialMaster) error
	HasThisMonthData() (bool, error)
	RemoveThisMonthData() (int, error)
}
