package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/mariotiara/sfe-data-pipe/internal/application/ezengagecalldetail"
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
	filePath := "data/202509_Call Detailed eZEngage.xlsx"
	defer db.Close()

	excelData, err := excel.ReadEZEngageCallDetailFromExcel(filePath)
	if err != nil {
		log.Fatal("error reading excel:", err)
	}

	// for _, data := range excelData {
	// 	fmt.Printf("%+v\n", data)
	// }

	fmt.Println(len((excelData)))
	repo := postgres.NewEZEngageCallDetailRepository(db)
	service := ezengagecalldetail.NewService(repo)

	if err := service.ImportCallDetails(excelData); err != nil {
		log.Fatal("error importing data:", err)
	}

	fmt.Println("Data importted successfully")
}
