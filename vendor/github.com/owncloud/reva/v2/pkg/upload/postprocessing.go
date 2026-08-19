// This file holds the postprocessing result collector: the second half of an
// async upload. finishUpload stages the bytes and publishes BytesReceived; the
// postprocessing service scans them and reports back here, and only then are the
// bytes committed through the driver seam.
//
// It is driver-agnostic on purpose: this logic used to live inside decomposedfs,
// which meant only that driver could offer async uploads.

package upload

import (
	"context"
	"errors"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/rs/zerolog"

	"github.com/owncloud/reva/v2/pkg/appctx"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rhttp/datatx/metrics"
	"github.com/owncloud/reva/v2/pkg/utils"
)

// RegisteredEvents are the postprocessing events the coordinator consumes.
var RegisteredEvents = []events.Unmarshaller{
	events.PostprocessingFinished{},
	events.PostprocessingStepFinished{},
	events.RestartPostprocessing{},
	events.CleanUpload{},
}

// StartPostprocessing subscribes to postprocessing results and switches the
// coordinator over to async uploads: from here on finished uploads stage their
// bytes and wait for a scan verdict instead of committing inline.
//
// The two go together on purpose. Deferring a commit is only safe if something
// will arrive to finish it, so there is no way to enable async without a running
// consumer, and none to run a consumer that never receives work.
//
// mountID is the storage id this provider serves. Postprocessing events are
// broadcast to every provider, so events for other storages must be dropped;
// pass "" only in tests, where a single provider sees a private stream.
//
// numConsumers goroutines share the subscription. Call once, before serving
// requests. Fails without a publisher: nothing would hand uploads to
// postprocessing, so every one of them would wait for a verdict that never comes.
func (c *coordinator) StartPostprocessing(stream events.Consumer, group, mountID string, numConsumers int) error {
	if c.pub == nil {
		return errors.New("coordinator: async uploads need an event publisher")
	}
	ch, err := events.Consume(stream, group, RegisteredEvents...)
	if err != nil {
		return err
	}
	if numConsumers <= 0 {
		numConsumers = 1
	}
	c.mountID = mountID
	c.async = true
	for i := 0; i < numConsumers; i++ {
		go c.Postprocessing(ch)
	}
	return nil
}

// Postprocessing consumes postprocessing results until ch is closed. Run it in
// its own goroutine, one per configured consumer.
func (c *coordinator) Postprocessing(ch <-chan events.Event) {
	for event := range ch {
		c.processEvent(context.Background(), event)
	}
}

// servesStorage reports whether an event is for the storage this coordinator
// serves. Postprocessing runs as a separate service and broadcasts its results to
// every storage provider, so each has to recognise its own. Events that name no
// storage predate this and are accepted.
func (c *coordinator) servesStorage(id *provider.ResourceId) bool {
	if c.mountID == "" || id.GetStorageId() == "" {
		return true
	}
	return id.GetStorageId() == c.mountID
}

func (c *coordinator) processEvent(ctx context.Context, event events.Event) {
	log := appctx.GetLogger(ctx)

	switch ev := event.Event.(type) {
	case events.PostprocessingFinished:
		if !c.servesStorage(ev.ResourceID) {
			return
		}
		c.onPostprocessingFinished(ctx, ev, log)
	case events.RestartPostprocessing:
		c.onRestartPostprocessing(ctx, ev, log)
	case events.CleanUpload:
		session, err := c.store.Get(ctx, ev.UploadID)
		if err != nil {
			log.Error().Err(err).Str("uploadid", ev.UploadID).Msg("CleanUpload: could not load session")
			return
		}
		if !ev.KeepUpload {
			c.rollbackPrepared(ctx, session, session.SizeDiff())
		} else {
			c.rollbackNode(ctx, session)
			metrics.UploadProcessing.Dec()
		}
	case events.PostprocessingStepFinished:
		if !c.servesStorage(ev.ResourceID) {
			return
		}
		if ev.FinishedStep != events.PPStepAntivirus {
			// only the antivirus result is recorded on the session
			return
		}
		res, ok := ev.Result.(events.VirusscanResult)
		if !ok || res.ErrorMsg != "" {
			// the scan itself failed; PostprocessingFinished decides the outcome
			return
		}
		session, err := c.store.Get(ctx, ev.UploadID)
		if err != nil {
			// an empty upload id means an on-demand scan, which has no session
			if ev.UploadID != "" {
				log.Error().Err(err).Str("uploadid", ev.UploadID).Msg("PostprocessingStepFinished: could not load session")
			}
			return
		}
		session.SetScanData(res.Description, res.Scandate)
		metrics.UploadSessionsScanned.Inc()
		if err := session.Persist(ctx); err != nil {
			log.Error().Err(err).Str("uploadid", ev.UploadID).Msg("could not persist scan result")
		}
	}
}

