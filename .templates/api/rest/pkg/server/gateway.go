package server

import (
	"context"
	{{[- if not .API.Config.Insecure ]}}
	"crypto/tls"
	{{[- end ]}}
	"fmt"
	"net/http"

	{{[- if .Example ]}}

	"{{[ .Project ]}}/contracts/events"
	{{[- end ]}}
	"{{[ .Project ]}}/contracts/info"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	{{[- if .API.Config.Insecure ]}}
	"google.golang.org/grpc"
	{{[- else ]}}
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	forward := fmt.Sprintf("localhost:%d", gw.cfg.Port)

	// Register REST/gRPC gateway
	{{[- if .API.Config.Insecure ]}}
	opts := []grpc.DialOption{grpc.WithInsecure()}
	{{[- else ]}}
	opts := gw.TLSOptions()
	{{[- end ]}}
	gateway := runtime.NewServeMux(
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{EmitDefaults: true},
		),
	)
	// Register all gateways
	if err := info.RegisterInfoHandlerFromEndpoint(
		ctx, gateway, forward, opts,
	); err != nil {
		return err
	}
	{{[- if .Example ]}}

	if err := events.RegisterEventsHandlerFromEndpoint(
		ctx, gateway, forward, opts,
	); err != nil {
		return err
	}
	{{[- end ]}}
	
	return gw.Serve(gateway)
}

// Shutdown process graceful shutdown for the gateway server
func (gw GatewayServer) Shutdown(ctx context.Context) error {
	if gw.srv != nil {
		return gw.srv.Shutdown(ctx)
	}

	return nil
}

{{[- if not .API.Config.Insecure ]}}

// TLSOptions gives TLS secure/insecure option
func (gw GatewayServer) TLSOptions() []grpc.DialOption {
	options := []grpc.DialOption{}

	if gw.cfg.Insecure {
		return append(options, grpc.WithInsecure())
	}

	return append(options, grpc.WithTransportCredentials(credentials.NewTLS(
		&tls.Config{
			// nolint: gosec
			InsecureSkipVerify: true,
		},
	)))
}
{{[- end ]}}

// Serve prepares server and listen
func (gw *GatewayServer) Serve(handler http.Handler) error {
	// Create gateway server
	gw.srv = &http.Server{
		// Listening http -> gRPC address
		Addr: fmt.Sprintf(":%d", gw.cfg.Gateway.Port),
	}

	// Add gateway handler
	mux := http.NewServeMux()
{{[- if .API.Config.Insecure ]}}
	mux.Handle("/", handler)
	gw.srv.Handler = mux

	return gw.srv.ListenAndServe()
}
{{[- else ]}}

if gw.cfg.Insecure {
		mux.Handle("/", handler)
		gw.srv.Handler = mux

		return gw.srv.ListenAndServe()
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=15768000; includeSubDomains")
		handler.ServeHTTP(w, r)
	})

	gw.srv.Handler = mux
	gw.srv.TLSConfig = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		NextProtos:               []string{"h2"},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
	gw.srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))

	// Configure HTTP2 server
	if err := http2.ConfigureServer(gw.srv, nil); err != nil {
		return fmt.Errorf("failed to configure HTTP2 server: %s", err)
	}

	return gw.srv.ListenAndServeTLS(
		gw.cfg.Certificates.Crt,
		gw.cfg.Certificates.Key,
	)
}
{{[- end ]}}
