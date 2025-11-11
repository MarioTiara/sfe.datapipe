package ezengagecalldetail

import (
	"fmt"
	"log"

	"github.com/mariotiara/sfe-data-pipe/internal/application/shared"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/ezengagecalldetail"
)

type Service struct {
	repo     ezengagecalldetail.Repository
	streamer shared.StreamDataSource
	mapper   CallDetailMapper
}

func NewService(repo ezengagecalldetail.Repository, streamer shared.StreamDataSource, mapper CallDetailMapper) *Service {
	return &Service{repo: repo, streamer: streamer, mapper: mapper}
}

func (s *Service) Run() error {
	// Stream rows from the loader
	rowsCh, loaderErrCh := s.streamer.StreamRows()

	// Collect any loader errors in a separate goroutine
	go func() {
		for err := range loaderErrCh {
			if err != nil {
				log.Printf("Loader error: %v\n", err)
			}
		}
	}()

	// Map rows to entities (this returns a slice now)
	entities, err := s.mapper.MapRowsToCallDetails(rowsCh)
	if err != nil {
		return fmt.Errorf("failed to map rows: %w", err)
	}

	log.Printf("Total entities collected: %d\n", len(entities))

	// Save entities to repository
	return s.saveData(entities)
}

func (s *Service) saveData(data []*ezengagecalldetail.EZEngageCallDetail) error {
	hasData, _ := s.repo.HasThisMonthData()
	if hasData {
		row, _ := s.repo.RemoveThisMonthData()
		ms := fmt.Sprintf("%d is removed", row)
		log.Println(ms)
	}

	return s.repo.SaveRange(data)
}
