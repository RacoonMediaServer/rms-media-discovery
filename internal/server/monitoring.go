package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalRequestsCount *prometheus.CounterVec
	panicCnt           prometheus.Counter
	limitReachedCnt    prometheus.Counter
	unauthorizedReqCnt prometheus.Counter
)

func init() {
	totalRequestsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "server_requests",
		Help:      "The total number of processed requests",
	}, []string{"method", "code"})

	panicCnt = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "panics",
		Help:      "The total number of panic happens",
	})

	limitReachedCnt = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "server_rejected_by_limit_requests",
		Help:      "The total number of rejected requests by rate limit",
	})

	unauthorizedReqCnt = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "server_unauthorized_requests",
		Help:      "The total number of unauthorized requests",
	})
}
