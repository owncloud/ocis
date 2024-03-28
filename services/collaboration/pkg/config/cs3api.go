package config

// WopiApp defines the available configuration in order to connect to a WOPI app.
type CS3Api struct {
	Gateway     Gateway     `yaml:"gateway"`
	DataGateway DataGateway `yaml:"datagateway"`
}

type Gateway struct {
	Name string `yaml: "name" env:"OCIS_REVA_GATEWAY;COLLABORATION_CS3API_GATEWAY_NAME" desc:"The service name of the CS3API gateway." introductionVersion:"5.1"`
}

type DataGateway struct {
	Insecure bool `yaml:"insecure" env:"COLLABORATION_CS3API_DATAGATEWAY_INSECURE" desc:"Connect to the CS3API data gateway insecurely." introductionVersion:"5.1"`
}
