package torrents

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider/2captcha"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider/rutracker"
	"github.com/apex/log"
)

type SearchTypeHint int

const (
	SearchType_Movies SearchTypeHint = iota
	SearchType_Music
	SearchType_Books
	SearchType_Other
)

type Service interface {
	Search(ctx context.Context, query string, hint SearchTypeHint, limit uint) ([]model.Torrent, error)
	Download(ctx context.Context, link string) ([]byte, error)
}

type service struct {
	provider provider.TorrentsProvider
	log      *log.Entry
}

func New(access model.AccessProvider) Service {
	return &service{
		provider: newAggregator([]provider.TorrentsProvider{
			rutracker.NewProvider(access, _captcha.NewSolver(access)),
		}),
		log: log.WithField("from", "torrents"),
	}
}
