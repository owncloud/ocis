package registry

import (
	"os"
	"strings"
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

const (
	registryEnv        = "MICRO_REGISTRY"
	registryAddressEnv = "MICRO_REGISTRY_ADDRESS"
)

var (
	once      sync.Once
	regPlugin string
	reg       registry.Registry
)

func Configure(plugin string) {
	if reg == nil {
		regPlugin = plugin
	}
}

// GetRegistry returns a configured micro registry based on Micro env vars.
// It defaults to mDNS, so mind that systems with mDNS disabled by default (i.e SUSE) will have a hard time
// and it needs to explicitly use etcd. Os awareness for providing a working registry out of the box should be done.
func GetRegistry() registry.Registry {
	once.Do(func() {
		addresses := strings.Split(os.Getenv(registryAddressEnv), ",")
		// prefer env of setting from Configure()
		plugin := os.Getenv(registryEnv)
		if plugin == "" {
			plugin = regPlugin
		}

		switch plugin {
		case "nats":
			reg = natsr.NewRegistry(
				registry.Addrs(addresses...),
			)
		case "kubernetes":
			reg = kubernetesr.NewRegistry(
				registry.Addrs(addresses...),
			)
		case "etcd":
			reg = etcdr.NewRegistry(
				registry.Addrs(addresses...),
			)
		case "consul":
			reg = consulr.NewRegistry(
				registry.Addrs(addresses...),
			)
		case "memory":
			reg = memr.NewRegistry()
		default:
			reg = mdnsr.NewRegistry()
		}
		// No cache needed for in-memory registry
		if plugin != "memory" {
			reg = cache.New(reg, cache.WithTTL(20*time.Second))
		}
	})
	// always use cached registry to prevent registry
	// lookup for every request
	return reg
}
