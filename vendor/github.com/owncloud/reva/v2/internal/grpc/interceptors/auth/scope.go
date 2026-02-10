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

package auth

import (
	"context"
	"fmt"
	"strings"

	appprovider "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	appregistry "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	ocmv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	registry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	"github.com/owncloud/reva/v2/pkg/appctx"
	"github.com/owncloud/reva/v2/pkg/auth/scope"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	statuspkg "github.com/owncloud/reva/v2/pkg/rgrpc/status"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/token"
	"github.com/owncloud/reva/v2/pkg/utils"
	"google.golang.org/grpc/metadata"
)

func expandAndVerifyScope(ctx context.Context, req interface{}, tokenScope map[string]*authpb.Scope, user *userpb.User, gatewayAddr string, mgr token.Manager) error {
	log := appctx.GetLogger(ctx)
	client, err := pool.GetGatewayServiceClient(gatewayAddr)
	if err != nil {
		return err
	}

	if ref, ok := extractRef(req, tokenScope); ok {
		// The request is for a storage reference. This can be the case for multiple scenarios:
		// - If the path is not empty, the request might be coming from a share where the accessor is
		//   trying to impersonate the owner, since the share manager doesn't know the
		//   share path.
		// - If the ID not empty, the request might be coming from
		//   - a resource present inside a shared folder, or
		//   - a share created for a lightweight account after the token was minted.
		log.Info().Msgf("resolving storage reference to check token scope %s", ref.String())
		for k := range tokenScope {
			switch {
			case strings.HasPrefix(k, "publicshare"):
				if err = resolvePublicShare(ctx, ref, tokenScope[k], client, mgr); err == nil {
					return nil
				}

			case strings.HasPrefix(k, "share"):
				if err = resolveUserShare(ctx, ref, tokenScope[k], client, mgr); err == nil {
					return nil
				}

			case strings.HasPrefix(k, "ocmshare"):
				if err = resolveOCMShare(ctx, ref, tokenScope[k], client, mgr); err == nil {
					return nil
				}
			}
			log.Err(err).Interface("ref", ref).Interface("scope", k).Msg("error resolving reference under scope")
		}
	}

	return errtypes.PermissionDenied(fmt.Sprintf("access to resource %+v not allowed within the assigned scope", req))
}

func resolvePublicShare(ctx context.Context, ref *provider.Reference, scope *authpb.Scope, client gateway.GatewayAPIClient, mgr token.Manager) error {
	var share link.PublicShare
	err := utils.UnmarshalJSONToProtoV1(scope.Resource.Value, &share)
	if err != nil {
		return err
	}

	return checkIfNestedResource(ctx, ref, share.ResourceId, client, mgr)
}

func resolveOCMShare(ctx context.Context, ref *provider.Reference, scope *authpb.Scope, client gateway.GatewayAPIClient, mgr token.Manager) error {
	var share ocmv1beta1.Share
	if err := utils.UnmarshalJSONToProtoV1(scope.Resource.Value, &share); err != nil {
		return err
	}

	// for ListOCMSharesRequest, the ref resource id is empty and we set path to . to indicate the root of the share
	if ref.GetResourceId() == nil && ref.Path == "." {
		ref.ResourceId = share.GetResourceId()
	}

	return checkIfNestedResource(ctx, ref, share.ResourceId, client, mgr)
}

func resolveUserShare(ctx context.Context, ref *provider.Reference, scope *authpb.Scope, client gateway.GatewayAPIClient, mgr token.Manager) error {
	var share collaboration.Share
	err := utils.UnmarshalJSONToProtoV1(scope.Resource.Value, &share)
	if err != nil {
		return err
	}

	return checkIfNestedResource(ctx, ref, share.ResourceId, client, mgr)
}

