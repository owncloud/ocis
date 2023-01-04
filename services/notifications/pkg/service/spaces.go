package service

import (
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
)

func (s eventsNotifier) handleSpaceShared(e events.SpaceShared) {
	logger := s.logger.With().
		Str("event", "SpaceShared").
		Str("itemid", e.ID.OpaqueId).
		Logger()

	ownerCtx, sharerDisplayName, err := s.impersonate(e.Executant)
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

	msg, subj, err := s.render("spaces/sharedSpace.email.body.tmpl", "spaces/sharedSpace.email.subject.tmpl", map[string]string{
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

func (s eventsNotifier) handleSpaceUnshared(e events.SpaceUnshared) {
	logger := s.logger.With().
		Str("event", "SpaceUnshared").
		Str("itemid", e.ID.OpaqueId).
		Logger()

	ownerCtx, sharerDisplayName, err := s.impersonate(e.Executant)
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

	msg, subj, err := s.render("spaces/unsharedSpace.email.body.tmpl", "spaces/unsharedSpace.email.subject.tmpl", map[string]string{
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
