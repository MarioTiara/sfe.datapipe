package excel

import (
	"fmt"
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/customerfe"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/ezengagecalldetail"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"
)

// MapRowsToEZEngageCallDetail maps rows to domain entities
func MapRowsToEZEngageCallDetail(rows <-chan []string) (<-chan *ezengagecalldetail.EZEngageCallDetail, <-chan error) {
	fmt.Println("start mapping")
	outCh := make(chan *ezengagecalldetail.EZEngageCallDetail)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)

		for row := range rows {
			if len(row) == 0 {
				continue // skip empty rows
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

			outCh <- e
		}
	}()

	fmt.Println("mapping finished")
	return outCh, errCh
}

func MapRowsToCustomerFE(rows <-chan []string) (<-chan *customerfe.CustomerFE, <-chan error) {
	outCh := make(chan *customerfe.CustomerFE)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)

		for row := range rows {
			if len(row) == 0 {
				continue // skip empty rows
			}

			var datePtr *time.Time
			if len(row) > 34 {
				if parsed, err := time.Parse("20060102", row[34]); err == nil {
					datePtr = &parsed
				}
			}

			c := &customerfe.CustomerFE{
				CustomerCode:      safeGet(row, 0),
				CustomerName:      safeGet(row, 1),
				Address:           safeGet(row, 2),
				CityName:          safeGet(row, 3),
				BankCountry:       safeGet(row, 4),
				Regio:             safeGet(row, 5),
				TransZone:         safeGet(row, 6),
				Telephone:         safeGet(row, 7),
				Bran1:             safeGet(row, 8),
				Bran2:             safeGet(row, 9),
				Bran3:             safeGet(row, 10),
				ChannelIc4:        safeGet(row, 11),
				Vtext:             safeGet(row, 12),
				Katr1:             safeGet(row, 13),
				Katr2:             safeGet(row, 14),
				Katr3:             safeGet(row, 15),
				Adrnr:             safeGet(row, 16),
				Vkorg:             safeGet(row, 17),
				SalesOffice:       safeGet(row, 18),
				SalesDistrict:     safeGet(row, 19),
				CustomerGrp1:      safeGet(row, 20),
				CustGrp1Desc:      safeGet(row, 21),
				CustomerGrp2:      safeGet(row, 22),
				ShippingCondition: safeGet(row, 23),
				Lprio:             safeGet(row, 24),
				Branch:            safeGet(row, 25),
				Eikto:             safeGet(row, 26),
				Ktokd:             safeGet(row, 27),
				FlagDeletion:      safeGet(row, 28),
				Sperr1:            safeGet(row, 29),
				Aufsd1:            safeGet(row, 30),
				Lifsd1:            safeGet(row, 31),
				Faksd1:            safeGet(row, 32),
				Cassd1:            safeGet(row, 33),
				Erdat:             datePtr,
				Ernam:             safeGet(row, 35),
				PostalCode:        safeGet(row, 36),
			}

			outCh <- c
		}
	}()

	return outCh, errCh
}

func MapRowsToHirarki(rows <-chan []string) (<-chan *hirarki.Hirarki, <-chan error) {
	outCh := make(chan *hirarki.Hirarki)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)

		for row := range rows {
			if len(row) == 0 {
				continue // skip empty rows
			}

			h := &hirarki.Hirarki{
				RayonCode:              safeGet(row, 0),
				Plant:                  safeGet(row, 1),
				RayonType:              safeGet(row, 2),
				BUM:                    safeGet(row, 3),
				BUMName:                safeGet(row, 4),
				NSM:                    safeGet(row, 5),
				NSMName:                safeGet(row, 6),
				ASM:                    safeGet(row, 7),
				ASMName:                safeGet(row, 8),
				FSS:                    safeGet(row, 9),
				FSSName:                safeGet(row, 10),
				SLM:                    safeGet(row, 11),
				SLMName:                safeGet(row, 12),
				SalesmanCategoryUpdate: safeGet(row, 13),
				BranchName:             safeGet(row, 14),
				MLO:                    safeGet(row, 15),
				Remarks:                safeGet(row, 16),
				TerrCode:               safeGet(row, 17),
				Username:               safeGet(row, 18),
				Change:                 safeGet(row, 19),
				CategoryRayon:          safeGet(row, 20),
				RayonDetail:            safeGet(row, 21),
			}

			outCh <- h
		}
	}()

	return outCh, errCh
}

func MapRowsToMaterialMaster(rows <-chan []string) (<-chan *materialmaster.MaterialMaster, <-chan error) {
	outCh := make(chan *materialmaster.MaterialMaster)
	errCh := make(chan error, 1)

	go func() {
		defer close(outCh)
		defer close(errCh)

		for row := range rows {
			if len(row) == 0 {
				continue // skip empty rows
			}

			m := &materialmaster.MaterialMaster{
				PrincipalCode:       safeGet(row, 0),
				PrincipalName:       safeGet(row, 1),
				MaterialCode:        safeGet(row, 2),
				MaterialName:        safeGet(row, 3),
				Uom:                 safeGet(row, 4),
				Division:            safeGet(row, 5),
				DivisionDescription: safeGet(row, 6),
				ProductHierarchy:    safeGet(row, 7),
				Brand:               safeGet(row, 8),
				MATKL:               safeGet(row, 9),
				PH2Code:             safeGet(row, 10),
				PH2Name:             safeGet(row, 11),
				PH3Code:             safeGet(row, 12),
				PH3Name:             safeGet(row, 13),
				ATCCode:             safeGet(row, 14),
				ATCDescription:      safeGet(row, 15),
				DGIndicatorCode:     safeGet(row, 16),
				DGIndicatorName:     safeGet(row, 17),
				GeneralDGCode:       safeGet(row, 18),
				GeneralDGName:       safeGet(row, 19),
			}

			outCh <- m
		}
	}()

	return outCh, errCh
}
