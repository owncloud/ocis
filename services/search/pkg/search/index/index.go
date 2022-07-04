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
	"math"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/v2/mapping"
	"google.golang.org/protobuf/types/known/timestamppb"

	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/utils"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
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
	mapping, err := BuildMapping()
	if err != nil {
		return nil, err
	}
	bi, err := bleve.New(path, mapping)
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

// Delete marks an entity from the index as deleten (still keeping it around)
func (i *Index) Delete(id *sprovider.ResourceId) error {
	return i.markAsDeleted(idToBleveId(id), true)
}

// Restore marks an entity from the index as not being deleted
func (i *Index) Restore(id *sprovider.ResourceId) error {
	return i.markAsDeleted(idToBleveId(id), false)
}

func (i *Index) markAsDeleted(id string, deleted bool) error {
	doc, err := i.updateEntity(id, func(doc *indexDocument) {
		doc.Deleted = deleted
	})
	if err != nil {
		return err
	}

	if doc.Type == uint64(sprovider.ResourceType_RESOURCE_TYPE_CONTAINER) {
		query := bleve.NewConjunctionQuery(
			bleve.NewQueryStringQuery("RootID:"+doc.RootID),
			bleve.NewQueryStringQuery("Path:"+queryEscape(doc.Path+"/*")),
		)
		bleveReq := bleve.NewSearchRequest(query)
		bleveReq.Size = math.MaxInt
		bleveReq.Fields = []string{"*"}
		res, err := i.bleveIndex.Search(bleveReq)
		if err != nil {
			return err
		}

		for _, h := range res.Hits {
			_, err := i.updateEntity(h.ID, func(doc *indexDocument) {
				doc.Deleted = deleted
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Index) updateEntity(id string, mutateFunc func(doc *indexDocument)) (*indexDocument, error) {
	doc, err := i.getEntity(id)
	if err != nil {
		return nil, err
	}
	mutateFunc(doc)
	err = i.bleveIndex.Index(doc.ID, doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (i *Index) getEntity(id string) (*indexDocument, error) {
	req := bleve.NewSearchRequest(bleve.NewDocIDQuery([]string{id}))
	req.Fields = []string{"*"}
	res, err := i.bleveIndex.Search(req)
	if err != nil {
		return nil, err
	}
	if res.Hits.Len() == 0 {
		return nil, errors.New("entity not found")
	}
	return fieldsToEntity(res.Hits[0].Fields), nil
}

// Purge removes an entity from the index
func (i *Index) Purge(id *sprovider.ResourceId) error {
	return i.bleveIndex.Delete(idToBleveId(id))
}

// Move update the path of an entry and all its children
func (i *Index) Move(id *sprovider.ResourceId, fullPath string) error {
	bleveId := idToBleveId(id)
	doc, err := i.getEntity(bleveId)
	if err != nil {
		return err
	}
	oldName := doc.Path
	newName := utils.MakeRelativePath(fullPath)

	doc, err = i.updateEntity(bleveId, func(doc *indexDocument) {
		doc.Path = newName
		doc.Name = path.Base(newName)
	})
	if err != nil {
		return err
	}

	if doc.Type == uint64(sprovider.ResourceType_RESOURCE_TYPE_CONTAINER) {
		query := bleve.NewConjunctionQuery(
			bleve.NewQueryStringQuery("RootID:"+doc.RootID),
			bleve.NewQueryStringQuery("Path:"+queryEscape(oldName+"/*")),
		)
		bleveReq := bleve.NewSearchRequest(query)
		bleveReq.Size = math.MaxInt
		bleveReq.Fields = []string{"*"}
		res, err := i.bleveIndex.Search(bleveReq)
		if err != nil {
			return err
		}

		for _, h := range res.Hits {
			_, err := i.updateEntity(h.ID, func(doc *indexDocument) {
				doc.Path = strings.Replace(doc.Path, oldName, newName, 1)
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Search searches the index according to the criteria specified in the given SearchIndexRequest
func (i *Index) Search(ctx context.Context, req *searchsvc.SearchIndexRequest) (*searchsvc.SearchIndexResponse, error) {
	deletedQuery := bleve.NewBoolFieldQuery(false)
	deletedQuery.SetField("Deleted")
	query := bleve.NewConjunctionQuery(
		bleve.NewQueryStringQuery(req.Query),
		deletedQuery, // Skip documents that have been marked as deleted
		bleve.NewQueryStringQuery("RootID:"+req.Ref.ResourceId.StorageId+"!"+req.Ref.ResourceId.OpaqueId),        // Limit search to the space
		bleve.NewQueryStringQuery("Path:"+queryEscape(utils.MakeRelativePath(path.Join(req.Ref.Path, "/"))+"*")), // Limit search to this directory in the space
	)
	bleveReq := bleve.NewSearchRequest(query)
	bleveReq.Size = 200
	if req.PageSize > 0 {
		bleveReq.Size = int(req.PageSize)
	}
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
func BuildMapping() (mapping.IndexMapping, error) {
	nameMapping := bleve.NewTextFieldMapping()
	nameMapping.Analyzer = "lowercaseKeyword"

	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("Name", nameMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultAnalyzer = keyword.Name
	indexMapping.DefaultMapping = docMapping
	err := indexMapping.AddCustomAnalyzer("lowercaseKeyword",
		map[string]interface{}{
			"type":      custom.Name,
			"tokenizer": single.Name,
			"token_filters": []string{
				lowercase.Name,
			},
		})
	if err != nil {
		return nil, err
	}

	return indexMapping, nil
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
		Deleted:  fields["Deleted"].(bool),
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

func queryEscape(s string) string {
	re := regexp.MustCompile(`([` + regexp.QuoteMeta(`+=&|><!(){}[]^\"~*?:\/`) + `\-\s])`)
	return re.ReplaceAllString(s, "\\$1")
}
