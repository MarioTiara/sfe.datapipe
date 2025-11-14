package customerfe

import (
	"context"
	"fmt"

	"github.com/mariotiara/sfe-data-pipe/internal/application/datastream"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/customerfe"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type Service struct {
	logger   logger.Logger
	repo     customerfe.Repository
	streamer datastream.DataStream
	mapper   CustomerFEMapper
}

func NewService(repo customerfe.Repository, streamer datastream.DataStream, mapper CustomerFEMapper, logger logger.Logger) *Service {
	return &Service{repo: repo, streamer: streamer, mapper: mapper, logger: logger}
}

func (s *Service) Run(ctx context.Context) error {
	// Stream rows from the loader
	rowsCh, loaderErrCh := s.streamer.StreamRows(ctx)

	// Collect any loader errors in a separate goroutine
	go func() {
		for err := range loaderErrCh {
			if err != nil {
				s.logger.Error(ctx, "Loader error: %v\n", err)
			}
		}
	}()

	// Map rows to entities (this returns a slice now)
	entities, err := s.mapper.MapRowsToCustomerFE(ctx, rowsCh)
	if err != nil {
		s.logger.Error(ctx, "failed to map rows: %v", err)
		return err
	}

	s.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(entities)))

	// Save entities to repository
	return s.saveData(ctx, entities)
}

func (s *Service) saveData(ctx context.Context, data []*customerfe.CustomerFE) error {
	hasData, _ := s.repo.HasThisMonthData(ctx)
	if hasData {
		s.logger.Info(ctx, "Existing data found for the same month; old records will be removed")
		s.repo.RemoveThisMonthData(ctx)

	}

	return s.repo.SaveRange(ctx, data)
}
