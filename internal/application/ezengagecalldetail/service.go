package ezengagecalldetail

import (
	"fmt"

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
	rowsCh, loaderErrCh := s.streamer.StreamRows()

	entityCh, mapErrCh := s.mapper.MapRowsToCallDetails(rowsCh)

	entities := []ezengagecalldetail.EZEngageCallDetail{}

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

	return s.repo.SaveRange(entities)
}

func (s *Service) saveData(data []ezengagecalldetail.EZEngageCallDetail) error {
	hasData, _ := s.repo.HasThisMonthData()
	if hasData {
		row, _ := s.repo.RemoveThisMonthData()
		ms := fmt.Sprintf("%d is removed", row)
		fmt.Println(ms)
	}

	return s.repo.SaveRange(data)
}
