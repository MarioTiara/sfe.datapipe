package masteroutlet

import (
	"context"
	"fmt"
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/application/ports"
	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type Service struct {
	logger        logger.Logger
	repo          masteroutlet.Repository
	fileSources   ports.FileSource
	mapper        MasterOutletMapper
	tabularReader ports.TabularFileReader
}

func NewService(repo masteroutlet.Repository, fileSources ports.FileSource, tabularReader ports.TabularFileReader, mapper MasterOutletMapper, logger logger.Logger) *Service {
	return &Service{repo: repo, fileSources: fileSources, tabularReader: tabularReader, mapper: mapper, logger: logger}
}

func (s *Service) Run(ctx context.Context) error {
	start := time.Now()
	// List files to process
	files, err := s.fileSources.List("Masterlist Outlet")
	if err != nil {
		s.logger.Error(ctx, "Error in loading file sources", err)
		return err
	}

	s.logger.Info(ctx, fmt.Sprintf("Found %d files", len(files)))
	fileCount := 0
	rowCount := 0

	for _, f := range files {
		fstart := time.Now()
		s.logger.Info(ctx, "Processing file",
			logger.Field{Key: "file_name", Value: f.Name},
		)

		// Stream rows from the file
		rowsCh, err := s.tabularReader.ReadRows(f.Name)
		// Map rows to entities
		entities, err := s.mapper.MapRowsToMasterOutlet(ctx, rowsCh)
		if err != nil {
			s.logger.Error(ctx, "Failed to map rows", err,
				logger.Field{Key: "file_name", Value: f.Name},
			)
			continue
		}

		rowCount += len(entities)
		// Save entities to database
		if err := s.saveData(ctx, entities); err != nil {
			s.logger.Error(ctx, "Pipeline failed to insert to database", err,
				logger.Field{Key: "file_name", Value: f.Name},
				logger.Field{Key: "process_time", Value: time.Since(fstart).String()},
			)
			continue
		}

		// Move file to archive
		if err := s.fileSources.Move(f.Name, "archive"); err != nil {
			s.logger.Error(ctx, "Failed to move file to archive", err,
				logger.Field{Key: "file_name", Value: f.Name},
			)
		}

		s.logger.Info(ctx, "Processing file done",
			logger.Field{Key: "file_name", Value: f.Name},
			logger.Field{Key: "process_time", Value: time.Since(fstart).String()},
			logger.Field{Key: "rows_processed", Value: len(entities)},
		)
		fileCount++
		rowCount += len(entities)
	}

	s.logger.Info(ctx, "Pipeline Process Done",
		logger.Field{Key: "process_time", Value: time.Since(start).String()},
		logger.Field{Key: "total_file_processed", Value: fileCount},
		logger.Field{Key: "total_rows_processed", Value: rowCount},
	)

	return nil
}

func (s *Service) saveData(ctx context.Context, data []*masteroutlet.MasterOutlet) error {
	hasData, _ := s.repo.HasThisMonthData(ctx)
	if hasData {
		s.logger.Info(ctx, "Existing data found for the same month; old records will be removed")
		s.repo.RemoveThisMonthData(ctx)
	}
	return s.repo.SaveRange(ctx, data)
}
