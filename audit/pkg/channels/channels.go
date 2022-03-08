// Package channels provides different communication channels to notify users.
package channels

import (
	"context"
	"net/smtp"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	groups "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/notifications/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/pkg/errors"
)

// Channel defines the methods of a communication channel.
type Channel interface {
	// SendMessage sends a message to users.
	SendMessage(userIDs []string, msg string) error
	// SendMessageToGroup sends a message to a group.
	SendMessageToGroup(groupdID *groups.GroupId, msg string) error
}

// NewMailChannel instantiates a new mail communication channel.
func NewMailChannel(cfg config.Config, logger log.Logger) (Channel, error) {
	gc, err := pool.GetGatewayServiceClient(cfg.Notifications.RevaGateway)
	if err != nil {
		logger.Error().Err(err).Msg("could not get gateway client")
		return nil, err
	}
	return Mail{
		gatewayClient: gc,
		conf:          cfg,
		logger:        logger,
	}, nil
}

// Mail is the communcation channel for email.
type Mail struct {
	gatewayClient gateway.GatewayAPIClient
	conf          config.Config
	logger        log.Logger
}

// SendMessage sends a message to all given users.
func (m Mail) SendMessage(userIDs []string, msg string) error {
	to, err := m.getReceiverAddresses(userIDs)
	if err != nil {
		return err
	}
	body := []byte(msg)

	smtpConf := m.conf.Notifications.SMTP
	auth := smtp.PlainAuth("", smtpConf.Sender, smtpConf.Password, smtpConf.Host)
	if err := smtp.SendMail(smtpConf.Host+":"+smtpConf.Port, auth, smtpConf.Sender, to, body); err != nil {
		return errors.Wrap(err, "could not send mail")
	}
	return nil
}

// SendMessageToGroup sends a message to all members of the given group.
func (m Mail) SendMessageToGroup(groupID *groups.GroupId, msg string) error {
	// TODO We need an authenticated context here...
	res, err := m.gatewayClient.GetGroup(context.Background(), &groups.GetGroupRequest{GroupId: groupID})
	if err != nil {
		return err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return errors.New("could not get group")
	}

	members := make([]string, 0, len(res.Group.Members))
	for _, id := range res.Group.Members {
		members = append(members, id.OpaqueId)
	}

	return m.SendMessage(members, msg)
}

func (m Mail) getReceiverAddresses(receivers []string) ([]string, error) {
	addresses := make([]string, 0, len(receivers))
	for _, id := range receivers {
		// Authenticate is too costly but at the moment our only option to get the user.
		// We don't have an authenticated context so calling `GetUser` doesn't work.
		res, err := m.gatewayClient.Authenticate(context.Background(), &gateway.AuthenticateRequest{
			Type:         "machine",
			ClientId:     "userid:" + id,
			ClientSecret: m.conf.Notifications.MachineAuthSecret,
		})
		if err != nil {
			return nil, err
		}
		if res.Status.Code != rpc.Code_CODE_OK {
			m.logger.Error().
				Interface("status", res.Status).
				Str("receiver_id", id).
				Msg("could not get user")
			continue
		}
		addresses = append(addresses, res.User.Mail)
	}

	return addresses, nil
}
