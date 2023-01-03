package svc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jellydator/ttlcache/v3"
	libregraph "github.com/owncloud/libre-graph-api-go"
	ocisldap "github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/store"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity/ldap"
	graphm "github.com/owncloud/ocis/v2/services/graph/pkg/middleware"
)

const (
	// HeaderPurge defines the header name for the purge header.
	HeaderPurge     = "Purge"
	displayNameAttr = "displayName"
)

// Service defines the service handlers.
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

	GetEducationSchools(http.ResponseWriter, *http.Request)
	GetEducationSchool(http.ResponseWriter, *http.Request)
	PostEducationSchool(http.ResponseWriter, *http.Request)
	PatchEducationSchool(http.ResponseWriter, *http.Request)
	DeleteEducationSchool(http.ResponseWriter, *http.Request)
	GetEducationSchoolUsers(http.ResponseWriter, *http.Request)
	PostEducationSchoolUser(http.ResponseWriter, *http.Request)
	DeleteEducationSchoolUser(http.ResponseWriter, *http.Request)

	GetEducationUsers(http.ResponseWriter, *http.Request)
	GetEducationUser(http.ResponseWriter, *http.Request)
	PostEducationUser(http.ResponseWriter, *http.Request)
	DeleteEducationUser(http.ResponseWriter, *http.Request)
	PatchEducationUser(http.ResponseWriter, *http.Request)

	GetDrives(w http.ResponseWriter, r *http.Request)
	GetSingleDrive(w http.ResponseWriter, r *http.Request)
	GetAllDrives(w http.ResponseWriter, r *http.Request)
	GetRootDriveChildren(w http.ResponseWriter, r *http.Request)
	CreateDrive(w http.ResponseWriter, r *http.Request)
	UpdateDrive(w http.ResponseWriter, r *http.Request)
	DeleteDrive(w http.ResponseWriter, r *http.Request)

	GetTags(w http.ResponseWriter, r *http.Request)
	AssignTags(w http.ResponseWriter, r *http.Request)
	UnassignTags(w http.ResponseWriter, r *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) (Graph, error) {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	spacePropertiesCache := ttlcache.New[string, interface{}](
		ttlcache.WithTTL[string, interface{}](
			time.Duration(options.Config.Spaces.ExtendedSpacePropertiesCacheTTL),
		),
	)
	go spacePropertiesCache.Start()

	usersCache := ttlcache.New[string, libregraph.User](
		ttlcache.WithTTL[string, libregraph.User](
			time.Duration(options.Config.Spaces.UsersCacheTTL),
		),
	)
	go usersCache.Start()

	groupsCache := ttlcache.New[string, libregraph.Group](
		ttlcache.WithTTL[string, libregraph.Group](
			time.Duration(options.Config.Spaces.GroupsCacheTTL),
		),
	)
	go groupsCache.Start()

	svc := Graph{
		config:                   options.Config,
		mux:                      m,
		logger:                   &options.Logger,
		spacePropertiesCache:     spacePropertiesCache,
		usersCache:               usersCache,
		groupsCache:              groupsCache,
		eventsPublisher:          options.EventsPublisher,
		gatewayClient:            options.GatewayClient,
		searchService:            options.SearchService,
		identityEducationBackend: options.IdentityEducationBackend,
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
					MinVersion: tls.VersionTLS12,
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
					tlsConf = &tls.Config{
						MinVersion: tls.VersionTLS12,
					}
				}
				certs := x509.NewCertPool()
				pemData, err := os.ReadFile(options.Config.Identity.LDAP.CACert)
				if err != nil {
					options.Logger.Error().Err(err).Msgf("Error initializing LDAP Backend")
					return svc, err
				}
				if !certs.AppendCertsFromPEM(pemData) {
					options.Logger.Error().Msgf("Error initializing LDAP Backend. Adding CA cert failed")
					return svc, err
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
			lb, err := identity.NewLDAPBackend(conn, options.Config.Identity.LDAP, &options.Logger)
			if err != nil {
				options.Logger.Error().Msgf("Error initializing LDAP Backend: '%s'", err)
				return svc, err
			}
			svc.identityBackend = lb
			if options.IdentityEducationBackend == nil {
				if options.Config.Identity.LDAP.EducationResourcesEnabled {
					svc.identityEducationBackend = lb
				} else {
					errEduBackend := &identity.ErrEducationBackend{}
					svc.identityEducationBackend = errEduBackend
				}
			}
		default:
			err := fmt.Errorf("Unknown Identity Backend: '%s'", options.Config.Identity.Backend)
			options.Logger.Err(err)
			return svc, err
		}
	} else {
		svc.identityBackend = options.IdentityBackend
	}

	if options.PermissionService == nil {
		svc.permissionsService = settingssvc.NewPermissionService("com.owncloud.api.settings", grpc.DefaultClient())
	} else {
		svc.permissionsService = options.PermissionService
	}

	roleManager := options.RoleManager
	if roleManager == nil {
		storeOptions := store.OcisStoreOptions{
			Type:    options.Config.CacheStore.Type,
			Address: options.Config.CacheStore.Address,
			Size:    options.Config.CacheStore.Size,
		}
		m := roles.NewManager(
			roles.StoreOptions(storeOptions),
			roles.Logger(options.Logger),
			roles.RoleService(options.RoleService),
		)
		roleManager = &m
	}

	var requireAdmin func(http.Handler) http.Handler
	if options.RequireAdminMiddleware == nil {
		requireAdmin = graphm.RequireAdmin(roleManager, options.Logger)
	} else {
		requireAdmin = options.RequireAdminMiddleware
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Route("/v1.0", func(r chi.Router) {
			r.Route("/extensions/org.libregraph", func(r chi.Router) {
				r.Get("/tags", svc.GetTags)
				r.Put("/tags", svc.AssignTags)
				r.Delete("/tags", svc.UnassignTags)
			})
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
					r.With(requireAdmin).Route("/appRoleAssignments", func(r chi.Router) {
						r.Get("/", svc.ListAppRoleAssignments)
						r.Post("/", svc.CreateAppRoleAssignment)
						r.Delete("/{appRoleAssignmentID}", svc.DeleteAppRoleAssignment)
					})
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
			r.With(requireAdmin).Route("/education", func(r chi.Router) {
				r.Route("/schools", func(r chi.Router) {
					r.Get("/", svc.GetEducationSchools)
					r.Post("/", svc.PostEducationSchool)
					r.Route("/{schoolID}", func(r chi.Router) {
						r.Get("/", svc.GetEducationSchool)
						r.Delete("/", svc.DeleteEducationSchool)
						r.Patch("/", svc.PatchEducationSchool)
						r.Route("/users", func(r chi.Router) {
							r.Get("/", svc.GetEducationSchoolUsers)
							r.Post("/$ref", svc.PostEducationSchoolUser)
							r.Delete("/{userID}/$ref", svc.DeleteEducationSchoolUser)
						})
					})
				})
				r.Route("/users", func(r chi.Router) {
					r.Get("/", svc.GetEducationUsers)
					r.Post("/", svc.PostEducationUser)
					r.Route("/{userID}", func(r chi.Router) {
						r.Get("/", svc.GetEducationUser)
						r.Delete("/", svc.DeleteEducationUser)
						r.Patch("/", svc.PatchEducationUser)
					})
				})
				r.Route("/classes", func(r chi.Router) {
					r.Get("/", svc.GetGroups)
					r.Post("/", svc.PostGroup)
					r.Route("/{groupID}", func(r chi.Router) {
						r.Get("/", svc.GetGroup)
						r.Delete("/", svc.DeleteGroup)
						r.Patch("/", svc.PatchGroup)
						r.Route("/members", func(r chi.Router) {
							r.Get("/", svc.GetGroupMembers)
							r.Post("/$ref", svc.PostGroupMember)
							r.Delete("/{memberID}/$ref", svc.DeleteGroupMember)
						})
					})
				})
			})
			r.Get("/applications/{applicationID}", svc.GetApplication)
		})
	})

	return svc, nil
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
