package masteroutlet

import (
	"fmt"

	"github.com/mariotiara/sfe-data-pipe/internal/application/shared"
	masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"
)

type Service struct {
	repo     masteroutlet.Repository
	streamer shared.StreamDataSource
	mapper   MasterOutletMapper
}

func NewService(repo masteroutlet.Repository, streamer shared.StreamDataSource, mapper MasterOutletMapper) *Service {
	return &Service{repo: repo, streamer: streamer, mapper: mapper}
}

func (s *Service) Run() error {
	rowsCh, loaderErrCh := s.streamer.StreamRows()

	entityCh, mapErrCh := s.mapper.MapRowsToMasterOutlet(rowsCh)

	entities := []masteroutlet.MasterOutlet{}

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

func (s *Service) saveData(data []masteroutlet.MasterOutlet) error {
	hasData, _ := s.repo.HasThisMonthData()
	if hasData {
		s.repo.RemoveThisMonthData()
	}
	return s.repo.SaveRange(data)
}
