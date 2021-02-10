package grpc

import (
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/ocis-pkg/service/http"
	"github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	svc "github.com/owncloud/ocis/thumbnails/pkg/service/v0"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/imgsource"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/storage"
	"github.com/owncloud/ocis/thumbnails/pkg/version"
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

	var thumbnail proto.ThumbnailServiceHandler
	{
		thumbnail = svc.NewService(
			svc.Config(options.Config),
			svc.Logger(options.Logger),
			svc.ThumbnailSource(imgsource.NewWebDavSource(options.Config.Thumbnail.WebDavSource)),
			svc.ThumbnailStorage(
				storage.NewFileSystemStorage(
					options.Config.Thumbnail.FileSystemStorage,
					options.Logger,
				),
			),
		)
		thumbnail = svc.NewInstrument(thumbnail, options.Metrics)
		thumbnail = svc.NewLogging(thumbnail, options.Logger)
	}

	_ = proto.RegisterThumbnailServiceHandler(
		service.Server(),
		thumbnail,
	)

	service.Init()
	http.M.Unlock()
	return service
}
