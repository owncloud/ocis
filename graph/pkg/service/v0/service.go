package svc

import (
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/owncloud/ocis/graph/pkg/identity"
	"github.com/owncloud/ocis/graph/pkg/identity/ldap"
	"github.com/owncloud/ocis/ocis-pkg/account"
	opkgm "github.com/owncloud/ocis/ocis-pkg/middleware"
)

const (
	// HeaderPurge defines the header name for the purge header.
	HeaderPurge = "Purge"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	GetMe(http.ResponseWriter, *http.Request)
	GetUsers(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
	PostUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	PatchUser(http.ResponseWriter, *http.Request)

	GetGroups(http.ResponseWriter, *http.Request)
	GetGroup(http.ResponseWriter, *http.Request)
	PostGroup(http.ResponseWriter, *http.Request)
	GetGroupMembers(http.ResponseWriter, *http.Request)
	PostGroupMember(http.ResponseWriter, *http.Request)

	GetDrives(w http.ResponseWriter, r *http.Request)
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
		config:               options.Config,
		mux:                  m,
		logger:               &options.Logger,
		identityBackend:      backend,
		spacePropertiesCache: ttlcache.NewCache(),
	}
	if options.GatewayClient == nil {
		var err error
		svc.gatewayClient, err = pool.GetGatewayServiceClient(options.Config.Reva.Address)
		if err != nil {
			options.Logger.Error().Err(err).Msg("Could not get gateway client")
			return nil
		}
	} else {
		svc.gatewayClient = options.GatewayClient
	}
	if options.HTTPClient == nil {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: options.Config.Spaces.Insecure, //nolint:gosec
		}
		svc.httpClient = &http.Client{}
	} else {
		svc.httpClient = options.HTTPClient
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
				r.Post("/", svc.PostUser)
				r.Route("/{userID}", func(r chi.Router) {
					r.Get("/", svc.GetUser)
					r.Delete("/", svc.DeleteUser)
					r.Patch("/", svc.PatchUser)
				})
			})
			r.Route("/groups", func(r chi.Router) {
				r.Get("/", svc.GetGroups)
				r.Post("/", svc.PostGroup)
				r.Route("/{groupID}", func(r chi.Router) {
					r.Get("/", svc.GetGroup)
					r.Route("/members", func(r chi.Router) {
						r.Get("/", svc.GetGroupMembers)
						r.Post("/$ref", svc.PostGroupMember)
					})
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
					r.Route("/{driveID}", func(r chi.Router) {
						r.Patch("/", svc.UpdateDrive)
						r.Get("/", svc.GetSingleDrive)
						r.Delete("/", svc.DeleteDrive)
					})
				})
			})
		})
	})

	return svc
}

// parseHeaderPurge parses the 'Purge' header.
// '1', 't', 'T', 'TRUE', 'true', 'True' are parsed as true
// all other values are false.
func parsePurgeHeader(h http.Header) bool {
	val := h.Get(HeaderPurge)

	if b, err := strconv.ParseBool(val); err == nil {
		return b
	}
	return false
}
