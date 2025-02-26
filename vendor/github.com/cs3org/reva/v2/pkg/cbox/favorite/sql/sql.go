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

package cbox

import (
	"context"
	"database/sql"
	"fmt"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/cbox/utils"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/storage/favorite"
	"github.com/cs3org/reva/v2/pkg/storage/favorite/registry"
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

// New returns an instance of the cbox sql favorites manager.
func New(m map[string]interface{}) (favorite.Manager, error) {
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

func (m *mgr) ListFavorites(ctx context.Context, userID *user.UserId) ([]*provider.ResourceId, error) {
	user := ctxpkg.ContextMustGetUser(ctx)
	infos := []*provider.ResourceId{}
	query := `SELECT fileid_prefix, fileid FROM cbox_metadata WHERE uid=? AND tag_key="fav"`
	rows, err := m.db.Query(query, user.Id.OpaqueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var info provider.ResourceId
		if err := rows.Scan(&info.SpaceId, &info.OpaqueId); err != nil {
			return nil, err
		}
		infos = append(infos, &info)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return infos, nil
}

func (m *mgr) SetFavorite(ctx context.Context, userID *user.UserId, resourceInfo *provider.ResourceInfo) error {
	user := ctxpkg.ContextMustGetUser(ctx)
	spaceID := resourceInfo.Id.SpaceId

	// The primary key is just the ID in the table, it should ideally be (uid, fileid_prefix, fileid, tag_key)
	// For the time being, just check if the favorite already exists. If it does, return early
	var id int
	query := `SELECT id FROM cbox_metadata WHERE uid=? AND fileid_prefix=? AND fileid=? AND tag_key="fav"`
	if err := m.db.QueryRow(query, user.Id.OpaqueId, spaceID, resourceInfo.Id.OpaqueId).Scan(&id); err == nil {
		// Favorite is already set, return
		return nil
	}

	query = `INSERT INTO cbox_metadata SET item_type=?, uid=?, fileid_prefix=?, fileid=?, tag_key="fav"`
	vals := []interface{}{utils.ResourceTypeToItemInt(resourceInfo.Type), user.Id.OpaqueId, spaceID, resourceInfo.Id.OpaqueId}
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(vals...); err != nil {
		return err
	}
	return nil
}

func (m *mgr) UnsetFavorite(ctx context.Context, userID *user.UserId, resourceInfo *provider.ResourceInfo) error {
	user := ctxpkg.ContextMustGetUser(ctx)
	spaceID := resourceInfo.Id.SpaceId

	stmt, err := m.db.Prepare(`DELETE FROM cbox_metadata WHERE uid=? AND fileid_prefix=? AND fileid=? AND tag_key="fav"`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(user.Id.OpaqueId, spaceID, resourceInfo.Id.OpaqueId)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
