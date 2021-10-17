package server

import (
	"context"
	"fmt"
	{{[- if .API.UI ]}}
	"mime"
	{{[- end ]}}
	"net/http"
	{{[- if .API.UI ]}}
	"strings"
	{{[- end ]}}

	{{[- if .Example ]}}

	"{{[ .Project ]}}/contracts/events"
	{{[- end ]}}
	"{{[ .Project ]}}/contracts/info"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	{{[- if .API.UI ]}}
	"github.com/rakyll/statik/fs"
	{{[- end ]}}
	"go.uber.org/zap"
	"google.golang.org/grpc"
	{{[- if .API.UI ]}}

	// OpenApi UI files.
	public "{{[ .Project ]}}/public/openapi"
	{{[- end ]}}
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

	if err := events.RegisterPublicHandlerFromEndpoint(
		ctx, gateway, forward, opts,
	); err != nil {
		return fmt.Errorf("failed to register public handler: %w", err)
	}
	{{[- end ]}}

	return gw.Serve(gateway{{[- if .API.UI ]}}, getOpenAPIHandler(){{[- end ]}})
}

// Shutdown process graceful shutdown for the gateway server.
func (gw GatewayServer) Shutdown(ctx context.Context) error {
	if gw.srv != nil {
		if err := gw.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown gateway: %w", err)
		}
	}

	return nil
}

// Serve prepares server and listen.
func (gw *GatewayServer) Serve(handler{{[- if .API.UI ]}}, openapi{{[- end ]}} http.Handler) error {
	// Create gateway server
	gw.srv = &http.Server{
		// Listening http -> gRPC address.
		Addr: fmt.Sprintf(":%d", gw.cfg.Gateway.Port),
	}

	// Add gateway handler.
	{{[- if .API.UI ]}}
	gw.srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/"+version.API) {
			handler.ServeHTTP(w, r)

			return
		}
		openapi.ServeHTTP(w, r)
	})
	{{[- else ]}}
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	gw.srv.Handler = mux
	{{[- end ]}}

	if err := gw.srv.ListenAndServe(); err != nil{
		return fmt.Errorf("failed to serve gateway: %w", err)
	}

	return nil
}

{{[- if .API.UI ]}}

// getOpenAPIHandler serves an OpenAPI UI for public namespace.
func getOpenAPIHandler() http.Handler {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		// Panic since this is a permanent error.
		panic("creating mime: " + err.Error())
	}

	sfs, err := fs.NewWithNamespace(public.Public)
	if err != nil {
		// Panic since this is a permanent error.
		panic("creating OpenAPI filesystem: " + err.Error())
	}

	return http.FileServer(sfs)
}
{{[- end ]}}
