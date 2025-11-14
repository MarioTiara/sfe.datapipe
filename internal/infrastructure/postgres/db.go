package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

func NewPostgressDB() (*sql.DB, error) {
	// config, _ := configs.Load()

	constr := config.ConnString
	// constr := "postgres://postgres:secret@localhost:5432/mydb?sslmode=disable"
	fmt.Println(constr)
	db, err := sql.Open("postgres", constr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
