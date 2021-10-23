package info_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"{{[ .Project ]}}/pkg/info"
	"{{[ .Project ]}}/pkg/logger"
	"{{[ .Project ]}}/pkg/version"
)

var ErrReturnError = errors.New("test of return Error")

func testHandler(
	t *testing.T, handler http.Handler, method, path string, code int, body string,
) {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), method, path, nil)
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
	t.Parallel()

	service := info.NewService(logger.New(new(logger.Config)))
	testHandler(
		t, service, "GET", "/healthz",
		http.StatusOK, "Ok",
	)
	testHandler(
		t, service, "GET", "/readyz",
		http.StatusOK, "Ok",
	)
	service.RegisterLivenessProbe(func() error {
		return ErrReturnError
	})
	service.RegisterReadinessProbe(func() error {
		return ErrReturnError
	})
	testHandler(
		t, service, "GET", "/healthz",
		http.StatusInternalServerError, ErrReturnError.Error()+"\n",
	)
	testHandler(
		t, service, "GET", "/readyz",
		http.StatusInternalServerError, ErrReturnError.Error()+"\n",
	)
}

func TestNotAllowed(t *testing.T) {
	t.Parallel()

	service := info.NewService(logger.New(new(logger.Config)))
	testHandler(
		t, service, "POST", "/",
		http.StatusMethodNotAllowed, "Only GET is allowed\n",
	)
}

func TestOptions(t *testing.T) {
	t.Parallel()

	service := info.NewService(logger.New(new(logger.Config)))
	testHandler(
		t, service, "OPTIONS", "/",
		http.StatusOK, "",
	)
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	service := info.NewService(logger.New(new(logger.Config)))
	testHandler(
		t, service, "GET", "/notfound",
		http.StatusNotFound, "404 page not found\n",
	)
	testHandler(
		t, service, "OPTIONS", "/notfound",
		http.StatusNotFound, "404 page not found\n",
	)
}

func TestAddHandler(t *testing.T) {
	t.Parallel()

	service := info.NewService(logger.New(new(logger.Config)))
	service.AddHandlerFunc(
		"/handler", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Handler"))
			if err != nil {
				t.Error(err)
			}
		},
	)
	testHandler(
		t, service, "GET", "/handler",
		http.StatusOK, "Handler",
	)
	testHandler(
		t, service, "OPTIONS", "/handler",
		http.StatusOK, "",
	)
}

func TestInfo(t *testing.T) {
	t.Parallel()

	service := info.NewService(logger.New(new(logger.Config)))
	data, err := json.Marshal(
		map[string]string{
			"version": version.RELEASE + "-" + version.COMMIT + "-" + version.BRANCH,
{{[- if .API.Enabled ]}}
			"API":     version.API,
{{[- end ]}}
			"date":    version.DATE,
			"repo":    version.REPO,
		},
	)

	if err != nil {
		t.Error(err)
	}

	testHandler(
		t, service, "GET", "/info",
		http.StatusOK, string(data),
	)
}

func TestRun(t *testing.T) {
	t.Parallel()

	service := info.NewService(logger.New(new(logger.Config)))

	if service.Run(":0") == nil {
		t.Error("Expected HTTP handler, got nil")
	}
}
