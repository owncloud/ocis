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

	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	searchmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"
)

type indexDocument struct {
	RootID string
	Path   string
	ID     string

	Name string
	Size uint64
}

// Index represents a bleve based search index
type Index struct {
	bleveIndex bleve.Index
}

// NewPersisted returns a new instance of Index with the data being persisted in the given directory
func NewPersisted(path string) (*Index, error) {
	bi, err := bleve.New(path, BuildMapping())
	if err != nil {
		return nil, err
	}
	return New(bi)
}

// New returns a new instance of Index using the given bleve Index as the backend
func New(bleveIndex bleve.Index) (*Index, error) {
	return &Index{
		bleveIndex: bleveIndex,
	}, nil
}

// Add adds a new entity to the Index
func (i *Index) Add(ref *sprovider.Reference, ri *sprovider.ResourceInfo) error {
	entity := toEntity(ref, ri)
	return i.bleveIndex.Index(idToBleveId(ri.Id), entity)
}

// Remove removes an entity from the index
func (i *Index) Remove(ri *sprovider.ResourceInfo) error {
	return i.bleveIndex.Delete(idToBleveId(ri.Id))
}

// Search searches the index according to the criteria specified in the given SearchIndexRequest
func (i *Index) Search(ctx context.Context, req *searchsvc.SearchIndexRequest) (*searchsvc.SearchIndexResponse, error) {
	query := bleve.NewConjunctionQuery(
		bleve.NewQueryStringQuery("Name:"+req.Query),
		bleve.NewQueryStringQuery("RootID:"+req.Ref.ResourceId.StorageId+"!"+req.Ref.ResourceId.OpaqueId), // Limit search to the space
		bleve.NewQueryStringQuery("Path:"+req.Ref.Path+"*"),                                               // Limit search to this directory in the space
	)
	bleveReq := bleve.NewSearchRequest(query)
	bleveReq.Fields = []string{"*"}
	res, err := i.bleveIndex.Search(bleveReq)
	if err != nil {
		return nil, err
	}

	matches := []*searchmsg.Match{}
	for _, h := range res.Hits {
		match, err := fromFields(h.Fields)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return &searchsvc.SearchIndexResponse{
		Matches: matches,
	}, nil
}

// BuildMapping builds a bleve index mapping which can be used for indexing
func BuildMapping() mapping.IndexMapping {
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultAnalyzer = keyword.Name
	return indexMapping
}

func toEntity(ref *sprovider.Reference, ri *sprovider.ResourceInfo) *indexDocument {
	return &indexDocument{
		RootID: idToBleveId(ref.ResourceId),
		Path:   ref.Path,
		ID:     idToBleveId(ri.Id),
		Name:   ri.Path,
		Size:   ri.Size,
	}
}

func fromFields(fields map[string]interface{}) (*searchmsg.Match, error) {
	rootIDParts := strings.SplitN(fields["RootID"].(string), "!", 2)
	IDParts := strings.SplitN(fields["ID"].(string), "!", 2)

	return &searchmsg.Match{
		Entity: &searchmsg.Entity{
			Ref: &searchmsg.Reference{
				ResourceId: &searchmsg.ResourceID{
					StorageId: rootIDParts[0],
					OpaqueId:  rootIDParts[1],
				},
				Path: fields["Path"].(string),
			},
			Id: &searchmsg.ResourceID{
				StorageId: IDParts[0],
				OpaqueId:  IDParts[1],
			},
			Name: fields["Name"].(string),
			Size: uint64(fields["Size"].(float64)),
		},
	}, nil
}

func idToBleveId(id *sprovider.ResourceId) string {
	return id.StorageId + "!" + id.OpaqueId
}
