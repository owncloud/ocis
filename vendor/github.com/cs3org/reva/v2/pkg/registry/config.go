package registry

import (
	"github.com/mitchellh/mapstructure"
	mRegistry "go-micro.dev/v4/registry"
)

// Configuration for the service registry.
type Configuration struct {
	Services []*mRegistry.Service
}

// ConfigurationFromMap returns a Configuration based on a map.
func ConfigurationFromMap(cm map[string]interface{}) (Configuration, error) {
	var c Configuration

	err := mapstructure.Decode(cm, &c)
	return c, err
}
