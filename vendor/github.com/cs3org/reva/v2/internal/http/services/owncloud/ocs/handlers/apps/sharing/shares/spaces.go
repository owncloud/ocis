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

package shares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	groupv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaborationv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	registry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/response"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	sdk "github.com/cs3org/reva/v2/pkg/sdk/common"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/pkg/errors"
)

func (h *Handler) getGrantee(ctx context.Context, name string) (provider.Grantee, error) {
	log := appctx.GetLogger(ctx)
	client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		return provider.Grantee{}, err
	}
	userRes, err := client.GetUserByClaim(ctx, &userpb.GetUserByClaimRequest{
		Claim: "username",
		Value: name,
	})
	if err == nil && userRes.Status.Code == rpc.Code_CODE_OK {
		return provider.Grantee{
			Type: provider.GranteeType_GRANTEE_TYPE_USER,
			Id:   &provider.Grantee_UserId{UserId: userRes.User.Id},
		}, nil
	}
	log.Debug().Str("name", name).Msg("no user found")

	groupRes, err := client.GetGroupByClaim(ctx, &groupv1beta1.GetGroupByClaimRequest{
		Claim:               "group_name",
		Value:               name,
		SkipFetchingMembers: true,
	})
	if err == nil && groupRes.Status.Code == rpc.Code_CODE_OK {
		return provider.Grantee{
			Type: provider.GranteeType_GRANTEE_TYPE_GROUP,
			Id:   &provider.Grantee_GroupId{GroupId: groupRes.Group.Id},
		}, nil
	}
	log.Debug().Str("name", name).Msg("no group found")

	return provider.Grantee{}, fmt.Errorf("no grantee found with name %s", name)
}

func (h *Handler) addSpaceMember(w http.ResponseWriter, r *http.Request, info *provider.ResourceInfo, role *conversions.Role, roleVal []byte) {
	ctx := r.Context()

	if info.Space.SpaceType == "personal" {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "can not add members to personal spaces", nil)
		return
	}

	shareWith := r.FormValue("shareWith")
	if shareWith == "" {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "missing shareWith", nil)
		return
	}

	grantee, err := h.getGrantee(ctx, shareWith)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "error getting grantee", err)
		return
	}

	client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "error getting gateway client", err)
		return
	}

	permissions := role.CS3ResourcePermissions()
	// All members of a space should be able to list shares inside that space.
	// The viewer role doesn't have the ListGrants permission so we set it here.
	permissions.ListGrants = true

	expireDate := r.PostFormValue("expireDate")
	var expirationTs *types.Timestamp
	if expireDate != "" {
		expiration, err := time.Parse(_iso8601, expireDate)
		if err != nil {
			// Web sends different formats when adding and when editing a space membership...
			// We need to fix this in a separate PR.
			expiration, err = time.Parse(time.RFC3339, expireDate)
			if err != nil {
				response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "could not parse expireDate", err)
				return
			}
		}
		expirationTs = &types.Timestamp{
			Seconds: uint64(expiration.UnixNano() / int64(time.Second)),
			Nanos:   uint32(expiration.UnixNano() % int64(time.Second)),
		}
	}

	if role.Name != conversions.RoleManager {
		ref := provider.Reference{ResourceId: info.GetId()}
		p, err := h.findProvider(ctx, &ref)
		if err != nil {
			response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "error getting storage provider", err)
			return
		}

		providerClient, err := h.getStorageProviderClient(p)
		if err != nil {
			response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "error getting storage provider client", err)
			return
		}

		lgRes, err := providerClient.ListGrants(ctx, &provider.ListGrantsRequest{Ref: &ref})
		if err != nil || lgRes.Status.Code != rpc.Code_CODE_OK {
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error listing space grants", err)
			return
		}

		if !isSpaceManagerRemaining(lgRes.Grants, grantee) {
			response.WriteOCSError(w, r, http.StatusForbidden, "the space must have at least one manager", nil)
			return
		}
	}

	createShareRes, err := client.CreateShare(ctx, &collaborationv1beta1.CreateShareRequest{
		ResourceInfo: info,
		Grant: &collaborationv1beta1.ShareGrant{
			Permissions: &collaborationv1beta1.SharePermissions{
				Permissions: permissions,
			},
			Grantee:    &grantee,
			Expiration: expirationTs,
		},
	})
	if err != nil || createShareRes.Status.Code != rpc.Code_CODE_OK {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "could not add space member", err)
		return
	}

	response.WriteOCSSuccess(w, r, nil)
}

func (h *Handler) isSpaceShare(r *http.Request, spaceID string) (*registry.ProviderInfo, bool) {
	ref, err := storagespace.ParseReference(spaceID)
	if err != nil {
		return nil, false
	}

	if ref.ResourceId.OpaqueId == "" {
		ref.ResourceId.OpaqueId = ref.ResourceId.SpaceId
	}

	p, err := h.findProvider(r.Context(), &ref)
	return p, err == nil
}

