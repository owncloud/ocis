package svc

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
	"github.com/libregraph/lico/bootstrap"
	guestBackendSupport "github.com/libregraph/lico/bootstrap/backends/guest"
	kcBackendSupport "github.com/libregraph/lico/bootstrap/backends/kc"
	ldapBackendSupport "github.com/libregraph/lico/bootstrap/backends/ldap"
	libreGraphBackendSupport "github.com/libregraph/lico/bootstrap/backends/libregraph"
	licoconfig "github.com/libregraph/lico/config"
	"github.com/libregraph/lico/server"
	"github.com/owncloud/ocis/v2/extensions/idp/pkg/assets"
	cs3BackendSupport "github.com/owncloud/ocis/v2/extensions/idp/pkg/backends/cs3/bootstrap"
	"github.com/owncloud/ocis/v2/extensions/idp/pkg/config"
	"github.com/owncloud/ocis/v2/extensions/idp/pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"gopkg.in/yaml.v2"
	"stash.kopano.io/kgol/rndm"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	ctx := context.Background()
	options := newOptions(opts...)
	logger := options.Logger.Logger
	assetVFS := assets.New(
		assets.Logger(options.Logger),
		assets.Config(options.Config),
	)

	if err := ldap.WaitForCA(options.Logger, options.Config.IDP.Insecure, options.Config.Ldap.TLSCACert); err != nil {
		logger.Fatal().Err(err).Msg("The configured LDAP CA cert does not exist")
	}
	if options.Config.IDP.Insecure {
		// force CACert to be empty to avoid lico try to load it
		options.Config.Ldap.TLSCACert = ""
	}

	if err := createTemporaryClientsConfig(
		options.Config.IDP.IdentifierRegistrationConf,
		options.Config.IDP.Iss,
		options.Config.Clients,
	); err != nil {
		logger.Fatal().Err(err).Msg("could not create default config")
	}

	switch options.Config.IDP.IdentityManager {
	case "cs3":
		cs3BackendSupport.MustRegister()
		if err := initCS3EnvVars(options.Config.Reva.Address, options.Config.MachineAuthAPIKey); err != nil {
			logger.Fatal().Err(err).Msg("could not initialize cs3 backend env vars")
		}
	case "ldap":
		ldapBackendSupport.MustRegister()
		if err := initLicoInternalLDAPEnvVars(&options.Config.Ldap); err != nil {
			logger.Fatal().Err(err).Msg("could not initialize ldap env vars")
		}
	default:
		guestBackendSupport.MustRegister()
		kcBackendSupport.MustRegister()
		libreGraphBackendSupport.MustRegister()
	}

	// https://play.golang.org/p/Mh8AVJCd593
	idpSettings := bootstrap.Settings(options.Config.IDP)
	bs, err := bootstrap.Boot(ctx, &idpSettings, &licoconfig.Config{
		Logger: log.LogrusWrap(logger),
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("could not bootstrap idp")
	}

	managers := bs.Managers()
	routes := []server.WithRoutes{managers.Must("identity").(server.WithRoutes)}
	handlers := managers.Must("handler").(http.Handler)

	svc := IDP{
		logger: options.Logger,
		config: options.Config,
		assets: assetVFS,
	}

	svc.initMux(ctx, routes, handlers, options)

	return svc
}

type temporaryClientConfig struct {
	Clients []config.Client `yaml:"clients"`
}

func createTemporaryClientsConfig(filePath, ocisURL string, clients []config.Client) error {

	folder := path.Dir(filePath)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0700); err != nil {
			return err
		}
	}

	for i, client := range clients {

		for i, entry := range client.RedirectURIs {
			client.RedirectURIs[i] = strings.ReplaceAll(entry, "{{OCIS_URL}}", strings.TrimRight(ocisURL, "/"))
		}
		for i, entry := range client.Origins {
			client.Origins[i] = strings.ReplaceAll(entry, "{{OCIS_URL}}", strings.TrimRight(ocisURL, "/"))
		}
		clients[i] = client
	}

	c := temporaryClientConfig{
		Clients: clients,
	}

	conf, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	confOnDisk, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer confOnDisk.Close()

	err = ioutil.WriteFile(filePath, conf, 0600)
	if err != nil {
		return err
	}

	return nil

}

