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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/config"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/response"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/publicshare"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/share"
	"github.com/cs3org/reva/v2/pkg/share/cache"
	cachereg "github.com/cs3org/reva/v2/pkg/share/cache/registry"
	warmupreg "github.com/cs3org/reva/v2/pkg/share/cache/warmup/registry"
	"github.com/cs3org/reva/v2/pkg/storage/utils/templates"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/pkg/errors"
)

const (
	storageIDPrefix string = "shared::"
)

var (
	errParsingSpaceReference = errors.New("could not parse space reference")
)

// Handler implements the shares part of the ownCloud sharing API
type Handler struct {
	gatewayAddr                           string
	machineAuthAPIKey                     string
	storageRegistryAddr                   string
	publicURL                             string
	sharePrefix                           string
	homeNamespace                         string
	skipUpdatingExistingSharesMountpoints bool
	additionalInfoTemplate                *template.Template
	userIdentifierCache                   *ttlcache.Cache
	resourceInfoCache                     cache.ResourceInfoCache
	resourceInfoCacheTTL                  time.Duration

	getClient GatewayClientGetter
}

// we only cache the minimal set of data instead of the full user metadata
type userIdentifiers struct {
	DisplayName string
	Username    string
	Mail        string
}

type ocsError struct {
	Error   error
	Code    int
	Message string
}

func getCacheWarmupManager(c *config.Config) (cache.Warmup, error) {
	if f, ok := warmupreg.NewFuncs[c.CacheWarmupDriver]; ok {
		return f(c.CacheWarmupDrivers[c.CacheWarmupDriver])
	}
	return nil, fmt.Errorf("driver not found: %s", c.CacheWarmupDriver)
}

// GatewayClientGetter is the function being used to retrieve a gateway client instance
type GatewayClientGetter func() (gateway.GatewayAPIClient, error)

func getCacheManager(c *config.Config) (cache.ResourceInfoCache, error) {
	if f, ok := cachereg.NewFuncs[c.ResourceInfoCacheDriver]; ok {
		return f(c.ResourceInfoCacheDrivers[c.ResourceInfoCacheDriver])
	}
	return nil, fmt.Errorf("driver not found: %s", c.ResourceInfoCacheDriver)
}

// Init initializes this and any contained handlers
func (h *Handler) Init(c *config.Config) {
	h.gatewayAddr = c.GatewaySvc
	h.machineAuthAPIKey = c.MachineAuthAPIKey
	h.storageRegistryAddr = c.StorageregistrySvc
	h.publicURL = c.Config.Host
	h.sharePrefix = c.SharePrefix
	h.homeNamespace = c.HomeNamespace
	h.skipUpdatingExistingSharesMountpoints = c.SkipUpdatingExistingSharesMountpoints

	h.additionalInfoTemplate, _ = template.New("additionalInfo").Parse(c.AdditionalInfoAttribute)
	h.resourceInfoCacheTTL = time.Second * time.Duration(c.ResourceInfoCacheTTL)

	h.userIdentifierCache = ttlcache.NewCache()
	_ = h.userIdentifierCache.SetTTL(time.Second * time.Duration(c.UserIdentifierCacheTTL))

	cache, err := getCacheManager(c)
	if err == nil {
		h.resourceInfoCache = cache
	}

	if h.resourceInfoCacheTTL > 0 {
		cwm, err := getCacheWarmupManager(c)
		if err == nil {
			go h.startCacheWarmup(cwm)
		}
	}
	h.getClient = h.getPoolClient
}

// InitWithGetter initializes the handler and adds the clientGetter
func (h *Handler) InitWithGetter(c *config.Config, clientGetter GatewayClientGetter) {
	h.Init(c)
	h.getClient = clientGetter
}

func (h *Handler) startCacheWarmup(c cache.Warmup) {
	time.Sleep(2 * time.Second)
	infos, err := c.GetResourceInfos()
	if err != nil {
		return
	}
	for _, r := range infos {
		key := storagespace.FormatResourceID(*r.Id)
		_ = h.resourceInfoCache.SetWithExpire(key, r, h.resourceInfoCacheTTL)
	}
}

func (h *Handler) extractReference(r *http.Request) (provider.Reference, error) {
	var ref provider.Reference

	// NOTE: space_ref is deprecated and will be removed in ~2 weeks (1.6.22)
	sr := r.FormValue("space_ref")
	if sr != "" {
		return storagespace.ParseReference(sr)
	}

	p, id := r.FormValue("path"), r.FormValue("space")
	if p == "" && id == "" {
		return ref, errors.New("need path or space to extract reference")
	}

	if p != "" {
		u := ctxpkg.ContextMustGetUser(r.Context())
		ref.Path = path.Join(h.getHomeNamespace(u), p)
	}

	if id != "" {
		rid, err := storagespace.ParseID(id)
		if err != nil {
			return ref, err
		}
		ref.ResourceId = &rid
	}

	return ref, nil
}

