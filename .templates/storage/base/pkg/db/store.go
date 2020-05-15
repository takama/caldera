package db

import (
	"context"
{{[- if .Example ]}}

	"{{[ .Project ]}}/pkg/db/provider"
{{[- end ]}}
)

// Store design database interface with providers.
type Store interface {
	Check() error
	Shutdown(ctx context.Context) error
	{{[- if .Example ]}}
	EventsProvider() provider.Events
	{{[- end ]}}
}
