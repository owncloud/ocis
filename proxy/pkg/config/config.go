package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;PROXY_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;PROXY_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;PROXY_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;PROXY_LOG_FILE"`
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"PROXY_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"PROXY_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"PROXY_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"PROXY_DEBUG_ZPAGES"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"PROXY_HTTP_ADDR"`
	Root      string `ocisConfig:"root" env:"PROXY_HTTP_ROOT"`
	Namespace string
	TLSCert   string `ocisConfig:"tls_cert" env:"PROXY_TRANSPORT_TLS_CERT"`
	TLSKey    string `ocisConfig:"tls_key" env:"PROXY_TRANSPORT_TLS_KEY"`
	TLS       bool   `ocisConfig:"tls" env:"PROXY_TLS"`
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;PROXY_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;PROXY_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;PROXY_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;PROXY_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"PROXY_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Policy enables us to use multiple directors.
type Policy struct {
	Name   string  `ocisConfig:"name"`
	Routes []Route `ocisConfig:"routes"`
}

// Route define forwarding routes
type Route struct {
	Type        RouteType `ocisConfig:"type"`
	Endpoint    string    `ocisConfig:"endpoint"`
	Backend     string    `ocisConfig:"backend"`
	ApacheVHost bool      `ocisConfig:"apache-vhost"`
}

// RouteType defines the type of a route
type RouteType string

const (
	// PrefixRoute are routes matched by a prefix
	PrefixRoute RouteType = "prefix"
	// QueryRoute are routes matched by a prefix and query parameters
	QueryRoute RouteType = "query"
	// RegexRoute are routes matched by a pattern
	RegexRoute RouteType = "regex"
	// DefaultRouteType is the PrefixRoute
	DefaultRouteType RouteType = PrefixRoute
)

var (
	// RouteTypes is an array of the available route types
	RouteTypes = []RouteType{QueryRoute, RegexRoute, PrefixRoute}
)

// Reva defines all available REVA configuration.
type Reva struct {
	Address    string     `ocisConfig:"address" env:"REVA_GATEWAY"`
	Middleware Middleware `ocisConfig:"middleware"`
}

// Middleware configures proxy middlewares.
type Middleware struct {
	Auth Auth `ocisConfig:"middleware"`
}

// Auth configures proxy http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `ocisConfig:""`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service `ocisConfig:"service"`

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	Policies              []Policy        `ocisConfig:"policies"`
	OIDC                  OIDC            `ocisConfig:"oidc"`
	TokenManager          TokenManager    `ocisConfig:"token_manager"`
	PolicySelector        *PolicySelector `ocisConfig:"policy_selector"`
	Reva                  Reva            `ocisConfig:"reva"`
	PreSignedURL          PreSignedURL    `ocisConfig:"pre_signed_url"`
	AccountBackend        string          `ocisConfig:"account_backend" env:"PROXY_ACCOUNT_BACKEND_TYPE"`
	UserOIDCClaim         string          `ocisConfig:"user_oidc_claim" env:"PROXY_USER_OIDC_CLAIM"`
	UserCS3Claim          string          `ocisConfig:"user_cs3_claim" env:"PROXY_USER_CS3_CLAIM"`
	MachineAuthAPIKey     string          `ocisConfig:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;PROXY_MACHINE_AUTH_API_KEY"`
	AutoprovisionAccounts bool            `ocisConfig:"auto_provision_accounts" env:"PROXY_AUTOPROVISION_ACCOUNTS"`
	EnableBasicAuth       bool            `ocisConfig:"enable_basic_auth" env:"PROXY_ENABLE_BASIC_AUTH"`
	InsecureBackends      bool            `ocisConfig:"insecure_backends" env:"PROXY_INSECURE_BACKENDS"`

	Context    context.Context
	Supervised bool
}

// OIDC is the config for the OpenID-Connect middleware. If set the proxy will try to authenticate every request
// with the configured oidc-provider
type OIDC struct {
	Issuer        string        `ocisConfig:"issuer" env:"OCIS_URL;PROXY_OIDC_ISSUER"`
	Insecure      bool          `ocisConfig:"insecure" env:"OCIS_INSECURE;PROXY_OIDC_INSECURE"`
	UserinfoCache UserinfoCache `ocisConfig:"user_info_cache"`
}

// UserinfoCache is a TTL cache configuration.
type UserinfoCache struct {
	Size int `ocisConfig:"size" env:"PROXY_OIDC_USERINFO_CACHE_SIZE"`
	TTL  int `ocisConfig:"ttl" env:"PROXY_OIDC_USERINFO_CACHE_TTL"`
}

// PolicySelector is the toplevel-configuration for different selectors
type PolicySelector struct {
	Static    *StaticSelectorConf    `ocisConfig:"static"`
	Migration *MigrationSelectorConf `ocisConfig:"migration"`
	Claims    *ClaimsSelectorConf    `ocisConfig:"claims"`
	Regex     *RegexSelectorConf     `ocisConfig:"regex"`
}

// StaticSelectorConf is the config for the static-policy-selector
type StaticSelectorConf struct {
	Policy string `ocisConfig:"policy"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET;PROXY_JWT_SECRET"`
}

