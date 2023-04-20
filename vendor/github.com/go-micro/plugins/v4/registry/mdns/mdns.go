// Package mdns provides a multicast dns registry
package mdns

import (
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/util/cmd"
)

func init() {
	cmd.DefaultRegistries["mdns"] = NewRegistry
}

// NewRegistry returns a new mdns registry.
func NewRegistry(opts ...registry.Option) registry.Registry {
	return registry.NewRegistry(opts...)
}
