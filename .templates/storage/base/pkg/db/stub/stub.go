package stub

import (
	"context"
	{{[- if .Prometheus.Enabled ]}}
	"database/sql"
	{{[- end ]}}

	"{{[ .Project ]}}/pkg/db"
	{{[- if .Example ]}}
	"{{[ .Project ]}}/pkg/db/provider"
	{{[- end ]}}
	{{[- if .Prometheus.Enabled ]}}
	"{{[ .Project ]}}/pkg/metrics"
	{{[- end ]}}

	"go.uber.org/zap"
)

const (
	// Driver defines database driver name.
	Driver = "stub"
)

// Stub contains stub configuration.
type Stub struct {
	cfg *db.Config
	log *zap.Logger
	{{[- if .Example ]}}
	// Contract providers
	events *eventsProvider
	{{[- end ]}}
}

// New creates new DB connection.
func New(cfg *db.Config, log *zap.Logger) (*Stub, error) {
	log.Info("DB", zap.String("driver", Driver))

	conn := &Stub{
		cfg: cfg,
		log: log,
	}

	{{[- if .Example ]}}

	conn.events = &eventsProvider{cfg: cfg}

	if err := conn.events.load(); err != nil {
		return nil, err
	}
	{{[- end ]}}

	return conn, nil
}

// Check DB.
func (s Stub) Check() error {
	return nil
}

// Shutdown process graceful shutdown for the storage.
func (s Stub) Shutdown(ctx context.Context) error {
	s.log.Debug("DB closed", zap.String("driver", Driver))
	return nil
}

{{[- if .Example ]}}

// EventsProvider returns data store provider for Events.
func (s Stub) EventsProvider() provider.Events {
	return s.events
}
{{[- end ]}}

{{[- if .Prometheus.Enabled ]}}

// MetricFunc returns a "stub" func to monitor connectivity to stub.
func (s Stub) MetricFunc() metrics.MetricFunc {
	return metrics.DBMetricFunc("", "", sql.DBStats{})
}
{{[- end ]}}
