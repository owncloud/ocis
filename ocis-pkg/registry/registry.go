package registry

import (
	"errors"
	"sync"
	"time"

	consulr "github.com/go-micro/plugins/v4/registry/consul"
	etcdr "github.com/go-micro/plugins/v4/registry/etcd"
	kubernetesr "github.com/go-micro/plugins/v4/registry/kubernetes"
	mdnsr "github.com/go-micro/plugins/v4/registry/mdns"
	memr "github.com/go-micro/plugins/v4/registry/memory"
	natsr "github.com/go-micro/plugins/v4/registry/nats"

	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/registry/cache"
)

var (
	once      sync.Once
	regPlugin string
	reg       registry.Registry
)

// Registry defines the parameters for a go micro registry
type Registry struct {
	Type      string
	Addresses []string
}

// GetRegistry returns a configured micro registry based on the given function parameters.
func GetRegistry(r Registry) (registry.Registry, error) {
	if regPlugin != "" && regPlugin != r.Type {
		return nil, errors.New("registry has already been configured to a different registry")
	}

	once.Do(func() {
		switch r.Type {
		case "nats":
			reg = natsr.NewRegistry(
				registry.Addrs(r.Addresses...),
			)
			regPlugin = "nats"
		case "kubernetes":
			reg = kubernetesr.NewRegistry(
				registry.Addrs(r.Addresses...),
			)
			regPlugin = "kubernetes"
		case "etcd":
			reg = etcdr.NewRegistry(
				registry.Addrs(r.Addresses...),
			)
			regPlugin = "etcd"
		case "consul":
			reg = consulr.NewRegistry(
				registry.Addrs(r.Addresses...),
			)
			regPlugin = "consul"
		case "memory":
			reg = memr.NewRegistry()
			regPlugin = "memory"
		case "mdns":
			reg = mdnsr.NewRegistry()
			regPlugin = "mdns"
		default:
			reg = nil
		}
		// No cache needed for in-memory registry
		if r.Type != "memory" {
			// otherwise use cached registry to prevent registry
			// lookup for every request
			reg = cache.New(reg, cache.WithTTL(20*time.Second))
		}
	})
	if reg == nil {
		return nil, errors.New("unknown registry")
	}

	return reg, nil
}
