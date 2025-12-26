package tabular

import (
	"context"
	"fmt"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/ezengagecalldetail"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type excelEZEnggaeCallDetailMapper struct {
	logger logger.Logger
}

func NewExcelEZEnggaeCallDetailMapper(logger logger.Logger) *excelEZEnggaeCallDetailMapper {
	return &excelEZEnggaeCallDetailMapper{logger: logger}
}

func (m *excelEZEnggaeCallDetailMapper) MapRowsToCallDetails(
	ctx context.Context,
	rows <-chan []string,
) ([]*ezengagecalldetail.EZEngageCallDetail, error) {
	var result []*ezengagecalldetail.EZEngageCallDetail

	for row := range rows {
		if len(row) == 0 {
			continue
		}

		e := &ezengagecalldetail.EZEngageCallDetail{
			Date:                       safeGet(row, 0),
			Team:                       safeGet(row, 1),
			Position:                   safeGet(row, 2),
			BusinessType:               safeGet(row, 3),
			SalesOrg:                   safeGet(row, 4),
			TerritoryCode:              safeGet(row, 5),
			SalesRepresentative:        safeGet(row, 6),
			Status:                     safeGet(row, 7),
			MissedCallRemark:           safeGet(row, 8),
			PlannedUnplanned:           safeGet(row, 9),
			Activity:                   safeGet(row, 10),
			OtherworkName:              safeGet(row, 11),
			OtherworkNote:              safeGet(row, 12),
			CustomerCode:               safeGet(row, 13),
			CustomerName:               safeGet(row, 14),
			CustomerShipToCode:         safeGet(row, 15),
			CustomerShipToName:         safeGet(row, 16),
			CustomerCity:               safeGet(row, 17),
			CustomerClass:              safeGet(row, 18),
			Frequency:                  safeGet(row, 19),
			Principal:                  safeGet(row, 20),
			Product:                    safeGet(row, 21),
			ProductCode:                safeGet(row, 22),
			StockInventoryQuantity:     safeGet(row, 23),
			Listed:                     safeGet(row, 24),
			ShelfStock:                 safeGet(row, 25),
			OutOfStock:                 safeGet(row, 26),
			Location:                   safeGet(row, 27),
			ShelfSpace:                 safeGet(row, 28),
			ShareOfSpace:               safeGet(row, 29),
			NoOfFacings:                safeGet(row, 30),
			NoOfPhotosUpload:           safeGet(row, 31),
			ProductEnlistment:          safeGet(row, 32),
			PromotionalActivity:        safeGet(row, 33),
			SalesRecommendationFromICO: safeGet(row, 34),
			PersonInCharge:             safeGet(row, 35),
			WithContract:               safeGet(row, 36),
			StartDate:                  safeGet(row, 37),
			EndDate:                    safeGet(row, 38),
			VirtualImages:              safeGet(row, 39),
			PresentationDurationTime:   safeGet(row, 40),
			PresentationFileName:       safeGet(row, 41),
			PerformCollection:          safeGet(row, 42),
			ModeOfCollection:           safeGet(row, 43),
			PlacedOrder:                safeGet(row, 44),
			ModeOfOrder:                safeGet(row, 45),
			ReasonOfNotUsingEzrx:       safeGet(row, 46),
			CheckInTime:                safeGet(row, 47),
			CheckInDate:                safeGet(row, 48),
			CheckOutTime:               safeGet(row, 49),
			CheckOutDate:               safeGet(row, 50),
			VisitDuration:              safeGet(row, 51),
			TravelDurationTime:         safeGet(row, 52),
			PreCallNotes:               safeGet(row, 53),
			PostCallNotes:              safeGet(row, 54),
			BarrierEncountered:         safeGet(row, 55),
			CallSource:                 safeGet(row, 56),
			RouteStatus:                safeGet(row, 57),
			ModalityOfCall:             safeGet(row, 58),
			SignatureImage:             safeGet(row, 59),
			SignatureCount:             safeGet(row, 60),
			SignatureStartTime:         safeGet(row, 61),
			SignatureEndTime:           safeGet(row, 62),
			SignatureDiffTime:          safeGet(row, 63),
			LongitudeData:              safeGet(row, 64),
			LatitudeData:               safeGet(row, 65),
			ActualLongitudeData:        safeGet(row, 66),
			ActualLatitudeData:         safeGet(row, 67),
			LocationAccuracy:           safeGet(row, 68),
			GeoLocationRemarks:         safeGet(row, 69),
			Mode:                       safeGet(row, 70),
			WorkWith:                   safeGet(row, 71),
		}

		result = append(result, e)
	}

	m.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(result)))
	return result, nil
}
