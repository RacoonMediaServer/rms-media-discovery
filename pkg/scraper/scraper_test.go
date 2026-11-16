package scraper

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
)

func newTestServer(gotUA *string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*gotUA = r.UserAgent()
		if r.URL.Path == "/challenge" {
			w.Header().Set("cf-mitigated", "challenge")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, _ = w.Write([]byte("<html><body><b id='x'>hello</b></body></html>"))
	}))
}

func TestScraper_StrictModeReturnsStatusError(t *testing.T) {
	var ua string
	srv := newTestServer(&ua)
	defer srv.Close()

	c := NewWithOptions("test", Options{UserAgent: "my-agent", FailOnRequestError: true})
	c.SetContext(context.Background())
	err := c.Select("#x", func(e *colly.HTMLElement, _ interface{}) {}).Get(srv.URL + "/challenge")

	var statusErr *requester.StatusError
	assert.True(t, errors.As(err, &statusErr), "got %v", err)
	assert.Equal(t, 403, statusErr.Code)
	assert.True(t, statusErr.IsChallenge())
	assert.Equal(t, "my-agent", ua)

	// ошибка не должна "залипать" на следующих запросах
	found := false
	err = c.Select("#x", func(e *colly.HTMLElement, _ interface{}) { found = true }).Get(srv.URL + "/ok")
	assert.NoError(t, err)
	assert.True(t, found)
}

func TestScraper_DefaultModeKeepsOldBehavior(t *testing.T) {
	var ua string
	srv := newTestServer(&ua)
	defer srv.Close()

	// без FailOnRequestError ошибка HTTP только логируется - остальные провайдеры (rutor, anidub) на это полагаются
	c := New("test")
	c.SetContext(context.Background())
	err := c.Select("#x", func(e *colly.HTMLElement, _ interface{}) {}).Get(srv.URL + "/challenge")
	assert.NoError(t, err)
	assert.Contains(t, ua, "Chrome/108")
}
