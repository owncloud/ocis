package engine

import (
	"context"
	"errors"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	storageProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	searchMessage "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchService "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Bleve represents a search engine which utilizes bleve to search and store resources.
type Bleve struct {
	index bleve.Index
}

// NewBleveIndex returns a new bleve index
// given path must exist.
func NewBleveIndex(root string) (bleve.Index, error) {
	destination := filepath.Join(root, "bleve")
	index, err := bleve.Open(destination)
	if errors.Is(bleve.ErrorIndexPathDoesNotExist, err) {
		m, err := BuildBleveMapping()
		if err != nil {
			return nil, err
		}
		index, err = bleve.New(destination, m)
		if err != nil {
			return nil, err
		}

		return index, nil
	}

	return index, err
}

// NewBleveEngine creates a new Bleve instance
func NewBleveEngine(index bleve.Index) *Bleve {
	return &Bleve{
		index: index,
	}
}

// BuildBleveMapping builds a bleve index mapping which can be used for indexing
func BuildBleveMapping() (mapping.IndexMapping, error) {
	lowercaseMapping := bleve.NewTextFieldMapping()
	lowercaseMapping.Analyzer = "lowercaseKeyword"

	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("Name", lowercaseMapping)
	docMapping.AddFieldMappingsAt("Tags", lowercaseMapping)
	docMapping.AddFieldMappingsAt("Content", lowercaseMapping)

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

// Search executes a search request operation within the index.
// Returns a SearchIndexResponse object or an error.
func (b *Bleve) Search(_ context.Context, sir *searchService.SearchIndexRequest) (*searchService.SearchIndexResponse, error) {
	q := bleve.NewConjunctionQuery(
		&query.QueryStringQuery{
			Query: b.buildQuery(sir.Query),
		},
		// Skip documents that have been marked as deleted
		&query.BoolFieldQuery{
			Bool:     false,
			FieldVal: "Deleted",
		},
		&query.QueryStringQuery{
			Query: fmt.Sprintf(
				"RootID:%s",
				storagespace.FormatResourceID(
					storageProvider.ResourceId{
						StorageId: sir.Ref.GetResourceId().GetStorageId(),
						SpaceId:   sir.Ref.GetResourceId().GetSpaceId(),
						OpaqueId:  sir.Ref.GetResourceId().GetOpaqueId(),
					},
				),
			),
		},
		// Limit search to this directory in the space
		&query.QueryStringQuery{
			Query: fmt.Sprintf("Path:%s*", escapeQuery(utils.MakeRelativePath(path.Join(sir.Ref.Path, "/")))),
		},
	)

	bleveReq := bleve.NewSearchRequest(q)

	switch {
	case sir.PageSize == -1:
		bleveReq.Size = math.MaxInt
	case sir.PageSize == 0:
		bleveReq.Size = 200
	default:
		bleveReq.Size = int(sir.PageSize)
	}

	bleveReq.Fields = []string{"*"}
	res, err := b.index.Search(bleveReq)
	if err != nil {
		return nil, err
	}

	matches := []*searchMessage.Match{}
	for _, hit := range res.Hits {
		rootID, err := storagespace.ParseID(getValue[string](hit.Fields, "RootID"))
		if err != nil {
			return nil, err
		}

		rID, err := storagespace.ParseID(getValue[string](hit.Fields, "ID"))
		if err != nil {
			return nil, err
		}

		match := &searchMessage.Match{
			Score: float32(hit.Score),
			Entity: &searchMessage.Entity{
				Ref: &searchMessage.Reference{
					ResourceId: resourceIDtoSearchID(rootID),
					Path:       getValue[string](hit.Fields, "Path"),
				},
				Id:       resourceIDtoSearchID(rID),
				Name:     getValue[string](hit.Fields, "Name"),
				Size:     uint64(getValue[float64](hit.Fields, "Size")),
				Type:     uint64(getValue[float64](hit.Fields, "Type")),
				MimeType: getValue[string](hit.Fields, "MimeType"),
				Deleted:  getValue[bool](hit.Fields, "Deleted"),
				Tags:     getSliceValue[string](hit.Fields, "Tags"),
			},
		}

		if mtime, err := time.Parse(time.RFC3339, getValue[string](hit.Fields, "Mtime")); err == nil {
			match.Entity.LastModifiedTime = &timestamppb.Timestamp{Seconds: mtime.Unix(), Nanos: int32(mtime.Nanosecond())}
		}

		matches = append(matches, match)
	}

	return &searchService.SearchIndexResponse{
		Matches:      matches,
		TotalMatches: int32(res.Total),
	}, nil
}

// Upsert indexes or stores Resource data fields.
func (b *Bleve) Upsert(id string, r Resource) error {
	return b.index.Index(id, r)
}

// Move updates the resource location and all of its necessary fields.
func (b *Bleve) Move(id string, target string) error {
	r, err := b.getResource(id)
	if err != nil {
		return err
	}
	oldName := r.Path
	newName := utils.MakeRelativePath(target)

	r, err = b.updateEntity(id, func(r *Resource) {
		r.Path = newName
		r.Name = path.Base(newName)
	})
	if err != nil {
		return err
	}

	if r.Type == uint64(storageProvider.ResourceType_RESOURCE_TYPE_CONTAINER) {
		q := bleve.NewConjunctionQuery(
			bleve.NewQueryStringQuery("RootID:"+r.RootID),
			bleve.NewQueryStringQuery("Path:"+escapeQuery(oldName+"/*")),
		)
		bleveReq := bleve.NewSearchRequest(q)
		bleveReq.Size = math.MaxInt
		bleveReq.Fields = []string{"*"}
		res, err := b.index.Search(bleveReq)
		if err != nil {
			return err
		}

		for _, h := range res.Hits {
			_, err := b.updateEntity(h.ID, func(r *Resource) {
				r.Path = strings.Replace(r.Path, oldName, newName, 1)
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete marks the resource as deleted.
// The resource object will stay in the bleve index,
// instead of removing the resource it just marks it as deleted!
// can be undone
func (b *Bleve) Delete(id string) error {
	return b.setDeleted(id, true)
}

// Restore is the counterpart to Delete.
// It restores the resource which makes it available again.
func (b *Bleve) Restore(id string) error {
	return b.setDeleted(id, false)
}

// Purge removes a resource from the index, irreversible operation.
func (b *Bleve) Purge(id string) error {
	return b.index.Delete(id)
}

// DocCount returns the number of resources in the index.
func (b *Bleve) DocCount() (uint64, error) {
	return b.index.DocCount()
}

func (b *Bleve) getResource(id string) (*Resource, error) {
	req := bleve.NewSearchRequest(bleve.NewDocIDQuery([]string{id}))
	req.Fields = []string{"*"}
	res, err := b.index.Search(req)
	if err != nil {
		return nil, err
	}
	if res.Hits.Len() == 0 {
		return nil, errors.New("entity not found")
	}

	fields := res.Hits[0].Fields

	return &Resource{
		ID:      getValue[string](fields, "ID"),
		RootID:  getValue[string](fields, "RootID"),
		Path:    getValue[string](fields, "Path"),
		Type:    uint64(getValue[float64](fields, "Type")),
		Deleted: getValue[bool](fields, "Deleted"),
		Document: content.Document{
			Name:     getValue[string](fields, "Name"),
			Size:     uint64(getValue[float64](fields, "Size")),
			Mtime:    getValue[string](fields, "Mtime"),
			MimeType: getValue[string](fields, "MimeType"),
			Tags:     getSliceValue[string](fields, "Tags"),
		},
	}, nil
}

func (b *Bleve) updateEntity(id string, mutateFunc func(r *Resource)) (*Resource, error) {
	it, err := b.getResource(id)
	if err != nil {
		return nil, err
	}

	mutateFunc(it)

	return it, b.index.Index(it.ID, it)
}

func (b *Bleve) setDeleted(id string, deleted bool) error {
	it, err := b.updateEntity(id, func(r *Resource) {
		r.Deleted = deleted
	})
	if err != nil {
		return err
	}

	if it.Type == uint64(storageProvider.ResourceType_RESOURCE_TYPE_CONTAINER) {
		q := bleve.NewConjunctionQuery(
			bleve.NewQueryStringQuery("RootID:"+it.RootID),
			bleve.NewQueryStringQuery("Path:"+escapeQuery(it.Path+"/*")),
		)
		bleveReq := bleve.NewSearchRequest(q)
		bleveReq.Size = math.MaxInt
		bleveReq.Fields = []string{"*"}
		res, err := b.index.Search(bleveReq)
		if err != nil {
			return err
		}

		for _, h := range res.Hits {
			_, err := b.updateEntity(h.ID, func(r *Resource) {
				r.Deleted = deleted
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bleve) buildQuery(qs string) string {
	fields := []string{"RootID", "Path", "ID", "Name", "Size", "Mtime", "MimeType", "Type", "Content", "Title", "Tags"}
	for _, field := range fields {
		qs = strings.ReplaceAll(qs, strings.ToLower(field)+":", field+":")
	}

	if strings.Contains(qs, ":") {
		return qs // Sophisticated field based search
	}

	// this is a basic filename search
	return "Name:*" + strings.ReplaceAll(strings.ToLower(qs), " ", `\ `) + "*"
}
