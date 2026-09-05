package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/reva/v2/pkg/autoprop"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/shamaton/msgpack/v2"
	microstore "go-micro.dev/v4/store"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/config"
)

// Nats runs into max payload exceeded errors at around 7k activities. Let's keep a buffer.
var _maxActivities = 6000

// RawActivity represents an activity as it is stored in the activitylog store
type RawActivity struct {
	EventID   string    `json:"event_id"`
	Depth     int       `json:"depth"`
	Timestamp time.Time `json:"timestamp"`
}

// Debouncer coalesces activities per resource id over a configurable window
// before flushing them to the store in a single read-modify-write. This removes
// most of the per-event marshal/unmarshal/write cycles that make the activitylog
// consumer fall behind under bursty load.
type Debouncer struct {
	after time.Duration
	flush func(id string, activities []RawActivity) error
	log   log.Logger

	mutex   sync.Mutex
	pending map[string]*queueItem
}

type queueItem struct {
	activities []RawActivity
	timer      *time.Timer
}

// NewDebouncer returns a new Debouncer. A duration of 0 flushes synchronously,
// which preserves the previous (un-buffered) behaviour and is convenient for tests.
func NewDebouncer(after time.Duration, logger log.Logger, flush func(id string, activities []RawActivity) error) *Debouncer {
	return &Debouncer{
		after:   after,
		flush:   flush,
		log:     logger,
		pending: make(map[string]*queueItem),
	}
}

