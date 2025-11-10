package hirarki

import (
	"fmt"

	"github.com/mariotiara/sfe-data-pipe/internal/application/shared"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"
)

type Service struct {
	repo     hirarki.Repository
	streamer shared.StreamDataSource
	mapper   HirarkiMapper
}

func NewService(repo hirarki.Repository, streamer shared.StreamDataSource, mapper HirarkiMapper) *Service {
	return &Service{repo: repo, streamer: streamer, mapper: mapper}
}

func (s *Service) Run() error {
	rowsCh, loaderErrCh := s.streamer.StreamRows()

	entityCh, mapErrCh := s.mapper.MapRowsToCallDetails(rowsCh)

	entities := []hirarki.Hirarki{}

	for {
		select {
		case e, ok := <-entityCh:
			if !ok {
				entityCh = nil
				continue
			}
			entities = append(entities, *e)
		case err, ok := <-loaderErrCh:
			if ok {
				fmt.Println("Loader error:", err)
			}
			loaderErrCh = nil

		case err, ok := <-mapErrCh:
			if ok {
				fmt.Println("Mapping error:", err)
			}
			mapErrCh = nil
		}

		if entityCh == nil && loaderErrCh == nil && mapErrCh == nil {
			break
		}
	}

	return s.saveData(entities)
}

func (s *Service) saveData(data []hirarki.Hirarki) error {
	hasData, _ := s.repo.HasThisMonthData()
	if hasData {
		s.repo.RemoveThisMonthData()
	}
	return s.repo.SaveRange(data)
}
