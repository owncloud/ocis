package search

import (
	"context"
	"sync"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
	"github.com/owncloud/ocis/v2/services/search/pkg/indexer"
)

// SpaceDebouncer debounces operations on spaces for a configurable amount of time
type SpaceDebouncer struct {
	after   time.Duration
	f       func(id *provider.StorageSpaceId, userID *user.UserId)
	pending map[string]*time.Timer

	mutex sync.Mutex
}

// NewSpaceDebouncer returns a new SpaceDebouncer instance
func NewSpaceDebouncer(d time.Duration, f func(id *provider.StorageSpaceId, userID *user.UserId)) *SpaceDebouncer {
	return &SpaceDebouncer{
		after:   d,
		f:       f,
		pending: map[string]*time.Timer{},
	}
}

// Debounce restars the debounce timer for the given space
func (d *SpaceDebouncer) Debounce(id *provider.StorageSpaceId, userID *user.UserId) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if t := d.pending[id.OpaqueId]; t != nil {
		t.Stop()
	}

	d.pending[id.OpaqueId] = time.AfterFunc(d.after, func() {
		d.f(id, userID)
	})
}

type eventHandler struct {
	logger    log.Logger
	engine    engine.Engine
	gateway   gateway.GatewayAPIClient
	extractor content.Extractor
	indexer   indexer.Indexer
	secret    string

	indexSpaceDebouncer *SpaceDebouncer
}

// HandleEvents listens to the needed events,
// it handles the whole resource indexing livecycle.
func HandleEvents(eng engine.Engine, extractor content.Extractor, gw gateway.GatewayAPIClient, bus events.Consumer, indexer indexer.Indexer, logger log.Logger, cfg *config.Config) error {
	evts := []events.Unmarshaller{
		events.ItemTrashed{},
		events.ItemRestored{},
		events.ItemMoved{},
		events.ContainerCreated{},
		events.FileTouched{},
		events.FileVersionRestored{},
		events.TagsAdded{},
		events.TagsRemoved{},
	}

	if cfg.Events.AsyncUploads {
		evts = append(evts, events.UploadReady{})
	} else {
		evts = append(evts, events.FileUploaded{})
	}

	ch, err := events.Consume(bus, "search", evts...)
	if err != nil {
		return err
	}

	if cfg.Events.NumConsumers == 0 {
		cfg.Events.NumConsumers = 1
	}

	for i := 0; i < cfg.Events.NumConsumers; i++ {
		go func(eh *eventHandler, ch <-chan interface{}) {
			for e := range ch {
				eh.logger.Debug().Interface("event", e).Msg("updating index")

				switch ev := e.(type) {
				case events.ItemTrashed:
					eh.trashItem(ev.ID)
					eh.reindexSpace(ev, ev.Ref, ev.Executant, ev.SpaceOwner)
				case events.ItemMoved:
					eh.moveItem(ev.Ref, ev.Executant)
					eh.reindexSpace(ev, ev.Ref, ev.Executant, ev.SpaceOwner)
				case events.ItemRestored:
					eh.restoreItem(ev.Ref, ev.Executant)
					eh.reindexSpace(ev, ev.Ref, ev.Executant, ev.SpaceOwner)
				case events.ContainerCreated:
					eh.reindexSpace(ev, ev.Ref, ev.Executant, ev.SpaceOwner)
				case events.FileTouched:
					eh.reindexSpace(ev, ev.Ref, ev.Executant, ev.SpaceOwner)
				case events.FileVersionRestored:
					eh.reindexSpace(ev, ev.Ref, ev.Executant, ev.SpaceOwner)
				case events.FileUploaded:
					eh.reindexSpace(ev, ev.Ref, ev.Executant, ev.SpaceOwner)
				case events.UploadReady:
					eh.reindexSpace(ev, ev.FileRef, ev.ExecutingUser.Id, ev.SpaceOwner)
				case events.TagsAdded:
					eh.reindexSpace(ev, ev.Ref, ev.Executant, ev.SpaceOwner)
				case events.TagsRemoved:
					eh.reindexSpace(ev, ev.Ref, ev.Executant, ev.SpaceOwner)
				}
			}
		}(
			&eventHandler{
				logger:    logger,
				engine:    eng,
				secret:    cfg.MachineAuthAPIKey,
				gateway:   gw,
				extractor: extractor,
				indexer:   indexer,
				indexSpaceDebouncer: NewSpaceDebouncer(50*time.Millisecond, func(id *provider.StorageSpaceId, userID *user.UserId) {
					err := indexer.IndexSpace(context.Background(), id, userID)
					if err != nil {
						logger.Error().Err(err).Interface("spaceID", id).Interface("userID", userID).Msg("error while indexing a space")
					}
				}),
			},
			ch,
		)
	}

	return nil
}

func (eh *eventHandler) trashItem(rid *provider.ResourceId) {
	err := eh.engine.Delete(storagespace.FormatResourceID(*rid))
	if err != nil {
		eh.logger.Error().Err(err).Interface("Id", rid).Msg("failed to remove item from index")
	}
}

func (eh *eventHandler) reindexSpace(ev interface{}, ref *provider.Reference, executant, owner *user.UserId) {
	eh.logger.Debug().Interface("event", ev).Msg("resource has been changed, scheduling a space resync")
	spaceID := &provider.StorageSpaceId{
		OpaqueId: storagespace.FormatResourceID(provider.ResourceId{
			StorageId: ref.GetResourceId().GetStorageId(),
			SpaceId:   ref.GetResourceId().GetSpaceId(),
		}),
	}
	if owner != nil {
		eh.indexSpaceDebouncer.Debounce(spaceID, owner)
	} else {
		eh.indexSpaceDebouncer.Debounce(spaceID, executant)
	}
}

func (eh *eventHandler) restoreItem(ref *provider.Reference, uid *user.UserId) {
	ctx, stat, path := eh.resInfo(uid, ref)
	if ctx == nil || stat == nil || path == "" {
		return
	}

	if err := eh.engine.Restore(storagespace.FormatResourceID(*stat.Info.Id)); err != nil {
		eh.logger.Error().Err(err).Msg("failed to restore the changed resource in the index")
	}
}

func (eh *eventHandler) moveItem(ref *provider.Reference, uid *user.UserId) {
	ctx, stat, path := eh.resInfo(uid, ref)
	if ctx == nil || stat == nil || path == "" {
		return
	}

	if err := eh.engine.Move(storagespace.FormatResourceID(*stat.GetInfo().GetId()), storagespace.FormatResourceID(*stat.GetInfo().GetParentId()), path); err != nil {
		eh.logger.Error().Err(err).Msg("failed to move the changed resource in the index")
	}
}

func (eh *eventHandler) resInfo(uid *user.UserId, ref *provider.Reference) (context.Context, *provider.StatResponse, string) {
	ownerCtx, err := getAuthContext(&user.User{Id: uid}, eh.gateway, eh.secret, eh.logger)
	if err != nil {
		return nil, nil, ""
	}

	statRes, err := statResource(ownerCtx, ref, eh.gateway, eh.logger)
	if err != nil {
		return nil, nil, ""
	}

	r, err := ResolveReference(ownerCtx, ref, statRes.GetInfo(), eh.gateway)
	if err != nil {
		return nil, nil, ""
	}

	return ownerCtx, statRes, r.GetPath()
}
