/*
 * Copyright 2017-2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package managers

import (
	"fmt"
)

// ServiceUsesManagers is an interface for service which register to managers.
type ServiceUsesManagers interface {
	RegisterManagers(mgrs *Managers) error
}

// Managers is a registry for named managers.
type Managers struct {
	registry map[string]interface{}
}

// New creates a new Managers.
func New() *Managers {
	return &Managers{
		registry: make(map[string]interface{}),
	}
}

// Set adds the provided manager with the provided name to the accociated
// Managers.
func (m *Managers) Set(name string, manager interface{}) {
	m.registry[name] = manager
}

// Get returns the manager identified by the given name from the accociated
// managers.
func (m *Managers) Get(name string) (interface{}, bool) {
	manager, ok := m.registry[name]

	return manager, ok
}

// Must returns the manager indentified by the given name or panics.
func (m *Managers) Must(name string) interface{} {
	manager, ok := m.Get(name)
	if !ok {
		panic(fmt.Errorf("manager %s not found", name))
	}

	return manager
}

// Apply registers the accociated manager's registered managers.
func (m *Managers) Apply() error {
	for _, manager := range m.registry {
		if service, ok := manager.(ServiceUsesManagers); ok {
			err := service.RegisterManagers(m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
