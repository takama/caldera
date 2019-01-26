package system

import (
	"context"
	"errors"
	"time"
)

// ErrNotImplemented declares error for method that isn't implemented
var ErrNotImplemented = errors.New("this method is not implemented")

// ErrEmptyServerPointer declares error for nil pointer
var ErrEmptyServerPointer = errors.New("server pointer should not be nil")

// Operations implements simplest Operator interface
type Operations struct {
	shutdowns []Shutdowner
}

// NewOperator creates operator
func NewOperator(sd ...Shutdowner) *Operations {
	service := new(Operations)
	service.shutdowns = append(service.shutdowns, sd...)

	return service
}

// Reload operation implementation
func (o Operations) Reload() error {
	return ErrNotImplemented
}

// Maintenance operation implementation
func (o Operations) Maintenance() error {
	return ErrNotImplemented
}

// Shutdown operation
func (o Operations) Shutdown() []error {
	var errs []error
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	for _, fn := range o.shutdowns {
		if err := fn.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
