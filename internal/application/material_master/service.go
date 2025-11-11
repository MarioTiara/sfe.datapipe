package materialmaster

import (
	"fmt"
	"log"

	"github.com/mariotiara/sfe-data-pipe/internal/application/shared"
	materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"
)

type Service struct {
	repo     materialmaster.Repository
	streamer shared.StreamDataSource
	mapper   MaterialMasterMapper
}

func NewService(repo materialmaster.Repository, streamer shared.StreamDataSource, mapper MaterialMasterMapper) *Service {
	return &Service{repo: repo, streamer: streamer, mapper: mapper}
}

func (s *Service) Run() error {
	rowsCh, loaderErrCh := s.streamer.StreamRows()
	go func() {
		for err := range loaderErrCh {
			if err != nil {
				log.Printf("Loader error: %v\n", err)
			}
		}
	}()

	entities, err := s.mapper.MapRowsToMaterialMaster(rowsCh)
	if err != nil {
		return fmt.Errorf("failed to map rows: %w", err)
	}

	log.Printf("Total entities collected: %d\n", len(entities))

	return s.saveData(entities)
}

func (s *Service) saveData(data []*materialmaster.MaterialMaster) error {
	hasData, _ := s.repo.HasThisMonthData()
	if hasData {
		s.repo.RemoveThisMonthData()
	}
	return s.repo.SaveRange(data)
}
