package rutracker

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/navigator"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/apex/log"
)

type session struct {
	credentials model.Credentials
	n           navigator.Navigator
	s           provider.CaptchaSolver
	authorized  bool
	l           *log.Entry
}

func newSession(cred model.Credentials, solver provider.CaptchaSolver) (*session, error) {
	n, err := navigator.New("rutracker")
	if err != nil {
		return nil, fmt.Errorf("cannot create session browser: %w", err)
	}
	return &session{
		credentials: cred,
		n:           n,
		s:           solver,
		l:           log.WithField("from", "rutracker").WithField("account", cred.AccountId),
	}, nil
}

func (s *session) authorize(ctx context.Context) error {
	p, err := s.n.NewPage(s.l, ctx)
	if err != nil {
		return fmt.Errorf("create browser page failed: %w", err)
	}
	defer p.Close()

	err = p.Batch("authorization").
		Goto("https://rutracker.org/forum/login.php?redirect=search.php").
		FetchContent().
		Error()

	if err != nil {
		return err
	}

	var captcha struct {
		required bool
		url      string
		code     string
	}

	doc := p.Document()
	captcha.url, captcha.required = doc.Find("#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(3) > td:nth-child(2) > div:nth-child(1) > img").
		First().
		Attr("src")

	if captcha.required {
		s.l.Debug("Captcha is required")

		code, err := s.s.Solve(ctx, provider.Captcha{
			Url:           captcha.url,
			CaseSensitive: false,
			MinLength:     4,
			MaxLength:     6,
		})
		if err != nil {
			return fmt.Errorf("cannot solve captcha: %w", err)
		}

		err = p.Batch("authorization with captcha").
			Goto("https://rutracker.org/forum/login.php?redirect=search.php").
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(1) > td:nth-child(2) > input[type=text]`, s.credentials.Login).
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > input[type=password]`, s.credentials.Password).
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(3) > td:nth-child(2) > div:nth-child(2) > input.reg-input`, code).
			Click(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(4) > td > input`).
			FetchContent().
			Error()

	} else {
		err = p.Batch("authorization").
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(1) > td:nth-child(2) > input[type=text]`, s.credentials.Login).
			Type(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > input[type=password]`, s.credentials.Password).
			Click(`#login-form-full > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(4) > td > input`).
			FetchContent().
			Error()
	}

	if err != nil {
		return err
	}

	doc = p.Document()

	doc.Find("#logged-in-username").Each(func(i int, selection *goquery.Selection) {
		s.authorized = true
	})

	if !s.authorized {
		p.RaiseError(errors.New("cannot parse account logged"))
		return errBadAccount
	}
	return nil
}

func (s *session) search(ctx context.Context, q model.SearchQuery) ([]model.Torrent, error) {
	l := utils.LogFromContext(ctx, "rutracker")
	p, err := s.n.NewPage(l, ctx)
	if err != nil {
		return []model.Torrent{}, fmt.Errorf("cannot create browser page: %w", err)
	}
	defer p.Close()
	torrents := make([]model.Torrent, 0, q.Limit)

	u := "https://rutracker.org/forum/tracker.php?nm=" + url.QueryEscape(q.Query)
	err = p.Batch("searching...").
		Goto(u).
		FetchContent().
		Error()

	if err != nil {
		return nil, err
	}

	doc := p.Document()
	doc.Find(`#tor-tbl > tbody > tr`).Each(func(i int, selection *goquery.Selection) {
		t := parseTorrent(selection)
		if t.IsValid() {
			torrents = append(torrents, t)
		}
	})

	return torrents, nil
}
