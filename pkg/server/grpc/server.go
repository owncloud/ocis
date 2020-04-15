package grpc

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	svc "github.com/owncloud/ocis-settings/pkg/service/v0"
	"github.com/owncloud/ocis-pkg/v2/service/grpc"
)

// NewService initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Name(options.Name),
		grpc.Context(options.Context),
		grpc.Address(options.Address),
		grpc.Namespace(options.Namespace),
		grpc.Logger(options.Logger),
	)

	hdlr := svc.NewService(options.Config)
	if err := proto.RegisterBundleServiceHandler(service.Server(), hdlr); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}

	service.Init()
	return service
}
