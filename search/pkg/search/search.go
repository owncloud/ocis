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

// SearchRequest represents a search request from a user to the search provider
type SearchRequest struct {
	Query string
}

// Match holds the information of a matched resource in a search
type Match struct {
	Reference *sprovider.Reference
	Info      *sprovider.ResourceInfo
}

// SearchResult contains the matches being returned for a search
type SearchResult struct {
	Matches []Match
}

// ProviderClient is the interface to the search provider service
type ProviderClient interface {
	Search(ctx context.Context, req *SearchRequest) (*SearchResult, error)
}

// SearchIndexRequest represents a search request to the index
type SearchIndexRequest struct {
	Reference *sprovider.Reference
	Query     string
}

// SearchResult contains the matches in the index being returned for a search
type SearchIndexResult struct {
	Matches []Match
}

// IndexClient is the interface to the search index
type IndexClient interface {
	Search(ctx context.Context, req *SearchIndexRequest) (*SearchIndexResult, error)
}
