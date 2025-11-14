package app

import (
	"database/sql"

	"github.com/mariotiara/sfe-data-pipe/configs"
	"github.com/mariotiara/sfe-data-pipe/internal/application/customerfe"
	"github.com/mariotiara/sfe-data-pipe/internal/application/ezengagecalldetail"
	"github.com/mariotiara/sfe-data-pipe/internal/application/hirarki"
	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/application/master_outlet"
	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/application/material_master"
	"github.com/mariotiara/sfe-data-pipe/internal/application/salesfe"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/excel"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/postgres"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/shared"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

func NewEZEngagePipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *ezengagecalldetail.Service {
	filepath := config.EZEngageCallPath
	repo := postgres.NewEZEngageCallDetailRepository(db, config, logger)
	parser := &excel.ExcelParser{}
	streamer := shared.NewLocalFolderLoader(logger, parser, filepath)
	mapper := excel.NewExcelEZEnggaeCallDetailMapper(logger)
	service := ezengagecalldetail.NewService(repo, streamer, mapper, logger)
	return service
}
func NewCustomerFEPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *customerfe.Service {
	filepath := config.CustomerfePath
	repo := postgres.NewCustomerRepository(db, config, logger)
	parser := &excel.ExcelParser{}
	streamer := shared.NewLocalFolderLoader(logger, parser, filepath)
	mapper := excel.NewExcelCustomerFEMapper(logger)
	service := customerfe.NewService(repo, streamer, mapper, logger)
	return service
}
func NewHirarkiPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *hirarki.Service {
	filepath := config.CustomerfePath
	repo := postgres.NewHirarkiRepository(db, config, logger)
	parser := &excel.ExcelParser{}
	streamer := shared.NewLocalFolderLoader(logger, parser, filepath)
	mapper := excel.NewHirarkiExcelMapper(logger)
	service := hirarki.NewService(repo, streamer, mapper, logger)
	return service
}
func NewSalesFEPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *salesfe.Service {
	filepath := config.SalesFEPath
	repo := postgres.NewSalesRepository(db, config, logger)
	parser := &excel.ExcelParser{}
	streamer := shared.NewLocalFolderLoader(logger, parser, filepath)
	mapper := excel.NewExcelSalesFEMapper(logger)
	service := salesfe.NewService(repo, streamer, mapper, logger)
	return service
}

func NewMaterialMasterPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *materialmaster.Service {
	filepath := config.MasterMaterialPath
	repo := postgres.NewMaterialMasterRepository(db, config, logger)
	parser := &excel.ExcelParser{}
	streamer := shared.NewLocalFolderLoader(logger, parser, filepath)
	mapper := excel.NewExcelMaterialMasterMapper(logger)
	service := materialmaster.NewService(repo, streamer, mapper, logger)
	return service
}

func NewMasterOutletPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *masteroutlet.Service {
	filepath := config.MasterOutletPath
	repo := postgres.NewMasterOutletRepository(db, config, logger)
	parser := &excel.ExcelParser{}
	streamer := shared.NewLocalFolderLoader(logger, parser, filepath)
	mapper := excel.NewExcelMasterOutletMapper(logger)
	service := masteroutlet.NewService(repo, streamer, mapper, logger)
	return service
}
