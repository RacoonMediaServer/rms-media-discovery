package utils

import (
	"context"
	"github.com/gocolly/colly/v2"
	"net/http"
	"time"
)

type contextTransport struct {
	ctx       context.Context
	transport *http.Transport
}

func (t *contextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.WithContext(t.ctx)
	return t.transport.RoundTrip(req)
}

func CollyWithContext(c *colly.Collector, ctx context.Context) {
	c.OnRequest(func(req *colly.Request) {
		select {
		case <-ctx.Done():
			req.Abort()
		default:
		}
	})

	transport := &contextTransport{
		ctx:       ctx,
		transport: &http.Transport{},
	}
	c.SetClient(&http.Client{
		Transport: transport,
		Timeout:   2 * time.Minute,
	})
}
