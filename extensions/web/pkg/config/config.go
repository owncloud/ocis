package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing,omitempty"`
	Log     *Log     `yaml:"log,omitempty"`
	Debug   Debug    `yaml:"debug,omitempty"`

	HTTP HTTP `yaml:"http,omitempty"`

	Asset Asset  `yaml:"asset,omitempty"`
	File  string `yaml:"file,omitempty" env:"WEB_UI_CONFIG"` // TODO: rename this to a more self explaining string
	Web   Web    `yaml:"web,omitempty"`

	Context context.Context `yaml:"-"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `yaml:"path" env:"WEB_ASSET_PATH"`
}

// WebConfig defines the available web configuration for a dynamically rendered config.json.
type WebConfig struct {
	Server        string                 `json:"server,omitempty" yaml:"server" env:"OCIS_URL;WEB_UI_CONFIG_SERVER"`
	Theme         string                 `json:"theme,omitempty" yaml:"theme" env:""`
	Version       string                 `json:"version,omitempty" yaml:"version" env:"WEB_UI_CONFIG_VERSION"`
	OpenIDConnect OIDC                   `json:"openIdConnect,omitempty" yaml:"oids"`
	Apps          []string               `json:"apps" yaml:"apps"`
	ExternalApps  []ExternalApp          `json:"external_apps,omitempty" yaml:"external_apps"`
	Options       map[string]interface{} `json:"options,omitempty" yaml:"options"`
}

// OIDC defines the available oidc configuration
type OIDC struct {
	MetadataURL  string `json:"metadata_url,omitempty" yaml:"metadata_url" env:"WEB_OIDC_METADATA_URL"`
	Authority    string `json:"authority,omitempty" yaml:"authority" env:"OCIS_URL;WEB_OIDC_AUTHORITY"`
	ClientID     string `json:"client_id,omitempty" yaml:"client_id" env:"WEB_OIDC_CLIENT_ID"`
	ResponseType string `json:"response_type,omitempty" yaml:"response_type" env:"WEB_OIDC_RESPONSE_TYPE"`
	Scope        string `json:"scope,omitempty" yaml:"scope" env:"WEB_OIDC_SCOPE"`
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
	ID   string `json:"id,omitempty" yaml:"id"`
	Path string `json:"path,omitempty" yaml:"path"`
	// Config is completely dynamic, because it depends on the extension
	Config map[string]interface{} `json:"config,omitempty" yaml:"config"`
}

// ExternalAppConfig defines an external web app configuration.
type ExternalAppConfig struct {
	URL string `json:"url,omitempty" yaml:"url" env:""`
}

// Web defines the available web configuration.
type Web struct {
	Path        string    `yaml:"path" env:"WEB_UI_PATH"`
	ThemeServer string    `yaml:"theme_server" env:"OCIS_URL;WEB_UI_THEME_SERVER"` // used to build Theme in WebConfig
	ThemePath   string    `yaml:"theme_path" env:"WEB_UI_THEME_PATH"`              // used to build Theme in WebConfig
	Config      WebConfig `yaml:"config"`
}
