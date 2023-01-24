package rutracker

import (
	"context"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/media"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"github.com/gocolly/colly/v2"
	"strings"
	"sync"
)

type scrapper struct {
	ctx context.Context
	wg  sync.WaitGroup
}

func newScrapper(ctx context.Context) *scrapper {
	return &scrapper{
		ctx: ctx,
	}
}

func (s *scrapper) scrapAsync(c *colly.Collector, torrent *model.Torrent, link string) {
	c = c.Clone()
	collyWithContext(c, s.ctx)

	t := *torrent

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		result, err := s.scrap(c, link, t)
		if err != nil {
			return
		}
		*torrent = result
	}()
}

func (s *scrapper) scrap(c *colly.Collector, link string, t model.Torrent) (model.Torrent, error) {
	gotFirst := false
	c.OnHTML(".post_body", func(e *colly.HTMLElement) {
		if !gotFirst {
			gotFirst = true
			_, mediaInfo, ok := strings.Cut(e.Text, "MediaInfo\n")
			if ok {
				t.Media = media.ParseInfo(mediaInfo)
			}
		}
	})
	err := c.Visit(link)
	if err != nil {
		return model.Torrent{}, err
	}
	c.Wait()

	return t, nil
}

func (s *scrapper) wait() {
	s.wg.Wait()
}
