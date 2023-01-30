package torrents

import (
	"context"
	"errors"
	model2 "git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	provider2 "git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/2captcha"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/aggregator"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/rutor"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider/rutracker"
	"github.com/apex/log"
	"github.com/teris-io/shortid"
	"sync"
	"time"
)

var ErrExpiredDownloadLink = errors.New("download link expired or not registered")

type Service interface {
	Search(ctx context.Context, query model2.SearchQuery) ([]model2.Torrent, error)
	Download(ctx context.Context, link string) ([]byte, error)
}

type service struct {
	provider provider2.TorrentsProvider
	log      *log.Entry
	gen      *shortid.Shortid
	links    sync.Map
}

func New(access model2.AccessProvider) Service {
	return &service{
		provider: aggregator.NewTorrentProvider(aggregator.FastPolicy, []provider2.TorrentsProvider{
			rutracker.NewProvider(access, provider2.NewCaptchaSolverMonitor(_captcha.NewSolver(access))),
			rutor.NewProvider(),
		}),
		log: log.WithField("from", "torrents"),
		gen: shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().Nanosecond())),
	}
}
