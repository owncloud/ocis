package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
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
			err = a.AddActivity(ev.FileRef, e.ID, utils.TSToTime(ev.Timestamp))
		case events.FileTouched:
			err = a.AddActivity(ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		case events.ContainerCreated:
			err = a.AddActivity(ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		case events.ItemTrashed:
			err = a.AddActivityTrashed(ev.ID, ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		case events.ItemPurged:
			err = a.AddActivity(ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		case events.ItemMoved:
			err = a.AddActivity(ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		case events.ShareCreated:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, utils.TSToTime(ev.CTime))
		case events.ShareUpdated:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, utils.TSToTime(ev.MTime))
		case events.ShareRemoved:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, ev.Timestamp)
		case events.LinkCreated:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, utils.TSToTime(ev.CTime))
		case events.LinkUpdated:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, utils.TSToTime(ev.CTime))
		case events.LinkRemoved:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, utils.TSToTime(ev.Timestamp))
		case events.SpaceShared:
			err = a.AddActivity(sToRef(ev.ID), e.ID, ev.Timestamp)
		case events.SpaceShareUpdated:
			err = a.AddActivity(sToRef(ev.ID), e.ID, ev.Timestamp)
		case events.SpaceUnshared:
			err = a.AddActivity(sToRef(ev.ID), e.ID, ev.Timestamp)
		}

		if err != nil {
			a.log.Error().Err(err).Interface("event", e).Msg("could not process event")
		}
	}
	return nil
}

// AddActivity addds the activity to the given resource and all its parents
func (a *ActivitylogService) AddActivity(initRef *provider.Reference, eventID string, timestamp time.Time) error {
	gwc, err := a.gws.Next()
	if err != nil {
		return fmt.Errorf("cant get gateway client: %w", err)
	}

	ctx, err := utils.GetServiceUserContext(a.cfg.ServiceAccount.ServiceAccountID, gwc, a.cfg.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return fmt.Errorf("cant get service user context: %w", err)
	}

	return a.addActivity(initRef, eventID, timestamp, func(ref *provider.Reference) (*provider.ResourceInfo, error) {
		return utils.GetResource(ctx, ref, gwc)
	})
}

// Activities returns the activities for the given reference
func (a *ActivitylogService) Activities(ref *provider.Reference) ([]Activity, error) {
	resourceID, err := storagespace.FormatReference(ref)
	if err != nil {
		return nil, fmt.Errorf("could not format reference: %w", err)
	}

	records, err := a.store.Read(resourceID)
	if err != nil && err != microstore.ErrNotFound {
		return nil, fmt.Errorf("could not read activities: %w", err)
	}

	if len(records) == 0 {
		return []Activity{}, nil
	}

	var activities []Activity
	if err := json.Unmarshal(records[0].Value, &activities); err != nil {
		return nil, fmt.Errorf("could not unmarshal activities: %w", err)
	}

	return activities, nil
}

// AddActivityTrashed adds the activity to trashed item
func (a *ActivitylogService) AddActivityTrashed(resourceID *provider.ResourceId, reference *provider.Reference, eventID string, timestamp time.Time) error {
	gwc, err := a.gws.Next()
	if err != nil {
		return fmt.Errorf("cant get gateway client: %w", err)
	}

	ctx, err := utils.GetServiceUserContext(a.cfg.ServiceAccount.ServiceAccountID, gwc, a.cfg.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return fmt.Errorf("cant get service user context: %w", err)
	}

	// store activity on trashed item
	if err := a.storeActivity(resourceID, eventID, 0, timestamp); err != nil {
		return fmt.Errorf("could not store activity: %w", err)
	}

	// get previous parent
	ref := &provider.Reference{
		ResourceId: reference.GetResourceId(),
		Path:       filepath.Dir(reference.GetPath()),
	}

	return a.addActivity(ref, eventID, timestamp, func(ref *provider.Reference) (*provider.ResourceInfo, error) {
		return utils.GetResource(ctx, ref, gwc)
	})
}

// note: getResource is abstracted to allow unit testing, in general this will just be utils.GetResource
func (a *ActivitylogService) addActivity(initRef *provider.Reference, eventID string, timestamp time.Time, getResource func(*provider.Reference) (*provider.ResourceInfo, error)) error {
	var (
		info  *provider.ResourceInfo
		err   error
		depth int
		ref   = initRef
	)
	for {
		info, err = getResource(ref)
		if err != nil {
			return fmt.Errorf("could not get resource info: %w", err)
		}

		if err := a.storeActivity(info.GetId(), eventID, depth, timestamp); err != nil {
			return fmt.Errorf("could not store activity: %w", err)
		}

		if info != nil && utils.IsSpaceRoot(info) {
			return nil
		}

		depth++
		ref = &provider.Reference{ResourceId: info.GetParentId()}
	}
}

func (a *ActivitylogService) storeActivity(rid *provider.ResourceId, eventID string, depth int, timestamp time.Time) error {
	resourceID := storagespace.FormatResourceID(*rid)

	records, err := a.store.Read(resourceID)
	if err != nil && err != microstore.ErrNotFound {
		return err
	}

	var activities []Activity
	if len(records) > 0 {
		if err := json.Unmarshal(records[0].Value, &activities); err != nil {
			return err
		}
	}

	// TODO: max len check?
	activities = append(activities, Activity{
		EventID:   eventID,
		Depth:     depth,
		Timestamp: timestamp,
	})

	b, err := json.Marshal(activities)
	if err != nil {
		return err
	}

	return a.store.Write(&microstore.Record{
		Key:   resourceID,
		Value: b,
	})
}

func toRef(r *provider.ResourceId) *provider.Reference {
	return &provider.Reference{
		ResourceId: r,
	}
}

func sToRef(s *provider.StorageSpaceId) *provider.Reference {
	return &provider.Reference{
		ResourceId: &provider.ResourceId{
			OpaqueId: s.GetOpaqueId(),
			SpaceId:  s.GetOpaqueId(),
		},
	}
}