// CreateShare handles POST requests on /apps/files_sharing/api/v1/shares
func (h *Handler) CreateShare(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shareType, err := strconv.Atoi(r.FormValue("shareType"))
	if err != nil {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "shareType must be an integer", nil)
		return
	}
	// get user permissions on the shared file

	client, err := h.getClient()
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting grpc gateway client", err)
		return
	}
	ref, err := h.extractReference(r)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, errParsingSpaceReference.Error(), errParsingSpaceReference)
		return
	}
	sublog := appctx.GetLogger(ctx).With().Interface("ref", ref).Logger()

	statReq := provider.StatRequest{Ref: &ref}
	statRes, err := client.Stat(ctx, &statReq)
	if err != nil {
		sublog.Debug().Err(err).Msg("CreateShare: error on stat call")
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "missing resource information", fmt.Errorf("error getting resource information"))
		return
	}

	if statRes.Status.Code != rpc.Code_CODE_OK {
		switch statRes.Status.Code {
		case rpc.Code_CODE_NOT_FOUND:
			response.WriteOCSError(w, r, http.StatusNotFound, "Not found", nil)
			w.WriteHeader(http.StatusNotFound)
		case rpc.Code_CODE_PERMISSION_DENIED:
			response.WriteOCSError(w, r, http.StatusNotFound, "No share permission", nil)
			w.WriteHeader(http.StatusForbidden)
		default:
			log.Error().Interface("status", statRes.Status).Msg("CreateShare: stat failed")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// check user has share permissions
	if !conversions.RoleFromResourcePermissions(statRes.Info.PermissionSet).OCSPermissions().Contain(conversions.PermissionShare) {
		response.WriteOCSError(w, r, http.StatusNotFound, "No share permission", nil)
		return
	}

	reqRole, reqPermissions := r.FormValue("role"), r.FormValue("permissions")
	switch shareType {
	case int(conversions.ShareTypeUser), int(conversions.ShareTypeGroup):
		// user collaborations default to Manager (=all permissions)
		role, val, ocsErr := h.extractPermissions(reqRole, reqPermissions, statRes.Info, conversions.NewManagerRole())
		if ocsErr != nil {
			response.WriteOCSError(w, r, ocsErr.Code, ocsErr.Message, ocsErr.Error)
			return
		}

		var share *collaboration.Share
		if shareType == int(conversions.ShareTypeUser) {
			share, ocsErr = h.createUserShare(w, r, statRes.Info, role, val)
		} else {
			share, ocsErr = h.createGroupShare(w, r, statRes.Info, role, val)
		}
		if ocsErr != nil {
			response.WriteOCSError(w, r, ocsErr.Code, ocsErr.Message, ocsErr.Error)
			return
		}

		s, err := conversions.CS3Share2ShareData(ctx, share)
		if err != nil {
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error mapping share data", err)
			return
		}

		err = h.addFileInfo(ctx, s, statRes.Info)
		if err != nil {
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error adding fileinfo to share", err)
			return
		}

		h.mapUserIds(ctx, client, s)

		if !h.skipUpdatingExistingSharesMountpoints {
			status, msg, err := h.updateExistingShareMountpoints(ctx, shareType, share, statRes.Info, client)
			if status != response.MetaOK.StatusCode {
				response.WriteOCSError(w, r, status, msg, err)
				return
			}
		}

		response.WriteOCSSuccess(w, r, s)
	case int(conversions.ShareTypePublicLink):
		// public links default to read only
		_, _, ocsErr := h.extractPermissions(reqRole, reqPermissions, statRes.Info, conversions.NewViewerRole())
		if ocsErr != nil {
			response.WriteOCSError(w, r, http.StatusNotFound, "No share permission", nil)
			return
		}
		share, ocsErr := h.createPublicLinkShare(w, r, statRes.Info)
		if ocsErr != nil {
			response.WriteOCSError(w, r, ocsErr.Code, ocsErr.Message, ocsErr.Error)
			return
		}

		s := conversions.PublicShare2ShareData(share, r, h.publicURL)
		err = h.addFileInfo(ctx, s, statRes.Info)
		if err != nil {
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error enhancing response with share data", err)
			return
		}
		h.mapUserIds(ctx, client, s)

		response.WriteOCSSuccess(w, r, s)
	case int(conversions.ShareTypeFederatedCloudShare):
		// federated shares default to read only
		if role, val, err := h.extractPermissions(reqRole, reqPermissions, statRes.Info, conversions.NewViewerRole()); err == nil {
			h.createFederatedCloudShare(w, r, statRes.Info, role, val)
		}
	case int(conversions.ShareTypeSpaceMembership):
		switch reqRole {
		// Note: we convert viewer and editor roles to spaceviewer and spaceditor to keep backwards compatibility
		// we can remove this switch when this behaviour is no longer wanted.
		case conversions.RoleViewer:
			reqRole = conversions.RoleSpaceViewer
		case conversions.RoleEditor:
			reqRole = conversions.RoleSpaceEditor
		}
		if role, val, err := h.extractPermissions(reqRole, reqPermissions, statRes.Info, conversions.NewSpaceViewerRole()); err == nil {
			switch role.Name {
			case conversions.RoleManager, conversions.RoleSpaceEditor, conversions.RoleSpaceViewer:
				h.addSpaceMember(w, r, statRes.Info, role, val)
			default:
				response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "invalid role for space member", nil)
				return
			}
		}
	default:
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "unknown share type", nil)
	}
}

func (h *Handler) updateExistingShareMountpoints(ctx context.Context, shareType int, share *collaboration.Share, info *provider.ResourceInfo, client gateway.GatewayAPIClient) (int, string, error) {
	if shareType == int(conversions.ShareTypeUser) {
		res, err := client.GetUser(ctx, &userpb.GetUserRequest{
			UserId: &userpb.UserId{
				OpaqueId: share.Grantee.GetUserId().GetOpaqueId(),
			},
		})
		if err != nil {
			return response.MetaServerError.StatusCode, "could not look up user", err
		}
		if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
			return response.MetaServerError.StatusCode, "get user call failed", nil
		}
		if res.User == nil {
			return response.MetaServerError.StatusCode, "grantee not found", nil
		}

		// Get auth
		granteeCtx := ctxpkg.ContextSetUser(context.Background(), res.User)

		authRes, err := client.Authenticate(granteeCtx, &gateway.AuthenticateRequest{
			Type:         "machine",
			ClientId:     res.User.Username,
			ClientSecret: h.machineAuthAPIKey,
		})
		if err != nil {
			return response.MetaServerError.StatusCode, "could not do machine authentication", err
		}
		if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
			return response.MetaServerError.StatusCode, "machine authentication failed", nil
		}
		granteeCtx = metadata.AppendToOutgoingContext(granteeCtx, ctxpkg.TokenHeader, authRes.Token)

		lrs, ocsResponse := getSharesList(granteeCtx, client)
		if ocsResponse != nil {
			return ocsResponse.OCS.Meta.StatusCode, ocsResponse.OCS.Meta.Message, nil
		}

		for _, s := range lrs.Shares {
			if s.GetShare().GetId() != share.Id && s.State == collaboration.ShareState_SHARE_STATE_ACCEPTED && utils.ResourceIDEqual(s.Share.ResourceId, info.GetId()) {
				updateRequest := &collaboration.UpdateReceivedShareRequest{
					Share: &collaboration.ReceivedShare{
						Share:      share,
						MountPoint: s.MountPoint,
						State:      collaboration.ShareState_SHARE_STATE_ACCEPTED,
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"state", "mount_point"}},
				}

				shareRes, err := client.UpdateReceivedShare(granteeCtx, updateRequest)
				if err != nil || shareRes.Status.Code != rpc.Code_CODE_OK {
					return response.MetaServerError.StatusCode, "grpc update received share request failed", err
				}
			}
		}
	}
	return response.MetaOK.StatusCode, "", nil
}

