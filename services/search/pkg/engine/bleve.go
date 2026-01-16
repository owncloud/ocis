package engine

import (
	"context"
	"errors"
	"math"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/token/porter"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	storageProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"

	libregraph "github.com/owncloud/libre-graph-api-go"

	searchMessage "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchService "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	bleveEngine "github.com/owncloud/ocis/v2/services/search/pkg/engine/bleve"
	searchQuery "github.com/owncloud/ocis/v2/services/search/pkg/query"
)

// Bleve represents a search engine which utilizes bleve to search and store resources.
type Bleve struct {
	indexGetter  bleveEngine.IndexGetter
	queryCreator searchQuery.Creator[query.Query]
}

// NewBleveEngine creates a new Bleve instance
// If scalable is set to true, one connection to the index is created and
// closed per operation, so multiple operations can be executed in parallel.
// If set to false, only one write connection is created for the whole
// service, which will lock the index for other processes. In this case,
// you must close the engine yourself.
func NewBleveEngine(indexGetter bleveEngine.IndexGetter, queryCreator searchQuery.Creator[query.Query]) *Bleve {
	return &Bleve{
		indexGetter:  indexGetter,
		queryCreator: queryCreator,
	}
}

// Close will get the index and close it. If the indexGetter is returning
// new instances, this method will close just the new returned instance but
// not any other instances that might be in use.
//
// This method is useful if "memory" and "persistent" (not "persistentScale")
// index getters are used.
func (b *Bleve) Close() error {
	// regardless of the implementation, we want to close the index
	bleveIndex, _, err := b.indexGetter.GetIndex()
	if err != nil {
		return err
	}
	return bleveIndex.Close()
}

