package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing *Tracing `ocisConfig:"tracing"`
	Log     *Log     `ocisConfig:"log"`
	Debug   Debug    `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	Asset Asset  `ocisConfig:"asset"`
	File  string `ocisConfig:"file" env:"WEB_UI_CONFIG"` // TODO: rename this to a more self explaining string
	Web   Web    `ocisConfig:"web"`

	Context context.Context
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
