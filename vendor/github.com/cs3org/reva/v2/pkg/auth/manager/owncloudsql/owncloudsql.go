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

package owncloudsql

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/auth"
	"github.com/cs3org/reva/v2/pkg/auth/manager/owncloudsql/accounts"
	"github.com/cs3org/reva/v2/pkg/auth/manager/registry"
	"github.com/cs3org/reva/v2/pkg/auth/scope"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	// Provides mysql drivers
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	registry.Register("owncloudsql", NewMysql)
}

type manager struct {
	c  *config
	db *accounts.Accounts
}

type config struct {
	DbUsername       string `mapstructure:"dbusername"`
	DbPassword       string `mapstructure:"dbpassword"`
	DbHost           string `mapstructure:"dbhost"`
	DbPort           int    `mapstructure:"dbport"`
	DbName           string `mapstructure:"dbname"`
	Idp              string `mapstructure:"idp"`
	Nobody           int64  `mapstructure:"nobody"`
	LegacySalt       string `mapstructure:"legacy_salt"`
	JoinUsername     bool   `mapstructure:"join_username"`
	JoinOwnCloudUUID bool   `mapstructure:"join_ownclouduuid"`
}

// NewMysql returns a new auth manager connection to an owncloud mysql database
func NewMysql(m map[string]interface{}) (auth.Manager, error) {
	mgr := &manager{}
	err := mgr.Configure(m)
	if err != nil {
		err = errors.Wrap(err, "error creating a new auth manager")
		return nil, err
	}

	mgr.db, err = accounts.NewMysql(
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mgr.c.DbUsername, mgr.c.DbPassword, mgr.c.DbHost, mgr.c.DbPort, mgr.c.DbName),
		mgr.c.JoinUsername,
		mgr.c.JoinOwnCloudUUID,
		false,
	)
	if err != nil {
		return nil, err
	}

	return mgr, nil
}

func (m *manager) Configure(ml map[string]interface{}) error {
	c, err := parseConfig(ml)
	if err != nil {
		return err
	}

	if c.Nobody == 0 {
		c.Nobody = 99
	}

	m.c = c
	return nil
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, &c); err != nil {
		return nil, err
	}
	return c, nil
}

func (m *manager) Authenticate(ctx context.Context, login, clientSecret string) (*user.User, map[string]*authpb.Scope, error) {
	log := appctx.GetLogger(ctx)

	// 1. find user by login

	account, err := m.db.GetAccountByLogin(ctx, login)
	if err != nil {
		return nil, nil, errtypes.NotFound(login)
	}
	// 2. verify the user password
	if !m.verify(clientSecret, account.PasswordHash) {
		return nil, nil, errtypes.InvalidCredentials(login)
	}

	userID := &user.UserId{
		Idp:      m.c.Idp,
		OpaqueId: account.OwnCloudUUID.String,
		Type:     user.UserType_USER_TYPE_PRIMARY, // TODO: assign the appropriate user type for guest accounts
	}

	u := &user.User{
		Id: userID,
		// TODO add more claims from the StandardClaims, eg EmailVerified and lastlogin
		Username:    account.Username.String,
		Mail:        account.Email.String,
		DisplayName: account.DisplayName.String,
		//UidNumber:   uidNumber,
		//GidNumber:   gidNumber,
	}

	if u.Groups, err = m.db.GetAccountGroups(ctx, account.UserID); err != nil {
		return nil, nil, err
	}

	var scopes map[string]*authpb.Scope
	if userID != nil && (userID.Type == user.UserType_USER_TYPE_LIGHTWEIGHT || userID.Type == user.UserType_USER_TYPE_FEDERATED) {
		scopes, err = scope.AddLightweightAccountScope(authpb.Role_ROLE_OWNER, nil)
		if err != nil {
			return nil, nil, err
		}
	} else {
		scopes, err = scope.AddOwnerScope(nil)
		if err != nil {
			return nil, nil, err
		}
	}
	// do not log password hash
	account.PasswordHash = "***redacted***"
	log.Debug().Interface("account", account).Interface("user", u).Msg("authenticated user")

	return u, scopes, nil
}

func (m *manager) verify(password, hash string) bool {
	splitHash := strings.SplitN(hash, "|", 2)
	switch len(splitHash) {
	case 2:
		if splitHash[0] == "1" {
			return m.verifyHashV1(password, splitHash[1])
		}
	case 1:
		return m.legacyHashVerify(password, hash)
	}
	return false
}

func (m *manager) legacyHashVerify(password, hash string) bool {
	// TODO rehash $newHash = $this->hash($message);
	switch len(hash) {
	case 60: // legacy PHPass hash
		return nil == bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+m.c.LegacySalt))
	case 40: // legacy sha1 hash
		h := sha1.Sum([]byte(password))
		return hmac.Equal([]byte(hash), []byte(hex.EncodeToString(h[:])))
	}
	return false
}
func (m *manager) verifyHashV1(password, hash string) bool {
	// TODO implement password_needs_rehash
	return nil == bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
