package rutracker

import (
	"context"
	"errors"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/apex/log"
	"net/http"
	"strings"
	"sync"
)

var (
	errBadAccount = errors.New("account is unaccessible")
)

type ruTrackerProvider struct {
	log    *log.Entry
	access model.AccessProvider
	s      provider.CaptchaSolver

	mu       sync.RWMutex
	sessions map[string]*session
}

func NewProvider(access model.AccessProvider, solver provider.CaptchaSolver) provider.TorrentsProvider {
	return &ruTrackerProvider{
		log:      log.WithField("from", "rutracker"),
		access:   access,
		sessions: make(map[string]*session),
		s:        solver,
	}
}

func (r *ruTrackerProvider) ID() string {
	return "rutracker"
}

func (r *ruTrackerProvider) SearchTorrents(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	for {
		cred, err := r.access.GetCredentials("rutracker")
		if err != nil {
			return nil, err
		}
		s, err := r.getOrCreateSession(ctx, cred)
		if err != nil {
			if errors.Is(err, errBadAccount) {
				r.access.MarkUnaccesible(cred.AccountId)
				continue
			}
			return nil, err
		}

		result, err := s.search(ctx, q)
		if err != nil {
			return nil, err
		}

		cookies, err := s.n.GetCookies()
		if err != nil {
			return nil, fmt.Errorf("extract cookies failed: %w", err)
		}
		for i := range result {
			t := &result[i]
			t.Downloader = r.newDownloader(t.Link, cookies)
		}

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

func (r *ruTrackerProvider) getOrCreateSession(ctx context.Context, cred model.Credentials) (*session, error) {
	if s, ok := r.getSession(cred.AccountId); ok {
		return s, nil
	}

	s, err := newSession(cred, r.s)
	if err != nil {
		return nil, fmt.Errorf("create new session failed: %w", err)
	}

	if err = s.authorize(ctx); err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	r.mu.Lock()
	r.sessions[cred.AccountId] = s
	r.mu.Unlock()

	return s, nil
}

func (r *ruTrackerProvider) getSession(accountId string) (*session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.sessions[accountId]
	if !ok {
		return nil, ok
	}
	return s, ok
}