// Init cs3 backend vars which are currently not accessible via idp api
func initCS3EnvVars(cs3Addr, machineAuthAPIKey string) error {
	var defaults = map[string]string{
		"CS3_GATEWAY":              cs3Addr,
		"CS3_MACHINE_AUTH_API_KEY": machineAuthAPIKey,
	}

	for k, v := range defaults {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("could not set cs3 env var %s=%s", k, v)
		}
	}

	return nil
}

// Init ldap backend vars which are currently not accessible via idp api
func initLicoInternalLDAPEnvVars(ldap *config.Ldap) error {
	filter := fmt.Sprintf("(objectclass=%s)", ldap.ObjectClass)
	if ldap.Filter != "" {
		filter = fmt.Sprintf("(&%s%s)", ldap.Filter, filter)
	}
	var defaults = map[string]string{
		"LDAP_URI":                 ldap.URI,
		"LDAP_BINDDN":              ldap.BindDN,
		"LDAP_BINDPW":              ldap.BindPassword,
		"LDAP_BASEDN":              ldap.BaseDN,
		"LDAP_SCOPE":               ldap.Scope,
		"LDAP_LOGIN_ATTRIBUTE":     ldap.LoginAttribute,
		"LDAP_EMAIL_ATTRIBUTE":     ldap.EmailAttribute,
		"LDAP_NAME_ATTRIBUTE":      ldap.NameAttribute,
		"LDAP_UUID_ATTRIBUTE":      ldap.UUIDAttribute,
		"LDAP_UUID_ATTRIBUTE_TYPE": ldap.UUIDAttributeType,
		"LDAP_FILTER":              filter,
	}

	if ldap.TLSCACert != "" {
		defaults["LDAP_TLS_CACERT"] = ldap.TLSCACert
	}

	for k, v := range defaults {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("could not set ldap env var %s=%s", k, v)
		}
	}

	return nil
}

// IDP defines implements the business logic for Service.
type IDP struct {
	logger log.Logger
	config *config.Config
	mux    *chi.Mux
	assets http.FileSystem
}

// initMux initializes the internal idp gorilla mux and mounts it in to a ocis chi-router
func (idp *IDP) initMux(ctx context.Context, r []server.WithRoutes, h http.Handler, options Options) {
	gm := mux.NewRouter()
	for _, route := range r {
		route.AddRoutes(ctx, gm)
	}

	// Delegate rest to provider which is also a handler.
	if h != nil {
		gm.NotFoundHandler = h
	}

	idp.mux = chi.NewMux()
	idp.mux.Use(options.Middleware...)

	idp.mux.Use(middleware.Static(
		"/signin/v1/",
		assets.New(
			assets.Logger(options.Logger),
			assets.Config(options.Config),
		),
	))

	// handle / | index.html with a template that needs to have the BASE_PREFIX replaced
	idp.mux.Get("/signin/v1/identifier", idp.Index())
	idp.mux.Get("/signin/v1/identifier/", idp.Index())
	idp.mux.Get("/signin/v1/identifier/index.html", idp.Index())

	idp.mux.Mount("/", gm)
}

// ServeHTTP implements the Service interface.
func (idp IDP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idp.mux.ServeHTTP(w, r)
}

// Index renders the static html with the
func (idp IDP) Index() http.HandlerFunc {

	f, err := idp.assets.Open("/identifier/index.html")
	if err != nil {
		idp.logger.Fatal().Err(err).Msg("Could not open index template")
	}

	template, err := ioutil.ReadAll(f)
	if err != nil {
		idp.logger.Fatal().Err(err).Msg("Could not read index template")
	}
	if err = f.Close(); err != nil {
		idp.logger.Fatal().Err(err).Msg("Could not close body")
	}

	// TODO add environment variable to make the path prefix configurable
	pp := "/signin/v1"
	indexHTML := bytes.Replace(template, []byte("__PATH_PREFIX__"), []byte(pp), 1)

	nonce := rndm.GenerateRandomString(32)
	indexHTML = bytes.Replace(indexHTML, []byte("__CSP_NONCE__"), []byte(nonce), 1)

	indexHTML = bytes.Replace(indexHTML, []byte("__PASSWORD_RESET_LINK__"), []byte(idp.config.Service.PasswordResetURI), 1)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(indexHTML); err != nil {
			idp.logger.Error().Err(err).Msg("could not write to response writer")
		}
	})
}
