package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Log defines the available logging configuration.
type Log struct {
	Level  string `ocisConfig:"level"`
	Pretty bool   `ocisConfig:"pretty"`
	Color  bool   `ocisConfig:"color"`
	File   string `ocisConfig:"file"`
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr    string `ocisConfig:"addr"`
	Root    string `ocisConfig:"root"`
	TLSCert string `ocisConfig:"tls_cert"`
	TLSKey  string `ocisConfig:"tls_key"`
	TLS     bool   `ocisConfig:"tls"`
}

// Service defines the available service configuration.
type Service struct {
	Name      string `ocisConfig:"name"`
	Namespace string `ocisConfig:"namespace"`
	Version   string `ocisConfig:"version"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
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
	Address    string     `ocisConfig:"address"`
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

// Cache is a TTL cache configuration.
type Cache struct {
	Size int `ocisConfig:"size"`
	TTL  int `ocisConfig:"ttl"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Log                   *shared.Log     `ocisConfig:"log"`
	Debug                 Debug           `ocisConfig:"debug"`
	HTTP                  HTTP            `ocisConfig:"http"`
	Service               Service         `ocisConfig:"service"`
	Tracing               Tracing         `ocisConfig:"tracing"`
	Policies              []Policy        `ocisConfig:"policies"`
	OIDC                  OIDC            `ocisConfig:"oidc"`
	TokenManager          TokenManager    `ocisConfig:"token_manager"`
	PolicySelector        *PolicySelector `ocisConfig:"policy_selector"`
	Reva                  Reva            `ocisConfig:"reva"`
	PreSignedURL          PreSignedURL    `ocisConfig:"pre_signed_url"`
	AccountBackend        string          `ocisConfig:"account_backend"`
	UserOIDCClaim         string          `ocisConfig:"user_oidc_claim"`
	UserCS3Claim          string          `ocisConfig:"user_cs3_claim"`
	MachineAuthAPIKey     string          `ocisConfig:"machine_auth_api_key"`
	AutoprovisionAccounts bool            `ocisConfig:"auto_provision_accounts"`
	EnableBasicAuth       bool            `ocisConfig:"enable_basic_auth"`
	InsecureBackends      bool            `ocisConfig:"insecure_backends"`

	Context    context.Context
	Supervised bool
}

// OIDC is the config for the OpenID-Connect middleware. If set the proxy will try to authenticate every request
// with the configured oidc-provider
type OIDC struct {
	Issuer        string `ocisConfig:"issuer"`
	Insecure      bool   `ocisConfig:"insecure"`
	UserinfoCache Cache  `ocisConfig:"user_info_cache"`
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
	JWTSecret string `ocisConfig:"jwt_secret"`
}

// PreSignedURL is the config for the presigned url middleware
type PreSignedURL struct {
	AllowedHTTPMethods []string `ocisConfig:"allowed_http_methods"`
	Enabled            bool     `ocisConfig:"enabled"`
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

// New initializes a new configuration
func New() *Config {
	return &Config{
		HTTP: HTTP{},
	}
}

// DefaultConfig provides with a working local configuration for a proxy service.
func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:  "0.0.0.0:9205",
			Token: "",
		},
		HTTP: HTTP{
			Addr:    "0.0.0.0:9200",
			Root:    "/",
			TLSCert: path.Join(defaults.BaseDataPath(), "proxy", "server.crt"),
			TLSKey:  path.Join(defaults.BaseDataPath(), "proxy", "server.key"),
			TLS:     true,
		},
		Service: Service{
			Name:      "proxy",
			Namespace: "com.owncloud.web",
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
			UserinfoCache: Cache{
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
		AccountBackend:    "accounts",
		UserOIDCClaim:     "email",
		UserCS3Claim:      "mail",
		MachineAuthAPIKey: "change-me-please",
		//AutoprovisionAccounts: false,
		//EnableBasicAuth:       false,
		//InsecureBackends:      false,
		Context: nil,
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
