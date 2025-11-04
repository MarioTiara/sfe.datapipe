package postgres

import (
	"context"
	"database/sql"

	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"
)

type MaterialMasterRepository struct {
	db *sql.DB
}

func NewMaterialMasterRepository(db *sql.DB) *MaterialMasterRepository {
	return &MaterialMasterRepository{db: db}
}

// Save inserts a single MaterialMaster record
func (r *MaterialMasterRepository) Save(m materialmaster.MaterialMaster) error {
	ctx := context.Background()
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
func (r *MaterialMasterRepository) SaveRange(materials []materialmaster.MaterialMaster) error {
	ctx := context.Background()
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

	for _, m := range materials {
		_, err := stmt.ExecContext(ctx,
			m.PrincipalCode, m.PrincipalName, m.MaterialCode, m.MaterialName,
			m.Uom, m.Division, m.DivisionDescription, m.ProductHierarchy, m.Brand,
			m.MATKL, m.PH2Code, m.PH2Name, m.PH3Code, m.PH3Name,
			m.ATCCode, m.ATCDescription, m.DGIndicatorCode, m.DGIndicatorName,
			m.GeneralDGCode, m.GeneralDGName,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
