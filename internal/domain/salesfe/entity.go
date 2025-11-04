package salesfe

import (
	"errors"
	"time"
)

type SalesFE struct {
	InvoiceDate           time.Time
	PONumber              string
	POType                string
	POTypeDesc            string
	CustomerCode          string
	CustomerName          string
	ChannelIC1            string
	ChannelIC1Description string
	ChannelIC4            string
	ChannelIC4Description string
	Plant                 string
	Branch                string
	Principal             string
	ProductGroup          string
	ItemCode              string
	ItemName              string
	Sales                 float64
	NetSales              float64
	SalesUnit             float64
	BonusUnit             float64
	BUN1                  string
}

func NewSalesFE(
	invoiceDate time.Time,
	poNumber, poType, poTypeDesc, custCode, custName,
	channelIC1, channelIC1Desc, channelIC4, channelIC4Desc,
	plant, branch, principal, productGroup, itemCode, itemName,
	bun1 string,
	sales, netSales, salesUnit, bonusUnit float64,

) (*SalesFE, error) {
	if poNumber == "" {
		return nil, errors.New("PO number is required")
	}
	if custCode == "" {
		return nil, errors.New("Customer code is required")
	}
	if itemCode == "" {
		return nil, errors.New("Item code is required")
	}

	return &SalesFE{
		InvoiceDate:           invoiceDate,
		PONumber:              poNumber,
		POType:                poType,
		POTypeDesc:            poTypeDesc,
		CustomerCode:          custCode,
		CustomerName:          custName,
		ChannelIC1:            channelIC1,
		ChannelIC1Description: channelIC1Desc,
		ChannelIC4:            channelIC4,
		ChannelIC4Description: channelIC4Desc,
		Plant:                 plant,
		Branch:                branch,
		Principal:             principal,
		ProductGroup:          productGroup,
		ItemCode:              itemCode,
		ItemName:              itemName,
		Sales:                 sales,
		NetSales:              netSales,
		SalesUnit:             salesUnit,
		BonusUnit:             bonusUnit,
		BUN1:                  bun1,
	}, nil
}
