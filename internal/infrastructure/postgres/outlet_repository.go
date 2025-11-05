package postgres

import (
	"context"
	"database/sql"

	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"
)

type MasterOutletRepository struct {
	db *sql.DB
}

func NewMasterOutletRepository(db *sql.DB) *MasterOutletRepository {
	return &MasterOutletRepository{db: db}
}

func (r *MasterOutletRepository) Save(m masteroutlet.MasterOutlet) error {
	ctx := context.Background()
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

func (r *MasterOutletRepository) SaveRange(ms []masteroutlet.MasterOutlet) error {
	ctx := context.Background()
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

	for _, m := range ms {
		_, err := stmt.ExecContext(ctx,
			m.CustomerCode, m.CustomerName, m.Channel, m.Plant, m.BranchName,
			m.NIKSalesman, m.NameSalesman, m.RayonCode, m.Rayon, m.NewClass,
			m.CallPlanFullMonth, m.TargetFreq, m.TerrCode, m.Username,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