func (h *Handler) extractPermissions(reqRole string, reqPermissions string, ri *provider.ResourceInfo, defaultPermissions *conversions.Role) (*conversions.Role, []byte, *ocsError) {
	var role *conversions.Role

	// the share role overrides the requested permissions
	if reqRole != "" {
		role = conversions.RoleFromName(reqRole)
	}

	// if the role is unknown - fall back to reqPermissions or defaultPermissions
	if role == nil || role.Name == conversions.RoleUnknown {
		// map requested permissions
		if reqPermissions == "" {
			// TODO default link vs user share
			role = defaultPermissions
		} else {
			pint, err := strconv.Atoi(reqPermissions)
			if err != nil {
				return nil, nil, &ocsError{
					Code:    response.MetaBadRequest.StatusCode,
					Message: "permissions must be an integer",
					Error:   err,
				}
			}
			perm, err := conversions.NewPermissions(pint)
			if err != nil {
				if err == conversions.ErrPermissionNotInRange {
					return nil, nil, &ocsError{
						Code:    http.StatusNotFound,
						Message: err.Error(),
						Error:   err,
					}
				}
				return nil, nil, &ocsError{
					Code:    response.MetaBadRequest.StatusCode,
					Message: err.Error(),
					Error:   err,
				}
			}
			role = conversions.RoleFromOCSPermissions(perm)
		}
	}

	permissions := role.OCSPermissions()
	if ri != nil && ri.Type == provider.ResourceType_RESOURCE_TYPE_FILE && permissions != conversions.PermissionInvalid {
		// Single file shares should never have delete or create permissions
		permissions &^= conversions.PermissionCreate
		permissions &^= conversions.PermissionDelete
		if permissions == conversions.PermissionInvalid {
			return nil, nil, &ocsError{
				Code:    response.MetaBadRequest.StatusCode,
				Message: "Cannot set the requested share permissions",
				Error:   errors.New("cannot set the requested share permissions"),
			}
		}
		role = conversions.RoleFromOCSPermissions(permissions)
	}

	existingPermissions := conversions.RoleFromResourcePermissions(ri.PermissionSet).OCSPermissions()
	if !existingPermissions.Contain(permissions) {
		return nil, nil, &ocsError{
			Code:    http.StatusNotFound,
			Message: "Cannot set the requested share permissions",
			Error:   errors.New("cannot set the requested share permissions"),
		}
	}

	roleMap := map[string]string{"name": role.Name}
	val, err := json.Marshal(roleMap)
	if err != nil {
		return nil, nil, &ocsError{
			Code:    response.MetaServerError.StatusCode,
			Message: "could not encode role",
			Error:   err,
		}
	}

	return role, val, nil
}

// PublicShareContextName represent cross boundaries context for the name of the public share
type PublicShareContextName string

// GetShare handles GET requests on /apps/files_sharing/api/v1/shares/(shareid)
func (h *Handler) GetShare(w http.ResponseWriter, r *http.Request) {
	var share *conversions.ShareData
	var resourceID *provider.ResourceId
	shareID := chi.URLParam(r, "shareid")
	ctx := r.Context()
	logger := appctx.GetLogger(r.Context())
	logger.Debug().Str("shareID", shareID).Msg("get share by id")
	client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting grpc gateway client", err)
		return
	}

	logger.Debug().Str("shareID", shareID).Msg("get public share by id")
	psRes, err := client.GetPublicShare(r.Context(), &link.GetPublicShareRequest{
		Ref: &link.PublicShareReference{
			Spec: &link.PublicShareReference_Id{
				Id: &link.PublicShareId{
					OpaqueId: shareID,
				},
			},
		},
	})

	// FIXME: the backend is returning an err when the public share is not found
	// the below code can be uncommented once error handling is normalized
	// to return Code_CODE_NOT_FOUND when a public share was not found
	/*
		if err != nil {
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error making GetPublicShare grpc request", err)
			return
		}

		if psRes.Status.Code != rpc.Code_CODE_OK && psRes.Status.Code != rpc.Code_CODE_NOT_FOUND {
			logger.Error().Err(err).Msgf("grpc get public share request failed, code: %v", psRes.Status.Code.String)
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc get public share request failed", err)
			return
		}

	*/

	if err == nil && psRes.GetShare() != nil {
		share = conversions.PublicShare2ShareData(psRes.Share, r, h.publicURL)
		resourceID = psRes.Share.ResourceId
	}

	if share == nil {
		// check if we have a user share
		logger.Debug().Str("shareID", shareID).Msg("get user share by id")
		uRes, err := client.GetShare(r.Context(), &collaboration.GetShareRequest{
			Ref: &collaboration.ShareReference{
				Spec: &collaboration.ShareReference_Id{
					Id: &collaboration.ShareId{
						OpaqueId: shareID,
					},
				},
			},
		})

		// FIXME: the backend is returning an err when the public share is not found
		// the below code can be uncommented once error handling is normalized
		// to return Code_CODE_NOT_FOUND when a public share was not found
		/*
			if err != nil {
				response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error making GetShare grpc request", err)
				return
			}

			if uRes.Status.Code != rpc.Code_CODE_OK && uRes.Status.Code != rpc.Code_CODE_NOT_FOUND {
				logger.Error().Err(err).Msgf("grpc get user share request failed, code: %v", uRes.Status.Code)
				response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc get user share request failed", err)
				return
			}
		*/

		if err == nil && uRes.GetShare() != nil {
			resourceID = uRes.Share.ResourceId
			share, err = conversions.CS3Share2ShareData(ctx, uRes.Share)
			if err != nil {
				response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error mapping share data", err)
				return
			}
		}
	}

	if share == nil {
		logger.Debug().Str("shareID", shareID).Msg("no share found with this id")
		response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "share not found", nil)
		return
	}

	info, status, err := h.getResourceInfoByID(ctx, client, resourceID)
	if err != nil {
		log.Error().Err(err).Msg("error mapping share data")
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error mapping share data", err)
		return
	}

	if status.Code != rpc.Code_CODE_OK {
		log.Error().Err(err).Str("status", status.Code.String()).Msg("error mapping share data")
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error mapping share data", err)
		return
	}

	err = h.addFileInfo(ctx, share, info)
	if err != nil {
		log.Error().Err(err).Msg("error mapping share data")
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error mapping share data", err)
	}
	h.mapUserIds(ctx, client, share)

	response.WriteOCSSuccess(w, r, []*conversions.ShareData{share})
}

