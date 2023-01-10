package service

import (
	"context"

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
func (t tracing) Webfinger(ctx context.Context, resource, rel string) (webfinger.JSONResourceDescriptor, error) {
	ctx, span := webfingertracing.TraceProvider.Tracer("webfinger").Start(ctx, "Webfinger", trace.WithAttributes(
		attribute.KeyValue{Key: "resource", Value: attribute.StringValue(resource)},
		attribute.KeyValue{Key: "rel", Value: attribute.StringValue(rel)},
	))
	defer span.End()

	return t.next.Webfinger(ctx, resource, rel)
}
