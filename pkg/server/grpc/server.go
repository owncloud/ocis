package grpc

import (
	"github.com/owncloud/ocis-pkg/v2/service/grpc"
	"github.com/owncloud/ocis-store/pkg/proto/v0"
	svc "github.com/owncloud/ocis-store/pkg/service/v0"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("store"),
		grpc.Context(options.Context),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Logger(options.Logger),
		grpc.Flags(options.Flags...),
	)

	hdlr, err := svc.New(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
	)
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("could not initialize service handler")
	}
	if err = proto.RegisterStoreHandler(service.Server(), hdlr); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}

	service.Init()
	return service
}
