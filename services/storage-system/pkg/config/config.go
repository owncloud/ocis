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

	GRPC   GRPCConfig `yaml:"grpc"`
	HTTP   HTTPConfig `yaml:"http"`
	Events Events     `yaml:"events"`

	TokenManager     *TokenManager `yaml:"token_manager"`
	Reva             *Reva         `yaml:"reva"`
	SystemUserID     string        `yaml:"system_user_id" env:"OCIS_SYSTEM_USER_ID" desc:"ID of the oCIS storage-system system user. Admins need to set the ID for the STORAGE-SYSTEM system user in this config option which is then used to reference the user. Any reasonable long string is possible, preferably this would be an UUIDv4 format."`
	SystemUserAPIKey string        `yaml:"system_user_api_key" env:"OCIS_SYSTEM_USER_API_KEY" desc:"API key for the STORAGE-SYSTEM system user."`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"STORAGE_SYSTEM_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token."`

	Driver        string  `yaml:"driver" env:"STORAGE_SYSTEM_DRIVER" desc:"The driver which should be used by the service."`
	Drivers       Drivers `yaml:"drivers"`
	DataServerURL string  `yaml:"data_server_url" env:"STORAGE_SYSTEM_DATA_SERVER_URL" desc:"URL of the data server, needs to be reachable by other services using this service."`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;STORAGE_SYSTEM_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;STORAGE_SYSTEM_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORAGE_SYSTEM_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;STORAGE_SYSTEM_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_SYSTEM_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_SYSTEM_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_SYSTEM_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_SYSTEM_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_SYSTEM_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"STORAGE_SYSTEM_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_SYSTEM_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_SYSTEM_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"STORAGE_SYSTEM_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"STORAGE_SYSTEM_GRPC_PROTOCOL" desc:"The transport protocol of the GPRC service."`
}

type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"STORAGE_SYSTEM_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"STORAGE_SYSTEM_HTTP_PROTOCOL" desc:"The transport protocol of the HTTP service."`
}

type Drivers struct {
	OCIS OCISDriver `yaml:"ocis"`
}

type OCISDriver struct {
	// Root is the absolute path to the location of the data
	Root string `yaml:"root" env:"STORAGE_SYSTEM_OCIS_ROOT" desc:"Path for the directory where the STORAGE-SYSTEM service stores it's persistent data."`
}

type Events struct {
	Addr      string `yaml:"endpoint" env:"STORAGE_SYSTEM_EVENTS_ENDPOINT" desc:"The address of the streaming service."`
	ClusterID string `yaml:"cluster" env:"STORAGE_SYSTEM_EVENTS_CLUSTER" desc:"The clusterID of the streaming service. Mandatory when using the NATS service."`
}
