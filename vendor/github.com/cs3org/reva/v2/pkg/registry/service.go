// Copyright 2018-2023 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package registry

import (
	mRegistry "go-micro.dev/v4/registry"
	"go-micro.dev/v4/selector"
)

// DiscoverServicesByAddress searches all available services for nodes that match the given address
func DiscoverServicesByAddress(address string) ([]*mRegistry.Service, error) {
	var services []*mRegistry.Service

	// registry not set, return an empty service map and re-try next time.
	if gRegistry == nil {
		return services, nil
	}

	availableServices, err := gRegistry.ListServices()
	if err != nil {
		return nil, err
	}

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
