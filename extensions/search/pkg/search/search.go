// Copyright 2018-2022 CERN
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

package search

import (
	"context"

	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"
)

//go:generate mockery --name=ProviderClient
//go:generate mockery --name=IndexClient

// ProviderClient is the interface to the search provider service
type ProviderClient interface {
	Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error)
	IndexSpace(ctx context.Context, req *searchsvc.IndexSpaceRequest) (*searchsvc.IndexSpaceResponse, error)
}

// IndexClient is the interface to the search index
type IndexClient interface {
	Search(ctx context.Context, req *searchsvc.SearchIndexRequest) (*searchsvc.SearchIndexResponse, error)
	Add(ref *providerv1beta1.Reference, ri *providerv1beta1.ResourceInfo) error
	Delete(ri *providerv1beta1.ResourceId) error
	Restore(ri *providerv1beta1.ResourceId) error
	Purge(ri *providerv1beta1.ResourceId) error
	DocCount() (uint64, error)
}
