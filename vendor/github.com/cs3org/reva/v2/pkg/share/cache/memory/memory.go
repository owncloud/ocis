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
	"time"

	"github.com/bluele/gcache"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/share/cache"
	"github.com/cs3org/reva/v2/pkg/share/cache/registry"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	registry.Register("memory", New)
}

type config struct {
	CacheSize int `mapstructure:"cache_size"`
}

type manager struct {
	cache gcache.Cache
}

// New returns an implementation of a resource info cache that stores the objects in memory
func New(m map[string]interface{}) (cache.ResourceInfoCache, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, errors.Wrap(err, "error decoding conf")
	}
	if c.CacheSize == 0 {
		c.CacheSize = 10000
	}

	return &manager{
		cache: gcache.New(c.CacheSize).LFU().Build(),
	}, nil
}

func (m *manager) Get(key string) (*provider.ResourceInfo, error) {
	infoIf, err := m.cache.Get(key)
	if err != nil {
		return nil, err
	}
	return infoIf.(*provider.ResourceInfo), nil
}

func (m *manager) GetKeys(keys []string) ([]*provider.ResourceInfo, error) {
	infos := make([]*provider.ResourceInfo, len(keys))
	for i, key := range keys {
		if infoIf, err := m.cache.Get(key); err == nil {
			infos[i] = infoIf.(*provider.ResourceInfo)
		}
	}
	return infos, nil
}

func (m *manager) Set(key string, info *provider.ResourceInfo) error {
	return m.cache.Set(key, info)
}

func (m *manager) SetWithExpire(key string, info *provider.ResourceInfo, expiration time.Duration) error {
	return m.cache.SetWithExpire(key, info, expiration)
}
