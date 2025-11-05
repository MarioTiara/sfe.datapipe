package ezengagecalldetail

import "github.com/mariotiara/sfe-data-pipe/internal/domain/ezengagecalldetail"

type Service struct {
	repo ezengagecalldetail.Repository
}

func NewService(repo ezengagecalldetail.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImportCallDetails(data []ezengagecalldetail.EZEngageCallDetail) error {
	return s.repo.SaveRange(data)
}
