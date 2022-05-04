package grpc

import (
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	svc "github.com/owncloud/ocis/v2/extensions/thumbnails/pkg/service/grpc/v0"
	"github.com/owncloud/ocis/v2/extensions/thumbnails/pkg/service/grpc/v0/decorators"
	"github.com/owncloud/ocis/v2/extensions/thumbnails/pkg/thumbnail/imgsource"
	"github.com/owncloud/ocis/v2/extensions/thumbnails/pkg/thumbnail/storage"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	thumbnailssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/thumbnails/v0"
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
		grpc.Version(version.String),
	)
	tconf := options.Config.Thumbnail
	gc, err := pool.GetGatewayServiceClient(tconf.RevaGateway)
	if err != nil {
		options.Logger.Error().Err(err).Msg("could not get gateway client")
		return grpc.Service{}
	}
	var thumbnail decorators.DecoratedService
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
		thumbnail = decorators.NewInstrument(thumbnail, options.Metrics)
		thumbnail = decorators.NewLogging(thumbnail, options.Logger)
		thumbnail = decorators.NewTracing(thumbnail)
	}

	_ = thumbnailssvc.RegisterThumbnailServiceHandler(
		service.Server(),
		thumbnail,
	)

	return service
}
