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

package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/ocm/share"

	"github.com/cs3org/reva/v2/pkg/ocm/share/manager/registry"
	"github.com/cs3org/reva/v2/pkg/ocm/share/sender"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/genproto/protobuf/field_mask"
)

func init() {
	registry.Register("json", New)
}

// New returns a new authorizer object.
func New(m map[string]interface{}) (share.Manager, error) {
	c, err := parseConfig(m)
	if err != nil {
		err = errors.Wrap(err, "error creating a new manager")
		return nil, err
	}
	c.init()

	// load or create file
	model, err := loadOrCreate(c.File)
	if err != nil {
		err = errors.Wrap(err, "error loading the file containing the shares")
		return nil, err
	}

	mgr := &mgr{
		c:     c,
		model: model,
	}

	return mgr, nil
}

func loadOrCreate(file string) (*shareModel, error) {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		if err := os.WriteFile(file, []byte("{}"), 0700); err != nil {
			err = errors.Wrap(err, "error creating the file: "+file)
			return nil, err
		}
	}

	fd, err := os.OpenFile(file, os.O_CREATE, 0644)
	if err != nil {
		err = errors.Wrap(err, "error opening the file: "+file)
		return nil, err
	}
	defer fd.Close()

	data, err := io.ReadAll(fd)
	if err != nil {
		err = errors.Wrap(err, "error reading the data")
		return nil, err
	}

	m := &shareModel{}
	if err := json.Unmarshal(data, m); err != nil {
		err = errors.Wrap(err, "error decoding data to json")
		return nil, err
	}

	if m.Shares == nil {
		m.Shares = map[string]interface{}{}
	}
	if m.ReceivedShares == nil {
		m.ReceivedShares = map[string]interface{}{}
	}
	m.file = file

	return m, nil
}

type shareModel struct {
	file           string
	Shares         map[string]interface{} `json:"shares"`
	ReceivedShares map[string]interface{} `json:"received_shares"`
}

type config struct {
	File                string `mapstructure:"file"`
	InsecureConnections bool   `mapstructure:"insecure_connections"`
}

func (c *config) init() {
	if c.File == "" {
		c.File = "/var/tmp/reva/ocm-shares.json"
	}
}

type mgr struct {
	c          *config
	sync.Mutex // concurrent access to the file
	model      *shareModel
}

func (m *shareModel) Save() error {
	data, err := json.Marshal(m)
	if err != nil {
		err = errors.Wrap(err, "error encoding to json")
		return err
	}

	if err := os.WriteFile(m.file, data, 0644); err != nil {
		err = errors.Wrap(err, "error writing to file: "+m.file)
		return err
	}

	return nil
}

func (m *shareModel) ReadFile() error {
	data, err := os.ReadFile(m.file)
	if err != nil {
		err = errors.Wrap(err, "error reading the data")
		return err
	}

	if err := json.Unmarshal(data, m); err != nil {
		err = errors.Wrap(err, "error decoding data to json")
		return err
	}

	return nil
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}
	return c, nil
}

func genID() string {
	return uuid.New().String()
}

