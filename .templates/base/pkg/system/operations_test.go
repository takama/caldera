package system_test

import (
	"errors"
	"net/http"
	"testing"

	"{{[ .Project ]}}/pkg/system"
)

func TestStubHandling(t *testing.T) {
	t.Parallel()

	operator := system.NewOperator(&http.Server{})

	if err := operator.Reload(); !errors.Is(err, system.ErrNotImplemented) {
		t.Error("Expected error", system.ErrNotImplemented, "got", err)
	}

	if err := operator.Maintenance(); !errors.Is(err, system.ErrNotImplemented) {
		t.Error("Expected error", system.ErrNotImplemented, "got", err)
	}

	if errs := operator.Shutdown(); len(errs) > 0 {
		t.Error("Expected success, got errors", errs)
	}
}
