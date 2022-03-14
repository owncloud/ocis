package http

import (
	"github.com/go-chi/chi/v5/middleware"
	ocismiddleware "github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/service/http"
	"github.com/owncloud/ocis/ocis-pkg/version"
	svc "github.com/owncloud/ocis/thumbnails/pkg/service/http/v0"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/storage"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Name(options.Config.Service.Name),
		http.Version(version.String),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
	)

	handle := svc.NewService(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.Middleware(
			middleware.RealIP,
			middleware.RequestID,
			// ocismiddleware.Secure,
			ocismiddleware.Version(
				options.Config.Service.Name,
				version.String,
			),
			ocismiddleware.Logger(options.Logger),
		),
		svc.ThumbnailStorage(
			storage.NewFileSystemStorage(
				options.Config.Thumbnail.FileSystemStorage,
				options.Logger,
			),
		),
	)

	{
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLogging(handle, options.Logger)
		handle = svc.NewTracing(handle)
	}

	if err := micro.RegisterHandler(service.Server(), handle); err != nil {
		return http.Service{}, err
	}

	return service, nil
}
