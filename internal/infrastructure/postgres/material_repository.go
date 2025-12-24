package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariotiara/sfe-data-pipe/configs"
	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type MaterialMasterRepository struct {
	logger  logger.Logger
	db      *sql.DB
	config  *configs.Config
	pgxPool *pgxpool.Pool
}

func NewMaterialMasterRepository(db *sql.DB, config *configs.Config, logger logger.Logger) (*MaterialMasterRepository, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.ConnString) // e.g. "postgres://user:pass@localhost/db"
	if err != nil {
		return nil, err
	}
	return &MaterialMasterRepository{db: db, pgxPool: pool, config: config, logger: logger}, nil
}

func (r *MaterialMasterRepository) RemoveThisMonthData(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		DELETE FROM material_master
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

	r.logger.Info(ctx, fmt.Sprintf("%d Rows has been removed", count))
	return count, nil
}

func (r *MaterialMasterRepository) HasThisMonthData(ctx context.Context) (bool, error) {

	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		SELECT COUNT(1)
		FROM material_master
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

// Save inserts a single MaterialMaster record
func (r *MaterialMasterRepository) Save(ctx context.Context, m *materialmaster.MaterialMaster) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO material_master (
			principal_code, principal_name, material_code, material_name,
			uom, division, division_description, product_hierarchy, brand,
			matkl, ph2_code, ph2_name, ph3_code, ph3_name,
			atc_code, atc_description, dg_indicator_code, dg_indicator_name,
			general_dg_code, general_dg_name
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17, $18,
			$19, $20
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		m.PrincipalCode, m.PrincipalName, m.MaterialCode, m.MaterialName,
		m.Uom, m.Division, m.DivisionDescription, m.ProductHierarchy, m.Brand,
		m.MATKL, m.PH2Code, m.PH2Name, m.PH3Code, m.PH3Name,
		m.ATCCode, m.ATCDescription, m.DGIndicatorCode, m.DGIndicatorName,
		m.GeneralDGCode, m.GeneralDGName,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SaveRange inserts multiple MaterialMaster records efficiently in a single transaction
func (r *MaterialMasterRepository) SaveRange(
	ctx context.Context,
	materials []*materialmaster.MaterialMaster,
) error {

	batchSize := r.config.DBBatchSize
	start := time.Now()
	totalInserted := 0

	for i := 0; i < len(materials); i += batchSize {
		batchStart := i
		batchEnd := i + batchSize
		if batchEnd > len(materials) {
			batchEnd = len(materials)
		}

		batch := materials[batchStart:batchEnd]
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
			pgx.Identifier{"material_master"},
			[]string{
				"principal_code", "principal_name", "material_code", "material_name",
				"uom", "division", "division_description", "product_hierarchy", "brand",
				"matkl", "ph2_code", "ph2_name", "ph3_code", "ph3_name",
				"atc_code", "atc_description", "dg_indicator_code", "dg_indicator_name",
				"general_dg_code", "general_dg_name",
			},
			pgx.CopyFromSlice(len(batch), func(i int) ([]any, error) {
				m := batch[i]
				return []any{
					m.PrincipalCode, m.PrincipalName, m.MaterialCode, m.MaterialName,
					m.Uom, m.Division, m.DivisionDescription, m.ProductHierarchy, m.Brand,
					m.MATKL, m.PH2Code, m.PH2Name, m.PH3Code, m.PH3Name,
					m.ATCCode, m.ATCDescription, m.DGIndicatorCode, m.DGIndicatorName,
					m.GeneralDGCode, m.GeneralDGName,
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
