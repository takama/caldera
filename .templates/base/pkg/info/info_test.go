package info

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"{{[ .Project ]}}/pkg/logger"
	"{{[ .Project ]}}/pkg/version"
)

var ErrReturnError = errors.New("test of return Error")

func testHandler(
	t *testing.T, handler http.HandlerFunc, method, path string, code int, body string) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Error(err)
	}

	trw := httptest.NewRecorder()
	handler.ServeHTTP(trw, req)

	if trw.Code != code {
		t.Error("Expected status code:", code, "got", trw.Code)
	}
	if trw.Body.String() != body {
		t.Error("Expected body", body, "got", trw.Body.String())
	}
}

func TestProbe(t *testing.T) {
	service := NewService(logger.New(new(logger.Config)))
	testHandler(
		t, service.ServeHTTP, "GET", "/healthz",
		http.StatusOK, "Ok",
	)
	testHandler(
		t, service.ServeHTTP, "GET", "/readyz",
		http.StatusOK, "Ok",
	)
	service.RegisterLivenessProbe(func() error {
		return ErrReturnError
	})
	service.RegisterReadinessProbe(func() error {
		return ErrReturnError
	})
	testHandler(
		t, service.ServeHTTP, "GET", "/healthz",
		http.StatusInternalServerError, ErrReturnError.Error()+"\n",
	)
	testHandler(
		t, service.ServeHTTP, "GET", "/readyz",
		http.StatusInternalServerError, ErrReturnError.Error()+"\n",
	)
}

func TestNotAllowed(t *testing.T) {
	service := NewService(logger.New(new(logger.Config)))
	testHandler(
		t, service.ServeHTTP, "POST", "/",
		http.StatusMethodNotAllowed, "Only GET is allowed\n",
	)
}

func TestOptions(t *testing.T) {
	service := NewService(logger.New(new(logger.Config)))
	testHandler(
		t, service.ServeHTTP, "OPTIONS", "/",
		http.StatusOK, "",
	)
}

func TestNotFound(t *testing.T) {
	service := NewService(logger.New(new(logger.Config)))
	testHandler(
		t, service.ServeHTTP, "GET", "/notfound",
		http.StatusNotFound, "404 page not found\n",
	)
	testHandler(
		t, service.ServeHTTP, "OPTIONS", "/notfound",
		http.StatusNotFound, "404 page not found\n",
	)
}

func TestAddHandler(t *testing.T) {
	service := NewService(logger.New(new(logger.Config)))
	service.AddHandler(
		// nolint: unparam
		"/handler", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Handler"))
			if err != nil {
				t.Error(err)
			}
		},
	)
	testHandler(
		t, service.ServeHTTP, "GET", "/handler",
		http.StatusOK, "Handler",
	)
	testHandler(
		t, service.ServeHTTP, "OPTIONS", "/handler",
		http.StatusOK, "",
	)
}

func TestInfo(t *testing.T) {
	service := NewService(logger.New(new(logger.Config)))
	data, err := json.Marshal(
		map[string]string{
			"version": version.RELEASE + "-" + version.COMMIT + "-" + version.BRANCH,
			"date":    version.DATE,
			"repo":    version.REPO,
		},
	)

	if err != nil {
		t.Error(err)
	}

	testHandler(
		t, service.ServeHTTP, "GET", "/info",
		http.StatusOK, string(data),
	)
}

func TestRun(t *testing.T) {
	service := NewService(logger.New(new(logger.Config)))
	h := service.Run(":0")
	if h == nil {
		t.Error("Expected HTTP handler, got nil")
	}
}
