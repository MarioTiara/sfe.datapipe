package hirarki

import (
	"context"
	"fmt"
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/application/ports"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type Service struct {
	logger        logger.Logger
	repo          hirarki.Repository
	fileSources   ports.FileSource
	mapper        HirarkiMapper
	tabularReader ports.TabularFileReader
}

func NewService(repo hirarki.Repository, fileSources ports.FileSource, tabularReader ports.TabularFileReader, mapper HirarkiMapper, logger logger.Logger) *Service {
	return &Service{repo: repo, fileSources: fileSources, tabularReader: tabularReader, mapper: mapper, logger: logger}
}

func (s *Service) Run(ctx context.Context) error {
	start := time.Now()
	files, err := s.fileSources.List("Hirarki")
	if err != nil {
		s.logger.Error(ctx, "Error in loading file sources", err)
		return err
	}

	s.logger.Info(ctx, fmt.Sprintf("Found %d files", len(files)))
	fileCount := 0
	rowCount := 0
	for _, f := range files {
		fstart := time.Now()
		s.logger.Info(
			ctx, "Processing file",
			logger.Field{Key: "file_name", Value: f.Name},
		)

		rowsCh, err := s.tabularReader.ReadRows(f.Name)
		if err != nil {
			s.logger.Error(ctx, "failed reads tabular data: %v", err)
			continue
		}
		entities, err := s.mapper.MapRowsToHirarki(ctx, rowsCh)
		if err != nil {
			s.logger.Error(ctx, "failed to map rows: %v", err)
			continue
		}
		s.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(entities)))
		err = s.saveData(ctx, entities)
		if err != nil {
			s.logger.Error(ctx, "Hirarki Pipeline Failed to insert to database", err,
				logger.Field{Key: "process_time", Value: time.Since(start).String()},
			)
			continue
		}

		s.fileSources.Move(f.Name, "archieve")
		s.logger.Info(
			ctx, "Processing file done",
			logger.Field{Key: "file_name", Value: f.Name},
			logger.Field{Key: "process_time", Value: time.Since(fstart).String()},
		)

		fileCount += 1
		rowCount += len(entities)
	}

	s.logger.Info(ctx, "Pipeline Process Done",
		logger.Field{Key: "process_time", Value: time.Since(start).String()},
		logger.Field{Key: "total_file_processed", Value: fileCount},
		logger.Field{Key: "total_rows_processed", Value: rowCount},
	)
	return nil
}

func (s *Service) saveData(ctx context.Context, data []*hirarki.Hirarki) error {
	hasData, _ := s.repo.HasThisMonthData(ctx)
	if hasData {
		s.logger.Info(ctx, "Existing data found for the same month; old records will be removed")
		s.repo.RemoveThisMonthData(ctx)
	}

	return s.repo.SaveRange(ctx, data)
}