// UpdateShare handles PUT requests on /apps/files_sharing/api/v1/shares/(shareid)
func (h *Handler) UpdateShare(w http.ResponseWriter, r *http.Request) {
	shareID := chi.URLParam(r, "shareid")
	// FIXME: isPublicShare is already doing a GetShare and GetPublicShare,
	// we should just reuse that object when doing updates
	if h.isPublicShare(r, shareID) {
		h.updatePublicShare(w, r, shareID)
		return
	}
	h.updateShare(w, r, shareID) // TODO PUT is used with incomplete data to update a share}
}

func (h *Handler) updateShare(w http.ResponseWriter, r *http.Request, shareID string) {
	ctx := r.Context()

	client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting grpc gateway client", err)
		return
	}

	shareR, err := client.GetShare(r.Context(), &collaboration.GetShareRequest{
		Ref: &collaboration.ShareReference{
			Spec: &collaboration.ShareReference_Id{
				Id: &collaboration.ShareId{
					OpaqueId: shareID,
				},
			},
		},
	})

	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc update share request", err)
		return
	}

	if shareR.Status.GetCode() != rpc.Code_CODE_OK {
		// TODO: error code from shareR response
		response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "cant find requested share", fmt.Errorf("Can't find share %s. Response code: %v", shareID, shareR.Status.GetCode()))
		return
	}

	info, status, err := h.getResourceInfoByID(ctx, client, shareR.Share.ResourceId)
	if err != nil || status.Code != rpc.Code_CODE_OK {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error mapping share data", err)
		return
	}

	role, _, ocsErr := h.extractPermissions(r.FormValue("role"), r.FormValue("permissions"), info, conversions.NewManagerRole())
	if ocsErr != nil {
		response.WriteOCSError(w, r, ocsErr.Code, ocsErr.Message, ocsErr.Error)
		return
	}

	uReq := &collaboration.UpdateShareRequest{
		Ref: &collaboration.ShareReference{
			Spec: &collaboration.ShareReference_Id{
				Id: &collaboration.ShareId{
					OpaqueId: shareID,
				},
			},
		},
		Field: &collaboration.UpdateShareRequest_UpdateField{
			Field: &collaboration.UpdateShareRequest_UpdateField_Permissions{
				Permissions: &collaboration.SharePermissions{
					// this completely overwrites the permissions for this user
					Permissions: role.CS3ResourcePermissions(),
				},
			},
		},
	}
	uRes, err := client.UpdateShare(ctx, uReq)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc update share request", err)
		return
	}

	if uRes.Status.Code != rpc.Code_CODE_OK {
		if uRes.Status.Code == rpc.Code_CODE_NOT_FOUND {
			response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "not found", nil)
			return
		}
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc update share request failed", err)
		return
	}

	share, err := conversions.CS3Share2ShareData(ctx, uRes.Share)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error mapping share data", err)
		return
	}

	statReq := provider.StatRequest{Ref: &provider.Reference{
		ResourceId: uRes.Share.ResourceId,
	}}

	statRes, err := client.Stat(r.Context(), &statReq)
	if err != nil {
		log.Debug().Err(err).Str("shares", "update user share").Msg("error during stat")
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "missing resource information", fmt.Errorf("error getting resource information"))
		return
	}

	if statRes.Status.Code != rpc.Code_CODE_OK {
		if statRes.Status.Code == rpc.Code_CODE_NOT_FOUND {
			response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "update user share: resource not found", err)
			return
		}

		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc stat request failed for stat after updating user share", err)
		return
	}

	err = h.addFileInfo(r.Context(), share, statRes.Info)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, err.Error(), err)
		return
	}
	h.mapUserIds(ctx, client, share)

	response.WriteOCSSuccess(w, r, share)
}

// RemoveShare handles DELETE requests on /apps/files_sharing/api/v1/shares/(shareid)
func (h *Handler) RemoveShare(w http.ResponseWriter, r *http.Request) {
	shareID := chi.URLParam(r, "shareid")
	switch {
	case h.isPublicShare(r, shareID):
		h.removePublicShare(w, r, shareID)
	case h.isUserShare(r, shareID):
		h.removeUserShare(w, r, shareID)
	default:
		// The request is a remove space member request.
		h.removeSpaceMember(w, r, shareID)
	}
}

