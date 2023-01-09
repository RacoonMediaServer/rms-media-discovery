package movies

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

type Service interface {
	Search(query string, limit uint) ([]model.Movie, error)
}

type service struct {
}

func New() Service {
	return &service{}
}

func (s *service) Search(query string, limit uint) ([]model.Movie, error) {
	return nil, nil
}
