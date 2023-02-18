package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	limitReachedRequestsCounter prometheus.Counter
	searchDurationMetric        *prometheus.HistogramVec
)

func init() {
	limitReachedRequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "rms",
		Name:      "server_rejected_by_limit_requests",
		Help:      "The total number of rejected requests by rate limit",
	})

	searchDurationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "rms",
		Name:      "server_search_requests_duration",
		Help:      "Duration of search requests",
	}, []string{"service"})
}