// ListShares handles GET requests on /apps/files_sharing/api/v1/shares
func (h *Handler) ListShares(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("shared_with_me") != "" {
		var err error
		listSharedWithMe, err := strconv.ParseBool(r.FormValue("shared_with_me"))
		if err != nil {
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error mapping share data", err)
		}
		if listSharedWithMe {
			h.listSharesWithMe(w, r)
			return
		}
	}
	h.listSharesWithOthers(w, r)
}

const (
	ocsStateUnknown  = -1
	ocsStateAccepted = 0
	ocsStatePending  = 1
	ocsStateRejected = 2
)

func (h *Handler) listSharesWithMe(w http.ResponseWriter, r *http.Request) {
	// which pending state to list
	stateFilter := getStateFilter(r.FormValue("state"))

	log := appctx.GetLogger(r.Context())
	client, err := h.getClient()
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting grpc gateway client", err)
		return
	}

	ctx := r.Context()

	var pinfo *provider.ResourceInfo
	p := r.URL.Query().Get("path")
	shareRef := r.URL.Query().Get("share_ref")
	// we need to lookup the resource id so we can filter the list of shares later
	if p != "" || shareRef != "" {
		ref, err := h.extractReference(r)
		if err != nil {
			response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, errParsingSpaceReference.Error(), errParsingSpaceReference)
			return
		}

		var status *rpc.Status
		pinfo, status, err = h.getResourceInfoByReference(ctx, client, &ref)
		if err != nil {
			response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc stat request", err)
			return
		}
		if status.Code != rpc.Code_CODE_OK {
			switch status.Code {
			case rpc.Code_CODE_NOT_FOUND:
				response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "path not found", nil)
			case rpc.Code_CODE_PERMISSION_DENIED:
				response.WriteOCSError(w, r, response.MetaUnauthorized.StatusCode, "permission denied", nil)
			default:
				response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc stat request failed", nil)
			}
			return
		}
	}

	filters := []*collaboration.Filter{}
	var shareTypes []string
	shareTypesParam := r.URL.Query().Get("share_types")
	if shareTypesParam != "" {
		shareTypes = strings.Split(shareTypesParam, ",")
	}
	for _, s := range shareTypes {
		if s == "" {
			continue
		}
		shareType, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "invalid share type", err)
			return
		}
		switch shareType {
		case int(conversions.ShareTypeUser):
			filters = append(filters, share.UserGranteeFilter())
		case int(conversions.ShareTypeGroup):
			filters = append(filters, share.GroupGranteeFilter())
		}
	}

	if len(shareTypes) != 0 && len(filters) == 0 {
		// If a share_types filter was set for anything other than user or group shares just return an empty response
		response.WriteOCSSuccess(w, r, []*conversions.ShareData{})
		return
	}

	lrsRes, err := client.ListReceivedShares(ctx, &collaboration.ListReceivedSharesRequest{Filters: filters})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc ListReceivedShares request", err)
		return
	}

	if lrsRes.Status.Code != rpc.Code_CODE_OK {
		if lrsRes.Status.Code == rpc.Code_CODE_NOT_FOUND {
			response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "not found", nil)
			return
		}
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc ListReceivedShares request failed", err)
		return
	}

	shares := make([]*conversions.ShareData, 0, len(lrsRes.GetShares()))

	// TODO(refs) filter out "invalid" shares
	for _, rs := range lrsRes.GetShares() {
		if stateFilter != ocsStateUnknown && rs.GetState() != stateFilter {
			continue
		}
		var info *provider.ResourceInfo
		if pinfo != nil {
			// check if the shared resource matches the path resource
			if !utils.ResourceIDEqual(rs.Share.ResourceId, pinfo.Id) {
				// try next share
				continue
			}
			// we can reuse the stat info
			info = pinfo
		} else {
			var status *rpc.Status
			// FIXME the ResourceID is the id of the resource, but we want the id of the mount point so we can fetch that path, well we have the mountpoint path in the receivedshare
			// first stat mount point, but the shares storage provider only handles accepted shares so we send the try to make the requests for only those
			if rs.State == collaboration.ShareState_SHARE_STATE_ACCEPTED {
				mountID := &provider.ResourceId{
					StorageId: storagespace.FormatStorageID(utils.ShareStorageProviderID, utils.ShareStorageSpaceID),
					OpaqueId:  rs.Share.Id.OpaqueId,
				}
				info, status, err = h.getResourceInfoByID(ctx, client, mountID)
			}
			if rs.State != collaboration.ShareState_SHARE_STATE_ACCEPTED || err != nil || status.Code != rpc.Code_CODE_OK {
				// fallback to unmounted resource
				info, status, err = h.getResourceInfoByID(ctx, client, rs.Share.ResourceId)
				if err != nil || status.Code != rpc.Code_CODE_OK {
					h.logProblems(status, err, "could not stat, skipping")
					continue
				}
			}
		}

		data, err := conversions.CS3Share2ShareData(r.Context(), rs.Share)
		if err != nil {
			log.Debug().Interface("share", rs.Share).Interface("shareData", data).Err(err).Msg("could not CS3Share2ShareData, skipping")
			continue
		}

		data.State = mapState(rs.GetState())

		if err := h.addFileInfo(ctx, data, info); err != nil {
			log.Debug().Interface("received_share", rs.Share.Id).Err(err).Msg("could not add file info, skipping")
			continue
		}
		h.mapUserIds(r.Context(), client, data)

		if data.State == ocsStateAccepted {
			// only accepted shares can be accessed when jailing users into their home.
			// in this case we cannot stat shared resources that are outside the users home (/home),
			// the path (/users/u-u-i-d/foo) will not be accessible

			// in a global namespace we can access the share using the full path
			// in a jailed namespace we have to point to the mount point in the users /Shares jail
			// - needed for oc10 hot migration
			// or use the /dav/spaces/<space id> endpoint?

			// list /Shares and match fileids with list of received shares
			// - only works for a /Shares folder jail
			// - does not work for freely mountable shares as in oc10 because we would need to iterate over the whole tree, there is no listing of mountpoints, yet

			// can we return the mountpoint when the gateway resolves the listing of shares?
			// - no, the gateway only sees the same list any has the same options as the ocs service
			// - we would need to have a list of mountpoints for the shares -> owncloudstorageprovider for hot migration migration

			// best we can do for now is stat the /Shares jail if it is set and return those paths

			// if we are in a jail and the current share has been accepted use the stat from the share jail
			// Needed because received shares can be jailed in a folder in the users home

			if h.sharePrefix != "/" {
				// if we have a mount point use it to build the path
				if rs.MountPoint != nil && rs.MountPoint.Path != "" {
					// override path with info from share jail
					data.FileTarget = path.Join(h.sharePrefix, rs.MountPoint.Path)
					data.Path = path.Join(h.sharePrefix, rs.MountPoint.Path)
				} else {
					data.FileTarget = path.Join(h.sharePrefix, path.Base(info.Path))
					data.Path = path.Join(h.sharePrefix, path.Base(info.Path))
				}
			} else {
				data.FileTarget = info.Path
				data.Path = info.Path
			}
		} else {
			// not accepted shares need their Path jailed to make the testsuite happy

			if h.sharePrefix != "/" {
				data.Path = path.Join("/", path.Base(info.Path))
			}

		}

		shares = append(shares, data)
		log.Debug().Msgf("share: %+v", *data)
	}

	response.WriteOCSSuccess(w, r, shares)
}

