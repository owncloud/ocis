package grpc

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	svc "github.com/owncloud/ocis/v2/services/search/pkg/service/grpc/v0"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) (grpc.Service, error) {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Name(options.Config.Service.Name),
		grpc.Context(options.Context),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Logger(options.Logger),
		grpc.Flags(options.Flags...),
		grpc.Version(version.GetString()),
	)

	handle, closer, err := svc.NewHandler(
		svc.Config(options.Config),
		svc.Logger(options.Logger),
	)
	defer closer()
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing search service")
		return grpc.Service{}, err
	}

	if err := searchsvc.RegisterSearchProviderHandler(
		service.Server(),
		handle,
	); err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error registering search provider handler")
		return grpc.Service{}, err
	}

	return service, nil
}
