package main

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/logger"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/postgres"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/processid"
)

func main() {
	fmt.Println("start")

	db, err := postgres.NewPostgressDB()
	if err != nil {
		// log.Fatal("failed to start connection to database")
		return
	}
	exLogger := logger.NewZapLogger()
	ctx := processid.WithContext(context.Background())
	ezengage := NewEZEngagePipeline(ctx, db, exLogger)
	ezengage.Run(ctx)

	// salesfe := NewSalesFEPipeline(db)
	// salesfe.Run()

	// hirarki := NewHirarkiPipeline(db)
	// hirarki.Run()

	// masteroutlet := NewMasterOutletPipeline(db)
	// masteroutlet.Run()

	// materialmaster := NewMaterialMasterPipeline(db)
	// materialmaster.Run()

	// customerfe := NewCustomerFEPipeline(db)
	// customerfe.Run()

}
