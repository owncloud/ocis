package search

import (
	"context"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
)

type eventHandler struct {
	logger    log.Logger
	engine    engine.Engine
	gateway   gateway.GatewayAPIClient
	extractor content.Extractor
	secret    string
}

// HandleEvents listens to the needed events,
// it handles the whole resource indexing livecycle.
func HandleEvents(eng engine.Engine, extractor content.Extractor, gw gateway.GatewayAPIClient, bus events.Stream, logger log.Logger, cfg *config.Config) error {
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

	go func(eh *eventHandler, ch <-chan interface{}) {
		for e := range ch {
			eh.logger.Debug().Interface("event", e).Msg("updating index")

			switch ev := e.(type) {
			case events.ItemTrashed:
				eh.trashItem(ev.ID)
			case events.ItemMoved:
				eh.moveItem(ev.Ref, ev.Executant)
			case events.ItemRestored:
				eh.restoreItem(ev.Ref, ev.Executant)
			case events.ContainerCreated:
				eh.upsertItem(ev.Ref, ev.Executant)
			case events.FileTouched:
				eh.upsertItem(ev.Ref, ev.Executant)
			case events.FileVersionRestored:
				eh.upsertItem(ev.Ref, ev.Executant)
			case events.FileUploaded:
				eh.upsertItem(ev.Ref, ev.Executant)
			case events.UploadReady:
				eh.upsertItem(ev.FileRef, ev.ExecutingUser.Id)
			case events.TagsAdded:
				eh.upsertItem(ev.Ref, ev.Executant)
			case events.TagsRemoved:
				eh.upsertItem(ev.Ref, ev.Executant)
			}
		}
	}(
		&eventHandler{
			logger:    logger,
			engine:    eng,
			secret:    cfg.MachineAuthAPIKey,
			gateway:   gw,
			extractor: extractor,
		},
		ch,
	)

	return nil
}

func (eh *eventHandler) trashItem(rid *provider.ResourceId) {
	err := eh.engine.Delete(storagespace.FormatResourceID(*rid))
	if err != nil {
		eh.logger.Error().Err(err).Interface("Id", rid).Msg("failed to remove item from index")
	}
}

func (eh *eventHandler) upsertItem(ref *provider.Reference, uid *user.UserId) {
	ctx, stat, path := eh.resInfo(uid, ref)
	if ctx == nil || stat == nil || path == nil {
		return
	}

	doc, err := eh.extractor.Extract(ctx, stat.Info)
	if err != nil {
		eh.logger.Error().Err(err).Msg("failed to extract resource content")
		return
	}

	r := engine.Resource{
		ID: storagespace.FormatResourceID(*stat.Info.Id),
		RootID: storagespace.FormatResourceID(provider.ResourceId{
			StorageId: stat.Info.Id.StorageId,
			OpaqueId:  stat.Info.Id.SpaceId,
			SpaceId:   stat.Info.Id.SpaceId,
		}),
		Path:     utils.MakeRelativePath(path.Path),
		Type:     uint64(stat.Info.Type),
		Document: doc,
	}

	if err = eh.engine.Upsert(r.ID, r); err != nil {
		eh.logger.Error().Err(err).Msg("error adding updating the resource in the index")
	} else {
		logDocCount(eh.engine, eh.logger)
	}
}

func (eh *eventHandler) restoreItem(ref *provider.Reference, uid *user.UserId) {
	ctx, stat, path := eh.resInfo(uid, ref)
	if ctx == nil || stat == nil || path == nil {
		return
	}

	if err := eh.engine.Restore(storagespace.FormatResourceID(*stat.Info.Id)); err != nil {
		eh.logger.Error().Err(err).Msg("failed to restore the changed resource in the index")
	}
}

func (eh *eventHandler) moveItem(ref *provider.Reference, uid *user.UserId) {
	ctx, stat, path := eh.resInfo(uid, ref)
	if ctx == nil || stat == nil || path == nil {
		return
	}

	if err := eh.engine.Move(storagespace.FormatResourceID(*stat.Info.Id), path.Path); err != nil {
		eh.logger.Error().Err(err).Msg("failed to move the changed resource in the index")
	}
}

func (eh *eventHandler) resInfo(uid *user.UserId, ref *provider.Reference) (context.Context, *provider.StatResponse, *provider.GetPathResponse) {
	ownerCtx, err := getAuthContext(&user.User{Id: uid}, eh.gateway, eh.secret, eh.logger)
	if err != nil {
		return nil, nil, nil
	}

	statRes, err := statResource(ownerCtx, ref, eh.gateway, eh.logger)
	if err != nil {
		return nil, nil, nil
	}

	pathRes, err := getPath(ownerCtx, statRes.Info.Id, eh.gateway, eh.logger)
	if err != nil {
		return nil, nil, nil
	}

	return ownerCtx, statRes, pathRes
}
