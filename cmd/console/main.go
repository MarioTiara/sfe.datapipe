package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/postgres"
)

func main() {
	fmt.Println("start")

	db, err := postgres.NewPostgressDB()
	if err != nil {
		log.Fatal("failed to start connection to database")
		return
	}

	// ezengage := NewEZEngagePipeline(db)
	// ezengage.Run()

	salesfe := NewSalesFEPipeline(db)
	salesfe.Run()

	// hirarki := NewHirarkiPipeline(db)
	// hirarki.Run()

	// masteroutlet := NewMasterOutletPipeline(db)
	// masteroutlet.Run()

	// materialmaster := NewMaterialMasterPipeline(db)
	// materialmaster.Run()

	// customerfe := NewCustomerFEPipeline(db)
	// customerfe.Run()

}
