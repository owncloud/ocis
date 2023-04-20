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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	permissionsv1beta1 "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/rs/zerolog/log"

	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/response"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/publicshare"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/pkg/errors"
)

func (h *Handler) createPublicLinkShare(w http.ResponseWriter, r *http.Request, statInfo *provider.ResourceInfo) (*link.PublicShare, *ocsError) {
	ctx := r.Context()
	log := appctx.GetLogger(ctx)

	c, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		return nil, &ocsError{
			Code:    response.MetaServerError.StatusCode,
			Message: "error getting grpc gateway client",
			Error:   err,
		}
	}

	user := ctxpkg.ContextMustGetUser(ctx)
	resp, err := c.CheckPermission(ctx, &permissionsv1beta1.CheckPermissionRequest{
		SubjectRef: &permissionsv1beta1.SubjectReference{
			Spec: &permissionsv1beta1.SubjectReference_UserId{
				UserId: user.Id,
			},
		},
		Permission: "PublicLink.Write",
	})
	if err != nil {
		return nil, &ocsError{
			Code:    response.MetaServerError.StatusCode,
			Message: "failed to check user permission",
			Error:   err,
		}
	}

	if resp.Status.Code != rpc.Code_CODE_OK {
		return nil, &ocsError{
			Code:    response.MetaForbidden.StatusCode,
			Message: "user is not allowed to create a public link",
		}
	}

	err = r.ParseForm()
	if err != nil {
		return nil, &ocsError{
			Code:    response.MetaBadRequest.StatusCode,
			Message: "Could not parse form from request",
			Error:   err,
		}
	}

	// check if a quicklink should be created
	quick, _ := strconv.ParseBool(r.FormValue("quicklink")) // no need to check the error - defaults to zero value!
	if quick {
		f := []*link.ListPublicSharesRequest_Filter{publicshare.ResourceIDFilter(statInfo.Id)}
		req := link.ListPublicSharesRequest{Filters: f}
		res, err := c.ListPublicShares(ctx, &req)
		if err != nil {
			return nil, &ocsError{
				Code:    response.MetaServerError.StatusCode,
				Message: "could not list public links",
				Error:   err,
			}
		}
		if res.Status.Code != rpc.Code_CODE_OK {
			return nil, &ocsError{
				Code:    int(res.Status.GetCode()),
				Message: "could not list public links",
			}
		}

		for _, l := range res.GetShare() {
			if l.Quicklink {
				return l, nil
			}
		}
	}

	newPermissions, err := permissionFromRequest(r, h)
	if err != nil {
		return nil, &ocsError{
			Code:    response.MetaBadRequest.StatusCode,
			Message: "Could not read permission from request",
			Error:   err,
		}
	}

	if newPermissions == nil {
		// default perms: read-only
		// TODO: the default might change depending on allowed permissions and configs
		newPermissions, err = ocPublicPermToCs3(1, h)
		if err != nil {
			return nil, &ocsError{
				Code:    response.MetaServerError.StatusCode,
				Message: "Could not convert default permissions",
				Error:   err,
			}
		}
	}

	if statInfo != nil && statInfo.Type == provider.ResourceType_RESOURCE_TYPE_FILE {
		// Single file shares should never have delete or create permissions
		role := conversions.RoleFromResourcePermissions(newPermissions, true)
		permissions := role.OCSPermissions()
		permissions &^= conversions.PermissionCreate
		permissions &^= conversions.PermissionDelete
		newPermissions = conversions.RoleFromOCSPermissions(permissions).CS3ResourcePermissions()
	}

	if !sufficientPermissions(statInfo.PermissionSet, newPermissions, true) {
		response.WriteOCSError(w, r, http.StatusNotFound, "no share permission", nil)
		return nil, &ocsError{
			Code:    http.StatusNotFound,
			Message: "Cannot set the requested share permissions",
			Error:   errors.New("cannot set the requested share permissions"),
		}
	}

	req := link.CreatePublicShareRequest{
		ResourceInfo: statInfo,
		Grant: &link.Grant{
			Permissions: &link.PublicSharePermissions{
				Permissions: newPermissions,
			},
			Password: r.FormValue("password"),
		},
	}

	expireTimeString, ok := r.Form["expireDate"]
	if ok {
		if expireTimeString[0] != "" {
			expireTime, err := conversions.ParseTimestamp(expireTimeString[0])
			if err != nil {
				return nil, &ocsError{
					Code:    response.MetaBadRequest.StatusCode,
					Message: err.Error(),
					Error:   err,
				}
			}
			if expireTime != nil {
				req.Grant.Expiration = expireTime
			}
		}
	}

	// set displayname and password protected as arbitrary metadata
	req.ResourceInfo.ArbitraryMetadata = &provider.ArbitraryMetadata{
		Metadata: map[string]string{
			"name":      r.FormValue("name"),
			"quicklink": r.FormValue("quicklink"),
			// "password": r.FormValue("password"),
		},
	}

	createRes, err := c.CreatePublicShare(ctx, &req)
	if err != nil {
		log.Debug().Err(err).Str("createShare", "shares").Msgf("error creating a public share to resource id: %v", statInfo.GetId())
		return nil, &ocsError{
			Code:    response.MetaServerError.StatusCode,
			Message: "error creating public share",
			Error:   fmt.Errorf("error creating a public share to resource id: %v", statInfo.GetId()),
		}
	}

	if createRes.Status.Code != rpc.Code_CODE_OK {
		log.Debug().Err(errors.New("create public share failed")).Str("shares", "createShare").Msgf("create public share failed with status code: %v", createRes.Status.Code.String())
		return nil, &ocsError{
			Code:    response.MetaServerError.StatusCode,
			Message: "grpc create public share request failed",
			Error:   nil,
		}
	}
	return createRes.Share, nil
}

