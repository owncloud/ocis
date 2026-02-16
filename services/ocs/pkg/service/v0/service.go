package svc

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	opkgm "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/config"
	ocsm "github.com/owncloud/ocis/v2/services/ocs/pkg/middleware"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/response"
	microstore "go-micro.dev/v4/store"
)

// Service defines the service handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Ocs{
		config: options.Config,
		mux:    m,
		logger: options.Logger,
		store:  options.Store,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.NotFound(svc.NotFound)
		r.Use(middleware.StripSlashes)
		r.Use(opkgm.ExtractAccountUUID(
			account.Logger(options.Logger),
			account.JWTSecret(options.Config.TokenManager.JWTSecret)),
		)
		r.Use(ocsm.OCSFormatCtx) // updates request Accept header according to format=(json|xml) query parameter
		r.Route("/v{version:(1|2)}.php", func(r chi.Router) {
			r.Use(response.VersionCtx) // stores version in context
			r.Get("/cloud/user/signing-key", svc.GetSigningKey)
		})
	})

	_ = chi.Walk(m, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		options.Logger.Debug().Str("method", method).Str("route", route).Int("middlewares", len(middlewares)).Msg("serving endpoint")
		return nil
	})

	return svc
}

// Ocs defines implements the business logic for Service.
type Ocs struct {
	config *config.Config
	logger log.Logger
	mux    *chi.Mux
	store  microstore.Store
}

// ServeHTTP implements the Service interface.
func (o Ocs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o.mux.ServeHTTP(w, r)
}

// NotFound uses ErrRender to always return a proper OCS payload
func (o Ocs) NotFound(w http.ResponseWriter, r *http.Request) {
	o.mustRender(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "not found"))
}

func (o Ocs) mustRender(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
	if err := render.Render(w, r, renderer); err != nil {
		o.logger.Err(err).Msgf("failed to write response for ocs request %s on %s", r.Method, r.URL)
	}
}
