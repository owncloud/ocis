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

	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/gomodule/redigo/redis"
)

const (
	groupPrefix           = "group:"
	idPrefix              = "id:"
	namePrefix            = "name:"
	gidPrefix             = "gid:"
	groupMembersPrefix    = "members:"
	groupInternalIDPrefix = "internal:"
)

func initRedisPool(address, username, password string) *redis.Pool {
	return &redis.Pool{

		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			var c redis.Conn
			var err error
			switch {
			case username != "":
				c, err = redis.Dial("tcp", address,
					redis.DialUsername(username),
					redis.DialPassword(password),
				)
			case password != "":
				c, err = redis.Dial("tcp", address,
					redis.DialPassword(password),
				)
			default:
				c, err = redis.Dial("tcp", address)
			}

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

func (m *manager) fetchCachedInternalID(gid *grouppb.GroupId) (string, error) {
	return m.getVal(groupPrefix + groupInternalIDPrefix + gid.OpaqueId)
}

func (m *manager) cacheInternalID(gid *grouppb.GroupId, internalID string) error {
	return m.setVal(groupPrefix+groupInternalIDPrefix+gid.OpaqueId, internalID, -1)
}

func (m *manager) findCachedGroups(query string) ([]*grouppb.Group, error) {
	conn := m.redisPool.Get()
	defer conn.Close()
	if conn != nil {
		query = fmt.Sprintf("%s*%s*", groupPrefix, strings.ReplaceAll(strings.ToLower(query), " ", "_"))
		keys, err := redis.Strings(conn.Do("KEYS", query))
		if err != nil {
			return nil, err
		}
		var args []interface{}
		for _, k := range keys {
			args = append(args, k)
		}

		// Fetch the groups for all these keys
		groupStrings, err := redis.Strings(conn.Do("MGET", args...))
		if err != nil {
			return nil, err
		}
		groupMap := make(map[string]*grouppb.Group)
		for _, group := range groupStrings {
			g := grouppb.Group{}
			if err = json.Unmarshal([]byte(group), &g); err == nil {
				groupMap[g.Id.OpaqueId] = &g
			}
		}

		var groups []*grouppb.Group
		for _, g := range groupMap {
			groups = append(groups, g)
		}

		return groups, nil
	}

	return nil, errors.New("rest: unable to get connection from redis pool")
}

func (m *manager) fetchCachedGroupDetails(gid *grouppb.GroupId) (*grouppb.Group, error) {
	group, err := m.getVal(groupPrefix + idPrefix + gid.OpaqueId)
	if err != nil {
		return nil, err
	}

	g := grouppb.Group{}
	if err = json.Unmarshal([]byte(group), &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (m *manager) cacheGroupDetails(g *grouppb.Group) error {
	encodedGroup, err := json.Marshal(&g)
	if err != nil {
		return err
	}
	if err = m.setVal(groupPrefix+idPrefix+strings.ToLower(g.Id.OpaqueId), string(encodedGroup), -1); err != nil {
		return err
	}

	if g.GidNumber != 0 {
		if err = m.setVal(groupPrefix+gidPrefix+strconv.FormatInt(g.GidNumber, 10), g.Id.OpaqueId, -1); err != nil {
			return err
		}
	}
	if g.DisplayName != "" {
		if err = m.setVal(groupPrefix+namePrefix+g.Id.OpaqueId+"_"+strings.ToLower(g.DisplayName), g.Id.OpaqueId, -1); err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) fetchCachedGroupByParam(field, claim string) (*grouppb.Group, error) {
	group, err := m.getVal(groupPrefix + field + ":" + strings.ToLower(claim))
	if err != nil {
		return nil, err
	}

	g := grouppb.Group{}
	if err = json.Unmarshal([]byte(group), &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (m *manager) fetchCachedGroupMembers(gid *grouppb.GroupId) ([]*userpb.UserId, error) {
	members, err := m.getVal(groupPrefix + groupMembersPrefix + strings.ToLower(gid.OpaqueId))
	if err != nil {
		return nil, err
	}
	u := []*userpb.UserId{}
	if err = json.Unmarshal([]byte(members), &u); err != nil {
		return nil, err
	}
	return u, nil
}

func (m *manager) cacheGroupMembers(gid *grouppb.GroupId, members []*userpb.UserId) error {
	u, err := json.Marshal(&members)
	if err != nil {
		return err
	}
	return m.setVal(groupPrefix+groupMembersPrefix+strings.ToLower(gid.OpaqueId), string(u), m.conf.GroupMembersCacheExpiration*60)
}
