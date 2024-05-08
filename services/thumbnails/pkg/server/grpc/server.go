package grpc

import (
	"github.com/cs3org/reva/v2/pkg/bytesize"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	thumbnailssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/thumbnails/v0"
	svc "github.com/owncloud/ocis/v2/services/thumbnails/pkg/service/grpc/v0"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/service/grpc/v0/decorators"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/thumbnail/imgsource"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/thumbnail/storage"
)

// NewService initializes the grpc service and server.
func NewService(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service, err := grpc.NewServiceWithClient(
		options.Config.GrpcClient,
		grpc.TLSEnabled(options.Config.GRPC.TLS.Enabled),
		grpc.TLSCert(
			options.Config.GRPC.TLS.Cert,
			options.Config.GRPC.TLS.Key,
		),
		grpc.Logger(options.Logger),
		grpc.Namespace(options.Namespace),
		grpc.Name(options.Name),
		grpc.Version(version.GetString()),
		grpc.Address(options.Address),
		grpc.Context(options.Context),
		grpc.Version(version.GetString()),
		grpc.TraceProvider(options.TraceProvider),
	)
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("Error creating thumbnail service")
		return grpc.Service{}
	}

	tconf := options.Config.Thumbnail
	tm, err := pool.StringToTLSMode(options.Config.GRPCClientTLS.Mode)
	if err != nil {
		options.Logger.Error().Err(err).Msg("could not get gateway client tls mode")
		return grpc.Service{}
	}

	gatewaySelector, err := pool.GatewaySelector(tconf.RevaGateway,
		pool.WithTLSCACert(options.Config.GRPCClientTLS.CACert),
		pool.WithTLSMode(tm),
		pool.WithRegistry(registry.GetRegistry()),
		pool.WithTracerProvider(options.TraceProvider),
	)
	if err != nil {
		options.Logger.Error().Err(err).Msg("could not get gateway selector")
		return grpc.Service{}
	}
	b, err := bytesize.Parse(tconf.MaxInputImageFileSize)
	if err != nil {
		options.Logger.Error().Err(err).Msg("could not parse MaxInputImageFileSize")
		return grpc.Service{}
	}

	var thumbnail decorators.DecoratedService
	{
		thumbnail = svc.NewService(
			svc.Config(options.Config),
			svc.Logger(options.Logger),
			svc.ThumbnailSource(imgsource.NewWebDavSource(tconf, b)),
			svc.ThumbnailStorage(
				storage.NewFileSystemStorage(
					tconf.FileSystemStorage,
					options.Logger,
				),
			),
			svc.CS3Source(imgsource.NewCS3Source(tconf, gatewaySelector, b)),
			svc.GatewaySelector(gatewaySelector),
		)
		thumbnail = decorators.NewInstrument(thumbnail, options.Metrics)
		thumbnail = decorators.NewLogging(thumbnail, options.Logger)
		thumbnail = decorators.NewTracing(thumbnail, options.TraceProvider)
	}

	_ = thumbnailssvc.RegisterThumbnailServiceHandler(
		service.Server(),
		thumbnail,
	)

	return service
}
