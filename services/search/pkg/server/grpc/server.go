package grpc

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	svc "github.com/owncloud/ocis/v2/services/search/pkg/service/v0"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service, err := grpc.NewService(
		grpc.TLSEnabled(options.Config.GRPC.TLS.Enabled),
		grpc.TLSCert(
			options.Config.GRPC.TLS.Cert,
			options.Config.GRPC.TLS.Key,
		),
		grpc.Name(options.Config.Service.Name),
		grpc.Context(options.Context),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Logger(options.Logger),
		grpc.Flags(options.Flags...),
		grpc.Version(version.GetString()),
	)
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("Error creating search service")
		return grpc.Service{}
	}

	handle, err := svc.NewHandler(
		svc.Config(options.Config),
		svc.Logger(options.Logger),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing search service")
		return grpc.Service{}
	}
	_ = searchsvc.RegisterSearchProviderHandler(
		service.Server(),
		handle,
	)
	return service
}
