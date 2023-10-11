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
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	conversions "github.com/cs3org/reva/v2/pkg/cbox/utils"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/publicshare"
	"github.com/cs3org/reva/v2/pkg/publicshare/manager/registry"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const publicShareType = 3

func init() {
	registry.Register("sql", New)
}

type config struct {
	SharePasswordHashCost      int    `mapstructure:"password_hash_cost"`
	JanitorRunInterval         int    `mapstructure:"janitor_run_interval"`
	EnableExpiredSharesCleanup bool   `mapstructure:"enable_expired_shares_cleanup"`
	DbUsername                 string `mapstructure:"db_username"`
	DbPassword                 string `mapstructure:"db_password"`
	DbHost                     string `mapstructure:"db_host"`
	DbPort                     int    `mapstructure:"db_port"`
	DbName                     string `mapstructure:"db_name"`
	GatewaySvc                 string `mapstructure:"gatewaysvc"`
}

type manager struct {
	c      *config
	db     *sql.DB
	client gatewayv1beta1.GatewayAPIClient
}

func (c *config) init() {
	if c.SharePasswordHashCost == 0 {
		c.SharePasswordHashCost = 11
	}
	if c.JanitorRunInterval == 0 {
		c.JanitorRunInterval = 3600
	}
}

func (m *manager) startJanitorRun() {
	if !m.c.EnableExpiredSharesCleanup {
		return
	}

	ticker := time.NewTicker(time.Duration(m.c.JanitorRunInterval) * time.Second)
	work := make(chan os.Signal, 1)
	signal.Notify(work, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-work:
			return
		case <-ticker.C:
			_ = m.cleanupExpiredShares()
		}
	}
}

// New returns a new public share manager.
func New(m map[string]interface{}) (publicshare.Manager, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}
	c.init()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.DbUsername, c.DbPassword, c.DbHost, c.DbPort, c.DbName))
	if err != nil {
		return nil, err
	}

	gw, err := pool.GetGatewayServiceClient(c.GatewaySvc)
	if err != nil {
		return nil, err
	}

	mgr := manager{
		c:      c,
		db:     db,
		client: gw,
	}
	go mgr.startJanitorRun()

	return &mgr, nil
}

