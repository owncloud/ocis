package http

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	svc "github.com/owncloud/ocis/v2/services/userlog/pkg/service"
	"go-micro.dev/v4"
)

// Service is the service interface
type Service interface {
}

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service, err := http.NewService(
		http.TLSConfig(options.Config.HTTP.TLS),
		http.Logger(options.Logger),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Name("userlog"),
		http.Version(version.GetString()),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return http.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	//middlewares := []func(stdhttp.Handler) stdhttp.Handler{
	//middleware.TraceContext,
	//chimiddleware.RequestID,
	//middleware.Version(
	//"userlog",
	//version.GetString(),
	//),
	//middleware.Logger(
	//options.Logger,
	//),
	//}

	handle, err := svc.NewUserlogService(options.Config, options.Consumer, options.Store, options.RegisteredEvents)
	if err != nil {
		return http.Service{}, err
	}

	{
		//handle = svc.NewInstrument(handle, options.Metrics)
		//handle = svc.NewLogging(handle, options.Logger)
		//handle = svc.NewTracing(handle)
	}

	if err := micro.RegisterHandler(service.Server(), handle); err != nil {
		return http.Service{}, err
	}

	return service, nil
}
