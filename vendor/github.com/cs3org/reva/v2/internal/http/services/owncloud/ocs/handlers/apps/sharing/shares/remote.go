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
	"net/http"
	"strconv"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	invitepb "github.com/cs3org/go-cs3apis/cs3/ocm/invite/v1beta1"
	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/go-chi/chi/v5"

	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/response"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
)

func (h *Handler) createFederatedCloudShare(w http.ResponseWriter, r *http.Request, statInfo *provider.ResourceInfo, role *conversions.Role, roleVal []byte) {
	ctx := r.Context()

	selector, err := pool.GatewaySelector(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting gateway selector", nil)
		return
	}
	c, err := selector.Next()
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error selecting next client", nil)
		return
	}

	shareWithUser, shareWithProvider := r.FormValue("shareWithUser"), r.FormValue("shareWithProvider")
	if shareWithUser == "" || shareWithProvider == "" {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "missing shareWith parameters", nil)
		return
	}

	providerInfoResp, err := c.GetInfoByDomain(ctx, &ocmprovider.GetInfoByDomainRequest{
		Domain: shareWithProvider,
	})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc get invite by domain info request", err)
		return
	}

	remoteUserRes, err := c.GetAcceptedUser(ctx, &invitepb.GetAcceptedUserRequest{
		RemoteUserId: &userpb.UserId{OpaqueId: shareWithUser, Idp: shareWithProvider, Type: userpb.UserType_USER_TYPE_PRIMARY},
	})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error searching recipient", err)
		return
	}
	if remoteUserRes.Status.Code != rpc.Code_CODE_OK {
		response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "user not found", err)
		return
	}

	createShareReq := &ocm.CreateOCMShareRequest{
		Opaque: &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				/* TODO extend the spec with role names?
				"role": {
					Decoder: "plain",
					Value:   []byte(role.Name),
				},
				*/
				"permissions": {
					Decoder: "plain",
					Value:   []byte(strconv.Itoa(int(role.OCSPermissions()))),
				},
				"name": {
					Decoder: "plain",
					Value:   []byte(statInfo.Path),
				},
			},
		},
		ResourceId: statInfo.Id,
		Grant: &ocm.ShareGrant{
			Grantee: &provider.Grantee{
				Type: provider.GranteeType_GRANTEE_TYPE_USER,
				Id:   &provider.Grantee_UserId{UserId: remoteUserRes.RemoteUser.GetId()},
			},
			Permissions: &ocm.SharePermissions{
				Permissions: role.CS3ResourcePermissions(),
			},
		},
		RecipientMeshProvider: providerInfoResp.ProviderInfo,
	}

	createShareResponse, err := c.CreateOCMShare(ctx, createShareReq)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc create ocm share request", err)
		return
	}
	if createShareResponse.Status.Code != rpc.Code_CODE_OK {
		if createShareResponse.Status.Code == rpc.Code_CODE_NOT_FOUND {
			response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "not found", nil)
			return
		}
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc create ocm share request failed", err)
		return
	}

	response.WriteOCSSuccess(w, r, "OCM Share created")
}

// GetFederatedShare handles GET requests on /apps/files_sharing/api/v1/shares/remote_shares/{shareid}
func (h *Handler) GetFederatedShare(w http.ResponseWriter, r *http.Request) {

	// TODO: Implement response with HAL schemating
	ctx := r.Context()

	shareID := chi.URLParam(r, "shareid")
	selector, err := pool.GatewaySelector(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting gateway selector", nil)
		return
	}
	gatewayClient, err := selector.Next()
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error selecting next client", nil)
		return
	}

	listOCMSharesRequest := &ocm.GetOCMShareRequest{
		Ref: &ocm.ShareReference{
			Spec: &ocm.ShareReference_Id{
				Id: &ocm.ShareId{
					OpaqueId: shareID,
				},
			},
		},
	}
	ocmShareResponse, err := gatewayClient.GetOCMShare(ctx, listOCMSharesRequest)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc get ocm share request", err)
		return
	}

	share := ocmShareResponse.GetShare()
	if share == nil {
		response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "share not found", err)
		return
	}
	response.WriteOCSSuccess(w, r, share)
}

// ListFederatedShares handles GET requests on /apps/files_sharing/api/v1/shares/remote_shares
func (h *Handler) ListFederatedShares(w http.ResponseWriter, r *http.Request) {

	// TODO Implement pagination.
	// TODO Implement response with HAL schemating
	ctx := r.Context()

	selector, err := pool.GatewaySelector(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting gateway selector", nil)
		return
	}
	gatewayClient, err := selector.Next()
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error selecting next client", nil)
		return
	}

	listOCMSharesResponse, err := gatewayClient.ListOCMShares(ctx, &ocm.ListOCMSharesRequest{})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc list ocm share request", err)
		return
	}

	shares := listOCMSharesResponse.GetShares()
	if shares == nil {
		shares = make([]*ocm.Share, 0)
	}
	response.WriteOCSSuccess(w, r, shares)
}
