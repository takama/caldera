package system

import (
	"net/http"
	"testing"
)

func TestStubHandling(t *testing.T) {
	operator := NewOperator(&http.Server{})
	err := operator.Reload()
	if err != ErrNotImplemented {
		t.Error("Expected error", ErrNotImplemented, "got", err)
	}
	err = operator.Maintenance()
	if err != ErrNotImplemented {
		t.Error("Expected error", ErrNotImplemented, "got", err)
	}
	errs := operator.Shutdown()
	if len(errs) > 0 {
		t.Error("Expected success, got errors", errs)
	}
}
