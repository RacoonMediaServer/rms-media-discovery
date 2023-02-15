package anidub

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/scraper"
	"github.com/apex/log"
	"net/http"
	"sync"
)

type anidubProvider struct {
	access  model.AccessProvider
	l       *log.Entry
	cookies sync.Map
}

func (a *anidubProvider) ID() string {
	return "anidub"
}

func (a *anidubProvider) SearchTorrents(ctx context.Context, query model.SearchQuery) ([]model.Torrent, error) {
	if query.Type == model.Music || query.Type == model.Books {
		return []model.Torrent{}, nil
	}

	cred, err := a.access.GetCredentials(a.ID())
	if err != nil {
		return nil, err
	}

	cookies, err := a.getCookies(ctx, cred)
	if err != nil {
		return nil, err
	}

	s := scraper.New(a.ID())
	s.SetContext(ctx)
	if err = s.SetCookies("https://tr.anidub.com/", cookies); err != nil {
		return nil, fmt.Errorf("set cookies failed: %w", err)
	}

	// TODO:
	var results []model.Torrent
	return results, nil
}

func (a *anidubProvider) getCookies(ctx context.Context, cred model.Credentials) (cookies []*http.Cookie, err error) {
	cookieAny, ok := a.cookies.Load(cred.AccountId)
	if !ok {
		cookies, err = a.authorize(ctx, cred)
		if err != nil {
			return nil, fmt.Errorf("auth failed: %w", err)
		}
		a.cookies.Store(cred.AccountId, cookies)
	} else {
		cookies = cookieAny.([]*http.Cookie)
	}

	return
}

func (a *anidubProvider) authorize(ctx context.Context, cred model.Credentials) ([]*http.Cookie, error) {
	s := scraper.New(a.ID())
	s.SetContext(ctx)

	isLogged := false
	err := s.Select("", loginChecker(&isLogged)).Post("https://tr.anidub.com/index.php?do=register", map[string]string{
		"login_name":     cred.Login,
		"login_password": cred.Password,
		"login":          "submit",
	})

	if err != nil {
		return nil, err
	}

	if !isLogged {
		return nil, fmt.Errorf("probably account is expired")
	}

	return s.Cookies("https://tr.anidub.com/"), nil
}

func New(access model.AccessProvider) provider.TorrentsProvider {
	return &anidubProvider{
		access: access,
		l:      log.WithField("from", "anidub"),
	}
}
