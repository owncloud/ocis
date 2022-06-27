// Copyright 2018-2021 CERN
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

// Registry provides with means for dynamically registering services.
type Registry interface {
	// Add registers a Service on the memoryRegistry. Repeated names is allowed, services are distinguished by their metadata.
	Add(Service) error

	// GetService retrieves a Service and all of its nodes by Service name. It returns []*Service because we can have
	// multiple versions of the same Service running alongside each others.
	GetService(string) (Service, error)
}

// Service defines a service.
type Service interface {
	Name() string
	Nodes() []Node
}

// Node defines nodes on a service.
type Node interface {
	// Address where the given node is running.
	Address() string

	// metadata is used in order to differentiate services implementations. For instance an AuthProvider Service could
	// have multiple implementations, basic, bearer ..., metadata would be used to select the Service type depending on
	// its implementation.
	Metadata() map[string]string

	// ID returns the node ID.
	ID() string
}
