// Package channels provides different communication channels to notify users.
package channels

import (
	"context"
	"net/smtp"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

type Channel interface {
	SendMessage(receiver, msg string) error
}

func NewMailChanel(logger log.Logger) (Channel, error) {
	gc, err := pool.GetGatewayServiceClient("localhost:9142")
	if err != nil {
		logger.Error().Err(err).Msg("could not get gateway client")
		return nil, err
	}
	return Mail{
		gatewayClient: gc,
	}, nil
}

type Mail struct {
	gatewayClient gateway.GatewayAPIClient
}

func (m Mail) SendMessage(receiver, msg string) error {
	res, err := m.gatewayClient.Authenticate(context.Background(), &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + receiver,
		ClientSecret: "change-me-please",
	})
	if err != nil {
		return err
	}

	from := "god"
	password := "godisdead"
	to := []string{res.User.Mail}
	host := "localhost"
	port := "1025"
	body := []byte(msg)
	auth := smtp.PlainAuth("", from, password, host)

	return smtp.SendMail(host+":"+port, auth, from, to, body)
}
