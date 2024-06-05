package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
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
	gws    pool.Selectable[gateway.GatewayAPIClient]
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
		gws:    o.GatewaySelector,
	}

	return s, nil
}

// Run runs the service
func (a *ActivitylogService) Run() error {
	for e := range a.events {
		var err error
		switch ev := e.Event.(type) {
		case events.UploadReady:
			err = a.addActivity(ev.FileRef, e.ID, utils.TSToTime(ev.Timestamp))
		}

		if err != nil {
			a.log.Error().Err(err).Interface("event", e).Msg("could not process event")
		}
	}
	return nil
}

func (a *ActivitylogService) addActivity(initRef *provider.Reference, eventID string, timestamp time.Time) error {
	gwc, err := a.gws.Next()
	if err != nil {
		return fmt.Errorf("cant get gateway client: %w", err)
	}

	ctx, err := utils.GetServiceUserContext(a.cfg.ServiceAccount.ServiceAccountID, gwc, a.cfg.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return fmt.Errorf("cant get service user context: %w", err)
	}

	var info *provider.ResourceInfo
	depth, ref := 0, initRef
	for {
		if err := a.addActivityToReference(ref, eventID, depth, timestamp); err != nil {
			return fmt.Errorf("could not store activity: %w", err)
		}

		if info != nil && utils.IsSpaceRoot(info) {
			return nil
		}

		info, err = utils.GetResource(ctx, ref, gwc)
		if err != nil {
			return fmt.Errorf("could not get resource info: %w", err)
		}

		depth++
		ref = &provider.Reference{ResourceId: info.GetParentId()}
	}
}

func (a *ActivitylogService) addActivityToReference(ref *provider.Reference, eventID string, depth int, timestamp time.Time) error {
	fileID, err := storagespace.FormatReference(ref)
	if err != nil {
		return err
	}

	return a.storeActivity(fileID, Activity{
		EventID:   eventID,
		Depth:     depth,
		Timestamp: timestamp,
	})
}

func (a *ActivitylogService) storeActivity(resourceID string, activity Activity) error {
	records, err := a.store.Read(resourceID)
	if err != nil {
		return err
	}

	var activities []Activity
	if len(records) > 0 {
		if err := json.Unmarshal(records[0].Value, &activities); err != nil {
			return err
		}
	}

	// TODO: max len check?
	activities = append(activities, activity)

	b, err := json.Marshal(activities)
	if err != nil {
		return err
	}

	return a.store.Write(&microstore.Record{
		Key:   resourceID,
		Value: b,
	})
}
