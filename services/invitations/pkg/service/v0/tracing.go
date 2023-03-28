package service

import (
	"context"

	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
	invitationstracing "github.com/owncloud/ocis/v2/services/invitations/pkg/tracing"
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

// Invite implements the Service interface.
func (t tracing) Invite(ctx context.Context, invitation *invitations.Invitation) (*invitations.Invitation, error) {
	ctx, span := invitationstracing.TraceProvider.Tracer("invitations").Start(ctx, "Invite", trace.WithAttributes(
		attribute.KeyValue{Key: "invitation", Value: attribute.StringValue(invitation.InvitedUserEmailAddress)},
	))
	defer span.End()

	return t.next.Invite(ctx, invitation)
}
