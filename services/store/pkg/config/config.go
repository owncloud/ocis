package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"go-micro.dev/v4/client"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	GRPCClientTLS *shared.GRPCClientTLS `yaml:"grpc_client_tls"`
	GrpcClient    client.Client         `yaml:"-"`

	Datapath string `yaml:"data_path" env:"STORE_DATA_PATH" desc:"The directory where the filesystem storage will store ocis settings. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/store." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`

	Context context.Context `yaml:"-"`
}
