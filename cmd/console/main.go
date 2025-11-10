package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("start")
	connStr := "postgres://postgres:secret@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	ezengage := NewEZEngagePipeline(db)
	ezengage.Run()

	salesfe := NewSalesFEPipeline(db)
	salesfe.Run()

	hirarki := NewHirarkiPipeline(db)
	hirarki.Run()

	masteroutlet := NewMasterOutletPipeline(db)
	masteroutlet.Run()

	materialmaster := NewMaterialMasterPipeline(db)
	materialmaster.Run()

	customerfe := NewCustomerFEPipeline(db)
	customerfe.Run()

}
