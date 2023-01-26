package users

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var userRequestsCounter *prometheus.CounterVec

func init() {
	userRequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "users_requests_by_user",
		Help:      "Amount of user's requests",
	}, []string{"key"})
}
