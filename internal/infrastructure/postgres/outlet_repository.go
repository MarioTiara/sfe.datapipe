package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariotiara/sfe-data-pipe/configs"
	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type MasterOutletRepository struct {
	logger  logger.Logger
	db      *sql.DB
	config  *configs.Config
	pgxPool *pgxpool.Pool
}

func NewMasterOutletRepository(db *sql.DB, config *configs.Config, logger logger.Logger) (*MasterOutletRepository, error) {

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.ConnString) // e.g. "postgres://user:pass@localhost/db"
	if err != nil {
		return nil, err
	}
	return &MasterOutletRepository{db: db, pgxPool: pool, config: config, logger: logger}, nil
}

func (r *MasterOutletRepository) RemoveThisMonthData(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		DELETE FROM master_outlet
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

func (r *MasterOutletRepository) HasThisMonthData(ctx context.Context) (bool, error) {
	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		SELECT COUNT(1)
		FROM master_outlet
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

func (r *MasterOutletRepository) Save(ctx context.Context, m *masteroutlet.MasterOutlet) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO master_outlet (
			customer_code, customer_name, channel, plant, branch_name,
			nik_salesman, name_salesman, rayon_code, rayon, new_class,
			call_plan_full_month, target_freq, terr_code, username
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		m.CustomerCode, m.CustomerName, m.Channel, m.Plant, m.BranchName,
		m.NIKSalesman, m.NameSalesman, m.RayonCode, m.Rayon, m.NewClass,
		m.CallPlanFullMonth, m.TargetFreq, m.TerrCode, m.Username,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *MasterOutletRepository) SaveRange(
	ctx context.Context,
	ms []*masteroutlet.MasterOutlet,
) error {

	batchSize := r.config.DBBatchSize
	start := time.Now()
	totalInserted := 0

	for i := 0; i < len(ms); i += batchSize {
		batchStart := i
		batchEnd := i + batchSize
		if batchEnd > len(ms) {
			batchEnd = len(ms)
		}

		batch := ms[batchStart:batchEnd]
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
			pgx.Identifier{"master_outlet"},
			[]string{
				"customer_code", "customer_name", "channel", "plant", "branch_name",
				"nik_salesman", "name_salesman", "rayon_code", "rayon", "new_class",
				"call_plan_full_month", "target_freq", "terr_code", "username",
			},
			pgx.CopyFromSlice(len(batch), func(i int) ([]any, error) {
				m := batch[i]
				return []any{
					m.CustomerCode, m.CustomerName, m.Channel, m.Plant, m.BranchName,
					m.NIKSalesman, m.NameSalesman, m.RayonCode, m.Rayon, m.NewClass,
					m.CallPlanFullMonth, m.TargetFreq, m.TerrCode, m.Username,
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
