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
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/response"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
)

const (
	_iso8601 = "2006-01-02T15:04:05Z0700"
)

func (h *Handler) createUserShare(w http.ResponseWriter, r *http.Request, statInfo *provider.ResourceInfo, role *conversions.Role, roleVal []byte) (*collaboration.Share, *ocsError) {
	ctx := r.Context()
	c, err := h.getClient()
	if err != nil {
		return nil, &ocsError{
			Code:    response.MetaServerError.StatusCode,
			Message: "error getting grpc gateway client",
			Error:   err,
		}
	}

	shareWith := r.FormValue("shareWith")
	if shareWith == "" {
		return nil, &ocsError{
			Code:    response.MetaBadRequest.StatusCode,
			Message: "missing shareWith",
			Error:   err,
		}
	}

	userRes, err := c.GetUserByClaim(ctx, &userpb.GetUserByClaimRequest{
		Claim:                  "username",
		Value:                  shareWith,
		SkipFetchingUserGroups: true,
	})
	if err != nil {
		return nil, &ocsError{
			Code:    response.MetaServerError.StatusCode,
			Message: "error searching recipient",
			Error:   err,
		}
	}

	if userRes.Status.Code != rpc.Code_CODE_OK {
		return nil, &ocsError{
			Code:    response.MetaNotFound.StatusCode,
			Message: "user not found",
			Error:   err,
		}
	}

	expireDate := r.PostFormValue("expireDate")
	var expirationTs *types.Timestamp
	if expireDate != "" {
		// FIXME: the web ui sends the RFC3339 format when updating a share but
		// initially on creating a share the format ISO 8601 is used.
		// OC10 uses RFC3339 in both cases so we should fix the web ui and change it here.
		expiration, err := time.Parse(_iso8601, expireDate)
		if err != nil {
			return nil, &ocsError{
				Code:    response.MetaBadRequest.StatusCode,
				Message: "could not parse expireDate",
				Error:   err,
			}
		}
		expirationTs = &types.Timestamp{
			Seconds: uint64(expiration.UnixNano() / int64(time.Second)),
			Nanos:   uint32(expiration.UnixNano() % int64(time.Second)),
		}
	}

	createShareReq := &collaboration.CreateShareRequest{
		Opaque: &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"role": {
					Decoder: "json",
					Value:   roleVal,
				},
			},
		},
		ResourceInfo: statInfo,
		Grant: &collaboration.ShareGrant{
			Grantee: &provider.Grantee{
				Type: provider.GranteeType_GRANTEE_TYPE_USER,
				Id:   &provider.Grantee_UserId{UserId: userRes.User.GetId()},
			},
			Permissions: &collaboration.SharePermissions{
				Permissions: role.CS3ResourcePermissions(),
			},
			Expiration: expirationTs,
		},
	}

	share, ocsErr := h.createCs3Share(ctx, w, r, c, createShareReq)
	if ocsErr != nil {
		return nil, ocsErr
	}
	return share, nil
}

func (h *Handler) isUserShare(r *http.Request, oid string) (*collaboration.Share, bool) {
	logger := appctx.GetLogger(r.Context())
	client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		logger.Err(err)
	}

	getShareRes, err := client.GetShare(r.Context(), &collaboration.GetShareRequest{
		Ref: &collaboration.ShareReference{
			Spec: &collaboration.ShareReference_Id{
				Id: &collaboration.ShareId{
					OpaqueId: oid,
				},
			},
		},
	})
	if err != nil {
		logger.Err(err)
		return nil, false
	}

	return getShareRes.GetShare(), getShareRes.GetShare() != nil
}

func (h *Handler) removeUserShare(w http.ResponseWriter, r *http.Request, share *collaboration.Share) {
	ctx := r.Context()

	uClient, err := h.getClient()
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting grpc gateway client", err)
		return
	}

	shareRef := &collaboration.ShareReference{
		Spec: &collaboration.ShareReference_Id{
			Id: share.Id,
		},
	}

	data, err := conversions.CS3Share2ShareData(ctx, share)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "deleting share failed", err)
		return
	}
	// A deleted share should not have an ID.
	data.ID = ""

	uReq := &collaboration.RemoveShareRequest{Ref: shareRef}
	uRes, err := uClient.RemoveShare(ctx, uReq)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc delete share request", err)
		return
	}

	if uRes.Status.Code != rpc.Code_CODE_OK {
		if uRes.Status.Code == rpc.Code_CODE_NOT_FOUND {
			response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "not found", nil)
			return
		}
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc delete share request failed", err)
		return
	}
	if currentUser, ok := ctxpkg.ContextGetUser(ctx); ok {
		h.statCache.RemoveStat(currentUser.Id, share.ResourceId)
	}
	response.WriteOCSSuccess(w, r, data)
}

func (h *Handler) listUserShares(r *http.Request, filters []*collaboration.Filter) ([]*conversions.ShareData, *rpc.Status, error) {
	ctx := r.Context()
	log := appctx.GetLogger(ctx)

	lsUserSharesRequest := collaboration.ListSharesRequest{
		Filters: filters,
	}

	ocsDataPayload := make([]*conversions.ShareData, 0)
	if h.gatewayAddr != "" {
		// get a connection to the users share provider
		client, err := h.getClient()
		if err != nil {
			return ocsDataPayload, nil, err
		}

		// do list shares request. filtered
		lsUserSharesResponse, err := client.ListShares(ctx, &lsUserSharesRequest)
		if err != nil {
			return ocsDataPayload, nil, err
		}
		if lsUserSharesResponse.Status.Code != rpc.Code_CODE_OK {
			return ocsDataPayload, lsUserSharesResponse.Status, nil
		}

		// build OCS response payload
		for _, s := range lsUserSharesResponse.Shares {
			data, err := conversions.CS3Share2ShareData(ctx, s)
			if err != nil {
				log.Debug().Interface("share", s).Interface("shareData", data).Err(err).Msg("could not CS3Share2ShareData, skipping")
				continue
			}

			info, status, err := h.getResourceInfoByID(ctx, client, s.ResourceId)
			if err != nil || status.Code != rpc.Code_CODE_OK {
				log.Debug().Interface("share", s).Interface("status", status).Interface("shareData", data).Err(err).Msg("could not stat share, skipping")
				continue
			}
			u := ctxpkg.ContextMustGetUser(ctx)
			// check if the user has the permission to list all shares on the resource
			if !utils.UserEqual(s.Creator, u.Id) && !info.GetPermissionSet().ListGrants {
				log.Debug().Interface("share", s).Interface("user", u).Msg("user has no permission to list all grants and is not the creator of this share")
				continue
			}

			h.addFileInfo(ctx, data, info)
			h.mapUserIds(ctx, client, data)

			// Filter out a share if ShareWith is not found because the user or group already deleted
			if data.ShareWith == "" {
				continue
			}

			log.Debug().Interface("share", s).Interface("info", info).Interface("shareData", data).Msg("mapped")
			ocsDataPayload = append(ocsDataPayload, data)
		}
	}

	return ocsDataPayload, nil, nil
}
