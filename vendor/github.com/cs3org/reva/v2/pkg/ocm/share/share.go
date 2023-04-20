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

package share

import (
	"context"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"google.golang.org/genproto/protobuf/field_mask"
)

// Manager is the interface that manipulates the OCM shares.
type Manager interface {
	// Create a new share in fn with the given acl.
	Share(ctx context.Context, md *provider.ResourceId, g *ocm.ShareGrant, name string,
		pi *ocmprovider.ProviderInfo, pm string, owner *userpb.UserId, token string, st ocm.Share_ShareType) (*ocm.Share, error)

	// GetShare gets the information for a share by the given ref.
	GetShare(ctx context.Context, ref *ocm.ShareReference) (*ocm.Share, error)

	// Unshare deletes the share pointed by ref.
	Unshare(ctx context.Context, ref *ocm.ShareReference) error

	// UpdateShare updates the mode of the given share.
	UpdateShare(ctx context.Context, ref *ocm.ShareReference, p *ocm.SharePermissions) (*ocm.Share, error)

	// ListShares returns the shares created by the user. If md is provided is not nil,
	// it returns only shares attached to the given resource.
	ListShares(ctx context.Context, filters []*ocm.ListOCMSharesRequest_Filter) ([]*ocm.Share, error)

	// ListReceivedShares returns the list of shares the user has access.
	ListReceivedShares(ctx context.Context) ([]*ocm.ReceivedShare, error)

	// GetReceivedShare returns the information for a received share the user has access.
	GetReceivedShare(ctx context.Context, ref *ocm.ShareReference) (*ocm.ReceivedShare, error)

	// UpdateReceivedShare updates the received share with share state.
	UpdateReceivedShare(ctx context.Context, share *ocm.ReceivedShare, fieldMask *field_mask.FieldMask) (*ocm.ReceivedShare, error)
}

// ResourceIDFilter is an abstraction for creating filter by resource id.
func ResourceIDFilter(id *provider.ResourceId) *ocm.ListOCMSharesRequest_Filter {
	return &ocm.ListOCMSharesRequest_Filter{
		Type: ocm.ListOCMSharesRequest_Filter_TYPE_RESOURCE_ID,
		Term: &ocm.ListOCMSharesRequest_Filter_ResourceId{
			ResourceId: id,
		},
	}
}
