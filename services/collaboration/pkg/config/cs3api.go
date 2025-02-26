package config

import "github.com/owncloud/ocis/v2/ocis-pkg/shared"

// CS3Api defines the available configuration in order to access to the CS3 gateway.
type CS3Api struct {
	Gateway       Gateway               `yaml:"gateway"`
	DataGateway   DataGateway           `yaml:"datagateway"`
	GRPCClientTLS *shared.GRPCClientTLS `yaml:"grpc_client_tls"`
}

// Gateway defines the available configuration for the CS3 API gateway
type Gateway struct {
	Name string `yaml:"name" env:"OCIS_REVA_GATEWAY" desc:"CS3 gateway used to look up user metadata." introductionVersion:"6.0.0"`
}

// DataGateway defines the available configuration for the CS3 API data gateway
type DataGateway struct {
	Insecure bool `yaml:"insecure" env:"COLLABORATION_CS3API_DATAGATEWAY_INSECURE" desc:"Connect to the CS3API data gateway insecurely." introductionVersion:"6.0.0"`
}
