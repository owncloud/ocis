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

import (
	"context"
	"sync"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storage/favorite"
	"github.com/cs3org/reva/v2/pkg/storage/favorite/registry"
)

func init() {
	registry.Register("memory", New)
}

type mgr struct {
	sync.RWMutex
	favorites map[string]map[string]*provider.ResourceId
}

// New returns an instance of the in-memory favorites manager.
func New(m map[string]interface{}) (favorite.Manager, error) {
	return &mgr{favorites: make(map[string]map[string]*provider.ResourceId)}, nil
}

func (m *mgr) ListFavorites(_ context.Context, userID *user.UserId) ([]*provider.ResourceId, error) {
	m.RLock()
	defer m.RUnlock()
	favorites := make([]*provider.ResourceId, 0, len(m.favorites[userID.OpaqueId]))
	for _, id := range m.favorites[userID.OpaqueId] {
		favorites = append(favorites, id)
	}
	return favorites, nil
}

func (m *mgr) SetFavorite(_ context.Context, userID *user.UserId, resourceInfo *provider.ResourceInfo) error {
	m.Lock()
	defer m.Unlock()
	if m.favorites[userID.OpaqueId] == nil {
		m.favorites[userID.OpaqueId] = make(map[string]*provider.ResourceId)
	}
	m.favorites[userID.OpaqueId][resourceInfo.Id.OpaqueId] = resourceInfo.Id
	return nil
}

func (m *mgr) UnsetFavorite(_ context.Context, userID *user.UserId, resourceInfo *provider.ResourceInfo) error {
	m.Lock()
	defer m.Unlock()
	delete(m.favorites[userID.OpaqueId], resourceInfo.Id.OpaqueId)
	return nil
}
