package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *shared.Reva  `yaml:"reva"`

	ExternalAddr string  `yaml:"external_addr" env:"APP_PROVIDER_EXTERNAL_ADDR" desc:"Address of the app provider, where the GATEWAY service can reach it."`
	Driver       string  `yaml:"driver" env:"APP_PROVIDER_DRIVER" desc:"Driver, the APP PROVIDER services uses. Only 'wopi' is supported as of now."`
	Drivers      Drivers `yaml:"drivers"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;APP_PROVIDER_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;APP_PROVIDER_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;APP_PROVIDER_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;APP_PROVIDER_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"APP_PROVIDER_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"APP_PROVIDER_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint"`
	Pprof  bool   `yaml:"pprof" env:"APP_PROVIDER_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling"`
	Zpages bool   `yaml:"zpages" env:"APP_PROVIDER_DEBUG_ZPAGES" desc:"Enables zpages, which can  be used for collecting and viewing traces in-memory."`
}

type Service struct {
	Name string `yaml:"name" env:"APP_PROVIDER_SERVICE_NAME" desc:"The name of the service. This needs to be changed when using more than one app provider. Each app provider configured needs to be identified by a unique service name. Possible examples are: 'app-provider-collabora', 'app-provider-onlyoffice', 'app-provider-office365'."`
}

type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"APP_PROVIDER_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"APP_PROVIDER_GRPC_PROTOCOL" desc:"The transport protocol of the GPRC service."`
}

type Drivers struct {
	WOPI WOPIDriver `yaml:"wopi" desc:"Driver for the CS3org WOPI server"`
}

type WOPIDriver struct {
	AppAPIKey                 string `yaml:"app_api_key" env:"APP_PROVIDER_WOPI_APP_API_KEY" desc:"API key for the wopi app."`
	AppDesktopOnly            bool   `yaml:"app_desktop_only" env:"APP_PROVIDER_WOPI_APP_DESKTOP_ONLY" desc:"Offer this app only on desktop."`
	AppIconURI                string `yaml:"app_icon_uri" env:"APP_PROVIDER_WOPI_APP_ICON_URI" desc:"URI to an app icon to be used by clients."`
	AppInternalURL            string `yaml:"app_internal_url" env:"APP_PROVIDER_WOPI_APP_INTERNAL_URL" desc:"Internal URL to the app, like in your DMZ."`
	AppName                   string `yaml:"app_name" env:"APP_PROVIDER_WOPI_APP_NAME" desc:"Human readable app name."`
	AppURL                    string `yaml:"app_url" env:"APP_PROVIDER_WOPI_APP_URL" desc:"URL for end users to access the app."`
	AppDisableChat            bool   `yaml:"app_disable_chat" env:"APP_PROVIDER_WOPI_DISABLE_CHAT" desc:"Disable the chat functionality of the office app."`
	Insecure                  bool   `yaml:"insecure" env:"APP_PROVIDER_WOPI_INSECURE" desc:"Disable TLS certificate validation for requests to the WOPI server and the web office application. Do not set this in production environments."`
	IopSecret                 string `yaml:"wopi_server_iop_secret" env:"APP_PROVIDER_WOPI_WOPI_SERVER_IOP_SECRET" desc:"Shared secret of the CS3org WOPI server."`
	WopiURL                   string `yaml:"wopi_server_external_url" env:"APP_PROVIDER_WOPI_WOPI_SERVER_EXTERNAL_URL" desc:"External url of the CS3org WOPI server."`
	WopiFolderURLBaseURL      string `yaml:"wopi_folder_url_base_url" env:"OCIS_URL;APP_PROVIDER_WOPI_FOLDER_URL_BASE_URL" desc:"Base url to navigate back from the app to the containing folder in the file list."`
	WopiFolderURLPathTemplate string `yaml:"wopi_folder_url_path_template" env:"APP_PROVIDER_WOPI_FOLDER_URL_PATH_TEMPLATE" desc:"Path template to navigate back from the app to the containing folder in the file list. Supported template variables are {{.ResourceInfo.ResourceID}}, {{.ResourceInfo.Mtime.Seconds}}, {{.ResourceInfo.Name}}, {{.ResourceInfo.Path}}, {{.ResourceInfo.Type}}, {{.ResourceInfo.Id.SpaceId}}, {{.ResourceInfo.Id.StorageId}}, {{.ResourceInfo.Id.OpaqueId}}, {{.ResourceInfo.MimeType}}"`
}
