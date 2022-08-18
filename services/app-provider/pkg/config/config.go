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
	Reva         *Reva         `yaml:"reva"`

	ExternalAddr string  `yaml:"external_addr" env:"APP_PROVIDER_EXTERNAL_ADDR" desc:"Address of the app provider, where the GATEWAY service can reach it."`
	Driver       string  `yaml:"driver" env:"APP_PROVIDER_DRIVER" desc:"Driver, the APP PROVIDER services uses. Only \"wopi\" is supported as of now."`
	Drivers      Drivers `yaml:"drivers"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}

type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;APP_PROVIDER_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;APP_PROVIDER_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;APP_PROVIDER_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;APP_PROVIDER_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;APP_PROVIDER_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;APP_PROVIDER_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;APP_PROVIDER_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;APP_PROVIDER_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"APP_PROVIDER_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"APP_PROVIDER_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint"`
	Pprof  bool   `yaml:"pprof" env:"APP_PROVIDER_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling"`
	Zpages bool   `yaml:"zpages" env:"APP_PROVIDER_DEBUG_ZPAGES" desc:"Enables zpages, which can  be used for collecting and viewing traces in-memory."`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"APP_PROVIDER_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"APP_PROVIDER_GRPC_PROTOCOL" desc:"The transport protocol of the GPRC service."`
}

type Drivers struct {
	WOPI WOPIDriver `yaml:"wopi" desc:"driver for the CS3org WOPI server"`
}

type WOPIDriver struct {
	AppAPIKey      string `yaml:"app_api_key" env:"APP_PROVIDER_WOPI_APP_API_KEY" desc:"API key for the wopi app."`
	AppDesktopOnly bool   `yaml:"app_desktop_only" env:"APP_PROVIDER_WOPI_APP_DESKTOP_ONLY" desc:"Offer this app only on desktop."`
	AppIconURI     string `yaml:"app_icon_uri" env:"APP_PROVIDER_WOPI_APP_ICON_URI" desc:"URI to an app icon to be used by clients."`
	AppInternalURL string `yaml:"app_internal_url" env:"APP_PROVIDER_WOPI_APP_INTERNAL_URL" desc:"Internal URL to the app, like in your DMZ."`
	AppName        string `yaml:"app_name" env:"APP_PROVIDER_WOPI_APP_NAME" desc:"Human readable app name."`
	AppURL         string `yaml:"app_url" env:"APP_PROVIDER_WOPI_APP_URL" desc:"URL for end users to access the app."`
	Insecure       bool   `yaml:"insecure" env:"APP_PROVIDER_WOPI_INSECURE" desc:"Allow insecure connections to the app."`
	IopSecret      string `yaml:"wopi_server_iop_secret" env:"APP_PROVIDER_WOPI_WOPI_SERVER_IOP_SECRET" desc:"Shared secret of the CS3org WOPI server."`
	WopiURL        string `yaml:"wopi_server_external_url" env:"APP_PROVIDER_WOPI_WOPI_SERVER_EXTERNAL_URL" desc:"External url of the CS3org WOPI server."`
}