func (m *manager) CreatePublicShare(ctx context.Context, u *user.User, rInfo *provider.ResourceInfo, g *link.Grant) (*link.PublicShare, error) {

	tkn := utils.RandString(15)
	now := time.Now().Unix()

	displayName, ok := rInfo.ArbitraryMetadata.Metadata["name"]
	if !ok {
		displayName = tkn
	}
	createdAt := &typespb.Timestamp{
		Seconds: uint64(now),
	}

	creator := conversions.FormatUserID(u.Id)
	owner := conversions.FormatUserID(rInfo.Owner)
	permissions := conversions.SharePermToInt(g.Permissions.Permissions)
	itemType := conversions.ResourceTypeToItem(rInfo.Type)
	prefix := rInfo.Id.SpaceId
	itemSource := rInfo.Id.OpaqueId
	fileSource, err := strconv.ParseUint(itemSource, 10, 64)
	if err != nil {
		// it can be the case that the item source may be a character string
		// we leave fileSource blank in that case
		fileSource = 0
	}

	query := "insert into oc_share set share_type=?,uid_owner=?,uid_initiator=?,item_type=?,fileid_prefix=?,item_source=?,file_source=?,permissions=?,stime=?,token=?,share_name=?"
	params := []interface{}{publicShareType, owner, creator, itemType, prefix, itemSource, fileSource, permissions, now, tkn, displayName}

	var passwordProtected bool
	password := g.Password
	if password != "" {
		password, err = hashPassword(password, m.c.SharePasswordHashCost)
		if err != nil {
			return nil, errors.Wrap(err, "could not hash share password")
		}
		passwordProtected = true

		query += ",share_with=?"
		params = append(params, password)
	}

	if g.Expiration != nil && g.Expiration.Seconds != 0 {
		t := time.Unix(int64(g.Expiration.Seconds), 0)
		query += ",expiration=?"
		params = append(params, t)
	}

	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(params...)
	if err != nil {
		return nil, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &link.PublicShare{
		Id: &link.PublicShareId{
			OpaqueId: strconv.FormatInt(lastID, 10),
		},
		Owner:             rInfo.GetOwner(),
		Creator:           u.Id,
		ResourceId:        rInfo.Id,
		Token:             tkn,
		Permissions:       g.Permissions,
		Ctime:             createdAt,
		Mtime:             createdAt,
		PasswordProtected: passwordProtected,
		Expiration:        g.Expiration,
		DisplayName:       displayName,
	}, nil
}

func (m *manager) UpdatePublicShare(ctx context.Context, u *user.User, req *link.UpdatePublicShareRequest) (*link.PublicShare, error) {
	query := "update oc_share set "
	paramsMap := map[string]interface{}{}
	params := []interface{}{}

	now := time.Now().Unix()
	uid := conversions.FormatUserID(u.Id)

	switch req.GetUpdate().GetType() {
	case link.UpdatePublicShareRequest_Update_TYPE_DISPLAYNAME:
		paramsMap["share_name"] = req.Update.GetDisplayName()
	case link.UpdatePublicShareRequest_Update_TYPE_PERMISSIONS:
		paramsMap["permissions"] = conversions.SharePermToInt(req.Update.GetGrant().GetPermissions().Permissions)
	case link.UpdatePublicShareRequest_Update_TYPE_EXPIRATION:
		paramsMap["expiration"] = time.Unix(int64(req.Update.GetGrant().Expiration.Seconds), 0)
	case link.UpdatePublicShareRequest_Update_TYPE_PASSWORD:
		if req.Update.GetGrant().Password == "" {
			paramsMap["share_with"] = ""
		} else {
			h, err := hashPassword(req.Update.GetGrant().Password, m.c.SharePasswordHashCost)
			if err != nil {
				return nil, errors.Wrap(err, "could not hash share password")
			}
			paramsMap["share_with"] = h
		}
	default:
		return nil, fmt.Errorf("invalid update type: %v", req.GetUpdate().GetType())
	}

	for k, v := range paramsMap {
		query += k + "=?"
		params = append(params, v)
	}

	switch {
	case req.Ref.GetId() != nil:
		query += ",stime=? where id=? AND (uid_owner=? or uid_initiator=?)"
		params = append(params, now, req.Ref.GetId().OpaqueId, uid, uid)
	case req.Ref.GetToken() != "":
		query += ",stime=? where token=? AND (uid_owner=? or uid_initiator=?)"
		params = append(params, now, req.Ref.GetToken(), uid, uid)
	default:
		return nil, errtypes.NotFound(req.Ref.String())
	}

	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	if _, err = stmt.Exec(params...); err != nil {
		return nil, err
	}

	return m.GetPublicShare(ctx, u, req.Ref, false)
}

func (m *manager) getByToken(ctx context.Context, token string, u *user.User) (*link.PublicShare, string, error) {
	s := conversions.DBShare{Token: token}
	query := "select coalesce(uid_owner, '') as uid_owner, coalesce(uid_initiator, '') as uid_initiator, coalesce(share_with, '') as share_with, coalesce(fileid_prefix, '') as fileid_prefix, coalesce(item_source, '') as item_source, coalesce(item_type, '') as item_type, coalesce(expiration, '') as expiration, coalesce(share_name, '') as share_name, id, stime, permissions FROM oc_share WHERE (orphan = 0 or orphan IS NULL) AND share_type=? AND token=?"
	if err := m.db.QueryRow(query, publicShareType, token).Scan(&s.UIDOwner, &s.UIDInitiator, &s.ShareWith, &s.Prefix, &s.ItemSource, &s.ItemType, &s.Expiration, &s.ShareName, &s.ID, &s.STime, &s.Permissions); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errtypes.NotFound(token)
		}
		return nil, "", err
	}
	share, err := conversions.ConvertToCS3PublicShare(ctx, m.client, s)
	if err != nil {
		return nil, "", err
	}
	return share, s.ShareWith, nil
}

func (m *manager) getByID(ctx context.Context, id *link.PublicShareId, u *user.User) (*link.PublicShare, string, error) {
	uid := conversions.FormatUserID(u.Id)
	s := conversions.DBShare{ID: id.OpaqueId}
	query := "select coalesce(uid_owner, '') as uid_owner, coalesce(uid_initiator, '') as uid_initiator, coalesce(share_with, '') as share_with, coalesce(fileid_prefix, '') as fileid_prefix, coalesce(item_source, '') as item_source, coalesce(item_type, '') as item_type, coalesce(token,'') as token, coalesce(expiration, '') as expiration, coalesce(share_name, '') as share_name, stime, permissions FROM oc_share WHERE (orphan = 0 or orphan IS NULL) AND share_type=? AND id=? AND (uid_owner=? OR uid_initiator=?)"
	if err := m.db.QueryRow(query, publicShareType, id.OpaqueId, uid, uid).Scan(&s.UIDOwner, &s.UIDInitiator, &s.ShareWith, &s.Prefix, &s.ItemSource, &s.ItemType, &s.Token, &s.Expiration, &s.ShareName, &s.STime, &s.Permissions); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errtypes.NotFound(id.OpaqueId)
		}
		return nil, "", err
	}
	share, err := conversions.ConvertToCS3PublicShare(ctx, m.client, s)
	if err != nil {
		return nil, "", err
	}
	return share, s.ShareWith, nil
}

