package info

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"{{[ .Project ]}}/pkg/version"

	"go.uber.org/zap"
)

// Config contains params to setup info handler
type Config struct {
	Port       int
	Statistics bool
}

// ProbeChecker defines simple function for probe checks
type ProbeChecker func() error

// Service contains info/health-check functionality
type Service struct {
	handlers        map[string]http.HandlerFunc
	logger          *zap.Logger
	livenessProbes  []ProbeChecker
	readinessProbes []ProbeChecker
}

// NewService creates new service with info/health-check handlers
func NewService(logger *zap.Logger) *Service {
	return &Service{
		logger:   logger,
		handlers: make(map[string]http.HandlerFunc),
	}
}

// AddHandler adds new handler with given path to info service
func (s *Service) AddHandler(path string, handler http.HandlerFunc) {
	s.handlers[path] = handler
}

// RegisterLivenessProbe defines liveness probe function
func (s *Service) RegisterLivenessProbe(checker ProbeChecker) {
	s.livenessProbes = append(s.livenessProbes, checker)
}

// RegisterReadinessProbe defines readiness probe function
func (s *Service) RegisterReadinessProbe(checker ProbeChecker) {
	s.readinessProbes = append(s.readinessProbes, checker)
}

// Run info/health-check service
func (s *Service) Run(addr string) *http.Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: s,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// Check for known errors
			if err != context.DeadlineExceeded &&
				err != context.Canceled &&
				err != http.ErrServerClosed {
				s.logger.Fatal(err.Error())
			}
			s.logger.Warn(err.Error())
		}
	}()

	return srv
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := strings.TrimRight(r.URL.Path, " /")
	switch r.Method {
	case "GET":
		switch route {
		case "", "/info":
			s.info(w)
		case "/healthz":
			s.liveness(w)
		case "/readyz":
			s.readiness(w)
		default:
			for path, handler := range s.handlers {
				if route == path {
					handler(w, r)
					return
				}
			}
			http.NotFound(w, r)
		}
	default:
		if r.Method == "OPTIONS" {
			switch route {
			case "", "/info", "/healthz", "/readyz":
				w.Header().Set("Allow", "GET")
			default:
				for path := range s.handlers {
					if route == path {
						w.Header().Set("Allow", "GET")
						return
					}
				}
				http.NotFound(w, r)
			}
			return
		}
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Service) info(w http.ResponseWriter) {
	data, err := json.Marshal(
		map[string]string{
			"version": version.RELEASE + "-" + version.COMMIT + "-" + version.BRANCH,
			"date":    version.DATE,
			"repo":    version.REPO,
		},
	)

	if err != nil {
		s.writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Service) liveness(w http.ResponseWriter) {
	for _, checker := range s.livenessProbes {
		if err := checker(); err != nil {
			s.writeError(w, err)
			return
		}
	}
	if _, err := io.WriteString(w, "Ok"); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Service) readiness(w http.ResponseWriter) {
	for _, checker := range s.readinessProbes {
		if err := checker(); err != nil {
			s.writeError(w, err)
			return
		}
	}
	if _, err := io.WriteString(w, "Ok"); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Service) writeError(w http.ResponseWriter, err error) {
	message := err.Error()
	s.logger.Error(message)
	http.Error(w, message, http.StatusInternalServerError)
}
