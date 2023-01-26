package rutracker

import (
	"context"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/media"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/provider"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/scraper"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/gocolly/colly/v2"
	"net/url"
	"strings"
)

type session struct {
	credentials model.Credentials
	c           scraper.Scraper
	s           provider.CaptchaSolver
	authorized  bool
}

func newSession(cred model.Credentials, solver provider.CaptchaSolver) *session {
	return &session{
		credentials: cred,
		c:           scraper.New("rutracker"),
		s:           solver,
	}
}

func (s *session) authorize(ctx context.Context) error {
	s.c.SetContext(ctx)

	var captcha struct {
		required bool
		url      string
		code     string
		sid      string
	}

	sel := s.c.Select("#logged-in-username", func(e *colly.HTMLElement, userData interface{}) {
		s.authorized = true
	})

	sel = sel.SelectResponse(func(response *colly.Response, userData interface{}) {
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

	err := sel.Post("https://rutracker.org/forum/login.php", map[string]string{
		"login_username": s.credentials.Login,
		"login_password": s.credentials.Password,
		"login":          "Вход",
	})

	if err != nil {
		return err
	}

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

		err = sel.Post("https://rutracker.org/forum/login.php", map[string]string{
			"login_username": s.credentials.Login,
			"login_password": s.credentials.Password,
			"login":          "Вход",
			"cap_sid":        captcha.sid,
			captcha.code:     code,
		})
		if err != nil {
			return fmt.Errorf("cannot login with captcha: %w", err)
		}
	}

	if !s.authorized {
		return errBadAccount
	}
	return nil
}

func (s *session) search(ctx context.Context, query string, limit uint) ([]model.Torrent, error) {
	torrents := make([]model.Torrent, 0, limit)

	u := "https://rutracker.org/forum/tracker.php?nm=" + url.QueryEscape(query)
	err := s.c.Select("#tor-tbl > tbody > tr", func(e *colly.HTMLElement, userData interface{}) {
		torrents = append(torrents, parseTorrent(e))
	}).Get(u)

	if err != nil {
		return nil, err
	}

	utils.SortTorrents(torrents)
	torrents = utils.Bound(torrents, limit)

	s.parseDetails(torrents)

	return torrents, nil
}

func (s *session) parseDetails(torrents []model.Torrent) {
	type scrapCtx struct {
		t    *model.Torrent
		done bool
	}

	c := s.c.Clone()
	sel := c.Select(".post_body", func(e *colly.HTMLElement, userData interface{}) {
		ctx := userData.(scrapCtx)
		if !ctx.done {
			ctx.done = true
			_, mediaInfo, ok := strings.Cut(e.Text, "MediaInfo\n")
			if ok {
				ctx.t.Media = media.ParseInfo(mediaInfo)
			}
		}
	})

	for i := range torrents {
		ctx := scrapCtx{t: &torrents[i]}
		sel.GetAsync("https://rutracker.org/forum/"+ctx.t.DetailLink, &ctx)
	}

	c.Wait()
}
