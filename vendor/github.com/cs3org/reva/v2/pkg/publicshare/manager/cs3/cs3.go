// Copyright 2018-2022 CERN
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

package cs3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/publicshare"
	"github.com/cs3org/reva/v2/pkg/publicshare/manager/registry"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storage/utils/indexer"
	indexerErrors "github.com/cs3org/reva/v2/pkg/storage/utils/indexer/errors"
	"github.com/cs3org/reva/v2/pkg/storage/utils/indexer/option"
	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
)

func init() {
	registry.Register("cs3", NewDefault)
}

// Manager implements a publicshare manager using a cs3 storage backend
type Manager struct {
	gatewayClient gateway.GatewayAPIClient
	sync.RWMutex

	storage          metadata.Storage
	indexer          indexer.Indexer
	passwordHashCost int

	initialized bool
}

type config struct {
	GatewayAddr       string `mapstructure:"gateway_addr"`
	ProviderAddr      string `mapstructure:"provider_addr"`
	ServiceUserID     string `mapstructure:"service_user_id"`
	ServiceUserIdp    string `mapstructure:"service_user_idp"`
	MachineAuthAPIKey string `mapstructure:"machine_auth_apikey"`
}

// NewDefault returns a new manager instance with default dependencies
func NewDefault(m map[string]interface{}) (publicshare.Manager, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error creating a new manager")
		return nil, err
	}

	s, err := metadata.NewCS3Storage(c.GatewayAddr, c.ProviderAddr, c.ServiceUserID, c.ServiceUserIdp, c.MachineAuthAPIKey)
	if err != nil {
		return nil, err
	}
	indexer := indexer.CreateIndexer(s)

	client, err := pool.GetGatewayServiceClient(c.GatewayAddr)
	if err != nil {
		return nil, err
	}

	return New(client, s, indexer, bcrypt.DefaultCost)
}

// New returns a new manager instance
func New(gatewayClient gateway.GatewayAPIClient, storage metadata.Storage, indexer indexer.Indexer, passwordHashCost int) (*Manager, error) {
	return &Manager{
		gatewayClient:    gatewayClient,
		storage:          storage,
		indexer:          indexer,
		passwordHashCost: passwordHashCost,
		initialized:      false,
	}, nil
}

func (m *Manager) initialize() error {
	if m.initialized {
		return nil
	}

	m.Lock()
	defer m.Unlock()

	if m.initialized { // check if initialization happened while grabbing the lock
		return nil
	}

	err := m.storage.Init(context.Background(), "public-share-manager-metadata")
	if err != nil {
		return err
	}
	if err := m.storage.MakeDirIfNotExist(context.Background(), "publicshares"); err != nil {
		return err
	}
	err = m.indexer.AddIndex(&link.PublicShare{}, option.IndexByField("Id.OpaqueId"), "Token", "publicshares", "unique", nil, true)
	if err != nil {
		return err
	}
	err = m.indexer.AddIndex(&link.PublicShare{}, option.IndexByFunc{
		Name: "Owner",
		Func: indexOwnerFunc,
	}, "Token", "publicshares", "non_unique", nil, true)
	if err != nil {
		return err
	}
	err = m.indexer.AddIndex(&link.PublicShare{}, option.IndexByFunc{
		Name: "Creator",
		Func: indexCreatorFunc,
	}, "Token", "publicshares", "non_unique", nil, true)
	if err != nil {
		return err
	}
	err = m.indexer.AddIndex(&link.PublicShare{}, option.IndexByFunc{
		Name: "ResourceId",
		Func: indexResourceIDFunc,
	}, "Token", "publicshares", "non_unique", nil, true)
	if err != nil {
		return err
	}
	m.initialized = true
	return nil
}

// Dump exports public shares to channels (e.g. during migration)
func (m *Manager) Dump(ctx context.Context, shareChan chan<- *publicshare.WithPassword) error {
	if err := m.initialize(); err != nil {
		return err
	}

	pshares, err := m.storage.ListDir(ctx, "publicshares")
	if err != nil {
		return err
	}

	for _, v := range pshares {
		var local publicshare.WithPassword
		ps, err := m.getByToken(ctx, v.Name)
		if err != nil {
			return err
		}
		local.Password = ps.Password
		local.PublicShare = ps.PublicShare

		shareChan <- &local
	}

	return nil
}