func (h *Handler) listPublicShares(r *http.Request, filters []*link.ListPublicSharesRequest_Filter) ([]*conversions.ShareData, *rpc.Status, error) {
	ctx := r.Context()
	log := appctx.GetLogger(ctx)

	ocsDataPayload := make([]*conversions.ShareData, 0)
	// TODO(refs) why is this guard needed? Are we moving towards a gateway only for service discovery? without a gateway this is dead code.
	if h.gatewayAddr != "" {
		client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
		if err != nil {
			return ocsDataPayload, nil, err
		}

		req := link.ListPublicSharesRequest{
			Filters: filters,
		}

		res, err := client.ListPublicShares(ctx, &req)
		if err != nil {
			return ocsDataPayload, nil, err
		}
		if res.Status.Code != rpc.Code_CODE_OK {
			return ocsDataPayload, res.Status, nil
		}

		for _, share := range res.GetShare() {
			info, status, err := h.getResourceInfoByID(ctx, client, share.ResourceId)
			if err != nil || status.Code != rpc.Code_CODE_OK {
				log.Debug().Interface("share", share).Interface("status", status).Err(err).Msg("could not stat share, skipping")
				continue
			}

			sData := conversions.PublicShare2ShareData(share, r, h.publicURL)

			sData.Name = share.DisplayName

			h.addFileInfo(ctx, sData, info)
			h.mapUserIds(ctx, client, sData)

			log.Debug().Interface("share", share).Interface("info", info).Interface("shareData", share).Msg("mapped")

			ocsDataPayload = append(ocsDataPayload, sData)

		}

		return ocsDataPayload, nil, nil
	}

	return ocsDataPayload, nil, errors.New("bad request")
}

func (h *Handler) isPublicShare(r *http.Request, oid string) bool {
	logger := appctx.GetLogger(r.Context())
	client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		logger.Err(err)
	}

	psRes, err := client.GetPublicShare(r.Context(), &link.GetPublicShareRequest{
		Ref: &link.PublicShareReference{
			Spec: &link.PublicShareReference_Id{
				Id: &link.PublicShareId{
					OpaqueId: oid,
				},
			},
		},
	})
	if err != nil {
		logger.Err(err)
		return false
	}

	return psRes.GetShare() != nil
}

