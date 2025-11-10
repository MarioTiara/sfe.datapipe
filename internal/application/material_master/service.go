package materialmaster

import (
	"fmt"

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

	entityCh, mapErrCh := s.mapper.MapRowsToMaterialMaster(rowsCh)

	entities := []materialmaster.MaterialMaster{}

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

func (s *Service) saveData(data []materialmaster.MaterialMaster) error {
	hasData, _ := s.repo.HasThisMonthData()
	if hasData {
		s.repo.RemoveThisMonthData()
	}
	return s.repo.SaveRange(data)
}
