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

import (
	"github.com/mitchellh/mapstructure"
)

// Config configures a registry
type Config struct {
	Services map[string]map[string]*service `mapstructure:"services"`
}

// service implements the Service interface. Attributes are exported so that mapstructure can unmarshal values onto them.
type service struct {
	Name  string `mapstructure:"name"`
	Nodes []node `mapstructure:"nodes"`
}

type node struct {
	Address  string            `mapstructure:"address"`
	Metadata map[string]string `mapstructure:"metadata"`
}

// ParseConfig translates Config file values into a Config struct for consumers.
func ParseConfig(m map[string]interface{}) (*Config, error) {
	c := &Config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}

	if len(c.Services) == 0 {
		c.Services = make(map[string]map[string]*service)
	}

	return c, nil
}
