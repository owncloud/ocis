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

	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

//go:generate mockery --name=ProviderClient
//go:generate mockery --name=IndexClient

type SearchRequest struct {
	Query string
}

type Match struct {
	Reference *sprovider.Reference
	Info      *sprovider.ResourceInfo
}

type SearchResult struct {
	Matches []Match
}

type SearchIndexRequest struct {
	// Reference is not a list because the Path is used as a filter which is
	// cut off in the matches by the provider. Multiple paths would not be
	// distinguishable.
	Reference *sprovider.Reference
	Query     string
}

type SearchIndexResult struct {
	Matches []Match
}

type ProviderClient interface {
	Search(ctx context.Context, req *SearchRequest) (*SearchResult, error)
}

type IndexClient interface {
	Search(ctx context.Context, req *SearchIndexRequest) (*SearchIndexResult, error)
}
