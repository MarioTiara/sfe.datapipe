package ezengagecalldetail

import (
	"context"
	"fmt"
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/application/datastream"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/ezengagecalldetail"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type Service struct {
	logger   logger.Logger
	repo     ezengagecalldetail.Repository
	streamer datastream.DataStream
	mapper   CallDetailMapper
}

func NewService(repo ezengagecalldetail.Repository, streamer datastream.DataStream, mapper CallDetailMapper, logger logger.Logger) *Service {
	return &Service{repo: repo, streamer: streamer, mapper: mapper, logger: logger}
}

func (s *Service) Run(ctx context.Context) error {
	// Stream rows from the loader
	strat := time.Now()
	s.logger.Info(ctx, "EZEngage Pipeline Started")

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
	entities, err := s.mapper.MapRowsToCallDetails(ctx, rowsCh)
	if err != nil {
		s.logger.Error(ctx, "failed to map rows: %v", err)
		return err
	}

	s.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(entities)))

	// Save entities to repository
	err = s.saveData(ctx, entities)
	if err != nil {
		s.logger.Error(ctx, "EZEngage Failed", err,
			logger.Field{Key: "process_time", Value: time.Since(strat).String()},
		)
		return err
	}

	s.logger.Info(ctx, "EZEngage Process Done", logger.Field{Key: "process_time", Value: time.Since(strat).String()})
	return nil

}

func (s *Service) saveData(ctx context.Context, data []*ezengagecalldetail.EZEngageCallDetail) error {
	hasData, _ := s.repo.HasThisMonthData(ctx)
	if hasData {
		s.logger.Info(ctx, "Existing data found for the same month; old records will be removed")
		s.repo.RemoveThisMonthData(ctx)
	}

	return s.repo.SaveRange(ctx, data)
}