func (h *Handler) updatePublicShare(w http.ResponseWriter, r *http.Request, shareID string) {
	updates := []*link.UpdatePublicShareRequest_Update{}
	logger := appctx.GetLogger(r.Context())

	gwC, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		log.Err(err).Str("shareID", shareID).Msg("updatePublicShare")
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "error getting a connection to the gateway service", nil)
		return
	}

	ctx := r.Context()

	user := ctxpkg.ContextMustGetUser(ctx)
	resp, err := gwC.CheckPermission(ctx, &permissionsv1beta1.CheckPermissionRequest{
		SubjectRef: &permissionsv1beta1.SubjectReference{
			Spec: &permissionsv1beta1.SubjectReference_UserId{
				UserId: user.Id,
			},
		},
		Permission: "PublicLink.Write",
	})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "failed to check user permission", err)
		return
	}

	if resp.Status.Code != rpc.Code_CODE_OK {
		response.WriteOCSError(w, r, response.MetaForbidden.StatusCode, "user is not allowed to create a public link", nil)
		return
	}

	before, err := gwC.GetPublicShare(r.Context(), &link.GetPublicShareRequest{
		Ref: &link.PublicShareReference{
			Spec: &link.PublicShareReference_Id{
				Id: &link.PublicShareId{
					OpaqueId: shareID,
				},
			},
		},
	})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "failed to get public share", nil)
		return
	}

	u := ctxpkg.ContextMustGetUser(r.Context())
	if !publicshare.IsCreatedByUser(*before.Share, u) {
		sRes, err := gwC.Stat(r.Context(), &provider.StatRequest{Ref: &provider.Reference{ResourceId: before.Share.ResourceId}})
		if err != nil {
			log.Err(err).Interface("resource_id", before.Share.ResourceId).Msg("failed to stat shared resource")
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "failed to get public share", nil)
			return
		}

		if !sRes.Info.PermissionSet.UpdateGrant {
			response.WriteOCSError(w, r, response.MetaUnauthorized.StatusCode, "missing permissions to update share", err)
			return
		}
	}

	err = r.ParseForm()
	if err != nil {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "Could not parse form from request", err)
		return
	}

	// indicates whether values to update were found,
	// to check if the request was valid,
	// not whether an actual update has been performed
	updatesFound := false

	newName, ok := r.Form["name"]
	if ok {
		updatesFound = true
		if newName[0] != before.Share.DisplayName {
			updates = append(updates, &link.UpdatePublicShareRequest_Update{
				Type:        link.UpdatePublicShareRequest_Update_TYPE_DISPLAYNAME,
				DisplayName: newName[0],
			})
		}
	}

	// Permissions
	newPermissions, err := permissionFromRequest(r, h)
	logger.Debug().Interface("newPermissions", newPermissions).Msg("Parsed permissions")
	if err != nil {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "invalid permissions", err)
		return
	}

	// update permissions if given
	if newPermissions != nil {
		updatesFound = true
		publicSharePermissions := &link.PublicSharePermissions{
			Permissions: newPermissions,
		}
		beforePerm, _ := json.Marshal(before.GetShare().Permissions)
		afterPerm, _ := json.Marshal(publicSharePermissions)
		if string(beforePerm) != string(afterPerm) {
			logger.Info().Str("shares", "update").Msgf("updating permissions from %v to: %v", string(beforePerm), string(afterPerm))
			updates = append(updates, &link.UpdatePublicShareRequest_Update{
				Type: link.UpdatePublicShareRequest_Update_TYPE_PERMISSIONS,
				Grant: &link.Grant{
					Permissions: publicSharePermissions,
					Password:    r.FormValue("password"),
				},
			})
		}
	}

	statReq := provider.StatRequest{Ref: &provider.Reference{ResourceId: before.Share.ResourceId}}
	statRes, err := gwC.Stat(r.Context(), &statReq)
	if err != nil {
		log.Debug().Err(err).Str("shares", "update public share").Msg("error during stat")
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "missing resource information", fmt.Errorf("error getting resource information"))
		return
	}

	// empty permissions mean internal link here - NOT denial. Hence we need an extra check
	if !sufficientPermissions(statRes.Info.PermissionSet, newPermissions, true) {
		response.WriteOCSError(w, r, http.StatusNotFound, "no share permission", nil)
		return
	}

	// ExpireDate
	expireTimeString, ok := r.Form["expireDate"]
	// check if value is set and must be updated or cleared
	if ok {
		updatesFound = true
		var newExpiration *types.Timestamp
		if expireTimeString[0] != "" {
			newExpiration, err = conversions.ParseTimestamp(expireTimeString[0])
			if err != nil {
				response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "invalid datetime format", err)
				return
			}
		}

		beforeExpiration, _ := json.Marshal(before.Share.Expiration)
		afterExpiration, _ := json.Marshal(newExpiration)
		if string(afterExpiration) != string(beforeExpiration) {
			logger.Debug().Str("shares", "update").Msgf("updating expire date from %v to: %v", string(beforeExpiration), string(afterExpiration))
			updates = append(updates, &link.UpdatePublicShareRequest_Update{
				Type: link.UpdatePublicShareRequest_Update_TYPE_EXPIRATION,
				Grant: &link.Grant{
					Expiration: newExpiration,
				},
			})
		}
	}

	// Password
	newPassword, ok := r.Form["password"]
	// update or clear password
	if ok {
		updatesFound = true
		logger.Info().Str("shares", "update").Msg("password updated")
		updates = append(updates, &link.UpdatePublicShareRequest_Update{
			Type: link.UpdatePublicShareRequest_Update_TYPE_PASSWORD,
			Grant: &link.Grant{
				Password: newPassword[0],
			},
		})
	}

	publicShare := before.Share

	// Updates are atomical. See: https://github.com/cs3org/cs3apis/pull/67#issuecomment-617651428 so in order to get the latest updated version
	if len(updates) > 0 {
		uRes := &link.UpdatePublicShareResponse{Share: before.Share}
		for k := range updates {
			uRes, err = gwC.UpdatePublicShare(r.Context(), &link.UpdatePublicShareRequest{
				Ref: &link.PublicShareReference{
					Spec: &link.PublicShareReference_Id{
						Id: &link.PublicShareId{
							OpaqueId: shareID,
						},
					},
				},
				Update: updates[k],
			})
			if err != nil {
				log.Err(err).Str("shareID", shareID).Msg("sending update request to public link provider")
				response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "Error sending update request to public link provider", err)
				return
			}
			if uRes.Status.Code != rpc.Code_CODE_OK {
				log.Debug().Str("shareID", shareID).Msgf("sending update request to public link provider failed: %s", uRes.Status.Message)
				response.WriteOCSError(w, r, response.MetaServerError.StatusCode, fmt.Sprintf("Error sending update request to public link provider: %s", uRes.Status.Message), nil)
				return
			}
		}
		publicShare = uRes.Share
	} else if !updatesFound {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "No updates specified in request", nil)
		return
	}

	s := conversions.PublicShare2ShareData(publicShare, r, h.publicURL)
	h.addFileInfo(r.Context(), s, statRes.Info)
	h.mapUserIds(r.Context(), gwC, s)

	response.WriteOCSSuccess(w, r, s)
}

