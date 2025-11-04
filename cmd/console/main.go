package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/application/material_master"
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
	filePath := "data/202509_SFE_Tablemaster_Material.xlsx"
	defer db.Close()

	excelData, err := excel.ReadMaterialMasterFromExcel(filePath)
	if err != nil {
		log.Fatal("error reading excel:", err)
	}

	for _, data := range excelData {
		fmt.Printf("%+v\n", data)
	}
	repo := postgres.NewMaterialMasterRepository(db)
	service := materialmaster.NewService(repo)

	if err := service.ImportMaterials(excelData); err != nil {
		log.Fatal("error importing data:", err)
	}

	fmt.Println("Data importted successfully")
}
