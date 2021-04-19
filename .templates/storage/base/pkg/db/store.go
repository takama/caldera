package db

import (
	"context"
{{[- if .Example ]}}

	"{{[ .Project ]}}/pkg/db/provider"
{{[- end ]}}
{{[- if .Prometheus.Enabled ]}}

	"{{[ .Project ]}}/pkg/metrics"
{{[- end ]}}
)

// Store design database interface with providers.
type Store interface {
	Check() error
	Shutdown(ctx context.Context) error
	{{[- if .Example ]}}
	EventsProvider() provider.Events
	{{[- end ]}}
	{{[- if .Prometheus.Enabled ]}}
	MetricFunc() metrics.MetricFunc
	{{[- end ]}}
}
