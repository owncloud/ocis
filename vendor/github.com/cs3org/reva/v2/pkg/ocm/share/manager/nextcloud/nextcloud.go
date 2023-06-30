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

// Package nextcloud verifies a clientID and clientSecret against a Nextcloud backend.
package nextcloud

import (
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/utils"

	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"

	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/ocm/share"
	"github.com/cs3org/reva/v2/pkg/ocm/share/manager/registry"
	"github.com/cs3org/reva/v2/pkg/ocm/share/sender"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/genproto/protobuf/field_mask"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	registry.Register("nextcloud", New)
}

// Manager is the Nextcloud-based implementation of the share.Manager interface
// see https://github.com/cs3org/reva/blob/v1.13.0/pkg/ocm/share/share.go#L30-L57
type Manager struct {
	client       *http.Client
	sharedSecret string
	endPoint     string
}

// ShareManagerConfig contains config for a Nextcloud-based ShareManager
type ShareManagerConfig struct {
	EndPoint     string `mapstructure:"endpoint" docs:";The Nextcloud backend endpoint for user check"`
	SharedSecret string `mapstructure:"shared_secret"`
	MockHTTP     bool   `mapstructure:"mock_http"`
}

// Action describes a REST request to forward to the Nextcloud backend
type Action struct {
	verb string
	argS string
}

// GranteeAltMap is an alternative map to JSON-unmarshal a Grantee
// Grantees are hard to unmarshal, so unmarshalling into a map[string]interface{} first,
// see also https://github.com/pondersource/sciencemesh-nextcloud/issues/27
type GranteeAltMap struct {
	ID *provider.Grantee_UserId `json:"id"`
}

// ShareAltMap is an alternative map to JSON-unmarshal a Share
type ShareAltMap struct {
	ID          *ocm.ShareId          `json:"id"`
	ResourceID  *provider.ResourceId  `json:"resource_id"`
	Permissions *ocm.SharePermissions `json:"permissions"`
	Grantee     *GranteeAltMap        `json:"grantee"`
	Owner       *userpb.UserId        `json:"owner"`
	Creator     *userpb.UserId        `json:"creator"`
	Ctime       *typespb.Timestamp    `json:"ctime"`
	Mtime       *typespb.Timestamp    `json:"mtime"`
}

// ReceivedShareAltMap is an alternative map to JSON-unmarshal a ReceivedShare
type ReceivedShareAltMap struct {
	Share *ShareAltMap   `json:"share"`
	State ocm.ShareState `json:"state"`
}

func (c *ShareManagerConfig) init() {
}

