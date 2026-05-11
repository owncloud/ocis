package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config is the configuration for the storage-kiteworks service
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *shared.Reva  `yaml:"reva"`

	GracefulShutdownTimeout int `yaml:"graceful_shutdown_timeout" env:"STORAGE_KITEWORKS_GRACEFUL_SHUTDOWN_TIMEOUT" desc:"The number of seconds to wait for the 'storage-kiteworks' service to shutdown cleanly before exiting with an error that gets logged." introductionVersion:"1.0.0"`

	Driver  KiteworksDriver `yaml:"driver"`
	MountID string          `yaml:"mount_id" env:"STORAGE_KITEWORKS_MOUNT_ID" desc:"Mount ID of this storage provider." introductionVersion:"1.0.0"`

	Context context.Context `yaml:"-"`
}

// Log configures the logging
type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_KITEWORKS_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"1.0.0"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_KITEWORKS_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"1.0.0"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_KITEWORKS_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"1.0.0"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_KITEWORKS_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"1.0.0"`
}

// Service holds general service configuration
type Service struct {
	Name string `yaml:"-" env:"STORAGE_KITEWORKS_SERVICE_NAME" desc:"Service name to use. Change this when starting an additional storage provider with a custom configuration to prevent it from colliding with the default 'storage-kiteworks' service." introductionVersion:"1.0.0"`
}

// Debug is the configuration for the debug server
type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_KITEWORKS_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"1.0.0"`
	Token  string `yaml:"token" env:"STORAGE_KITEWORKS_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"1.0.0"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_KITEWORKS_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"1.0.0"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_KITEWORKS_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"1.0.0"`
}

// GRPCConfig is the configuration for the grpc server
type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"STORAGE_KITEWORKS_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"1.0.0"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"OCIS_GRPC_PROTOCOL;STORAGE_KITEWORKS_GRPC_PROTOCOL" desc:"The transport protocol of the GPRC service." introductionVersion:"1.0.0"`
}

// KiteworksDriver holds the Kiteworks-specific driver configuration
type KiteworksDriver struct {
	Endpoint  string `yaml:"endpoint"   env:"STORAGE_KITEWORKS_ENDPOINT"   desc:"Base URL of the Kiteworks server, e.g. https://kiteworks.example.com." introductionVersion:"1.0.0"`
	Insecure  bool   `yaml:"insecure"   env:"STORAGE_KITEWORKS_INSECURE"   desc:"Skip TLS certificate verification (development only)." introductionVersion:"1.0.0"`
	ChunkSize int64  `yaml:"chunk_size" env:"STORAGE_KITEWORKS_CHUNK_SIZE" desc:"Upload chunk size in bytes. Default 5242880 (5 MB)." introductionVersion:"1.0.0"`
}
