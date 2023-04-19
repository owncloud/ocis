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

	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
)

func receivedShareScope(_ context.Context, scope *authpb.Scope, resource interface{}, logger *zerolog.Logger) (bool, error) {
	var share collaboration.ReceivedShare
	err := utils.UnmarshalJSONToProtoV1(scope.Resource.Value, &share)
	if err != nil {
		return false, err
	}

	switch v := resource.(type) {
	case *collaboration.GetReceivedShareRequest:
		return checkShareRef(share.Share, v.GetRef()), nil
	case *collaboration.UpdateReceivedShareRequest:
		return checkShare(share.Share, v.GetShare().GetShare()), nil
	case string:
		return checkSharePath(v) || checkResourcePath(v), nil
	}

	msg := fmt.Sprintf("resource type assertion failed: %+v", resource)
	logger.Debug().Str("scope", "receivedShareScope").Msg(msg)
	return false, errtypes.InternalError(msg)
}

// AddReceivedShareScope adds the scope to allow access to a received user/group share and
// the shared resource.
func AddReceivedShareScope(share *collaboration.ReceivedShare, role authpb.Role, scopes map[string]*authpb.Scope) (map[string]*authpb.Scope, error) {
	// Create a new "scope share" to only expose the required fields to the scope.
	scopeShare := &collaboration.Share{Id: share.Share.Id, Owner: share.Share.Owner, Creator: share.Share.Creator, ResourceId: share.Share.ResourceId}

	val, err := utils.MarshalProtoV1ToJSON(&collaboration.ReceivedShare{Share: scopeShare})
	if err != nil {
		return nil, err
	}
	if scopes == nil {
		scopes = make(map[string]*authpb.Scope)
	}
	scopes["receivedshare:"+share.Share.Id.OpaqueId] = &authpb.Scope{
		Resource: &types.OpaqueEntry{
			Decoder: "json",
			Value:   val,
		},
		Role: role,
	}
	return scopes, nil
}
