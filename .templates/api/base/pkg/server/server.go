package server

import (
	"context"
	{{[- if not .API.Config.Insecure ]}}
	"crypto/tls"
	{{[- end ]}}
	"fmt"
	"net"

	{{[- if .Example ]}}

	"{{[ .Project ]}}/contracts/events"
	"{{[ .Project ]}}/pkg/db/provider"
	{{[- end ]}}

	"go.uber.org/zap"
	"google.golang.org/grpc"
	{{[- if not .API.Config.Insecure ]}}
	"google.golang.org/grpc/credentials"
	{{[- end ]}}
)

// Server contains core functionality of the service
type Server struct {
	cfg *Config
	log *zap.Logger
	srv *grpc.Server
	is  *infoServer
	{{[- if .Example ]}}
	// Contract provider servers
	es *eventsServer
	{{[- end ]}}
}

// New creates a new core server
func New(ctx context.Context, cfg *Config, log *zap.Logger) (*Server, error) {
	return &Server{
		cfg: cfg,
		log: log,
		is:  new(infoServer),
		{{[- if .Example ]}}
		es:  new(eventsServer),
		{{[- end ]}}
	}, nil
}

{{[- if .Example ]}}

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
	{{[- if .Example ]}}
	if err := s.checkProviders(); err != nil {
		return err
	}
	{{[- end ]}}
	// Register gRPC server
	{{[- if .API.Config.Insecure ]}}
	s.srv = grpc.NewServer()
	{{[- else ]}}
	s.srv = grpc.NewServer(s.ServerOptions()...)
	{{[- end ]}}
	info.RegisterInfoServer(s.srv, s.is)
	{{[- if .Example ]}}
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

{{[- if .Example ]}}

func (s Server) checkProviders() error {
	if s.es.Events == nil {
		return ErrEventsProviderEmpty
	}

	return nil
}
{{[- end ]}}

{{[- if not .API.Config.Insecure ]}}

// ServerOptions gives server authentication and secure/insecure options
func (s Server) ServerOptions() []grpc.ServerOption {
	options := []grpc.ServerOption{}

	if !s.cfg.Insecure {
		cert, err := tls.LoadX509KeyPair(s.cfg.Certificates.Crt, s.cfg.Certificates.Key)
		if err != nil {
			s.log.Fatal("Failed to load key pair", zap.Error(err))
		}
		// Enable TLS for all incoming connections.
		options = append(options, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	}

	return options
}
{{[- end ]}}
