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

package memory

import "github.com/cs3org/reva/v2/pkg/registry"

// NewService creates a new memory registry.Service.
func NewService(name string, nodes []interface{}) registry.Service {
	n := make([]node, 0)
	for i := 0; i < len(nodes); i++ {
		n = append(n, node{
			// explicit type conversions because types are not exported to prevent from circular dependencies until released.
			id:      nodes[i].(map[string]interface{})["id"].(string),
			address: nodes[i].(map[string]interface{})["address"].(string),
			//metadata: nodes[i].(map[string]interface{})["metadata"].(map[string]string),
		})
	}

	return service{
		name:  name,
		nodes: n,
	}
}

// service implements the Service interface
type service struct {
	name  string
	nodes []node
}

// Name implements the service interface.
func (s service) Name() string {
	return s.name
}

// Nodes implements the service interface.
func (s service) Nodes() []registry.Node {
	ret := make([]registry.Node, 0)
	for i := range s.nodes {
		ret = append(ret, s.nodes[i])
	}
	return ret
}

func (s *service) mergeNodes(n1, n2 []registry.Node) {
	n1 = append(n1, n2...)
	for _, n := range n1 {
		s.nodes = append(s.nodes, node{
			id:       n.ID(),
			address:  n.Address(),
			metadata: n.Metadata(),
		})
	}
}
