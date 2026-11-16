package scraper

import (
	"context"
	"net/http"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/gocolly/colly/v2"
)

// contextTransport привязывает все запросы коллектора к контексту
type contextTransport struct {
	ctx context.Context
}

func (t *contextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.WithContext(t.ctx)
	// Берем http.DefaultTransport при каждом запросе, а не создаем свой http.Transport: иначе скрапер не видит настроек
	// процесса (выбор IPv6, HTTP(S)_PROXY, таймауты dial/TLS/idle), которыми пользуется остальной код (requester)
	return http.DefaultTransport.RoundTrip(req)
}

func setCollyContext(c *colly.Collector, ctx context.Context) {
	c.SetClient(&http.Client{
		Transport: &contextTransport{ctx: ctx},
		Timeout:   requester.Timeout,
	})
}
