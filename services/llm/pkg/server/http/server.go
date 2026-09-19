// services/llm/pkg/server/http/server.go
package http

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go-micro.dev/v4"

	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	ohttp "github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/llm/pkg/service"
)

// Server initializes the http service and server.
func Server(opts ...Option) (ohttp.Service, error) {
	options := newOptions(opts...)

	newService, err := ohttp.NewService(
		ohttp.TLSConfig(options.Config.HTTP.TLS),
		ohttp.Logger(options.Logger),
		ohttp.Namespace(options.Config.HTTP.Namespace),
		ohttp.Name(options.Config.Service.Name),
		ohttp.Version(version.GetString()),
		ohttp.Address(options.Config.HTTP.Addr),
		ohttp.Context(options.Context),
		ohttp.Flags(options.Flags...),
		ohttp.TraceProvider(options.TraceProvider),
	)
	if err != nil {
		options.Logger.Error().Err(err).Msg("Error initializing http service")
		return ohttp.Service{}, err
	}

	svc := service.New(options.Config, options.Limiter, options.Logger)

	mux := chi.NewMux()
	mux.Use(
		middleware.GetOtelhttpMiddleware(options.Config.Service.Name, options.TraceProvider),
		chimiddleware.RequestID,
		middleware.Version(options.Config.Service.Name, version.GetString()),
		middleware.Logger(options.Logger),
		middleware.ExtractAccountUUID(
			account.Logger(options.Logger),
			account.JWTSecret(options.Config.TokenManager.JWTSecret),
		),
		middleware.Cors(
			cors.Logger(options.Logger),
			cors.AllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
			cors.AllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
			cors.AllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
			cors.AllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
		),
	)

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {
		// proxy forwards the original request path unchanged (it does not
		// strip the matched prefix — see activitylog's route registration
		// for the same pattern), so this must be the full external path.
		r.Post("/graph/v1beta1/extensions/org.libregraph/llm/chat/completions", svc.HandleChatCompletions)
	})

	if err := micro.RegisterHandler(newService.Server(), mux); err != nil {
		options.Logger.Error().Err(err).Msg("failed to register the handler")
		return ohttp.Service{}, err
	}

	newService.Init()
	return newService, nil
}
