// Copyright 2018-2023 CERN
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
	"time"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	invitepb "github.com/cs3org/go-cs3apis/cs3/ocm/invite/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	conversions "github.com/cs3org/reva/v2/pkg/cbox/utils"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/ocm/invite"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils/cfg"
	"github.com/go-sql-driver/mysql"

	"github.com/cs3org/reva/v2/pkg/ocm/invite/repository/registry"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	"github.com/pkg/errors"
)

// This module implement the invite.Repository interface as a mysql driver.
//
// The OCM Invitation tokens are saved in the table:
//     ocm_tokens(*token*, initiator, expiration, description)
//
// The OCM remote user are saved in the table:
//     ocm_remote_users(*initiator*, *opaque_user_id*, *idp*, email, display_name)

func init() {
	registry.Register("sql", New)
}

type mgr struct {
	c      *config
	db     *sql.DB
	client gatewayv1beta1.GatewayAPIClient
}

type config struct {
	DBUsername string `mapstructure:"db_username"`
	DBPassword string `mapstructure:"db_password"`
	DBAddress  string `mapstructure:"db_address"`
	DBName     string `mapstructure:"db_name"`
	GatewaySvc string `mapstructure:"gatewaysvc"`
}

func (c *config) ApplyDefaults() {
	c.GatewaySvc = sharedconf.GetGatewaySVC(c.GatewaySvc)
}

// New creates a sql repository for ocm tokens and users.
func New(m map[string]interface{}) (invite.Repository, error) {
	var c config
	if err := cfg.Decode(m, &c); err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", c.DBUsername, c.DBPassword, c.DBAddress, c.DBName))
	if err != nil {
		return nil, errors.Wrap(err, "sql: error opening connection to mysql database")
	}

	gw, err := pool.GetGatewayServiceClient(c.GatewaySvc)
	if err != nil {
		return nil, err
	}

	mgr := mgr{
		c:      &c,
		db:     db,
		client: gw,
	}
	return &mgr, nil
}

// AddToken stores the token in the repository.
func (m *mgr) AddToken(ctx context.Context, token *invitepb.InviteToken) error {
	query := "INSERT INTO ocm_tokens SET token=?,initiator=?,expiration=?,description=?"
	_, err := m.db.ExecContext(ctx, query, token.Token, conversions.FormatUserID(token.UserId), timestampToTime(token.Expiration), token.Description)
	return err
}

func timestampToTime(t *types.Timestamp) time.Time {
	return time.Unix(int64(t.Seconds), int64(t.Nanos))
}

type dbToken struct {
	Token       string
	Initiator   string
	Expiration  time.Time
	Description string
}

// GetToken gets the token from the repository.
func (m *mgr) GetToken(ctx context.Context, token string) (*invitepb.InviteToken, error) {
	query := "SELECT token, initiator, expiration, description FROM ocm_tokens where token=?"

	var tkn dbToken
	if err := m.db.QueryRowContext(ctx, query, token).Scan(&tkn.Token, &tkn.Initiator, &tkn.Expiration, &tkn.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, invite.ErrTokenNotFound
		}
		return nil, err
	}
	return m.convertToInviteToken(ctx, tkn)
}

func (m *mgr) convertToInviteToken(ctx context.Context, tkn dbToken) (*invitepb.InviteToken, error) {
	user, err := conversions.ExtractUserID(ctx, m.client, tkn.Initiator)
	if err != nil {
		return nil, err
	}
	return &invitepb.InviteToken{
		Token:  tkn.Token,
		UserId: user,
		Expiration: &types.Timestamp{
			Seconds: uint64(tkn.Expiration.Unix()),
		},
		Description: tkn.Description,
	}, nil
}

