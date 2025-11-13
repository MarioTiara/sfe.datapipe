package materialmaster

import (
	"context"
	"fmt"
	"log"

	"github.com/mariotiara/sfe-data-pipe/internal/application/datastream"
	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type Service struct {
	logger   logger.Logger
	repo     materialmaster.Repository
	streamer datastream.DataStream
	mapper   MaterialMasterMapper
}

func NewService(repo materialmaster.Repository, streamer datastream.DataStream, mapper MaterialMasterMapper, logger logger.Logger) *Service {
	return &Service{repo: repo, streamer: streamer, mapper: mapper}
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

	entities, err := s.mapper.MapRowsToMaterialMaster(rowsCh)
	if err != nil {
		s.logger.Error(ctx, "failed to map rows: %v", err)
		return err
	}

	log.Printf("Total entities collected: %d\n", len(entities))
	s.logger.Info(ctx, fmt.Sprintf("Total entities collected: %d\n", len(entities)))
	return s.saveData(ctx, entities)
}

func (s *Service) saveData(ctx context.Context, data []*materialmaster.MaterialMaster) error {
	hasData, _ := s.repo.HasThisMonthData(ctx)
	if hasData {
		s.repo.RemoveThisMonthData(ctx)
	}
	return s.repo.SaveRange(ctx, data)
}
