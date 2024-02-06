package config

// WopiApp defines the available configuration in order to connect to a WOPI app.
type CS3Api struct {
	Gateway     Gateway     `yaml:"gateway"`
	DataGateway DataGateway `yaml:"datagateway"`
}

type Gateway struct {
	Name string `yaml: "name" env:"COLLABORATION_CS3API_GATEWAY_NAME" desc:"service name of the CS3API gateway"`
}

type DataGateway struct {
	Insecure bool `yaml:"insecure" env:"COLLABORATION_CS3API_DATAGATEWAY_INSECURE" desc:"connect to the CS3API data gateway insecurely"`
}
