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

// baseURL - адрес сайта (переменная для подмены в тестах)
var baseURL = "https://rutracker.org"

var (
	errBadAccount = errors.New("account is unaccessible")

	// ErrChallenge - сайт отдал проверку Cloudflare вместо страницы. Автоматически ее пройти нельзя:
	// нужно зайти на сайт в обычном браузере и передать сервису свежие cookies (bb_session, cf_clearance) и User-Agent
	// через поля cookies / userAgent аккаунта.
	ErrChallenge = errors.New("blocked by Cloudflare challenge")
)

type ruTrackerProvider struct {
	log    *log.Entry
	access model.AccessProvider
	s      provider.CaptchaSolver

	mu      sync.RWMutex
	cookies map[string][]*http.Cookie
}

// session - состояние авторизации аккаунта, с которым выполняются запросы к сайту
type session struct {
	accountId string
	cookies   []*http.Cookie
	userAgent string

	// provided - cookies заданы в аккаунте вручную (а не получены через логин в браузере)
	provided bool
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

	var lastErr error
	for {
		cred, err := r.access.GetCredentials("rutracker")
		if err != nil {
			if lastErr != nil {
				// чтобы вместо "account not found" было видно, почему закончились аккаунты
				return nil, fmt.Errorf("%w: %w", err, lastErr)
			}
			return nil, err
		}

		sess, err := r.openSession(ctx, cred)
		if err != nil {
			if errors.Is(err, errBadAccount) {
				l.Warnf("Account '%s' is unusable: %s", cred.AccountId, err)
				r.access.MarkUnaccesible(cred.AccountId)
				lastErr = err
				continue
			}
			return nil, err
		}

		result, err := search(ctx, q, sess.cookies, sess.userAgent)
		if err != nil {
			l.Errorf("search failed: %s", err)
			if r.handleSessionError(l, sess, err) {
				lastErr = err
				continue
			}
			return nil, err
		}

		for i := range result {
			t := &result[i]
			t.Downloader = r.newDownloader(t.Link, sess)
		}

		l.Debugf("Got %d results", len(result))
		return result, nil
	}
}

// handleSessionError реагирует на ошибки, говорящие о том, что сессия аккаунта непригодна.
// Возвращает true, если стоит повторить запрос с другим аккаунтом.
func (r *ruTrackerProvider) handleSessionError(l *log.Entry, sess session, err error) (retry bool) {
	switch {
	case errors.Is(err, ErrChallenge):
		if sess.provided {
			// cf_clearance привязан к IP и User-Agent: ниже все, что нужно, чтобы понять, что именно не совпало
			rayID := ""
			var statusErr *requester.StatusError
			if errors.As(err, &statusErr) {
				rayID = statusErr.Header.Get("cf-ray")
			}
			names := make([]string, 0, len(sess.cookies))
			for _, c := range sess.cookies {
				names = append(names, c.Name)
			}
			l.Warnf("Account '%s': cookies rejected by Cloudflare (userAgent=%q, cookies=[%s], cf-ray=%s). "+
				"cf_clearance is bound to the User-Agent and the IP it was issued for: check that userAgent is exactly the browser's one "+
				"and that the service reaches rutracker.org from the same IP as the browser (VPN/proxy/container network), "+
				"otherwise it has expired. Update 'cookies' and 'userAgent' of the account",
				sess.accountId, sess.userAgent, strings.Join(names, ","), rayID)
			r.access.MarkUnaccesible(sess.accountId)
			return true
		}
		// куки, полученные при логине в браузере, устарели
		r.dropCookies(sess.accountId)
		return false

	case errors.Is(err, errBadAccount):
		if sess.provided {
			l.Warnf("Account '%s': session is not authorized, bb_session has expired. Update 'cookies' of the account", sess.accountId)
			r.access.MarkUnaccesible(sess.accountId)
			return true
		}
		// сессия, полученная при логине, протухла - при следующем запросе авторизуемся заново
		r.dropCookies(sess.accountId)
		return false
	}

	return false
}

func (r *ruTrackerProvider) newDownloader(link string, sess session) model.DownloadFunc {
	return func(ctx context.Context) ([]byte, error) {
		req := requester.NewWithUserAgent(r, sess.userAgent)
		req.SetCookies(sess.cookies)
		data, contentType, err := req.Download(ctx, baseURL+"/forum/"+link)
		if err != nil {
			var statusErr *requester.StatusError
			if errors.As(err, &statusErr) && statusErr.IsChallenge() {
				err = fmt.Errorf("%w: %w", ErrChallenge, err)
				r.handleSessionError(utils.LogFromContext(ctx, r.ID()), sess, err)
			}
			return nil, err
		}
		if !strings.HasPrefix(contentType, "application/x-bittorrent") {
			return nil, errors.New("unexpected Content-Type")
		}
		return data, err
	}
}

// openSession возвращает сессию аккаунта: заданные в аккаунте cookies либо результат логина в браузере
func (r *ruTrackerProvider) openSession(ctx context.Context, cred model.Credentials) (session, error) {
	if cred.Cookies != "" {
		cookies := parseCookies(cred.Cookies)
		if len(cookies) == 0 {
			return session{}, fmt.Errorf("%w: cookies are empty or malformed", errBadAccount)
		}
		return session{accountId: cred.AccountId, cookies: cookies, userAgent: cred.UserAgent, provided: true}, nil
	}

	cookies, err := r.login(ctx, cred)
	if err != nil {
		return session{}, err
	}
	return session{accountId: cred.AccountId, cookies: cookies}, nil
}

// parseCookies разбирает строку в формате заголовка Cookie ("a=b; c=d")
func parseCookies(s string) []*http.Cookie {
	req := http.Request{Header: http.Header{"Cookie": []string{s}}}
	return req.Cookies()
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

func (r *ruTrackerProvider) dropCookies(accountId string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cookies, accountId)
}
