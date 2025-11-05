package ezengagecalldetail

type EZEngageCallDetail struct {
	Date                       string
	Team                       string
	Position                   string
	BusinessType               string
	SalesOrg                   string
	TerritoryCode              string
	SalesRepresentative        string
	Status                     string
	MissedCallRemark           string
	PlannedUnplanned           string
	Activity                   string
	OtherworkName              string
	OtherworkNote              string
	CustomerCode               string
	CustomerName               string
	CustomerShipToCode         string
	CustomerShipToName         string
	CustomerCity               string
	CustomerClass              string
	Frequency                  string
	Principal                  string
	Product                    string
	ProductCode                string
	StockInventoryQuantity     string
	Listed                     string
	ShelfStock                 string
	OutOfStock                 string
	Location                   string
	ShelfSpace                 string
	ShareOfSpace               string
	NoOfFacings                string
	NoOfPhotosUpload           string
	ProductEnlistment          string
	PromotionalActivity        string
	SalesRecommendationFromICO string
	PersonInCharge             string
	WithContract               string
	StartDate                  string
	EndDate                    string
	VirtualImages              string
	PresentationDurationTime   string
	PresentationFileName       string
	PerformCollection          string
	ModeOfCollection           string
	PlacedOrder                string
	ModeOfOrder                string
	ReasonOfNotUsingEzrx       string
	CheckInTime                string
	CheckInDate                string
	CheckOutTime               string
	CheckOutDate               string
	VisitDuration              string
	TravelDurationTime         string
	PreCallNotes               string
	PostCallNotes              string
	BarrierEncountered         string
	CallSource                 string
	RouteStatus                string
	ModalityOfCall             string
	SignatureImage             string
	SignatureCount             string
	SignatureStartTime         string
	SignatureEndTime           string
	SignatureDiffTime          string
	LongitudeData              string
	LatitudeData               string
	ActualLongitudeData        string
	ActualLatitudeData         string
	LocationAccuracy           string
	GeoLocationRemarks         string
	Mode                       string
	WorkWith                   string
}

func NewEZEngageCallDetail(
	date, team, position, businessType, salesOrg, territoryCode, salesRep, status,
	missedCallRemark, plannedUnplanned, activity, otherworkName, otherworkNote,
	customerCode, customerName, customerShipToCode, customerShipToName, customerCity,
	customerClass, frequency, principal, product, productCode, stockInventoryQuantity,
	listed, shelfStock, outOfStock, location, shelfSpace, shareOfSpace, noOfFacings,
	noOfPhotosUpload, productEnlistment, promotionalActivity, salesRecommendationFromICO,
	personInCharge, withContract, startDate, endDate, virtualImages, presentationDurationTime,
	presentationFileName, performCollection, modeOfCollection, placedOrder, modeOfOrder,
	reasonOfNotUsingEzrx, checkInTime, checkInDate, checkOutTime, checkOutDate, visitDuration,
	travelDurationTime, preCallNotes, postCallNotes, barrierEncountered, callSource, routeStatus,
	modalityOfCall, signatureImage, signatureCount, signatureStartTime, signatureEndTime,
	signatureDiffTime, longitudeData, latitudeData, actualLongitudeData, actualLatitudeData,
	locationAccuracy, geoLocationRemarks, mode, workWith string,
) EZEngageCallDetail {
	return EZEngageCallDetail{
		Date:                       date,
		Team:                       team,
		Position:                   position,
		BusinessType:               businessType,
		SalesOrg:                   salesOrg,
		TerritoryCode:              territoryCode,
		SalesRepresentative:        salesRep,
		Status:                     status,
		MissedCallRemark:           missedCallRemark,
		PlannedUnplanned:           plannedUnplanned,
		Activity:                   activity,
		OtherworkName:              otherworkName,
		OtherworkNote:              otherworkNote,
		CustomerCode:               customerCode,
		CustomerName:               customerName,
		CustomerShipToCode:         customerShipToCode,
		CustomerShipToName:         customerShipToName,
		CustomerCity:               customerCity,
		CustomerClass:              customerClass,
		Frequency:                  frequency,
		Principal:                  principal,
		Product:                    product,
		ProductCode:                productCode,
		StockInventoryQuantity:     stockInventoryQuantity,
		Listed:                     listed,
		ShelfStock:                 shelfStock,
		OutOfStock:                 outOfStock,
		Location:                   location,
		ShelfSpace:                 shelfSpace,
		ShareOfSpace:               shareOfSpace,
		NoOfFacings:                noOfFacings,
		NoOfPhotosUpload:           noOfPhotosUpload,
		ProductEnlistment:          productEnlistment,
		PromotionalActivity:        promotionalActivity,
		SalesRecommendationFromICO: salesRecommendationFromICO,
		PersonInCharge:             personInCharge,
		WithContract:               withContract,
		StartDate:                  startDate,
		EndDate:                    endDate,
		VirtualImages:              virtualImages,
		PresentationDurationTime:   presentationDurationTime,
		PresentationFileName:       presentationFileName,
		PerformCollection:          performCollection,
		ModeOfCollection:           modeOfCollection,
		PlacedOrder:                placedOrder,
		ModeOfOrder:                modeOfOrder,
		ReasonOfNotUsingEzrx:       reasonOfNotUsingEzrx,
		CheckInTime:                checkInTime,
		CheckInDate:                checkInDate,
		CheckOutTime:               checkOutTime,
		CheckOutDate:               checkOutDate,
		VisitDuration:              visitDuration,
		TravelDurationTime:         travelDurationTime,
		PreCallNotes:               preCallNotes,
		PostCallNotes:              postCallNotes,
		BarrierEncountered:         barrierEncountered,
		CallSource:                 callSource,
		RouteStatus:                routeStatus,
		ModalityOfCall:             modalityOfCall,
		SignatureImage:             signatureImage,
		SignatureCount:             signatureCount,
		SignatureStartTime:         signatureStartTime,
		SignatureEndTime:           signatureEndTime,
		SignatureDiffTime:          signatureDiffTime,
		LongitudeData:              longitudeData,
		LatitudeData:               latitudeData,
		ActualLongitudeData:        actualLongitudeData,
		ActualLatitudeData:         actualLatitudeData,
		LocationAccuracy:           locationAccuracy,
		GeoLocationRemarks:         geoLocationRemarks,
		Mode:                       mode,
		WorkWith:                   workWith,
	}
}
