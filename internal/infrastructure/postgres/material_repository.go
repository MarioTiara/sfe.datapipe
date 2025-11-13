package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mariotiara/sfe-data-pipe/configs"
	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type MaterialMasterRepository struct {
	logger logger.Logger
	db     *sql.DB
	config *configs.Config
}

func NewMaterialMasterRepository(db *sql.DB, config *configs.Config, logger logger.Logger) *MaterialMasterRepository {
	return &MaterialMasterRepository{db: db, config: config, logger: logger}
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
func (r *MaterialMasterRepository) SaveRange(ctx context.Context, materials []*materialmaster.MaterialMaster) error {
	batchSize := r.config.DBBatchSize
	for i := 0; i < len(materials); i += batchSize {
		end := i + batchSize
		if end > len(materials) {
			end = len(materials)
		}

		r.logger.Info(ctx, fmt.Sprintf("insert %d of %d", i, len(materials)))
		batch := materials[i:end]

		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

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
			tx.Rollback()
			return fmt.Errorf("prepare stmt: %w", err)
		}

		for _, m := range batch {
			_, err := stmt.ExecContext(ctx,
				m.PrincipalCode, m.PrincipalName, m.MaterialCode, m.MaterialName,
				m.Uom, m.Division, m.DivisionDescription, m.ProductHierarchy, m.Brand,
				m.MATKL, m.PH2Code, m.PH2Name, m.PH3Code, m.PH3Name,
				m.ATCCode, m.ATCDescription, m.DGIndicatorCode, m.DGIndicatorName,
				m.GeneralDGCode, m.GeneralDGName,
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

		r.logger.Info(ctx, fmt.Sprintf("✅ Committed batch %d–%d successfully\n", i, end))
	}

	return nil
}
