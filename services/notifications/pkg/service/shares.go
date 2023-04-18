package service

import (
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/email"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (s eventsNotifier) handleShareCreated(e events.ShareCreated) {
	logger := s.logger.With().
		Str("event", "ShareCreated").
		Str("itemid", e.ItemID.OpaqueId).
		Logger()

	ownerCtx, owner, err := utils.Impersonate(e.Sharer, s.gwClient, s.machineAuthAPIKey)
	if err != nil {
		logger.Error().Err(err).Msg("Could not impersonate sharer")
		return
	}

	resourceInfo, err := s.getResourceInfo(ownerCtx, e.ItemID, &fieldmaskpb.FieldMask{Paths: []string{"name"}})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not stat resource")
		return
	}

	shareLink, err := urlJoinPath(s.ocisURL, "files/shares/with-me")
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not create link to the share")
		return
	}

	shareGrantee, err := s.getGranteeName(ownerCtx, e.GranteeUserID, e.GranteeGroupID)
	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("Could not get grantee name")
		return
	}

	sharerDisplayName := owner.GetDisplayName()
	subj, msg, err := s.render(email.ShareCreated, map[string]interface{}{
		"ShareGrantee": shareGrantee,
		"ShareSharer":  sharerDisplayName,
		"ShareFolder":  resourceInfo.Name,
		"ShareLink":    shareLink,
	})

	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("Could not render E-Mail body template for shares")
	}

	if err := s.send(ownerCtx, e.GranteeUserID, e.GranteeGroupID, msg, subj, sharerDisplayName); err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("failed to send a message")
	}

}

func (s eventsNotifier) handleShareExpired(e events.ShareExpired) {
	logger := s.logger.With().
		Str("event", "ShareExpired").
		Str("itemid", e.ItemID.GetOpaqueId()).
		Logger()

	ctx, owner, err := utils.Impersonate(e.ShareOwner, s.gwClient, s.machineAuthAPIKey)
	if err != nil {
		logger.Error().Err(err).Msg("Could not impersonate sharer")
		return
	}

	resourceInfo, err := s.getResourceInfo(ctx, e.ItemID, &fieldmaskpb.FieldMask{Paths: []string{"name"}})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not stat resource")
		return
	}

	shareGrantee, err := s.getGranteeName(ctx, e.GranteeUserID, e.GranteeGroupID)
	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("Could not get grantee name")
		return
	}

	subj, msg, err := s.render(email.ShareExpired, map[string]interface{}{
		"ShareGrantee": shareGrantee,
		"ShareFolder":  resourceInfo.GetName(),
		"ExpiredAt":    e.ExpiredAt.Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("Could not render E-Mail body template for shares")
	}

	if err := s.send(ctx, e.GranteeUserID, e.GranteeGroupID, msg, subj, owner.GetDisplayName()); err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("failed to send a message")
	}

}
