package torrents

import (
	"context"
	"errors"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/aggregator"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/thepiratebay"
	"github.com/apex/log"
	"github.com/teris-io/shortid"
	"sync"
	"time"
)

var ErrExpiredDownloadLink = errors.New("download link expired or not registered")

type Service interface {
	Search(ctx context.Context, query model.SearchQuery) ([]model.Torrent, error)
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
		provider: aggregator.NewTorrentProvider(aggregator.FastPolicy, []provider.TorrentsProvider{
			//rutracker.NewProvider(access, provider.NewCaptchaSolverMonitor(_captcha.NewSolver(access))),
			//rutor.NewProvider(),
			thepiratebay.New(),
		}),
		log: log.WithField("from", "torrents"),
		gen: shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().Nanosecond())),
	}
}