func (h *Handler) listSharesWithOthers(w http.ResponseWriter, r *http.Request) {
	shares := make([]*conversions.ShareData, 0)

	filters := []*collaboration.Filter{}
	linkFilters := []*link.ListPublicSharesRequest_Filter{}
	var e error

	// shared with others
	p := r.URL.Query().Get("path")
	s := r.URL.Query().Get("space")
	spaceRef := r.URL.Query().Get("space_ref")
	if p != "" || s != "" || spaceRef != "" {
		ref, err := h.extractReference(r)
		if err != nil {
			response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, errParsingSpaceReference.Error(), errParsingSpaceReference)
			return
		}
		filters, linkFilters, e = h.addFilters(w, r, &ref)
		if e != nil {
			// result has been written as part of addFilters
			return
		}
	}

	var shareTypes []string
	shareTypesParam := r.URL.Query().Get("share_types")
	if shareTypesParam != "" {
		shareTypes = strings.Split(shareTypesParam, ",")
	}

	listPublicShares := len(shareTypes) == 0 // if no share_types filter was set we want to list all share by default
	listUserShares := len(shareTypes) == 0   // if no share_types filter was set we want to list all share by default
	for _, s := range shareTypes {
		if s == "" {
			continue
		}
		shareType, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "invalid share type", err)
			return
		}

		switch shareType {
		case int(conversions.ShareTypeUser):
			listUserShares = true
			filters = append(filters, share.UserGranteeFilter())
		case int(conversions.ShareTypeGroup):
			listUserShares = true
			filters = append(filters, share.GroupGranteeFilter())
		case int(conversions.ShareTypePublicLink):
			listPublicShares = true
		}
	}

	if listPublicShares {
		publicShares, status, err := h.listPublicShares(r, linkFilters)
		h.logProblems(status, err, "could not listPublicShares")
		shares = append(shares, publicShares...)
	}
	if listUserShares {
		userShares, status, err := h.listUserShares(r, filters)
		h.logProblems(status, err, "could not listUserShares")
		shares = append(shares, userShares...)
	}

	response.WriteOCSSuccess(w, r, shares)
}

func (h *Handler) logProblems(s *rpc.Status, e error, msg string) {
	if e != nil {
		// errors need to be taken care of
		log.Error().Err(e).Msg(msg)
		return
	}
	if s != nil && s.Code != rpc.Code_CODE_OK {
		switch s.Code {
		// not found and permission denied can happen during normal operations
		case rpc.Code_CODE_NOT_FOUND:
			log.Debug().Interface("status", s).Msg(msg)
		case rpc.Code_CODE_PERMISSION_DENIED:
			log.Debug().Interface("status", s).Msg(msg)
		default:
			// anything else should not happen, someone needs to dig into it
			log.Error().Interface("status", s).Msg(msg)
		}
	}
}

func (h *Handler) addFilters(w http.ResponseWriter, r *http.Request, ref *provider.Reference) ([]*collaboration.Filter, []*link.ListPublicSharesRequest_Filter, error) {
	ctx := r.Context()

	// first check if the file exists
	client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting grpc gateway client", err)
		return nil, nil, err
	}

	info, status, err := h.getResourceInfoByReference(ctx, client, ref)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error sending a grpc stat request", err)
		return nil, nil, err
	}

	if status.Code != rpc.Code_CODE_OK {
		err = errors.New(status.Message)
		if status.Code == rpc.Code_CODE_NOT_FOUND {
			response.WriteOCSError(w, r, response.MetaNotFound.StatusCode, "not found", err)
			return nil, nil, err
		}
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "grpc stat request failed", err)
		return nil, nil, err
	}

	collaborationFilters := []*collaboration.Filter{share.ResourceIDFilter(info.Id)}
	linkFilters := []*link.ListPublicSharesRequest_Filter{publicshare.ResourceIDFilter(info.Id)}

	return collaborationFilters, linkFilters, nil
}

