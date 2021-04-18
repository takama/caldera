package server

import (
	"context"
	"fmt"
	"net"

	"{{[ .Project ]}}/contracts/info"
	{{[- if .Example ]}}
	"{{[ .Project ]}}/contracts/events"
	"{{[ .Project ]}}/pkg/db/provider"
	{{[- end ]}}

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server contains core functionality of the service.
type Server struct {
	cfg *Config
	log *zap.Logger
	srv *grpc.Server
	hs  *healthServer
	is  *infoServer
	{{[- if .Example ]}}
	// Contract provider servers.
	es *eventsServer
	{{[- end ]}}
}

// New creates a new core server.
func New(ctx context.Context, cfg *Config, log *zap.Logger) (*Server, error) {
	return &Server{
		cfg: cfg,
		log: log,
		hs:  new(healthServer),
		is:  new(infoServer),
		{{[- if .Example ]}}
		es:  new(eventsServer),
		{{[- end ]}}
	}, nil
}

{{[- if .Example ]}}

// RegisterEventsProvider assign data store provider for Events.
func (s *Server) RegisterEventsProvider(provider provider.Events) {
	if provider != nil {
		s.es.Events = provider
	}
}
{{[- end ]}}

// LivenessProbe returns liveness probe of the server.
func (s Server) LivenessProbe() error {
	return nil
}

// ReadinessProbe returns readiness probe for the server.
func (s Server) ReadinessProbe() error {
	return nil
}

// Run starts the server.
func (s *Server) Run(ctx context.Context) error {
	{{[- if .Example ]}}
	if err := s.checkProviders(); err != nil {
		return err
	}

	{{[- end ]}}
	// Register gRPC server
	s.srv = grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(s.srv, s.hs)
	info.RegisterInfoServer(s.srv, s.is)
	{{[- if .Example ]}}
	events.RegisterPublicServer(s.srv, s.es)
	{{[- end ]}}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to create server listener: %w", err)
	}

	return s.srv.Serve(listener)
}

// Shutdown process graceful shutdown for the server.
func (s Server) Shutdown(ctx context.Context) error {
	if s.srv != nil {
		s.srv.GracefulStop()
	}

	return nil
}

{{[- if .Example ]}}

func (s Server) checkProviders() error {
	if s.es.Events == nil {
		return ErrEventsProviderEmpty
	}

	return nil
}
{{[- end ]}}