func parseConfig(m map[string]interface{}) (*ShareManagerConfig, error) {
	c := &ShareManagerConfig{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

func getUser(ctx context.Context) (*userpb.User, error) {
	u, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		err := errors.Wrap(errtypes.UserRequired(""), "nextcloud storage driver: error getting user from ctx")
		return nil, err
	}
	return u, nil
}

// New returns a share manager implementation that verifies against a Nextcloud backend.
func New(m map[string]interface{}) (share.Manager, error) {
	c, err := parseConfig(m)
	if err != nil {
		return nil, err
	}
	c.init()

	return NewShareManager(c)
}

// NewShareManager returns a new Nextcloud-based ShareManager
func NewShareManager(c *ShareManagerConfig) (*Manager, error) {
	var client *http.Client
	if c.MockHTTP {
		// called := make([]string, 0)
		// nextcloudServerMock := GetNextcloudServerMock(&called)
		// client, _ = TestingHTTPClient(nextcloudServerMock)

		// Wait for SetHTTPClient to be called later
		client = nil
	} else {
		if len(c.EndPoint) == 0 {
			return nil, errors.New("Please specify 'endpoint' in '[grpc.services.ocmshareprovider.drivers.nextcloud]' and  '[grpc.services.ocmcore.drivers.nextcloud]'")
		}
		client = &http.Client{}
	}

	return &Manager{
		endPoint:     c.EndPoint, // e.g. "http://nc/apps/sciencemesh/"
		sharedSecret: c.SharedSecret,
		client:       client,
	}, nil
}

// SetHTTPClient sets the HTTP client
func (sm *Manager) SetHTTPClient(c *http.Client) {
	sm.client = c
}

func getUsername(ctx context.Context) string {
	user, err := getUser(ctx)
	if err != nil {
		return "unknown"
	}
	if len(user.Username) > 0 {
		return user.Username
	}
	if len(user.Id.OpaqueId) > 0 {
		return user.Id.OpaqueId
	}

	return "empty-username"
}

func (sm *Manager) do(ctx context.Context, a Action, username string) (int, []byte, error) {
	url := sm.endPoint + "~" + username + "/api/ocm/" + a.verb

	log := appctx.GetLogger(ctx)
	log.Info().Msgf("am.do %s %s", url, a.argS)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(a.argS))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("X-Reva-Secret", sm.sharedSecret)

	req.Header.Set("Content-Type", "application/json")
	resp, err := sm.client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	// curl -i -H 'application/json' -H 'X-Reva-Secret: shared-secret-1' -d '{"md":{"opaque_id":"fileid-/other/q/as"},"g":{"grantee":{"type":1,"Id":{"UserId":{"idp":"revanc2.docker","opaque_id":"marie"}}},"permissions":{"permissions":{"get_path":true,"initiate_file_download":true,"list_container":true,"list_file_versions":true,"stat":true}}},"provider_domain":"cern.ch","resource_type":"file","provider_id":2,"owner_opaque_id":"einstein","owner_display_name":"Albert Einstein","protocol":{"name":"webdav","options":{"sharedSecret":"secret","permissions":"webdav-property"}}}' https://nc1.docker/index.php/apps/sciencemesh/~/api/ocm/addSentShare

	log.Info().Msgf("am.do response %d %s", resp.StatusCode, body)
	return resp.StatusCode, body, nil
}

