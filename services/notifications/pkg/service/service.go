package service

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/email"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Service interface {
	Run() error
}

// NewEventsNotifier provides a new eventsNotifier
func NewEventsNotifier(events <-chan interface{}, channel channels.Channel, logger log.Logger, gwClient gateway.GatewayAPIClient, machineAuthAPIKey, emailTemplatePath string) Service {
	return eventsNotifier{
		logger:            logger,
		channel:           channel,
		events:            events,
		signals:           make(chan os.Signal, 1),
		gwClient:          gwClient,
		machineAuthAPIKey: machineAuthAPIKey,
		emailTemplatePath: emailTemplatePath,
	}
}

type eventsNotifier struct {
	logger            log.Logger
	channel           channels.Channel
	events            <-chan interface{}
	signals           chan os.Signal
	gwClient          gateway.GatewayAPIClient
	machineAuthAPIKey string
	emailTemplatePath string
}

func (s eventsNotifier) Run() error {
	signal.Notify(s.signals, syscall.SIGINT, syscall.SIGTERM)
	s.logger.Debug().
		Msg("eventsNotifier started")
	for {
		select {
		case evt := <-s.events:
			go func() {
				switch e := evt.(type) {
				case events.SpaceShared:
					s.handleSpaceShared(e)
				case events.ShareCreated:
					s.handleShareCreated(e)
				}
			}()
		case <-s.signals:
			s.logger.Debug().
				Msg("eventsNotifier stopped")
			return nil
		}
	}
}

func (s eventsNotifier) handleSpaceShared(e events.SpaceShared) {
	userResponse, err := s.gwClient.GetUser(context.Background(), &userv1beta1.GetUserRequest{
		UserId: e.Owner,
	})
	if err != nil || userResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "SpaceCreated").
			Msg("Could not get user response from gatway client")
		return
	}
	// Get auth context
	ownerCtx := ctxpkg.ContextSetUser(context.Background(), userResponse.User)
	authRes, err := s.gwClient.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + e.Executant.OpaqueId,
		ClientSecret: s.machineAuthAPIKey,
	})
	if err != nil || authRes.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "SpaceCreated").
			Msg("Could not impersonate sharer")
		return
	}

	if authRes.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "SpaceCreated").
			Msg("could not get authenticated context for user")
		return
	}
	ownerCtx = metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token)

	resourceID, err := storagespace.ParseID(e.ID.OpaqueId)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event", "SpaceCreated").
			Str("itemid", e.ID.OpaqueId).
			Msg("could not parse resourceid from ItemID ")
		return
	}
	// TODO: maybe cache this stat to reduce storage iops
	md, err := s.gwClient.Stat(ownerCtx, &providerv1beta1.StatRequest{
		Ref: &providerv1beta1.Reference{
			ResourceId: &resourceID,
		},
		FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"name"}},
	})

	if err != nil || md.Status.Code != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Str("itemid", e.ID.OpaqueId).
			Msg("could not stat resource")
		return
	}

	if md.Status.Code != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Str("itemid", e.ID.OpaqueId).
			Str("rpc status", md.Status.Code.String()).
			Msg("could not stat resource")
		return
	}

	// old code
	msg, err := email.RenderEmailTemplate("sharedSpace.email.tmpl", map[string]string{
		// TODO: add additional fields here (like link etc.)
		"SpaceSharer": "spacesharer",
		"SpaceName":   md.Info.Space.Name,
	}, s.emailTemplatePath)

	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event", "SpaceCreated").
			Msg("Could not render E-Mail template for spaces")
	}
	if e.Executant != nil {
		err = s.channel.SendMessage([]string{e.Executant.OpaqueId}, msg, "You were invited to join a space")
	}
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event", "SpaceCreated").
			Msg("failed to send a message")
	}
}

func (s eventsNotifier) handleShareCreated(e events.ShareCreated) {
	userResponse, err := s.gwClient.GetUser(context.Background(), &userv1beta1.GetUserRequest{
		UserId: e.Sharer,
	})
	if err != nil || userResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Msg("Could not get user response from gatway client")
		return
	}

	// Get auth context
	ownerCtx := ctxpkg.ContextSetUser(context.Background(), userResponse.User)
	authRes, err := s.gwClient.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + e.Sharer.OpaqueId,
		ClientSecret: s.machineAuthAPIKey,
	})
	if err != nil || authRes.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Msg("Could not impersonate sharer")
		return
	}

	if authRes.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Msg("could not get authenticated context for user")
		return
	}
	ownerCtx = metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token)

	// TODO: maybe cache this stat to reduce storage iops
	md, err := s.gwClient.Stat(ownerCtx, &providerv1beta1.StatRequest{
		Ref: &providerv1beta1.Reference{
			ResourceId: e.ItemID,
		},
		FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"name"}},
	})

	if err != nil || md.Status.Code != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Str("itemid", e.ItemID.OpaqueId).
			Msg("could not stat resource")
		return
	}

	if md.Status.Code != rpcv1beta1.Code_CODE_OK {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Str("itemid", e.ItemID.OpaqueId).
			Str("rpc status", md.Status.Code.String()).
			Msg("could not stat resource")
		return
	}

	shareLink, err := url.JoinPath(e.Executant.Idp, "files/shares/with-me")

	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Msg("could not create link to the share")
		return
	}

	msg, err := email.RenderEmailTemplate("shareCreated.email.tmpl", map[string]string{
		// TODO: add additional fields here (like link etc.)
		"ShareSharer": userResponse.User.DisplayName,
		"ShareFolder": md.Info.Name,
		"ShareLink":   shareLink,
	}, s.emailTemplatePath)

	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Msg("Could not render E-Mail template for shares")
	}
	if e.GranteeUserID != nil {
		err = s.channel.SendMessage([]string{e.GranteeUserID.OpaqueId}, msg, "You have received a share.")
	} else if e.GranteeGroupID != nil {
		err = s.channel.SendMessageToGroup(e.GranteeGroupID, msg, "You have received a share.")
	}
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Msg("failed to send a message")
	}
}
