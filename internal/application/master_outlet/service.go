package masteroutlet

import masteroutlet "github.com/mariotiara/sfe-data-pipe/internal/domain/master_outlet"

type Service struct {
	repo masteroutlet.Repository
}

func NewService(repo masteroutlet.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImporOutlets(data []masteroutlet.MasterOutlet) error {
	return s.repo.SaveRange(data)
}
