package service

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path"
	"syscall"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	groupv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
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
func NewEventsNotifier(
	events <-chan interface{},
	channel channels.Channel,
	logger log.Logger,
	gwClient gateway.GatewayAPIClient,
	machineAuthAPIKey, emailTemplatePath, ocisURL string) Service {
	return eventsNotifier{
		logger:            logger,
		channel:           channel,
		events:            events,
		signals:           make(chan os.Signal, 1),
		gwClient:          gwClient,
		machineAuthAPIKey: machineAuthAPIKey,
		emailTemplatePath: emailTemplatePath,
		ocisURL:           ocisURL,
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
	ocisURL           string
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

func (s eventsNotifier) handleShareCreated(e events.ShareCreated) {
	logger := s.logger.With().
		Str("event", "ShareCreated").
		Str("itemid", e.ItemID.OpaqueId).
		Logger()

	impersonateRes, err := s.impersonate(e.Sharer)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not impersonate sharer")
		return
	}
	ownerCtx := metadata.AppendToOutgoingContext(context.Background(), revactx.TokenHeader, impersonateRes.Token)

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

	shareGrantee := ""
	switch {
	// Note: We're using the 'ownerCtx' (authenticated as the share owner) here for requesting
	// the Grantees of the shares. Ideally the notfication service would use some kind of service
	// user for this.
	case e.GranteeUserID != nil:
		granteeUserResponse, err := s.gwClient.GetUser(ownerCtx, &userv1beta1.GetUserRequest{
			UserId: e.GranteeUserID,
		})
		if err != nil || granteeUserResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
			s.logger.Error().
				Err(err).
				Str("event", "ShareCreated").
				Msg("Could not get user response from gatway client")
			return
		}
		shareGrantee = granteeUserResponse.GetUser().DisplayName
	case e.GranteeGroupID != nil:
		granteeGroupResponse, err := s.gwClient.GetGroup(ownerCtx, &groupv1beta1.GetGroupRequest{
			GroupId: e.GranteeGroupID,
		})
		if err != nil || granteeGroupResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
			s.logger.Error().
				Err(err).
				Str("event", "ShareCreated").
				Msg("Could not get group response from gatway client")
			return
		}
		shareGrantee = granteeGroupResponse.GetGroup().DisplayName
	default:
		s.logger.Error().
			Str("event", "ShareCreated").
			Msg("Event 'ShareCreated' has no grantee")
		return
	}

	sharerDisplayName := impersonateRes.GetUser().DisplayName
	msg, err := email.RenderEmailTemplate("shares/shareCreated.email.body.tmpl", map[string]string{
		"ShareGrantee": shareGrantee,
		"ShareSharer":  sharerDisplayName,
		"ShareFolder":  resourceInfo.Name,
		"ShareLink":    shareLink,
	}, s.emailTemplatePath)

	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Msg("Could not render E-Mail body template for shares")
	}

	emailSubject, err := email.RenderEmailTemplate("shares/shareCreated.email.subject.tmpl", map[string]string{
		"ShareSharer": sharerDisplayName,
		"ShareFolder": resourceInfo.Name,
	}, s.emailTemplatePath)

	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event", "SpaceCreated").
			Msg("Could not render E-Mail subject template for shares")
	}

	if e.GranteeUserID != nil {
		err = s.channel.SendMessage(ownerCtx, []string{e.GranteeUserID.OpaqueId}, msg, emailSubject, sharerDisplayName)
	} else if e.GranteeGroupID != nil {
		err = s.channel.SendMessageToGroup(ownerCtx, e.GranteeGroupID, msg, emailSubject, sharerDisplayName)
	}
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event", "ShareCreated").
			Msg("failed to send a message")
	}
}

func (s eventsNotifier) impersonate(userID *userv1beta1.UserId) (*gateway.AuthenticateResponse, error) {
	getUserResponse, err := s.gwClient.GetUser(context.Background(), &userv1beta1.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	if getUserResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
		return nil, fmt.Errorf("error getting user: %s", getUserResponse.Status.Message)
	}

	// Get auth context
	ownerCtx := revactx.ContextSetUser(context.Background(), getUserResponse.User)
	authRes, err := s.gwClient.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + userID.OpaqueId,
		ClientSecret: s.machineAuthAPIKey,
	})
	if err != nil {
		return nil, err
	}
	if authRes.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		return nil, fmt.Errorf("error impersonating user: %s", authRes.Status.Message)
	}
	return authRes, nil
}

// TODO: this function is a backport for go1.19 url.JoinPath, upon go bump, replace this
func urlJoinPath(base string, elements ...string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(append([]string{u.Path}, elements...)...)
	return u.String(), nil
}

func (s eventsNotifier) getResourceInfo(ctx context.Context, resourceID *providerv1beta1.ResourceId, fieldmask *fieldmaskpb.FieldMask) (*providerv1beta1.ResourceInfo, error) {
	// TODO: maybe cache this stat to reduce storage iops
	md, err := s.gwClient.Stat(ctx, &providerv1beta1.StatRequest{
		Ref: &providerv1beta1.Reference{
			ResourceId: resourceID,
		},
		FieldMask: fieldmask,
	})

	if err != nil {
		return nil, err
	}

	if md.Status.Code != rpcv1beta1.Code_CODE_OK {
		return nil, fmt.Errorf("could not resource info: %s", md.Status.Message)
	}
	return md.GetInfo(), nil
}
