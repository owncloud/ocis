package engine

import (
	"context"
	"errors"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/v2/mapping"
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

type Bleve struct {
	index bleve.Index
}

func NewBleveIndex(root string) (bleve.Index, error) {
	destination := filepath.Join(root, "index.bleve")
	index, err := bleve.Open(destination)
	if err != nil {
		m, err := BuildBleveMapping()
		if err != nil {
			return nil, err
		}
		index, err = bleve.New(destination, m)
		if err != nil {
			return nil, err
		}
	}

	return index, nil
}

func NewBleveEngine(index bleve.Index) *Bleve {
	return &Bleve{
		index: index,
	}
}

// BuildBleveMapping builds a bleve index mapping which can be used for indexing
func BuildBleveMapping() (mapping.IndexMapping, error) {
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

func (b *Bleve) Search(ctx context.Context, sir *searchService.SearchIndexRequest) (*searchService.SearchIndexResponse, error) {
	deletedQuery := bleve.NewBoolFieldQuery(false)
	deletedQuery.SetField("Deleted")
	query := bleve.NewConjunctionQuery(
		bleve.NewQueryStringQuery(b.buildQuery(sir.Query)),
		deletedQuery, // Skip documents that have been marked as deleted
		bleve.NewQueryStringQuery("RootID:"+storagespace.FormatResourceID(
			storageProvider.ResourceId{
				StorageId: sir.Ref.GetResourceId().GetStorageId(),
				SpaceId:   sir.Ref.GetResourceId().GetSpaceId(),
				OpaqueId:  sir.Ref.GetResourceId().GetOpaqueId(),
			})), // Limit search to the space
		bleve.NewQueryStringQuery("Path:"+escapeQuery(utils.MakeRelativePath(path.Join(sir.Ref.Path, "/"))+"*")), // Limit search to this directory in the space
	)

	bleveReq := bleve.NewSearchRequest(query)
	bleveReq.Size = 200
	if sir.PageSize > 0 {
		bleveReq.Size = int(sir.PageSize)
	}
	bleveReq.Fields = []string{"*"}
	res, err := b.index.Search(bleveReq)
	if err != nil {
		return nil, err
	}

	matches := []*searchMessage.Match{}
	for _, hit := range res.Hits {
		rootID, err := storagespace.ParseID(hit.Fields["RootID"].(string))
		if err != nil {
			return nil, err
		}

		rID, err := storagespace.ParseID(hit.Fields["ID"].(string))
		if err != nil {
			return nil, err
		}

		match := &searchMessage.Match{
			Score: float32(hit.Score),
			Entity: &searchMessage.Entity{
				Ref: &searchMessage.Reference{
					ResourceId: resourceIDtoSearchID(rootID),
					Path:       hit.Fields["Path"].(string),
				},
				Id:       resourceIDtoSearchID(rID),
				Name:     hit.Fields["Name"].(string),
				Size:     uint64(hit.Fields["Size"].(float64)),
				Type:     uint64(hit.Fields["Type"].(float64)),
				MimeType: hit.Fields["MimeType"].(string),
				Deleted:  hit.Fields["Deleted"].(bool),
			},
		}

		if mtime, err := time.Parse(time.RFC3339, hit.Fields["Mtime"].(string)); err == nil {
			match.Entity.LastModifiedTime = &timestamppb.Timestamp{Seconds: mtime.Unix(), Nanos: int32(mtime.Nanosecond())}
		}

		matches = append(matches, match)
	}

	return &searchService.SearchIndexResponse{
		Matches:      matches,
		TotalMatches: int32(res.Total),
	}, nil
}

func (b *Bleve) Upsert(id string, ent Entity) error {
	return b.index.Index(id, ent)
}

func (b *Bleve) Move(id string, target string) error {
	ent, err := b.getEntity(id)
	if err != nil {
		return err
	}
	oldName := ent.Path
	newName := utils.MakeRelativePath(target)

	ent, err = b.updateEntity(id, func(ent *Entity) {
		ent.Path = newName
		ent.Name = path.Base(newName)
	})
	if err != nil {
		return err
	}

	if ent.Type == uint64(storageProvider.ResourceType_RESOURCE_TYPE_CONTAINER) {
		query := bleve.NewConjunctionQuery(
			bleve.NewQueryStringQuery("RootID:"+ent.RootID),
			bleve.NewQueryStringQuery("Path:"+escapeQuery(oldName+"/*")),
		)
		bleveReq := bleve.NewSearchRequest(query)
		bleveReq.Size = math.MaxInt
		bleveReq.Fields = []string{"*"}
		res, err := b.index.Search(bleveReq)
		if err != nil {
			return err
		}

		for _, h := range res.Hits {
			_, err := b.updateEntity(h.ID, func(ent *Entity) {
				ent.Path = strings.Replace(ent.Path, oldName, newName, 1)
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bleve) Delete(id string) error {
	return b.setDeleted(id, true)
}

func (b *Bleve) Restore(id string) error {
	return b.setDeleted(id, false)
}

func (b *Bleve) Purge(id string) error {
	return b.index.Delete(id)
}

func (b *Bleve) DocCount() (uint64, error) {
	return b.index.DocCount()
}

func (b *Bleve) getEntity(id string) (*Entity, error) {
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

	return &Entity{
		ID:      fields["ID"].(string),
		RootID:  fields["RootID"].(string),
		Path:    fields["Path"].(string),
		Type:    uint64(fields["Type"].(float64)),
		Deleted: fields["Deleted"].(bool),
		Document: content.Document{
			Name:     fields["Name"].(string),
			Size:     uint64(fields["Size"].(float64)),
			Mtime:    fields["Mtime"].(string),
			MimeType: fields["MimeType"].(string),
		},
	}, nil
}

func (b *Bleve) updateEntity(id string, mutateFunc func(ent *Entity)) (*Entity, error) {
	it, err := b.getEntity(id)
	if err != nil {
		return nil, err
	}

	mutateFunc(it)

	err = b.index.Index(it.ID, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func (b *Bleve) setDeleted(id string, deleted bool) error {
	it, err := b.updateEntity(id, func(ent *Entity) {
		ent.Deleted = deleted
	})
	if err != nil {
		return err
	}

	if it.Type == uint64(storageProvider.ResourceType_RESOURCE_TYPE_CONTAINER) {
		query := bleve.NewConjunctionQuery(
			bleve.NewQueryStringQuery("RootID:"+it.RootID),
			bleve.NewQueryStringQuery("Path:"+escapeQuery(it.Path+"/*")),
		)
		bleveReq := bleve.NewSearchRequest(query)
		bleveReq.Size = math.MaxInt
		bleveReq.Fields = []string{"*"}
		res, err := b.index.Search(bleveReq)
		if err != nil {
			return err
		}

		for _, h := range res.Hits {
			_, err := b.updateEntity(h.ID, func(ent *Entity) {
				ent.Deleted = deleted
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b Bleve) buildQuery(q string) string {
	query := q
	fields := []string{"RootID", "Path", "ID", "Name", "Size", "Mtime", "MimeType", "Type", "Content", "Title"}
	for _, field := range fields {
		query = strings.ReplaceAll(query, strings.ToLower(field)+":", field+":")
	}

	if strings.Contains(query, ":") {
		return query // Sophisticated field based search
	}

	// this is a basic filename search
	return "Name:*" + strings.ReplaceAll(strings.ToLower(query), " ", `\ `) + "*"
}
