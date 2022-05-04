package grpc

import (
	svc "github.com/owncloud/ocis/v2/extensions/store/pkg/service/v0"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	storesvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/store/v0"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Name(options.Config.Service.Name),
		grpc.Version(version.String),
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
	if err = storesvc.RegisterStoreHandler(service.Server(), hdlr); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}

	return service
}
