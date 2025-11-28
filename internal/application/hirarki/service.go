package hirarki

import (
	"context"
	"time"

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
	start := time.Now()
	s.logger.Info(ctx, "Hirarki Pipeline Started")
	rowsCh, loaderErrCh := s.streamer.StreamRows(ctx)
	go func() {
		for err := range loaderErrCh {
			if err != nil {
				s.logger.Error(ctx, "Loader error: %v\n", err)
			}
		}
	}()

	entities, err := s.mapper.MapRowsToHirarki(ctx, rowsCh)
	if err != nil {
		s.logger.Error(ctx, "failed to map rows: %v", err)
		return err
	}

	err = s.saveData(ctx, entities)
	if err != nil {
		s.logger.Error(ctx, "Hirarki Pipeline Failed", err,
			logger.Field{Key: "process_time", Value: time.Since(start).String()},
		)
		return err
	}

	s.logger.Info(ctx, "Hirarki Pipeline Process Done", logger.Field{Key: "process_time", Value: time.Since(start).String()})
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
