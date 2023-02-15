package scraper

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/gocolly/colly/v2"
	"net/http"
)

type contextTransport struct {
	ctx       context.Context
	transport *http.Transport
}

func (t *contextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.WithContext(t.ctx)
	return t.transport.RoundTrip(req)
}

func setCollyContext(c *colly.Collector, ctx context.Context) {
	transport := &contextTransport{
		ctx:       ctx,
		transport: &http.Transport{},
	}
	c.SetClient(&http.Client{
		Transport: transport,
		Timeout:   requester.Timeout,
	})
}
