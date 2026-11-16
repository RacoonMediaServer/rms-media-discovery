package rutracker

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
)

const (
	goodUA      = "UA-good"
	goodCookies = "bb_session=good; cf_clearance=cf"

	resultsPage = `<html><body><span id="logged-in-username">user</span>
<table id="tor-tbl"><tbody>
<tr><td><a class="tLink" href="viewtopic.php?t=1">Матрица</a></td>
<td><a class="tr-dl" href="dl.php?t=1">11 MB ↓</a></td>
<td><b class="seedmed">7</b></td><td class="leechmed">2</td></tr>
</tbody></table></body></html>`

	// так выглядит ответ сайта, если сессия не авторизована: редирект на страницу логина
	loginPage = `<html><body><form id="login-form-full"></form></body></html>`
)

// fakeSite эмулирует rutracker за Cloudflare: cf_clearance привязан к User-Agent
func fakeSite(t *testing.T) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cf, err := r.Cookie("cf_clearance")
		if err != nil || cf.Value != "cf" || r.UserAgent() != goodUA {
			w.Header().Set("cf-mitigated", "challenge")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		session, err := r.Cookie("bb_session")
		authorized := err == nil && session.Value == "good"

		switch r.URL.Path {
		case "/forum/tracker.php":
			if !authorized {
				_, _ = w.Write([]byte(loginPage))
				return
			}
			_, _ = w.Write([]byte(resultsPage))
		case "/forum/dl.php":
			w.Header().Set("Content-Type", "application/x-bittorrent; charset=Windows-1251")
			_, _ = w.Write([]byte("d8:announce"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	prev := baseURL
	baseURL = srv.URL
	t.Cleanup(func() {
		baseURL = prev
		srv.Close()
	})
	return srv
}

// fakeAccess выдает аккаунты по порядку, пропуская помеченные недоступными (как настоящий сервис аккаунтов)
type fakeAccess struct {
	accounts []model.Credentials
	marked   []string
}

func (a *fakeAccess) GetCredentials(string) (model.Credentials, error) {
	for _, acc := range a.accounts {
		if !a.isMarked(acc.AccountId) {
			return acc, nil
		}
	}
	return model.Credentials{}, errors.New("account not found")
}

func (a *fakeAccess) GetApiKey(string) (model.ApiKey, error) { return model.ApiKey{}, nil }

func (a *fakeAccess) MarkUnaccesible(id string) { a.marked = append(a.marked, id) }

func (a *fakeAccess) isMarked(id string) bool {
	for _, m := range a.marked {
		if m == id {
			return true
		}
	}
	return false
}

func TestSearchTorrents_ProvidedCookies(t *testing.T) {
	fakeSite(t)
	access := &fakeAccess{accounts: []model.Credentials{{AccountId: "rutracker.1", Cookies: goodCookies, UserAgent: goodUA}}}
	p := NewProvider(access, nil)

	result, err := p.SearchTorrents(context.Background(), model.SearchQuery{Query: "Матрица"})
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Матрица", result[0].Title)
	assert.Equal(t, uint64(11), result[0].SizeMB)
	assert.Empty(t, access.marked)

	// скачивание использует ту же сессию и UA
	data, err := result[0].Downloader(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "d8:announce", string(data))
}

func TestSearchTorrents_ChallengeSkipsToNextAccount(t *testing.T) {
	fakeSite(t)
	access := &fakeAccess{accounts: []model.Credentials{
		// cf_clearance получен на другом User-Agent - Cloudflare его отклонит
		{AccountId: "rutracker.stale", Cookies: goodCookies, UserAgent: "another-UA"},
		{AccountId: "rutracker.ok", Cookies: goodCookies, UserAgent: goodUA},
	}}
	p := NewProvider(access, nil)

	result, err := p.SearchTorrents(context.Background(), model.SearchQuery{Query: "Матрица"})
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, []string{"rutracker.stale"}, access.marked)
}

func TestSearchTorrents_ChallengeOnlyAccount(t *testing.T) {
	fakeSite(t)
	access := &fakeAccess{accounts: []model.Credentials{{AccountId: "rutracker.1", Cookies: "bb_session=good", UserAgent: goodUA}}}
	p := NewProvider(access, nil)

	_, err := p.SearchTorrents(context.Background(), model.SearchQuery{Query: "Матрица"})
	// причина должна быть видна в ошибке, а не только "account not found"
	assert.ErrorIs(t, err, ErrChallenge)
	assert.Equal(t, []string{"rutracker.1"}, access.marked)
}

func TestSearchTorrents_ExpiredSession(t *testing.T) {
	fakeSite(t)
	// cf_clearance в порядке, но bb_session протух - сайт показывает страницу логина вместо результатов
	access := &fakeAccess{accounts: []model.Credentials{{AccountId: "rutracker.1", Cookies: "bb_session=old; cf_clearance=cf", UserAgent: goodUA}}}
	p := NewProvider(access, nil)

	_, err := p.SearchTorrents(context.Background(), model.SearchQuery{Query: "Матрица"})
	assert.ErrorIs(t, err, errBadAccount)
	assert.Equal(t, []string{"rutracker.1"}, access.marked)
}

func TestSearchTorrents_MalformedCookies(t *testing.T) {
	fakeSite(t)
	access := &fakeAccess{accounts: []model.Credentials{{AccountId: "rutracker.1", Cookies: " ;; ", UserAgent: goodUA}}}
	p := NewProvider(access, nil)

	_, err := p.SearchTorrents(context.Background(), model.SearchQuery{Query: "Матрица"})
	assert.ErrorIs(t, err, errBadAccount)
	assert.Equal(t, []string{"rutracker.1"}, access.marked)
}

func TestDownload_ChallengeMarksAccount(t *testing.T) {
	fakeSite(t)
	access := &fakeAccess{accounts: []model.Credentials{{AccountId: "rutracker.1", Cookies: goodCookies, UserAgent: goodUA}}}
	p := NewProvider(access, nil)

	result, err := p.SearchTorrents(context.Background(), model.SearchQuery{Query: "Матрица"})
	require.NoError(t, err)
	require.Len(t, result, 1)

	// cf_clearance протухает уже после поиска
	prev := baseURL
	broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cf-mitigated", "challenge")
		w.WriteHeader(http.StatusForbidden)
	}))
	defer broken.Close()
	baseURL = broken.URL
	defer func() { baseURL = prev }()

	_, err = result[0].Downloader(context.Background())
	assert.ErrorIs(t, err, ErrChallenge)
	assert.Equal(t, []string{"rutracker.1"}, access.marked)
}