// Load imports public shares and received shares from channels (e.g. during migration)
func (m *Manager) Load(ctx context.Context, shareChan <-chan *publicshare.WithPassword) error {
	log := appctx.GetLogger(ctx)
	if err := m.initialize(); err != nil {
		return err
	}
	for ps := range shareChan {
		if err := m.persist(context.Background(), ps); err != nil {
			log.Error().Err(err).Interface("publicshare", ps).Msg("error loading public share")
		}
	}
	return nil
}

// CreatePublicShare creates a new public share
func (m *Manager) CreatePublicShare(ctx context.Context, u *user.User, ri *provider.ResourceInfo, g *link.Grant) (*link.PublicShare, error) {
	if err := m.initialize(); err != nil {
		return nil, err
	}

	id := &link.PublicShareId{
		OpaqueId: utils.RandString(15),
	}

	tkn := utils.RandString(15)
	now := time.Now().UnixNano()

	displayName, quicklink := tkn, false
	if ri.ArbitraryMetadata != nil {
		metadataName, ok := ri.ArbitraryMetadata.Metadata["name"]
		if ok {
			displayName = metadataName
		}

		quicklink, _ = strconv.ParseBool(ri.ArbitraryMetadata.Metadata["quicklink"])
	}

	var passwordProtected bool
	password := g.Password
	if password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(password), m.passwordHashCost)
		if err != nil {
			return nil, errors.Wrap(err, "could not hash share password")
		}
		password = string(h)
		passwordProtected = true
	}

	createdAt := &typespb.Timestamp{
		Seconds: uint64(now / int64(time.Second)),
		Nanos:   uint32(now % int64(time.Second)),
	}

	s := &publicshare.WithPassword{
		PublicShare: link.PublicShare{
			Id:                id,
			Owner:             ri.GetOwner(),
			Creator:           u.Id,
			ResourceId:        ri.Id,
			Token:             tkn,
			Permissions:       g.Permissions,
			Ctime:             createdAt,
			Mtime:             createdAt,
			PasswordProtected: passwordProtected,
			Expiration:        g.Expiration,
			DisplayName:       displayName,
			Quicklink:         quicklink,
		},
		Password: password,
	}

	err := m.persist(ctx, s)
	if err != nil {
		return nil, err
	}

	return &s.PublicShare, nil
}

// UpdatePublicShare updates an existing public share
func (m *Manager) UpdatePublicShare(ctx context.Context, u *user.User, req *link.UpdatePublicShareRequest) (*link.PublicShare, error) {
	if err := m.initialize(); err != nil {
		return nil, err
	}

	ps, err := m.getWithPassword(ctx, req.Ref)
	if err != nil {
		return nil, err
	}

	switch req.Update.Type {
	case link.UpdatePublicShareRequest_Update_TYPE_DISPLAYNAME:
		ps.PublicShare.DisplayName = req.Update.DisplayName
	case link.UpdatePublicShareRequest_Update_TYPE_PERMISSIONS:
		ps.PublicShare.Permissions = req.Update.Grant.Permissions
	case link.UpdatePublicShareRequest_Update_TYPE_EXPIRATION:
		ps.PublicShare.Expiration = req.Update.Grant.Expiration
	case link.UpdatePublicShareRequest_Update_TYPE_PASSWORD:
		if req.Update.Grant.Password == "" {
			ps.Password = ""
			ps.PublicShare.PasswordProtected = false
		} else {
			h, err := bcrypt.GenerateFromPassword([]byte(req.Update.Grant.Password), m.passwordHashCost)
			if err != nil {
				return nil, errors.Wrap(err, "could not hash share password")
			}
			ps.Password = string(h)
			ps.PublicShare.PasswordProtected = true
		}
	default:
		return nil, errtypes.BadRequest("no valid update type given")
	}

	err = m.persist(ctx, ps)
	if err != nil {
		return nil, err
	}

	return &ps.PublicShare, nil
}

// GetPublicShare returns an existing public share
func (m *Manager) GetPublicShare(ctx context.Context, u *user.User, ref *link.PublicShareReference, sign bool) (*link.PublicShare, error) {
	if err := m.initialize(); err != nil {
		return nil, err
	}

	ps, err := m.getWithPassword(ctx, ref)
	if err != nil {
		return nil, err
	}

	if ps.PublicShare.PasswordProtected && sign {
		err = publicshare.AddSignature(&ps.PublicShare, ps.Password)
		if err != nil {
			return nil, err
		}
	}

	return &ps.PublicShare, nil
}

