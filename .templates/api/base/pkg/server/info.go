package server

import (
	"context"

	"{{[ .Project ]}}/contracts/info"
	"{{[ .Project ]}}/pkg/version"

	"github.com/golang/protobuf/ptypes/empty"
)

type infoServer struct {
}

// GetInfo returns the server information.
func (is infoServer) GetInfo(
	ctx context.Context,
	empty *empty.Empty,
) (*info.Report, error) {
	return &info.Report{
		Version:     version.RELEASE + "-" + version.COMMIT + "-" + version.BRANCH,
		API:         version.API,
		Date:        version.DATE,
		Repository:  version.REPO,
		Description: version.DESC,
	}, nil
}

// AuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`.
func (is infoServer) AuthFuncOverride(
	ctx context.Context,
	fullMethodName string,
) (context.Context, error) {
	return ctx, nil
}
