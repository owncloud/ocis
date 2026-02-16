package http

import (
	"fmt"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	svc "github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0"
	"github.com/owncloud/reva/v2/pkg/store"
	"go-micro.dev/v4"
	microstore "go-micro.dev/v4/store"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service, err := http.NewService(
		http.TLSConfig(options.Config.HTTP.TLS),
		http.Logger(options.Logger),
		http.Name(options.Config.Service.Name),
		http.Version(version.GetString()),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
		http.TraceProvider(options.TraceProvider),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return http.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	signingKeyStore := store.Create(
		store.Store(options.Config.SigningKeys.Store),
		store.TTL(options.Config.SigningKeys.TTL),
		microstore.Nodes(options.Config.SigningKeys.Nodes...),
		microstore.Database("proxy"),
		microstore.Table("signing-keys"),
		store.Authentication(options.Config.SigningKeys.AuthUsername, options.Config.SigningKeys.AuthPassword),
	)

	handle := svc.NewService(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.Middleware(
			middleware.GetOtelhttpMiddleware(options.Config.Service.Name, options.TraceProvider),
			chimiddleware.RealIP,
			chimiddleware.RequestID,
			middleware.NoCache,
			middleware.Cors(
				cors.Logger(options.Logger),
				cors.AllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
				cors.AllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
				cors.AllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
				cors.AllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
			),
			middleware.Version(
				options.Config.Service.Name,
				version.GetString(),
			),
			middleware.Logger(options.Logger),
		),
		svc.Store(signingKeyStore),
	)

	if err := micro.RegisterHandler(service.Server(), handle); err != nil {
		return http.Service{}, err
	}

	return service, nil
}
