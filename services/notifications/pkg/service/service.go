package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
)

type Service interface {
	Run() error
}

func NewEventsNotifier(events <-chan interface{}, channel channels.Channel, logger log.Logger) Service {
	return eventsNotifier{
		logger:  logger,
		channel: channel,
		events:  events,
		signals: make(chan os.Signal, 1),
	}
}

type eventsNotifier struct {
	logger  log.Logger
	channel channels.Channel
	events  <-chan interface{}
	signals chan os.Signal
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
				case events.ShareCreated:
					msg := "You got a share!"
					var err error
					if e.GranteeUserID != nil {
						err = s.channel.SendMessage([]string{e.GranteeUserID.OpaqueId}, msg)
					} else if e.GranteeGroupID != nil {
						err = s.channel.SendMessageToGroup(e.GranteeGroupID, msg)
					}
					if err != nil {
						s.logger.Error().
							Err(err).
							Str("event", "ShareCreated").
							Msg("failed to send a message")
					}
				case events.VirusscanFinished:
					if !e.Infected {
						// no need to annoy the user
						return
					}

					if e.ExecutingUser == nil {
						s.logger.Error().Str("events", "VirusscanFinished").Str("uploadid", e.UploadID).Msg("no executing user")
						return
					}
					m := "Dear %s,\nThe virusscan of file '%s' discovered it is infected with '%s'.\nThe system is configured to handle infected files like: %s.\nContact your administrator for more information."
					msg := fmt.Sprintf(m, e.ExecutingUser.GetUsername(), e.Filename, e.Description, e.Outcome)
					if err := s.channel.SendMessage([]string{e.ExecutingUser.GetId().GetOpaqueId()}, msg); err != nil {
						s.logger.Error().Err(err).Str("event", "VirusScanFinished").Msg("failed to send a message")
					}
				}
			}()
		case <-s.signals:
			s.logger.Debug().
				Msg("eventsNotifier stopped")
			return nil
		}
	}
}
