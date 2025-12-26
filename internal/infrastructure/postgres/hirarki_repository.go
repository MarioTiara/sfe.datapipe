package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariotiara/sfe-data-pipe/configs"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type HirarkiRepository struct {
	logger  logger.Logger
	db      *sql.DB
	config  *configs.Config
	pgxPool *pgxpool.Pool
}

func NewHirarkiRepository(db *sql.DB, config *configs.Config, logger logger.Logger) (*HirarkiRepository, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.ConnString) // e.g. "postgres://user:pass@localhost/db"
	if err != nil {
		return nil, err
	}
	return &HirarkiRepository{db: db, pgxPool: pool, config: config, logger: logger}, nil
}

func (r *HirarkiRepository) HasThisMonthData(ctx context.Context) (bool, error) {
	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		SELECT COUNT(1)
		FROM hirarki
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

func (r *HirarkiRepository) RemoveThisMonthData(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		DELETE FROM hirarki
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

func (r *HirarkiRepository) Save(ctx context.Context, h *hirarki.Hirarki) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO hirarki (
			rayon_code, plant, rayon_type, bum, bum_name, nsm, nsm_name,
			asm, asm_name, fss, fss_name, slm, slm_name,
			salesman_category_update, branch_name, mlo, remarks,
			terr_code, username, change, category_rayon, rayon_detail
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17,
			$18, $19, $20, $21, $22
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		h.RayonCode, h.Plant, h.RayonType, h.BUM, h.BUMName, h.NSM, h.NSMName,
		h.ASM, h.ASMName, h.FSS, h.FSSName, h.SLM, h.SLMName,
		h.SalesmanCategoryUpdate, h.BranchName, h.MLO, h.Remarks,
		h.TerrCode, h.Username, h.Change, h.CategoryRayon, h.RayonDetail,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *HirarkiRepository) SaveRange(
	ctx context.Context,
	hs []*hirarki.Hirarki,
) error {

	batchSize := r.config.DBBatchSize
	start := time.Now()
	totalInserted := 0

	for i := 0; i < len(hs); i += batchSize {
		batchStart := i
		batchEnd := i + batchSize
		if batchEnd > len(hs) {
			batchEnd = len(hs)
		}

		batch := hs[batchStart:batchEnd]
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
			pgx.Identifier{"hirarki"},
			[]string{
				"rayon_code", "plant", "rayon_type", "bum", "bum_name", "nsm", "nsm_name",
				"asm", "asm_name", "fss", "fss_name", "slm", "slm_name",
				"salesman_category_update", "branch_name", "mlo", "remarks",
				"terr_code", "username", "change", "category_rayon", "rayon_detail",
			},
			pgx.CopyFromSlice(len(batch), func(i int) ([]any, error) {
				h := batch[i]
				return []any{
					h.RayonCode, h.Plant, h.RayonType, h.BUM, h.BUMName, h.NSM, h.NSMName,
					h.ASM, h.ASMName, h.FSS, h.FSSName, h.SLM, h.SLMName,
					h.SalesmanCategoryUpdate, h.BranchName, h.MLO, h.Remarks,
					h.TerrCode, h.Username, h.Change, h.CategoryRayon, h.RayonDetail,
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
