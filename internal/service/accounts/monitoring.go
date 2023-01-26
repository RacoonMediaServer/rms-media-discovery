package accounts

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	unaccessibleAccountsGauge *prometheus.GaugeVec
)

func init() {
	unaccessibleAccountsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "accounts_unaccessible_count",
		Help:      "The count of unaccessible accounts",
	}, []string{"service"})
}
