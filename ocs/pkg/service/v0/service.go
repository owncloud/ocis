package svc

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/micro/go-micro/v2/client/grpc"

	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocs/pkg/config"
	ocsm "github.com/owncloud/ocis/ocs/pkg/middleware"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/response"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

var defaultClient = grpc.NewClient()

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	GetConfig(http.ResponseWriter, *http.Request)
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
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.NotFound(svc.NotFound)
		r.Use(middleware.StripSlashes)
		r.Use(ocsm.AccessToken(
			ocsm.Logger(options.Logger),
			ocsm.TokenManagerConfig(options.Config.TokenManager),
		))
		r.Use(ocsm.OCSFormatCtx) // updates request Accept header according to format=(json|xml) query parameter
		r.Route("/v{version:(1|2)}.php", func(r chi.Router) {
			r.Use(response.VersionCtx) // stores version in context
			r.Route("/apps/files_sharing/api/v1", func(r chi.Router) {})
			r.Route("/apps/notifications/api/v1", func(r chi.Router) {})
			r.Route("/cloud", func(r chi.Router) {
				r.Route("/capabilities", func(r chi.Router) {})
				r.Route("/user", func(r chi.Router) {
					r.Get("/", svc.GetUser)
					r.Get("/signing-key", svc.GetSigningKey)
				})
				r.Route("/users", func(r chi.Router) {
					r.Get("/", svc.ListUsers)
					r.Post("/", svc.AddUser)
					r.Get("/{userid}", svc.GetUser)
					r.Put("/{userid}", svc.EditUser)
					r.Delete("/{userid}", svc.DeleteUser)

					r.Route("/{userid}/groups", func(r chi.Router) {
						r.Get("/", svc.ListUserGroups)
						r.Post("/", svc.AddToGroup)
						r.Delete("/", svc.RemoveFromGroup)
					})
				})
				r.Route("/groups", func(r chi.Router) {
					r.Get("/", svc.ListGroups)
					r.Post("/", svc.AddGroup)
					r.Delete("/{groupid}", svc.DeleteGroup)
					r.Get("/{groupid}", svc.GetGroupMembers)
				})
			})
			r.Route("/config", func(r chi.Router) {
				r.Get("/", svc.GetConfig)
			})
		})
	})

	return svc
}

// Ocs defines implements the business logic for Service.
type Ocs struct {
	config *config.Config
	logger log.Logger
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (o Ocs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o.mux.ServeHTTP(w, r)
}

// NotFound uses ErrRender to always return a proper OCS payload
func (o Ocs) NotFound(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "not found"))
}

func (o Ocs) getAccountService() accounts.AccountsService {
	return accounts.NewAccountsService("com.owncloud.api.accounts", defaultClient)
}

func (o Ocs) getGroupsService() accounts.GroupsService {
	return accounts.NewGroupsService("com.owncloud.api.accounts", defaultClient)
}
