package config

// Service defines the available service configuration.
type Service struct {
	Name string `yaml:"-"`

	ApplicationID          string `yaml:"application_id" env:"GRAPH_APPLICATION_ID" desc:"The ocis web application id"` // TODO actually this is the application id for ocis web, and ocis web also needs to know it
	ApplicationDisplayName string `yaml:"application_displayname" env:"GRAPH_APPLICATION_DISPLAYNAME" desc:"The ocis web application name"`
}
