package rutracker

import (
	"context"
	"errors"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/media"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/navigator"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
	"net/url"
	"strings"
	"sync"
	"time"
)

type session struct {
	credentials model.Credentials
	n           navigator.Navigator
	s           provider.CaptchaSolver
	authorized  bool
	l           *log.Entry
}

const emulateDelay = 300 * time.Millisecond
const parseDetailsTimeout = 10 * time.Second

func newSession(cred model.Credentials, solver provider.CaptchaSolver) (*session, error) {
	n, err := navigator.New()
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

func (s *session) search(ctx context.Context, query string, limit uint) ([]model.Torrent, error) {
	l := utils.LogFromContext(ctx, "rutracker", s.l)
	p, err := s.n.NewPage(l, ctx)
	if err != nil {
		return []model.Torrent{}, fmt.Errorf("cannot create browser page: %w", err)
	}
	torrents := make([]model.Torrent, 0, limit)

	u := "https://rutracker.org/forum/tracker.php?nm=" + url.QueryEscape(query)
	err = p.Batch("searching...").
		Goto(u).
		FetchContent().
		Error()

	if err != nil {
		return nil, err
	}

	doc := p.Document()
	doc.Find(`#tor-tbl > tbody > tr`).Each(func(i int, selection *goquery.Selection) {
		torrents = append(torrents, parseTorrent(selection))
	})

	utils.SortTorrents(torrents)
	torrents = utils.Bound(torrents, limit)

	s.parseDetails(ctx, torrents)

	return torrents, nil
}

func (s *session) parseDetails(ctx context.Context, torrents []model.Torrent) {
	wg := sync.WaitGroup{}
	childCtx, cancel := context.WithTimeout(ctx, parseDetailsTimeout)
	defer cancel()

	for i := range torrents {
		t := &torrents[i]
		p, err := s.n.NewPage(s.l, childCtx)
		if err != nil {
			s.l.Warnf("Cannot create new page: %s", err)
			continue
		}

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err = p.Batch(fmt.Sprintf("parse details #%d", i)).
				Goto("https://rutracker.org/forum/" + t.DetailLink).
				FetchContent().
				Error()
			if err != nil {
				s.l.Warnf("Cannot parse details: %s", err)
				return
			}

			post := p.Document().Find(`.post_body`).First()
			_, mediaInfo, ok := strings.Cut(post.Text(), "MediaInfo\n")
			if ok {
				t.Media = media.ParseInfo(mediaInfo)
			}
		}(i)
	}
	wg.Wait()
}
