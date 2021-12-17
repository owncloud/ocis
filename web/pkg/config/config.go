package config

import (
	"context"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"WEB_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"WEB_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"WEB_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"WEB_DEBUG_ZPAGES"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"WEB_HTTP_ADDR"`
	Namespace string
	Root      string `ocisConfig:"root" env:"WEB_HTTP_ROOT"`
	CacheTTL  int    `ocisConfig:"cache_ttl" env:"WEB_CACHE_TTL"`
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;WEB_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;WEB_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;WEB_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;WEB_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"WEB_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;WEB_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;WEB_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;WEB_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;WEB_LOG_FILE"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path" env:"WEB_ASSET_PATH"`
}

// WebConfig defines the available web configuration for a dynamically rendered config.json.
type WebConfig struct {
	Server        string                 `json:"server,omitempty" ocisConfig:"server" env:"OCIS_URL;WEB_UI_CONFIG_SERVER"`
	Theme         string                 `json:"theme,omitempty" ocisConfig:"theme" env:""`
	Version       string                 `json:"version,omitempty" ocisConfig:"version" env:"WEB_UI_CONFIG_VERSION"`
	OpenIDConnect OIDC                   `json:"openIdConnect,omitempty" ocisConfig:"oids"`
	Apps          []string               `json:"apps" ocisConfig:"apps"`
	ExternalApps  []ExternalApp          `json:"external_apps,omitempty" ocisConfig:"external_apps"`
	Options       map[string]interface{} `json:"options,omitempty" ocisConfig:"options"`
}

// OIDC defines the available oidc configuration
type OIDC struct {
	MetadataURL  string `json:"metadata_url,omitempty" ocisConfig:"metadata_url" env:"WEB_OIDC_METADATA_URL"`
	Authority    string `json:"authority,omitempty" ocisConfig:"authority" env:"OCIS_URL;WEB_OIDC_AUTHORITY"`
	ClientID     string `json:"client_id,omitempty" ocisConfig:"client_id" env:"WEB_OIDC_CLIENT_ID"`
	ResponseType string `json:"response_type,omitempty" ocisConfig:"response_type" env:"WEB_OIDC_RESPONSE_TYPE"`
	Scope        string `json:"scope,omitempty" ocisConfig:"scope" env:"WEB_OIDC_SCOPE"`
}

// ExternalApp defines an external web app.
// {
//	"name": "hello",
//	"path": "http://localhost:9105/hello.js",
//	  "config": {
//	    "url": "http://localhost:9105"
//	  }
//  }
type ExternalApp struct {
	ID   string `json:"id,omitempty" ocisConfig:"id"`
	Path string `json:"path,omitempty" ocisConfig:"path"`
	// Config is completely dynamic, because it depends on the extension
	Config map[string]interface{} `json:"config,omitempty" ocisConfig:"config"`
}

// ExternalAppConfig defines an external web app configuration.
type ExternalAppConfig struct {
	URL string `json:"url,omitempty" ocisConfig:"url" env:""`
}

// Web defines the available web configuration.
type Web struct {
	Path        string    `ocisConfig:"path" env:"WEB_UI_PATH"`
	ThemeServer string    `ocisConfig:"theme_server" env:"OCIS_URL;WEB_UI_THEME_SERVER"` // used to build Theme in WebConfig
	ThemePath   string    `ocisConfig:"theme_path" env:"WEB_UI_THEME_PATH"`              // used to build Theme in WebConfig
	Config      WebConfig `ocisConfig:"config"`
}

// Config combines all available configuration parts.
type Config struct {
	Service Service `ocisConfig:"service"`

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	Asset Asset  `ocisConfig:"asset"`
	File  string `ocisConfig:"file" env:"WEB_UI_CONFIG"` // TODO: rename this to a more self explaining string
	Web   Web    `ocisConfig:"web"`

	Context    context.Context
	Supervised bool
}

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9104",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9100",
			Root:      "/",
			Namespace: "com.owncloud.web",
			CacheTTL:  604800, // 7 days
		},
		Service: Service{
			Name: "web",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "web",
		},
		Asset: Asset{
			Path: "",
		},
		Web: Web{
			Path:        "",
			ThemeServer: "https://localhost:9200",
			ThemePath:   "/themes/owncloud/theme.json",
			Config: WebConfig{
				Server:  "https://localhost:9200",
				Theme:   "",
				Version: "0.1.0",
				OpenIDConnect: OIDC{
					MetadataURL:  "",
					Authority:    "https://localhost:9200",
					ClientID:     "web",
					ResponseType: "code",
					Scope:        "openid profile email",
				},
				Apps: []string{"files", "search", "media-viewer", "external"},
			},
		},
	}
}
