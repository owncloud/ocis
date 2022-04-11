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

package index

import (
	"context"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/owncloud/ocis/search/pkg/search"

	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

type Index struct {
	bleveIndex bleve.Index
}

type Entity struct {
	RootID string
	Path   string
	ID     string

	Name string
	Size uint64
}

func NewPersisted(path string) (*Index, error) {
	bi, err := bleve.New(path, BuildMapping())
	if err != nil {
		return nil, err
	}
	return &Index{
		bleveIndex: bi,
	}, nil
}

func New(bleveIndex bleve.Index) (*Index, error) {
	return &Index{
		bleveIndex: bleveIndex,
	}, nil
}

func (i *Index) Add(ref *sprovider.Reference, ri *sprovider.ResourceInfo) error {
	entity := toEntity(ref, ri)
	return i.bleveIndex.Index(entity.ID, entity)
}

func (i *Index) Remove(ri *sprovider.ResourceInfo) error {
	return i.bleveIndex.Delete(ri.Id.GetStorageId() + ":" + ri.Id.GetOpaqueId())
}

func (i *Index) Search(ctx context.Context, req *search.SearchIndexRequest) (*search.SearchIndexResult, error) {
	query := bleve.NewConjunctionQuery(
		bleve.NewQueryStringQuery(req.Query),
		bleve.NewQueryStringQuery("Path:"+req.Reference.Path+"*"), // Limit search to this directory in the space
	)
	bleveReq := bleve.NewSearchRequest(query)
	bleveReq.Fields = []string{"*"}
	res, err := i.bleveIndex.Search(bleveReq)
	if err != nil {
		return nil, err
	}

	matches := []search.Match{}
	for _, h := range res.Hits {
		match, err := fromFields(h.Fields)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return &search.SearchIndexResult{
		Matches: matches,
	}, nil
}

func BuildMapping() mapping.IndexMapping {
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultAnalyzer = keyword.Name
	return indexMapping
}

func toEntity(ref *sprovider.Reference, ri *sprovider.ResourceInfo) *Entity {
	return &Entity{
		RootID: ref.ResourceId.GetStorageId() + ":" + ref.ResourceId.GetOpaqueId(),
		Path:   ref.Path,
		ID:     ri.Id.GetStorageId() + ":" + ri.Id.GetOpaqueId(),
		Name:   ri.Path,
		Size:   ri.Size,
	}
}

func fromFields(fields map[string]interface{}) (search.Match, error) {
	rootIDParts := strings.SplitN(fields["RootID"].(string), ":", 2)
	IDParts := strings.SplitN(fields["ID"].(string), ":", 2)

	return search.Match{
		Reference: &sprovider.Reference{
			ResourceId: &sprovider.ResourceId{
				StorageId: rootIDParts[0],
				OpaqueId:  rootIDParts[1],
			},
			Path: fields["Path"].(string),
		},
		Info: &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: IDParts[0],
				OpaqueId:  IDParts[1],
			},
			Path: fields["Name"].(string),
			Size: uint64(fields["Size"].(float64)),
		},
	}, nil
}
