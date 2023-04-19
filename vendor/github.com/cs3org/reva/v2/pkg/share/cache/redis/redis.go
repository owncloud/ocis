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

package redis

import (
	"encoding/json"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/share/cache"
	"github.com/cs3org/reva/v2/pkg/share/cache/registry"
	"github.com/gomodule/redigo/redis"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	registry.Register("redis", New)
}

type config struct {
	RedisAddress  string `mapstructure:"redis_address"`
	RedisUsername string `mapstructure:"redis_username"`
	RedisPassword string `mapstructure:"redis_password"`
}

type manager struct {
	redisPool *redis.Pool
}

// New returns an implementation of a resource info cache that stores the objects in a redis cluster
func New(m map[string]interface{}) (cache.ResourceInfoCache, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, errors.Wrap(err, "error decoding conf")
	}

	if c.RedisAddress == "" {
		c.RedisAddress = "localhost:6379"
	}

	pool := &redis.Pool{
		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			var opts []redis.DialOption
			if c.RedisUsername != "" {
				opts = append(opts, redis.DialUsername(c.RedisUsername))
			}
			if c.RedisPassword != "" {
				opts = append(opts, redis.DialPassword(c.RedisPassword))
			}

			c, err := redis.Dial("tcp", c.RedisAddress, opts...)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &manager{
		redisPool: pool,
	}, nil
}

func (m *manager) Get(key string) (*provider.ResourceInfo, error) {
	infos, err := m.getVals([]string{key})
	if err != nil {
		return nil, err
	}
	return infos[0], nil
}

func (m *manager) GetKeys(keys []string) ([]*provider.ResourceInfo, error) {
	return m.getVals(keys)
}

func (m *manager) Set(key string, info *provider.ResourceInfo) error {
	return m.setVal(key, info, -1)
}

func (m *manager) SetWithExpire(key string, info *provider.ResourceInfo, expiration time.Duration) error {
	return m.setVal(key, info, int(expiration.Seconds()))
}

func (m *manager) setVal(key string, info *provider.ResourceInfo, expiration int) error {
	conn := m.redisPool.Get()
	defer conn.Close()
	if conn != nil {
		encodedInfo, err := json.Marshal(&info)
		if err != nil {
			return err
		}

		args := []interface{}{key, encodedInfo}
		if expiration != -1 {
			args = append(args, "EX", expiration)
		}

		if _, err := conn.Do("SET", args); err != nil {
			return err
		}
		return nil
	}
	return errors.New("cache: unable to get connection from redis pool")
}

func (m *manager) getVals(keys []string) ([]*provider.ResourceInfo, error) {
	conn := m.redisPool.Get()
	defer conn.Close()

	if conn != nil {
		vals, err := redis.Strings(conn.Do("MGET", keys))
		if err != nil {
			return nil, err
		}

		infos := make([]*provider.ResourceInfo, len(keys))
		for i, v := range vals {
			if v != "" {
				if err = json.Unmarshal([]byte(v), &infos[i]); err != nil {
					infos[i] = nil
				}
			}
		}
		return infos, nil
	}
	return nil, errors.New("cache: unable to get connection from redis pool")
}