// PreSignedURL is the config for the presigned url middleware
type PreSignedURL struct {
	AllowedHTTPMethods []string `ocisConfig:"allowed_http_methods"`
	Enabled            bool     `ocisConfig:"enabled" env:"PROXY_ENABLE_PRESIGNEDURLS"`
}

// MigrationSelectorConf is the config for the migration-selector
type MigrationSelectorConf struct {
	AccFoundPolicy        string `ocisConfig:"acc_found_policy"`
	AccNotFoundPolicy     string `ocisConfig:"acc_not_found_policy"`
	UnauthenticatedPolicy string `ocisConfig:"unauthenticated_policy"`
}

// ClaimsSelectorConf is the config for the claims-selector
type ClaimsSelectorConf struct {
	DefaultPolicy         string `ocisConfig:"default_policy"`
	UnauthenticatedPolicy string `ocisConfig:"unauthenticated_policy"`
	SelectorCookieName    string `ocisConfig:"selector_cookie_name"`
}

// RegexSelectorConf is the config for the regex-selector
type RegexSelectorConf struct {
	DefaultPolicy         string          `ocisConfig:"default_policy"`
	MatchesPolicies       []RegexRuleConf `ocisConfig:"matches_policies"`
	UnauthenticatedPolicy string          `ocisConfig:"unauthenticated_policy"`
	SelectorCookieName    string          `ocisConfig:"selector_cookie_name"`
}

type RegexRuleConf struct {
	Priority int    `ocisConfig:"priority"`
	Property string `ocisConfig:"property"`
	Match    string `ocisConfig:"match"`
	Policy   string `ocisConfig:"policy"`
}

// DefaultConfig provides with a working local configuration for a proxy service.
func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:  "0.0.0.0:9205",
			Token: "",
		},
		HTTP: HTTP{
			Addr:      "0.0.0.0:9200",
			Root:      "/",
			Namespace: "com.owncloud.web",
			TLSCert:   path.Join(defaults.BaseDataPath(), "proxy", "server.crt"),
			TLSKey:    path.Join(defaults.BaseDataPath(), "proxy", "server.key"),
			TLS:       true,
		},
		Service: Service{
			Name: "proxy",
		},
		Tracing: Tracing{
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "proxy",
		},
		OIDC: OIDC{
			Issuer:   "https://localhost:9200",
			Insecure: true,
			//Insecure: true,
			UserinfoCache: UserinfoCache{
				Size: 1024,
				TTL:  10,
			},
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		PolicySelector: nil,
		Reva: Reva{
			Address: "127.0.0.1:9142",
		},
		PreSignedURL: PreSignedURL{
			AllowedHTTPMethods: []string{"GET"},
			Enabled:            true,
		},
		AccountBackend:        "accounts",
		UserOIDCClaim:         "email",
		UserCS3Claim:          "mail",
		MachineAuthAPIKey:     "change-me-please",
		AutoprovisionAccounts: false,
		EnableBasicAuth:       false,
		InsecureBackends:      false,
		// TODO: enable
		//Policies: defaultPolicies(),
	}
}

func DefaultPolicies() []Policy {
	return []Policy{
		{
			Name: "ocis",
			Routes: []Route{
				{
					Endpoint: "/",
					Backend:  "http://localhost:9100",
				},
				{
					Endpoint: "/.well-known/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/konnect/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/signin/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/archiver",
					Backend:  "http://localhost:9140",
				},
				{
					Type:     RegexRoute,
					Endpoint: "/ocs/v[12].php/cloud/(users?|groups)", // we have `user`, `users` and `groups` in ocis-ocs
					Backend:  "http://localhost:9110",
				},
				{
					Endpoint: "/ocs/",
					Backend:  "http://localhost:9140",
				},
				{
					Type:     QueryRoute,
					Endpoint: "/remote.php/?preview=1",
					Backend:  "http://localhost:9115",
				},
				{
					Endpoint: "/remote.php/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/dav/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/webdav/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/status.php",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/index.php/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/data",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/app/",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/graph/",
					Backend:  "http://localhost:9120",
				},
				{
					Endpoint: "/graph-explorer",
					Backend:  "http://localhost:9135",
				},
				// if we were using the go micro api gateway we could look up the endpoint in the registry dynamically
				{
					Endpoint: "/api/v0/accounts",
					Backend:  "http://localhost:9181",
				},
				// TODO the lookup needs a better mechanism
				{
					Endpoint: "/accounts.js",
					Backend:  "http://localhost:9181",
				},
				{
					Endpoint: "/api/v0/settings",
					Backend:  "http://localhost:9190",
				},
				{
					Endpoint: "/settings.js",
					Backend:  "http://localhost:9190",
				},
			},
		},
		{
			Name: "oc10",
			Routes: []Route{
				{
					Endpoint: "/",
					Backend:  "http://localhost:9100",
				},
				{
					Endpoint: "/.well-known/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/konnect/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/signin/",
					Backend:  "http://localhost:9130",
				},
				{
					Endpoint: "/archiver",
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint:    "/ocs/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/remote.php/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/dav/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/webdav/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/status.php",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/index.php/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				{
					Endpoint:    "/data",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
			},
		},
	}
}
