package registry

import (
	"os"
	"strings"

	etcdr "github.com/asim/go-micro/plugins/registry/etcd/v4"
	kubernetesr "github.com/asim/go-micro/plugins/registry/kubernetes/v4"
	mdnsr "github.com/asim/go-micro/plugins/registry/mdns/v4"
	natsr "github.com/asim/go-micro/plugins/registry/nats/v4"

	"go-micro.dev/v4/registry"
)

var (
	registryEnv        = "MICRO_REGISTRY"
	registryAddressEnv = "MICRO_REGISTRY_ADDRESS"
)

// GetRegistry returns a configured micro registry based on Micro env vars.
// It defaults to mDNS, so mind that systems with mDNS disabled by default (i.e SUSE) will have a hard time
// and it needs to explicitly use etcd. Os awareness for providing a working registry out of the box should be done.
func GetRegistry() registry.Registry {
	addresses := strings.Split(os.Getenv(registryAddressEnv), ",")

	var r registry.Registry
	switch os.Getenv(registryEnv) {
	case "nats":
		r = natsr.NewRegistry(
			registry.Addrs(addresses...),
		)
	case "kubernetes":
		r = kubernetesr.NewRegistry(
			registry.Addrs(addresses...),
		)
	case "etcd":
		r = etcdr.NewRegistry(
			registry.Addrs(addresses...),
		)
	default:
		r = mdnsr.NewRegistry()
	}

	return r
}
