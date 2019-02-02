package server

import (
	"context"

	"{{[ .Project ]}}/contracts/info"
	"{{[ .Project ]}}/pkg/version"

	"github.com/golang/protobuf/ptypes/empty"
)

type infoServer struct {
}

// GetInfo returns the server information
func (is infoServer) GetInfo(
	ctx context.Context,
	empty *empty.Empty,
) (*info.Report, error) {
	return &info.Report{
		Version: version.RELEASE + "-" + version.COMMIT + "-" + version.BRANCH,
		Date:    version.DATE,
		Repo:    version.REPO,
	}, nil
}

// GetHealth returns the server health information
func (is infoServer) GetHealth(
	ctx context.Context,
	empty *empty.Empty,
) (*info.Health, error) {
	return &info.Health{
		Alive: true,
	}, nil
}

// AuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`
func (is infoServer) AuthFuncOverride(
	ctx context.Context,
	fullMethodName string,
) (context.Context, error) {
	return ctx, nil
}
