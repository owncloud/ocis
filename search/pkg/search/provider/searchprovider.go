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

package provider

import (
	"context"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/search/pkg/search"
)

type Provider struct {
	gwClient    gateway.GatewayAPIClient
	indexClient search.IndexClient
}

func New(gwClient gateway.GatewayAPIClient, indexClient search.IndexClient) *Provider {
	return &Provider{
		gwClient:    gwClient,
		indexClient: indexClient,
	}
}

func (p *Provider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResult, error) {
	if req.Query == "" {
		return nil, errtypes.PreconditionFailed("empty query provided")
	}

	listSpacesRes, err := p.gwClient.ListStorageSpaces(ctx, &providerv1beta1.ListStorageSpacesRequest{
		Opaque: &typesv1beta1.Opaque{Map: map[string]*typesv1beta1.OpaqueEntry{
			"path": {
				Decoder: "plain",
				Value:   []byte("/"),
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	matches := []search.Match{}
	for _, space := range listSpacesRes.StorageSpaces {
		pathPrefix := ""
		if space.SpaceType == "grant" {
			gpRes, err := p.gwClient.GetPath(ctx, &providerv1beta1.GetPathRequest{
				ResourceId: space.Root,
			})
			if err != nil {
				return nil, err
			}
			if gpRes.Status.Code != rpcv1beta1.Code_CODE_OK {
				return nil, errtypes.NewErrtypeFromStatus(gpRes.Status)
			}
			pathPrefix = utils.MakeRelativePath(gpRes.Path)
		}

		res, err := p.indexClient.Search(ctx, &search.SearchIndexRequest{
			Query: req.Query,
			Reference: &providerv1beta1.Reference{
				ResourceId: space.Root,
				Path:       pathPrefix,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, match := range res.Matches {
			if pathPrefix != "" {
				match.Reference.Path = utils.MakeRelativePath(strings.TrimPrefix(match.Reference.Path, pathPrefix))
			}
			matches = append(matches, match)
		}
	}

	return &search.SearchResult{Matches: matches}, nil
}
