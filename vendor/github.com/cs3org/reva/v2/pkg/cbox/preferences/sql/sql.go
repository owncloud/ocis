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

package sql

import (
	"context"
	"database/sql"
	"fmt"

	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/preferences"
	"github.com/cs3org/reva/v2/pkg/preferences/registry"
	"github.com/mitchellh/mapstructure"
)

func init() {
	registry.Register("sql", New)
}

type config struct {
	DbUsername string `mapstructure:"db_username"`
	DbPassword string `mapstructure:"db_password"`
	DbHost     string `mapstructure:"db_host"`
	DbPort     int    `mapstructure:"db_port"`
	DbName     string `mapstructure:"db_name"`
}

type mgr struct {
	c  *config
	db *sql.DB
}

// New returns an instance of the cbox sql preferences manager.
func New(m map[string]interface{}) (preferences.Manager, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.DbUsername, c.DbPassword, c.DbHost, c.DbPort, c.DbName))
	if err != nil {
		return nil, err
	}

	return &mgr{
		c:  c,
		db: db,
	}, nil
}

func (m *mgr) SetKey(ctx context.Context, key, namespace, value string) error {
	user, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		return errtypes.UserRequired("preferences: error getting user from ctx")
	}
	query := `INSERT INTO oc_preferences(userid, appid, configkey, configvalue) values(?, ?, ?, ?) ON DUPLICATE KEY UPDATE configvalue = ?`
	params := []interface{}{user.Id.OpaqueId, namespace, key, value, value}
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(params...); err != nil {
		return err
	}
	return nil
}

func (m *mgr) GetKey(ctx context.Context, key, namespace string) (string, error) {
	user, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		return "", errtypes.UserRequired("preferences: error getting user from ctx")
	}
	query := `SELECT configvalue FROM oc_preferences WHERE userid=? AND appid=? AND configkey=?`
	var val string
	if err := m.db.QueryRow(query, user.Id.OpaqueId, namespace, key).Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return "", errtypes.NotFound(namespace + ":" + key)
		}
		return "", err
	}
	return val, nil
}
