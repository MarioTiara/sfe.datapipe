package masteroutlet

type MasterOutlet struct {
	ID                string
	CustomerCode      string
	CustomerName      string
	Channel           string
	Plant             string
	BranchName        string
	NIKSalesman       string
	NameSalesman      string
	RayonCode         string
	Rayon             string
	NewClass          string
	CallPlanFullMonth string
	TargetFreq        string
	TerrCode          string
	Username          string
}

func NewMasterOutlet(
	customerCode, customerName, channel, plant, branchName,
	nikSalesman, nameSalesman, rayonCode, rayon, newClass,
	callPlanFullMonth, targetFreq, terrCode, username string,
) MasterOutlet {
	return MasterOutlet{
		CustomerCode:      customerCode,
		CustomerName:      customerName,
		Channel:           channel,
		Plant:             plant,
		BranchName:        branchName,
		NIKSalesman:       nikSalesman,
		NameSalesman:      nameSalesman,
		RayonCode:         rayonCode,
		Rayon:             rayon,
		NewClass:          newClass,
		CallPlanFullMonth: callPlanFullMonth,
		TargetFreq:        targetFreq,
		TerrCode:          terrCode,
		Username:          username,
	}
}
