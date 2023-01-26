package requester

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	outgoingRequestsCounter *prometheus.CounterVec
	outgoingRequestsMetric  *prometheus.HistogramVec
)

func init() {
	outgoingRequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "outgoing_requests_count",
		Help:      "Total amount of outgoing requests",
	}, []string{"code", "service"})

	outgoingRequestsMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "outgoing_request_duration",
		Help:      "Duration of outgoing requests",
	}, []string{"service"})
}
