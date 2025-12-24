package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariotiara/sfe-data-pipe/configs"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/salesfe"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

var config, _ = configs.Load()

type SalesRepository struct {
	logger  logger.Logger
	db      *sql.DB
	config  *configs.Config
	pgxPool *pgxpool.Pool
}

func NewSalesRepository(db *sql.DB, config *configs.Config, logger logger.Logger) (*SalesRepository, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.ConnString) // e.g. "postgres://user:pass@localhost/db"
	if err != nil {
		return nil, err
	}
	return &SalesRepository{db: db, pgxPool: pool, config: config, logger: logger}, nil
}

func (r *SalesRepository) RemoveThisMonthData(ctx context.Context) (int, error) {
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

func (r *SalesRepository) HasThisMonthData(ctx context.Context) (bool, error) {

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

	r.logger.Info(ctx, fmt.Sprintf("%d Rows has been removed", count))
	return count > 0, nil

}

func (r *SalesRepository) Save(ctx context.Context, sale *salesfe.SalesFE) error {
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

func (r *SalesRepository) SaveRange(
	ctx context.Context,
	sales []*salesfe.SalesFE,
) error {

	batchSize := r.config.DBBatchSize
	start := time.Now()
	totalInserted := 0

	for i := 0; i < len(sales); i += batchSize {
		batchStart := i
		batchEnd := i + batchSize
		if batchEnd > len(sales) {
			batchEnd = len(sales)
		}

		batch := sales[batchStart:batchEnd]
		batchNumber := (i / batchSize) + 1
		batchStartTime := time.Now()

		conn, err := r.pgxPool.Acquire(ctx)
		if err != nil {
			r.logger.Error(ctx, "failed to acquire pgx connection", err,
				logger.Field{Key: "batch_number", Value: batchNumber},
				logger.Field{Key: "error", Value: err.Error()},
			)
			continue
		}

		rowsInserted, err := conn.CopyFrom(
			ctx,
			pgx.Identifier{"sales_fe"},
			[]string{
				"invoice_date", "po_number", "po_type", "po_type_desc", "customer_code", "customer_name",
				"channel_ic1", "channel_ic1_description", "channel_ic4", "channel_ic4_description",
				"plant", "branch", "principal", "product_group", "item_code", "item_name",
				"sales", "net_sales", "sales_unit", "bonus_unit", "bun1",
			},
			pgx.CopyFromSlice(len(batch), func(i int) ([]any, error) {
				s := batch[i]
				return []any{
					s.InvoiceDate, s.PONumber, s.POType, s.POTypeDesc, s.CustomerCode,
					s.CustomerName, s.ChannelIC1, s.ChannelIC1Description, s.ChannelIC4,
					s.ChannelIC4Description, s.Plant, s.Branch, s.Principal, s.ProductGroup,
					s.ItemCode, s.ItemName, s.Sales, s.NetSales, s.SalesUnit, s.BonusUnit, s.BUN1,
				}, nil
			}),
		)

		conn.Release()

		if err != nil {
			r.logger.Error(ctx, "copy from batch failed", err,
				logger.Field{Key: "batch_number", Value: batchNumber},
				logger.Field{Key: "batch_start", Value: batchStart},
				logger.Field{Key: "batch_end", Value: batchEnd},
				logger.Field{Key: "batch_size", Value: len(batch)},
				logger.Field{Key: "error", Value: err.Error()},
			)
			continue
		}

		totalInserted += int(rowsInserted)

		r.logger.Info(ctx, "batch copy succeeded",
			logger.Field{Key: "batch_number", Value: batchNumber},
			logger.Field{Key: "inserted_rows", Value: rowsInserted},
			logger.Field{Key: "batch_duration", Value: time.Since(batchStartTime).String()},
		)
	}

	r.logger.Info(ctx, "bulk insert finished",
		logger.Field{Key: "total_inserted", Value: totalInserted},
		logger.Field{Key: "total_duration", Value: time.Since(start).String()},
	)

	return nil
}
