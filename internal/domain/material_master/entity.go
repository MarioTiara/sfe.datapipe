package materialmaster

type MaterialMaster struct {
	PrincipalCode       string
	PrincipalName       string
	MaterialCode        string
	MaterialName        string
	Uom                 string
	Division            string
	DivisionDescription string
	ProductHierarchy    string
	Brand               string
	MATKL               string
	PH2Code             string
	PH2Name             string
	PH3Code             string
	PH3Name             string
	ATCCode             string
	ATCDescription      string
	DGIndicatorCode     string
	DGIndicatorName     string
	GeneralDGCode       string
	GeneralDGName       string
}

func NewMaterialMaster(
	principalCode, principalName, materialCode, materialName, uom,
	division, divisionDescription, productHierarchy, brand, matkl,
	ph2Code, ph2Name, ph3Code, ph3Name, atcCode, atcDescription,
	dgIndicatorCode, dgIndicatorName, generalDGCode, generalDGName string,
) MaterialMaster {
	return MaterialMaster{
		PrincipalCode:       principalCode,
		PrincipalName:       principalName,
		MaterialCode:        materialCode,
		MaterialName:        materialName,
		Uom:                 uom,
		Division:            division,
		DivisionDescription: divisionDescription,
		ProductHierarchy:    productHierarchy,
		Brand:               brand,
		MATKL:               matkl,
		PH2Code:             ph2Code,
		PH2Name:             ph2Name,
		PH3Code:             ph3Code,
		PH3Name:             ph3Name,
		ATCCode:             atcCode,
		ATCDescription:      atcDescription,
		DGIndicatorCode:     dgIndicatorCode,
		DGIndicatorName:     dgIndicatorName,
		GeneralDGCode:       generalDGCode,
		GeneralDGName:       generalDGName,
	}
}
