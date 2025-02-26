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

package accounts

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/pkg/errors"
)

// Accounts represents oc10-style Accounts
type Accounts struct {
	driver                                     string
	db                                         *sql.DB
	joinUsername, joinUUID, enableMedialSearch bool
	selectSQL                                  string
}

// NewMysql returns a new Cache instance connecting to a MySQL database
func NewMysql(dsn string, joinUsername, joinUUID, enableMedialSearch bool) (*Accounts, error) {
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to the database")
	}

	// FIXME make configurable
	sqldb.SetConnMaxLifetime(time.Minute * 3)
	sqldb.SetConnMaxIdleTime(time.Second * 30)
	sqldb.SetMaxOpenConns(100)
	sqldb.SetMaxIdleConns(10)

	err = sqldb.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to the database")
	}

	return New("mysql", sqldb, joinUsername, joinUUID, enableMedialSearch)
}

// New returns a new Cache instance connecting to the given sql.DB
func New(driver string, sqldb *sql.DB, joinUsername, joinUUID, enableMedialSearch bool) (*Accounts, error) {

	sel := "SELECT id, email, user_id, display_name, quota, last_login, backend, home, state"
	from := `
		FROM oc_accounts a
		`
	if joinUsername {
		sel += ", p.configvalue AS username"
		from += `LEFT JOIN oc_preferences p
						ON a.user_id=p.userid
						AND p.appid='core'
						AND p.configkey='username'`
	} else {
		// fallback to user_id as username
		sel += ", user_id AS username"
	}
	if joinUUID {
		sel += ", p2.configvalue AS ownclouduuid"
		from += `LEFT JOIN oc_preferences p2
						ON a.user_id=p2.userid
						AND p2.appid='core'
						AND p2.configkey='ownclouduuid'`
	} else {
		// fallback to user_id as ownclouduuid
		sel += ", user_id AS ownclouduuid"
	}

	return &Accounts{
		driver:             driver,
		db:                 sqldb,
		joinUsername:       joinUsername,
		joinUUID:           joinUUID,
		enableMedialSearch: enableMedialSearch,
		selectSQL:          sel + from,
	}, nil
}

// Account stores information about accounts.
type Account struct {
	ID           uint64
	Email        sql.NullString
	UserID       string
	DisplayName  sql.NullString
	Quota        sql.NullString
	LastLogin    int
	Backend      string
	Home         string
	State        int8
	Username     sql.NullString // optional comes from the oc_preferences
	OwnCloudUUID sql.NullString // optional comes from the oc_preferences
}

func (as *Accounts) rowToAccount(ctx context.Context, row Scannable) (*Account, error) {
	a := Account{}
	if err := row.Scan(&a.ID, &a.Email, &a.UserID, &a.DisplayName, &a.Quota, &a.LastLogin, &a.Backend, &a.Home, &a.State, &a.Username, &a.OwnCloudUUID); err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("could not scan row, skipping")
		return nil, err
	}

	return &a, nil
}

// Scannable describes the interface providing a Scan method
type Scannable interface {
	Scan(...interface{}) error
}

// GetAccountByClaim fetches an account by mail, username or userid
func (as *Accounts) GetAccountByClaim(ctx context.Context, claim, value string) (*Account, error) {
	// TODO align supported claims with rest driver and the others, maybe refactor into common mapping
	var row *sql.Row
	var where string
	switch claim {
	case "mail":
		where = "WHERE a.email=?"
	// case "uid":
	//	claim = m.c.Schema.UIDNumber
	// case "gid":
	//	claim = m.c.Schema.GIDNumber
	case "username":
		if as.joinUsername {
			where = "WHERE p.configvalue=?"
		} else {
			// use user_id as username
			where = "WHERE a.user_id=?"
		}
	case "userid":
		if as.joinUUID {
			where = "WHERE p2.configvalue=?"
		} else {
			// use user_id as uuid
			where = "WHERE a.user_id=?"
		}
	default:
		return nil, errors.New("owncloudsql: invalid field " + claim)
	}

	row = as.db.QueryRowContext(ctx, as.selectSQL+where, value)

	return as.rowToAccount(ctx, row)
}

func sanitizeWildcards(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "%", `\%`), "_", `\_`)
}

// FindAccounts searches userid, displayname and email using the given query. The Wildcard caracters % and _ are escaped.
func (as *Accounts) FindAccounts(ctx context.Context, query string) ([]Account, error) {
	if as.enableMedialSearch {
		query = "%" + sanitizeWildcards(query) + "%"
	}
	// TODO join oc_account_terms
	where := "WHERE a.user_id LIKE ? OR a.display_name LIKE ? OR a.email LIKE ?"
	args := []interface{}{query, query, query}

	if as.joinUsername {
		where += " OR p.configvalue LIKE ?"
		args = append(args, query)
	}
	if as.joinUUID {
		where += " OR p2.configvalue LIKE ?"
		args = append(args, query)
	}

	rows, err := as.db.QueryContext(ctx, as.selectSQL+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []Account{}
	for rows.Next() {
		a := Account{}
		if err := rows.Scan(&a.ID, &a.Email, &a.UserID, &a.DisplayName, &a.Quota, &a.LastLogin, &a.Backend, &a.Home, &a.State, &a.Username, &a.OwnCloudUUID); err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Msg("could not scan row, skipping")
			continue
		}
		accounts = append(accounts, a)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// GetAccountGroups lasts the groups for an account
func (as *Accounts) GetAccountGroups(ctx context.Context, uid string) ([]string, error) {
	rows, err := as.db.QueryContext(ctx, "SELECT gid FROM oc_group_user WHERE uid=?", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []string{}
	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Msg("could not scan row, skipping")
			continue
		}
		groups = append(groups, group)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return groups, nil
}
