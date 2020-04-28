package db

import (
	"context"

	"{{[ .Project ]}}/pkg/db/provider"
)

// Store design database interface with providers.
type Store interface {
	Check() error
	Shutdown(ctx context.Context) error
	{{[- if .Example ]}}
	EventsProvider() provider.Events
	{{[- end ]}}
}
