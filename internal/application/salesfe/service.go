package salesfe

import (
	"github.com/mariotiara/sfe-data-pipe/internal/domain/salesfe"
)

type Service struct {
	repo salesfe.Repository
}

func NewService(repo salesfe.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImportSales(data []salesfe.SalesFE) error {
	return s.repo.SaveRange(data)
}
