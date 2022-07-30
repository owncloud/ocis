package search

import (
	"context"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
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

func HandleEvents(eng engine.Engine, extractor content.Extractor, gw gateway.GatewayAPIClient, bus events.Stream, logger log.Logger, cfg *config.Config) error {
	ch, err := events.Consume(
		bus,
		"search",
		events.ItemTrashed{},
		events.ItemRestored{},
		events.ItemMoved{},
		events.ContainerCreated{},
		events.FileTouched{},
		events.FileVersionRestored{},
		events.FileUploaded{},
		events.UploadReady{},
	)
	if err != nil {
		return err
	}

	go func(eh *eventHandler, ch <-chan interface{}) {
		for e := range ch {
			switch ev := e.(type) {
			case events.ItemTrashed:
				eh.logger.Debug().Interface("event", e).Msg("marking document as deleted")
				eh.trashItem(ev.ID)
			case events.ItemMoved:
				eh.logger.Debug().Interface("event", e).Msg("resource has been moved, updating the document")
				eh.moveItem(ev.Ref, &user.User{Id: ev.Executant})
			case events.ItemRestored:
				eh.logger.Debug().Interface("event", e).Msg("marking document as restored")
				eh.restoreItem(ev.Ref, &user.User{Id: ev.Executant})
			case events.ContainerCreated:
				eh.logger.Debug().Interface("event", e).Msg("resource container created, updating the document")
				eh.upsertItem(ev.Ref, &user.User{Id: ev.Executant})
			case events.FileTouched:
				eh.logger.Debug().Interface("event", e).Msg("resource has been changed, updating the document")
				eh.upsertItem(ev.Ref, &user.User{Id: ev.Executant})
			case events.FileVersionRestored:
				eh.logger.Debug().Interface("event", e).Msg("resource version restored, updating the document")
				eh.upsertItem(ev.Ref, &user.User{Id: ev.Executant})
			case events.FileUploaded:
				if cfg.Events.AsyncUploads {
					return
				}

				eh.logger.Debug().Interface("event", e).Msg("resource upload ready, updating the document")
				eh.upsertItem(ev.Ref, &user.User{Id: ev.Executant})
			case events.UploadReady:
				eh.logger.Debug().Interface("event", e).Msg("async resource upload ready, updating the document")
				eh.upsertItem(ev.FileRef, ev.ExecutingUser)
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

func (eh *eventHandler) upsertItem(ref *provider.Reference, owner *user.User) {
	if ref == nil || owner == nil {
		return
	}

	statRes, err := statResource(ref, owner, eh.gateway, eh.secret, eh.logger)
	if err != nil {
		eh.logger.Error().Err(err).Msg("failed to stat the changed resource")
		return
	}
	if statRes.Status.Code != rpc.Code_CODE_OK {
		eh.logger.Error().Interface("statRes", statRes).Msg("failed to stat the changed resource")
		return
	}

	doc, err := eh.extractor.Extract(context.TODO(), ref, statRes.Info)
	if err != nil {
		eh.logger.Error().Err(err).Msg("failed to extract resource content")
		return
	}

	ent := engine.Entity{
		ID:       storagespace.FormatResourceID(*statRes.Info.Id),
		RootID:   storagespace.FormatResourceID(*ref.ResourceId),
		Path:     ref.Path,
		Type:     uint64(statRes.Info.Type),
		Document: doc,
	}

	err = eh.engine.Upsert(ent.ID, ent)
	if err != nil {
		eh.logger.Error().Err(err).Msg("error adding updating the resource in the index")
	} else {
		logDocCount(eh.engine, eh.logger)
	}
}

func (eh *eventHandler) restoreItem(ref *provider.Reference, owner *user.User) {
	if ref == nil || owner == nil {
		return
	}

	statRes, err := statResource(ref, owner, eh.gateway, eh.secret, eh.logger)
	if err != nil {
		eh.logger.Error().Err(err).Msg("failed to stat the changed resource")
		return
	}

	if statRes.Status.Code != rpc.Code_CODE_OK {
		eh.logger.Error().Interface("statRes", statRes).Msg("failed to stat the changed resource")
		return
	}

	err = eh.engine.Restore(storagespace.FormatResourceID(*statRes.Info.Id))
	if err != nil {
		eh.logger.Error().Err(err).Msg("failed to restore the changed resource in the index")
	}
}

func (eh *eventHandler) moveItem(ref *provider.Reference, owner *user.User) {
	if ref == nil || owner == nil {
		return
	}

	statRes, err := statResource(ref, owner, eh.gateway, eh.secret, eh.logger)
	if err != nil {
		eh.logger.Error().Err(err).Msg("failed to stat the moved resource")
		return
	}
	if statRes.Status.Code != rpc.Code_CODE_OK {
		eh.logger.Error().Interface("statRes", statRes).Msg("failed to stat the moved resource")
		return
	}

	gpRes, err := getPath(statRes.Info.Id, owner, eh.gateway, eh.secret, eh.logger)
	if err != nil {
		eh.logger.Error().Err(err).Interface("ref", ref).Msg("failed to get path for moved resource")
		return
	}
	if gpRes.Status.Code != rpcv1beta1.Code_CODE_OK {
		eh.logger.Error().Interface("status", gpRes.Status).Interface("ref", ref).Msg("failed to get path for moved resource")
		return
	}

	err = eh.engine.Move(storagespace.FormatResourceID(*statRes.Info.Id), gpRes.Path)
	if err != nil {
		eh.logger.Error().Err(err).Msg("failed to move the changed resource in the index")
	}
}
