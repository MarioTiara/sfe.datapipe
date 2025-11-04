package materialmaster

import materialmaster "github.com/mariotiara/sfe-data-pipe/internal/domain/material_master"

type Service struct {
	repo materialmaster.Repository
}

func NewService(repo materialmaster.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImportMaterials(data []materialmaster.MaterialMaster) error {
	return s.repo.SaveRange(data)
}
