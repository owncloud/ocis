package config

// Application defines the available graph application configuration.
type Application struct {
	ID          string // is read from store
	DisplayName string `yaml:"displayname" env:"GRAPH_APPLICATION_DISPLAYNAME" desc:"The oCIS application name"`
}
