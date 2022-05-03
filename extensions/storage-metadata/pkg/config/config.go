package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing"`
	Log             *Log     `yaml:"log"`
	Debug           Debug    `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`
	HTTP HTTPConfig `yaml:"http"`

	TokenManager      *TokenManager `yaml:"token_manager"`
	Reva              *Reva         `yaml:"reva"`
	MachineAuthAPIKey string        `yaml:"machine_auth_api_key" env:"STORAGE_METADATA_MACHINE_AUTH_API_KEY"`
	MetadataUserID    string        `yaml:"metadata_user_id"`

	SkipUserGroupsInToken bool `yaml:"-"`

	Driver               string  `yaml:"driver" env:"STORAGE_METADATA_DRIVER" desc:"The driver which should be used by the service"`
	Drivers              Drivers `yaml:"drivers"`
	DataServerURL        string  `yaml:"data_server_url" env:"STORAGE_METADATA_DATA_SERVER_URL"`
	TempFolder           string  `yaml:"temp_folder" env:"STORAGE_METADATA_TEMP_FOLDER"`
	DataProviderInsecure bool    `yaml:"data_provider_insecure" env:"OCIS_INSECURE;STORAGE_METADATA_DATAPROVIDER_INSECURE"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;STORAGE_METADATA_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;STORAGE_METADATA_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORAGE_METADATA_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;STORAGE_METADATA_TRACING_COLLECTOR"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_METADATA_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_METADATA_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_METADATA_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_METADATA_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_METADATA_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"STORAGE_METADATA_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_METADATA_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_METADATA_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"STORAGE_METADATA_GRPC_ADDR" desc:"The address of the grpc service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"STORAGE_METADATA_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"STORAGE_METADATA_HTTP_ADDR" desc:"The address of the http service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"STORAGE_METADATA_HTTP_PROTOCOL" desc:"The transport protocol of the http service."`
}

type Drivers struct {
	OCIS OCISDriver `yaml:"ocis"`
}

type OCISDriver struct {
	// Root is the absolute path to the location of the data
	Root                string `yaml:"root" env:"STORAGE_METADATA_OCIS_ROOT"`
	UserLayout          string `yaml:"-"`
	PermissionsEndpoint string `yaml:"permissions_endpoint" env:"STORAGE_METADATA_OCIS_PERMISSIONS_ENDPOINT"`
}
