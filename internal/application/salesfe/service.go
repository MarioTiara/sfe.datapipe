package salesfe

import (
	"fmt"

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

	entityCh, mapErrCh := s.mapper.MapRowsToSaleFE(rowsCh)

	entities := []salesfe.SalesFE{}

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

func (s *Service) saveData(data []salesfe.SalesFE) error {
	hasData, _ := s.repo.HasThisMonthData()
	if hasData {
		s.repo.RemoveThisMonthData()
	}
	return s.repo.SaveRange(data)
}
