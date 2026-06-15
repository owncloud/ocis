package svc

import (
	"bytes"
	"context"
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

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
		logger:           options.Logger,
		config:           options.Config,
		assets:           assetVFS,
		tp:               options.TraceProvider,
		translator:       translator,
		bgImgURL:         safeBgImgURL(options.Config.Asset.LoginBackgroundUrl),
		passwordResetURI: safePasswordResetURI(options.Config.Service.PasswordResetURI),
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
	logger           log.Logger
	config           *config.Config
	mux              *chi.Mux
	assets           http.FileSystem
	tp               trace.TracerProvider
	translator       l10n.Translator
	bgImgURL         template.CSS
	passwordResetURI template.URL
	logoURL          string
	logoOnce         sync.Once
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
		idp.assets,
		idp.tp,
	))

	// Login and static pages (must be before Mount so chi matches first)
	idp.mux.Get("/signin/v1/identifier", idp.Index())
	idp.mux.Get("/signin/v1/identifier/", idp.Index())
	idp.mux.Get("/signin/v1/identifier/index.html", idp.Index())
	idp.mux.Get("/signin/v1/welcome", idp.Welcome())
	idp.mux.Get("/signin/v1/goodbye", idp.Goodbye())
	idp.mux.Get("/signin/v1/consent", idp.Consent())
	idp.mux.Get("/signin/v1/loginerror", idp.LoginError())
	idp.mux.Get("/signin/v1/chooseaccount", idp.Chooseaccount())

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

type basePageData struct {
	Lang, Title, Headline, Nonce, PathPrefix string
	BgImgURL                                 template.CSS
	LogoURL                                  string
}

type indexData struct {
	basePageData
	LabelUsername, LabelPassword                   string
	ButtonSignIn, ButtonSigningIn                  string
	ErrRequired, ErrInvalid, ErrFailed, ErrDefault string
	PasswordResetURI                               template.URL
	ResetLabel                                     string
}

type pageData struct {
	basePageData
	Message string
}

type consentData struct {
	basePageData
	AllowLabel, DenyLabel, ConsentConsequence string
}

type chooseaccountData struct {
	basePageData
	HeadlineSub     string
	UseAnotherLabel string
}

func (idp *IDP) basePage(r *http.Request, title, headline string) (basePageData, l10n.OcisLocale) {
	idp.logoOnce.Do(func() {
		insecure := idp.config.IDP.Insecure || os.Getenv("OCIS_INSECURE") == "true"
		idp.logoURL = fetchLoginLogoURL(idp.config.Commons.OcisURL, insecure)
	})
	lang := detectLocale(r)
	t := idp.translator.Locale(lang)
	return basePageData{
		Lang:       lang,
		Title:      t.Get(title),
		Headline:   t.Get(headline),
		Nonce:      rndm.GenerateRandomString(32),
		PathPrefix: "/signin/v1",
		BgImgURL:   idp.bgImgURL,
		LogoURL:    idp.logoURL,
	}, t
}

func (idp *IDP) renderPage(w http.ResponseWriter, tpl *template.Template, nonce string, data any) {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	writeSecureHTML(w, buf.Bytes(), nonce)
}

func (idp *IDP) mustParseTemplate(name, errMsg string) *template.Template {
	tpl, err := idp.parseTemplate(name)
	if err != nil {
		idp.logger.Fatal().Err(err).Msg(errMsg)
	}
	return tpl
}

// Chooseaccount renders the account picker page.
func (idp *IDP) Chooseaccount() http.HandlerFunc {
	tpl := idp.mustParseTemplate("/identifier/chooseaccount.html", "Could not load chooseaccount template")
	return func(w http.ResponseWriter, r *http.Request) {
		base, t := idp.basePage(r, "Choose an account - ownCloud", "Choose an account")
		data := chooseaccountData{
			basePageData:    base,
			HeadlineSub:     t.Get("to sign in"),
			UseAnotherLabel: t.Get("Use another account"),
		}
		idp.renderPage(w, tpl, base.Nonce, data)
	}
}

// Index renders the login page with templated variables.
func (idp *IDP) Index() http.HandlerFunc {
	tpl := idp.mustParseTemplate("/identifier/index.html", "Could not load index template")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("flow") == "consent" {
			http.Redirect(w, r, "/signin/v1/consent?"+r.URL.RawQuery, http.StatusFound)
			return
		}
		base, t := idp.basePage(r, "Sign in - ownCloud", "Login")
		data := indexData{
			basePageData:     base,
			LabelUsername:    t.Get("Username"),
			LabelPassword:    t.Get("Password"),
			ButtonSignIn:     t.Get("Log in"),
			ButtonSigningIn:  t.Get("Logging in…"),
			ErrRequired:      t.Get("Username and password are required."),
			ErrInvalid:       t.Get("Login failed. Invalid username or password."),
			ErrFailed:        t.Get("Login failed. Invalid username or password."),
			ErrDefault:       t.Get("Login failed. Invalid username or password."),
			PasswordResetURI: idp.passwordResetURI,
			ResetLabel:       t.Get("Reset password"),
		}
		idp.renderPage(w, tpl, base.Nonce, data)
	}
}

