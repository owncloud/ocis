// Copyright 2018-2024 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package upload

import (
	"context"
	"os"
	"path/filepath"
	"time"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"

	"github.com/owncloud/reva/v2/pkg/autoprop"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rhttp/datatx/metrics"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/utils"
)

// Start subscribes to the event stream and launches numConsumers goroutines
// that process postprocessing events.
func (c *coordinator) Start(stream events.Consumer) error {
	ch, err := events.Consume(
		stream,
		c.conGroup,
		events.PostprocessingFinished{},
		events.PostprocessingStepFinished{},
		events.RestartPostprocessing{},
		events.CleanUpload{},
	)
	if err != nil {
		return err
	}
	for i := 0; i < c.numConc; i++ {
		go c.postprocessingLoop(ch)
	}
	return nil
}

func (c *coordinator) postprocessingLoop(ch <-chan events.Event) {
	for event := range ch {
		c.processEvent(context.Background(), event)
	}
}

func (c *coordinator) processEvent(evCtx context.Context, event events.Event) {
	ctx, span := events.TraceEventConsumerWithTracer(evCtx, tracer, event)
	ctx = autoprop.SetMetaToContext(ctx, event.ExtraInfo)
	defer span.End()

	switch ev := event.Event.(type) {
	case events.PostprocessingFinished:
		c.handlePostprocessingFinished(ctx, ev)
	case events.PostprocessingStepFinished:
		c.handlePostprocessingStepFinished(ctx, ev)
	case events.RestartPostprocessing:
		c.handleRestartPostprocessing(ctx, ev)
	case events.CleanUpload:
		c.handleCleanUpload(ctx, ev)
	default:
		c.log.Error().Interface("event", ev).Msg("coordinator: unknown event")
	}
}

func (c *coordinator) handlePostprocessingFinished(ctx context.Context, ev events.PostprocessingFinished) {
	log := c.log.With().Str("event", "PostprocessingFinished").Str("uploadid", ev.UploadID).Logger()
	if ev.ResourceID != nil && ev.ResourceID.GetStorageId() != "" && ev.ResourceID.GetStorageId() != c.mountID {
		log.Debug().Msg("ignoring event for different storage")
		return
	}
	session, err := c.store.Get(ctx, ev.UploadID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get upload")
		// Session file gone (e.g. coordinator restarted mid-postprocessing).
		// Clear the processing flag directly using the node ID from the event so
		// the node does not stay stuck returning 429 Too Early forever.
		if ev.ResourceID != nil && ev.ResourceID.GetOpaqueId() != "" {
			ref := provider.Reference{ResourceId: ev.ResourceID}
			if mpErr := c.fs.MarkProcessing(ctx, &ref, false, ev.UploadID); mpErr != nil {
				log.Error().Err(mpErr).Msg("could not unmark processing after lost session")
			}
		}
		return
	}

	ctx = session.Context(ctx)

	log = c.log.With().Str("spaceid", session.SpaceID()).Str("nodeid", session.NodeID()).Logger()
	ref := session.Reference()
	if _, mdErr := c.fs.GetMD(ctx, &ref, []string{}, []string{}); mdErr != nil {
		if _, notFound := mdErr.(errtypes.IsNotFound); notFound {
			log.Debug().Err(mdErr).Msg("node deleted during postprocessing; cleaning up")
			session.Cleanup(true, true)
			if err := c.fs.MarkProcessing(ctx, &ref, false, session.ID()); err != nil {
				log.Error().Err(err).Msg("could not unmark processing during cleanup of deleted node")
			}
			return
		}
	}

	var (
		failed             bool
		revertNodeMetadata bool
		keepUpload         bool
		retryCommit        bool
	)

	switch ev.Outcome {
	default:
		log.Error().Str("outcome", string(ev.Outcome)).Msg("unknown postprocessing outcome - aborting")
		fallthrough
	case events.PPOutcomeAbort:
		failed = true
		revertNodeMetadata = !session.NodeExists()
		keepUpload = true
		metrics.UploadSessionsAborted.Inc()
	case events.PPOutcomeContinue:
		f, fopenErr := os.Open(session.BinPath())
		if fopenErr != nil {
			log.Error().Err(fopenErr).Msg("could not open staged binary for CommitUpload")
			failed = true
			keepUpload = true
			retryCommit = true
		} else {
			defer f.Close()
			commitRef := session.Reference()
			_, commitErr := c.fs.CommitUpload(ctx, &commitRef, storage.UploadSource{
				Body:      f,
				Length:    session.Size(),
				Metadata:  session.Metadata(),
				Checksums: session.Checksums(),
			})
			if commitErr != nil {
				log.Error().Err(commitErr).Msg("could not commit upload")
				failed = true
				keepUpload = true
				retryCommit = true
			} else {
				metrics.UploadSessionsFinalized.Inc()
			}
		}
	case events.PPOutcomeDelete:
		failed = true
		revertNodeMetadata = !session.NodeExists()
		metrics.UploadSessionsDeleted.Inc()
	}

	now := time.Now()

	session.Cleanup(!keepUpload, !keepUpload)

	nodeRef := session.Reference()
	if !retryCommit {
		if err := c.fs.MarkProcessing(ctx, &nodeRef, false, session.ID()); err != nil {
			log.Error().Err(err).Msg("could not unmark processing after postprocessing finished")
		}
		if revertNodeMetadata {
			if _, delErr := c.fs.Delete(ctx, &nodeRef); delErr != nil {
				if _, ok := delErr.(errtypes.NotFound); !ok {
					log.Error().Err(delErr).Msg("could not delete placeholder node on abort")
				}
			}
		}
	}

	var isVersion bool
	if session.NodeExists() {
		info, err := session.GetInfo(ctx)
		if err == nil && info.MetaData["versionsPath"] != "" {
			isVersion = true
		}
	}

	if err := events.Publish(
		ctx,
		c.pub,
		events.UploadReady{
			UploadID:      ev.UploadID,
			Failed:        failed,
			ExecutingUser: ev.ExecutingUser,
			Filename:      ev.Filename,
			FileRef: &provider.Reference{
				ResourceId: &provider.ResourceId{
					StorageId: session.ProviderID(),
					SpaceId:   session.SpaceID(),
					OpaqueId:  session.SpaceID(),
				},
				Path: utils.MakeRelativePath(filepath.Join(session.Dir(), session.Filename())),
			},
			ResourceID: &provider.ResourceId{
				StorageId: session.ProviderID(),
				SpaceId:   session.SpaceID(),
				OpaqueId:  session.NodeID(),
			},
			Timestamp:         utils.TimeToTS(now),
			SpaceOwner:        session.SpaceOwner(),
			IsVersion:         isVersion,
			ImpersonatingUser: ev.ImpersonatingUser,
		},
	); err != nil {
		log.Error().Err(err).Msg("Failed to publish UploadReady event")
	}
}

