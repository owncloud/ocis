package config

// CS3Api defines the available configuration in order to access to the CS3 gateway.
type CS3Api struct {
	Gateway     Gateway     `yaml:"gateway"`
	DataGateway DataGateway `yaml:"datagateway"`
}

// Gateway defines the available configuration for the CS3 API gateway
type Gateway struct {
	Name string `yaml:"name" env:"OCIS_REVA_GATEWAY;COLLABORATION_CS3API_GATEWAY_NAME" desc:"CS3 gateway used to look up user metadata." introductionVersion:"6.0.0" deprecationVersion:"6.0" removalVersion:"%%NEXT_PRODUCTION_VERSION%%" deprecationInfo:"COLLABORATION_CS3API_GATEWAY_NAME removed for simplicity."`
}

// DataGateway defines the available configuration for the CS3 API data gateway
type DataGateway struct {
	Insecure bool `yaml:"insecure" env:"COLLABORATION_CS3API_DATAGATEWAY_INSECURE" desc:"Connect to the CS3API data gateway insecurely." introductionVersion:"6.0.0"`
}
