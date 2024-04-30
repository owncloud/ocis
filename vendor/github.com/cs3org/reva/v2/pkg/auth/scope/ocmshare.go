// Copyright 2018-2023 CERN
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

package scope

import (
	"context"
	"path/filepath"
	"strings"

	appprovider "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	appregistry "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ocmv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	registry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
)

// FIXME: the namespace here is hardcoded
// find a way to pass it from the config.
const ocmNamespace = "/ocm"

func ocmShareScope(_ context.Context, scope *authpb.Scope, resource interface{}, _ *zerolog.Logger) (bool, error) {
	var share ocmv1beta1.Share
	if err := utils.UnmarshalJSONToProtoV1(scope.Resource.Value, &share); err != nil {
		return false, err
	}

	switch v := resource.(type) {
	// viewer role
	case *registry.ListStorageProvidersRequest:
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
		return checkStorageRefForOCMShare(&share, ref, ocmNamespace), nil
	case *registry.GetStorageProvidersRequest:
		return checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.StatRequest:
		return checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.ListContainerRequest:
		return checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.InitiateFileDownloadRequest:
		return checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *appprovider.OpenInAppRequest:
		return checkStorageRefForOCMShare(&share, &provider.Reference{ResourceId: v.ResourceInfo.Id}, ocmNamespace), nil
	case *gateway.OpenInAppRequest:
		return checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.GetLockRequest:
		return checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil

	// editor role
	case *provider.CreateContainerRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.TouchFileRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.DeleteRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.MoveRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetSource(), ocmNamespace) && checkStorageRefForOCMShare(&share, v.GetDestination(), ocmNamespace), nil
	case *provider.InitiateFileUploadRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.SetArbitraryMetadataRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.UnsetArbitraryMetadataRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.SetLockRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.RefreshLockRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil
	case *provider.UnlockRequest:
		return hasRoleEditor(*scope) && checkStorageRefForOCMShare(&share, v.GetRef(), ocmNamespace), nil

	// App provider requests
	case *appregistry.GetDefaultAppProviderForMimeTypeRequest:
		return true, nil
	case *appregistry.GetAppProvidersRequest:
		return true, nil
	case *userv1beta1.GetUserByClaimRequest:
		return true, nil
	case *userv1beta1.GetUserRequest:
		return true, nil
	case *provider.ListStorageSpacesRequest:
		return true, nil

	case *ocmv1beta1.GetOCMShareRequest:
		return checkOCMShareRef(&share, v.GetRef()), nil
	case *ocmv1beta1.GetOCMShareByTokenRequest:
		return share.Token == v.GetToken(), nil
	case string:
		return checkResourcePath(v), nil
	}
	return false, nil
}

func checkStorageRefForOCMShare(s *ocmv1beta1.Share, r *provider.Reference, ns string) bool {
	if r.ResourceId != nil {
		return utils.ResourceIDEqual(s.ResourceId, r.GetResourceId()) || strings.HasPrefix(r.ResourceId.OpaqueId, s.Token)
	}

	// FIXME: the paths here are hardcoded
	if strings.HasPrefix(r.GetPath(), "/public/"+s.Token) {
		return true
	}
	return strings.HasPrefix(r.GetPath(), filepath.Join(ns, s.Token))
}

func checkOCMShareRef(s *ocmv1beta1.Share, ref *ocmv1beta1.ShareReference) bool {
	return ref.GetToken() == s.Token
}

// AddOCMShareScope adds the scope to allow access to an OCM share and the share resource.
func AddOCMShareScope(share *ocmv1beta1.Share, role authpb.Role, scopes map[string]*authpb.Scope) (map[string]*authpb.Scope, error) {
	// Create a new "scope share" to only expose the required fields `ResourceId` and `Token` to the scope.
	scopeShare := ocmv1beta1.Share{ResourceId: share.ResourceId, Token: share.Token}
	val, err := utils.MarshalProtoV1ToJSON(&scopeShare)
	if err != nil {
		return nil, err
	}
	if scopes == nil {
		scopes = make(map[string]*authpb.Scope)
	}

	scopes["ocmshare:"+share.Id.OpaqueId] = &authpb.Scope{
		Resource: &types.OpaqueEntry{
			Decoder: "json",
			Value:   val,
		},
		Role: role,
	}
	return scopes, nil
}

// GetOCMSharesFromScopes returns all OCM shares in the given scope.
func GetOCMSharesFromScopes(scopes map[string]*authpb.Scope) ([]*ocmv1beta1.Share, error) {
	var shares []*ocmv1beta1.Share
	for k, s := range scopes {
		if strings.HasPrefix(k, "ocmshare:") {
			res := s.Resource
			if res.Decoder != "json" {
				return nil, errtypes.InternalError("resource should be json encoded")
			}
			var share ocmv1beta1.Share
			err := utils.UnmarshalJSONToProtoV1(res.Value, &share)
			if err != nil {
				return nil, err
			}
			shares = append(shares, &share)
		}
	}
	return shares, nil
}
