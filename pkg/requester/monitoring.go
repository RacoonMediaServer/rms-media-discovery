package requester

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OutgoingRequestsCounter *prometheus.CounterVec
	OutgoingRequestsMetric  *prometheus.HistogramVec
)

func init() {
	OutgoingRequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "outgoing_requests_count",
		Help:      "Total amount of outgoing requests",
	}, []string{"code", "service"})

	OutgoingRequestsMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "outgoing_request_duration",
		Help:      "Duration of outgoing requests",
	}, []string{"service"})
}
