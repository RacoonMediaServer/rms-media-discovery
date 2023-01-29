package torrents

import (
	"context"
	"errors"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	_captcha "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider/2captcha"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider/rutracker"
	"github.com/apex/log"
	"github.com/teris-io/shortid"
	"sync"
	"time"
)

type SearchTypeHint int

const (
	SearchType_Movies SearchTypeHint = iota
	SearchType_Music
	SearchType_Books
	SearchType_Other
)

var ErrExpiredDownloadLink = errors.New("download link expired or not registered")

type Service interface {
	Search(ctx context.Context, query string, hint SearchTypeHint, limit uint) ([]model.Torrent, error)
	Download(ctx context.Context, link string) ([]byte, error)
}

type service struct {
	provider provider.TorrentsProvider
	log      *log.Entry
	gen      *shortid.Shortid
	links    sync.Map
}

func New(access model.AccessProvider) Service {
	return &service{
		provider: newAggregator([]provider.TorrentsProvider{
			rutracker.NewProvider(access, provider.NewCaptchaSolverMonitor(_captcha.NewSolver(access))),
			//rutor.NewProvider(),
		}),
		log: log.WithField("from", "torrents"),
		gen: shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().Nanosecond())),
	}
}
