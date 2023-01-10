package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	ohttp "github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (ohttp.Service, error) {
	options := newOptions(opts...)
	service := options.Service

	svc, err := ohttp.NewService(
		ohttp.TLSConfig(options.Config.HTTP.TLS),
		ohttp.Logger(options.Logger),
		ohttp.Namespace(options.Config.HTTP.Namespace),
		ohttp.Name(options.Config.Service.Name),
		ohttp.Version(version.GetString()),
		ohttp.Address(options.Config.HTTP.Addr),
		ohttp.Context(options.Context),
		ohttp.Flags(options.Flags...),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return ohttp.Service{}, err
	}

	mux := chi.NewMux()

	mux.Use(chimiddleware.RealIP)
	mux.Use(chimiddleware.RequestID)
	mux.Use(middleware.TraceContext)
	mux.Use(middleware.NoCache)
	mux.Use(
		middleware.Cors(
			cors.Logger(options.Logger),
			cors.AllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
			cors.AllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
			cors.AllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
			cors.AllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
		))
	mux.Use(middleware.Secure)

	mux.Use(middleware.Version(
		options.Name,
		version.String,
	))

	// this logs http request related data
	mux.Use(middleware.Logger(
		options.Logger,
	))

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {

		r.Get("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// from https://www.rfc-editor.org/rfc/rfc7033#section-4.2
			//
			// If the "resource" parameter is a value for which the server has no
			// information, the server MUST indicate that it was unable to match the
			// request as per Section 10.4.5 of RFC 2616.
			// TODO the MUST might be a problem, is a guest instance ok enough?
			resource := r.URL.Query().Get("resource")
			rel := r.URL.Query().Get("rel")

			jrd, err := service.Webfinger(ctx, resource, rel)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			w.Header().Set("Content-type", "application/jrd+json")
			render.Status(r, http.StatusOK)
			render.JSON(w, r, jrd)
		})
	})

	err = micro.RegisterHandler(svc.Server(), mux)
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("failed to register the handler")
	}

	svc.Init()
	return svc, nil
}
