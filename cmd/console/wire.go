package main

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
)

var config, _ = configs.Load()

func NewEZEngagePipeline(db *sql.DB) *ezengagecalldetail.Service {
	filepath := config.EZEngageCallPath
	repo := postgres.NewEZEngageCallDetailRepository(db, config)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: filepath}
	mapper := excel.NewExcelEZEnggaeCallDetailMapper()
	service := ezengagecalldetail.NewService(repo, streamer, mapper)
	return service
}
func NewCustomerFEPipeline(db *sql.DB) *customerfe.Service {
	filepath := config.CustomerfePath
	repo := postgres.NewCustomerRepository(db, config)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: filepath}
	mapper := excel.NewExcelCustomerFEMapper()
	service := customerfe.NewService(repo, streamer, mapper)
	return service
}
func NewHirarkiPipeline(db *sql.DB) *hirarki.Service {
	filepath := config.CustomerfePath
	repo := postgres.NewHirarkiRepository(db, config)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: filepath}
	mapper := excel.NewHirarkiExcelMapper()
	service := hirarki.NewService(repo, streamer, mapper)
	return service
}
func NewSalesFEPipeline(db *sql.DB) *salesfe.Service {
	filepath := config.SalesFEPath
	repo := postgres.NewSalesRepository(db, config)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: filepath}
	mapper := excel.NewExcelSalesFEMapper()
	service := salesfe.NewService(repo, streamer, mapper)
	return service
}

func NewMaterialMasterPipeline(db *sql.DB) *materialmaster.Service {
	filepath := config.MasterMaterialPath
	repo := postgres.NewMaterialMasterRepository(db, config)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: filepath}
	mapper := excel.NewExcelMaterialMasterMapper()
	service := materialmaster.NewService(repo, streamer, mapper)
	return service
}

func NewMasterOutletPipeline(db *sql.DB) *masteroutlet.Service {
	filepath := config.MasterOutletPath
	repo := postgres.NewMasterOutletRepository(db, config)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: filepath}
	mapper := excel.NewExcelMasterOutletMapper()
	service := masteroutlet.NewService(repo, streamer, mapper)
	return service
}
