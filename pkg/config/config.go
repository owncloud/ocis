package config

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
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
	Addr      string
	Namespace string
	Root      string
	TLSCert   string
	TLSKey    string
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

// Config combines all available configuration parts.
type Config struct {
	File     string
	Log      Log
	Debug    Debug
	HTTP     HTTP
	Tracing  Tracing
	Asset    Asset
	Policies []Policy
	OIDC     *OIDC
}

// OIDC is the config for the OpenID-Connect middleware. If set the proxy will try to authenticate every request
// with the configured oidc-provider
type OIDC struct {
	Endpoint    string
	Realm       string
	SigningAlgs []string
	Insecure    bool
}

// New initializes a new configuration
func New() *Config {
	return &Config{}
}
