package postgres

import (
	"context"
	"database/sql"

	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
)

type HirarkiRepository struct {
	db *sql.DB
}

func NewHirarkiRepository(db *sql.DB) *HirarkiRepository {
	return &HirarkiRepository{db: db}
}

func (r *HirarkiRepository) Save(h hirarki.Hirarki) error {
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

func (r *HirarkiRepository) SaveRange(hs []hirarki.Hirarki) error {
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

	for _, h := range hs {
		_, err := stmt.ExecContext(ctx,
			h.RayonCode, h.Plant, h.RayonType, h.BUM, h.BUMName, h.NSM, h.NSMName,
			h.ASM, h.ASMName, h.FSS, h.FSSName, h.SLM, h.SLMName,
			h.SalesmanCategoryUpdate, h.BranchName, h.MLO, h.Remarks,
			h.TerrCode, h.Username, h.Change, h.CategoryRayon, h.RayonDetail,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