// Called from both grpc CreateOCMShare for outgoing
// and http /ocm/shares for incoming
// pi is provider info
// pm is permissions
func (m *mgr) Share(ctx context.Context, md *provider.ResourceId, g *ocm.ShareGrant, name string,
	pi *ocmprovider.ProviderInfo, pm string, owner *userpb.UserId, token string, st ocm.Share_ShareType) (*ocm.Share, error) {
	id := genID()
	now := time.Now().UnixNano()
	ts := &typespb.Timestamp{
		Seconds: uint64(now / 1000000000),
		Nanos:   uint32(now % 1000000000),
	}

	// Since both OCMCore and OCMShareProvider use the same package, we distinguish
	// between calls received from them on the basis of whether they provide info
	// about the remote provider on which the share is to be created.
	// If this info is provided, this call is on the owner's mesh provider and so
	// we call the CreateOCMCoreShare method on the remote provider as well,
	// else this is received from another provider and we only create a local share.
	var isOwnersMeshProvider bool
	if pi != nil {
		isOwnersMeshProvider = true
	}

	var userID *userpb.UserId
	if !isOwnersMeshProvider {
		// Since this call is on the remote provider, the owner of the resource is expected to be specified.
		if owner == nil {
			return nil, errors.New("json: owner of resource not provided")
		}
		userID = owner
		g.Grantee.Opaque = &typespb.Opaque{
			Map: map[string]*typespb.OpaqueEntry{
				"token": {
					Decoder: "plain",
					Value:   []byte(token),
				},
			},
		}
	} else {
		userID = ctxpkg.ContextMustGetUser(ctx).GetId()
	}

	// do not allow share to myself if share is for a user
	if g.Grantee.Type == provider.GranteeType_GRANTEE_TYPE_USER && utils.UserEqual(g.Grantee.GetUserId(), userID) {
		return nil, errors.New("json: user and grantee are the same")
	}

	// check if share already exists.
	key := &ocm.ShareKey{
		Owner:      userID,
		ResourceId: md,
		Grantee:    g.Grantee,
	}
	_, err := m.getByKey(ctx, key)

	// share already exists
	if isOwnersMeshProvider && err == nil {
		return nil, errtypes.AlreadyExists(key.String())
	}

	s := &ocm.Share{
		Id: &ocm.ShareId{
			OpaqueId: id,
		},
		Name:        name,
		ResourceId:  md,
		Permissions: g.Permissions,
		Grantee:     g.Grantee,
		Owner:       userID,
		Creator:     userID,
		Ctime:       ts,
		Mtime:       ts,
		ShareType:   st,
	}

	if isOwnersMeshProvider {
		protocol := map[string]interface{}{
			"name": "webdav",
			"options": map[string]string{
				"permissions": pm,
				"token":       ctxpkg.ContextMustGetToken(ctx),
			},
		}
		if st == ocm.Share_SHARE_TYPE_TRANSFER {
			protocol["name"] = "datatx"
		}

		requestBodyMap := map[string]interface{}{
			"shareWith":    g.Grantee.GetUserId().OpaqueId,
			"name":         name,
			"providerId":   fmt.Sprintf("%s:%s", md.StorageId, md.OpaqueId),
			"owner":        userID.OpaqueId,
			"protocol":     protocol,
			"meshProvider": userID.Idp, // FIXME: move this into the 'owner' string?
		}
		err = sender.Send(requestBodyMap, pi)
		if err != nil {
			err = errors.Wrap(err, "error sending OCM POST")
			return nil, err
		}

	}

	m.Lock()
	if err := m.model.ReadFile(); err != nil {
		err = errors.Wrap(err, "error reading model")
		return nil, err
	}

	if isOwnersMeshProvider {
		encShare, err := utils.MarshalProtoV1ToJSON(s)
		if err != nil {
			return nil, err
		}
		m.model.Shares[s.Id.OpaqueId] = string(encShare)
	} else {
		encShare, err := utils.MarshalProtoV1ToJSON(&ocm.ReceivedShare{
			Share: s,
			State: ocm.ShareState_SHARE_STATE_PENDING,
		})
		if err != nil {
			return nil, err
		}
		m.model.ReceivedShares[s.Id.OpaqueId] = string(encShare)
	}

	if err := m.model.Save(); err != nil {
		err = errors.Wrap(err, "error saving model")
		return nil, err
	}
	m.Unlock()

	return s, nil
}

func (m *mgr) getByID(ctx context.Context, id *ocm.ShareId) (*ocm.Share, error) {
	m.Lock()
	defer m.Unlock()

	if err := m.model.ReadFile(); err != nil {
		err = errors.Wrap(err, "error reading model")
		return nil, err
	}

	if s, ok := m.model.Shares[id.OpaqueId]; ok {
		var share ocm.Share
		if err := utils.UnmarshalJSONToProtoV1([]byte(s.(string)), &share); err != nil {
			return nil, err
		}
		return &share, nil
	}

	return nil, errtypes.NotFound(id.String())
}

func (m *mgr) getByKey(ctx context.Context, key *ocm.ShareKey) (*ocm.Share, error) {
	m.Lock()
	defer m.Unlock()

	if err := m.model.ReadFile(); err != nil {
		err = errors.Wrap(err, "error reading model")
		return nil, err
	}

	for _, s := range m.model.Shares {
		var share ocm.Share
		if err := utils.UnmarshalJSONToProtoV1([]byte(s.(string)), &share); err != nil {
			continue
		}
		if (utils.UserEqual(key.Owner, share.Owner) || utils.UserEqual(key.Owner, share.Creator)) &&
			utils.ResourceIDEqual(key.ResourceId, share.ResourceId) && utils.GranteeEqual(key.Grantee, share.Grantee) {
			return &share, nil
		}
	}
	return nil, errtypes.NotFound(key.String())
}

