package service

import (
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/email"
)

func (s eventsNotifier) handleSpaceShared(e events.SpaceShared) {
	logger := s.logger.With().
		Str("event", "SpaceShared").
		Str("itemid", e.ID.OpaqueId).
		Logger()

	ownerCtx, owner, err := utils.Impersonate(e.Executant, s.gwClient, s.machineAuthAPIKey)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not handle space shared event")
		return
	}

	resourceID, err := storagespace.ParseID(e.ID.OpaqueId)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not parse resourceid from ItemID ")
		return
	}

	resourceInfo, err := s.getResourceInfo(ownerCtx, &resourceID, nil)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not get space info")
		return
	}

	shareLink, err := urlJoinPath(s.ocisURL, "f", e.ID.OpaqueId)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not create link to the share")
		return
	}

	// Note: We're using the 'ownerCtx' (authenticated as the share owner) here for requesting
	// the Grantees of the shares. Ideally the notfication service would use some kind of service
	// user for this.
	spaceGrantee, err := s.getGranteeName(ownerCtx, e.GranteeUserID, e.GranteeGroupID)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get grantee name")
		return
	}

	sharerDisplayName := owner.GetDisplayName()
	subj, msg, err := s.render(email.SharedSpace, map[string]interface{}{
		"SpaceGrantee": spaceGrantee,
		"SpaceSharer":  sharerDisplayName,
		"SpaceName":    resourceInfo.GetSpace().GetName(),
		"ShareLink":    shareLink,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Could not render E-Mail template for spaces")
		return
	}

	if err := s.send(ownerCtx, e.GranteeUserID, e.GranteeGroupID, msg, subj, sharerDisplayName); err != nil {
		logger.Error().Err(err).Msg("failed to send a message")
	}
}

func (s eventsNotifier) handleSpaceUnshared(e events.SpaceUnshared) {
	logger := s.logger.With().
		Str("event", "SpaceUnshared").
		Str("itemid", e.ID.OpaqueId).
		Logger()

	ownerCtx, owner, err := utils.Impersonate(e.Executant, s.gwClient, s.machineAuthAPIKey)
	if err != nil {
		logger.Error().Err(err).Msg("could not handle space unshared event")
		return
	}

	resourceID, err := storagespace.ParseID(e.ID.OpaqueId)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not parse resourceid from ItemID ")
		return
	}

	resourceInfo, err := s.getResourceInfo(ownerCtx, &resourceID, nil)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not get space info")
		return
	}

	shareLink, err := urlJoinPath(s.ocisURL, "f", e.ID.OpaqueId)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not create link to the share")
		return
	}

	// Note: We're using the 'ownerCtx' (authenticated as the share owner) here for requesting
	// the Grantees of the shares. Ideally the notfication service would use some kind of service
	// user for this.
	spaceGrantee, err := s.getGranteeName(ownerCtx, e.GranteeUserID, e.GranteeGroupID)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get grantee name")
		return
	}

	sharerDisplayName := owner.GetDisplayName()
	subj, msg, err := s.render(email.UnsharedSpace, map[string]interface{}{
		"SpaceGrantee": spaceGrantee,
		"SpaceSharer":  sharerDisplayName,
		"SpaceName":    resourceInfo.GetSpace().Name,
		"ShareLink":    shareLink,
	})

	if err != nil {
		logger.Error().Err(err).Msg("Could not render E-Mail template for spaces")
		return
	}

	if err := s.send(ownerCtx, e.GranteeUserID, e.GranteeGroupID, msg, subj, sharerDisplayName); err != nil {
		logger.Error().Err(err).Msg("failed to send a message")
	}
}

func (s eventsNotifier) handleSpaceMembershipExpired(e events.SpaceMembershipExpired) {
	logger := s.logger.With().
		Str("event", "SpaceMembershipExpired").
		Str("itemid", e.SpaceID.GetOpaqueId()).
		Logger()

	ctx, owner, err := utils.Impersonate(e.SpaceOwner, s.gwClient, s.machineAuthAPIKey)
	if err != nil {
		logger.Error().Err(err).Msg("Could not impersonate sharer")
		return
	}

	shareGrantee, err := s.getGranteeName(ctx, e.GranteeUserID, e.GranteeGroupID)
	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("Could not get grantee name")
		return
	}

	subj, msg, err := s.render(email.MembershipExpired, map[string]interface{}{
		"SpaceGrantee": shareGrantee,
		"SpaceName":    e.SpaceName,
		"ExpiredAt":    e.ExpiredAt.Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("Could not render E-Mail body template for shares")
	}

	if err := s.send(ctx, e.GranteeUserID, e.GranteeGroupID, msg, subj, owner.GetDisplayName()); err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("failed to send a message")
	}

}
