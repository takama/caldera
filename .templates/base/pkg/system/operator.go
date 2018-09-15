package system

import "context"

// Operator defines reload, maintenance and shutdown interface
type Operator interface {
	Reload() error
	Maintenance() error
	Shutdown() []error
}

// Shutdowner defines Shutdown interface
type Shutdowner interface {
	Shutdown(context.Context) error
}

// Checker defines simple probe checkers
type Checker interface {
	Check() error
}
