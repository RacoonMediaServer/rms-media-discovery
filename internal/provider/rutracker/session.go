package rutracker

import (
	"context"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"net/http"
	"net/url"
)

type session struct {
	credentials model.Credentials
	c           *colly.Collector
	s           provider.CaptchaSolver
	authorized  bool
}

func newSession(cred model.Credentials, solver provider.CaptchaSolver) *session {
	return &session{
		credentials: cred,
		c: colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
			colly.AllowURLRevisit(),
		),
		s: solver,
	}
}

func (s *session) authorize(ctx context.Context) error {
	s.c.SetDebugger(&debug.LogDebugger{})
	collyWithContext(s.c, ctx)

	var captcha struct {
		required bool
		url      string
		code     string
		sid      string
	}

	s.c.OnHTML("#logged-in-username", func(e *colly.HTMLElement) {
		s.authorized = true
	})

	s.c.OnResponse(func(response *colly.Response) {
		content := string(response.Body)

		matches := captchaUrlExpr.FindStringSubmatch(content)
		if len(matches) < 2 {
			return
		}

		captcha.required = true
		captcha.url = matches[1]

		matches = captchaCodeExpr.FindStringSubmatch(content)
		if len(matches) >= 2 {
			captcha.code = matches[1]
		}

		matches = captchaSidExpr.FindStringSubmatch(content)
		if len(matches) >= 2 {
			captcha.sid = matches[1]
		}
	})

	err := s.c.Post("https://rutracker.org/forum/login.php", map[string]string{
		"login_username": s.credentials.Login,
		"login_password": s.credentials.Password,
		"login":          "Вход",
	})

	if err != nil {
		return err
	}
	s.c.Wait()

	if captcha.required {
		code, err := s.s.Solve(ctx, provider.Captcha{
			Url:           captcha.url,
			CaseSensitive: false,
			MinLength:     4,
			MaxLength:     6,
		})
		if err != nil {
			return fmt.Errorf("cannot solve captcha: %w", err)
		}

		err = s.c.Post("https://rutracker.org/forum/login.php", map[string]string{
			"login_username": s.credentials.Login,
			"login_password": s.credentials.Password,
			"login":          "Вход",
			"cap_sid":        captcha.sid,
			captcha.code:     code,
		})
		if err != nil {
			return fmt.Errorf("cannot login with captcha: %w", err)
		}
		s.c.Wait()
	}

	if !s.authorized {
		return errBadAccount
	}
	return nil
}

func (s *session) search(ctx context.Context, query string) ([]model.Torrent, error) {
	torrents := make([]model.Torrent, 0)
	collyWithContext(s.c, ctx)

	scrapper := newScrapper(ctx)
	s.c.OnHTML("#tor-tbl > tbody > tr", func(e *colly.HTMLElement) {
		torrents = append(torrents, parseTorrent(e))
		href, ok := e.DOM.Find(`a.tLink`).Attr("href")
		if ok {
			scrapper.scrapAsync(s.c, &torrents[len(torrents)-1], "https://rutracker.org/forum/"+href)
		}
	})

	if err := s.c.Visit("https://rutracker.org/forum/tracker.php?nm=" + url.QueryEscape(query)); err != nil {
		return nil, err
	}
	s.c.Wait()

	scrapper.wait()

	return torrents, nil
}

type contextTransport struct {
	ctx       context.Context
	transport *http.Transport
}

func (t *contextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.WithContext(t.ctx)
	return t.transport.RoundTrip(req)
}

func collyWithContext(c *colly.Collector, ctx context.Context) {
	c.OnRequest(func(req *colly.Request) {
		select {
		case <-ctx.Done():
			req.Abort()
		default:
		}
	})

	trans := &contextTransport{
		ctx:       ctx,
		transport: &http.Transport{},
	}
	c.WithTransport(trans)
}
