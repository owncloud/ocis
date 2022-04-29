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
	"errors"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/mapping"
	"google.golang.org/protobuf/types/known/timestamppb"

	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	searchmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"
)

type indexDocument struct {
	RootID string
	Path   string
	ID     string

	Name     string
	Size     uint64
	Mtime    string
	MimeType string
	Type     uint64

	Deleted bool
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

// DocCount returns the number of elemenst in the index
func (i *Index) DocCount() (uint64, error) {
	return i.bleveIndex.DocCount()
}

// Add adds a new entity to the Index
func (i *Index) Add(ref *sprovider.Reference, ri *sprovider.ResourceInfo) error {
	entity := toEntity(ref, ri)
	return i.bleveIndex.Index(idToBleveId(ri.Id), entity)
}

// Delete marks an entity from the index as delete (still keeping it around)
func (i *Index) Delete(id *sprovider.ResourceId) error {
	return i.markAsDeleted(idToBleveId(id))
}

func (i *Index) markAsDeleted(id string) error {
	req := bleve.NewSearchRequest(bleve.NewDocIDQuery([]string{id}))
	req.Fields = []string{"*"}
	res, err := i.bleveIndex.Search(req)
	if err != nil {
		return err
	}
	if res.Hits.Len() == 0 {
		return errors.New("entity not found")
	}
	entity := fieldsToEntity(res.Hits[0].Fields)
	entity.Deleted = true

	if entity.Type == uint64(sprovider.ResourceType_RESOURCE_TYPE_CONTAINER) {
		query := bleve.NewConjunctionQuery(
			bleve.NewQueryStringQuery("RootID:"+entity.RootID),
			bleve.NewQueryStringQuery("Path:"+entity.Path+"/*"),
		)
		bleveReq := bleve.NewSearchRequest(query)
		bleveReq.Fields = []string{"*"}
		res, err := i.bleveIndex.Search(bleveReq)
		if err != nil {
			return err
		}

		for _, h := range res.Hits {
			i.markAsDeleted(h.ID)
		}
	}

	return i.bleveIndex.Index(entity.ID, entity)
}

// Purge removes an entity from the index
func (i *Index) Purge(id *sprovider.ResourceId) error {
	return i.bleveIndex.Delete(idToBleveId(id))
}

// Search searches the index according to the criteria specified in the given SearchIndexRequest
func (i *Index) Search(ctx context.Context, req *searchsvc.SearchIndexRequest) (*searchsvc.SearchIndexResponse, error) {
	deletedQuery := bleve.NewBoolFieldQuery(false)
	deletedQuery.SetField("Deleted")
	query := bleve.NewConjunctionQuery(
		bleve.NewQueryStringQuery("Name:"+req.Query),
		deletedQuery, // Skip documents that have been marked as deleted
		bleve.NewQueryStringQuery("RootID:"+req.Ref.ResourceId.StorageId+"!"+req.Ref.ResourceId.OpaqueId), // Limit search to the space
		bleve.NewQueryStringQuery("Path:"+req.Ref.Path+"*"),                                               // Limit search to this directory in the space
	)
	bleveReq := bleve.NewSearchRequest(query)
	bleveReq.Size = 200
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
	doc := &indexDocument{
		RootID:   idToBleveId(ref.ResourceId),
		Path:     ref.Path,
		ID:       idToBleveId(ri.Id),
		Name:     ri.Path,
		Size:     ri.Size,
		MimeType: ri.MimeType,
		Type:     uint64(ri.Type),
		Deleted:  false,
	}

	if ri.Mtime != nil {
		doc.Mtime = time.Unix(int64(ri.Mtime.Seconds), int64(ri.Mtime.Nanos)).UTC().Format(time.RFC3339)
	}

	return doc
}

func fieldsToEntity(fields map[string]interface{}) *indexDocument {
	doc := &indexDocument{
		RootID:   fields["RootID"].(string),
		Path:     fields["Path"].(string),
		ID:       fields["ID"].(string),
		Name:     fields["Name"].(string),
		Size:     uint64(fields["Size"].(float64)),
		Mtime:    fields["Mtime"].(string),
		MimeType: fields["MimeType"].(string),
		Type:     uint64(fields["Type"].(float64)),
	}
	return doc
}

func fromFields(fields map[string]interface{}) (*searchmsg.Match, error) {
	rootIDParts := strings.SplitN(fields["RootID"].(string), "!", 2)
	IDParts := strings.SplitN(fields["ID"].(string), "!", 2)

	match := &searchmsg.Match{
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
			Name:     fields["Name"].(string),
			Size:     uint64(fields["Size"].(float64)),
			Type:     uint64(fields["Type"].(float64)),
			MimeType: fields["MimeType"].(string),
			Deleted:  fields["Deleted"].(bool),
		},
	}

	if mtime, err := time.Parse(time.RFC3339, fields["Mtime"].(string)); err == nil {
		match.Entity.LastModifiedTime = &timestamppb.Timestamp{Seconds: mtime.Unix(), Nanos: int32(mtime.Nanosecond())}
	}

	return match, nil
}

func idToBleveId(id *sprovider.ResourceId) string {
	if id == nil {
		return ""
	}
	return id.StorageId + "!" + id.OpaqueId
}
