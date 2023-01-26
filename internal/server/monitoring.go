package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalRequestsCounter        *prometheus.CounterVec
	panicCounter                prometheus.Counter
	limitReachedRequestsCounter prometheus.Counter
	unauthorizedRequstsCounter  prometheus.Counter
	searchDurationMetric        *prometheus.HistogramVec
)

func init() {
	totalRequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "server_requests",
		Help:      "The total number of processed requests",
	}, []string{"method", "code"})

	panicCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "panics",
		Help:      "The total number of panic happens",
	})

	limitReachedRequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "server_rejected_by_limit_requests",
		Help:      "The total number of rejected requests by rate limit",
	})

	unauthorizedRequstsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "server_unauthorized_requests",
		Help:      "The total number of unauthorized requests",
	})

	searchDurationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "server_search_requests_duration",
		Help:      "Duration of search requests",
	}, []string{"service"})
}
