package service

import (
	"context"

	groupv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/email"
	"google.golang.org/grpc/metadata"
)

func (s eventsNotifier) handleSpaceShared(e events.SpaceShared) {
	logger := s.logger.With().
		Str("event", "SpaceShared").
		Str("itemid", e.ID.OpaqueId).
		Logger()

	impersonateRes, err := s.impersonate(e.Executant)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not handle space shared event")
		return
	}
	ownerCtx := metadata.AppendToOutgoingContext(context.Background(), revactx.TokenHeader, impersonateRes.Token)

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

	spaceGrantee := ""
	switch {
	// Note: We're using the 'ownerCtx' (authenticated as the share owner) here for requesting
	// the Grantees of the shares. Ideally the notfication service would use some kind of service
	// user for this.
	case e.GranteeUserID != nil:
		granteeUserResponse, err := s.gwClient.GetUser(ownerCtx, &userv1beta1.GetUserRequest{
			UserId: e.GranteeUserID,
		})
		if err != nil || granteeUserResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
			logger.Error().
				Err(err).
				Msg("Could not get user response from gatway client")
			return
		}
		spaceGrantee = granteeUserResponse.GetUser().GetDisplayName()
	case e.GranteeGroupID != nil:
		granteeGroupResponse, err := s.gwClient.GetGroup(ownerCtx, &groupv1beta1.GetGroupRequest{
			GroupId: e.GranteeGroupID,
		})
		if err != nil || granteeGroupResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
			logger.Error().
				Err(err).
				Msg("Could not get group response from gatway client")
			return
		}
		spaceGrantee = granteeGroupResponse.GetGroup().GetDisplayName()
	default:
		logger.Error().
			Msg("Event 'SpaceShared' has no grantee")
		return
	}

	sharerDisplayName := impersonateRes.GetUser().GetDisplayName()
	msg, err := email.RenderEmailTemplate("spaces/sharedSpace.email.body.tmpl", map[string]string{
		"SpaceGrantee": spaceGrantee,
		"SpaceSharer":  sharerDisplayName,
		"SpaceName":    resourceInfo.GetSpace().Name,
		"ShareLink":    shareLink,
	}, s.emailTemplatePath)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not render E-Mail body template for spaces")
	}

	emailSubject, err := email.RenderEmailTemplate("spaces/sharedSpace.email.subject.tmpl", map[string]string{
		"SpaceSharer": sharerDisplayName,
		"SpaceName":   resourceInfo.GetSpace().Name,
	}, s.emailTemplatePath)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not render E-Mail subject template for spaces")
	}

	if e.GranteeUserID != nil {
		err = s.channel.SendMessage(ownerCtx, []string{e.GranteeUserID.OpaqueId}, msg, emailSubject, sharerDisplayName)
	} else if e.GranteeGroupID != nil {
		err = s.channel.SendMessageToGroup(ownerCtx, e.GranteeGroupID, msg, emailSubject, sharerDisplayName)
	}
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to send a message")
	}
}

func (s eventsNotifier) handleSpaceUnshared(e events.SpaceUnshared) {
	logger := s.logger.With().
		Str("event", "SpaceUnshared").
		Str("itemid", e.ID.OpaqueId).
		Logger()

	impersonateRes, err := s.impersonate(e.Executant)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not handle space unshared event")
		return
	}
	ownerCtx := metadata.AppendToOutgoingContext(context.Background(), revactx.TokenHeader, impersonateRes.Token)

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

	spaceGrantee := ""
	switch {
	// Note: We're using the 'ownerCtx' (authenticated as the share owner) here for requesting
	// the Grantees of the shares. Ideally the notfication service would use some kind of service
	// user for this.
	case e.GranteeUserID != nil:
		granteeUserResponse, err := s.gwClient.GetUser(ownerCtx, &userv1beta1.GetUserRequest{
			UserId: e.GranteeUserID,
		})
		if err != nil || granteeUserResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
			logger.Error().
				Err(err).
				Msg("Could not get user response from gatway client")
			return
		}
		spaceGrantee = granteeUserResponse.GetUser().GetDisplayName()
	case e.GranteeGroupID != nil:
		granteeGroupResponse, err := s.gwClient.GetGroup(ownerCtx, &groupv1beta1.GetGroupRequest{
			GroupId: e.GranteeGroupID,
		})
		if err != nil || granteeGroupResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
			logger.Error().
				Err(err).
				Msg("Could not get group response from gatway client")
			return
		}
		spaceGrantee = granteeGroupResponse.GetGroup().GetDisplayName()
	default:
		logger.Error().
			Msg("Event 'SpaceShared' has no grantee")
		return
	}

	sharerDisplayName := impersonateRes.GetUser().GetDisplayName()
	msg, err := email.RenderEmailTemplate("spaces/unsharedSpace.email.body.tmpl", map[string]string{
		"SpaceGrantee": spaceGrantee,
		"SpaceSharer":  sharerDisplayName,
		"SpaceName":    resourceInfo.GetSpace().Name,
		"ShareLink":    shareLink,
	}, s.emailTemplatePath)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not render E-Mail body template for spaces")
	}

	emailSubject, err := email.RenderEmailTemplate("spaces/unsharedSpace.email.subject.tmpl", map[string]string{
		"SpaceSharer": sharerDisplayName,
		"SpaceName":   resourceInfo.GetSpace().Name,
	}, s.emailTemplatePath)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not render E-Mail subject template for spaces")
	}

	if e.GranteeUserID != nil {
		err = s.channel.SendMessage(ownerCtx, []string{e.GranteeUserID.OpaqueId}, msg, emailSubject, sharerDisplayName)
	} else if e.GranteeGroupID != nil {
		err = s.channel.SendMessageToGroup(ownerCtx, e.GranteeGroupID, msg, emailSubject, sharerDisplayName)
	}
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to send a message")
	}
}
