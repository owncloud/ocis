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

	GRPC       GRPCConfig    `yaml:"grpc"`
	GrpcClient client.Client `yaml:"-"`

	TokenManager *TokenManager `yaml:"token_manager"`

	Reva                       *shared.Reva          `yaml:"reva"`
	GRPCClientTLS              *shared.GRPCClientTLS `yaml:"grpc_client_tls"`
	Events                     Events                `yaml:"events"`
	Engine                     Engine                `yaml:"engine"`
	Extractor                  Extractor             `yaml:"extractor"`
	ContentExtractionSizeLimit uint64                `yaml:"content_extraction_size_limit" env:"SEARCH_CONTENT_EXTRACTION_SIZE_LIMIT" desc:"Maximum file size in bytes that is allowed for content extraction."`

	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;SEARCH_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services."`

	Context context.Context `yaml:"-"`
}
