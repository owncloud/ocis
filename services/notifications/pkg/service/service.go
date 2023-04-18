package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path"
	"syscall"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/email"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// Service should be named `Runner`
type Service interface {
	Run() error
}

// NewEventsNotifier provides a new eventsNotifier
func NewEventsNotifier(
	events <-chan events.Event,
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
	events            <-chan events.Event
	signals           chan os.Signal
	gwClient          gateway.GatewayAPIClient
	machineAuthAPIKey string
	emailTemplatePath string
	translationPath   string
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
				switch e := evt.Event.(type) {
				case events.SpaceShared:
					s.handleSpaceShared(e)
				case events.SpaceUnshared:
					s.handleSpaceUnshared(e)
				case events.SpaceMembershipExpired:
					s.handleSpaceMembershipExpired(e)
				case events.ShareCreated:
					s.handleShareCreated(e)
				case events.ShareExpired:
					s.handleShareExpired(e)
				}
			}()
		case <-s.signals:
			s.logger.Debug().
				Msg("eventsNotifier stopped")
			return nil
		}
	}
}

func (s eventsNotifier) render(template email.MessageTemplate, values map[string]interface{}) (string, string, error) {
	// The locate have to come from the user setting
	return email.RenderEmailTemplate(template, "en", s.emailTemplatePath, s.translationPath, values)
}

func (s eventsNotifier) send(ctx context.Context, u *user.UserId, g *group.GroupId, msg, subj, sender string) error {
	if u != nil {
		return s.channel.SendMessage(ctx, []string{u.GetOpaqueId()}, msg, subj, sender)

	}

	if g != nil {
		return s.channel.SendMessageToGroup(ctx, g, msg, subj, sender)
	}

	return nil
}

func (s eventsNotifier) getGranteeName(ctx context.Context, u *user.UserId, g *group.GroupId) (string, error) {
	switch {
	case u != nil:
		r, err := s.gwClient.GetUser(ctx, &user.GetUserRequest{UserId: u})
		if err != nil {
			return "", err
		}

		if r.Status.Code != rpc.Code_CODE_OK {
			return "", fmt.Errorf("unexpected status code from gateway client: %d", r.GetStatus().GetCode())
		}

		return r.GetUser().GetDisplayName(), nil
	case g != nil:
		r, err := s.gwClient.GetGroup(ctx, &group.GetGroupRequest{GroupId: g})
		if err != nil {
			return "", err
		}

		if r.GetStatus().GetCode() != rpc.Code_CODE_OK {
			return "", fmt.Errorf("unexpected status code from gateway client: %d", r.GetStatus().GetCode())
		}

		return r.GetGroup().GetDisplayName(), nil
	default:
		return "", errors.New("Need at least one non-nil grantee")
	}

}

func (s eventsNotifier) getResourceInfo(ctx context.Context, resourceID *provider.ResourceId, fieldmask *fieldmaskpb.FieldMask) (*provider.ResourceInfo, error) {
	// TODO: maybe cache this stat to reduce storage iops
	md, err := s.gwClient.Stat(ctx, &provider.StatRequest{
		Ref: &provider.Reference{
			ResourceId: resourceID,
		},
		FieldMask: fieldmask,
	})

	if err != nil {
		return nil, err
	}

	if md.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("could not resource info: %s", md.Status.Message)
	}
	return md.GetInfo(), nil
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
