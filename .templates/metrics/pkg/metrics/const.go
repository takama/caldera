// nolint: gochecknoglobals
package metrics

import (
	"runtime"
	"time"

	"{{[ .Project ]}}/pkg/version"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	buildTimestamp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "build_timestamp",
			Help: "Build timestamp with additional info.",
		},
		[]string{"go_version", "version", "revision"},
	)
)

func registerConstMetrics() {
	time, err := time.Parse("2006-01-02T15:04:05+07", version.DATE)
	if err != nil {
		return
	}

	buildTimestamp.With(prometheus.Labels{
		"go_version": runtime.Version(),
		"version":    version.RELEASE,
		"revision":   version.COMMIT,
	}).Set(float64(time.Unix()))

	prometheus.MustRegister(buildTimestamp)
}
