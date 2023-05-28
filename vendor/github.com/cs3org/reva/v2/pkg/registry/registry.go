package registry

import (
	"time"

	mRegistry "go-micro.dev/v4/registry"
)

var (
	gRegistry mRegistry.Registry
)

// Init prepares the service registry
func Init(cfg Configuration, nRegistry mRegistry.Registry) error {
	// fixme: get rid of global registry
	// first come first serves, the first service defines the registry type.
	if gRegistry == nil && nRegistry != nil {
		gRegistry = nRegistry
	}

	rOpts := []mRegistry.RegisterOption{mRegistry.RegisterTTL(time.Minute)}
	for _, service := range cfg.Services {
		if err := gRegistry.Register(service, rOpts...); err != nil {
			return err
		}
	}

	return nil
}

func Ready() bool {
	return gRegistry != nil
}
