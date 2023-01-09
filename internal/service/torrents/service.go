package torrents

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

type SearchTypeHint int

const (
	SearchType_Movies SearchTypeHint = iota
	SearchType_Music
	SearchType_Books
	SearchType_Other
)

type Service interface {
	Search(query string, hint SearchTypeHint, limit uint) ([]model.Torrent, error)
	Download(link string) ([]byte, error)
}

type service struct {
}

func New() Service {
	return &service{}
}

func (s *service) Search(query string, hint SearchTypeHint, limit uint) ([]model.Torrent, error) {
	return nil, nil
}

func (s *service) Download(link string) ([]byte, error) {
	return nil, nil
}
