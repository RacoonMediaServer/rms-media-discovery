package rutracker

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/scraper"
	"github.com/gocolly/colly/v2"
)

func search(ctx context.Context, q model.SearchQuery, cookies []*http.Cookie, userAgent string) ([]model.Torrent, error) {
	c := scraper.NewWithOptions("rutracker", scraper.Options{UserAgent: userAgent, FailOnRequestError: true})
	c.SetContext(ctx)

	if err := c.SetCookies(baseURL+"/", cookies); err != nil {
		return []model.Torrent{}, fmt.Errorf("set login cookies failed: %w", err)
	}

	torrents := make([]model.Torrent, 0, q.Limit)
	loggedIn := false

	u := baseURL + "/forum/tracker.php?nm=" + url.QueryEscape(q.Query)
	err := c.Select(`#logged-in-username`, func(e *colly.HTMLElement, userData interface{}) { loggedIn = true }).
		Select(`#tor-tbl > tbody > tr`, getTorrentParser(&torrents)).
		Get(u)
	if err != nil {
		var statusErr *requester.StatusError
		if errors.As(err, &statusErr) && statusErr.IsChallenge() {
			return []model.Torrent{}, fmt.Errorf("%w: %w", ErrChallenge, err)
		}
		return []model.Torrent{}, err
	}

	// при невалидной сессии сайт редиректит на страницу логина, и без этой проверки мы бы молча вернули пустой результат
	if !loggedIn {
		return []model.Torrent{}, fmt.Errorf("%w: session is not authorized", errBadAccount)
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
