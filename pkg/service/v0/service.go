package svc

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/gorilla/mux"
	"github.com/owncloud/ocis-konnectd/pkg/config"
	"github.com/owncloud/ocis-konnectd/pkg/log"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"stash.kopano.io/kc/konnect/bootstrap"
	kcconfig "stash.kopano.io/kc/konnect/config"
	"stash.kopano.io/kc/konnect/server"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Dummy(http.ResponseWriter, *http.Request)
}

// Konnectd defines implements the business logic for Service.
type Konnectd struct {
	config *config.Config
	mux    http.Handler
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	ctx := context.Background()
	options := newOptions(opts...)
	logger := options.Logger.Logger
	initKonnectInternalEnvVars(logger)

	bs, err := bootstrap.Boot(ctx, &options.Config.Konnectd, &kcconfig.Config{
		Logger: log.Wrap(logger),
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not bootstrap konnectd")
	}

	managers := bs.Managers()
	routes := []server.WithRoutes{managers.Must("identity").(server.WithRoutes)}
	handlers := managers.Must("handler").(http.Handler)

	return Konnectd{
		config: options.Config,
		mux:    newMux(ctx, routes, handlers, options.Middleware),
	}
}

// Init vars which are currently not accessible via konnectd api
func initKonnectInternalEnvVars(l zerolog.Logger) {
	var defaults = map[string]string{
		"LDAP_URI":                 "ldap://localhost:9125",
		"LDAP_BINDDN":              "cn=admin,dc=example,dc=org",
		"LDAP_BINDPW":              "admin",
		"LDAP_BASEDN":              "ou=users,dc=example,dc=org",
		"LDAP_SCOPE":               "sub",
		"LDAP_LOGIN_ATTRIBUTE":     "uid",
		"LDAP_EMAIL_ATTRIBUTE":     "mail",
		"LDAP_NAME_ATTRIBUTE":      "cn",
		"LDAP_UUID_ATTRIBUTE":      "customuid",
		"LDAP_UUID_ATTRIBUTE_TYPE": "text",
		"LDAP_FILTER":              "(objectClass=person)",
	}

	for k, v := range defaults {
		if _, exists := os.LookupEnv(k); !exists {
			if err := os.Setenv(k, v); err != nil {
				l.Fatal().Err(err).Msgf("Could not set env var %s=%s", k, v)
			}
		}
	}
}

// newMux initializes the internal konnectd gorilla mux and mounts it in to a ocis chi-router
func newMux(ctx context.Context, r []server.WithRoutes, h http.Handler, middleware []func(http.Handler) http.Handler) http.Handler {
	gm := mux.NewRouter()
	for _, route := range r {
		route.AddRoutes(ctx, gm)
	}

	// Delegate rest to provider which is also a handler.
	if h != nil {
		gm.NotFoundHandler = h
	}

	m := chi.NewMux()
	m.Use(middleware...)
	m.Mount("/", gm)

	return m
}

// ServeHTTP implements the Service interface.
func (g Konnectd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Dummy implements the Service interface.
func (g Konnectd) Dummy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))
}