func (m *mgr) ListTokens(ctx context.Context, initiator *userpb.UserId) ([]*invitepb.InviteToken, error) {
	query := "SELECT token, initiator, expiration, description FROM ocm_tokens WHERE initiator=? AND expiration > NOW()"

	tokens := []*invitepb.InviteToken{}
	rows, err := m.db.QueryContext(ctx, query, conversions.FormatUserID(initiator))
	if err != nil {
		return nil, err
	}

	var tkn dbToken
	for rows.Next() {
		if err := rows.Scan(&tkn.Token, &tkn.Initiator, &tkn.Expiration, &tkn.Description); err != nil {
			continue
		}
		token, err := m.convertToInviteToken(ctx, tkn)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// AddRemoteUser stores the remote user.
func (m *mgr) AddRemoteUser(ctx context.Context, initiator *userpb.UserId, remoteUser *userpb.User) error {
	query := "INSERT INTO ocm_remote_users SET initiator=?, opaque_user_id=?, idp=?, email=?, display_name=?"
	if _, err := m.db.ExecContext(ctx, query, conversions.FormatUserID(initiator), conversions.FormatUserID(remoteUser.Id), remoteUser.Id.Idp, remoteUser.Mail, remoteUser.DisplayName); err != nil {
		// check if the user already exist in the db
		// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_dup_entry
		var e *mysql.MySQLError
		if errors.As(err, &e) && e.Number == 1062 {
			return invite.ErrUserAlreadyAccepted
		}
		return err
	}
	return nil
}

type dbOCMUser struct {
	OpaqueUserID string
	Idp          string
	Email        string
	DisplayName  string
}

// GetRemoteUser retrieves details about a remote user who has accepted an invite to share.
func (m *mgr) GetRemoteUser(ctx context.Context, initiator *userpb.UserId, remoteUserID *userpb.UserId) (*userpb.User, error) {
	query := "SELECT opaque_user_id, idp, email, display_name FROM ocm_remote_users WHERE initiator=? AND opaque_user_id=? AND idp=?"

	var user dbOCMUser
	if err := m.db.QueryRowContext(ctx, query, conversions.FormatUserID(initiator), conversions.FormatUserID(remoteUserID), remoteUserID.Idp).
		Scan(&user.OpaqueUserID, &user.Idp, &user.Email, &user.DisplayName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errtypes.NotFound(remoteUserID.OpaqueId)
		}
		return nil, err
	}
	return user.toCS3User(), nil
}

func (u *dbOCMUser) toCS3User() *userpb.User {
	return &userpb.User{
		Id: &userpb.UserId{
			Idp:      u.Idp,
			OpaqueId: u.OpaqueUserID,
			Type:     userpb.UserType_USER_TYPE_FEDERATED,
		},
		Mail:        u.Email,
		DisplayName: u.DisplayName,
	}
}

// FindRemoteUsers finds remote users who have accepted invites based on their attributes.
func (m *mgr) FindRemoteUsers(ctx context.Context, initiator *userpb.UserId, attr string) ([]*userpb.User, error) {
	// TODO: (gdelmont) this query can get really slow in case the number of rows is too high.
	// For the time being this is not expected, but if in future this happens, consider to add
	// a fulltext index.
	query := "SELECT opaque_user_id, idp, email, display_name FROM ocm_remote_users WHERE initiator=? AND (opaque_user_id LIKE ? OR idp LIKE ? OR email LIKE ? OR display_name LIKE ?)"
	s := "%" + attr + "%"
	params := []any{conversions.FormatUserID(initiator), s, s, s, s}

	rows, err := m.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	var u dbOCMUser
	var users []*userpb.User
	for rows.Next() {
		if err := rows.Scan(&u.OpaqueUserID, &u.Idp, &u.Email, &u.DisplayName); err != nil {
			continue
		}
		users = append(users, u.toCS3User())
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *mgr) DeleteRemoteUser(ctx context.Context, initiator *userpb.UserId, remoteUser *userpb.UserId) error {
	query := "DELETE FROM ocm_remote_users WHERE initiator=? AND opaque_user_id=? AND idp=?"
	_, err := m.db.ExecContext(ctx, query, conversions.FormatUserID(initiator), conversions.FormatUserID(remoteUser), remoteUser.Idp)
	return err
}
