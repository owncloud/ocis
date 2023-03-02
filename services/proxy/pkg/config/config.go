package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `mask:"struct" yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `mask:"struct" yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	Reva          *shared.Reva          `yaml:"reva"`
	GRPCClientTLS *shared.GRPCClientTLS `yaml:"grpc_client_tls"`

	RoleQuotas            map[string]uint64  `yaml:"role_quotas"`
	Policies              []Policy           `yaml:"policies"`
	OIDC                  OIDC               `yaml:"oidc"`
	TokenManager          *TokenManager      `mask:"struct" yaml:"token_manager"`
	PolicySelector        *PolicySelector    `yaml:"policy_selector"`
	PreSignedURL          PreSignedURL       `yaml:"pre_signed_url"`
	AccountBackend        string             `yaml:"account_backend" env:"PROXY_ACCOUNT_BACKEND_TYPE" desc:"Account backend the PROXY service should use. Currently only 'cs3' is possible here."`
	UserOIDCClaim         string             `yaml:"user_oidc_claim" env:"PROXY_USER_OIDC_CLAIM" desc:"The name of an OpenID Connect claim that is used for resolving users with the account backend. The value of the claim must hold a per user unique, stable and non re-assignable identifier. The availability of claims depends on your Identity Provider. There are common claims available for most Identity providers like 'email' or 'preferred_user' but you can also add your own claim."`
	UserCS3Claim          string             `yaml:"user_cs3_claim" env:"PROXY_USER_CS3_CLAIM" desc:"The name of a CS3 user attribute (claim) that should be mapped to the 'user_oidc_claim'. Supported values are 'username', 'mail' and 'userid'."`
	MachineAuthAPIKey     string             `mask:"password" yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;PROXY_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary to access resources from other services."`
	AutoprovisionAccounts bool               `yaml:"auto_provision_accounts" env:"PROXY_AUTOPROVISION_ACCOUNTS" desc:"Set this to 'true' to automatically provision users that do not yet exist in the users service on-demand upon first sign-in. To use this a write-enabled libregraph user backend needs to be setup an running."`
	EnableBasicAuth       bool               `yaml:"enable_basic_auth" env:"PROXY_ENABLE_BASIC_AUTH" desc:"Set this to true to enable 'basic authentication' (username/password)."`
	InsecureBackends      bool               `yaml:"insecure_backends" env:"PROXY_INSECURE_BACKENDS" desc:"Disable TLS certificate validation for all HTTP backend connections."`
	BackendHTTPSCACert    string             `yaml:"backend_https_cacert" env:"PROXY_HTTPS_CACERT" desc:"Path/File for the root CA certificate used to validate the server’s TLS certificate for https enabled backend services."`
	AuthMiddleware        AuthMiddleware     `yaml:"auth_middleware"`
	PoliciesMiddleware    PoliciesMiddleware `yaml:"policies_middleware"`

	Context context.Context `yaml:"-" json:"-"`
}

// Policy enables us to use multiple directors.
type Policy struct {
	Name   string  `yaml:"name"`
	Routes []Route `yaml:"routes"`
}

