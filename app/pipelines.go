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

	ports "github.com/mariotiara/sfe-data-pipe/internal/infrastructure/filesources"
	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/postgres"

	"github.com/mariotiara/sfe-data-pipe/internal/infrastructure/tabular"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

func NewEZEngagePipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *ezengagecalldetail.Service {
	filepath := config.EZEngageCallPath
	// archievePath := config.ArchivePath
	repo, _ := postgres.NewEZEngageCallDetailRepository(db, config, logger)
	mapper := tabular.NewExcelEZEnggaeCallDetailMapper(logger)
	fileSources := &ports.LocalFileSource{BaseDir: filepath}
	tabularReader := tabular.New(fileSources)

	service := ezengagecalldetail.NewService(repo, fileSources, tabularReader, mapper, logger)
	return service
}
func NewCustomerFEPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *customerfe.Service {
	filepath := config.CustomerfePath
	// archievePath := config.ArchivePath
	repo, _ := postgres.NewCustomerRepository(db, config, logger)
	mapper := tabular.NewExcelCustomerFEMapper(logger)

	fileSource, _ := ports.NewSFTPFileSource(
		config.SFT_Username,
		config.SFTP_Password,
		config.SFT_Host,
		config.SFT_Port,
		filepath,
	)

	tabularReader := tabular.New(fileSource)

	service := customerfe.NewService(repo, fileSource, tabularReader, mapper, logger)
	return service
}

func NewHirarkiPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *hirarki.Service {
	filepath := config.HirarkiPath
	// archievePath := config.ArchivePath
	repo, _ := postgres.NewHirarkiRepository(db, config, logger)
	mapper := tabular.NewHirarkiExcelMapper(logger)

	fileSources := &ports.LocalFileSource{BaseDir: filepath}
	tabularReader := tabular.New(fileSources)

	service := hirarki.NewService(repo, fileSources, tabularReader, mapper, logger)
	return service
}

func NewSalesFEPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *salesfe.Service {
	filepath := config.SalesFEPath
	// archievePath := config.ArchivePath
	repo, _ := postgres.NewSalesRepository(db, config, logger)
	mapper := tabular.NewExcelSalesFEMapper(logger)
	fileSources := &ports.LocalFileSource{BaseDir: filepath}
	tabularReader := tabular.New(fileSources)

	service := salesfe.NewService(repo, fileSources, tabularReader, mapper, logger)
	return service
}

func NewMaterialMasterPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *materialmaster.Service {
	filepath := config.MasterMaterialPath
	// archievePath := config.ArchivePath
	repo, _ := postgres.NewMaterialMasterRepository(db, config, logger)
	mapper := tabular.NewExcelMaterialMasterMapper(logger)
	fileSources := &ports.LocalFileSource{BaseDir: filepath}
	tabularReader := tabular.New(fileSources)

	service := materialmaster.NewService(repo, fileSources, tabularReader, mapper, logger)
	return service
}

func NewMasterOutletPipeline(db *sql.DB, config *configs.Config, logger logger.Logger) *masteroutlet.Service {
	filepath := config.MasterOutletPath
	// archievePath := config.ArchivePath
	repo, _ := postgres.NewMasterOutletRepository(db, config, logger)
	mapper := tabular.NewExcelMasterOutletMapper(logger)
	fileSources := &ports.LocalFileSource{BaseDir: filepath}
	tabularReader := tabular.New(fileSources)

	service := masteroutlet.NewService(repo, fileSources, tabularReader, mapper, logger)
	return service
}