func (m *manager) GetPublicShare(ctx context.Context, u *user.User, ref *link.PublicShareReference, sign bool) (*link.PublicShare, error) {
	var s *link.PublicShare
	var pw string
	var err error
	switch {
	case ref.GetId() != nil:
		s, pw, err = m.getByID(ctx, ref.GetId(), u)
	case ref.GetToken() != "":
		s, pw, err = m.getByToken(ctx, ref.GetToken(), u)
	default:
		err = errtypes.NotFound(ref.String())
	}
	if err != nil {
		return nil, err
	}

	if expired(s) {
		if err := m.cleanupExpiredShares(); err != nil {
			return nil, err
		}
		return nil, errtypes.NotFound(ref.String())
	}

	if s.PasswordProtected && sign {
		if err := publicshare.AddSignature(s, pw); err != nil {
			return nil, err
		}

	}

	return s, nil
}

func (m *manager) ListPublicShares(ctx context.Context, u *user.User, filters []*link.ListPublicSharesRequest_Filter, sign bool) ([]*link.PublicShare, error) {
	uid := conversions.FormatUserID(u.Id)
	query := "select coalesce(uid_owner, '') as uid_owner, coalesce(uid_initiator, '') as uid_initiator, coalesce(share_with, '') as share_with, coalesce(fileid_prefix, '') as fileid_prefix, coalesce(item_source, '') as item_source, coalesce(item_type, '') as item_type, coalesce(token,'') as token, coalesce(expiration, '') as expiration, coalesce(share_name, '') as share_name, id, stime, permissions FROM oc_share WHERE (orphan = 0 or orphan IS NULL) AND (uid_owner=? or uid_initiator=?) AND (share_type=?)"
	var resourceFilters, ownerFilters, creatorFilters string
	var resourceParams, ownerParams, creatorParams []interface{}
	params := []interface{}{uid, uid, publicShareType}

	for _, f := range filters {
		switch f.Type {
		case link.ListPublicSharesRequest_Filter_TYPE_RESOURCE_ID:
			if len(resourceFilters) != 0 {
				resourceFilters += " OR "
			}
			resourceFilters += "(fileid_prefix=? AND item_source=?)"
			resourceParams = append(resourceParams, f.GetResourceId().SpaceId, f.GetResourceId().OpaqueId)
		case link.ListPublicSharesRequest_Filter_TYPE_OWNER:
			if len(ownerFilters) != 0 {
				ownerFilters += " OR "
			}
			ownerFilters += "(uid_owner=?)"
			ownerParams = append(ownerParams, conversions.FormatUserID(f.GetOwner()))
		case link.ListPublicSharesRequest_Filter_TYPE_CREATOR:
			if len(creatorFilters) != 0 {
				creatorFilters += " OR "
			}
			creatorFilters += "(uid_initiator=?)"
			creatorParams = append(creatorParams, conversions.FormatUserID(f.GetCreator()))
		}
	}
	if resourceFilters != "" {
		query = fmt.Sprintf("%s AND (%s)", query, resourceFilters)
		params = append(params, resourceParams...)
	}
	if ownerFilters != "" {
		query = fmt.Sprintf("%s AND (%s)", query, ownerFilters)
		params = append(params, ownerParams...)
	}
	if creatorFilters != "" {
		query = fmt.Sprintf("%s AND (%s)", query, creatorFilters)
		params = append(params, creatorParams...)
	}

	rows, err := m.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var s conversions.DBShare
	shares := []*link.PublicShare{}
	for rows.Next() {
		if err := rows.Scan(&s.UIDOwner, &s.UIDInitiator, &s.ShareWith, &s.Prefix, &s.ItemSource, &s.ItemType, &s.Token, &s.Expiration, &s.ShareName, &s.ID, &s.STime, &s.Permissions); err != nil {
			continue
		}
		cs3Share, err := conversions.ConvertToCS3PublicShare(ctx, m.client, s)
		if err != nil {
			return nil, err
		}
		if expired(cs3Share) {
			_ = m.cleanupExpiredShares()
		} else {
			if cs3Share.PasswordProtected && sign {
				if err := publicshare.AddSignature(cs3Share, s.ShareWith); err != nil {
					return nil, err
				}
			}
			shares = append(shares, cs3Share)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shares, nil
}

func (m *manager) RevokePublicShare(ctx context.Context, u *user.User, ref *link.PublicShareReference) error {
	uid := conversions.FormatUserID(u.Id)
	query := "delete from oc_share where "
	params := []interface{}{}

	switch {
	case ref.GetId() != nil && ref.GetId().OpaqueId != "":
		query += "id=? AND (uid_owner=? or uid_initiator=?)"
		params = append(params, ref.GetId().OpaqueId, uid, uid)
	case ref.GetToken() != "":
		query += "token=? AND (uid_owner=? or uid_initiator=?)"
		params = append(params, ref.GetToken(), uid, uid)
	default:
		return errtypes.NotFound(ref.String())
	}

	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(params...)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt == 0 {
		return errtypes.NotFound(ref.String())
	}
	return nil
}

func (m *manager) GetPublicShareByToken(ctx context.Context, token string, auth *link.PublicShareAuthentication, sign bool) (*link.PublicShare, error) {
	s := conversions.DBShare{Token: token}
	query := "select coalesce(uid_owner, '') as uid_owner, coalesce(uid_initiator, '') as uid_initiator, coalesce(share_with, '') as share_with, coalesce(fileid_prefix, '') as fileid_prefix, coalesce(item_source, '') as item_source, coalesce(item_type, '') as item_type, coalesce(expiration, '') as expiration, coalesce(share_name, '') as share_name, id, stime, permissions FROM oc_share WHERE share_type=? AND token=?"
	if err := m.db.QueryRow(query, publicShareType, token).Scan(&s.UIDOwner, &s.UIDInitiator, &s.ShareWith, &s.Prefix, &s.ItemSource, &s.ItemType, &s.Expiration, &s.ShareName, &s.ID, &s.STime, &s.Permissions); err != nil {
		if err == sql.ErrNoRows {
			return nil, errtypes.NotFound(token)
		}
		return nil, err
	}
	cs3Share, err := conversions.ConvertToCS3PublicShare(ctx, m.client, s)
	if err != nil {
		return nil, err
	}
	if s.ShareWith != "" {
		if !authenticate(cs3Share, s.ShareWith, auth) {
			// if check := checkPasswordHash(auth.Password, s.ShareWith); !check {
			return nil, errtypes.InvalidCredentials(token)
		}

		if sign {
			if err := publicshare.AddSignature(cs3Share, s.ShareWith); err != nil {
				return nil, err
			}
		}
	}

	if expired(cs3Share) {
		if err := m.cleanupExpiredShares(); err != nil {
			return nil, err
		}
		return nil, errtypes.NotFound(token)
	}

	return cs3Share, nil
}

func (m *manager) cleanupExpiredShares() error {
	if !m.c.EnableExpiredSharesCleanup {
		return nil
	}

	query := "update oc_share set orphan = 1 where expiration IS NOT NULL AND expiration < ?"
	params := []interface{}{time.Now().Format("2006-01-02 03:04:05")}

	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(params...); err != nil {
		return err
	}
	return nil
}

func expired(s *link.PublicShare) bool {
	if s.Expiration != nil {
		if t := time.Unix(int64(s.Expiration.GetSeconds()), int64(s.Expiration.GetNanos())); t.Before(time.Now()) {
			return true
		}
	}
	return false
}

func hashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return "1|" + string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(strings.TrimPrefix(hash, "1|")), []byte(password))
	return err == nil
}

func authenticate(share *link.PublicShare, pw string, auth *link.PublicShareAuthentication) bool {
	switch {
	case auth.GetPassword() != "":
		return checkPasswordHash(auth.GetPassword(), pw)
	case auth.GetSignature() != nil:
		sig := auth.GetSignature()
		now := time.Now()
		expiration := time.Unix(int64(sig.GetSignatureExpiration().GetSeconds()), int64(sig.GetSignatureExpiration().GetNanos()))
		if now.After(expiration) {
			return false
		}
		s, err := publicshare.CreateSignature(share.Token, pw, expiration)
		if err != nil {
			// TODO(labkode): pass context to call to log err.
			// No we are blind
			return false
		}
		return sig.GetSignature() == s
	}
	return false
}
