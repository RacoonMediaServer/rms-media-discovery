package scraper

import (
	"fmt"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/requester"
	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *scraper) getContextLog(ctx *colly.Context) *log.Entry {
	l, ok := ctx.GetAny("log").(*log.Entry)
	if !ok {
		l = log.WithField("from", s.service)
	}
	return l
}

func (s *scraper) setCallbacks() {
	s.c.OnRequest(func(r *colly.Request) {
		select {
		case <-s.ctx.Done():
			r.Abort()
		default:
		}

		l := utils.LogFromContext(s.ctx, s.service).WithField("url", r.URL.String())
		l.Debugf("Fetching...")
		timer := prometheus.NewTimer(requester.OutgoingRequestsMetric.WithLabelValues(s.service))
		r.Ctx.Put("timer", timer)
		r.Ctx.Put("log", l)
	})

	s.c.OnError(func(response *colly.Response, err error) {
		l := s.getContextLog(response.Ctx)
		timer, ok := response.Ctx.GetAny("timer").(*prometheus.Timer)
		if !ok {
			l.Error("Cannot extract timer!")
			return
		}
		timer.ObserveDuration()

		l.Errorf("Request failed: %s", err)
	})

	s.c.OnResponse(func(response *colly.Response) {
		l := s.getContextLog(response.Ctx)
		timer, ok := response.Ctx.GetAny("timer").(*prometheus.Timer)
		if !ok {
			l.Error("Cannot extract timer!")
			return
		}
		timer.ObserveDuration()

		l.Debugf("Got response: %d", response.StatusCode)
		requester.OutgoingRequestsCounter.WithLabelValues(fmt.Sprintf("%d", response.StatusCode), s.service).Inc()
	})
}
