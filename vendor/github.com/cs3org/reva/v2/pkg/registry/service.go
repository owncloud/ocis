package registry

import (
	mRegistry "go-micro.dev/v4/registry"
	"go-micro.dev/v4/selector"
)

// GetServiceByAddress searches all available services for nodes that match the given address
func GetServiceByAddress(address string) ([]*mRegistry.Service, error) {
	// registry not yet ready, return an empty service map and re-try next time.
	if gRegistry == nil {
		return []*mRegistry.Service{}, nil
	}

	availableServices, err := gRegistry.ListServices()
	if err != nil {
		return nil, err
	}

	var services []*mRegistry.Service
	for _, service := range availableServices {
		for _, node := range service.Nodes {
			if node.Address != address {
				continue
			}

			services = append(services, service)
		}
	}

	return services, nil
}

// GetNodeAddress returns a random address from the service nodes
func GetNodeAddress(services []*mRegistry.Service) (string, error) {
	// fixme: roundRobin would be nice, but we need to persist the next closure somehow.
	next := selector.Random(services)
	node, err := next()
	if err != nil {
		return "", err
	}

	return node.Address, err
}
