package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr"`
	Root      string `ocisConfig:"root"`
	Namespace string `ocisConfig:"namespace"`
	CacheTTL  int    `ocisConfig:"cache_ttl"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path"`
}

// WebConfig defines the available web configuration for a dynamically rendered config.json.
type WebConfig struct {
	Server        string                 `json:"server,omitempty" ocisConfig:"server"`
	Theme         string                 `json:"theme,omitempty" ocisConfig:"theme"`
	Version       string                 `json:"version,omitempty" ocisConfig:"version"`
	OpenIDConnect OIDC                   `json:"openIdConnect,omitempty" ocisConfig:"oids"`
	Apps          []string               `json:"apps" ocisConfig:"apps"`
	ExternalApps  []ExternalApp          `json:"external_apps,omitempty" ocisConfig:"external_apps"`
	Options       map[string]interface{} `json:"options,omitempty" ocisConfig:"options"`
}

// OIDC defines the available oidc configuration
type OIDC struct {
	MetadataURL  string `json:"metadata_url,omitempty" ocisConfig:"metadata_url"`
	Authority    string `json:"authority,omitempty" ocisConfig:"authority"`
	ClientID     string `json:"client_id,omitempty" ocisConfig:"client_id"`
	ResponseType string `json:"response_type,omitempty" ocisConfig:"response_type"`
	Scope        string `json:"scope,omitempty" ocisConfig:"scope"`
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
	URL string `json:"url,omitempty" ocisConfig:"url"`
}

// Web defines the available web configuration.
type Web struct {
	Path        string    `ocisConfig:"path"`
	ThemeServer string    `ocisConfig:"theme_server"` // used to build Theme in WebConfig
	ThemePath   string    `ocisConfig:"theme_path"`   // used to build Theme in WebConfig
	Config      WebConfig `ocisConfig:"config"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File    string      `ocisConfig:"file"`
	Log     *shared.Log `ocisConfig:"log"`
	Debug   Debug       `ocisConfig:"debug"`
	HTTP    HTTP        `ocisConfig:"http"`
	Tracing Tracing     `ocisConfig:"tracing"`
	Asset   Asset       `ocisConfig:"asset"`
	Web     Web         `ocisConfig:"web"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
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
