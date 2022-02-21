// Package channels provides different communication channels to notify users.
package channels

import (
	"context"
	"net/smtp"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/notifications/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Channel defines the methods of a communication channel.
type Channel interface {
	// Todo(c0rby): Do we need a PrepareMessage method?
	// Maybe channels need to format the message or will the caller
	// of SendMessage do that?
	// SendMessage sends a message in a channel specific way.
	SendMessage(receiver, msg string) error
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
	}, nil
}

// Mail is the communcation channel for email.
type Mail struct {
	gatewayClient gateway.GatewayAPIClient
	conf          config.Config
}

func (m Mail) SendMessage(receiver, msg string) error {
	smtpConf := m.conf.Notifications.SMTP
	res, err := m.gatewayClient.Authenticate(context.Background(), &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + receiver,
		ClientSecret: m.conf.Notifications.MachineAuthSecret,
	})
	if err != nil {
		return err
	}

	to := []string{res.User.Mail}
	body := []byte(msg)
	auth := smtp.PlainAuth("", smtpConf.Sender, smtpConf.Password, smtpConf.Host)
	return smtp.SendMail(smtpConf.Host+":"+smtpConf.Port, auth, smtpConf.Sender, to, body)
}