// Welcome renders the signed-in confirmation page.
func (idp *IDP) Welcome() http.HandlerFunc {
	tpl := idp.mustParseTemplate("/identifier/welcome.html", "Could not load welcome template")
	return func(w http.ResponseWriter, r *http.Request) {
		base, t := idp.basePage(r, "Signed in - ownCloud", "Signed in")
		data := pageData{
			basePageData: base,
			Message:      t.Get("You are signed in. You can close this window and return to the application."),
		}
		idp.renderPage(w, tpl, base.Nonce, data)
	}
}

// Goodbye renders the signed-out confirmation page.
func (idp *IDP) Goodbye() http.HandlerFunc {
	tpl := idp.mustParseTemplate("/identifier/goodbye.html", "Could not load goodbye template")
	return func(w http.ResponseWriter, r *http.Request) {
		base, t := idp.basePage(r, "Signed out - ownCloud", "Signed out")
		data := pageData{
			basePageData: base,
			Message:      t.Get("You are signed out. You can close this window."),
		}
		idp.renderPage(w, tpl, base.Nonce, data)
	}
}

// LoginError renders the login error page.
func (idp *IDP) LoginError() http.HandlerFunc {
	tpl := idp.mustParseTemplate("/identifier/loginerror.html", "Could not load loginerror template")
	return func(w http.ResponseWriter, r *http.Request) {
		base, t := idp.basePage(r, "Login Error - ownCloud", "Login Error")
		allowedMessages := map[string]string{
			"access_denied":    t.Get("Access denied."),
			"session_expired":  t.Get("Your session has expired."),
			"interaction_required": t.Get("Login required."),
		}
		msg := allowedMessages[r.URL.Query().Get("code")]
		if msg == "" {
			msg = t.Get("Login Error")
		}
		data := pageData{basePageData: base, Message: msg}
		idp.renderPage(w, tpl, base.Nonce, data)
	}
}

// Consent renders the OAuth2 consent page.
func (idp *IDP) Consent() http.HandlerFunc {
	tpl := idp.mustParseTemplate("/identifier/consent.html", "Could not load consent template")
	return func(w http.ResponseWriter, r *http.Request) {
		base, t := idp.basePage(r, "Authorize - ownCloud", "Authorize")
		data := consentData{
			basePageData:       base,
			AllowLabel:         t.Get("Allow"),
			DenyLabel:          t.Get("Cancel"),
			ConsentConsequence: t.Get("By clicking Allow, you allow this app to use your information."),
		}
		idp.renderPage(w, tpl, base.Nonce, data)
	}
}

func (idp *IDP) parseTemplate(name string) (*template.Template, error) {
	f, err := idp.assets.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return template.New(name).Parse(string(b))
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

// safeExternalURL validates that raw is an absolute http/https URL with a
// non-empty host, then returns a normalised form built from the parsed
// structure so that encoding tricks in the raw string cannot bypass the check.
func safeExternalURL(raw string) *url.URL {
	if raw == "" {
		return nil
	}
	u, err := url.ParseRequestURI(raw) // rejects relative refs and opaque URIs
	if err != nil || u.Host == "" {
		return nil
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		return u
	}
	return nil
}

func safeBgImgURL(raw string) template.CSS {
	if u := safeExternalURL(raw); u != nil {
		return template.CSS(u.String()) //nolint:gosec
	}
	return ""
}

func safePasswordResetURI(raw string) template.URL {
	if u := safeExternalURL(raw); u != nil {
		return template.URL(u.String()) //nolint:gosec
	}
	return ""
}

// fetchLoginLogoURL fetches the ownCloud theme.json from the web service at startup
// and returns the login logo URL. Falls back to the default ownCloud logo if unreachable.
func fetchLoginLogoURL(ocisURL string, insecure bool) string {
	const fallback = "/themes/owncloud/assets/logo.svg"
	if ocisURL == "" {
		return fallback
	}
	themeURL := strings.TrimRight(ocisURL, "/") + "/themes/owncloud/theme.json"
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}} //nolint:gosec
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	resp, err := client.Get(themeURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fallback
	}
	defer resp.Body.Close()
	var theme themeJSON
	if err := json.NewDecoder(resp.Body).Decode(&theme); err != nil {
		return fallback
	}
	logo := theme.Common.Logo
	if logo == "" {
		return fallback
	}
	// theme.json returns relative paths like "themes/owncloud/assets/logo.svg"
	if !strings.HasPrefix(logo, "/") && !strings.HasPrefix(logo, "http") {
		logo = "/" + logo
	}
	return logo
}

type themeJSON struct {
	Common struct {
		Logo string `json:"logo"`
	} `json:"common"`
}

func writeSecureHTML(w http.ResponseWriter, body []byte, nonce string) {
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Content-Security-Policy", fmt.Sprintf(
		"default-src 'self'; script-src 'nonce-%s' 'strict-dynamic'; style-src 'self' 'nonce-%s'; img-src 'self'; font-src 'self'; connect-src 'self'; object-src 'none'; form-action 'self'; base-uri 'none'; frame-ancestors 'none';",
		nonce, nonce))
	h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	h.Set("X-Frame-Options", "DENY")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Referrer-Policy", "origin")
	h.Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.WriteHeader(http.StatusOK)
	w.Write(body) //nolint:errcheck
}
