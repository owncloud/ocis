package http

import (
	"github.com/cs3org/reva/v2/pkg/events/server"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-micro/plugins/v4/events/natsjs"
	graphMiddleware "github.com/owncloud/ocis/v2/extensions/graph/pkg/middleware"
	svc "github.com/owncloud/ocis/v2/extensions/graph/pkg/service/v0"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/pkg/errors"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Name("graph"),
		http.Version(version.GetString()),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)

	publisher, err := server.NewNatsStream(
		natsjs.Address(options.Config.Events.Endpoint),
		natsjs.ClusterID(options.Config.Events.Cluster),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing events publisher")
		return http.Service{}, errors.Wrap(err, "could not initialize events publisher")
	}

	handle := svc.NewService(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.Middleware(
			chimiddleware.RequestID,
			middleware.Version(
				"graph",
				version.GetString(),
			),
			middleware.Logger(
				options.Logger,
			),
			graphMiddleware.Auth(
				account.Logger(options.Logger),
				account.JWTSecret(options.Config.TokenManager.JWTSecret),
			),
		),
		svc.EventsPublisher(publisher),
	)

	if handle == nil {
		return http.Service{}, errors.New("could not initialize graph service")
	}

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
