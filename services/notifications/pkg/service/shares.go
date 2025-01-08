package service

import (
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/email"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (s eventsNotifier) handleShareCreated(e events.ShareCreated) {
	logger := s.logger.With().
		Str("event", "ShareCreated").
		Str("itemid", e.ItemID.OpaqueId).
		Logger()

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		logger.Error().Err(err).Msg("could not select next gateway client")
		return
	}

	ctx, err := utils.GetServiceUserContext(s.serviceAccountID, gatewayClient, s.serviceAccountSecret)
	if err != nil {
		logger.Error().Err(err).Msg("Could not impersonate service user")
		return
	}

	resourceInfo, err := s.getResourceInfo(ctx, e.ItemID, &fieldmaskpb.FieldMask{Paths: []string{"name"}})
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

	owner, err := utils.GetUser(e.Sharer, gatewayClient)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get user")
		return
	}

	granteeList := s.ensureGranteeList(ctx, owner.GetId(), e.GranteeUserID, e.GranteeGroupID)
	filteredGrantees := s.filter.execute(ctx, granteeList, defaults.SettingUUIDProfileEventShareCreated)
	if filteredGrantees == nil {
		return
	}

	sharerDisplayName := owner.GetDisplayName()
	recipientList, err := s.render(ctx, email.ShareCreated,
		"ShareGrantee",
		map[string]string{
			"ShareSharer": sharerDisplayName,
			"ShareFolder": resourceInfo.Name,
			"ShareLink":   shareLink,
		}, filteredGrantees, sharerDisplayName)
	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("could not get render the email")
		return
	}
	s.send(ctx, recipientList)
}

func (s eventsNotifier) handleShareExpired(e events.ShareExpired) {
	logger := s.logger.With().
		Str("event", "ShareExpired").
		Str("itemid", e.ItemID.GetOpaqueId()).
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

	resourceInfo, err := s.getResourceInfo(ctx, e.ItemID, &fieldmaskpb.FieldMask{Paths: []string{"name"}})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not stat resource")
		return
	}

	owner, err := utils.GetUser(e.ShareOwner, gatewayClient)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get user")
		return
	}

	granteeList := s.ensureGranteeList(ctx, owner.GetId(), e.GranteeUserID, e.GranteeGroupID)
	filteredGrantees := s.filter.execute(ctx, granteeList, defaults.SettingUUIDProfileEventShareExpired)
	if filteredGrantees == nil {
		return
	}

	recipientList, err := s.render(ctx, email.ShareExpired,
		"ShareGrantee",
		map[string]string{
			"ShareFolder": resourceInfo.GetName(),
			"ExpiredAt":   e.ExpiredAt.Format("2006-01-02 15:04:05"),
		}, filteredGrantees, owner.GetDisplayName())
	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareExpired").Msg("could not get render the email")
		return
	}
	s.send(ctx, recipientList)
}
