package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mariotiara/sfe-data-pipe/configs"
)

func NewPostgressDB(config *configs.Config) (*sql.DB, error) {
	constr := config.ConnString
	if constr == "" {
		return nil, fmt.Errorf("connetion string is empty")
	}

	db, err := sql.Open("postgres", constr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
