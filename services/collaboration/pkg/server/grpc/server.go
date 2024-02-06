package grpc

import (
	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	svc "github.com/owncloud/ocis/v2/services/collaboration/pkg/service/grpc/v0"
	"google.golang.org/grpc"
)

// Server initializes a new grpc service ready to run
// THIS SERVICE IS REGISTERED AGAINST REVA, NOT GO-MICRO
func Server(opts ...Option) (*grpc.Server, func(), error) {
	grpcOpts := []grpc.ServerOption{}
	options := newOptions(opts...)
	grpcServer := grpc.NewServer(grpcOpts...)

	handle, teardown, err := svc.NewHandler(
		svc.Config(options.Config),
		svc.Logger(options.Logger),
		svc.AppURLs(options.App.AppURLs),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing collaboration service")
		return grpcServer, teardown, err
	}

	// register the app provider interface / OpenInApp call
	appproviderv1beta1.RegisterProviderAPIServer(grpcServer, handle)

	return grpcServer, teardown, nil
}