func (h *Handler) addFileInfo(ctx context.Context, s *conversions.ShareData, info *provider.ResourceInfo) error {
	log := appctx.GetLogger(ctx)
	if info != nil {
		// TODO The owner is not set in the storage stat metadata ...
		parsedMt, _, err := mime.ParseMediaType(info.MimeType)
		if err != nil {
			// Should never happen. We log anyways so that we know if it happens.
			log.Warn().Err(err).Msg("failed to parse mimetype")
		}
		s.MimeType = parsedMt
		// TODO STime:     &types.Timestamp{Seconds: info.Mtime.Seconds, Nanos: info.Mtime.Nanos},
		// TODO Storage: int
		s.ItemSource = storagespace.FormatResourceID(*info.Id)
		s.FileSource = s.ItemSource
		switch {
		case h.sharePrefix == "/":
			s.FileTarget = info.Path
			s.Path = info.Path
			client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
			if err == nil {
				gpRes, err := client.GetPath(ctx, &provider.GetPathRequest{
					ResourceId: info.Id,
				})
				if err == nil && gpRes.Status.Code == rpc.Code_CODE_OK {
					// TODO log error?
					s.Path = gpRes.Path

					// cut off configured home namespace, paths in ocs shares are relative to it
					identifier := h.mustGetIdentifiers(ctx, client, info.GetOwner().GetOpaqueId(), false)
					u := &userpb.User{
						Id:          info.Owner,
						Username:    identifier.Username,
						DisplayName: identifier.DisplayName,
						Mail:        identifier.Mail,
					}
					s.Path = strings.TrimPrefix(s.Path, h.getHomeNamespace(u))
				}
			}
		default:
			s.FileTarget = path.Join(h.sharePrefix, path.Base(info.Path))
			if s.ShareType == conversions.ShareTypePublicLink {
				s.FileTarget = path.Join("/", path.Base(info.Path))
			}
			s.Path = info.Path
			client, err := pool.GetGatewayServiceClient(h.gatewayAddr)
			if err == nil {
				gpRes, err := client.GetPath(ctx, &provider.GetPathRequest{
					ResourceId: info.Id,
				})
				if err == nil && gpRes.Status.Code == rpc.Code_CODE_OK {
					// TODO log error?
					s.Path = gpRes.Path
				}

				// cut off configured home namespace, paths in ocs shares are relative to it
				identifier := h.mustGetIdentifiers(ctx, client, info.GetOwner().GetOpaqueId(), false)
				u := &userpb.User{
					Id:          info.Owner,
					Username:    identifier.Username,
					DisplayName: identifier.DisplayName,
					Mail:        identifier.Mail,
				}
				s.Path = strings.TrimPrefix(s.Path, h.getHomeNamespace(u))
			}
		}
		s.StorageID = storageIDPrefix + s.FileTarget
		// TODO FileParent:
		// item type
		s.ItemType = conversions.ResourceType(info.GetType()).String()

		// file owner might not yet be set. Use file info
		if s.UIDFileOwner == "" {
			s.UIDFileOwner = info.GetOwner().GetOpaqueId()
		}
		// share owner might not yet be set. Use file info
		if s.UIDOwner == "" {
			s.UIDOwner = info.GetOwner().GetOpaqueId()
		}
	}
	return nil
}

// mustGetIdentifiers always returns a struct with identifiers, if the user or group could not be found they will all be empty
func (h *Handler) mustGetIdentifiers(ctx context.Context, client gateway.GatewayAPIClient, id string, isGroup bool) *userIdentifiers {
	sublog := appctx.GetLogger(ctx).With().Str("id", id).Logger()
	if id == "" {
		return &userIdentifiers{}
	}

	if idIf, err := h.userIdentifierCache.Get(id); err == nil {
		sublog.Debug().Msg("cache hit")
		return idIf.(*userIdentifiers)
	}

	sublog.Debug().Msg("cache miss")
	var ui *userIdentifiers

	if isGroup {
		res, err := client.GetGroup(ctx, &grouppb.GetGroupRequest{
			GroupId: &grouppb.GroupId{
				OpaqueId: id,
			},
			SkipFetchingMembers: true,
		})
		if err != nil {
			sublog.Err(err).Msg("could not look up group")
			return &userIdentifiers{}
		}
		if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
			sublog.Err(err).
				Int32("code", int32(res.GetStatus().GetCode())).
				Str("message", res.GetStatus().GetMessage()).
				Msg("get group call failed")
			return &userIdentifiers{}
		}
		if res.Group == nil {
			sublog.Debug().
				Int32("code", int32(res.GetStatus().GetCode())).
				Str("message", res.GetStatus().GetMessage()).
				Msg("group not found")
			return &userIdentifiers{}
		}
		ui = &userIdentifiers{
			DisplayName: res.Group.DisplayName,
			Username:    res.Group.GroupName,
			Mail:        res.Group.Mail,
		}
	} else {
		res, err := client.GetUser(ctx, &userpb.GetUserRequest{
			UserId: &userpb.UserId{
				OpaqueId: id,
			},
			SkipFetchingUserGroups: true,
		})
		if err != nil {
			sublog.Err(err).Msg("could not look up user")
			return &userIdentifiers{}
		}
		if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
			sublog.Err(err).
				Int32("code", int32(res.GetStatus().GetCode())).
				Str("message", res.GetStatus().GetMessage()).
				Msg("get user call failed")
			return &userIdentifiers{}
		}
		if res.User == nil {
			sublog.Debug().
				Int32("code", int32(res.GetStatus().GetCode())).
				Str("message", res.GetStatus().GetMessage()).
				Msg("user not found")
			return &userIdentifiers{}
		}
		ui = &userIdentifiers{
			DisplayName: res.User.DisplayName,
			Username:    res.User.Username,
			Mail:        res.User.Mail,
		}
	}
	_ = h.userIdentifierCache.Set(id, ui)
	log.Debug().Str("id", id).Msg("cache update")
	return ui
}

