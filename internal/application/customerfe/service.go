package customerfe

import (
	"context"
	"fmt"
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/application/ports"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/customerfe"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/fileingestion"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type Service struct {
	logger        logger.Logger
	repo          customerfe.Repository
	fileSources   ports.FileSource
	mapper        CustomerFEMapper
	tabularReader ports.TabularFileReader
	policy        fileingestion.RegexFileNamePolicy
	format        fileingestion.FileNameFormat
}

func NewService(repo customerfe.Repository,
	fileSources ports.FileSource,
	tabularReader ports.TabularFileReader,
	mapper CustomerFEMapper,
	policy fileingestion.RegexFileNamePolicy,
	format fileingestion.FileNameFormat,
	logger logger.Logger) *Service {
	return &Service{repo: repo,
		fileSources:   fileSources,
		tabularReader: tabularReader,
		mapper:        mapper,
		policy:        policy,
		format:        format,
		logger:        logger}
}

func (s *Service) Run(ctx context.Context) error {

	start := time.Now()
	files, err := s.fileSources.List("SFE_CustomersList")
	if err != nil {
		return err
	}

	s.logger.Info(ctx, fmt.Sprintf("SFE_CustomersList Pipeline found %d files", len(files)))
	fileCount := 0
	rowCount := 0

	for _, f := range files {
		fstart := time.Now()
		s.logger.Info(
			ctx, "Processing file",
			logger.Field{Key: "file_name", Value: f.Name},
		)

		var match, err = s.policy.Match(f.Name)
		if err != nil || !match {
			s.logger.Error(ctx, "file name is not match regex policy", err)
			continue
		}

		// ingestion, err := s.format.Parse(f.Name)
		// if err != nil {
		// 	s.logger.Error(ctx, "invalid file name", err)
		// 	continue
		// }

		rowsCh, err := s.tabularReader.ReadRows(f.Name)
		if err != nil {
			s.logger.Error(ctx, "failed reads tabular data: %v", err)
			continue
		}
		entities, err := s.mapper.MapRowsToCustomerFE(ctx, rowsCh)
		if err != nil {
			s.logger.Error(ctx, "failed to map rows: %v", err)
			continue
		}

		s.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(entities)))
		err = s.saveData(ctx, entities)

		if err != nil {
			s.logger.Error(ctx, "Failed to save data into database", err,
				logger.Field{Key: "process_time", Value: time.Since(start).String()})
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

	s.logger.Info(ctx, "SFE_CustomersList Pipeline Process Done",
		logger.Field{Key: "process_time", Value: time.Since(start).String()},
		logger.Field{Key: "total_file_processed", Value: fileCount},
		logger.Field{Key: "total_rows_processed", Value: rowCount},
	)
	return nil
}

func (s *Service) saveData(ctx context.Context, data []*customerfe.CustomerFE) error {

	return s.repo.SaveRange(ctx, data)
}

// func (s *)
