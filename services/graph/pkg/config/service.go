package config

// Service defines the available service configuration.
type Service struct {
	Name string `yaml:"name" env:"GRAPH_SERVICE_NAME" desc:"The name of the service." introductionVersion:"%%NEXT%%"`
}
