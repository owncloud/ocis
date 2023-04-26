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

	granteeList, err := s.getGranteeList(ownerCtx, owner.GetId(), e.GranteeUserID, e.GranteeGroupID)
	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("Could not get grantee list")
		return
	}

	sharerDisplayName := owner.GetDisplayName()
	recipientList, err := s.render(ownerCtx, email.ShareCreated,
		"ShareGrantee",
		map[string]interface{}{
			"ShareSharer": sharerDisplayName,
			"ShareFolder": resourceInfo.Name,
			"ShareLink":   shareLink,
		}, granteeList, sharerDisplayName)
	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareCreated").Msg("could not get render the email")
		return
	}
	s.send(ownerCtx, recipientList)
}

func (s eventsNotifier) handleShareExpired(e events.ShareExpired) {
	logger := s.logger.With().
		Str("event", "ShareExpired").
		Str("itemid", e.ItemID.GetOpaqueId()).
		Logger()

	ownerCtx, owner, err := utils.Impersonate(e.ShareOwner, s.gwClient, s.machineAuthAPIKey)
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

	granteeList, err := s.getGranteeList(ownerCtx, owner.GetId(), e.GranteeUserID, e.GranteeGroupID)
	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareExpired").Msg("Could not get grantee name")
		return
	}

	recipientList, err := s.render(ownerCtx, email.ShareExpired,
		"ShareGrantee",
		map[string]interface{}{
			"ShareFolder": resourceInfo.GetName(),
			"ExpiredAt":   e.ExpiredAt.Format("2006-01-02 15:04:05"),
		}, granteeList, owner.GetDisplayName())
	if err != nil {
		s.logger.Error().Err(err).Str("event", "ShareExpired").Msg("could not get render the email")
		return
	}
	s.send(ownerCtx, recipientList)
}
