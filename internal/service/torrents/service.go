package torrents

import (
	"context"
	"errors"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	_captcha "github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/2captcha"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/aggregator"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/rutor"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/rutracker"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/thepiratebay"
	"github.com/apex/log"
	"github.com/teris-io/shortid"
	"sync"
	"time"
)

var ErrExpiredDownloadLink = errors.New("download link expired or not registered")

type TaskStatus struct {
	Status      model.TaskStatus
	Results     []model.Torrent
	ContentType media.ContentType
	Err         error
}

type Service interface {
	Search(ctx context.Context, query model.SearchQuery) ([]model.Torrent, error)
	SearchAsync(query model.SearchQuery) (taskID string, err error)
	Status(taskID string) (TaskStatus, error)
	Cancel(taskID string) error
	Download(ctx context.Context, link string) ([]byte, error)
}

type service struct {
	provider provider.TorrentsProvider
	log      *log.Entry
	gen      *shortid.Shortid
	links    sync.Map

	tasks  sync.Map
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(access model.AccessProvider) Service {
	s := service{
		provider: aggregator.NewTorrentProvider(aggregator.PriorityPolicy, []provider.TorrentsProvider{
			rutracker.NewProvider(access, provider.NewCaptchaSolverMonitor(_captcha.NewSolver(access))),
			rutor.NewProvider(),
			thepiratebay.New(),
		}),
		log: log.WithField("from", "torrents"),
		gen: shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().Nanosecond())),
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	return &s
}
