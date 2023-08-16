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

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		logger.Error().Err(err).Msg("could not select next gateway client")
		return
	}

	ctx, err := utils.GetServiceUserContext(s.serviceAccountID, gatewayClient, s.serviceAccountSecret)
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

	resourceInfo, err := s.getResourceInfo(ctx, &resourceID, nil)
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

	executant, err := utils.GetUser(e.Executant, gatewayClient)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not get user")
		return
	}

	// Note: We're using the 'executantCtx' (authenticated as the share executant) here for requesting
	// the Grantees of the shares. Ideally the notfication service would use some kind of service
	// user for this.
	granteeList := s.ensureGranteeList(ctx, executant.GetId(), e.GranteeUserID, e.GranteeGroupID)
	if granteeList == nil {
		return
	}

	sharerDisplayName := executant.GetDisplayName()
	recipientList, err := s.render(ctx, email.SharedSpace,
		"SpaceGrantee",
		map[string]string{
			"SpaceSharer": sharerDisplayName,
			"SpaceName":   resourceInfo.GetSpace().GetName(),
			"ShareLink":   shareLink,
		}, granteeList, sharerDisplayName)
	if err != nil {
		s.logger.Error().Err(err).Str("event", "SharedSpace").Msg("could not get render the email")
		return
	}
	s.send(ctx, recipientList)
}

func (s eventsNotifier) handleSpaceUnshared(e events.SpaceUnshared) {
	logger := s.logger.With().
		Str("event", "SpaceUnshared").
		Str("itemid", e.ID.OpaqueId).
		Logger()

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		logger.Error().Err(err).Msg("could not select next gateway client")
		return
	}

	ctx, err := utils.GetServiceUserContext(s.serviceAccountID, gatewayClient, s.serviceAccountSecret)
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

	resourceInfo, err := s.getResourceInfo(ctx, &resourceID, nil)
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

	executant, err := utils.GetUser(e.Executant, gatewayClient)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not get user")
		return
	}

	// Note: We're using the 'executantCtx' (authenticated as the share executant) here for requesting
	// the Grantees of the shares. Ideally the notfication service would use some kind of service
	// user for this.
	granteeList := s.ensureGranteeList(ctx, executant.GetId(), e.GranteeUserID, e.GranteeGroupID)
	if granteeList == nil {
		return
	}

	sharerDisplayName := executant.GetDisplayName()
	recipientList, err := s.render(ctx, email.UnsharedSpace,
		"SpaceGrantee",
		map[string]string{
			"SpaceSharer": sharerDisplayName,
			"SpaceName":   resourceInfo.GetSpace().Name,
			"ShareLink":   shareLink,
		}, granteeList, sharerDisplayName)
	if err != nil {
		s.logger.Error().Err(err).Str("event", "UnsharedSpace").Msg("Could not get render the email")
		return
	}
	s.send(ctx, recipientList)
}

func (s eventsNotifier) handleSpaceMembershipExpired(e events.SpaceMembershipExpired) {
	logger := s.logger.With().
		Str("event", "SpaceMembershipExpired").
		Str("itemid", e.SpaceID.GetOpaqueId()).
		Logger()

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		logger.Error().Err(err).Msg("could not select next gateway client")
		return
	}

	ctx, err := utils.GetServiceUserContext(s.serviceAccountID, gatewayClient, s.serviceAccountSecret)
	if err != nil {
		logger.Error().Err(err).Msg("Could not impersonate sharer")
		return
	}

	owner, err := utils.GetUser(e.SpaceOwner, gatewayClient)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not get user")
		return
	}

	granteeList := s.ensureGranteeList(ctx, owner.GetId(), e.GranteeUserID, e.GranteeGroupID)
	if granteeList == nil {
		return
	}

	recipientList, err := s.render(ctx, email.MembershipExpired,
		"SpaceGrantee",
		map[string]string{
			"SpaceName": e.SpaceName,
			"ExpiredAt": e.ExpiredAt.Format("2006-01-02 15:04:05"),
		}, granteeList, owner.GetDisplayName())
	if err != nil {
		s.logger.Error().Err(err).Str("event", "SpaceUnshared").Msg("could not get render the email")
		return
	}
	s.send(ctx, recipientList)
}
