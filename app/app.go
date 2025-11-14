package app

import (
	"context"

	"github.com/mariotiara/sfe-data-pipe/configs"
	"github.com/mariotiara/sfe-data-pipe/internal/application/customerfe"
	"github.com/mariotiara/sfe-data-pipe/internal/application/ezengagecalldetail"
	"github.com/mariotiara/sfe-data-pipe/internal/application/hirarki"
	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/application/master_outlet"
	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/application/material_master"
	"github.com/mariotiara/sfe-data-pipe/internal/application/salesfe"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/logger"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/postgres"
)

type App struct {
	EZEngagePipeLine     *ezengagecalldetail.Service
	CustomerFEPipeline   *customerfe.Service
	HirarkiPipeLine      *hirarki.Service
	MasterOutletPipeline *masteroutlet.Service
	SalesFEPipeLine      *salesfe.Service
	MaterialPipeLine     *materialmaster.Service
}

func NewApp(ctx context.Context, config *configs.Config) (*App, error) {
	db, err := postgres.NewPostgressDB()
	if err != nil {
		return nil, err
	}

	log := logger.NewZapLogger()

	return &App{
		EZEngagePipeLine:     NewEZEngagePipeline(db, config, log),
		CustomerFEPipeline:   NewCustomerFEPipeline(db, config, log),
		HirarkiPipeLine:      NewHirarkiPipeline(db, config, log),
		MasterOutletPipeline: NewMasterOutletPipeline(db, config, log),
		SalesFEPipeLine:      NewSalesFEPipeline(db, config, log),
	}, nil
}