func (m *Manager) getWithPassword(ctx context.Context, ref *link.PublicShareReference) (*publicshare.WithPassword, error) {
	switch {
	case ref.GetToken() != "":
		return m.getByToken(ctx, ref.GetToken())
	case ref.GetId().GetOpaqueId() != "":
		return m.getByID(ctx, ref.GetId().GetOpaqueId())
	default:
		return nil, errtypes.BadRequest("neither id nor token given")
	}
}

func (m *Manager) getByID(ctx context.Context, id string) (*publicshare.WithPassword, error) {
	tokens, err := m.indexer.FindBy(&link.PublicShare{},
		indexer.NewField("Id.OpaqueId", id),
	)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, errtypes.NotFound("publicshare with the given id not found")
	}
	return m.getByToken(ctx, tokens[0])
}

func (m *Manager) getByToken(ctx context.Context, token string) (*publicshare.WithPassword, error) {
	fn := path.Join("publicshares", token)
	data, err := m.storage.SimpleDownload(ctx, fn)
	if err != nil {
		return nil, err
	}

	ps := &publicshare.WithPassword{}
	err = json.Unmarshal(data, ps)
	if err != nil {
		return nil, err
	}
	id := storagespace.UpdateLegacyResourceID(*ps.PublicShare.ResourceId)
	ps.PublicShare.ResourceId = &id
	return ps, nil
}

// ListPublicShares lists existing public shares matching the given filters
func (m *Manager) ListPublicShares(ctx context.Context, u *user.User, filters []*link.ListPublicSharesRequest_Filter, sign bool) ([]*link.PublicShare, error) {
	if err := m.initialize(); err != nil {
		return nil, err
	}

	log := appctx.GetLogger(ctx)
	var rIDs []*provider.ResourceId
	if len(filters) != 0 {
		grouped := publicshare.GroupFiltersByType(filters)
		for _, g := range grouped {
			for _, f := range g {
				if f.GetResourceId() != nil {
					rIDs = append(rIDs, f.GetResourceId())
				}
			}
		}
	}
	var (
		createdShareTokens []string
		err                error
	)

	// in spaces, always use the resourceId
	if len(rIDs) != 0 {
		for _, rID := range rIDs {
			shareTokens, err := m.indexer.FindBy(&link.PublicShare{},
				indexer.NewField("ResourceId", resourceIDToIndex(rID)),
			)
			if err != nil {
				return nil, err
			}
			createdShareTokens = append(createdShareTokens, shareTokens...)
		}
	} else {
		// fallback for legacy use
		createdShareTokens, err = m.indexer.FindBy(&link.PublicShare{},
			indexer.NewField("Owner", userIDToIndex(u.Id)),
			indexer.NewField("Creator", userIDToIndex(u.Id)),
		)
		if err != nil {
			return nil, err
		}
	}

	// We use shareMem as a temporary lookup store to check which shares were
	// already added. This is to prevent duplicates.
	shareMem := make(map[string]struct{})
	result := []*link.PublicShare{}
	for _, token := range createdShareTokens {
		ps, err := m.getByToken(ctx, token)
		if err != nil {
			return nil, err
		}

		if !publicshare.MatchesFilters(ps.PublicShare, filters) {
			continue
		}

		if publicshare.IsExpired(ps.PublicShare) {
			ref := &link.PublicShareReference{
				Spec: &link.PublicShareReference_Id{
					Id: ps.PublicShare.Id,
				},
			}
			if err := m.RevokePublicShare(ctx, u, ref); err != nil {
				log.Error().Err(err).
					Str("public_share_token", ps.PublicShare.Token).
					Str("public_share_id", ps.PublicShare.Id.OpaqueId).
					Msg("failed to revoke expired public share")
			}
			continue
		}

		if ps.PublicShare.PasswordProtected && sign {
			err = publicshare.AddSignature(&ps.PublicShare, ps.Password)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, &ps.PublicShare)
		shareMem[ps.PublicShare.Token] = struct{}{}
	}

	// If a user requests to list shares which have not been created by them
	// we have to explicitly fetch these shares and check if the user is
	// allowed to list the shares.
	// Only then can we add these shares to the result.
	grouped := publicshare.GroupFiltersByType(filters)
	idFilter, ok := grouped[link.ListPublicSharesRequest_Filter_TYPE_RESOURCE_ID]
	if !ok {
		return result, nil
	}

	var tokens []string
	if len(idFilter) > 0 {
		idFilters := make([]indexer.Field, 0, len(idFilter))
		for _, filter := range idFilter {
			resourceID := filter.GetResourceId()
			idFilters = append(idFilters, indexer.NewField("ResourceId", resourceIDToIndex(resourceID)))
		}
		tokens, err = m.indexer.FindBy(&link.PublicShare{}, idFilters...)
		if err != nil {
			return nil, err
		}
	}

	// statMem is used as a local cache to prevent statting resources which
	// already have been checked.
	statMem := make(map[string]struct{})
	for _, token := range tokens {
		if _, handled := shareMem[token]; handled {
			// We don't want to add a share multiple times when we added it
			// already.
			continue
		}

		s, err := m.getByToken(ctx, token)
		if err != nil {
			return nil, err
		}

		if _, checked := statMem[resourceIDToIndex(s.PublicShare.GetResourceId())]; !checked {
			sReq := &provider.StatRequest{
				Ref: &provider.Reference{ResourceId: s.PublicShare.GetResourceId()},
			}
			sRes, err := m.gatewayClient.Stat(ctx, sReq)
			if err != nil {
				continue
			}
			if sRes.Status.Code != rpc.Code_CODE_OK {
				continue
			}
			if !sRes.Info.PermissionSet.ListGrants {
				continue
			}
			statMem[resourceIDToIndex(s.PublicShare.GetResourceId())] = struct{}{}
		}

		if publicshare.MatchesFilters(s.PublicShare, filters) {
			result = append(result, &s.PublicShare)
			shareMem[s.PublicShare.Token] = struct{}{}
		}
	}
	return result, nil
}

