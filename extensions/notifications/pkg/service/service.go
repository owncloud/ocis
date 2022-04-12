package service

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/extensions/notifications/pkg/channels"
	"github.com/owncloud/ocis/ocis-pkg/log"
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
				}
			}()
		case <-s.signals:
			s.logger.Debug().
				Msg("eventsNotifier stopped")
			return nil
		}
	}
}