func (h *Handler) removePublicShare(w http.ResponseWriter, r *http.Request, shareID string) {
	ctx := r.Context()

	c, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting grpc gateway client", err)
		return
	}

	before, err := c.GetPublicShare(r.Context(), &link.GetPublicShareRequest{
		Ref: &link.PublicShareReference{
			Spec: &link.PublicShareReference_Id{
				Id: &link.PublicShareId{
					OpaqueId: shareID,
				},
			},
		},
	})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "failed to get public share", nil)
		return
	}

	u := ctxpkg.ContextMustGetUser(ctx)
	if !publicshare.IsCreatedByUser(*before.Share, u) {
		sRes, err := c.Stat(r.Context(), &provider.StatRequest{Ref: &provider.Reference{ResourceId: before.Share.ResourceId}})
		if err != nil {
			log.Err(err).Interface("resource_id", before.Share.ResourceId).Msg("failed to stat shared resource")
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "failed to get public share", nil)
			return
		}

		if !sRes.Info.PermissionSet.RemoveGrant {
			response.WriteOCSError(w, r, response.MetaUnauthorized.StatusCode, "missing permissions to remove share", err)
			return
		}
	}

	req := &link.RemovePublicShareRequest{
		Ref: &link.PublicShareReference{
			Spec: &link.PublicShareReference_Id{
				Id: &link.PublicShareId{
					OpaqueId: shareID,
				},
			},
		},
	}

	res, err := c.RemovePublicShare(ctx, req)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc delete share request", err)
		return
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		if res.Status.Code == rpc.Code_CODE_NOT_FOUND {
			response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "not found", nil)
			return
		}
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc delete share request failed", err)
		return
	}

	response.WriteOCSSuccess(w, r, nil)
}

