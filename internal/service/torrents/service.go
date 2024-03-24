package torrents

import (
	"context"
	"errors"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/anidub"
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

type Service struct {
	provider provider.TorrentsProvider
	log      *log.Entry
	gen      *shortid.Shortid
	links    sync.Map

	tasks  sync.Map
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(access model.AccessProvider) *Service {
	s := Service{
		//provider: aggregator.NewTorrentProvider(aggregator.PriorityPolicy, []provider.TorrentsProvider{
		//	rutracker.NewProvider(access, provider.NewCaptchaSolverMonitor(_captcha.NewSolver(access))),
		//	rutor.NewProvider(),
		//	thepiratebay.New(),
		//}),
		provider: anidub.New(),
		log:      log.WithField("from", "torrents"),
		gen:      shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().Nanosecond())),
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.gcResourcesProcess()
	}()

	return &s
}

func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
}
