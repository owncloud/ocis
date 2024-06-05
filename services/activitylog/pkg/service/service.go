package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/config"
	microstore "go-micro.dev/v4/store"
)

// Activity represents an activity
type Activity struct {
	EventID   string    `json:"event_id"`
	Depth     int       `json:"depth"`
	Timestamp time.Time `json:"timestamp"`
}

// ActivitylogService logs events per resource
type ActivitylogService struct {
	cfg    *config.Config
	log    log.Logger
	events <-chan events.Event
	store  microstore.Store
}

// New is what you need to implement.
func New(opts ...Option) (*ActivitylogService, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.Stream == nil {
		return nil, errors.New("stream is required")
	}

	if o.Store == nil {
		return nil, errors.New("store is required")
	}

	ch, err := events.Consume(o.Stream, o.Config.Service.Name, o.RegisteredEvents...)
	if err != nil {
		return nil, err
	}

	s := &ActivitylogService{
		log:    o.Logger,
		cfg:    o.Config,
		events: ch,
		store:  o.Store,
	}

	return s, nil
}

// Run runs the service
func (a *ActivitylogService) Run() error {
	for e := range a.events {
		switch ev := e.Event.(type) {
		case events.PostprocessingFinished:
			fmt.Println("PostprocessingFinished event received", ev)
		}
	}
	return nil
}
