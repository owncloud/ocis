package config

import (
	"context"
)

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
	File   string
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr    string
	Root    string
	TLSCert string
	TLSKey  string
	TLS     bool
}

// Service defines the available service configuration.
type Service struct {
	Name      string
	Namespace string
	Version   string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string
}

// Policy enables us to use multiple directors.
type Policy struct {
	Name   string
	Routes []Route
}

// Route define forwarding routes
type Route struct {
	Type        RouteType
	Endpoint    string
	Backend     string
	ApacheVHost bool `mapstructure:"apache-vhost"`
}

// RouteType defines the type of a route
type RouteType string

const (
	// PrefixRoute are routes matched by a prefix
	PrefixRoute RouteType = "prefix"
	// QueryRoute are routes machted by a prefix and query parameters
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
	Address    string
	Middleware Middleware
}

// Middleware configures proxy middlewares.
type Middleware struct {
	Auth Auth
}

// Auth configures proxy http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string
}

// Cache is a TTL cache configuration.
type Cache struct {
	Size int
	TTL  int
}

// Config combines all available configuration parts.
type Config struct {
	File                  string
	Log                   Log
	Debug                 Debug
	HTTP                  HTTP
	Service               Service
	Tracing               Tracing
	Asset                 Asset
	Policies              []Policy
	OIDC                  OIDC
	TokenManager          TokenManager
	PolicySelector        *PolicySelector `mapstructure:"policy_selector"`
	Reva                  Reva
	PreSignedURL          PreSignedURL
	AccountBackend        string
	UserOIDCClaim         string
	UserCS3Claim          string
	AutoprovisionAccounts bool
	EnableBasicAuth       bool
	InsecureBackends      bool

	Context    context.Context
	Supervised bool
}

// OIDC is the config for the OpenID-Connect middleware. If set the proxy will try to authenticate every request
// with the configured oidc-provider
type OIDC struct {
	Issuer        string
	Insecure      bool
	UserinfoCache Cache
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
	Policy string
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string
}

// PreSignedURL is the config for the presigned url middleware
type PreSignedURL struct {
	AllowedHTTPMethods []string
	Enabled            bool
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
