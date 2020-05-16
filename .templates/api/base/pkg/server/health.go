package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Server for the Health Check gRPC API.
type healthServer struct{}

// Check is used for health checks.
func (hs *healthServer) Check(
	ctx context.Context,
	in *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	// This is where we can implement checks of our service status
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is not implemented.
func (hs *healthServer) Watch(
	in *grpc_health_v1.HealthCheckRequest,
	srv grpc_health_v1.Health_WatchServer,
) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}
