package svc

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
	"github.com/libregraph/lico/bootstrap"
	guestBackendSupport "github.com/libregraph/lico/bootstrap/backends/guest"
	ldapBackendSupport "github.com/libregraph/lico/bootstrap/backends/ldap"
	libreGraphBackendSupport "github.com/libregraph/lico/bootstrap/backends/libregraph"
	licoconfig "github.com/libregraph/lico/config"
	"github.com/libregraph/lico/server"
	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
	"github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/idp/pkg/assets"
	cs3BackendSupport "github.com/owncloud/ocis/v2/services/idp/pkg/backends/cs3/bootstrap"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
	"github.com/owncloud/ocis/v2/services/idp/pkg/middleware"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"stash.kopano.io/kgol/rndm"
)

//go:embed l10n/locale
var _translationFS embed.FS

// Service defines the service handlers.
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

	if err := createTemporaryClientsConfig(
		options.Config.IDP.IdentifierRegistrationConf,
		options.Config.Commons.OcisURL,
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

		if err := ldap.WaitForCA(options.Logger, options.Config.IDP.Insecure, options.Config.Ldap.TLSCACert); err != nil {
			logger.Fatal().Err(err).Msg("The configured LDAP CA cert does not exist")
		}
		if options.Config.IDP.Insecure {
			// force CACert to be empty to avoid lico try to load it
			options.Config.Ldap.TLSCACert = ""
		}

		ldapBackendSupport.MustRegister()
		if err := initLicoInternalLDAPEnvVars(&options.Config.Ldap); err != nil {
			logger.Fatal().Err(err).Msg("could not initialize ldap env vars")
		}
	default:
		guestBackendSupport.MustRegister()
		libreGraphBackendSupport.MustRegister()
	}

	idpSettings := bootstrap.Settings{
		Iss:                               options.Config.IDP.Iss,
		IdentityManager:                   options.Config.IDP.IdentityManager,
		URIBasePath:                       options.Config.IDP.URIBasePath,
		SignInURI:                         options.Config.IDP.SignInURI,
		SignedOutURI:                      options.Config.IDP.SignedOutURI,
		AuthorizationEndpointURI:          options.Config.IDP.AuthorizationEndpointURI,
		EndsessionEndpointURI:             options.Config.IDP.EndsessionEndpointURI,
		Insecure:                          options.Config.IDP.Insecure,
		TrustedProxy:                      options.Config.IDP.TrustedProxy,
		AllowScope:                        options.Config.IDP.AllowScope,
		AllowClientGuests:                 options.Config.IDP.AllowClientGuests,
		AllowDynamicClientRegistration:    options.Config.IDP.AllowDynamicClientRegistration,
		EncryptionSecretFile:              options.Config.IDP.EncryptionSecretFile,
		Listen:                            options.Config.IDP.Listen,
		IdentifierClientDisabled:          options.Config.IDP.IdentifierClientDisabled,
		IdentifierClientPath:              options.Config.IDP.IdentifierClientPath,
		IdentifierRegistrationConf:        options.Config.IDP.IdentifierRegistrationConf,
		IdentifierScopesConf:              options.Config.IDP.IdentifierScopesConf,
		IdentifierDefaultBannerLogo:       options.Config.IDP.IdentifierDefaultBannerLogo,
		IdentifierDefaultSignInPageText:   options.Config.IDP.IdentifierDefaultSignInPageText,
		IdentifierDefaultUsernameHintText: options.Config.IDP.IdentifierDefaultUsernameHintText,
		IdentifierUILocales:               options.Config.IDP.IdentifierUILocales,
		SigningKid:                        options.Config.IDP.SigningKid,
		SigningMethod:                     options.Config.IDP.SigningMethod,
		SigningPrivateKeyFiles:            options.Config.IDP.SigningPrivateKeyFiles,
		ValidationKeysPath:                options.Config.IDP.ValidationKeysPath,
		CookieBackendURI:                  options.Config.IDP.CookieBackendURI,
		CookieNames:                       options.Config.IDP.CookieNames,
		CookieSameSite:                    options.Config.IDP.CookieSameSite,
		AccessTokenDurationSeconds:        options.Config.IDP.AccessTokenDurationSeconds,
		IDTokenDurationSeconds:            options.Config.IDP.IDTokenDurationSeconds,
		RefreshTokenDurationSeconds:       options.Config.IDP.RefreshTokenDurationSeconds,
		DyamicClientSecretDurationSeconds: options.Config.IDP.DynamicClientSecretDurationSeconds,
	}
	bs, err := bootstrap.Boot(ctx, &idpSettings, &licoconfig.Config{
		Logger: log.LogrusWrap(logger),
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("could not bootstrap idp")
	}

	managers := bs.Managers()
	routes := []server.WithRoutes{managers.Must("identity").(server.WithRoutes)}
	handlers := managers.Must("handler").(http.Handler)

	var translationFS fs.FS
	translationFS, _ = fs.Sub(_translationFS, "l10n/locale")
	translator := l10n.NewTranslator("en", "idp", translationFS)

	svc := &IDP{
		logger:     options.Logger,
		config:     options.Config,
		assets:     assetVFS,
		tp:         options.TraceProvider,
		translator: translator,
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
		if err := os.MkdirAll(folder, 0o700); err != nil {
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

	err = os.WriteFile(filePath, conf, 0o600)
	if err != nil {
		return err
	}

	return nil
}

// Init cs3 backend vars which are currently not accessible via idp api
func initCS3EnvVars(cs3Addr, machineAuthAPIKey string) error {
	defaults := map[string]string{
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

	var needsAnd bool
	if ldap.Filter != "" {
		filter += ldap.Filter
		needsAnd = true
	}

	if ldap.UserEnabledAttribute != "" {
		// Using a (!(enabled=FALSE)) filter here to allow user without
		// any value for the enable flag to log in
		filter += fmt.Sprintf("(!(%s=FALSE))", ldap.UserEnabledAttribute)
		needsAnd = true
	}

	if needsAnd {
		filter = fmt.Sprintf("(&%s)", filter)
	}

	defaults := map[string]string{
		"LDAP_URI":                 ldap.URI,
		"LDAP_BINDDN":              ldap.BindDN,
		"LDAP_BINDPW":              ldap.BindPassword,
		"LDAP_BASEDN":              ldap.BaseDN,
		"LDAP_SCOPE":               ldap.Scope,
		"LDAP_LOGIN_ATTRIBUTE":     ldap.LoginAttribute,
		"LDAP_EMAIL_ATTRIBUTE":     ldap.EmailAttribute,
		"LDAP_NAME_ATTRIBUTE":      ldap.NameAttribute,
		"LDAP_UUID_ATTRIBUTE":      ldap.UUIDAttribute,
		"LDAP_SUB_ATTRIBUTES":      ldap.UUIDAttribute,
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
	logger     log.Logger
	config     *config.Config
	mux        *chi.Mux
	assets     http.FileSystem
	tp         trace.TracerProvider
	translator l10n.Translator
}

// initMux initializes the internal idp gorilla mux and mounts it in to an ocis chi-router
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
		idp.tp,
	))

	// Login and static pages (must be before Mount so chi matches first)
	idp.mux.Get("/signin/v1/identifier", idp.Index())
	idp.mux.Get("/signin/v1/identifier/", idp.Index())
	idp.mux.Get("/signin/v1/identifier/index.html", idp.Index())
	idp.mux.Get("/signin/v1/welcome", idp.Welcome())
	idp.mux.Get("/signin/v1/goodbye", idp.Goodbye())

	idp.mux.Mount("/", gm)

	_ = chi.Walk(idp.mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		options.Logger.Debug().Str("method", method).Str("route", route).Int("middlewares", len(middlewares)).Msg("serving endpoint")
		return nil
	})
}

