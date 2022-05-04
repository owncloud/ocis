package grpc

import (
	accountssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/accounts/v0"

	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)
	handler := options.Handler

	service := grpc.NewService(
		grpc.Name(options.Config.Service.Name),
		grpc.Context(options.Context),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Logger(options.Logger),
		grpc.Flags(options.Flags...),
		grpc.Version(version.String),
	)

	if err := accountssvc.RegisterAccountsServiceHandler(service.Server(), handler); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}
	if err := accountssvc.RegisterGroupsServiceHandler(service.Server(), handler); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register groups handler")
	}
	if err := accountssvc.RegisterIndexServiceHandler(service.Server(), handler); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register index handler")
	}

	return service
}
