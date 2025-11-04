package customerfe

import "github.com/mariotiara/sfe-data-pipe/internal/domain/customerfe"

type Service struct {
	repo customerfe.Repository
}

func NewService(repo customerfe.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImportCustomers(data []customerfe.CustomerFE) error {
	return s.repo.SaveRange(data)
}
