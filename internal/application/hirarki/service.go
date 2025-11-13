package hirarki

import (
	"context"
	"fmt"

	"github.com/mariotiara/sfe-data-pipe/internal/application/datastream"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type Service struct {
	logger   logger.Logger
	repo     hirarki.Repository
	streamer datastream.DataStream
	mapper   HirarkiMapper
}

func NewService(repo hirarki.Repository, streamer datastream.DataStream, mapper HirarkiMapper, logger logger.Logger) *Service {
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

	entities, err := s.mapper.MapRowsToHirarki(rowsCh)
	if err != nil {
		s.logger.Error(ctx, "failed to map rows: %v", err)
		return err
	}

	s.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(entities)))

	return s.saveData(ctx, entities)
}

func (s *Service) saveData(ctx context.Context, data []*hirarki.Hirarki) error {
	hasData, _ := s.repo.HasThisMonthData(ctx)
	if hasData {
		s.repo.RemoveThisMonthData(ctx)
	}
	return s.repo.SaveRange(ctx, data)
}
