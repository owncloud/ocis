// Package channels provides different communication channels to notify users.
package channels

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/config"
	"github.com/pkg/errors"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Channel defines the methods of a communication channel.
type Channel interface {
	// SendMessage sends a message to users.
	SendMessage(ctx context.Context, message *Message) error
}

// Message represent the already rendered message including the user id opaqueID
type Message struct {
	Sender       string
	Recipient    []string
	Subject      string
	TextBody     string
	HTMLBody     string
	AttachInline map[string][]byte
}

// NewMailChannel instantiates a new mail communication channel.
func NewMailChannel(cfg config.Config, logger log.Logger) (Channel, error) {
	return Mail{
		conf:   cfg,
		logger: logger,
	}, nil
}

// Mail is the communication channel for email.
type Mail struct {
	conf   config.Config
	logger log.Logger
}

func (m Mail) getMailClient() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()
	server.Host = m.conf.Notifications.SMTP.Host
	server.Port = m.conf.Notifications.SMTP.Port
	server.Username = m.conf.Notifications.SMTP.Username
	if server.Username == "" {
		// compatibility fallback
		server.Username = m.conf.Notifications.SMTP.Sender
	}
	server.Password = m.conf.Notifications.SMTP.Password
	if server.TLSConfig == nil {
		server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	server.TLSConfig.InsecureSkipVerify = m.conf.Notifications.SMTP.Insecure

	switch strings.ToLower(m.conf.Notifications.SMTP.Authentication) {
	case "login":
		server.Authentication = mail.AuthLogin
	case "plain":
		server.Authentication = mail.AuthPlain
	case "crammd5":
		server.Authentication = mail.AuthCRAMMD5
	case "none":
		server.Authentication = mail.AuthNone
	case "auto", "":
		server.Authentication = mail.AuthAuto
	default:
		return nil, errors.New("unknown mail authentication method")
	}

	switch strings.ToLower(m.conf.Notifications.SMTP.Encryption) {
	case "tls":
		server.Encryption = mail.EncryptionTLS
		server.TLSConfig.ServerName = m.conf.Notifications.SMTP.Host
	case "starttls":
		server.Encryption = mail.EncryptionSTARTTLS
		server.TLSConfig.ServerName = m.conf.Notifications.SMTP.Host
	case "ssl":
		server.Encryption = mail.EncryptionSSL
	case "ssltls":
		server.Encryption = mail.EncryptionSSLTLS
	case "none":
		server.Encryption = mail.EncryptionNone
	default:
		return nil, errors.New("unknown mail encryption method")
	}

	smtpClient, err := server.Connect()
	if err != nil {
		return nil, err
	}

	return smtpClient, nil
}

// SendMessage sends a message to all given users.
func (m Mail) SendMessage(ctx context.Context, message *Message) error {
	if m.conf.Notifications.SMTP.Host == "" {
		m.logger.Info().Str("mail", "SendMessage").Msg("failed to send a message. SMTP host is  not set")
		return nil
	}

	smtpClient, err := m.getMailClient()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	if message.Sender != "" {
		tmpSender := fmt.Sprintf("%s via %s", message.Sender, m.conf.Notifications.SMTP.Sender)
		tmpParts := strings.SplitN(tmpSender, "<", 2)
		tmpSender = fmt.Sprintf("\"%s\" <%s", tmpParts[0], tmpParts[1])
		email.SetFrom(tmpSender).AddTo(message.Recipient...)
	} else {
		email.SetFrom(m.conf.Notifications.SMTP.Sender).AddTo(message.Recipient...)
	}
	email.SetSubject(message.Subject)
	email.SetBody(mail.TextPlain, message.TextBody)
	if message.HTMLBody != "" {
		email.AddAlternative(mail.TextHTML, message.HTMLBody)
		for filename, data := range message.AttachInline {
			email.Attach(&mail.File{Data: data, Name: filename, Inline: true})
		}
	}

	return email.Send(smtpClient)
}
