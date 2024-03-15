package torrents

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"time"
)

type downloadLink struct {
	created    time.Time
	downloader model.DownloadFunc
}

const linkExpiredTime = 10 * time.Minute

func (s *Service) processTorrentLink(t *model.Torrent) {
	id, err := s.gen.Generate()
	if err != nil {
		s.log.Errorf("Generate unique id failed: %s", err)
		return
	}
	t.Link = id
	s.links.Store(id, &downloadLink{
		created:    time.Now(),
		downloader: t.Downloader,
	})

	s.log.Debugf("Link generated: %s", id)
	linksRegisteredLinksGauge.Inc()
}

func (s *Service) cleanExpiredLinks() {
	now := time.Now()

	tmp := map[any]struct{}{}
	s.links.Range(func(key, value any) bool {
		dl, ok := value.(*downloadLink)
		if ok && now.Sub(dl.created) >= linkExpiredTime {
			tmp[key] = struct{}{}
			linksRegisteredLinksGauge.Dec()
		}
		return true
	})
	for k, _ := range tmp {
		s.log.Debugf("Link '%s' is expired", k.(string))
		s.links.Delete(k)
	}
}
