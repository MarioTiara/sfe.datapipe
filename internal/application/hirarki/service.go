package hirarki

import "github.com/mariotiara/sfe-data-pipe/internal/domain/hirarki"

type Service struct {
	repo hirarki.Repository
}

func NewService(repo hirarki.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImportHirarkiList(data []hirarki.Hirarki) error {
	return s.repo.SaveRange(data)
}