// ServeHTTP implements the Service interface.
func (idp *IDP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idp.mux.ServeHTTP(w, r)
}

// Index renders the login page with templated variables.
func (idp *IDP) Index() http.HandlerFunc {
	tpl, err := idp.readTemplate("/identifier/index.html")
	if err != nil {
		idp.logger.Fatal().Err(err).Msg("Could not load index template")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		lang := detectLocale(r)
		t := idp.translator.Locale(lang)
		replacements := map[string]string{
			"__TITLE__":             html.EscapeString(t.Get("Sign in - ownCloud")),
			"__SR_HEADLINE__":       html.EscapeString(t.Get("Login")),
			"__LABEL_USERNAME__":    html.EscapeString(t.Get("Username")),
			"__LABEL_PASSWORD__":    html.EscapeString(t.Get("Password")),
			"__BUTTON_SIGNIN__":     html.EscapeString(t.Get("Sign in")),
			"__BUTTON_SIGNING_IN__": html.EscapeString(t.Get("Signing in…")),
			"__ERR_REQUIRED__":      html.EscapeString(t.Get("Username and password are required.")),
			"__ERR_INVALID__":       html.EscapeString(t.Get("Invalid username or password.")),
			"__ERR_FAILED__":        html.EscapeString(t.Get("Login failed. Please try again.")),
			"__ERR_DEFAULT__":       html.EscapeString(t.Get("Login failed.")),
		}
		body, nonce := idp.renderTemplate(tpl, idp.config.Service.PasswordResetURI, r, replacements)
		writeSecureHTML(w, body, nonce)
	}
}

