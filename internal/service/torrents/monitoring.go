package torrents

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var linksRegisteredLinksGauge prometheus.Gauge

func init() {
	linksRegisteredLinksGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "rms",
		Subsystem: "media_discovery",
		Name:      "torrents_registered_links",
		Help:      "The count of registered download links",
	})
}
