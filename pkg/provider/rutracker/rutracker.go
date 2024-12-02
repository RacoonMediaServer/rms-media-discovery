package rutracker

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/apex/log"
)

var (
	errBadAccount = errors.New("account is unaccessible")
)

type ruTrackerProvider struct {
	log    *log.Entry
	access model.AccessProvider
	s      provider.CaptchaSolver

	mu      sync.RWMutex
	cookies map[string][]*http.Cookie
}

func NewProvider(access model.AccessProvider, solver provider.CaptchaSolver) provider.TorrentsProvider {
	return &ruTrackerProvider{
		log:     log.WithField("from", "rutracker"),
		access:  access,
		cookies: make(map[string][]*http.Cookie),
		s:       solver,
	}
}

func (r *ruTrackerProvider) ID() string {
	return "rutracker"
}

func applySearchHints(q *model.SearchQuery) {
	// применяем дополнительные параметры поиска так, как это лучше всего будет работать на конкретном трекере
	if q.Type == media.Movies {
		if q.Year != nil {
			q.Query += fmt.Sprintf(" %d", *q.Year)
		}
		if q.Season != nil {
			q.Query += fmt.Sprintf(" сезон %d", *q.Season)
		}
	}
	if q.Type == media.Music && q.Discography {
		q.Query += " дискография"
	}
}

func (r *ruTrackerProvider) SearchTorrents(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	l := utils.LogFromContext(ctx, r.ID())
	applySearchHints(&q)

	for {
		cred, err := r.access.GetCredentials("rutracker")
		if err != nil {
			return nil, err
		}

		cookies, err := r.login(ctx, cred)
		if err != nil {
			if errors.Is(err, errBadAccount) {
				r.access.MarkUnaccesible(cred.AccountId)
				continue
			}
			return nil, err
		}

		result, err := search(ctx, q, cookies)
		if err != nil {
			l.Errorf("search failed: %s", err)
			return nil, err
		}

		for i := range result {
			t := &result[i]
			t.Downloader = r.newDownloader(t.Link, cookies)
		}

		l.Debugf("Got %d results", len(result))
		return result, err
	}
}

func (r *ruTrackerProvider) newDownloader(link string, cookies []*http.Cookie) model.DownloadFunc {
	return func(ctx context.Context) ([]byte, error) {
		r := requester.New(r)
		r.SetCookies(cookies)
		data, contentType, err := r.Download(ctx, "https://rutracker.org/forum/"+link)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(contentType, "application/x-bittorrent") {
			return nil, errors.New("unexpected Content-Type")
		}
		return data, err
	}
}

func (r *ruTrackerProvider) login(ctx context.Context, cred model.Credentials) ([]*http.Cookie, error) {
	if s, ok := r.getCookies(cred.AccountId); ok {
		return s, nil
	}

	cookies, err := authorize(ctx, cred, r.s)
	if err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	r.mu.Lock()
	r.cookies[cred.AccountId] = cookies
	r.mu.Unlock()

	return cookies, nil
}

func (r *ruTrackerProvider) getCookies(accountId string) ([]*http.Cookie, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cookies, ok := r.cookies[accountId]
	if !ok {
		return nil, ok
	}
	return cookies, ok
}
