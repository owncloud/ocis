package service

import (
	"context"
	"net/url"

	webfingertracing "github.com/owncloud/ocis/v2/services/webfinger/pkg/tracing"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/webfinger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next Service) Service {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next Service
}

// Webfinger implements the Service interface.
func (t tracing) Webfinger(ctx context.Context, queryTarget *url.URL, rels []string) (webfinger.JSONResourceDescriptor, error) {
	ctx, span := webfingertracing.TraceProvider.Tracer("webfinger").Start(ctx, "Webfinger", trace.WithAttributes(
		attribute.KeyValue{Key: "query_target", Value: attribute.StringValue(queryTarget.String())},
		attribute.KeyValue{Key: "rels", Value: attribute.StringSliceValue(rels)},
	))
	defer span.End()

	return t.next.Webfinger(ctx, queryTarget, rels)
}
