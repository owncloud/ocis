package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	Asset Asset  `yaml:"asset"`
	File  string `yaml:"file" env:"WEB_UI_CONFIG_FILE" desc:"Read the ownCloud Web json based configuration from this path/file. The config file takes precedence over WEB_OPTION_xxx environment variables. See the text description for more details." introductionVersion:"pre5.0"`
	Web   Web    `yaml:"web"`
	Apps  map[string]App

	TokenManager *TokenManager `yaml:"token_manager"`

	GatewayAddress string          `yaml:"gateway_addr" env:"WEB_GATEWAY_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	Context        context.Context `yaml:"-"`
}

// Asset defines the available asset configuration.
type Asset struct {
	DeprecatedPath string `yaml:"path" env:"WEB_ASSET_PATH" desc:"Serve ownCloud Web assets from a path on the filesystem instead of the builtin assets." introductionVersion:"pre5.0" deprecationVersion:"5.1.0" removalVersion:"6.0.0" deprecationInfo:"The WEB_ASSET_PATH is deprecated and will be removed in the future." deprecationReplacement:"Use WEB_ASSET_CORE_PATH instead."`
	CorePath       string `yaml:"core_path" env:"WEB_ASSET_CORE_PATH" desc:"Serve ownCloud Web assets from a path on the filesystem instead of the builtin assets. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/web/assets/core" introductionVersion:"5.1"`
	ThemesPath     string `yaml:"themes_path" env:"WEB_ASSET_THEMES_PATH" desc:"Serve ownCloud Web themes from a path on the filesystem instead of the builtin assets. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/web/assets/themes" introductionVersion:"5.1"`
	AppsPath       string `yaml:"apps_path" env:"WEB_ASSET_APPS_PATH" desc:"Serve ownCloud Web apps assets from a path on the filesystem instead of the builtin assets. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/web/assets/apps" introductionVersion:"5.1"`
}

// CustomStyle references additional css to be loaded into ownCloud Web.
type CustomStyle struct {
	Href string `json:"href" yaml:"href"`
}

// CustomScript references an additional script to be loaded into ownCloud Web.
type CustomScript struct {
	Src   string `json:"src" yaml:"src"`
	Async bool   `json:"async,omitempty" yaml:"async"`
}

// CustomTranslation references a json file for overwriting translations in ownCloud Web.
type CustomTranslation struct {
	Url string `json:"url" yaml:"url"`
}

// WebConfig defines the available web configuration for a dynamically rendered config.json.
type WebConfig struct {
	Server        string              `json:"server,omitempty" yaml:"server" env:"OCIS_URL;WEB_UI_CONFIG_SERVER" desc:"URL, where the oCIS APIs are reachable for ownCloud Web." introductionVersion:"pre5.0"`
	Theme         string              `json:"theme,omitempty" yaml:"-"`
	OpenIDConnect OIDC                `json:"openIdConnect,omitempty" yaml:"oidc"`
	Apps          []string            `json:"apps" yaml:"apps"`
	Applications  []Application       `json:"applications,omitempty" yaml:"applications"`
	ExternalApps  []ExternalApp       `json:"external_apps,omitempty" yaml:"external_apps"`
	Options       Options             `json:"options,omitempty" yaml:"options"`
	Styles        []CustomStyle       `json:"styles,omitempty" yaml:"styles"`
	Scripts       []CustomScript      `json:"scripts,omitempty" yaml:"scripts"`
	Translations  []CustomTranslation `json:"customTranslations,omitempty" yaml:"custom_translations"`
}

// OIDC defines the available oidc configuration
type OIDC struct {
	MetadataURL           string `json:"metadata_url,omitempty" yaml:"metadata_url" env:"WEB_OIDC_METADATA_URL" desc:"URL for the OIDC well-known configuration endpoint. Defaults to the oCIS API URL + '/.well-known/openid-configuration'." introductionVersion:"pre5.0"`
	Authority             string `json:"authority,omitempty" yaml:"authority" env:"OCIS_URL;OCIS_OIDC_ISSUER;WEB_OIDC_AUTHORITY" desc:"URL of the OIDC issuer. It defaults to URL of the builtin IDP." introductionVersion:"pre5.0"`
	ClientID              string `json:"client_id,omitempty" yaml:"client_id" env:"OCIS_OIDC_CLIENT_ID;WEB_OIDC_CLIENT_ID" desc:"The OIDC client ID which ownCloud Web uses. This client needs to be set up in your IDP. Note that this setting has no effect when using the builtin IDP." introductionVersion:"pre5.0"`
	ResponseType          string `json:"response_type,omitempty" yaml:"response_type" env:"WEB_OIDC_RESPONSE_TYPE" desc:"The OIDC response type to use for authentication." introductionVersion:"pre5.0"`
	Scope                 string `json:"scope,omitempty" yaml:"scope" env:"WEB_OIDC_SCOPE" desc:"OIDC scopes to request during authentication to authorize access to user details. Defaults to 'openid profile email'. Values are separated by blank. More example values but not limited to are 'address' or 'phone' etc." introductionVersion:"pre5.0"`
	PostLogoutRedirectURI string `json:"post_logout_redirect_uri,omitempty" yaml:"post_logout_redirect_uri" env:"WEB_OIDC_POST_LOGOUT_REDIRECT_URI" desc:"This value needs to point to a valid and reachable web page. The web client will trigger a redirect to that page directly after the logout action. The default value is empty and redirects to the login page." introductionVersion:"pre5.0"`
}

// Application defines an application for the Web app switcher.
type Application struct {
	Icon   string            `json:"icon,omitempty" yaml:"icon"`
	Target string            `json:"target,omitempty" yaml:"target"`
	Title  map[string]string `json:"title,omitempty" yaml:"title"`
	Menu   string            `json:"menu,omitempty" yaml:"menu"`
	URL    string            `json:"url,omitempty" yaml:"url"`
}

// ExternalApp defines an external web app.
//
//	{
//		"name": "hello",
//		"path": "http://localhost:9105/hello.js",
//		  "config": {
//		    "url": "http://localhost:9105"
//		  }
//	 }
type ExternalApp struct {
	ID   string `json:"id,omitempty" yaml:"id"`
	Path string `json:"path,omitempty" yaml:"path"`
	// Config is completely dynamic, because it depends on the extension
	Config map[string]interface{} `json:"config,omitempty" yaml:"config"`
}

// ExternalAppConfig defines an external web app configuration.
type ExternalAppConfig struct {
	URL string `json:"url,omitempty" yaml:"url"`
}

// Web defines the available web configuration.
type Web struct {
	ThemeServer string    `yaml:"theme_server" env:"OCIS_URL;WEB_UI_THEME_SERVER" desc:"Base URL to load themes from. Will be prepended to the theme path." introductionVersion:"pre5.0"`  // used to build Theme in WebConfig
	ThemePath   string    `yaml:"theme_path" env:"WEB_UI_THEME_PATH" desc:"Subpath/file to load the theme. Will be appended to the URL of the theme server." introductionVersion:"pre5.0"` // used to build Theme in WebConfig
	Config      WebConfig `yaml:"config"`
}

// App defines the individual app configuration.
type App struct {
	Disabled bool           `yaml:"disabled"`
	Config   map[string]any `yaml:"config"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;WEB_JWT_SECRET" desc:"The secret to mint and validate jwt tokens." introductionVersion:"pre5.0"`
}