func (m *mgr) get(ctx context.Context, ref *ocm.ShareReference) (s *ocm.Share, err error) {
	switch {
	case ref.GetId() != nil:
		s, err = m.getByID(ctx, ref.GetId())
	case ref.GetKey() != nil:
		s, err = m.getByKey(ctx, ref.GetKey())
	default:
		err = errtypes.NotFound(ref.String())
	}

	if err != nil {
		return nil, err
	}

	// check if we are the owner
	user := ctxpkg.ContextMustGetUser(ctx)
	if utils.UserEqual(user.Id, s.Owner) || utils.UserEqual(user.Id, s.Creator) {
		return s, nil
	}

	// we return not found to not disclose information
	return nil, errtypes.NotFound(ref.String())
}

func (m *mgr) GetShare(ctx context.Context, ref *ocm.ShareReference) (*ocm.Share, error) {
	share, err := m.get(ctx, ref)
	if err != nil {
		return nil, err
	}

	return share, nil
}

func (m *mgr) Unshare(ctx context.Context, ref *ocm.ShareReference) error {
	m.Lock()
	defer m.Unlock()

	if err := m.model.ReadFile(); err != nil {
		err = errors.Wrap(err, "error reading model")
		return err
	}

	user := ctxpkg.ContextMustGetUser(ctx)
	for id, s := range m.model.Shares {
		var share ocm.Share
		if err := utils.UnmarshalJSONToProtoV1([]byte(s.(string)), &share); err != nil {
			continue
		}
		if sharesEqual(ref, &share) {
			if utils.UserEqual(user.Id, share.Owner) || utils.UserEqual(user.Id, share.Creator) {
				delete(m.model.Shares, id)
				if err := m.model.Save(); err != nil {
					err = errors.Wrap(err, "error saving model")
					return err
				}
				return nil
			}
		}
	}
	return errtypes.NotFound(ref.String())
}

func sharesEqual(ref *ocm.ShareReference, s *ocm.Share) bool {
	if ref.GetId() != nil && s.Id != nil {
		if ref.GetId().OpaqueId == s.Id.OpaqueId {
			return true
		}
	} else if ref.GetKey() != nil {
		if (utils.UserEqual(ref.GetKey().Owner, s.Owner) || utils.UserEqual(ref.GetKey().Owner, s.Creator)) &&
			utils.ResourceIDEqual(ref.GetKey().ResourceId, s.ResourceId) && utils.GranteeEqual(ref.GetKey().Grantee, s.Grantee) {
			return true
		}
	}
	return false
}

func (m *mgr) UpdateShare(ctx context.Context, ref *ocm.ShareReference, p *ocm.SharePermissions) (*ocm.Share, error) {
	m.Lock()
	defer m.Unlock()

	if err := m.model.ReadFile(); err != nil {
		err = errors.Wrap(err, "error reading model")
		return nil, err
	}

	user := ctxpkg.ContextMustGetUser(ctx)
	for id, s := range m.model.Shares {
		var share ocm.Share
		if err := utils.UnmarshalJSONToProtoV1([]byte(s.(string)), &share); err != nil {
			continue
		}
		if sharesEqual(ref, &share) {
			if utils.UserEqual(user.Id, share.Owner) || utils.UserEqual(user.Id, share.Creator) {
				now := time.Now().UnixNano()
				share.Permissions = p
				share.Mtime = &typespb.Timestamp{
					Seconds: uint64(now / 1000000000),
					Nanos:   uint32(now % 1000000000),
				}
				encShare, err := utils.MarshalProtoV1ToJSON(&share)
				if err != nil {
					return nil, err
				}
				m.model.Shares[id] = string(encShare)
				if err := m.model.Save(); err != nil {
					err = errors.Wrap(err, "error saving model")
					return nil, err
				}
				return &share, nil
			}
		}
	}
	return nil, errtypes.NotFound(ref.String())
}

