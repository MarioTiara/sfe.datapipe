package hirarki

type Hirarki struct {
	RayonCode              string
	Plant                  string
	RayonType              string
	BUM                    string
	BUMName                string
	NSM                    string
	NSMName                string
	ASM                    string
	ASMName                string
	FSS                    string
	FSSName                string
	SLM                    string
	SLMName                string
	SalesmanCategoryUpdate string
	BranchName             string
	MLO                    string
	Remarks                string
	TerrCode               string
	Username               string
	Change                 string
	CategoryRayon          string
	RayonDetail            string
}

func NewHirarki(
	rayonCode, plant, rayonType, bum, bumName, nsm, nsmName, asm, asmName,
	fss, fssName, slm, slmName, salesmanCatUpdate, branchName, mlo, remarks,
	terrCode, username, change, categoryRayon, rayonDetail string,
) Hirarki {
	return Hirarki{
		RayonCode:              rayonCode,
		Plant:                  plant,
		RayonType:              rayonType,
		BUM:                    bum,
		BUMName:                bumName,
		NSM:                    nsm,
		NSMName:                nsmName,
		ASM:                    asm,
		ASMName:                asmName,
		FSS:                    fss,
		FSSName:                fssName,
		SLM:                    slm,
		SLMName:                slmName,
		SalesmanCategoryUpdate: salesmanCatUpdate,
		BranchName:             branchName,
		MLO:                    mlo,
		Remarks:                remarks,
		TerrCode:               terrCode,
		Username:               username,
		Change:                 change,
		CategoryRayon:          categoryRayon,
		RayonDetail:            rayonDetail,
	}
}