// for public links oc10 api decreases all permissions to read: stay compatible!
func decreasePermissionsIfNecessary(perm int) int {
	if perm == int(conversions.PermissionAll) {
		perm = int(conversions.PermissionRead)
	}
	return perm
}

func ocPublicPermToCs3(permKey int, h *Handler) (*provider.ResourcePermissions, error) {

	permKey = decreasePermissionsIfNecessary(permKey)

	// TODO refactor this ocPublicPermToRole[permKey] check into a conversions.NewPublicSharePermissions?
	// not all permissions are possible for public shares
	_, ok := ocPublicPermToRole[permKey]
	if !ok {
		log.Error().Str("ocPublicPermToCs3", "shares").Int("perm", permKey).Msg("invalid public share permission")
		return nil, fmt.Errorf("invalid public share permission: %d", permKey)
	}

	perm, err := conversions.NewPermissions(permKey)
	if err != nil && err != conversions.ErrZeroPermission { // we allow empty permissions for public links
		return nil, err
	}

	return conversions.RoleFromOCSPermissions(perm).CS3ResourcePermissions(), nil
}

func permissionFromRequest(r *http.Request, h *Handler) (*provider.ResourcePermissions, error) {
	var err error
	// phoenix sends: {"permissions": 15}. See ocPublicPermToRole struct for mapping

	permKey := 1

	// note: "permissions" value has higher priority than "publicUpload"

	// handle legacy "publicUpload" arg that overrides permissions differently depending on the scenario
	// https://github.com/owncloud/core/blob/v10.4.0/apps/files_sharing/lib/Controller/Share20OcsController.php#L447
	publicUploadString, ok := r.Form["publicUpload"]
	if ok {
		publicUploadFlag, err := strconv.ParseBool(publicUploadString[0])
		if err != nil {
			log.Error().Err(err).Str("publicUpload", publicUploadString[0]).Msg("could not parse publicUpload argument")
			return nil, err
		}

		if publicUploadFlag {
			// all perms except reshare
			permKey = 15
		}
	} else {
		permissionsString, ok := r.Form["permissions"]
		if !ok {
			// no permission values given
			return nil, nil
		}

		permKey, err = strconv.Atoi(permissionsString[0])
		if err != nil {
			log.Error().Str("permissionFromRequest", "shares").Msgf("invalid type: %T", permKey)
			return nil, fmt.Errorf("invalid type: %T", permKey)
		}
	}

	p, err := ocPublicPermToCs3(permKey, h)
	if err != nil {
		return nil, err
	}
	return p, err
}

// TODO: add mapping for user share permissions to role

// Maps oc10 public link permissions to roles
var ocPublicPermToRole = map[int]string{
	// Recipients can do nothing
	0: "none",
	// Recipients can view and download contents.
	1: "viewer",
	// Recipients can view, download and edit single files.
	3: "file-editor",
	// Recipients can view, download, edit, delete and upload contents
	15: "editor",
	// Recipients can upload but existing contents are not revealed
	4: "uploader",
	// Recipients can view, download and upload contents
	5: "contributor",
}