// Route defines forwarding routes
type Route struct {
	Type RouteType `yaml:"type,omitempty"`
	// Method optionally limits the route to this HTTP method
	Method   string `yaml:"method,omitempty"`
	Endpoint string `yaml:"endpoint,omitempty"`
	// Backend is a static URL to forward the request to
	Backend string `yaml:"backend,omitempty"`
	// Service name to look up in the registry
	Service     string `yaml:"service,omitempty"`
	ApacheVHost bool   `yaml:"apache_vhost,omitempty"`
	Unprotected bool   `yaml:"unprotected,omitempty"`
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

// PoliciesMiddleware configures the proxy policies middleware.
type PoliciesMiddleware struct {
	Enabled bool   `yaml:"enabled" env:"PROXY_POLICIES_ENABLED" desc:"Activates policy evaluation."`
	Query   string `yaml:"query" env:"PROXY_POLICIES_QUERY" desc:"Defines the policy query."`
}

const (
	AccessTokenVerificationNone = "none"
	AccessTokenVerificationJWT  = "jwt"
	// tdb:
	// AccessTokenVerificationIntrospect = "introspect"
)

// OIDC is the config for the OpenID-Connect middleware. If set the proxy will try to authenticate every request
// with the configured oidc-provider
type OIDC struct {
	Issuer                  string        `yaml:"issuer" env:"OCIS_URL;OCIS_OIDC_ISSUER;PROXY_OIDC_ISSUER" desc:"URL of the OIDC issuer. It defaults to URL of the builtin IDP."`
	Insecure                bool          `yaml:"insecure" env:"OCIS_INSECURE;PROXY_OIDC_INSECURE" desc:"Disable TLS certificate validation for connections to the IDP. Note that this is not recommended for production environments."`
	AccessTokenVerifyMethod string        `yaml:"access_token_verify_method" env:"PROXY_OIDC_ACCESS_TOKEN_VERIFY_METHOD" desc:"Sets how OIDC access tokens should be verified. Possible values are 'none' and 'jwt'. When using 'none', no special validation apart from using it for accessing the IPD's userinfo endpoint will be done. When using 'jwt', it tries to parse the access token as a jwt token and verifies the signature using the keys published on the IDP's 'jwks_uri'."`
	UserinfoCache           UserinfoCache `yaml:"user_info_cache"`
	JWKS                    JWKS          `yaml:"jwks"`
	RewriteWellKnown        bool          `yaml:"rewrite_well_known" env:"PROXY_OIDC_REWRITE_WELLKNOWN" desc:"Enables rewriting the /.well-known/openid-configuration to the configured OIDC issuer. Needed by the Desktop Client, Android Client and iOS Client to discover the OIDC provider."`
}

type JWKS struct {
	RefreshInterval   uint64 `yaml:"refresh_interval" env:"PROXY_OIDC_JWKS_REFRESH_INTERVAL" desc:"The interval for refreshing the JWKS (JSON Web Key Set) in minutes in the background via a new HTTP request to the IDP."`
	RefreshTimeout    uint64 `yaml:"refresh_timeout" env:"PROXY_OIDC_JWKS_REFRESH_TIMEOUT" desc:"The timeout in seconds for an outgoing JWKS request."`
	RefreshRateLimit  uint64 `yaml:"refresh_limit" env:"PROXY_OIDC_JWKS_REFRESH_RATE_LIMIT" desc:"Limits the rate in seconds at which refresh requests are performed for unknown keys. This is used to prevent malicious clients from imposing high network load on the IDP via ocis."`
	RefreshUnknownKID bool   `yaml:"refresh_unknown_kid" env:"PROXY_OIDC_JWKS_REFRESH_UNKNOWN_KID" desc:"If set to 'true', the JWKS refresh request will occur every time an unknown KEY ID (KID) is seen. Always set a 'refresh_limit' when enabling this."`
}

// UserinfoCache is a TTL cache configuration.
type UserinfoCache struct {
	Size int `yaml:"size" env:"PROXY_OIDC_USERINFO_CACHE_SIZE" desc:"Cache size for OIDC user info."`
	TTL  int `yaml:"ttl" env:"PROXY_OIDC_USERINFO_CACHE_TTL" desc:"Max TTL in seconds for the OIDC user info cache."`
}

// PolicySelector is the toplevel-configuration for different selectors
type PolicySelector struct {
	Static *StaticSelectorConf `yaml:"static"`
	Claims *ClaimsSelectorConf `yaml:"claims"`
	Regex  *RegexSelectorConf  `yaml:"regex"`
}

// StaticSelectorConf is the config for the static-policy-selector
type StaticSelectorConf struct {
	Policy string `yaml:"policy"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `mask:"password" yaml:"jwt_secret" env:"OCIS_JWT_SECRET;PROXY_JWT_SECRET" desc:"The secret to mint and validate JWT tokens."`
}

// PreSignedURL is the config for the presigned url middleware
type PreSignedURL struct {
	AllowedHTTPMethods []string `yaml:"allowed_http_methods"`
	Enabled            bool     `yaml:"enabled" env:"PROXY_ENABLE_PRESIGNEDURLS" desc:"Allow OCS to get a signing key to sign requests."`
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
