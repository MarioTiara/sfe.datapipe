package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mariotiara/sfe-data-pipe/configs"
	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"
)

type MasterOutletRepository struct {
	db     *sql.DB
	config *configs.Config
}

func NewMasterOutletRepository(db *sql.DB, config *configs.Config) *MasterOutletRepository {
	return &MasterOutletRepository{db: db, config: config}
}

func (r *MasterOutletRepository) RemoveThisMonthData() (int, error) {
	ctx := context.Background()

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

	return count, nil
}

func (r *MasterOutletRepository) HasThisMonthData() (bool, error) {
	ctx := context.Background()

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

func (r *MasterOutletRepository) Save(m *masteroutlet.MasterOutlet) error {
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

func (r *MasterOutletRepository) SaveRange(ms []*masteroutlet.MasterOutlet) error {
	batchSize := r.config.DBBatchSize

	ctx := context.Background()
	for i := 0; i < len(ms); i += batchSize {
		end := i + batchSize
		if end > len(ms) {
			end = len(ms)
		}

		msg := fmt.Sprintf("insert %d of %d", i, len(ms))
		log.Println(msg)
		batch := ms[i:end]

		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

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
			tx.Rollback()
			return fmt.Errorf("prepare stmt: %w", err)
		}

		for _, m := range batch {
			_, err := stmt.ExecContext(ctx,
				m.CustomerCode, m.CustomerName, m.Channel, m.Plant, m.BranchName,
				m.NIKSalesman, m.NameSalesman, m.RayonCode, m.Rayon, m.NewClass,
				m.CallPlanFullMonth, m.TargetFreq, m.TerrCode, m.Username,
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
