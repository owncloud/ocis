package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing"`
	Logging         *Logging `yaml:"log"`
	Debug           Debug    `yaml:"debug"`
	Supervised      bool     `yaml:"-"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	SkipUserGroupsInToken bool    `yaml:"skip_user_groups_in_token"`
	ExternalAddr          string  `yaml:"external_addr"`
	Driver                string  `yaml:"driver"`
	Drivers               Drivers `yaml:"drivers"`
}

type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;APP_PROVIDER_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;APP_PROVIDER_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;APP_PROVIDER_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;APP_PROVIDER_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;APP_PROVIDER_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;APP_PROVIDER_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;APP_PROVIDER_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;APP_PROVIDER_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"APP_PROVIDER_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"APP_PROVIDER_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"APP_PROVIDER_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"APP_PROVIDER_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr     string `yaml:"addr" env:"APP_PROVIDER_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"APP_PROVIDER_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type Drivers struct {
	WOPI WOPIDriver
}

type WOPIDriver struct {
	AppAPIKey      string `yaml:"app_api_key"`
	AppDesktopOnly bool   `yaml:"app_desktop_only"`
	AppIconURI     string `yaml:"app_icon_uri"`
	AppInternalURL string `yaml:"app_internal_url"`
	AppName        string `yaml:"app_name"`
	AppURL         string `yaml:"app_url"`
	Insecure       bool   `yaml:"insecure"`
	IopSecret      string `yaml:"ipo_secret"`
	WopiURL        string `yaml:"wopi_url"`
}
