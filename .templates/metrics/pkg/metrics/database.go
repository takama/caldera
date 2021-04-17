// nolint: gochecknoglobals
package metrics

{{[- if .Storage.Enabled ]}}
import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	maxOpenDBConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "max_open_db_connections",
			Help: "Max open connections to database.",
		},
		[]string{"host", "database"},
	)
	openDBConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "open_db_connections",
			Help: "Open connections to database.",
		},
		[]string{"host", "database"},
	)
	inUseDBConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "in_use_db_connections",
			Help: "In use connections to database.",
		},
		[]string{"host", "database"},
	)
	idleDBConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "idle_db_connections",
			Help: "Idle connections to database.",
		},
		[]string{"host", "database"},
	)
	waitDBConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wait_db_connections",
			Help: "Wait connections to database.",
		},
		[]string{"host", "database"},
	)
)

// DBMetricFunc returns function that updates database-related metrics.
func DBMetricFunc(
	host string,
	database string,
	stats sql.DBStats,
) MetricFunc {
	return func() error {
		labels := prometheus.Labels{
			"host":     host,
			"database": database,
		}

		maxOpenDBConnections.With(labels).Set(float64(stats.MaxOpenConnections))
		openDBConnections.With(labels).Set(float64(stats.OpenConnections))
		inUseDBConnections.With(labels).Set(float64(stats.InUse))
		idleDBConnections.With(labels).Set(float64(stats.Idle))
		waitDBConnections.With(labels).Set(float64(stats.WaitCount))

		return nil
	}
}

func registerDatabaseMetrics() {
	prometheus.MustRegister(
		maxOpenDBConnections,
		openDBConnections,
		inUseDBConnections,
		idleDBConnections,
		waitDBConnections,
	)
}
{{[- else ]}}

// no database storage to expose metrics
{{[- end ]}}
