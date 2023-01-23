package rutracker

import (
	"context"
	"errors"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"github.com/apex/log"
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

func (r *ruTrackerProvider) SearchTorrents(ctx context.Context, query string) ([]model.Torrent, error) {
	for {
		cred, err := r.access.GetCredentials("rutracker")
		if err != nil {
			return nil, err
		}
		session, err := r.getOrCreateSession(ctx, cred)
		if err != nil {
			if errors.Is(err, errBadAccount) {
				r.access.MarkUnaccesible(cred.AccountId)
				continue
			}
			return nil, err
		}

		return session.search(ctx, query)
	}
}

func (r *ruTrackerProvider) Download(ctx context.Context, link string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ruTrackerProvider) getOrCreateSession(ctx context.Context, cred model.Credentials) (*session, error) {
	if session, ok := r.getSession(cred.AccountId); ok {
		return session, nil
	}

	session := newSession(cred, r.s)

	if err := session.authorize(ctx); err != nil {
		return nil, fmt.Errorf("auth failed: %w", err)
	}

	r.mu.Lock()
	r.sessions[cred.AccountId] = session
	r.mu.Unlock()

	newSession := *session
	newSession.c = session.c.Clone()
	return &newSession, nil
}

func (r *ruTrackerProvider) getSession(accountId string) (*session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, ok := r.sessions[accountId]
	if !ok {
		return nil, ok
	}
	newSession := *session
	newSession.c = session.c.Clone()
	return &newSession, ok
}
