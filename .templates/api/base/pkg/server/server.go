package server

import (
	"context"
	"fmt"
	"net"

	{{[- if .Contract ]}}

	"{{[ .Project ]}}/contracts/events"
	"{{[ .Project ]}}/pkg/db/provider"
	{{[- end ]}}

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server contains core functionality of the service
type Server struct {
	cfg *Config
	log *zap.Logger
	srv *grpc.Server
	{{[- if .Contract ]}}
	// Contract provider servers
	es *eventsServer
	{{[- end ]}}
}

// New creates a new core server
func New(ctx context.Context, cfg *Config, log *zap.Logger) (*Server, error) {
	return &Server{
		cfg: cfg,
		log: log,
		{{[- if .Contract ]}}
		es:  new(eventsServer),
		{{[- end ]}}
	}, nil
}

{{[- if .Contract ]}}

// RegisterEventsProvider assign data store provider for Events
func (s *Server) RegisterEventsProvider(provider provider.Events) {
	if provider != nil {
		s.es.Events = provider
	}
}
{{[- end ]}}

// LivenessProbe returns liveness probe of the server
func (s Server) LivenessProbe() error {
	return nil
}

// ReadinessProbe returns readiness probe for the server
func (s Server) ReadinessProbe() error {
	return nil
}

// Run starts the server
func (s *Server) Run(ctx context.Context) error {
	{{[- if .Contract ]}}
	if err := s.checkProviders(); err != nil {
		return err
	}
	{{[- end ]}}

	// Register gRPC server
	s.srv = grpc.NewServer()
	{{[- if .Contract ]}}
	events.RegisterEventsServer(s.srv, s.es)
	{{[- end ]}}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return err
	}

	return s.srv.Serve(listener)
}

// Shutdown process graceful shutdown for the server
func (s Server) Shutdown(ctx context.Context) error {
	if s.srv != nil {
		s.srv.GracefulStop()
	}

	return nil
}

{{[- if .Contract ]}}

func (s Server) checkProviders() error {
	if s.es.Events == nil {
		return ErrEventsProviderEmpty
	}

	return nil
}
{{[- end ]}}
