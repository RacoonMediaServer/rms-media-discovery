package torrents

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"time"
)

type downloadLink struct {
	created    time.Time
	downloader model.DownloadFunc
}

const linkExpiredTime = 10 * time.Minute

func (s *service) processTorrentLink(t *model.Torrent) {
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
}

func (s *service) cleanExpiredLinks() {
	now := time.Now()

	tmp := map[any]struct{}{}
	s.links.Range(func(key, value any) bool {
		dl, ok := value.(*downloadLink)
		if ok && now.Sub(dl.created) >= linkExpiredTime {
			tmp[key] = struct{}{}
		}
		return true
	})
	for k, _ := range tmp {
		s.links.Delete(k)
	}
}