// Debounce queues an activity for the given resource id. With a zero window the
// activity is flushed immediately. Otherwise the flush happens once the window
// has elapsed since the first queued activity for that id; further activities in
// the window are appended and flushed together.
func (d *Debouncer) Debounce(id string, ra RawActivity) {
	if d.after == 0 {
		if err := d.flush(id, []RawActivity{ra}); err != nil {
			d.log.Error().Err(err).Str("resourceid", id).Msg("could not store activity")
		}
		return
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	item, ok := d.pending[id]
	if !ok {
		item = &queueItem{}
		d.pending[id] = item
		item.timer = time.AfterFunc(d.after, func() { d.flushItem(id) })
	}
	item.activities = append(item.activities, ra)
}

// flushItem removes the queued activities for id under the lock and flushes them
// outside of it, so the (potentially slow) store write never blocks Debounce and
// the activities slice is never read and appended to concurrently.
func (d *Debouncer) flushItem(id string) {
	d.mutex.Lock()
	item, ok := d.pending[id]
	if !ok {
		d.mutex.Unlock()
		return
	}
	delete(d.pending, id)
	activities := item.activities
	d.mutex.Unlock()

	if err := d.flush(id, activities); err != nil {
		d.log.Error().Err(err).Str("resourceid", id).Msg("could not store buffered activities")
	}
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

	debouncer *Debouncer

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
	s.debouncer = NewDebouncer(o.Config.WriteBufferDuration, o.Logger, s.storeActivity)

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

		// trace provider is available here, otherwise the activitylog service should have crashed
		tp, _ := tracing.GetServiceTraceProvider(a.cfg.Tracing, a.cfg.Service.Name)
		evCtx := context.Background()
		ctx, span := events.TraceEventConsumer(evCtx, tp, e)
		ctx = autoprop.SetMetaToContext(ctx, e.ExtraInfo)
		// span is closed at the end of the loop

		switch ev := e.Event.(type) {
		case events.UploadReady:
			err = a.AddActivity(ctx, ev.FileRef, e.ID, utils.TSToTime(ev.Timestamp))
		case events.FileTouched:
			err = a.AddActivity(ctx, ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		// Disabled https://github.com/owncloud/ocis/issues/10293
		//case events.FileDownloaded:
		// we are only interested in public link downloads - so no need to store others.
		//if ev.ImpersonatingUser.GetDisplayName() == "Public" {
		//	err = a.AddActivity(ctx, ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		//}
		case events.ContainerCreated:
			err = a.AddActivity(ctx, ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		case events.ItemTrashed:
			err = a.AddActivityTrashed(ctx, ev.ID, ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		case events.ItemPurged:
			err = a.RemoveResource(ev.ID) // no ctx needed at the moment
		case events.ItemMoved:
			err = a.AddActivity(ctx, ev.Ref, e.ID, utils.TSToTime(ev.Timestamp))
		case events.ShareCreated:
			err = a.AddActivity(ctx, toRef(ev.ItemID), e.ID, utils.TSToTime(ev.CTime))
		case events.ShareUpdated:
			if ev.Sharer != nil && ev.ItemID != nil && ev.Sharer.GetOpaqueId() != ev.ItemID.GetSpaceId() {
				err = a.AddActivity(ctx, toRef(ev.ItemID), e.ID, utils.TSToTime(ev.MTime))
			}
		case events.ShareRemoved:
			err = a.AddActivity(ctx, toRef(ev.ItemID), e.ID, ev.Timestamp)
		case events.LinkCreated:
			err = a.AddActivity(ctx, toRef(ev.ItemID), e.ID, utils.TSToTime(ev.CTime))
		case events.LinkUpdated:
			if ev.Sharer != nil && ev.ItemID != nil && ev.Sharer.GetOpaqueId() != ev.ItemID.GetSpaceId() {
				err = a.AddActivity(ctx, toRef(ev.ItemID), e.ID, utils.TSToTime(ev.MTime))
			}
		case events.LinkRemoved:
			err = a.AddActivity(ctx, toRef(ev.ItemID), e.ID, utils.TSToTime(ev.Timestamp))
		case events.SpaceShared:
			err = a.AddSpaceActivity(ev.ID, e.ID, ev.Timestamp) // no ctx needed at the moment
		case events.SpaceUnshared:
			err = a.AddSpaceActivity(ev.ID, e.ID, ev.Timestamp) // no ctx needed at the moment
		}

		if err != nil {
			a.log.Error().Err(err).Interface("event", e).Msg("could not process event")
		}

		span.End()
	}
}

// AddActivity adds the activity to the given resource and all its parents
func (a *ActivitylogService) AddActivity(ctx context.Context, initRef *provider.Reference, eventID string, timestamp time.Time) error {
	gwc, err := a.gws.Next()
	if err != nil {
		return fmt.Errorf("cant get gateway client: %w", err)
	}

	ctx2, err := utils.GetServiceUserContextWithContext(ctx, gwc, a.cfg.ServiceAccount.ServiceAccountID, a.cfg.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return fmt.Errorf("cant get service user context: %w", err)
	}

	return a.addActivity(initRef, eventID, timestamp, func(ref *provider.Reference) (*provider.ResourceInfo, error) {
		return utils.GetResource(ctx2, ref, gwc)
	})
}

// AddActivityTrashed adds the activity to given trashed resource and all its former parents
func (a *ActivitylogService) AddActivityTrashed(ctx context.Context, resourceID *provider.ResourceId, reference *provider.Reference, eventID string, timestamp time.Time) error {
	gwc, err := a.gws.Next()
	if err != nil {
		return fmt.Errorf("cant get gateway client: %w", err)
	}

	ctx2, err := utils.GetServiceUserContextWithContext(ctx, gwc, a.cfg.ServiceAccount.ServiceAccountID, a.cfg.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return fmt.Errorf("cant get service user context: %w", err)
	}

	// store activity on trashed item
	if err := a.storeActivity(storagespace.FormatResourceID(resourceID), []RawActivity{{
		EventID:   eventID,
		Depth:     0,
		Timestamp: timestamp,
	}}); err != nil {
		return fmt.Errorf("could not store activity: %w", err)
	}

	// get previous parent
	ref := &provider.Reference{
		ResourceId: reference.GetResourceId(),
		Path:       filepath.Dir(reference.GetPath()),
	}

	return a.addActivity(ref, eventID, timestamp, func(ref *provider.Reference) (*provider.ResourceInfo, error) {
		return utils.GetResource(ctx2, ref, gwc)
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
	return a.storeActivity(storagespace.FormatResourceID(&rid), []RawActivity{{
		EventID:   eventID,
		Depth:     0,
		Timestamp: timestamp,
	}})
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

	b, err := msgpack.Marshal(acts)
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
	if err := unmarshalActivities(records[0].Value, &activities); err != nil {
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

		a.debouncer.Debounce(storagespace.FormatResourceID(info.GetId()), RawActivity{
			EventID:   eventID,
			Depth:     depth,
			Timestamp: timestamp,
		})

		if info != nil && utils.IsSpaceRoot(info) {
			return nil
		}

		depth++
		ref = &provider.Reference{ResourceId: info.GetParentId()}
	}
}

func (a *ActivitylogService) storeActivity(resourceID string, activities []RawActivity) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	records, err := a.store.Read(resourceID)
	if err != nil && err != microstore.ErrNotFound {
		return err
	}

	var existing []RawActivity
	if len(records) > 0 {
		if err := unmarshalActivities(records[0].Value, &existing); err != nil {
			return err
		}
	}

	existing = append(existing, activities...)
	if l := len(existing); l > _maxActivities {
		existing = existing[l-_maxActivities:]
	}

	b, err := msgpack.Marshal(existing)
	if err != nil {
		return err
	}

	return a.store.Write(&microstore.Record{
		Key:   resourceID,
		Value: b,
	})
}

// unmarshalActivities decodes a stored activity list. New records are written
// with msgpack; the json fallback keeps records written before the upgrade
// readable.
func unmarshalActivities(b []byte, activities *[]RawActivity) error {
	if err := msgpack.Unmarshal(b, activities); err != nil {
		return json.Unmarshal(b, activities)
	}
	return nil
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
