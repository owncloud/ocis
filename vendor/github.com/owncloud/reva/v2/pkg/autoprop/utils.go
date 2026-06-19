package autoprop

import (
	"context"
	"net/http"
	"strings"

	micrometa "go-micro.dev/v4/metadata"
	"google.golang.org/grpc/metadata"
)

const (
	GRPCAutoPropPrefix  = "autoprop-"
	HTTPAutoPropPrefix  = "X-Ocis-Autoprop-"
	MicroAutoPropPrefix = "M-Ocis-Autoprop-"
)

// moveOcisMetaToOutgoingContext will copy the values from the oCIS meta to
// the outgoing context so the data is sent through the GRPC call.
// Keys from the oCIS metadata will have the AutoPropPrefix prepended so
// they're easier to identify.
func moveOcisMetaToOutgoingContext(ctx context.Context) context.Context {
	md := metadata.Pairs()

	meta := GetMetaFromContext(ctx)
	if meta != nil {
		for key, values := range meta.CreateCopyAsMap(GRPCAutoPropPrefix) {
			md.Set(key, values...)
		}
	}

	md2, exists := metadata.FromOutgoingContext(ctx)
	if exists {
		md = metadata.Join(md, md2)
	}

	return metadata.NewOutgoingContext(ctx, md)
}

// moveIncomingContextToOcisMeta will copy the incoming GRPC context values to
// the oCIS meta. The oCIS meta will be set in the new returned context.
// Only keys with the AutoPropPrefix will be copied.
func moveIncomingContextToOcisMeta(ctx context.Context) context.Context {
	meta, isNew := GetMetaFromContext(ctx), false
	if meta == nil {
		meta = NewMeta()
		isNew = true
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for key, values := range md {
			if unprefixedKey, hasPrefix := strings.CutPrefix(key, GRPCAutoPropPrefix); hasPrefix {
				for _, value := range values {
					meta.AppendMeta(unprefixedKey, value)
				}
			}
		}
	}

	if isNew {
		return SetMetaToContext(ctx, meta)
	}
	// No need to create new context if there is metadata in the context
	// because it's already updated. Just return the previous context
	return ctx
}

func moveOcisMetaToHttpHeaders(r *http.Request, ctx context.Context) {
	meta := GetMetaFromContext(ctx)
	if meta != nil {
		for key, values := range meta.CreateCopyAsMap(HTTPAutoPropPrefix) {
			for _, value := range values {
				// http.CanonicalHeaderKey(...) is done while adding the key-value
				r.Header.Add(key, value)
			}
		}
	}
}

func moveHttpHeadersToOcisMeta(r *http.Request, ctx context.Context) context.Context {
	meta, isNew := GetMetaFromContext(ctx), false
	if meta == nil {
		meta = NewMeta()
		isNew = true
	}

	canonicalPrefix := http.CanonicalHeaderKey(HTTPAutoPropPrefix)
	for key, values := range r.Header {
		if unprefixedKey, hasPrefix := strings.CutPrefix(key, canonicalPrefix); hasPrefix {
			for _, value := range values {
				meta.AppendMeta(unprefixedKey, value)
			}
		}
	}

	if isNew {
		return SetMetaToContext(ctx, meta)
	}
	// No need to create new context if there is metadata in the context
	// because it's already updated. Just return the same request
	return ctx
}

func moveOcisMetaToGoMicroMetadata(ctx context.Context) context.Context {
	meta := GetMetaFromContext(ctx)
	if meta != nil {
		md := make(micrometa.Metadata, meta.Len())
		for key, values := range meta.CreateCopyAsMap(MicroAutoPropPrefix) {
			md.Set(key, strings.Join(values, "|||"))
		}
		return micrometa.MergeContext(ctx, md, true)
	}
	return ctx
}

func moveGoMicroMetadataToOcisMeta(ctx context.Context) context.Context {
	meta, isNew := GetMetaFromContext(ctx), false
	if meta == nil {
		meta = NewMeta()
		isNew = true
	}

	md, ok := micrometa.FromContext(ctx)
	if ok {
		for key, values := range md {
			if unprefixedKey, hasPrefix := strings.CutPrefix(key, MicroAutoPropPrefix); hasPrefix {
				for _, value := range strings.Split(values, "|||") {
					meta.AppendMeta(unprefixedKey, value)
				}
			}
		}
	}

	if isNew {
		return SetMetaToContext(ctx, meta)
	}
	// No need to create new context if there is metadata in the context
	// because it's already updated. Just return the previous context
	return ctx
}
