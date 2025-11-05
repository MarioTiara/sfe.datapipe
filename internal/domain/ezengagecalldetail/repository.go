package ezengagecalldetail

type Repository interface {
	Save(calldetail EZEngageCallDetail) error
	SaveRange(calldetail_list []EZEngageCallDetail) error
}