func TestIsChallengePage(t *testing.T) {
	parse := func(html string) *goquery.Document {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		require.NoError(t, err)
		return doc
	}

	// страница, которую rutracker отдавал при проверке Cloudflare
	assert.True(t, isChallengePage(parse(`<html><head><title>Один момент…</title></head><body>
		<input type="hidden" name="cf-turnstile-response" id="cf-chl-widget-x_response"></body></html>`)))
	assert.True(t, isChallengePage(parse(`<html><head><title>Just a moment...</title></head></html>`)))
	assert.True(t, isChallengePage(parse(`<html><body><script src="/cdn-cgi/challenge-platform/h/b/orchestrate/chl_page/v1"></script></body></html>`)))

	assert.False(t, isChallengePage(parse(`<html><head><title>Вход</title></head><body><form id="login-form-full"></form></body></html>`)))
}

func TestParseCookies(t *testing.T) {
	cookies := parseCookies("bb_guid=a; bb_t=a%3A4%3A%7B%7D; cf_clearance=x-y_z")
	require.Len(t, cookies, 3)
	assert.Equal(t, "bb_guid", cookies[0].Name)
	assert.Equal(t, "a%3A4%3A%7B%7D", cookies[1].Value)
	assert.Equal(t, "x-y_z", cookies[2].Value)

	assert.Empty(t, parseCookies(""))
}
