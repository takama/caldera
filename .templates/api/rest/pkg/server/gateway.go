package server

import (
	"context"
	"fmt"
	"net/http"

	{{[- if .Contract ]}}

	"{{[ .Project ]}}/contracts/events"
	{{[- end ]}}

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	{{[- if .Contract ]}}
	"google.golang.org/grpc"
	{{[- end ]}}
)

// GatewayServer contains gateway functionality of the service
type GatewayServer struct {
	cfg *Config
	log *zap.Logger
	srv *http.Server
}

// NewGateway creates a new gateway server
func NewGateway(ctx context.Context, cfg *Config, log *zap.Logger) (*GatewayServer, error) {
	return &GatewayServer{
		cfg: cfg,
		log: log,
	}, nil
}

// LivenessProbe returns liveness probe of the server
func (gw GatewayServer) LivenessProbe() error {
	return nil
}

// ReadinessProbe returns readiness probe for the server
func (gw GatewayServer) ReadinessProbe() error {
	return nil
}

// Run starts the gateway server
func (gw *GatewayServer) Run(ctx context.Context) error {
	// Listening http -> gRPC address
	addr := fmt.Sprintf(":%d", gw.cfg.Gateway.Port)

	// Register REST/gRPC gateway
	gateway := runtime.NewServeMux()
	{{[- if .Contract ]}}
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := events.RegisterEventsHandlerFromEndpoint(
		ctx, gateway, fmt.Sprintf("localhost:%d", gw.cfg.Port), opts,
	); err != nil {
		return err
	}
	{{[- end ]}}

	// Add gateway handler
	mux := http.NewServeMux()
	mux.Handle("/", gateway)

	// Create gateway server
	gw.srv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return gw.srv.ListenAndServe()
}

// Shutdown process graceful shutdown for the gateway server
func (gw GatewayServer) Shutdown(ctx context.Context) error {
	if gw.srv != nil {
		return gw.srv.Shutdown(ctx)
	}

	return nil
}
