package http

import (
	"fmt"

	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/internal/app"
	"github.com/riandyrn/otelchi"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service, err := http.NewService(
		http.TLSConfig(options.Config.HTTP.TLS),
		http.Logger(options.Logger),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Name(options.Config.Service.Name),
		http.Version(version.GetString()),
		http.Address(options.Config.HTTP.BindAddr),
		http.Context(options.Context),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return http.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	middlewares := []func(stdhttp.Handler) stdhttp.Handler{
		chimiddleware.RequestID,
		middleware.Version(
			"userlog",
			version.GetString(),
		),
		middleware.Logger(
			options.Logger,
		),
		middleware.ExtractAccountUUID(
			account.Logger(options.Logger),
			account.JWTSecret(options.Config.Secret), // previously, secret came from Config.TokenManager.JWTSecret
		),
		/*
			// Need CORS? not in the original server
			// Also, CORS isn't part of the config right now
			middleware.Cors(
				cors.Logger(options.Logger),
				cors.AllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
				cors.AllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
				cors.AllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
				cors.AllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
			),
		*/
		middleware.Secure,
	}

	mux := chi.NewMux()
	mux.Use(middlewares...)

	mux.Use(
		otelchi.Middleware(
			"collaboration",
			otelchi.WithChiRoutes(mux),
			otelchi.WithTracerProvider(options.TracerProvider),
			otelchi.WithPropagators(tracing.GetPropagator()),
		),
	)

	prepareRoutes(mux, options.App)

	if err := micro.RegisterHandler(service.Server(), mux); err != nil {
		return http.Service{}, err
	}

	return service, nil
}

func prepareRoutes(r *chi.Mux, demoapp *app.DemoApp) {
	r.Route("/wopi", func(r chi.Router) {

		r.Get("/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			app.WopiInfoHandler(demoapp, w, r)
		})

		r.Route("/files/{fileid}", func(r chi.Router) {

			r.Use(func(h stdhttp.Handler) stdhttp.Handler {
				// authentication and wopi context
				return app.WopiContextAuthMiddleware(demoapp, h)
			},
			)

			r.Get("/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				app.CheckFileInfo(demoapp, w, r)
			})

			r.Post("/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				action := r.Header.Get("X-WOPI-Override")
				switch action {

				case "LOCK":
					app.Lock(demoapp, w, r)
				case "GET_LOCK":
					app.GetLock(demoapp, w, r)
				case "REFRESH_LOCK":
					app.RefreshLock(demoapp, w, r)
				case "UNLOCK":
					app.UnLock(demoapp, w, r)

				case "PUT_USER_INFO":
					// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/putuserinfo
					stdhttp.Error(w, stdhttp.StatusText(stdhttp.StatusNotImplemented), stdhttp.StatusNotImplemented)
				case "PUT_RELATIVE":
					// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/putrelativefile
					stdhttp.Error(w, stdhttp.StatusText(stdhttp.StatusNotImplemented), stdhttp.StatusNotImplemented)
				case "RENAME_FILE":
					// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/renamefile
					stdhttp.Error(w, stdhttp.StatusText(stdhttp.StatusNotImplemented), stdhttp.StatusNotImplemented)
				case "DELETE":
					// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/deletefile
					stdhttp.Error(w, stdhttp.StatusText(stdhttp.StatusNotImplemented), stdhttp.StatusNotImplemented)

				default:
					stdhttp.Error(w, stdhttp.StatusText(stdhttp.StatusInternalServerError), stdhttp.StatusInternalServerError)
				}
			})

			r.Route("/contents", func(r chi.Router) {
				r.Get("/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
					app.GetFile(demoapp, w, r)
				})

				r.Post("/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
					action := r.Header.Get("X-WOPI-Override")
					switch action {

					case "PUT":
						app.PutFile(demoapp, w, r)

					default:
						stdhttp.Error(w, stdhttp.StatusText(stdhttp.StatusInternalServerError), stdhttp.StatusInternalServerError)
					}
				})
			})
		})
	})
}