// RevokePublicShare revokes an existing public share
func (m *Manager) RevokePublicShare(ctx context.Context, u *user.User, ref *link.PublicShareReference) error {
	if err := m.initialize(); err != nil {
		return err
	}

	ps, err := m.GetPublicShare(ctx, u, ref, false)
	if err != nil {
		return err
	}

	err = m.storage.Delete(ctx, path.Join("publicshares", ps.Token))
	if err != nil {
		if _, ok := err.(errtypes.NotFound); !ok {
			return err
		}
	}

	return m.indexer.Delete(ps)
}

// GetPublicShareByToken gets an existing public share in an unauthenticated context using either a password or a signature
func (m *Manager) GetPublicShareByToken(ctx context.Context, token string, auth *link.PublicShareAuthentication, sign bool) (*link.PublicShare, error) {
	if err := m.initialize(); err != nil {
		return nil, err
	}

	ps, err := m.getByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if publicshare.IsExpired(ps.PublicShare) {
		return nil, errtypes.NotFound("public share has expired")
	}

	if ps.PublicShare.PasswordProtected {
		if !publicshare.Authenticate(&ps.PublicShare, ps.Password, auth) {
			return nil, errtypes.InvalidCredentials("access denied")
		}
	}

	return &ps.PublicShare, nil
}

func indexOwnerFunc(v interface{}) (string, error) {
	ps, ok := v.(*link.PublicShare)
	if !ok {
		return "", fmt.Errorf("given entity is not a public share")
	}
	return userIDToIndex(ps.Owner), nil
}

func indexCreatorFunc(v interface{}) (string, error) {
	ps, ok := v.(*link.PublicShare)
	if !ok {
		return "", fmt.Errorf("given entity is not a public share")
	}
	return userIDToIndex(ps.Creator), nil
}

func indexResourceIDFunc(v interface{}) (string, error) {
	ps, ok := v.(*link.PublicShare)
	if !ok {
		return "", fmt.Errorf("given entity is not a public share")
	}
	return resourceIDToIndex(ps.ResourceId), nil
}

func userIDToIndex(id *user.UserId) string {
	return url.QueryEscape(id.Idp + ":" + id.OpaqueId)
}

func resourceIDToIndex(id *provider.ResourceId) string {
	return strings.Join([]string{id.StorageId, id.OpaqueId}, "!")
}

func (m *Manager) persist(ctx context.Context, ps *publicshare.WithPassword) error {
	data, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	fn := path.Join("publicshares", ps.PublicShare.Token)
	err = m.storage.SimpleUpload(ctx, fn, data)
	if err != nil {
		return err
	}

	_, err = m.indexer.Add(&ps.PublicShare)
	if err != nil {
		if _, ok := err.(*indexerErrors.AlreadyExistsErr); ok {
			return nil
		}
		err = m.indexer.Delete(&ps.PublicShare)
		if err != nil {
			return err
		}
		_, err = m.indexer.Add(&ps.PublicShare)
		return err
	}

	return nil
}
