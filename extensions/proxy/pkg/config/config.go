package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	Reva Reva `yaml:"reva"`

	Policies              []Policy        `yaml:"policies"`
	OIDC                  OIDC            `yaml:"oidc"`
	TokenManager          TokenManager    `yaml:"token_manager"`
	PolicySelector        *PolicySelector `yaml:"policy_selector"`
	PreSignedURL          PreSignedURL    `yaml:"pre_signed_url"`
	AccountBackend        string          `yaml:"account_backend" env:"PROXY_ACCOUNT_BACKEND_TYPE"`
	UserOIDCClaim         string          `yaml:"user_oidc_claim" env:"PROXY_USER_OIDC_CLAIM"`
	UserCS3Claim          string          `yaml:"user_cs3_claim" env:"PROXY_USER_CS3_CLAIM"`
	MachineAuthAPIKey     string          `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;PROXY_MACHINE_AUTH_API_KEY"`
	AutoprovisionAccounts bool            `yaml:"auto_provision_accounts" env:"PROXY_AUTOPROVISION_ACCOUNTS"`
	EnableBasicAuth       bool            `yaml:"enable_basic_auth" env:"PROXY_ENABLE_BASIC_AUTH"`
	InsecureBackends      bool            `yaml:"insecure_backends" env:"PROXY_INSECURE_BACKENDS"`
	AuthMiddleware        AuthMiddleware  `yaml:"auth_middleware"`

	Context context.Context `yaml:"-"`
}

// Policy enables us to use multiple directors.
type Policy struct {
	Name   string  `yaml:"name"`
	Routes []Route `yaml:"routes"`
}

// Route defines forwarding routes
type Route struct {
	Type RouteType `yaml:"type"`
	// Method optionally limits the route to this HTTP method
	Method   string `yaml:"method"`
	Endpoint string `yaml:"endpoint"`
	// Backend is a static URL to forward the request to
	Backend string `yaml:"backend"`
	// Service name to look up in the registry
	Service     string `yaml:"service"`
	ApacheVHost bool   `yaml:"apache-vhost"`
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

// AuthMiddleware configures the proxy http auth middleware.
type AuthMiddleware struct {
	CredentialsByUserAgent map[string]string `yaml:"credentials_by_user_agent"`
}

// OIDC is the config for the OpenID-Connect middleware. If set the proxy will try to authenticate every request
// with the configured oidc-provider
type OIDC struct {
	Issuer        string        `yaml:"issuer" env:"OCIS_URL;PROXY_OIDC_ISSUER"`
	Insecure      bool          `yaml:"insecure" env:"OCIS_INSECURE;PROXY_OIDC_INSECURE"`
	UserinfoCache UserinfoCache `yaml:"user_info_cache"`
}

// UserinfoCache is a TTL cache configuration.
type UserinfoCache struct {
	Size int `yaml:"size" env:"PROXY_OIDC_USERINFO_CACHE_SIZE"`
	TTL  int `yaml:"ttl" env:"PROXY_OIDC_USERINFO_CACHE_TTL"`
}

// PolicySelector is the toplevel-configuration for different selectors
type PolicySelector struct {
	Static    *StaticSelectorConf    `yaml:"static"`
	Migration *MigrationSelectorConf `yaml:"migration"`
	Claims    *ClaimsSelectorConf    `yaml:"claims"`
	Regex     *RegexSelectorConf     `yaml:"regex"`
}

// StaticSelectorConf is the config for the static-policy-selector
type StaticSelectorConf struct {
	Policy string `yaml:"policy"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;PROXY_JWT_SECRET"`
}

// PreSignedURL is the config for the presigned url middleware
type PreSignedURL struct {
	AllowedHTTPMethods []string `yaml:"allowed_http_methods"`
	Enabled            bool     `yaml:"enabled" env:"PROXY_ENABLE_PRESIGNEDURLS"`
}

// MigrationSelectorConf is the config for the migration-selector
type MigrationSelectorConf struct {
	AccFoundPolicy        string `yaml:"acc_found_policy"`
	AccNotFoundPolicy     string `yaml:"acc_not_found_policy"`
	UnauthenticatedPolicy string `yaml:"unauthenticated_policy"`
}

// ClaimsSelectorConf is the config for the claims-selector
type ClaimsSelectorConf struct {
	DefaultPolicy         string `yaml:"default_policy"`
	UnauthenticatedPolicy string `yaml:"unauthenticated_policy"`
	SelectorCookieName    string `yaml:"selector_cookie_name"`
}

// RegexSelectorConf is the config for the regex-selector
type RegexSelectorConf struct {
	DefaultPolicy         string          `yaml:"default_policy"`
	MatchesPolicies       []RegexRuleConf `yaml:"matches_policies"`
	UnauthenticatedPolicy string          `yaml:"unauthenticated_policy"`
	SelectorCookieName    string          `yaml:"selector_cookie_name"`
}

type RegexRuleConf struct {
	Priority int    `yaml:"priority"`
	Property string `yaml:"property"`
	Match    string `yaml:"match"`
	Policy   string `yaml:"policy"`
}