func (h *Handler) mapUserIds(ctx context.Context, client gateway.GatewayAPIClient, s *conversions.ShareData) {
	if s.UIDOwner != "" {
		owner := h.mustGetIdentifiers(ctx, client, s.UIDOwner, false)
		s.UIDOwner = owner.Username
		if s.DisplaynameOwner == "" {
			s.DisplaynameOwner = owner.DisplayName
		}
		if s.AdditionalInfoFileOwner == "" {
			s.AdditionalInfoFileOwner = h.getAdditionalInfoAttribute(ctx, owner)
		}
	}

	if s.UIDFileOwner != "" {
		fileOwner := h.mustGetIdentifiers(ctx, client, s.UIDFileOwner, false)
		s.UIDFileOwner = fileOwner.Username
		if s.DisplaynameFileOwner == "" {
			s.DisplaynameFileOwner = fileOwner.DisplayName
		}
		if s.AdditionalInfoOwner == "" {
			s.AdditionalInfoOwner = h.getAdditionalInfoAttribute(ctx, fileOwner)
		}
	}

	if s.ShareWith != "" && s.ShareWith != "***redacted***" {
		shareWith := h.mustGetIdentifiers(ctx, client, s.ShareWith, s.ShareType == conversions.ShareTypeGroup)
		s.ShareWith = shareWith.Username
		if s.ShareWithDisplayname == "" {
			s.ShareWithDisplayname = shareWith.DisplayName
		}
		if s.ShareWithAdditionalInfo == "" {
			s.ShareWithAdditionalInfo = h.getAdditionalInfoAttribute(ctx, shareWith)
		}
	}
}

func (h *Handler) getAdditionalInfoAttribute(ctx context.Context, u *userIdentifiers) string {
	var buf bytes.Buffer
	if err := h.additionalInfoTemplate.Execute(&buf, u); err != nil {
		log := appctx.GetLogger(ctx)
		log.Warn().Err(err).Msg("failed to parse additional info template")
		return ""
	}
	return buf.String()
}

func (h *Handler) getResourceInfoByReference(ctx context.Context, client gateway.GatewayAPIClient, ref *provider.Reference) (*provider.ResourceInfo, *rpc.Status, error) {
	var key string
	if ref.ResourceId == nil {
		// This is a path based reference
		key = ref.Path
	} else {
		var err error
		key, err = storagespace.FormatReference(ref)
		if err != nil {
			return nil, nil, err
		}
	}
	return h.getResourceInfo(ctx, client, key, ref)
}

func (h *Handler) getResourceInfoByID(ctx context.Context, client gateway.GatewayAPIClient, id *provider.ResourceId) (*provider.ResourceInfo, *rpc.Status, error) {
	return h.getResourceInfo(ctx, client, storagespace.FormatResourceID(*id), &provider.Reference{ResourceId: id, Path: "."})
}

// getResourceInfo retrieves the resource info to a target.
// This method utilizes caching if it is enabled.
func (h *Handler) getResourceInfo(ctx context.Context, client gateway.GatewayAPIClient, key string, ref *provider.Reference) (*provider.ResourceInfo, *rpc.Status, error) {
	logger := appctx.GetLogger(ctx)

	var pinfo *provider.ResourceInfo
	var status *rpc.Status
	var err error
	var foundInCache bool
	if h.resourceInfoCacheTTL > 0 && h.resourceInfoCache != nil {
		if pinfo, err = h.resourceInfoCache.Get(key); err == nil {
			logger.Debug().Msgf("cache hit for resource %+v", key)
			status = &rpc.Status{Code: rpc.Code_CODE_OK}
			foundInCache = true
		}
	}
	if !foundInCache {
		logger.Debug().Msgf("cache miss for resource %+v, statting", key)
		statReq := &provider.StatRequest{
			Ref: ref,
		}

		statRes, err := client.Stat(ctx, statReq)
		if err != nil {
			return nil, nil, err
		}

		if statRes.Status.Code != rpc.Code_CODE_OK {
			return nil, statRes.Status, nil
		}

		pinfo = statRes.GetInfo()
		status = statRes.Status
		if h.resourceInfoCacheTTL > 0 {
			_ = h.resourceInfoCache.SetWithExpire(key, pinfo, h.resourceInfoCacheTTL)
		}
	}

	return pinfo, status, nil
}

func (h *Handler) createCs3Share(ctx context.Context, w http.ResponseWriter, r *http.Request, client gateway.GatewayAPIClient, req *collaboration.CreateShareRequest) (*collaboration.Share, *ocsError) {
	createShareResponse, err := client.CreateShare(ctx, req)
	if err != nil {
		return nil, &ocsError{
			Code:    response.MetaServerError.StatusCode,
			Message: "error sending a grpc create share request",
			Error:   err,
		}
	}
	if createShareResponse.Status.Code != rpc.Code_CODE_OK {
		if createShareResponse.Status.Code == rpc.Code_CODE_NOT_FOUND {
			return nil, &ocsError{
				Code:    response.MetaNotFound.StatusCode,
				Message: "not found",
				Error:   nil,
			}
		}
		return nil, &ocsError{
			Code:    response.MetaServerError.StatusCode,
			Message: "grpc create share request failed",
			Error:   nil,
		}
	}
	return createShareResponse.Share, nil
}

func mapState(state collaboration.ShareState) int {
	var mapped int
	switch state {
	case collaboration.ShareState_SHARE_STATE_PENDING:
		mapped = ocsStatePending
	case collaboration.ShareState_SHARE_STATE_ACCEPTED:
		mapped = ocsStateAccepted
	case collaboration.ShareState_SHARE_STATE_REJECTED:
		mapped = ocsStateRejected
	default:
		mapped = ocsStateUnknown
	}
	return mapped
}

func getStateFilter(s string) collaboration.ShareState {
	var stateFilter collaboration.ShareState
	switch s {
	case "all":
		stateFilter = ocsStateUnknown // no filter
	case "0": // accepted
		stateFilter = collaboration.ShareState_SHARE_STATE_ACCEPTED
	case "1": // pending
		stateFilter = collaboration.ShareState_SHARE_STATE_PENDING
	case "2": // rejected
		stateFilter = collaboration.ShareState_SHARE_STATE_REJECTED
	default:
		stateFilter = collaboration.ShareState_SHARE_STATE_ACCEPTED
	}
	return stateFilter
}

func (h *Handler) getPoolClient() (gateway.GatewayAPIClient, error) {
	return pool.GetGatewayServiceClient(h.gatewayAddr)
}

func (h *Handler) getHomeNamespace(u *userpb.User) string {
	return templates.WithUser(u, h.homeNamespace)
}