// Share is called from both grpc CreateOCMShare for outgoing
// and http /ocm/shares for incoming
// pi is provider info
// pm is permissions
func (sm *Manager) Share(ctx context.Context, md *provider.ResourceId, g *ocm.ShareGrant, name string,
	pi *ocmprovider.ProviderInfo, pm string, owner *userpb.UserId, token string, st ocm.Share_ShareType) (*ocm.Share, error) {

	// Since both OCMCore and OCMShareProvider use the same package, we distinguish
	// between calls received from them on the basis of whether they provide info
	// about the remote provider on which the share is to be created.
	// If this info is provided, this call is on the owner's mesh provider and so
	// we call the CreateOCMCoreShare method on the remote provider as well,
	// else this is received from another provider and we only create a local share.
	var isOwnersMeshProvider bool
	var apiMethod string
	var username string
	if pi != nil {
		isOwnersMeshProvider = true
		apiMethod = "addSentShare"
		username = getUsername(ctx)
		token = randSeq(10)
	} else {
		apiMethod = "addReceivedShare"
		username = g.Grantee.GetUserId().OpaqueId
	}

	var userID *userpb.UserId
	if !isOwnersMeshProvider {
		// Since this call is on the remote provider, the owner of the resource is expected to be specified.
		if owner == nil {
			return nil, errors.New("nextcloud: owner of resource not provided")
		}
		userID = owner
	} else {
		userID = ctxpkg.ContextMustGetUser(ctx).GetId()
	}

	// do not allow share to myself if share is for a user
	if g.Grantee.Type == provider.GranteeType_GRANTEE_TYPE_USER && utils.UserEqual(g.Grantee.GetUserId(), userID) {
		return nil, errors.New("nextcloud: user and grantee are the same")
	}

	s := &ocm.Share{
		Name:        name,
		ResourceId:  md,
		Permissions: g.Permissions,
		Grantee:     g.Grantee,
		Owner:       userID,
		Creator:     userID,
		ShareType:   st,
	}

	var encShare []byte
	var err error

	if isOwnersMeshProvider {
		// adding the webdav sharedSecret in the Grantee because Share itself doesn't have an Opaque field,
		// see https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.Grantee
		// and https://cs3org.github.io/cs3apis/#cs3.sharing.ocm.v1beta1.Share
		s.Grantee.Opaque = &typespb.Opaque{
			Map: map[string]*typespb.OpaqueEntry{
				"sharedSecret": {
					Decoder: "plain",
					Value:   []byte(token),
				},
			},
		}

		encShare, err = utils.MarshalProtoV1ToJSON(s)
		if err != nil {
			return nil, err
		}
	} else {
		// adding the webdav sharedSecret and remote share id (called the "ProviderID" in OCM) in the Grantee because Share itself doesn't have an Opaque field,
		// see https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.Grantee
		// and https://cs3org.github.io/cs3apis/#cs3.sharing.ocm.v1beta1.Share
		s.Grantee.Opaque = &typespb.Opaque{
			Map: map[string]*typespb.OpaqueEntry{
				"sharedSecret": {
					Decoder: "plain",
					Value:   []byte(token),
				},
				"remoteShareId": {
					Decoder: "plain",
					Value:   g.Grantee.Opaque.Map["remoteShareId"].Value,
				},
			},
		}

		encShare, err = utils.MarshalProtoV1ToJSON(&ocm.ReceivedShare{
			Share: s,
			State: ocm.ShareState_SHARE_STATE_PENDING,
		})
		if err != nil {
			return nil, err
		}
	}

	_, body, err := sm.do(ctx, Action{apiMethod, string(encShare)}, username)

	s.Id = &ocm.ShareId{
		OpaqueId: string(body),
	}
	now := time.Now().UnixNano()
	s.Ctime = &typespb.Timestamp{
		Seconds: uint64(now / 1000000000),
		Nanos:   uint32(now % 1000000000),
	}
	if err != nil {
		return nil, err
	}

	if isOwnersMeshProvider {
		// token, ok := ctxpkg.ContextGetToken(ctx)
		// if !ok {
		// 	return nil, errors.New("Could not get token from context")
		// }
		var protocol map[string]interface{}
		if st == ocm.Share_SHARE_TYPE_TRANSFER {
			protocol = map[string]interface{}{
				"name": "datatx",
				"options": map[string]string{
					"permissions": pm,
					"token":       token, // FIXME: Where is the token for datatx generated?
				},
			}
		} else {
			protocol = map[string]interface{}{
				"name": "webdav",
				"options": map[string]string{
					"permissions":  pm,
					"sharedSecret": token,
				},
			}
		}
		requestBodyMap := map[string]interface{}{
			"shareWith":    g.Grantee.GetUserId().OpaqueId,
			"name":         name,
			"providerId":   s.Id.OpaqueId,
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

	return s, nil
}

// GetShare as defined in the ocm.share.Manager interface
// https://github.com/cs3org/reva/blob/v1.13.0/pkg/ocm/share/share.go#L30-L57
func (sm *Manager) GetShare(ctx context.Context, ref *ocm.ShareReference) (*ocm.Share, error) {
	bodyStr, err := json.Marshal(ref)
	if err != nil {
		return nil, err
	}
	_, body, err := sm.do(ctx, Action{"GetShare", string(bodyStr)}, getUsername(ctx))
	if err != nil {
		return nil, err
	}

	altResult := &ShareAltMap{}
	err = json.Unmarshal(body, &altResult)
	if altResult == nil {
		return nil, err
	}
	return &ocm.Share{
		Id:          altResult.ID,
		ResourceId:  altResult.ResourceID,
		Permissions: altResult.Permissions,
		Grantee: &provider.Grantee{
			Id: altResult.Grantee.ID,
		},
		Owner:   altResult.Owner,
		Creator: altResult.Creator,
		Ctime:   altResult.Ctime,
		Mtime:   altResult.Mtime,
	}, err
}

// Unshare as defined in the ocm.share.Manager interface
// https://github.com/cs3org/reva/blob/v1.13.0/pkg/ocm/share/share.go#L30-L57
func (sm *Manager) Unshare(ctx context.Context, ref *ocm.ShareReference) error {
	bodyStr, err := json.Marshal(ref)
	if err != nil {
		return err
	}

	_, _, err = sm.do(ctx, Action{"Unshare", string(bodyStr)}, getUsername(ctx))
	return err
}

// UpdateShare as defined in the ocm.share.Manager interface
// https://github.com/cs3org/reva/blob/v1.13.0/pkg/ocm/share/share.go#L30-L57
func (sm *Manager) UpdateShare(ctx context.Context, ref *ocm.ShareReference, p *ocm.SharePermissions) (*ocm.Share, error) {
	type paramsObj struct {
		Ref *ocm.ShareReference   `json:"ref"`
		P   *ocm.SharePermissions `json:"p"`
	}
	bodyObj := &paramsObj{
		Ref: ref,
		P:   p,
	}
	bodyStr, err := json.Marshal(bodyObj)
	if err != nil {
		return nil, err
	}

	_, body, err := sm.do(ctx, Action{"UpdateShare", string(bodyStr)}, getUsername(ctx))

	if err != nil {
		return nil, err
	}

	altResult := &ShareAltMap{}
	err = json.Unmarshal(body, &altResult)
	if altResult == nil {
		return nil, err
	}
	return &ocm.Share{
		Id:          altResult.ID,
		ResourceId:  altResult.ResourceID,
		Permissions: altResult.Permissions,
		Grantee: &provider.Grantee{
			Id: altResult.Grantee.ID,
		},
		Owner:   altResult.Owner,
		Creator: altResult.Creator,
		Ctime:   altResult.Ctime,
		Mtime:   altResult.Mtime,
	}, err
}

// ListShares as defined in the ocm.share.Manager interface
// https://github.com/cs3org/reva/blob/v1.13.0/pkg/ocm/share/share.go#L30-L57
func (sm *Manager) ListShares(ctx context.Context, filters []*ocm.ListOCMSharesRequest_Filter) ([]*ocm.Share, error) {
	bodyStr, err := json.Marshal(filters)
	if err != nil {
		return nil, err
	}

	_, respBody, err := sm.do(ctx, Action{"ListShares", string(bodyStr)}, getUsername(ctx))
	if err != nil {
		return nil, err
	}

	var respArr []ShareAltMap
	err = json.Unmarshal(respBody, &respArr)
	if err != nil {
		return nil, err
	}

	var pointers = make([]*ocm.Share, len(respArr))
	for i := 0; i < len(respArr); i++ {
		altResult := respArr[i]
		pointers[i] = &ocm.Share{
			Id:          altResult.ID,
			ResourceId:  altResult.ResourceID,
			Permissions: altResult.Permissions,
			Grantee: &provider.Grantee{
				Id: altResult.Grantee.ID,
			},
			Owner:   altResult.Owner,
			Creator: altResult.Creator,
			Ctime:   altResult.Ctime,
			Mtime:   altResult.Mtime,
		}
	}
	return pointers, err
}

// ListReceivedShares as defined in the ocm.share.Manager interface
// https://github.com/cs3org/reva/blob/v1.13.0/pkg/ocm/share/share.go#L30-L57
func (sm *Manager) ListReceivedShares(ctx context.Context) ([]*ocm.ReceivedShare, error) {
	_, respBody, err := sm.do(ctx, Action{"ListReceivedShares", string("")}, getUsername(ctx))
	if err != nil {
		return nil, err
	}

	var respArr []ReceivedShareAltMap
	err = json.Unmarshal(respBody, &respArr)
	if err != nil {
		return nil, err
	}
	var pointers = make([]*ocm.ReceivedShare, len(respArr))
	for i := 0; i < len(respArr); i++ {
		altResultShare := respArr[i].Share
		if altResultShare == nil {
			pointers[i] = &ocm.ReceivedShare{
				Share: nil,
				State: respArr[i].State,
			}
		} else {
			pointers[i] = &ocm.ReceivedShare{
				Share: &ocm.Share{
					Id:          altResultShare.ID,
					ResourceId:  altResultShare.ResourceID,
					Permissions: altResultShare.Permissions,
					Grantee: &provider.Grantee{
						Id: altResultShare.Grantee.ID,
					},
					Owner:   altResultShare.Owner,
					Creator: altResultShare.Creator,
					Ctime:   altResultShare.Ctime,
					Mtime:   altResultShare.Mtime,
				},
				State: respArr[i].State,
			}
		}
	}
	return pointers, err

}

// GetReceivedShare as defined in the ocm.share.Manager interface
// https://github.com/cs3org/reva/blob/v1.13.0/pkg/ocm/share/share.go#L29-L54
func (sm *Manager) GetReceivedShare(ctx context.Context, ref *ocm.ShareReference) (*ocm.ReceivedShare, error) {
	bodyStr, err := json.Marshal(ref)
	if err != nil {
		return nil, err
	}

	_, respBody, err := sm.do(ctx, Action{"GetReceivedShare", string(bodyStr)}, getUsername(ctx))
	if err != nil {
		return nil, err
	}

	var altResult ReceivedShareAltMap
	err = json.Unmarshal(respBody, &altResult)
	if err != nil {
		return nil, err
	}
	altResultShare := altResult.Share
	if altResultShare == nil {
		return &ocm.ReceivedShare{
			Share: nil,
			State: altResult.State,
		}, err
	}
	return &ocm.ReceivedShare{
		Share: &ocm.Share{
			Id:          altResultShare.ID,
			ResourceId:  altResultShare.ResourceID,
			Permissions: altResultShare.Permissions,
			Grantee: &provider.Grantee{
				Id: altResultShare.Grantee.ID,
			},
			Owner:   altResultShare.Owner,
			Creator: altResultShare.Creator,
			Ctime:   altResultShare.Ctime,
			Mtime:   altResultShare.Mtime,
		},
		State: altResult.State,
	}, err
}

// UpdateReceivedShare as defined in the ocm.share.Manager interface
// https://github.com/cs3org/reva/blob/v1.13.0/pkg/ocm/share/share.go#L30-L57
func (sm Manager) UpdateReceivedShare(ctx context.Context, receivedShare *ocm.ReceivedShare, fieldMask *field_mask.FieldMask) (*ocm.ReceivedShare, error) {
	type paramsObj struct {
		ReceivedShare *ocm.ReceivedShare    `json:"received_share"`
		FieldMask     *field_mask.FieldMask `json:"field_mask"`
	}

	bodyObj := &paramsObj{
		ReceivedShare: receivedShare,
		FieldMask:     fieldMask,
	}
	bodyStr, err := json.Marshal(bodyObj)
	if err != nil {
		return nil, err
	}

	_, respBody, err := sm.do(ctx, Action{"UpdateReceivedShare", string(bodyStr)}, getUsername(ctx))
	if err != nil {
		return nil, err
	}

	var altResult ReceivedShareAltMap
	err = json.Unmarshal(respBody, &altResult)
	if err != nil {
		return nil, err
	}
	altResultShare := altResult.Share
	if altResultShare == nil {
		return &ocm.ReceivedShare{
			Share: nil,
			State: altResult.State,
		}, err
	}
	return &ocm.ReceivedShare{
		Share: &ocm.Share{
			Id:          altResultShare.ID,
			ResourceId:  altResultShare.ResourceID,
			Permissions: altResultShare.Permissions,
			Grantee: &provider.Grantee{
				Id: altResultShare.Grantee.ID,
			},
			Owner:   altResultShare.Owner,
			Creator: altResultShare.Creator,
			Ctime:   altResultShare.Ctime,
			Mtime:   altResultShare.Mtime,
		},
		State: altResult.State,
	}, err
}
