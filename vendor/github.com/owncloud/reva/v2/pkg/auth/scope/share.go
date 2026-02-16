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

package scope

import (
	"context"
	"fmt"
	"strings"

	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	registry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
)

func shareScope(_ context.Context, scope *authpb.Scope, resource interface{}, logger *zerolog.Logger) (bool, error) {
	var share collaboration.Share
	err := utils.UnmarshalJSONToProtoV1(scope.Resource.Value, &share)
	if err != nil {
		return false, err
	}

	switch v := resource.(type) {
	// Viewer role
	case *registry.GetStorageProvidersRequest:
		return checkShareStorageRef(&share, v.GetRef()), nil
	case *provider.StatRequest:
		return checkShareStorageRef(&share, v.GetRef()), nil
	case *provider.ListContainerRequest:
		return checkShareStorageRef(&share, v.GetRef()), nil
	case *provider.InitiateFileDownloadRequest:
		return checkShareStorageRef(&share, v.GetRef()), nil

		// Editor role
		// TODO(ishank011): Add role checks,
		// need to return appropriate status codes in the ocs/ocdav layers.
	case *provider.CreateContainerRequest:
		return checkShareStorageRef(&share, v.GetRef()), nil
	case *provider.TouchFileRequest:
		return checkShareStorageRef(&share, v.GetRef()), nil
	case *provider.DeleteRequest:
		return checkShareStorageRef(&share, v.GetRef()), nil
	case *provider.MoveRequest:
		return checkShareStorageRef(&share, v.GetSource()) && checkShareStorageRef(&share, v.GetDestination()), nil
	case *provider.InitiateFileUploadRequest:
		return checkShareStorageRef(&share, v.GetRef()), nil

	case *collaboration.ListReceivedSharesRequest:
		return true, nil
	case *collaboration.GetReceivedShareRequest:
		return checkShareRef(&share, v.GetRef()), nil
	case string:
		return checkSharePath(v) || checkResourcePath(v), nil
	}

	msg := fmt.Sprintf("resource type assertion failed: %+v", resource)
	logger.Debug().Str("scope", "shareScope").Msg(msg)
	return false, errtypes.InternalError(msg)
}

func checkShareStorageRef(s *collaboration.Share, r *provider.Reference) bool {
	// ref: <id:<storage_id:$storageID opaque_id:$opaqueID > >
	if r.GetResourceId() != nil && r.Path == "" { // path must be empty
		return utils.ResourceIDEqual(s.ResourceId, r.GetResourceId())
	}
	return false
}

func checkShareRef(s *collaboration.Share, ref *collaboration.ShareReference) bool {
	if ref.GetId() != nil {
		return ref.GetId().OpaqueId == s.Id.OpaqueId
	}
	if key := ref.GetKey(); key != nil {
		return (utils.UserEqual(key.Owner, s.Owner) || utils.UserEqual(key.Owner, s.Creator)) &&
			utils.ResourceIDEqual(key.ResourceId, s.ResourceId) && utils.GranteeEqual(key.Grantee, s.Grantee)
	}
	return false
}
func checkShare(s1 *collaboration.Share, s2 *collaboration.Share) bool {
	if s2.GetId() != nil {
		return s2.GetId().OpaqueId == s1.Id.OpaqueId
	}
	return false
}

func checkSharePath(path string) bool {
	paths := []string{
		"/ocs/v2.php/apps/files_sharing/api/v1/shares",
		"/ocs/v1.php/apps/files_sharing/api/v1/shares",
		"/remote.php/webdav",
		"/webdav",
		"/remote.php/dav/files",
		"/dav/files",
	}
	for _, p := range paths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// AddShareScope adds the scope to allow access to a user/group share and
// the shared resource.
func AddShareScope(share *collaboration.Share, role authpb.Role, scopes map[string]*authpb.Scope) (map[string]*authpb.Scope, error) {
	// Create a new "scope share" to only expose the required fields to the scope.
	scopeShare := &collaboration.Share{Id: share.Id, Owner: share.Owner, Creator: share.Creator, ResourceId: share.ResourceId}

	val, err := utils.MarshalProtoV1ToJSON(scopeShare)
	if err != nil {
		return nil, err
	}
	if scopes == nil {
		scopes = make(map[string]*authpb.Scope)
	}
	scopes["share:"+share.Id.OpaqueId] = &authpb.Scope{
		Resource: &types.OpaqueEntry{
			Decoder: "json",
			Value:   val,
		},
		Role: role,
	}
	return scopes, nil
}
