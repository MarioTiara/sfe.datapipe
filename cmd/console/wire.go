package main

import (
	"database/sql"

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

func NewEZEngagePipeline(db *sql.DB) *ezengagecalldetail.Service {
	repo := postgres.NewEZEngageCallDetailRepository(db)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: "data/ezengagecalldetail"}
	mapper := excel.NewExcelEZEnggaeCallDetailMapper()
	service := ezengagecalldetail.NewService(repo, streamer, mapper)
	return service
}
func NewCustomerFEPipeline(db *sql.DB) *customerfe.Service {
	repo := postgres.NewCustomerRepository(db)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: "data/ezengagecalldetail"}
	mapper := excel.NewExcelCustomerFEMapper()
	service := customerfe.NewService(repo, streamer, mapper)
	return service
}
func NewHirarkiPipeline(db *sql.DB) *hirarki.Service {
	repo := postgres.NewHirarkiRepository(db)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: "data/ezengagecalldetail"}
	mapper := excel.NewHirarkiExcelMapper()
	service := hirarki.NewService(repo, streamer, mapper)
	return service
}
func NewSalesFEPipeline(db *sql.DB) *salesfe.Service {
	repo := postgres.NewSalesRepository(db)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: "data/ezengagecalldetail"}
	mapper := excel.NewExcelSalesFEMapper()
	service := salesfe.NewService(repo, streamer, mapper)
	return service
}

func NewMaterialMasterPipeline(db *sql.DB) *materialmaster.Service {
	repo := postgres.NewMaterialMasterRepository(db)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: "data/ezengagecalldetail"}
	mapper := excel.NewExcelMaterialMasterMapper()
	service := materialmaster.NewService(repo, streamer, mapper)
	return service
}

func NewMasterOutletPipeline(db *sql.DB) *masteroutlet.Service {
	repo := postgres.NewMasterOutletRepository(db)
	parser := &excel.ExcelParser{}
	streamer := &shared.LocalFolderLoader{Parser: parser, Path: "data/ezengagecalldetail"}
	mapper := excel.NewExcelMasterOutletMapper()
	service := masteroutlet.NewService(repo, streamer, mapper)
	return service
}
