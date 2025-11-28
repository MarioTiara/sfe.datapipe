package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mariotiara/sfe-data-pipe/configs"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/customerfe"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type CustomerRepository struct {
	logger logger.Logger
	db     *sql.DB
	config *configs.Config
}

func NewCustomerRepository(db *sql.DB, config *configs.Config, logger logger.Logger) *CustomerRepository {
	return &CustomerRepository{db: db, config: config, logger: logger}
}

func (r *CustomerRepository) RemoveThisMonthData(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		DELETE FROM customer_fe
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

func (r *CustomerRepository) HasThisMonthData(ctx context.Context) (bool, error) {
	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		SELECT COUNT(1)
		FROM customer_fe
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

// Save inserts a single customer record
func (r *CustomerRepository) Save(ctx context.Context, customer *customerfe.CustomerFE) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO customer_fe (
			customer_code, customer_name, address, city_name, bank_country,
			regio, trans_zone, telephone, bran1, bran2, bran3,
			channel_ic4, vtext, katr1, katr2, katr3, adrnr, vkorg,
			sales_office, sales_district, customer_grp1, cust_grp1_desc,
			customer_grp2, shipping_condition, lprio, branch, eikto,
			ktokd, flag_deletion, sperr1, aufsd1, lifsd1, faksd1,
			cassd1, erdat, ernam, postal_code
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22,
			$23, $24, $25, $26, $27,
			$28, $29, $30, $31, $32, $33,
			$34, $35, $36, $37
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		customer.CustomerCode, customer.CustomerName, customer.Address, customer.CityName, customer.BankCountry,
		customer.Regio, customer.TransZone, customer.Telephone, customer.Bran1, customer.Bran2, customer.Bran3,
		customer.ChannelIc4, customer.Vtext, customer.Katr1, customer.Katr2, customer.Katr3, customer.Adrnr, customer.Vkorg,
		customer.SalesOffice, customer.SalesDistrict, customer.CustomerGrp1, customer.CustGrp1Desc,
		customer.CustomerGrp2, customer.ShippingCondition, customer.Lprio, customer.Branch, customer.Eikto,
		customer.Ktokd, customer.FlagDeletion, customer.Sperr1, customer.Aufsd1, customer.Lifsd1, customer.Faksd1,
		customer.Cassd1, customer.Erdat, customer.Ernam, customer.PostalCode,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SaveRange inserts multiple customer records efficiently in a single transaction
func (r *CustomerRepository) SaveRange(ctx context.Context, customers []*customerfe.CustomerFE) error {
	batchSize := r.config.DBBatchSize
	tinserted := 0
	start := time.Now()
	for i := 0; i < len(customers); i += batchSize {
		end := i + batchSize
		if end > len(customers) {
			end = len(customers)
		}
		batch := customers[i:end]

		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

		stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO customer_fe (
			customer_code, customer_name, address, city_name, bank_country,
			regio, trans_zone, telephone, bran1, bran2, bran3,
			channel_ic4, vtext, katr1, katr2, katr3, adrnr, vkorg,
			sales_office, sales_district, customer_grp1, cust_grp1_desc,
			customer_grp2, shipping_condition, lprio, branch, eikto,
			ktokd, flag_deletion, sperr1, aufsd1, lifsd1, faksd1,
			cassd1, erdat, ernam, postal_code
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22,
			$23, $24, $25, $26, $27,
			$28, $29, $30, $31, $32, $33,
			$34, $35, $36, $37
		)
	`)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("prepare stmt: %w", err)
		}

		for _, c := range batch {
			_, err := stmt.ExecContext(ctx,
				c.CustomerCode, c.CustomerName, c.Address, c.CityName, c.BankCountry,
				c.Regio, c.TransZone, c.Telephone, c.Bran1, c.Bran2, c.Bran3,
				c.ChannelIc4, c.Vtext, c.Katr1, c.Katr2, c.Katr3, c.Adrnr, c.Vkorg,
				c.SalesOffice, c.SalesDistrict, c.CustomerGrp1, c.CustGrp1Desc,
				c.CustomerGrp2, c.ShippingCondition, c.Lprio, c.Branch, c.Eikto,
				c.Ktokd, c.FlagDeletion, c.Sperr1, c.Aufsd1, c.Lifsd1, c.Faksd1,
				c.Cassd1, c.Erdat, c.Ernam, c.PostalCode,
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

		tinserted = end
	}
	r.logger.Info(ctx, fmt.Sprintf("%d rows successfully inserted", tinserted),
		logger.Field{Key: "process_time", Value: time.Since(start).String()})
	return nil
}
