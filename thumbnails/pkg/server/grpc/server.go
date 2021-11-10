package grpc

import (
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	svc "github.com/owncloud/ocis/thumbnails/pkg/service/v0"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/imgsource"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/storage"
)

// NewService initializes the grpc service and server.
func NewService(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Logger(options.Logger),
		grpc.Namespace(options.Namespace),
		grpc.Name(options.Name),
		grpc.Version(version.String),
		grpc.Address(options.Address),
		grpc.Context(options.Context),
		grpc.Flags(options.Flags...),
		grpc.Version(options.Config.Server.Version),
	)
	tconf := options.Config.Thumbnail
	gc, err := pool.GetGatewayServiceClient(tconf.RevaGateway) //TODO: insecure defaults to true, https://github.com/cs3org/reva/issues/2216
	if err != nil {
		options.Logger.Error().Err(err).Msg("could not get gateway client")
		return grpc.Service{}
	}
	var thumbnail proto.ThumbnailServiceHandler
	{
		thumbnail = svc.NewService(
			svc.Config(options.Config),
			svc.Logger(options.Logger),
			svc.ThumbnailSource(imgsource.NewWebDavSource(tconf)),
			svc.ThumbnailStorage(
				storage.NewFileSystemStorage(
					tconf.FileSystemStorage,
					options.Logger,
				),
			),
			svc.CS3Source(imgsource.NewCS3Source(tconf, gc)),
			svc.CS3Client(gc),
		)
		thumbnail = svc.NewInstrument(thumbnail, options.Metrics)
		thumbnail = svc.NewLogging(thumbnail, options.Logger)
		thumbnail = svc.NewTracing(thumbnail)
	}

	_ = proto.RegisterThumbnailServiceHandler(
		service.Server(),
		thumbnail,
	)

	return service
}
