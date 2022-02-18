package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/owncloud/ocis/ocis-pkg/log"
)

type Service interface {
	Run() error
}

func NewEventsNotifier(events <-chan interface{}, logger log.Logger) Service {
	return eventsNotifier{
		logger:  logger,
		events:  events,
		signals: make(chan os.Signal, 1),
	}
}

type eventsNotifier struct {
	logger  log.Logger
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
				fmt.Println(evt)
			}()
		case <-s.signals:
			s.logger.Debug().
				Msg("eventsNotifier stopped")
			return nil
		}
	}
}
