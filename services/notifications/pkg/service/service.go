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
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/email"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	"go-micro.dev/v4/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var _defaultLocale = "en"

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
	valueService settingssvc.ValueService,
	machineAuthAPIKey, emailTemplatePath, ocisURL string) Service {

	return eventsNotifier{
		logger:            logger,
		channel:           channel,
		events:            events,
		signals:           make(chan os.Signal, 1),
		gwClient:          gwClient,
		valueService:      valueService,
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
	valueService      settingssvc.ValueService
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

// recipient represent the already rendered message including the user id opaqueID
type recipient struct {
	opaqueID string
	subject  string
	msg      string
}

func (s eventsNotifier) render(ctx context.Context, template email.MessageTemplate,
	granteeFieldName string, fields map[string]interface{}, granteeList []*user.UserId) ([]recipient, error) {
	// Render the Email Template for each user
	recipientList := make([]recipient, len(granteeList))
	for i, userID := range granteeList {
		locale := s.getUserLang(ctx, userID)
		grantee, err := s.getUserName(ctx, userID)
		if err != nil {
			return nil, err
		}
		fields[granteeFieldName] = grantee

		subj, msg, err := email.RenderEmailTemplate(template, locale, s.emailTemplatePath, s.translationPath, fields)
		if err != nil {
			return nil, err
		}
		recipientList[i] = recipient{opaqueID: userID.GetOpaqueId(), subject: subj, msg: msg}
	}
	return recipientList, nil
}

func (s eventsNotifier) send(ctx context.Context, recipientList []recipient, sender string) {
	for _, r := range recipientList {
		err := s.channel.SendMessage(ctx, []string{r.opaqueID}, r.msg, r.subject, sender)
		if err != nil {
			s.logger.Error().Err(err).Str("event", "SendEmail").Msg("failed to send a message")
		}
	}
}

func (s eventsNotifier) getGranteeList(ctx context.Context, executant, u *user.UserId, g *group.GroupId) ([]*user.UserId, error) {
	switch {
	case u != nil:
		if s.disableEmails(ctx, u) {
			return []*user.UserId{}, nil
		}
		return []*user.UserId{u}, nil
	case g != nil:
		res, err := s.gwClient.GetGroup(ctx, &group.GetGroupRequest{GroupId: g})
		if err != nil {
			return nil, err
		}
		if res.Status.Code != rpc.Code_CODE_OK {
			return nil, errors.New("could not get group")
		}

		var grantees []*user.UserId
		for _, userID := range res.GetGroup().GetMembers() {
			// don't add the executant
			if userID.GetOpaqueId() == executant.GetOpaqueId() {
				continue
			}

			// don't add users who opted out
			if s.disableEmails(ctx, userID) {
				continue
			}

			grantees = append(grantees, userID)
		}
		return grantees, nil
	default:
		return nil, errors.New("need at least one non-nil grantee")
	}
}

func (s eventsNotifier) getUserName(ctx context.Context, u *user.UserId) (string, error) {
	if u == nil {
		return "", errors.New("need at least one non-nil grantee")
	}
	r, err := s.gwClient.GetUser(ctx, &user.GetUserRequest{UserId: u})
	if err != nil {
		return "", err
	}
	if r.Status.Code != rpc.Code_CODE_OK {
		return "", fmt.Errorf("unexpected status code from gateway client: %d", r.GetStatus().GetCode())
	}
	return r.GetUser().GetDisplayName(), nil
}

func (s eventsNotifier) getUserLang(ctx context.Context, u *user.UserId) string {
	granteeCtx := metadata.Set(ctx, middleware.AccountID, u.OpaqueId)
	if resp, err := s.valueService.GetValueByUniqueIdentifiers(granteeCtx,
		&settingssvc.GetValueByUniqueIdentifiersRequest{
			AccountUuid: u.OpaqueId,
			SettingId:   defaults.SettingUUIDProfileLanguage,
		},
	); err == nil {
		val := resp.GetValue().GetValue().GetListValue().GetValues()
		if len(val) > 0 && val[0] != nil {
			return val[0].GetStringValue()
		}
	}
	return _defaultLocale
}

func (s eventsNotifier) disableEmails(ctx context.Context, u *user.UserId) bool {
	granteeCtx := metadata.Set(ctx, middleware.AccountID, u.OpaqueId)
	if resp, err := s.valueService.GetValueByUniqueIdentifiers(granteeCtx,
		&settingssvc.GetValueByUniqueIdentifiersRequest{
			AccountUuid: u.OpaqueId,
			SettingId:   defaults.SettingUUIDProfileDisableNotifications,
		},
	); err == nil {
		return resp.GetValue().GetValue().GetBoolValue()

	}
	return false
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
