package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Log defines the available logging configuration.
type Log struct {
	Level  string `mapstructure:"log_level"`
	Pretty bool   `mapstructure:"log_pretty"`
	Color  bool   `mapstructure:"log_color"`
	File   string `mapstructure:"log_file"`
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `mapstructure:"debug_addr"`
	Token  string `mapstructure:"debug_token"`
	Pprof  bool   `mapstructure:"debug_pprof"`
	Zpages bool   `mapstructure:"debug_zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr    string `mapstructure:"http_addr"`
	Root    string `mapstructure:"http_root"`
	TLSCert string `mapstructure:"http_tls_cert"`
	TLSKey  string `mapstructure:"http_tls_key"`
	TLS     bool   `mapstructure:"http_tls"`
}

// Service defines the available service configuration.
type Service struct {
	Name      string `mapstructure:"service_name"`
	Namespace string `mapstructure:"service_namespace"`
	Version   string `mapstructure:"service_version"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"tracing_enabled"`
	Type      string `mapstructure:"tracing_type"`
	Endpoint  string `mapstructure:"tracing_endpoint"`
	Collector string `mapstructure:"tracing_collector"`
	Service   string `mapstructure:"tracing_service"`
}

// Policy enables us to use multiple directors.
type Policy struct {
	Name   string  `mapstructure:"policy_name"`
	Routes []Route `mapstructure:"policy_routes"`
}

// Route define forwarding routes
type Route struct {
	Type        RouteType `mapstructure:"route_type"`
	Endpoint    string    `mapstructure:"route_endpoint"`
	Backend     string    `mapstructure:"route_backend"`
	ApacheVHost bool      `mapstructure:"apache-vhost"`
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
	RouteTypes []RouteType = []RouteType{QueryRoute, RegexRoute, PrefixRoute}
)

// Reva defines all available REVA configuration.
type Reva struct {
	Address    string     `mapstructure:"reva_address"`
	Middleware Middleware `mapstructure:"reva_middleware"`
}

// Middleware configures proxy middlewares.
type Middleware struct {
	Auth Auth `mapstructure:""`
}

// Auth configures proxy http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `mapstructure:""`
}

// Cache is a TTL cache configuration.
type Cache struct {
	Size int `mapstructure:"cache_size"`
	TTL  int `mapstructure:"cache_ttl"`
}

// Config combines all available configuration parts.
type Config struct {
	File                  string          `mapstructure:"file"`
	Log                   Log             `mapstructure:"log"`
	Debug                 Debug           `mapstructure:"debug"`
	HTTP                  HTTP            `mapstructure:"http"`
	Service               Service         `mapstructure:"service"`
	Tracing               Tracing         `mapstructure:"tracing"`
	Policies              []Policy        `mapstructure:"policies"`
	OIDC                  OIDC            `mapstructure:"oidc"`
	TokenManager          TokenManager    `mapstructure:"token_manager"`
	PolicySelector        *PolicySelector `mapstructure:"policy_selector"`
	Reva                  Reva            `mapstructure:"reva"`
	PreSignedURL          PreSignedURL    `mapstructure:"pre_signed_url"`
	AccountBackend        string          `mapstructure:"account_backend"`
	UserOIDCClaim         string          `mapstructure:"user_oidc_claim"`
	UserCS3Claim          string          `mapstructure:"user_cs3_claim"`
	MachineAuthAPIKey     string          `mapstructure:"machine_auth_api_key"`
	AutoprovisionAccounts bool            `mapstructure:"auto_provision_accounts"`
	EnableBasicAuth       bool            `mapstructure:"enable_basic_auth"`
	InsecureBackends      bool            `mapstructure:"insecure_backends"`

	Context    context.Context
	Supervised bool
}

// OIDC is the config for the OpenID-Connect middleware. If set the proxy will try to authenticate every request
// with the configured oidc-provider
type OIDC struct {
	Issuer        string `mapstructure:"oidc_issuer"`
	Insecure      bool   `mapstructure:"oidc_insecure"`
	UserinfoCache Cache  `mapstructure:"oidc_user_info_cache"`
}

// PolicySelector is the toplevel-configuration for different selectors
type PolicySelector struct {
	Static    *StaticSelectorConf
	Migration *MigrationSelectorConf
	Claims    *ClaimsSelectorConf
	Regex     *RegexSelectorConf
}

// StaticSelectorConf is the config for the static-policy-selector
type StaticSelectorConf struct {
	Policy string `mapstructure:"static_selector_policy"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

// PreSignedURL is the config for the presigned url middleware
type PreSignedURL struct {
	AllowedHTTPMethods []string `mapstructure:"allowed_http_methods"`
	Enabled            bool     `mapstructure:"enabled"`
}

// MigrationSelectorConf is the config for the migration-selector
type MigrationSelectorConf struct {
	AccFoundPolicy        string `mapstructure:"acc_found_policy"`
	AccNotFoundPolicy     string `mapstructure:"acc_not_found_policy"`
	UnauthenticatedPolicy string `mapstructure:"unauthenticated_policy"`
}

// ClaimsSelectorConf is the config for the claims-selector
type ClaimsSelectorConf struct {
	DefaultPolicy         string `mapstructure:"default_policy"`
	UnauthenticatedPolicy string `mapstructure:"unauthenticated_policy"`
	SelectorCookieName    string `mapstructure:"selector_cookie_name"`
}

// RegexSelectorConf is the config for the regex-selector
type RegexSelectorConf struct {
	DefaultPolicy         string          `mapstructure:"default_policy"`
	MatchesPolicies       []RegexRuleConf `mapstructure:"matches_policies"`
	UnauthenticatedPolicy string          `mapstructure:"unauthenticated_policy"`
	SelectorCookieName    string          `mapstructure:"selector_cookie_name"`
}

type RegexRuleConf struct {
	Priority int    `mapstructure:"priority"`
	Property string `mapstructure:"property"`
	Match    string `mapstructure:"match"`
	Policy   string `mapstructure:"policy"`
}

// New initializes a new configuration
func New() *Config {
	return &Config{
		HTTP: HTTP{},
	}
}

// DefaultConfig are values stored in the flag set, but moved to a struct.
func DefaultConfig() *Config {
	return &Config{
		File: "",
		Log:  Log{}, // logging config is inherited.
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
	}
}
