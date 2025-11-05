package masteroutlet

type Repository interface {
	Save(material MasterOutlet) error
	SaveRange(materials []MasterOutlet) error
}
