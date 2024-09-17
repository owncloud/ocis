package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	microstore "go-micro.dev/v4/store"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/config"
)

// RawActivity represents an activity as it is stored in the activitylog store
type RawActivity struct {
	EventID   string    `json:"event_id"`
	Depth     int       `json:"depth"`
	Timestamp time.Time `json:"timestamp"`
}

// ActivitylogService logs events per resource
type ActivitylogService struct {
	cfg        *config.Config
	log        log.Logger
	events     <-chan events.Event
	store      microstore.Store
	gws        pool.Selectable[gateway.GatewayAPIClient]
	mux        *chi.Mux
	evHistory  ehsvc.EventHistoryService
	valService settingssvc.ValueService
	lock       sync.RWMutex

	registeredEvents map[string]events.Unmarshaller
}

// New creates a new ActivitylogService
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
		log:              o.Logger,
		cfg:              o.Config,
		events:           ch,
		store:            o.Store,
		gws:              o.GatewaySelector,
		mux:              o.Mux,
		evHistory:        o.HistoryClient,
		valService:       o.ValueClient,
		lock:             sync.RWMutex{},
		registeredEvents: make(map[string]events.Unmarshaller),
	}

	s.mux.Get("/graph/v1beta1/extensions/org.libregraph/activities", s.HandleGetItemActivities)

	for _, e := range o.RegisteredEvents {
		typ := reflect.TypeOf(e)
		s.registeredEvents[typ.String()] = e
	}

	go s.Run()

	return s, nil
}

// Run runs the service
func (a *ActivitylogService) Run() {
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
			err = a.RemoveResource(ev.ID)
		case events.ItemMoved:
			err = a.AddActivity(ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		case events.ShareCreated:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, utils.TSToTime(ev.CTime))
		case events.ShareRemoved:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, ev.Timestamp)
		case events.LinkCreated:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, utils.TSToTime(ev.CTime))
		case events.LinkUpdated:
			if ev.Sharer != nil && ev.ItemID != nil && ev.Sharer.GetOpaqueId() != ev.ItemID.GetSpaceId() {
				err = a.AddActivity(toRef(ev.ItemID), e.ID, utils.TSToTime(ev.CTime))
			}
		case events.LinkRemoved:
			err = a.AddActivity(toRef(ev.ItemID), e.ID, utils.TSToTime(ev.Timestamp))
		case events.SpaceShared:
			err = a.AddSpaceActivity(ev.ID, e.ID, ev.Timestamp)
		case events.SpaceUnshared:
			err = a.AddSpaceActivity(ev.ID, e.ID, ev.Timestamp)
		}

		if err != nil {
			a.log.Error().Err(err).Interface("event", e).Msg("could not process event")
		}
	}
}

// AddActivity adds the activity to the given resource and all its parents
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

// AddActivityTrashed adds the activity to given trashed resource and all its former parents
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
	if err := a.storeActivity(storagespace.FormatResourceID(resourceID), eventID, 0, timestamp); err != nil {
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

// AddSpaceActivity adds the activity to the given spaceroot
func (a *ActivitylogService) AddSpaceActivity(spaceID *provider.StorageSpaceId, eventID string, timestamp time.Time) error {
	// spaceID is in format <providerid>$<spaceid>
	// activitylog service uses format <providerid>$<spaceid>!<resourceid>
	// lets do some converting, shall we?
	rid, err := storagespace.ParseID(spaceID.GetOpaqueId())
	if err != nil {
		return fmt.Errorf("could not parse space id: %w", err)
	}
	rid.OpaqueId = rid.GetSpaceId()
	return a.storeActivity(storagespace.FormatResourceID(&rid), eventID, 0, timestamp)

}

// Activities returns the activities for the given resource
func (a *ActivitylogService) Activities(rid *provider.ResourceId) ([]RawActivity, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.activities(rid)
}

// RemoveActivities removes the activities from the given resource
func (a *ActivitylogService) RemoveActivities(rid *provider.ResourceId, toDelete map[string]struct{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	curActivities, err := a.activities(rid)
	if err != nil {
		return err
	}

	var acts []RawActivity
	for _, a := range curActivities {
		if _, ok := toDelete[a.EventID]; !ok {
			acts = append(acts, a)
		}
	}

	b, err := json.Marshal(acts)
	if err != nil {
		return err
	}

	return a.store.Write(&microstore.Record{
		Key:   storagespace.FormatResourceID(rid),
		Value: b,
	})
}

// RemoveResource removes the resource from the store
func (a *ActivitylogService) RemoveResource(rid *provider.ResourceId) error {
	if rid == nil {
		return fmt.Errorf("resource id is required")
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	return a.store.Delete(storagespace.FormatResourceID(rid))
}

func (a *ActivitylogService) activities(rid *provider.ResourceId) ([]RawActivity, error) {
	resourceID := storagespace.FormatResourceID(rid)

	records, err := a.store.Read(resourceID)
	if err != nil && err != microstore.ErrNotFound {
		return nil, fmt.Errorf("could not read activities: %w", err)
	}

	if len(records) == 0 {
		return []RawActivity{}, nil
	}

	var activities []RawActivity
	if err := json.Unmarshal(records[0].Value, &activities); err != nil {
		return nil, fmt.Errorf("could not unmarshal activities: %w", err)
	}

	return activities, nil
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

		if err := a.storeActivity(storagespace.FormatResourceID(info.GetId()), eventID, depth, timestamp); err != nil {
			return fmt.Errorf("could not store activity: %w", err)
		}

		if info != nil && utils.IsSpaceRoot(info) {
			return nil
		}

		depth++
		ref = &provider.Reference{ResourceId: info.GetParentId()}
	}
}

func (a *ActivitylogService) storeActivity(resourceID string, eventID string, depth int, timestamp time.Time) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	records, err := a.store.Read(resourceID)
	if err != nil && err != microstore.ErrNotFound {
		return err
	}

	var activities []RawActivity
	if len(records) > 0 {
		if err := json.Unmarshal(records[0].Value, &activities); err != nil {
			return err
		}
	}

	// TODO: max len check?
	activities = append(activities, RawActivity{
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

func toSpace(r *provider.Reference) *provider.StorageSpaceId {
	return &provider.StorageSpaceId{
		OpaqueId: storagespace.FormatStorageID(r.GetResourceId().GetStorageId(), r.GetResourceId().GetSpaceId()),
	}
}