// BuildBleveMapping builds a bleve index mapping which can be used for indexing
func BuildBleveMapping() (mapping.IndexMapping, error) {
	nameMapping := bleve.NewTextFieldMapping()
	nameMapping.Analyzer = "lowercaseKeyword"

	lowercaseMapping := bleve.NewTextFieldMapping()
	lowercaseMapping.IncludeInAll = false
	lowercaseMapping.Analyzer = "lowercaseKeyword"

	fulltextFieldMapping := bleve.NewTextFieldMapping()
	fulltextFieldMapping.Analyzer = "fulltext"
	fulltextFieldMapping.IncludeInAll = false

	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("Name", nameMapping)
	docMapping.AddFieldMappingsAt("Tags", lowercaseMapping)
	docMapping.AddFieldMappingsAt("Content", fulltextFieldMapping)

	// Add explicit Photo field mappings to ensure fields are stored
	photoMapping := bleve.NewDocumentMapping()
	photoStringMapping := bleve.NewTextFieldMapping()
	photoStringMapping.Store = true
	photoStringMapping.Index = true
	photoStringMapping.Analyzer = keyword.Name
	photoNumericMapping := bleve.NewNumericFieldMapping()
	photoNumericMapping.Store = true
	photoNumericMapping.Index = true
	photoDateMapping := bleve.NewDateTimeFieldMapping()
	photoDateMapping.Store = true
	photoDateMapping.Index = true

	photoMapping.AddFieldMappingsAt("cameraMake", photoStringMapping)
	photoMapping.AddFieldMappingsAt("cameraModel", photoStringMapping)
	photoMapping.AddFieldMappingsAt("exposureDenominator", photoNumericMapping)
	photoMapping.AddFieldMappingsAt("exposureNumerator", photoNumericMapping)
	photoMapping.AddFieldMappingsAt("fNumber", photoNumericMapping)
	photoMapping.AddFieldMappingsAt("focalLength", photoNumericMapping)
	photoMapping.AddFieldMappingsAt("iso", photoNumericMapping)
	photoMapping.AddFieldMappingsAt("orientation", photoNumericMapping)
	photoMapping.AddFieldMappingsAt("takenDateTime", photoDateMapping)
	docMapping.AddSubDocumentMapping("photo", photoMapping)

	// Add explicit Location field mappings to ensure fields are stored
	locationMapping := bleve.NewDocumentMapping()
	locationNumericMapping := bleve.NewNumericFieldMapping()
	locationNumericMapping.Store = true
	locationNumericMapping.Index = true
	locationMapping.AddFieldMappingsAt("latitude", locationNumericMapping)
	locationMapping.AddFieldMappingsAt("longitude", locationNumericMapping)
	locationMapping.AddFieldMappingsAt("altitude", locationNumericMapping)
	docMapping.AddSubDocumentMapping("location", locationMapping)

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
		},
	)
	if err != nil {
		return nil, err
	}

	err = indexMapping.AddCustomAnalyzer("fulltext",
		map[string]interface{}{
			"type":      custom.Name,
			"tokenizer": unicode.Name,
			"token_filters": []string{
				lowercase.Name,
				porter.Name,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return indexMapping, nil
}

// Search executes a search request operation within the index.
// Returns a SearchIndexResponse object or an error.
func (b *Bleve) Search(ctx context.Context, sir *searchService.SearchIndexRequest) (*searchService.SearchIndexResponse, error) {
	bleveIndex, closeFn, err := b.indexGetter.GetIndex(bleveEngine.ReadOnly(true))
	if err != nil {
		return nil, err
	}
	defer closeFn()

	createdQuery, err := b.queryCreator.Create(sir.Query)
	if err != nil {
		if searchQuery.IsValidationError(err) {
			return nil, errtypes.BadRequest(err.Error())
		}
		return nil, err
	}

	q := bleve.NewConjunctionQuery(
		// Skip documents that have been marked as deleted
		&query.BoolFieldQuery{
			Bool:     false,
			FieldVal: "Deleted",
		},
		createdQuery,
	)

	if sir.Ref != nil {
		rootIDValue := storagespace.FormatResourceID(
			&storageProvider.ResourceId{
				StorageId: sir.Ref.GetResourceId().GetStorageId(),
				SpaceId:   sir.Ref.GetResourceId().GetSpaceId(),
				OpaqueId:  sir.Ref.GetResourceId().GetOpaqueId(),
			},
		)
		q.Conjuncts = append(
			q.Conjuncts,
			&query.TermQuery{
				FieldVal: "RootID",
				Term:     rootIDValue,
			},
		)
	}

	bleveReq := bleve.NewSearchRequest(q)
	bleveReq.Highlight = bleve.NewHighlight()

	switch {
	case sir.PageSize == -1:
		bleveReq.Size = math.MaxInt
	case sir.PageSize == 0:
		bleveReq.Size = 200
	default:
		bleveReq.Size = int(sir.PageSize)
	}

	bleveReq.Fields = []string{"*"}
	res, err := bleveIndex.Search(bleveReq)
	if err != nil {
		return nil, err
	}

	matches := make([]*searchMessage.Match, 0, len(res.Hits))
	totalMatches := res.Total
	for _, hit := range res.Hits {
		if sir.Ref != nil {
			hitPath := strings.TrimSuffix(getFieldValue[string](hit.Fields, "Path"), "/")
			requestedPath := utils.MakeRelativePath(sir.Ref.Path)
			isRoot := hitPath == requestedPath

			if !isRoot && requestedPath != "." && !strings.HasPrefix(hitPath, requestedPath+"/") {
				totalMatches--
				continue
			}
		}

		rootID, err := storagespace.ParseID(getFieldValue[string](hit.Fields, "RootID"))
		if err != nil {
			return nil, err
		}

		rID, err := storagespace.ParseID(getFieldValue[string](hit.Fields, "ID"))
		if err != nil {
			return nil, err
		}

		pID, _ := storagespace.ParseID(getFieldValue[string](hit.Fields, "ParentID"))
		match := &searchMessage.Match{
			Score: float32(hit.Score),
			Entity: &searchMessage.Entity{
				Ref: &searchMessage.Reference{
					ResourceId: resourceIDtoSearchID(rootID),
					Path:       getFieldValue[string](hit.Fields, "Path"),
				},
				Id:         resourceIDtoSearchID(rID),
				Name:       getFieldValue[string](hit.Fields, "Name"),
				ParentId:   resourceIDtoSearchID(pID),
				Size:       uint64(getFieldValue[float64](hit.Fields, "Size")),
				Type:       uint64(getFieldValue[float64](hit.Fields, "Type")),
				MimeType:   getFieldValue[string](hit.Fields, "MimeType"),
				Deleted:    getFieldValue[bool](hit.Fields, "Deleted"),
				Tags:       getFieldSliceValue[string](hit.Fields, "Tags"),
				Highlights: getFragmentValue(hit.Fragments, "Content", 0),
				Audio:      getAudioValue[searchMessage.Audio](hit.Fields),
				Image:      getImageValue[searchMessage.Image](hit.Fields),
				Location:   getLocationValue[searchMessage.GeoCoordinates](hit.Fields),
				Photo:      getPhotoValue[searchMessage.Photo](hit.Fields),
			},
		}

		if mtime, err := time.Parse(time.RFC3339, getFieldValue[string](hit.Fields, "Mtime")); err == nil {
			match.Entity.LastModifiedTime = &timestamppb.Timestamp{Seconds: mtime.Unix(), Nanos: int32(mtime.Nanosecond())}
		}

		matches = append(matches, match)
	}

	return &searchService.SearchIndexResponse{
		Matches:      matches,
		TotalMatches: int32(totalMatches),
	}, nil
}

// Upsert indexes or stores Resource data fields.
func (b *Bleve) Upsert(id string, r Resource) error {
	bleveIndex, closeFn, err := b.indexGetter.GetIndex()
	if err != nil {
		return err
	}
	defer closeFn()

	return bleveIndex.Index(id, r)
}

// Move updates the resource location and all of its necessary fields.
func (b *Bleve) Move(id string, parentid string, target string) error {
	bleveIndex, closeFn, err := b.indexGetter.GetIndex()
	if err != nil {
		return err
	}
	defer closeFn()

	r, err := b.getResource(bleveIndex, id)
	if err != nil {
		return err
	}
	currentPath := r.Path
	nextPath := utils.MakeRelativePath(target)

	r, err = b.updateEntity(bleveIndex, id, func(r *Resource) {
		r.Path = nextPath
		r.Name = path.Base(nextPath)
		r.ParentID = parentid
	})
	if err != nil {
		return err
	}

	if r.Type == uint64(storageProvider.ResourceType_RESOURCE_TYPE_CONTAINER) {
		q := bleve.NewConjunctionQuery(
			bleve.NewQueryStringQuery("RootID:"+r.RootID),
			bleve.NewQueryStringQuery("Path:"+escapeQuery(currentPath+"/*")),
		)
		bleveReq := bleve.NewSearchRequest(q)
		bleveReq.Size = math.MaxInt
		bleveReq.Fields = []string{"*"}
		res, err := bleveIndex.Search(bleveReq)
		if err != nil {
			return err
		}

		for _, h := range res.Hits {
			_, err := b.updateEntity(bleveIndex, h.ID, func(r *Resource) {
				r.Path = strings.Replace(r.Path, currentPath, nextPath, 1)
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
	bleveIndex, closeFn, err := b.indexGetter.GetIndex()
	if err != nil {
		return err
	}
	defer closeFn()

	return b.setDeleted(bleveIndex, id, true)
}

// Restore is the counterpart to Delete.
// It restores the resource which makes it available again.
func (b *Bleve) Restore(id string) error {
	bleveIndex, closeFn, err := b.indexGetter.GetIndex()
	if err != nil {
		return err
	}
	defer closeFn()

	return b.setDeleted(bleveIndex, id, false)
}

// Purge removes a resource from the index, irreversible operation.
func (b *Bleve) Purge(id string) error {
	bleveIndex, closeFn, err := b.indexGetter.GetIndex()
	if err != nil {
		return err
	}
	defer closeFn()

	return bleveIndex.Delete(id)
}

// DocCount returns the number of resources in the index.
func (b *Bleve) DocCount() (uint64, error) {
	bleveIndex, closeFn, err := b.indexGetter.GetIndex(bleveEngine.ReadOnly(true))
	if err != nil {
		return 0, err
	}
	defer closeFn()

	return bleveIndex.DocCount()
}

func (b *Bleve) getResource(bleveIndex bleve.Index, id string) (*Resource, error) {
	req := bleve.NewSearchRequest(bleve.NewDocIDQuery([]string{id}))
	req.Fields = []string{"*"}
	res, err := bleveIndex.Search(req)
	if err != nil {
		return nil, err
	}
	if res.Hits.Len() == 0 {
		return nil, errors.New("entity not found")
	}

	fields := res.Hits[0].Fields

	return &Resource{
		ID:       getFieldValue[string](fields, "ID"),
		RootID:   getFieldValue[string](fields, "RootID"),
		Path:     getFieldValue[string](fields, "Path"),
		ParentID: getFieldValue[string](fields, "ParentID"),
		Type:     uint64(getFieldValue[float64](fields, "Type")),
		Deleted:  getFieldValue[bool](fields, "Deleted"),
		Document: content.Document{
			Name:     getFieldValue[string](fields, "Name"),
			Title:    getFieldValue[string](fields, "Title"),
			Size:     uint64(getFieldValue[float64](fields, "Size")),
			Mtime:    getFieldValue[string](fields, "Mtime"),
			MimeType: getFieldValue[string](fields, "MimeType"),
			Content:  getFieldValue[string](fields, "Content"),
			Tags:     getFieldSliceValue[string](fields, "Tags"),
			Audio:    getAudioValue[libregraph.Audio](fields),
			Image:    getImageValue[libregraph.Image](fields),
			Location: getLocationValue[libregraph.GeoCoordinates](fields),
			Photo:    getPhotoValue[libregraph.Photo](fields),
		},
	}, nil
}

func newPointerOfType[T any]() *T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	ptr := reflect.New(t).Interface()
	return ptr.(*T)
}

func getFieldName(structField reflect.StructField) string {
	tag := structField.Tag.Get("json")
	if tag == "" {
		return structField.Name
	}

	return strings.Split(tag, ",")[0]
}

func unmarshalInterfaceMap(out any, flatMap map[string]interface{}, prefix string) bool {
	nonEmpty := false
	obj := reflect.ValueOf(out).Elem()
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		structField := obj.Type().Field(i)
		mapKey := prefix + getFieldName(structField)

		if value, ok := flatMap[mapKey]; ok {
			if field.Kind() == reflect.Ptr {
				alloc := reflect.New(field.Type().Elem())
				elemType := field.Type().Elem()

				// convert time strings from index for search requests
				if elemType == reflect.TypeOf(timestamppb.Timestamp{}) {
					if strValue, ok := value.(string); ok {
						if parsedTime, err := time.Parse(time.RFC3339, strValue); err == nil {
							alloc.Elem().Set(reflect.ValueOf(*timestamppb.New(parsedTime)))
							field.Set(alloc)
							nonEmpty = true
						}
					}
					continue
				}

				// convert time strings from index for libregraph structs when updating resources
				if elemType == reflect.TypeOf(time.Time{}) {
					if strValue, ok := value.(string); ok {
						if parsedTime, err := time.Parse(time.RFC3339, strValue); err == nil {
							alloc.Elem().Set(reflect.ValueOf(parsedTime))
							field.Set(alloc)
							nonEmpty = true
						}
					}
					continue
				}

				alloc.Elem().Set(reflect.ValueOf(value).Convert(elemType))
				field.Set(alloc)
				nonEmpty = true
			}
		}
	}

	return nonEmpty
}

func getAudioValue[T any](fields map[string]interface{}) *T {
	if !strings.HasPrefix(getFieldValue[string](fields, "MimeType"), "audio/") {
		return nil
	}

	var audio = newPointerOfType[T]()
	if ok := unmarshalInterfaceMap(audio, fields, "audio."); ok {
		return audio
	}

	return nil
}

func getImageValue[T any](fields map[string]interface{}) *T {
	var image = newPointerOfType[T]()
	if ok := unmarshalInterfaceMap(image, fields, "image."); ok {
		return image
	}

	return nil
}

func getLocationValue[T any](fields map[string]interface{}) *T {
	var location = newPointerOfType[T]()
	if ok := unmarshalInterfaceMap(location, fields, "location."); ok {
		return location
	}

	return nil
}

func getPhotoValue[T any](fields map[string]interface{}) *T {
	var photo = newPointerOfType[T]()

	if ok := unmarshalInterfaceMap(photo, fields, "photo."); ok {
		return photo
	}

	return nil
}

func (b *Bleve) updateEntity(bleveIndex bleve.Index, id string, mutateFunc func(r *Resource)) (*Resource, error) {
	it, err := b.getResource(bleveIndex, id)
	if err != nil {
		return nil, err
	}

	mutateFunc(it)

	return it, bleveIndex.Index(it.ID, it)
}

func (b *Bleve) setDeleted(bleveIndex bleve.Index, id string, deleted bool) error {
	it, err := b.updateEntity(bleveIndex, id, func(r *Resource) {
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
		res, err := bleveIndex.Search(bleveReq)
		if err != nil {
			return err
		}

		for _, h := range res.Hits {
			_, err := b.updateEntity(bleveIndex, h.ID, func(r *Resource) {
				r.Deleted = deleted
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
