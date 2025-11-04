package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/mariotiara/sfe-data-pipe/internal/application/customerfe"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/excel"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/postgres"
)

func main() {
	fmt.Println("start")
	connStr := "postgres://postgres:secret@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	filePath := "data/202509_SFE_CustomersList.XLSX"
	defer db.Close()

	excelData, err := excel.ReadFeCustomerFromExcel(filePath)
	if err != nil {
		log.Fatal("error reading excel:", err)
	}

	for _, data := range excelData {
		fmt.Printf("%+v\n", data)
	}
	repo := postgres.NewCustomerRepository(db)
	service := customerfe.NewService(repo)

	if err := service.ImportCustomers(excelData); err != nil {
		log.Fatal("error importing data:", err)
	}

	fmt.Println("Data importted successfully")
}
