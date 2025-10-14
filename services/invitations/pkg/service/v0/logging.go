package service

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
)

// NewLogging returns a service that logs messages.
func NewLogging(next Service, logger log.Logger) Service {
	return logging{
		next:   next,
		logger: logger,
	}
}

type logging struct {
	next   Service
	logger log.Logger
}

// Invite implements the Service interface.
func (l logging) Invite(ctx context.Context, invitation *invitations.Invitation) (*invitations.Invitation, error) {
	l.logger.Debug().
		Interface("invitation", invitation).
		Msg("Invite")

	return l.next.Invite(ctx, invitation)
}