func (h *Handler) removeSpaceMember(w http.ResponseWriter, r *http.Request, spaceID string, prov *registry.ProviderInfo) {
	ctx := r.Context()

	shareWith := r.URL.Query().Get("shareWith")
	if shareWith == "" {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "missing shareWith", nil)
		return
	}

	grantee, err := h.getGrantee(ctx, shareWith)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "error getting grantee", err)
		return
	}

	ref, err := storagespace.ParseReference(spaceID)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "could not parse space id", err)
		return
	}

	if ref.ResourceId.OpaqueId == "" {
		ref.ResourceId.OpaqueId = ref.ResourceId.SpaceId
	}

	providerClient, err := h.getStorageProviderClient(prov)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "error getting storage provider client", err)
		return
	}

	lgRes, err := providerClient.ListGrants(ctx, &provider.ListGrantsRequest{Ref: &ref})
	if err != nil || lgRes.Status.Code != rpc.Code_CODE_OK {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error listing space grants", err)
		return
	}

	if len(lgRes.Grants) == 1 || !isSpaceManagerRemaining(lgRes.Grants, grantee) {
		response.WriteOCSError(w, r, http.StatusForbidden, "can't remove the last manager", nil)
		return
	}

	gatewayClient, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "error getting gateway client", err)
		return
	}

	removeShareRes, err := gatewayClient.RemoveShare(ctx, &collaborationv1beta1.RemoveShareRequest{
		Ref: &collaborationv1beta1.ShareReference{
			Spec: &collaborationv1beta1.ShareReference_Key{
				Key: &collaborationv1beta1.ShareKey{
					ResourceId: ref.ResourceId,
					Grantee:    &grantee,
				},
			},
		},
	})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error removing grant", err)
		return
	}
	if removeShareRes.Status.Code != rpc.Code_CODE_OK {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error removing grant", err)
		return
	}

	response.WriteOCSSuccess(w, r, nil)
}

func (h *Handler) getStorageProviderClient(p *registry.ProviderInfo) (provider.ProviderAPIClient, error) {
	c, err := pool.GetStorageProviderServiceClient(p.Address)
	if err != nil {
		err = errors.Wrap(err, "shares spaces: error getting a storage provider client")
		return nil, err
	}

	return c, nil
}

func (h *Handler) findProvider(ctx context.Context, ref *provider.Reference) (*registry.ProviderInfo, error) {
	c, err := pool.GetStorageRegistryClient(h.storageRegistryAddr)
	if err != nil {
		return nil, errors.Wrap(err, "shares spaces: error getting storage registry client")
	}

	filters := map[string]string{}
	if ref.Path != "" {
		filters["path"] = ref.Path
	}
	if ref.ResourceId != nil {
		filters["storage_id"] = ref.ResourceId.StorageId
		filters["space_id"] = ref.ResourceId.SpaceId
		filters["opaque_id"] = ref.ResourceId.OpaqueId
	}

	listReq := &registry.ListStorageProvidersRequest{
		Opaque: &types.Opaque{},
	}
	sdk.EncodeOpaqueMap(listReq.Opaque, filters)

	res, err := c.ListStorageProviders(ctx, listReq)

	if err != nil {
		return nil, errors.Wrap(err, "shares spaces: error calling ListStorageProviders")
	}

	if res.Status.Code != rpc.Code_CODE_OK {
		switch res.Status.Code {
		case rpc.Code_CODE_NOT_FOUND:
			return nil, errtypes.NotFound("shares spaces: storage provider not found for reference:" + ref.String())
		case rpc.Code_CODE_PERMISSION_DENIED:
			return nil, errtypes.PermissionDenied("shares spaces: " + res.Status.Message + " for " + ref.String() + " with code " + res.Status.Code.String())
		case rpc.Code_CODE_INVALID_ARGUMENT, rpc.Code_CODE_FAILED_PRECONDITION, rpc.Code_CODE_OUT_OF_RANGE:
			return nil, errtypes.BadRequest("shares spaces: " + res.Status.Message + " for " + ref.String() + " with code " + res.Status.Code.String())
		case rpc.Code_CODE_UNIMPLEMENTED:
			return nil, errtypes.NotSupported("shares spaces: " + res.Status.Message + " for " + ref.String() + " with code " + res.Status.Code.String())
		default:
			return nil, status.NewErrorFromCode(res.Status.Code, "shares spaces")
		}
	}

	if len(res.Providers) < 1 {
		return nil, errtypes.NotFound("shares spaces: no provider found")
	}

	return res.Providers[0], nil
}

func isSpaceManagerRemaining(grants []*provider.Grant, grantee provider.Grantee) bool {
	for _, g := range grants {
		// RemoveGrant is currently the way to check for the manager role
		// If it is not set than the current grant is not for a manager and
		// we can just continue with the next one.
		if g.Permissions.RemoveGrant && !isEqualGrantee(*g.Grantee, grantee) {
			return true
		}
	}
	return false
}

func isEqualGrantee(a, b provider.Grantee) bool {
	// Ideally we would want to use utils.GranteeEqual()
	// but the grants stored in the decomposedfs aren't complete (missing usertype and idp)
	// because of that the check would fail so we can only check the ... for now.
	if a.Type != b.Type {
		return false
	}

	var aID, bID string
	switch a.Type {
	case provider.GranteeType_GRANTEE_TYPE_GROUP:
		aID = a.GetGroupId().OpaqueId
		bID = b.GetGroupId().OpaqueId
	case provider.GranteeType_GRANTEE_TYPE_USER:
		aID = a.GetUserId().OpaqueId
		bID = b.GetUserId().OpaqueId
	}
	return aID == bID
}