func checkIfNestedResource(ctx context.Context, ref *provider.Reference, shareRoot *provider.ResourceId, client gateway.GatewayAPIClient, mgr token.Manager) error {
	// Since the resource ID is obtained from the scope, the current token
	// has access to it.
	rootStat, err := client.Stat(ctx, &provider.StatRequest{Ref: &provider.Reference{ResourceId: shareRoot}})
	if err != nil {
		return err
	}
	if rootStat.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return statuspkg.NewErrorFromCode(rootStat.Status.Code, "auth interceptor")
	}

	rootInfo := rootStat.GetInfo()

	// We need to find out if the requested resource is a child of the `shareRoot` (coming from token scope)
	// We mint a token as the owner of the public share and try to stat the reference

	var user *userpb.User
	if rootInfo.GetOwner().GetType() == userpb.UserType_USER_TYPE_SPACE_OWNER {
		// fake a space owner user
		user = &userpb.User{
			Id: rootInfo.GetOwner(),
		}
	} else {
		userResp, err := client.GetUser(ctx, &userpb.GetUserRequest{UserId: rootInfo.GetOwner(), SkipFetchingUserGroups: true})
		if err != nil || userResp.Status.Code != rpc.Code_CODE_OK {
			return err
		}
		user = userResp.User
	}

	scope, err := scope.AddOwnerScope(map[string]*authpb.Scope{})
	if err != nil {
		return err
	}
	token, err := mgr.MintToken(ctx, user, scope)
	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(context.Background(), ctxpkg.TokenHeader, token)

	resourceStat, err := client.Stat(ctx, &provider.StatRequest{Ref: ref})
	if err != nil {
		return err
	}
	if resourceStat.GetStatus().GetCode() == rpc.Code_CODE_NOT_FOUND && ref.GetPath() != "" && ref.GetPath() != "." {
		// The resource does not seem to exist (yet?). We might be part of an initiate upload request.
		// Stat the parent to get its path and check that against the root path.
		resourceStat, err = client.Stat(ctx, &provider.StatRequest{Ref: &provider.Reference{ResourceId: ref.GetResourceId()}})
		if err != nil {
			return err
		}
	}
	if resourceStat.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return statuspkg.NewErrorFromCode(resourceStat.Status.Code, "auth interceptor")
	}

	// Check if the resource and the share root are in the same storage space
	ci := resourceStat.GetInfo()
	if ci.GetId().GetStorageId() != rootInfo.GetId().GetStorageId() ||
		ci.GetId().GetSpaceId() != rootInfo.GetId().GetSpaceId() {
		return errtypes.PermissionDenied("invalid resource")
	}

	// check if the resource path is subpath of the share root path
	rootPath, err := getPath(ctx, rootStat.GetInfo().GetId(), client)
	if err != nil {
		return err
	}

	resourcePath, err := getPath(ctx, resourceStat.GetInfo().GetId(), client)
	if err != nil {
		return err
	}

	if rootPath == "/" || resourcePath == rootPath || strings.HasPrefix(resourcePath, rootPath+"/") {
		return nil
	}

	return errtypes.PermissionDenied("invalid resource")
}

func getPath(ctx context.Context, resourceId *provider.ResourceId, client gateway.GatewayAPIClient) (string, error) {
	pathResp, err := client.GetPath(ctx, &provider.GetPathRequest{ResourceId: resourceId})
	if err != nil {
		return "", err
	}
	if pathResp.Status.Code != rpc.Code_CODE_OK {
		return "", statuspkg.NewErrorFromCode(pathResp.Status.Code, "auth interceptor")
	}
	return pathResp.Path, nil
}

func extractRefFromListProvidersReq(v *registry.ListStorageProvidersRequest) (*provider.Reference, bool) {
	ref := &provider.Reference{}
	if v.Opaque != nil && v.Opaque.Map != nil {
		if e, ok := v.Opaque.Map["storage_id"]; ok {
			if ref.ResourceId == nil {
				ref.ResourceId = &provider.ResourceId{}
			}
			ref.ResourceId.StorageId = string(e.Value)
		}
		if e, ok := v.Opaque.Map["space_id"]; ok {
			if ref.ResourceId == nil {
				ref.ResourceId = &provider.ResourceId{}
			}
			ref.ResourceId.SpaceId = string(e.Value)
		}
		if e, ok := v.Opaque.Map["opaque_id"]; ok {
			if ref.ResourceId == nil {
				ref.ResourceId = &provider.ResourceId{}
			}
			ref.ResourceId.OpaqueId = string(e.Value)
		}
		if e, ok := v.Opaque.Map["path"]; ok {
			ref.Path = string(e.Value)
		}
	}
	return ref, true
}

