package decorators

import (
	"context"
	"strings"

	thumbnailssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/thumbnails/v0"
	"github.com/owncloud/reva/v2/pkg/autoprop"
	"google.golang.org/grpc/metadata"
)

// NewTracing returns a service that instruments traces.
func NewAutoProp(next DecoratedService) DecoratedService {
	return tracing{
		Decorator: Decorator{next: next},
	}
}

type autoProp struct {
	Decorator
}

// GetThumbnail implements the ThumbnailServiceHandler interface.
func (a autoProp) GetThumbnail(ctx context.Context, req *thumbnailssvc.GetThumbnailRequest, rsp *thumbnailssvc.GetThumbnailResponse) error {
	// copied from autoprop.moveIncomingContextToOcisMeta
	// TODO: consider to refactor the thumbnail service so we can use
	// standard middlewares / interceptors.
	meta, isNew := autoprop.GetMetaFromContext(ctx), false
	if meta == nil {
		meta = autoprop.NewMeta()
		isNew = true
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for key, values := range md {
			if unprefixedKey, hasPrefix := strings.CutPrefix(key, autoprop.GRPCAutoPropPrefix); hasPrefix {
				for _, value := range values {
					meta.AppendMeta(unprefixedKey, value)
				}
			}
		}
	}

	newctx := ctx
	if isNew {
		newctx = autoprop.SetMetaToContext(ctx, meta)
	}

	return a.next.GetThumbnail(newctx, req, rsp)
}
