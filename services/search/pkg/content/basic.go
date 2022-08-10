package content

import (
	"context"
	storageProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/tags"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"time"
)

// Basic is the simplest Extractor implementation.
type Basic struct {
	logger log.Logger
}

// NewBasicExtractor creates a new Basic instance.
func NewBasicExtractor(logger log.Logger) (*Basic, error) {
	return &Basic{logger: logger}, nil
}

// Extract literally just rearranges the inputs and processes them into a Document.
func (b Basic) Extract(_ context.Context, ri *storageProvider.ResourceInfo) (Document, error) {
	doc := Document{
		Name:     ri.Path,
		Size:     ri.Size,
		MimeType: ri.MimeType,
	}

	if m := ri.ArbitraryMetadata.GetMetadata(); m != nil {
		if t, ok := m["tags"]; ok {
			doc.Tags = tags.FromList(t).AsSlice()
		}
	}

	if ri.Mtime != nil {
		doc.Mtime = time.Unix(int64(ri.Mtime.Seconds), int64(ri.Mtime.Nanos)).UTC().Format(time.RFC3339)
	}

	return doc, nil
}