// Welcome renders the signed-in confirmation page.
func (idp *IDP) Welcome() http.HandlerFunc {
	tpl, err := idp.readTemplate("/identifier/welcome.html")
	if err != nil {
		idp.logger.Fatal().Err(err).Msg("Could not load welcome template")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		lang := detectLocale(r)
		t := idp.translator.Locale(lang)
		replacements := map[string]string{
			"__TITLE__":           html.EscapeString(t.Get("Signed in - ownCloud")),
			"__SR_HEADLINE__":     html.EscapeString(t.Get("Signed in")),
			"__WELCOME_MESSAGE__": html.EscapeString(t.Get("You are signed in. You can close this window and return to the application.")),
		}
		body, nonce := idp.renderTemplate(tpl, "", r, replacements)
		writeSecureHTML(w, body, nonce)
	}
}

// Goodbye renders the signed-out confirmation page.
func (idp *IDP) Goodbye() http.HandlerFunc {
	tpl, err := idp.readTemplate("/identifier/goodbye.html")
	if err != nil {
		idp.logger.Fatal().Err(err).Msg("Could not load goodbye template")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		lang := detectLocale(r)
		t := idp.translator.Locale(lang)
		replacements := map[string]string{
			"__TITLE__":           html.EscapeString(t.Get("Signed out - ownCloud")),
			"__SR_HEADLINE__":     html.EscapeString(t.Get("Signed out")),
			"__GOODBYE_MESSAGE__": html.EscapeString(t.Get("You are signed out. You can close this window.")),
		}
		body, nonce := idp.renderTemplate(tpl, "", r, replacements)
		writeSecureHTML(w, body, nonce)
	}
}

func (idp *IDP) readTemplate(name string) ([]byte, error) {
	f, err := idp.assets.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

var supportedLangs = []language.Tag{
	language.English,
	language.German,
	language.French,
	language.Dutch,
}
var langMatcher = language.NewMatcher(supportedLangs)

func detectLocale(r *http.Request) string {
	accept := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(langMatcher, accept)
	base, _ := tag.Base()
	return base.String()
}

func (idp *IDP) renderTemplate(tpl []byte, passwordResetURI string, r *http.Request, replacements map[string]string) ([]byte, string) {
	lang := detectLocale(r)
	t := idp.translator.Locale(lang)

	pp := []byte("/signin/v1")
	nonce := rndm.GenerateRandomString(32)
	bg := []byte(idp.config.Asset.LoginBackgroundUrl)

	var resetHTML []byte
	if passwordResetURI != "" {
		resetLabel := html.EscapeString(t.Get("Reset password"))
		resetHTML = []byte(`<p><a href="` + html.EscapeString(passwordResetURI) + `">` + resetLabel + `</a></p>`)
	}

	out := bytes.ReplaceAll(tpl, []byte("__PATH_PREFIX__"), pp)
	out = bytes.ReplaceAll(out, []byte("__CSP_NONCE__"), []byte(nonce))
	out = bytes.ReplaceAll(out, []byte("__BG_IMG_URL__"), bg)
	out = bytes.ReplaceAll(out, []byte("__LANG__"), []byte(lang))
	out = bytes.ReplaceAll(out, []byte("__PASSWORD_RESET_LINK_HTML__"), resetHTML)

	for placeholder, value := range replacements {
		out = bytes.ReplaceAll(out, []byte(placeholder), []byte(value))
	}
	return out, nonce
}

func writeSecureHTML(w http.ResponseWriter, body []byte, nonce string) {
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Content-Security-Policy", fmt.Sprintf(
		"default-src 'self'; script-src 'nonce-%s'; style-src 'self' 'nonce-%s'; img-src 'self' data:; font-src 'self'; base-uri 'none'; frame-ancestors 'none';",
		nonce, nonce))
	h.Set("X-Frame-Options", "DENY")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Referrer-Policy", "origin")
	w.WriteHeader(http.StatusOK)
	w.Write(body) //nolint:errcheck
}
