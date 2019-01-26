package server

import "errors"

var (
	// ErrEventsProviderEmpty defines error for nilable Events provider
	ErrEventsProviderEmpty = errors.New("events provider should be registered before running the service")
)
