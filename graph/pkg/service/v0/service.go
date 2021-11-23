package svc

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/owncloud/ocis/graph/pkg/identity"
	"github.com/owncloud/ocis/graph/pkg/identity/ldap"
	"github.com/owncloud/ocis/ocis-pkg/account"
	opkgm "github.com/owncloud/ocis/ocis-pkg/middleware"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	GetMe(http.ResponseWriter, *http.Request)
	GetUsers(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	var backend identity.Backend
	switch options.Config.Identity.Backend {
	case "cs3":
		backend = &identity.CS3{
			Config: &options.Config.Reva,
			Logger: &options.Logger,
		}
	case "ldap":
		var err error
		conn := ldap.NewLDAPWithReconnect(&options.Logger,
			options.Config.Identity.LDAP.URI,
			options.Config.Identity.LDAP.BindDN,
			options.Config.Identity.LDAP.BindPassword,
		)
		if backend, err = identity.NewLDAPBackend(conn, options.Config.Identity.LDAP, &options.Logger); err != nil {
			options.Logger.Error().Msgf("Error initializing LDAP Backend: '%s'", err)
			return nil
		}
	default:
		options.Logger.Error().Msgf("Unknown Identity Backend: '%s'", options.Config.Identity.Backend)
		return nil
	}

	svc := Graph{
		config:          options.Config,
		mux:             m,
		logger:          &options.Logger,
		identityBackend: backend,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Route("/v1.0", func(r chi.Router) {
			r.Route("/me", func(r chi.Router) {
				r.Get("/", svc.GetMe)
				r.Get("/drives", svc.GetDrives)
				r.Get("/drive/root/children", svc.GetRootDriveChildren)
			})
			r.Route("/users", func(r chi.Router) {
				r.Get("/", svc.GetUsers)
				r.Route("/{userID}", func(r chi.Router) {
					r.Get("/", svc.GetUser)
				})
			})
			r.Route("/groups", func(r chi.Router) {
				r.Get("/", svc.GetGroups)
				r.Route("/{groupID}", func(r chi.Router) {
					r.Get("/", svc.GetGroup)
				})
			})
			r.Group(func(r chi.Router) {
				r.Use(opkgm.ExtractAccountUUID(
					account.Logger(options.Logger),
					account.JWTSecret(options.Config.TokenManager.JWTSecret)),
				)
				r.Route("/drives", func(r chi.Router) {
					r.Get("/", svc.GetDrives)
					r.Post("/", svc.CreateDrive)
				})
				r.Route("/Drive({firstSegmentIdentifier})", func(r chi.Router) {
					r.Patch("/*", svc.UpdateDrive)
				})
			})
		})
	})

	return svc
}
