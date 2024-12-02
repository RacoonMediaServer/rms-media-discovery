package rutracker

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/scraper"
	"github.com/gocolly/colly/v2"
)

func search(ctx context.Context, q model.SearchQuery, cookies []*http.Cookie) ([]model.Torrent, error) {
	c := scraper.New("rutracker")
	c.SetContext(ctx)

	if err := c.SetCookies("https://rutracker.org/", cookies); err != nil {
		return []model.Torrent{}, fmt.Errorf("set login cookies failed: %w", err)
	}

	torrents := make([]model.Torrent, 0, q.Limit)

	u := "https://rutracker.org/forum/tracker.php?nm=" + url.QueryEscape(q.Query)
	err := c.Select(`#tor-tbl > tbody > tr`, getTorrentParser(&torrents)).Get(u)
	if err != nil {
		return []model.Torrent{}, err
	}
	return torrents, nil
}

func getTorrentParser(result *[]model.Torrent) scraper.HTMLCallback {
	return func(e *colly.HTMLElement, userData interface{}) {
		t := parseTorrent(e)
		if t.IsValid() {
			*result = append(*result, t)
		}
	}
}
