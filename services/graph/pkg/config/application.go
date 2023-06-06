package config

// Application defines the available graph application configuration.
type Application struct {
	ID          string `yaml:"id" env:"GRAPH_APPLICATION_ID" desc:"The ocis application ID shown in the graph. All app roles are tied to this ID."`
	DisplayName string `yaml:"displayname" env:"GRAPH_APPLICATION_DISPLAYNAME" desc:"The oCIS application name"`
}
