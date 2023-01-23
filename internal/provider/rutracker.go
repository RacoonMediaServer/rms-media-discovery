package provider

import (
	"context"
	"errors"
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"net/url"
	"regexp"
	"strconv"
	"sync"
)

type ruTrackerProvider struct {
	log    *log.Entry
	access model.AccessProvider
	s      CaptchaSolver

	mu       sync.RWMutex
	sessions map[string]*ruTrackerSession
}

type ruTrackerSession struct {
	credentials model.Credentials
	c           *colly.Collector
	s           CaptchaSolver
	authorized  bool
}

var (
	captchaSidExpr  = regexp.MustCompile(`<input[^>]*name="cap_sid"[^>]*value="([^"]+)"[^>]*>`)
	captchaCodeExpr = regexp.MustCompile(`<input[^>]*name="(cap_code_[^"]+)"[^>]*value="[^"]*"[^>]*>`)
	captchaUrlExpr  = regexp.MustCompile(`<img[^>]*src="([^"]+\/captcha\/[^"]+)"[^>]*>`)
	extractSizeExpr = regexp.MustCompile(`^(\d+(.\d+)?) (MB|GB)`)
)

func NewRuTrackerProvider(access model.AccessProvider, solver CaptchaSolver) TorrentsProvider {
	return &ruTrackerProvider{
		log:      log.WithField("from", "rutracker"),
		access:   access,
		sessions: make(map[string]*ruTrackerSession),
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

func (r *ruTrackerProvider) getOrCreateSession(ctx context.Context, cred model.Credentials) (*ruTrackerSession, error) {
	if session, ok := r.getSession(cred.AccountId); ok {
		return session, nil
	}

	session := newRuTrackerSession(cred, r.s)

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

func (r *ruTrackerProvider) getSession(accountId string) (*ruTrackerSession, bool) {
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

func newRuTrackerSession(cred model.Credentials, solver CaptchaSolver) *ruTrackerSession {
	return &ruTrackerSession{
		credentials: cred,
		c: colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
			colly.AllowURLRevisit(),
		),
		s: solver,
	}
}

func (s *ruTrackerSession) authorize(ctx context.Context) error {
	s.c.SetDebugger(&debug.LogDebugger{})
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
		code, err := s.s.Solve(ctx, Captcha{
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

func (s *ruTrackerSession) search(ctx context.Context, query string) ([]model.Torrent, error) {
	torrents := make([]model.Torrent, 0)

	wg := sync.WaitGroup{}
	s.c.OnHTML("#tor-tbl > tbody > tr", func(e *colly.HTMLElement) {
		torrents = append(torrents, extractTorrent(e))
		href, ok := e.DOM.Find(`a.tLink`).Attr("href")
		if ok {
			//t := &torrents[len(torrents)-1]
			c := s.c.Clone()
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = c.Visit("https://rutracker.org/forum/" + href)
				c.Wait()
			}()
		}
	})
	wg.Wait()

	if err := s.c.Visit("https://rutracker.org/forum/tracker.php?nm=" + url.QueryEscape(query)); err != nil {
		return nil, err
	}
	s.c.Wait()

	return torrents, nil
}

func parseTorrentSize(text string) float32 {
	matches := extractSizeExpr.FindStringSubmatch(text)
	if matches != nil {
		result, err := strconv.ParseFloat(matches[1], 32)
		if err != nil {
			return 0
		}
		if matches[3] == "GB" {
			result *= 1024.
		}
		return float32(result)
	}

	return 0
}

func extractTorrent(e *colly.HTMLElement) model.Torrent {
	torrent := model.Torrent{}
	torrent.Title = e.DOM.Find(`a.tLink`).Text()

	dl := e.DOM.Find(`a.tr-dl`)
	link, _ := dl.Attr("href")
	torrent.Link = link
	torrent.SizeMB = parseTorrentSize(dl.Text())

	seeds := e.DOM.Find(`b.seedmed`).Text()
	seedersCount, _ := strconv.ParseUint(seeds, 10, 32)
	torrent.Seeders = uint(seedersCount)

	leechs := e.DOM.Find(`td.leechmed`).Text()
	peers, _ := strconv.Atoi(leechs)
	torrent.Seeders += uint(peers)

	return torrent
}
