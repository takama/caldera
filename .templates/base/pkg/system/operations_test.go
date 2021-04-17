package system_test

import (
	"errors"
	"net/http"
	"testing"

	"{{[ .Project ]}}/pkg/system"
)

func TestStubHandling(t *testing.T) {
	operator := system.NewOperator(&http.Server{})
	err := operator.Reload()

	if !errors.Is(err, system.ErrNotImplemented) {
		t.Error("Expected error", system.ErrNotImplemented, "got", err)
	}

	err = operator.Maintenance()

	if !errors.Is(err, system.ErrNotImplemented) {
		t.Error("Expected error", system.ErrNotImplemented, "got", err)
	}

	errs := operator.Shutdown()

	if len(errs) > 0 {
		t.Error("Expected success, got errors", errs)
	}
}
