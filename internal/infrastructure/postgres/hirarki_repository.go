package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mariotiara/sfe-data-pipe/configs"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
)

type HirarkiRepository struct {
	db     *sql.DB
	config *configs.Config
}

func NewHirarkiRepository(db *sql.DB, config *configs.Config) *HirarkiRepository {
	return &HirarkiRepository{db: db, config: config}
}

func (r *HirarkiRepository) HasThisMonthData() (bool, error) {
	ctx := context.Background()

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

func (r *HirarkiRepository) RemoveThisMonthData() (int, error) {
	ctx := context.Background()

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

	return count, nil
}

func (r *HirarkiRepository) Save(h *hirarki.Hirarki) error {
	ctx := context.Background()
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

func (r *HirarkiRepository) SaveRange(hs []*hirarki.Hirarki) error {
	batchSize := r.config.DBBatchSize
	ctx := context.Background()
	for i := 0; i < len(hs); i += batchSize {
		end := i + batchSize
		if end > len(hs) {
			end = len(hs)
		}

		msg := fmt.Sprintf("insert %d of %d", i, len(hs))
		log.Println(msg)
		batch := hs[i:end]

		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

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
			tx.Rollback()
			return fmt.Errorf("prepare stmt: %w", err)
		}

		for _, h := range batch {
			_, err := stmt.ExecContext(ctx,
				h.RayonCode, h.Plant, h.RayonType, h.BUM, h.BUMName, h.NSM, h.NSMName,
				h.ASM, h.ASMName, h.FSS, h.FSSName, h.SLM, h.SLMName,
				h.SalesmanCategoryUpdate, h.BranchName, h.MLO, h.Remarks,
				h.TerrCode, h.Username, h.Change, h.CategoryRayon, h.RayonDetail,
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
