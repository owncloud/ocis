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

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/gomodule/redigo/redis"
)

const (
	userPrefix       = "user:"
	usernamePrefix   = "username:"
	userIDPrefix     = "userid:"
	namePrefix       = "name:"
	mailPrefix       = "mail:"
	uidPrefix        = "uid:"
	userGroupsPrefix = "groups:"
)

func initRedisPool(address, username, password string) *redis.Pool {
	return &redis.Pool{

		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			var opts []redis.DialOption
			if username != "" {
				opts = append(opts, redis.DialUsername(username))
			}
			if password != "" {
				opts = append(opts, redis.DialPassword(password))
			}

			c, err := redis.Dial("tcp", address, opts...)
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
}

func (m *manager) setVal(key, val string, expiration int) error {
	conn := m.redisPool.Get()
	defer conn.Close()
	if conn != nil {
		args := []interface{}{key, val}
		if expiration != -1 {
			args = append(args, "EX", expiration)
		}
		if _, err := conn.Do("SET", args...); err != nil {
			return err
		}
		return nil
	}
	return errors.New("rest: unable to get connection from redis pool")
}

func (m *manager) getVal(key string) (string, error) {
	conn := m.redisPool.Get()
	defer conn.Close()
	if conn != nil {
		val, err := redis.String(conn.Do("GET", key))
		if err != nil {
			return "", err
		}
		return val, nil
	}
	return "", errors.New("rest: unable to get connection from redis pool")
}

func (m *manager) findCachedUsers(query string) ([]*userpb.User, error) {
	conn := m.redisPool.Get()
	defer conn.Close()
	if conn != nil {
		query = fmt.Sprintf("%s*%s*", userPrefix, strings.ReplaceAll(strings.ToLower(query), " ", "_"))
		keys, err := redis.Strings(conn.Do("KEYS", query))
		if err != nil {
			return nil, err
		}
		var args []interface{}
		for _, k := range keys {
			args = append(args, k)
		}

		// Fetch the users for all these keys
		userStrings, err := redis.Strings(conn.Do("MGET", args...))
		if err != nil {
			return nil, err
		}
		userMap := make(map[string]*userpb.User)
		for _, user := range userStrings {
			u := userpb.User{}
			if err = json.Unmarshal([]byte(user), &u); err == nil {
				userMap[u.Id.OpaqueId] = &u
			}
		}

		var users []*userpb.User
		for _, u := range userMap {
			users = append(users, u)
		}

		return users, nil
	}

	return nil, errors.New("rest: unable to get connection from redis pool")
}

func (m *manager) fetchCachedUserDetails(uid *userpb.UserId) (*userpb.User, error) {
	user, err := m.getVal(userPrefix + usernamePrefix + strings.ToLower(uid.OpaqueId))
	if err != nil {
		return nil, err
	}

	u := userpb.User{}
	if err = json.Unmarshal([]byte(user), &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *manager) cacheUserDetails(u *userpb.User) error {
	encodedUser, err := json.Marshal(&u)
	if err != nil {
		return err
	}
	if err = m.setVal(userPrefix+usernamePrefix+strings.ToLower(u.Id.OpaqueId), string(encodedUser), -1); err != nil {
		return err
	}
	if err = m.setVal(userPrefix+userIDPrefix+strings.ToLower(u.Id.OpaqueId), string(encodedUser), -1); err != nil {
		return err
	}

	if u.Mail != "" {
		if err = m.setVal(userPrefix+mailPrefix+strings.ToLower(u.Mail), string(encodedUser), -1); err != nil {
			return err
		}
	}
	if u.DisplayName != "" {
		if err = m.setVal(userPrefix+namePrefix+u.Id.OpaqueId+"_"+strings.ReplaceAll(strings.ToLower(u.DisplayName), " ", "_"), string(encodedUser), -1); err != nil {
			return err
		}
	}
	if u.UidNumber != 0 {
		if err = m.setVal(userPrefix+uidPrefix+strconv.FormatInt(u.UidNumber, 10), string(encodedUser), -1); err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) fetchCachedUserByParam(field, claim string) (*userpb.User, error) {
	user, err := m.getVal(userPrefix + field + ":" + strings.ToLower(claim))
	if err != nil {
		return nil, err
	}

	u := userpb.User{}
	if err = json.Unmarshal([]byte(user), &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *manager) fetchCachedUserGroups(uid *userpb.UserId) ([]string, error) {
	groups, err := m.getVal(userPrefix + userGroupsPrefix + strings.ToLower(uid.OpaqueId))
	if err != nil {
		return nil, err
	}
	g := []string{}
	if err = json.Unmarshal([]byte(groups), &g); err != nil {
		return nil, err
	}
	return g, nil
}

func (m *manager) cacheUserGroups(uid *userpb.UserId, groups []string) error {
	g, err := json.Marshal(&groups)
	if err != nil {
		return err
	}
	return m.setVal(userPrefix+userGroupsPrefix+strings.ToLower(uid.OpaqueId), string(g), m.conf.UserGroupsCacheExpiration*60)
}
