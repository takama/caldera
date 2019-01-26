package provider

import "errors"

var (
	// ErrNotDefinedTransaction defines error for nilable transaction
	ErrNotDefinedTransaction = errors.New("transaction does not defined")
	{{[- if .Contract ]}}
	// ErrNotDefinedID defines error when ID does not defined
	ErrNotDefinedID = errors.New("ID does not defined")
	// ErrNotDefinedName defines error when Name does not defined
	ErrNotDefinedName = errors.New("name does not defined")
	// ErrAlreadyExistingID defines error for existing ID
	ErrAlreadyExistingID = errors.New("ID already exists in database")
	// ErrNotExistingEvent defines error for absent Event
	ErrNotExistingEvent = errors.New("event does not exist in database")
	{{[- end ]}}
)
