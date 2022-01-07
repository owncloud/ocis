package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing Tracing `ocisConfig:"tracing"`
	Log     *Log    `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	Reva Reva `ocisConfig:"reva"`

	Policies              []Policy        `ocisConfig:"policies"`
	OIDC                  OIDC            `ocisConfig:"oidc"`
	TokenManager          TokenManager    `ocisConfig:"token_manager"`
	PolicySelector        *PolicySelector `ocisConfig:"policy_selector"`
	PreSignedURL          PreSignedURL    `ocisConfig:"pre_signed_url"`
	AccountBackend        string          `ocisConfig:"account_backend" env:"PROXY_ACCOUNT_BACKEND_TYPE"`
	UserOIDCClaim         string          `ocisConfig:"user_oidc_claim" env:"PROXY_USER_OIDC_CLAIM"`
	UserCS3Claim          string          `ocisConfig:"user_cs3_claim" env:"PROXY_USER_CS3_CLAIM"`
	MachineAuthAPIKey     string          `ocisConfig:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;PROXY_MACHINE_AUTH_API_KEY"`
	AutoprovisionAccounts bool            `ocisConfig:"auto_provision_accounts" env:"PROXY_AUTOPROVISION_ACCOUNTS"`
	EnableBasicAuth       bool            `ocisConfig:"enable_basic_auth" env:"PROXY_ENABLE_BASIC_AUTH"`
	InsecureBackends      bool            `ocisConfig:"insecure_backends" env:"PROXY_INSECURE_BACKENDS"`
	AuthMiddleware        AuthMiddleware  `ocisConfig:"auth_middleware"`

	Context context.Context
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

// AuthMiddleware configures the proxy http auth middleware.
type AuthMiddleware struct {
	CredentialsByUserAgent map[string]string `ocisConfig:"credentials_by_user_agent"`
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
