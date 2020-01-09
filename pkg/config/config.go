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
	Root      string
	Namespace string
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

// PhoenixConfig defines the available phoenix configuration for a dynamically rendered config.json.
type PhoenixConfig struct {
	Server        string `json:"server,omitempty"`
	Theme         string `json:"theme,omitempty"`
	Version       string `json:"version,omitempty"` // TODO what is version used for?
	OpenIDConnect OIDC   `json:"openIdConnect,omitempty"`
	// TODO add nilasempty when https://go-review.googlesource.com/c/go/+/205897/ is released
	Apps         []string      `json:"apps"`
	ExternalApps []ExternalApp `json:"external_apps,omitempty"`
}

// OIDC defines the available oidc configuration
type OIDC struct {
	MetadataURL  string `json:"metadata_url,omitempty"`
	Authority    string `json:"authority,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ResponseType string `json:"response_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// ExternalApp defines an external phoenix app.
// {
//	"name": "hello",
//	"path": "http://localhost:9105/hello.js",
//	  "config": {
//	    "url": "http://localhost:9105"
//	  }
//  }
type ExternalApp struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
	// Config is completely dynamic, because it depends on the extension
	Config map[string]interface{} `json:"config,omitempty"`
}

// ExternalAppConfig defines an external phoenix app configuration.
type ExternalAppConfig struct {
	URL string `json:"url,omitempty"`
}

// Phoenix defines the available phoenix configuration.
type Phoenix struct {
	Path   string
	Config PhoenixConfig
}

// Config combines all available configuration parts.
type Config struct {
	File    string
	Log     Log
	Debug   Debug
	HTTP    HTTP
	Tracing Tracing
	Asset   Asset
	OIDC    OIDC
	Phoenix Phoenix
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
