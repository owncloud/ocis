package svc

import (
	"net/http"
	"time"

	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	opkgm "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/config"
	ocsm "github.com/owncloud/ocis/v2/services/ocs/pkg/middleware"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/response"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
)

// Service defines the service handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	GetConfig(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	roleService := options.RoleService
	if roleService == nil {
		roleService = settingssvc.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient)
	}
	roleManager := options.RoleManager
	if roleManager == nil {
		m := roles.NewManager(
			roles.CacheSize(1024),
			roles.CacheTTL(time.Hour*24*7),
			roles.Logger(options.Logger),
			roles.RoleService(roleService),
		)
		roleManager = &m
	}

	svc := Ocs{
		config:      options.Config,
		mux:         m,
		RoleManager: roleManager,
		logger:      options.Logger,
	}

	if svc.config.AccountBackend == "" {
		svc.config.AccountBackend = "cs3"
	}

	requireUser := ocsm.RequireUser(
		ocsm.Logger(options.Logger),
	)

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
			r.Route("/apps/files_sharing/api/v1", func(r chi.Router) {})
			r.Route("/apps/notifications/api/v1", func(r chi.Router) {})
			r.Route("/cloud", func(r chi.Router) {
				r.Route("/capabilities", func(r chi.Router) {})
				// TODO /apps
				r.Route("/user", func(r chi.Router) {
					r.Get("/signing-key", svc.GetSigningKey)
				})
			})
			r.Route("/config", func(r chi.Router) {
				r.With(requireUser).Get("/", svc.GetConfig)
			})
		})
	})

	return svc
}

// Ocs defines implements the business logic for Service.
type Ocs struct {
	config      *config.Config
	logger      log.Logger
	RoleService settingssvc.RoleService
	RoleManager *roles.Manager
	mux         *chi.Mux
}

// ServeHTTP implements the Service interface.
func (o Ocs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o.mux.ServeHTTP(w, r)
}

// NotFound uses ErrRender to always return a proper OCS payload
func (o Ocs) NotFound(w http.ResponseWriter, r *http.Request) {
	o.mustRender(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "not found"))
}

func (o Ocs) getCS3Backend() backend.UserBackend {
	revaClient, err := pool.GetGatewayServiceClient(o.config.Reva.Address)
	if err != nil {
		o.logger.Fatal().Msgf("could not get reva client at address %s", o.config.Reva.Address)
	}
	return backend.NewCS3UserBackend(nil, revaClient, o.config.MachineAuthAPIKey, "", nil, o.logger)
}

// NotImplementedStub returns a not implemented error
func (o Ocs) NotImplementedStub(w http.ResponseWriter, r *http.Request) {
	o.mustRender(w, r, response.ErrRender(data.MetaUnknownError.StatusCode, "Not implemented"))
}

func (o Ocs) mustRender(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
	if err := render.Render(w, r, renderer); err != nil {
		o.logger.Err(err).Msgf("failed to write response for ocs request %s on %s", r.Method, r.URL)
	}
}
