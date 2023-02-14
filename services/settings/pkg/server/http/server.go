package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	ohttp "github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/permissions"
	svc "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (ohttp.Service, error) {
	options := newOptions(opts...)

	service, err := ohttp.NewService(
		ohttp.TLSConfig(options.Config.HTTP.TLS),
		ohttp.Logger(options.Logger),
		ohttp.Name(options.Name),
		ohttp.Version(version.GetString()),
		ohttp.Address(options.Config.HTTP.Addr),
		ohttp.Namespace(options.Config.HTTP.Namespace),
		ohttp.Context(options.Context),
		ohttp.Flags(options.Flags...),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return ohttp.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	handle := svc.NewService(options.Config, options.Logger)

	{
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLogging(handle, options.Logger)
		handle = svc.NewTracing(handle)
	}

	mux := chi.NewMux()

	mux.Use(chimiddleware.RealIP)
	mux.Use(chimiddleware.RequestID)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.Cors(
		cors.Logger(options.Logger),
		cors.AllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
		cors.AllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
		cors.AllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
		cors.AllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
	))
	mux.Use(middleware.Secure)
	mux.Use(middleware.ExtractAccountUUID(
		account.Logger(options.Logger),
		account.JWTSecret(options.Config.TokenManager.JWTSecret)),
	)

	mux.Use(middleware.Version(
		options.Name,
		version.GetString(),
	))

	mux.Use(middleware.Logger(
		options.Logger,
	))

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {
		settingssvc.RegisterBundleServiceWeb(r, handle)
		settingssvc.RegisterValueServiceWeb(r, handle)
		settingssvc.RegisterRoleServiceWeb(r, handle)
		settingssvc.RegisterPermissionServiceWeb(r, handle)
		r.MethodFunc("POST", "/api/v0/settings/permissions-list", func(w http.ResponseWriter, r *http.Request) {

			req := &permissions.ListPermissionsRequest{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				http.Error(w, err.Error(), http.StatusPreconditionFailed)
				return
			}

			us, ok := revactx.ContextGetUser(r.Context())
			if !ok {
				http.Error(w, "invalid user", http.StatusUnauthorized)
				return
			}

			if us.GetId().GetOpaqueId() != req.UserID {
				http.Error(w, fmt.Sprintf("user %s not found", req.UserID), http.StatusNotFound)
				return
			}

			resp, err := handle.ListPermissions(r.Context(), req)
			if err != nil {
				if errors.Is(err, settings.ErrNotFound) {
					http.Error(w, err.Error(), http.StatusNotFound)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			render.Status(r, http.StatusOK)
			render.JSON(w, r, resp)
		})
	})

	_ = chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		options.Logger.Debug().Str("method", method).Str("route", route).Int("middlewares", len(middlewares)).Msg("serving endpoint")
		return nil
	})

	micro.RegisterHandler(service.Server(), mux)

	return service, nil
}
