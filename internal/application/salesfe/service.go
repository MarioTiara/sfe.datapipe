package salesfe

import (
	"fmt"
	"log"

	"github.com/mariotiara/sfe-data-pipe/internal/application/shared"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/salesfe"
)

type Service struct {
	repo     salesfe.Repository
	streamer shared.StreamDataSource
	mapper   SalesFEMapper
}

func NewService(repo salesfe.Repository, streamer shared.StreamDataSource, mapper SalesFEMapper) *Service {
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

	entities, err := s.mapper.MapRowsToSaleFE(rowsCh)
	if err != nil {
		return fmt.Errorf("failed to map rows: %w", err)
	}

	log.Printf("Total entities collected: %d\n", len(entities))

	return s.saveData(entities)
}

func (s *Service) saveData(data []*salesfe.SalesFE) error {
	hasData, _ := s.repo.HasThisMonthData()
	if hasData {
		s.repo.RemoveThisMonthData()
	}
	return s.repo.SaveRange(data)
}