func (c *coordinator) handleRestartPostprocessing(ctx context.Context, ev events.RestartPostprocessing) {
	log := c.log.With().Str("event", "RestartPostprocessing").Str("uploadid", ev.UploadID).Logger()
	session, err := c.store.Get(ctx, ev.UploadID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get upload")
		return
	}
	ctx = session.Context(ctx)
	log = c.log.With().Str("spaceid", session.SpaceID()).Str("nodeid", session.NodeID()).Logger()
	s, err := session.URL(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not create url")
		return
	}

	metrics.UploadSessionsRestarted.Inc()

	if err := events.Publish(ctx, c.pub, events.BytesReceived{
		UploadID:      session.ID(),
		URL:           s,
		SpaceOwner:    session.SpaceOwner(),
		ExecutingUser: &user.User{Id: &user.UserId{OpaqueId: "postprocessing-restart"}},
		ResourceID: &provider.ResourceId{
			SpaceId:  session.SpaceID(),
			OpaqueId: session.NodeID(),
		},
		Filename: session.Filename(),
		Filesize: uint64(session.Size()),
	}); err != nil {
		log.Error().Err(err).Msg("Failed to publish BytesReceived event")
	}
}

func (c *coordinator) handleCleanUpload(ctx context.Context, ev events.CleanUpload) {
	log := c.log.With().Str("event", "CleanUpload").Str("uploadid", ev.UploadID).Logger()
	session, err := c.store.Get(ctx, ev.UploadID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get upload")
		return
	}
	ctx = session.Context(ctx)
	session.Cleanup(!ev.KeepUpload, !ev.KeepUpload)
	nodeRef := session.Reference()
	if err := c.fs.MarkProcessing(ctx, &nodeRef, false, session.ID()); err != nil {
		log.Error().Err(err).Msg("could not unmark processing during CleanUpload")
	}
	if !session.NodeExists() {
		if _, delErr := c.fs.Delete(ctx, &nodeRef); delErr != nil {
			if _, ok := delErr.(errtypes.NotFound); !ok {
				log.Error().Err(delErr).Msg("could not delete placeholder node during CleanUpload")
			}
		}
	}
}

func (c *coordinator) handlePostprocessingStepFinished(ctx context.Context, ev events.PostprocessingStepFinished) {
	log := c.log.With().Str("event", "PostprocessingStepFinished").Str("uploadid", ev.UploadID).Logger()
	if ev.ResourceID != nil && ev.ResourceID.GetStorageId() != "" && ev.ResourceID.GetStorageId() != c.mountID {
		log.Debug().Msg("ignoring event for different storage")
		return
	}
	if ev.FinishedStep != events.PPStepAntivirus {
		return
	}

	res, ok := ev.Result.(events.VirusscanResult)
	if !ok {
		log.Error().Msgf("coordinator: unexpected antivirus result type %T", ev.Result)
		return
	}
	if res.ErrorMsg != "" {
		return
	}
	log = c.log.With().Str("scan_description", res.Description).Bool("infected", res.Infected).Logger()

	if ev.UploadID == "" {
		// on-demand scanning not supported
		return
	}

	session, err := c.store.Get(ctx, ev.UploadID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get upload")
		return
	}
	log = c.log.With().Str("spaceid", session.SpaceID()).Str("nodeid", session.NodeID()).Logger()

	session.SetScanData(res.Description, res.Scandate)
	if err := session.Persist(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to persist scan results")
	}

	ctx = session.Context(ctx)
	ref := session.Reference()
	if err := c.fs.SetArbitraryMetadata(ctx, &ref, &provider.ArbitraryMetadata{
		Metadata: map[string]string{
			"scanstatus": res.Description,
			"scandate":   res.Scandate.Format(time.RFC3339Nano),
		},
	}); err != nil {
		log.Error().Err(err).Msg("Failed to write scan results to node")
	}

	metrics.UploadSessionsScanned.Inc()
}
