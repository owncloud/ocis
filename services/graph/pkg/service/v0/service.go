package svc

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	ocisldap "github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity/ldap"
	graphm "github.com/owncloud/ocis/v2/services/graph/pkg/middleware"
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
	ChangeOwnPassword(http.ResponseWriter, *http.Request)

	GetGroups(http.ResponseWriter, *http.Request)
	GetGroup(http.ResponseWriter, *http.Request)
	PostGroup(http.ResponseWriter, *http.Request)
	PatchGroup(http.ResponseWriter, *http.Request)
	DeleteGroup(http.ResponseWriter, *http.Request)
	GetGroupMembers(http.ResponseWriter, *http.Request)
	PostGroupMember(http.ResponseWriter, *http.Request)
	DeleteGroupMember(http.ResponseWriter, *http.Request)

	GetDrives(w http.ResponseWriter, r *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Graph{
		config:               options.Config,
		mux:                  m,
		logger:               &options.Logger,
		spacePropertiesCache: ttlcache.NewCache(),
		eventsPublisher:      options.EventsPublisher,
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
	if options.IdentityBackend == nil {
		switch options.Config.Identity.Backend {
		case "cs3":
			svc.identityBackend = &identity.CS3{
				Config: options.Config.Reva,
				Logger: &options.Logger,
			}
		case "ldap":
			var err error

			var tlsConf *tls.Config
			if options.Config.Identity.LDAP.Insecure {
				// When insecure is set to true then we don't need a certificate.
				options.Config.Identity.LDAP.CACert = ""
				tlsConf = &tls.Config{
					//nolint:gosec // We need the ability to run with "insecure" (dev/testing)
					InsecureSkipVerify: options.Config.Identity.LDAP.Insecure,
				}
			}

			if options.Config.Identity.LDAP.CACert != "" {
				if err := ocisldap.WaitForCA(options.Logger,
					options.Config.Identity.LDAP.Insecure,
					options.Config.Identity.LDAP.CACert); err != nil {
					options.Logger.Fatal().Err(err).Msg("The configured LDAP CA cert does not exist")
				}
				if tlsConf == nil {
					tlsConf = &tls.Config{}
				}
				certs := x509.NewCertPool()
				pemData, err := ioutil.ReadFile(options.Config.Identity.LDAP.CACert)
				if err != nil {
					options.Logger.Error().Err(err).Msgf("Error initializing LDAP Backend")
					return nil
				}
				if !certs.AppendCertsFromPEM(pemData) {
					options.Logger.Error().Msgf("Error initializing LDAP Backend. Adding CA cert failed")
					return nil
				}
				tlsConf.RootCAs = certs
			}

			conn := ldap.NewLDAPWithReconnect(&options.Logger,
				ldap.Config{
					URI:          options.Config.Identity.LDAP.URI,
					BindDN:       options.Config.Identity.LDAP.BindDN,
					BindPassword: options.Config.Identity.LDAP.BindPassword,
					TLSConfig:    tlsConf,
				},
			)
			if svc.identityBackend, err = identity.NewLDAPBackend(conn, options.Config.Identity.LDAP, &options.Logger); err != nil {
				options.Logger.Error().Msgf("Error initializing LDAP Backend: '%s'", err)
				return nil
			}
		default:
			options.Logger.Error().Msgf("Unknown Identity Backend: '%s'", options.Config.Identity.Backend)
			return nil
		}
	} else {
		svc.identityBackend = options.IdentityBackend
	}

	if options.RoleService == nil {
		svc.roleService = settingssvc.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient)
	} else {
		svc.roleService = options.RoleService
	}

	roleManager := options.RoleManager
	if roleManager == nil {
		m := roles.NewManager(
			roles.CacheSize(1024),
			roles.CacheTTL(time.Hour),
			roles.Logger(options.Logger),
			roles.RoleService(svc.roleService),
		)
		roleManager = &m
	}

	requireAdmin := graphm.RequireAdmin(roleManager, options.Logger)

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Route("/v1.0", func(r chi.Router) {
			r.Route("/me", func(r chi.Router) {
				r.Get("/", svc.GetMe)
				r.Get("/drives", svc.GetDrives)
				r.Get("/drive/root/children", svc.GetRootDriveChildren)
				r.Post("/changePassword", svc.ChangeOwnPassword)
			})
			r.Route("/users", func(r chi.Router) {
				r.With(requireAdmin).Get("/", svc.GetUsers)
				r.With(requireAdmin).Post("/", svc.PostUser)
				r.Route("/{userID}", func(r chi.Router) {
					r.Get("/", svc.GetUser)
					r.With(requireAdmin).Delete("/", svc.DeleteUser)
					r.With(requireAdmin).Patch("/", svc.PatchUser)
				})
			})
			r.Route("/groups", func(r chi.Router) {
				r.With(requireAdmin).Get("/", svc.GetGroups)
				r.With(requireAdmin).Post("/", svc.PostGroup)
				r.Route("/{groupID}", func(r chi.Router) {
					r.Get("/", svc.GetGroup)
					r.With(requireAdmin).Delete("/", svc.DeleteGroup)
					r.With(requireAdmin).Patch("/", svc.PatchGroup)
					r.Route("/members", func(r chi.Router) {
						r.With(requireAdmin).Get("/", svc.GetGroupMembers)
						r.With(requireAdmin).Post("/$ref", svc.PostGroupMember)
						r.With(requireAdmin).Delete("/{memberID}/$ref", svc.DeleteGroupMember)
					})
				})
			})
			r.Route("/drives", func(r chi.Router) {
				r.Get("/", svc.GetAllDrives)
				r.Post("/", svc.CreateDrive)
				r.Route("/{driveID}", func(r chi.Router) {
					r.Patch("/", svc.UpdateDrive)
					r.Get("/", svc.GetSingleDrive)
					r.Delete("/", svc.DeleteDrive)
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