// onPostprocessingFinished completes or discards an upload according to the
// outcome postprocessing reported.
func (c *coordinator) onPostprocessingFinished(ctx context.Context, ev events.PostprocessingFinished, log *zerolog.Logger) {
	session, err := c.store.Get(ctx, ev.UploadID)
	if err != nil {
		// Without the session we cannot reach the staged bytes, so they are leaked
		// here. Housekeeping cleans them up later.
		log.Error().Err(err).Str("uploadid", ev.UploadID).Msg("PostprocessingFinished: could not load session")
		return
	}
	ctx = session.Context(ctx)
	log = appctx.GetLogger(ctx)

	switch ev.Outcome {
	case events.PPOutcomeContinue:
		if err := c.finishAsync(ctx, session); err != nil {
			// Deliberately left as it is: the node stays marked and the session on
			// disk, so this can be retried with RestartPostprocessing or discarded
			// with CleanUpload rather than silently losing the bytes.
			log.Error().Err(err).Str("uploadid", ev.UploadID).Msg("could not commit upload after postprocessing")
			c.publishUploadFailed(ctx, session, ev)
		}
		return
	case events.PPOutcomeAbort:
		metrics.UploadSessionsAborted.Inc()
		// Keep the staged bytes: an abort is a transient failure and the upload can
		// be restarted with RestartPostprocessing.
		c.rollbackNode(ctx, session)
	case events.PPOutcomeDelete:
		metrics.UploadSessionsDeleted.Inc()
		c.rollbackPrepared(ctx, session, session.SizeDiff())
	default:
		log.Error().Str("outcome", string(ev.Outcome)).Str("uploadid", ev.UploadID).Msg("unknown postprocessing outcome, aborting")
		metrics.UploadSessionsAborted.Inc()
		c.rollbackNode(ctx, session)
	}

	c.publishUploadFailed(ctx, session, ev)
}

// onRestartPostprocessing re-publishes BytesReceived so a previously aborted
// upload gets another postprocessing run.
func (c *coordinator) onRestartPostprocessing(ctx context.Context, ev events.RestartPostprocessing, log *zerolog.Logger) {
	session, err := c.store.Get(ctx, ev.UploadID)
	if err != nil {
		log.Error().Err(err).Str("uploadid", ev.UploadID).Msg("RestartPostprocessing: could not load session")
		return
	}
	metrics.UploadSessionsRestarted.Inc()
	if err := c.publishBytesReceived(session.Context(ctx), session); err != nil {
		log.Error().Err(err).Str("uploadid", ev.UploadID).Msg("could not restart postprocessing")
	}
}

// rollbackNode reverts the node to its pre-upload state but keeps the staged
// bytes and the session, so postprocessing can be restarted.
func (c *coordinator) rollbackNode(ctx context.Context, session Session) {
	ref := session.Reference()
	if err := c.fs.RollbackUpload(ctx, &ref, session.ID(), rollbackInfo(session, session.SizeDiff())); err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Str("uploadid", session.ID()).Msg("could not roll back upload")
	}
	c.unmarkProcessing(ctx, session, &ref)
}

// publishUploadFailed tells consumers the upload will not become available.
// Clients wait on UploadReady, so staying silent would leave them hanging.
func (c *coordinator) publishUploadFailed(ctx context.Context, session Session, ev events.PostprocessingFinished) {
	if c.pub == nil {
		return
	}
	if err := events.Publish(ctx, c.pub, events.UploadReady{
		UploadID:      session.ID(),
		Failed:        true,
		Filename:      session.Filename(),
		SpaceOwner:    session.SpaceOwner(),
		ExecutingUser: ev.ExecutingUser,
		FileRef:       c.uploadRef(session),
		ResourceID: &provider.ResourceId{
			StorageId: session.ProviderID(),
			SpaceId:   session.SpaceID(),
			OpaqueId:  session.NodeID(),
		},
		Timestamp:         utils.TSNow(),
		IsVersion:         session.VersionCreated(),
		ImpersonatingUser: ev.ImpersonatingUser,
	}); err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Str("uploadid", session.ID()).Msg("failed to publish UploadReady event")
	}
}