func extractRefForReaderRole(req interface{}) (*provider.Reference, bool) {
	switch v := req.(type) {
	// Read requests
	case *registry.GetStorageProvidersRequest:
		return v.GetRef(), true
	case *registry.ListStorageProvidersRequest:
		return extractRefFromListProvidersReq(v)
	case *provider.StatRequest:
		return v.GetRef(), true
	case *provider.ListContainerRequest:
		return v.GetRef(), true
	case *provider.InitiateFileDownloadRequest:
		return v.GetRef(), true

	// App provider requests
	case *appregistry.GetAppProvidersRequest:
		return &provider.Reference{ResourceId: v.ResourceInfo.Id}, true
	case *appprovider.OpenInAppRequest:
		return &provider.Reference{ResourceId: v.ResourceInfo.Id}, true
	case *gateway.OpenInAppRequest:
		return v.GetRef(), true

	// Locking
	case *provider.GetLockRequest:
		return v.GetRef(), true
	case *provider.SetLockRequest:
		return v.GetRef(), true
	case *provider.RefreshLockRequest:
		return v.GetRef(), true
	case *provider.UnlockRequest:
		return v.GetRef(), true

	// OCM shares
	case *ocmv1beta1.ListReceivedOCMSharesRequest:
		return &provider.Reference{Path: "."}, true // we will try to stat the shared node

	}

	return nil, false

}

func extractRefForUploaderRole(req interface{}) (*provider.Reference, bool) {
	switch v := req.(type) {
	// Write Requests
	case *registry.GetStorageProvidersRequest:
		return v.GetRef(), true
	case *registry.ListStorageProvidersRequest:
		return extractRefFromListProvidersReq(v)
	case *provider.StatRequest:
		return v.GetRef(), true
	case *provider.CreateContainerRequest:
		return v.GetRef(), true
	case *provider.TouchFileRequest:
		return v.GetRef(), true
	case *provider.InitiateFileUploadRequest:
		return v.GetRef(), true

	// App provider requests
	case *appregistry.GetAppProvidersRequest:
		return &provider.Reference{ResourceId: v.ResourceInfo.Id}, true
	case *appprovider.OpenInAppRequest:
		return &provider.Reference{ResourceId: v.ResourceInfo.Id}, true
	case *gateway.OpenInAppRequest:
		return v.GetRef(), true

	// Locking
	case *provider.GetLockRequest:
		return v.GetRef(), true
	case *provider.SetLockRequest:
		return v.GetRef(), true
	case *provider.RefreshLockRequest:
		return v.GetRef(), true
	case *provider.UnlockRequest:
		return v.GetRef(), true
	}

	return nil, false

}

func extractRefForEditorRole(req interface{}) (*provider.Reference, bool) {
	switch v := req.(type) {
	// Remaining edit Requests
	case *provider.DeleteRequest:
		return v.GetRef(), true
	case *provider.MoveRequest:
		return v.GetSource(), true
	case *provider.SetArbitraryMetadataRequest:
		return v.GetRef(), true
	case *provider.UnsetArbitraryMetadataRequest:
		return v.GetRef(), true
	}

	return nil, false

}

func extractRef(req interface{}, tokenScope map[string]*authpb.Scope) (*provider.Reference, bool) {
	var readPerm, uploadPerm, editPerm bool
	for _, v := range tokenScope {
		if v.Role == authpb.Role_ROLE_OWNER || v.Role == authpb.Role_ROLE_EDITOR || v.Role == authpb.Role_ROLE_VIEWER {
			readPerm = true
		}
		if v.Role == authpb.Role_ROLE_OWNER || v.Role == authpb.Role_ROLE_EDITOR || v.Role == authpb.Role_ROLE_UPLOADER {
			uploadPerm = true
		}
		if v.Role == authpb.Role_ROLE_OWNER || v.Role == authpb.Role_ROLE_EDITOR {
			editPerm = true
		}
	}

	if readPerm {
		ref, ok := extractRefForReaderRole(req)
		if ok {
			return ref, true
		}
	}
	if uploadPerm {
		ref, ok := extractRefForUploaderRole(req)
		if ok {
			return ref, true
		}
	}
	if editPerm {
		ref, ok := extractRefForEditorRole(req)
		if ok {
			return ref, true
		}
	}

	return nil, false
}
