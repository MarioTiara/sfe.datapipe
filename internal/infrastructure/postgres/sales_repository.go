package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mariotiara/sfe-data-pipe/configs"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/salesfe"
)

var config, _ = configs.Load()

type SalesRepository struct {
	db     *sql.DB
	config *configs.Config
}

func NewSalesRepository(db *sql.DB, config *configs.Config) *SalesRepository {
	return &SalesRepository{db: db, config: config}
}

func (r *SalesRepository) RemoveThisMonthData() (int, error) {
	ctx := context.Background()

	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		DELETE FROM sales_fe
		WHERE EXTRACT (YEAR FROM created_at)=$1
		AND EXTRACT (MONTH FROM created_at)=$2
		RETURNING 1
	`

	rows, err := r.db.QueryContext(ctx, query, year, month)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	return count, nil
}

func (r *SalesRepository) HasThisMonthData() (bool, error) {
	ctx := context.Background()

	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		SELECT COUNT(1)
		FROM sales_fe
		WHERE EXTRACT (YEAR FROM created_at)=$1
		AND EXTRACT (MONTH FROM created_at)=$2
	`
	var count int

	err := r.db.QueryRowContext(ctx, query, year, month).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *SalesRepository) Save(sale *salesfe.SalesFE) error {
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

func (r *SalesRepository) SaveRange(sales []*salesfe.SalesFE) error {
	batchSize := r.config.DBBatchSize

	ctx := context.Background()
	for i := 0; i < len(sales); i += batchSize {
		end := i + batchSize
		if end > len(sales) {
			end = len(sales)
		}

		msg := fmt.Sprintf("insert %d of %d", i, len(sales))
		log.Println(msg)
		batch := sales[i:end]

		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

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
			tx.Rollback()
			return fmt.Errorf("prepare stmt: %w", err)
		}

		for _, s := range batch {
			_, err := stmt.ExecContext(ctx,
				s.InvoiceDate, s.PONumber, s.POType, s.POTypeDesc, s.CustomerCode,
				s.CustomerName, s.ChannelIC1, s.ChannelIC1Description, s.ChannelIC4,
				s.ChannelIC4Description, s.Plant, s.Branch, s.Principal, s.ProductGroup,
				s.ItemCode, s.ItemName, s.Sales, s.NetSales, s.SalesUnit, s.BonusUnit, s.BUN1,
			)
			if err != nil {
				stmt.Close()
				tx.Rollback()
				return fmt.Errorf("exec batch insert: %w", err)
			}
		}

		stmt.Close()
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit batch: %w", err)
		}

		log.Printf("✅ Committed batch %d–%d successfully\n", i, end)
	}

	return nil
}
