package engine

import (
	"context"
	storageProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	searchMessage "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchSearvice "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"regexp"
)

//go:generate mockery --name=Engine

// Engine is the interface to the search engine
type Engine interface {
	Search(ctx context.Context, req *searchSearvice.SearchIndexRequest) (*searchSearvice.SearchIndexResponse, error)
	Upsert(id string, r Resource) error
	Move(id string, target string) error
	Delete(id string) error
	Restore(id string) error
	Purge(id string) error
	DocCount() (uint64, error)
}

// Resource is the entity that is stored in the index.
type Resource struct {
	content.Document

	ID      string
	RootID  string
	Path    string
	Type    uint64
	Deleted bool
}

func resourceIDtoSearchID(id storageProvider.ResourceId) *searchMessage.ResourceID {
	return &searchMessage.ResourceID{
		StorageId: id.GetStorageId(),
		SpaceId:   id.GetSpaceId(),
		OpaqueId:  id.GetOpaqueId()}
}

func escapeQuery(s string) string {
	re := regexp.MustCompile(`([` + regexp.QuoteMeta(`+=&|><!(){}[]^\"~*?:\/`) + `\-\s])`)
	return re.ReplaceAllString(s, "\\$1")
}

func getValue[T any](m map[string]interface{}, key string) (out T) {
	val, ok := m[key]
	if !ok {
		return
	}

	out, _ = val.(T)

	return
}
