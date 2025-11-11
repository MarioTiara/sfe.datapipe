package excel

import (
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/customerfe"
)

type excelCustomerFEMapper struct{}

func NewExcelCustomerFEMapper() *excelCustomerFEMapper {
	return &excelCustomerFEMapper{}
}

func (m *excelCustomerFEMapper) MapRowsToCustomerFE(rows <-chan []string) ([]*customerfe.CustomerFE, error) {
	var result []*customerfe.CustomerFE
	for row := range rows {
		if len(row) == 0 {
			continue
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
		result = append(result, c)

	}

	return result, nil
}
