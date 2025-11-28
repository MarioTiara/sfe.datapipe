package main

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/mariotiara/sfe-data-pipe/app"
	"github.com/mariotiara/sfe-data-pipe/configs"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/logger"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/processid"
)

func main() {

	log := logger.NewZapLogger()
	ctx := context.Background()
	config, err := configs.Load()
	if err != nil {
		log.Error(ctx, "Failed to load config %v", err)
		return
	}

	app, err := app.NewApp(ctx, config)
	if err != nil {
		log.Error(ctx, "Failed to start app %v", err)
	}

	ezctx := processid.WithContext(ctx)
	app.EZEngagePipeLine.Run(ezctx)

	salesctx := processid.WithContext(ctx)
	app.SalesFEPipeLine.Run(salesctx)

	hctx := processid.WithContext(ctx)
	app.HirarkiPipeLine.Run(hctx)

	moctx := processid.WithContext(ctx)
	app.MasterOutletPipeline.Run(moctx)

	mmctx := processid.WithContext(ctx)
	app.MaterialPipeLine.Run(mmctx)

	csctx := processid.WithContext(ctx)
	app.CustomerFEPipeline.Run(csctx)
}
