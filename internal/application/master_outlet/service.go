package masteroutlet

import (
	"context"
	"fmt"

	"github.com/mariotiara/sfe-data-pipe/internal/application/datastream"
	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type Service struct {
	logger   logger.Logger
	repo     masteroutlet.Repository
	streamer datastream.DataStream
	mapper   MasterOutletMapper
}

func NewService(repo masteroutlet.Repository, streamer datastream.DataStream, mapper MasterOutletMapper, logger logger.Logger) *Service {
	return &Service{repo: repo, streamer: streamer, mapper: mapper, logger: logger}
}

func (s *Service) Run(ctx context.Context) error {
	rowsCh, loaderErrCh := s.streamer.StreamRows(ctx)
	go func() {
		for err := range loaderErrCh {
			if err != nil {
				s.logger.Error(ctx, "Loader error: %v\n", err)
			}
		}
	}()

	entities, err := s.mapper.MapRowsToMasterOutlet(ctx, rowsCh)
	if err != nil {
		s.logger.Error(ctx, "failed to map rows: %v", err)
		return err
	}

	s.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(entities)))
	return s.saveData(ctx, entities)
}

func (s *Service) saveData(ctx context.Context, data []*masteroutlet.MasterOutlet) error {
	hasData, _ := s.repo.HasThisMonthData(ctx)
	if hasData {
		s.logger.Info(ctx, "Existing data found for the same month; old records will be removed")
		s.repo.RemoveThisMonthData(ctx)
	}
	return s.repo.SaveRange(ctx, data)
}
