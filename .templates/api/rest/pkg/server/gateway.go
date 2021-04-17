package server

import (
	"context"
	"fmt"
	"net/http"

	{{[- if .Example ]}}

	"{{[ .Project ]}}/contracts/events"
	{{[- end ]}}
	"{{[ .Project ]}}/contracts/info"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GatewayServer contains gateway functionality of the service.
type GatewayServer struct {
	cfg *Config
	log *zap.Logger
	srv *http.Server
}

// NewGateway creates a new gateway server.
func NewGateway(ctx context.Context, cfg *Config, log *zap.Logger) (*GatewayServer, error) {
	return &GatewayServer{
		cfg: cfg,
		log: log,
	}, nil
}

// LivenessProbe returns liveness probe of the server.
func (gw GatewayServer) LivenessProbe() error {
	return nil
}

// ReadinessProbe returns readiness probe for the server.
func (gw GatewayServer) ReadinessProbe() error {
	return nil
}

// Run starts the gateway server.
func (gw *GatewayServer) Run(ctx context.Context) error {
	forward := fmt.Sprintf("localhost:%d", gw.cfg.Port)

	// Register REST/gRPC gateway
	opts := []grpc.DialOption{grpc.WithInsecure()}
	gateway := runtime.NewServeMux()

	// Register all gateways
	if err := info.RegisterInfoHandlerFromEndpoint(
		ctx, gateway, forward, opts,
	); err != nil {
		return fmt.Errorf("failed to register info handler: %w", err)
	}
	{{[- if .Example ]}}

	if err := events.RegisterEventsHandlerFromEndpoint(
		ctx, gateway, forward, opts,
	); err != nil {
		return err
	}
	{{[- end ]}}

	return gw.Serve(cors.Default().Handler(gateway))
}

// Shutdown process graceful shutdown for the gateway server.
func (gw GatewayServer) Shutdown(ctx context.Context) error {
	if gw.srv != nil {
		return gw.srv.Shutdown(ctx)
	}

	return nil
}

// Serve prepares server and listen.
func (gw *GatewayServer) Serve(handler http.Handler) error {
	// Create gateway server
	gw.srv = &http.Server{
		// Listening http -> gRPC address.
		Addr: fmt.Sprintf(":%d", gw.cfg.Gateway.Port),
	}

	// Add gateway handler.
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	gw.srv.Handler = mux

	return gw.srv.ListenAndServe()
}