func (m *mgr) ListShares(ctx context.Context, filters []*ocm.ListOCMSharesRequest_Filter) ([]*ocm.Share, error) {
	var ss []*ocm.Share
	m.Lock()
	defer m.Unlock()

	if err := m.model.ReadFile(); err != nil {
		err = errors.Wrap(err, "error reading model")
		return nil, err
	}

	user := ctxpkg.ContextMustGetUser(ctx)
	for _, s := range m.model.Shares {
		var share ocm.Share
		if err := utils.UnmarshalJSONToProtoV1([]byte(s.(string)), &share); err != nil {
			continue
		}
		if utils.UserEqual(user.Id, share.Owner) || utils.UserEqual(user.Id, share.Creator) {
			// no filter we return earlier
			if len(filters) == 0 {
				ss = append(ss, &share)
			} else {
				// check filters
				// TODO(labkode): add the rest of filters.
				for _, f := range filters {
					if f.Type == ocm.ListOCMSharesRequest_Filter_TYPE_RESOURCE_ID {
						if utils.ResourceIDEqual(share.ResourceId, f.GetResourceId()) {
							ss = append(ss, &share)
						}
					}
				}
			}
		}
	}
	return ss, nil
}

func (m *mgr) ListReceivedShares(ctx context.Context) ([]*ocm.ReceivedShare, error) {
	var rss []*ocm.ReceivedShare
	m.Lock()
	defer m.Unlock()

	if err := m.model.ReadFile(); err != nil {
		err = errors.Wrap(err, "error reading model")
		return nil, err
	}

	user := ctxpkg.ContextMustGetUser(ctx)
	for _, s := range m.model.ReceivedShares {
		var rs ocm.ReceivedShare
		if err := utils.UnmarshalJSONToProtoV1([]byte(s.(string)), &rs); err != nil {
			continue
		}
		share := rs.Share
		if utils.UserEqual(user.Id, share.Owner) || utils.UserEqual(user.Id, share.Creator) {
			// omit shares created by me
			continue
		}
		if share.Grantee.Type == provider.GranteeType_GRANTEE_TYPE_USER && utils.UserEqual(user.Id, share.Grantee.GetUserId()) {
			rss = append(rss, &rs)
		}
	}
	return rss, nil
}

func (m *mgr) GetReceivedShare(ctx context.Context, ref *ocm.ShareReference) (*ocm.ReceivedShare, error) {
	return m.getReceived(ctx, ref)
}

func (m *mgr) getReceived(ctx context.Context, ref *ocm.ShareReference) (*ocm.ReceivedShare, error) {
	m.Lock()
	defer m.Unlock()

	if err := m.model.ReadFile(); err != nil {
		err = errors.Wrap(err, "error reading model")
		return nil, err
	}

	user := ctxpkg.ContextMustGetUser(ctx)
	for _, s := range m.model.ReceivedShares {
		var rs ocm.ReceivedShare
		if err := utils.UnmarshalJSONToProtoV1([]byte(s.(string)), &rs); err != nil {
			continue
		}
		share := rs.Share
		if sharesEqual(ref, share) {
			if share.Grantee.Type == provider.GranteeType_GRANTEE_TYPE_USER && utils.UserEqual(user.Id, share.Grantee.GetUserId()) {
				return &rs, nil
			}
		}
	}
	return nil, errtypes.NotFound(ref.String())
}

func (m *mgr) UpdateReceivedShare(ctx context.Context, share *ocm.ReceivedShare, fieldMask *field_mask.FieldMask) (*ocm.ReceivedShare, error) {
	rs, err := m.getReceived(ctx, &ocm.ShareReference{Spec: &ocm.ShareReference_Id{Id: share.Share.Id}})
	if err != nil {
		return nil, err
	}

	m.Lock()
	defer m.Unlock()

	for i := range fieldMask.Paths {
		switch fieldMask.Paths[i] {
		case "state":
			rs.State = share.State
		case "mount_point":
			rs.MountPoint = share.MountPoint
		default:
			return nil, errtypes.NotSupported("updating " + fieldMask.Paths[i] + " is not supported")
		}
	}

	if err := m.model.ReadFile(); err != nil {
		err = errors.Wrap(err, "error reading model")
		return nil, err
	}

	encShare, err := utils.MarshalProtoV1ToJSON(rs)
	if err != nil {
		return nil, err
	}
	m.model.ReceivedShares[rs.Share.Id.GetOpaqueId()] = string(encShare)

	if err := m.model.Save(); err != nil {
		err = errors.Wrap(err, "error saving model")
		return nil, err
	}

	return rs, nil
}
