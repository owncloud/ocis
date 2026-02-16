package service

import (
	"context"

	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next Service, tp trace.TracerProvider) Service {
	return tracing{
		next: next,
		tp:   tp,
	}
}

type tracing struct {
	next Service
	tp   trace.TracerProvider
}

// Invite implements the Service interface.
func (t tracing) Invite(ctx context.Context, invitation *invitations.Invitation) (*invitations.Invitation, error) {
	spanOpts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.KeyValue{
				Key: "invitation", Value: attribute.StringValue(invitation.InvitedUserEmailAddress),
			}),
	}
	ctx, span := t.tp.Tracer("invitations").Start(ctx, "Invite", spanOpts...)
	defer span.End()

	return t.next.Invite(ctx, invitation)
}
