package postgres

import (
	"context"
	"database/sql"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/salesfe"
)

type SalesRepository struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) *SalesRepository {
	return &SalesRepository{db: db}
}

func (r *SalesRepository) Save(sale salesfe.SalesFE) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO sales_fe (
			invoice_date, po_number, po_type, po_type_desc, customer_code, customer_name,
			channel_ic1, channel_ic1_description, channel_ic4, channel_ic4_description,
			plant, branch, principal, product_group, item_code, item_name,
			sales, net_sales, sales_unit, bonus_unit, bun1
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21
		)

	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		sale.InvoiceDate, sale.PONumber, sale.POType, sale.POTypeDesc, sale.CustomerCode,
		sale.CustomerName, sale.ChannelIC1, sale.ChannelIC1Description, sale.ChannelIC4,
		sale.ChannelIC4Description, sale.Plant, sale.Branch, sale.Principal, sale.ProductGroup,
		sale.ItemCode, sale.ItemName, sale.Sales, sale.NetSales, sale.SalesUnit, sale.BonusUnit, sale.BUN1,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SalesRepository) SaveRange(sales []salesfe.SalesFE) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO sales_fe (
			invoice_date, po_number, po_type, po_type_desc, customer_code, customer_name,
			channel_ic1, channel_ic1_description, channel_ic4, channel_ic4_description,
			plant, branch, principal, product_group, item_code, item_name,
			sales, net_sales, sales_unit, bonus_unit, bun1
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21
		)

	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, s := range sales {
		_, err := stmt.ExecContext(ctx,
			s.InvoiceDate, s.PONumber, s.POType, s.POTypeDesc, s.CustomerCode,
			s.CustomerName, s.ChannelIC1, s.ChannelIC1Description, s.ChannelIC4,
			s.ChannelIC4Description, s.Plant, s.Branch, s.Principal, s.ProductGroup,
			s.ItemCode, s.ItemName, s.Sales, s.NetSales, s.SalesUnit, s.BonusUnit, s.BUN1,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
