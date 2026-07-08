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

// Package upload provides the driver-agnostic upload coordinator:
// TUS session management, postprocessing event loop, lifecycle event
// publishing.
package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/rs/zerolog"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/owncloud/reva/v2/pkg/autoprop"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rhttp/datatx/metrics"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/utils"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/owncloud/reva/pkg/upload")
}

// impersonatingUser extracts the impersonating user from the context user's opaque field.
// Returns nil when no impersonation is in effect.
func impersonatingUser(ctx context.Context) *user.User {
	u, ok := ctxpkg.ContextGetUser(ctx)
	if !ok || u == nil {
		return nil
	}
	if !utils.ExistsInOpaque(u.Opaque, "impersonating-user") {
		return nil
	}
	iu := &user.User{}
	if err := utils.ReadJSONFromOpaque(u.Opaque, "impersonating-user", iu); err != nil {
		return nil
	}
	return iu
}

var errNotImplemented = tusd.NewError("ERR_NOT_IMPLEMENTED", "use InitiateUpload on the CS3 API to start a new upload", 501)

// rollback unmarks processing, cleans up session files, and deletes the placeholder
// node if it was created by this upload (NodeExists=false at initiation).
func (c *coordinator) rollback(ctx context.Context, session Session) {
	ref := session.Reference()
	_ = c.fs.MarkProcessing(ctx, &ref, false, session.ID())
	session.Cleanup(true, true)
	if !session.NodeExists() {
		_, _ = c.fs.Delete(ctx, &ref)
	}
}

// commitSync runs CommitUpload inline and cleans up the session.
// Used by the sync path (async=false) from FinishUpload (TUS) and Upload (simple PUT).
func (c *coordinator) commitSync(ctx context.Context, session Session) error {
	ref := session.Reference()
	f, err := os.Open(session.BinPath())
	if err != nil {
		c.rollback(ctx, session)
		return err
	}
	if _, err := c.fs.CommitUpload(ctx, &ref, storage.UploadSource{
		Body:      f,
		Length:    session.Size(),
		Metadata:  session.Metadata(),
		Checksums: session.Checksums(),
	}); err != nil {
		c.rollback(ctx, session)
		return err
	}
	_ = c.fs.MarkProcessing(ctx, &ref, false, session.ID())
	session.Cleanup(true, false)
	metrics.UploadSessionsFinalized.Inc()
	return nil
}

// Coordinator owns the full upload lifecycle: session initiation, TUS data transfer,
// postprocessing event loop, and UploadReady publishing.
type Coordinator interface {
	InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error)
	Upload(ctx context.Context, req storage.UploadRequest, uff storage.UploadFinishedFunc) (*provider.ResourceInfo, error)
	GetUpload(ctx context.Context, id string) (tusd.Upload, error)
	UseIn(composer *tusd.StoreComposer)
	ListUploadSessions(ctx context.Context, filter storage.UploadSessionFilter) ([]storage.UploadSession, error)
	Start(stream events.Consumer) error
}

// coordinator is the concrete implementation of Coordinator.
type coordinator struct {
	fs       storage.FS
	store    SessionStore
	pub      events.Publisher
	async    bool
	mountID  string
	numConc  int
	conGroup string
	log      *zerolog.Logger
}

// NewCoordinator constructs a Coordinator. Call Start to begin consuming events.
// async=true requires a non-nil pub.
func NewCoordinator(
	fs storage.FS,
	store SessionStore,
	pub events.Publisher,
	async bool,
	mountID string,
	consumerGroup string,
	numConsumers int,
	log *zerolog.Logger,
) (Coordinator, error) {
	if async && pub == nil {
		return nil, fmt.Errorf("need event stream for async upload processing")
	}
	if numConsumers <= 0 {
		numConsumers = 1
	}
	return &coordinator{
		fs:       fs,
		store:    store,
		pub:      pub,
		async:    async,
		mountID:  mountID,
		numConc:  numConsumers,
		conGroup: consumerGroup,
		log:      log,
	}, nil
}

// Start subscribes to the event stream and launches numConsumers goroutines
// that process postprocessing events.
func (c *coordinator)Start(stream events.Consumer) error {
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

func (c *coordinator)postprocessingLoop(ch <-chan events.Event) {
	for event := range ch {
		c.processEvent(context.Background(), event)
	}
}

func (c *coordinator)processEvent(evCtx context.Context, event events.Event) {
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

func (c *coordinator)handlePostprocessingFinished(ctx context.Context, ev events.PostprocessingFinished) {
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
		// Only revert node metadata for new files. For overwrites the node still
		// holds the previous content.
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
		// Only revert node metadata for new files. For overwrites the node still
		// holds the previous content.
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

func (c *coordinator)handleRestartPostprocessing(ctx context.Context, ev events.RestartPostprocessing) {
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

func (c *coordinator)handleCleanUpload(ctx context.Context, ev events.CleanUpload) {
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

func (c *coordinator)handlePostprocessingStepFinished(ctx context.Context, ev events.PostprocessingStepFinished) {
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

// InitiateUpload creates a node placeholder via TouchFile and builds an upload session.
func (c *coordinator)InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	existing, err := c.fs.GetMD(ctx, ref, []string{}, []string{})
	var nodeExists bool
	switch err.(type) {
	case nil:
		nodeExists = true
	case errtypes.IsNotFound:
		nodeExists = false
	default:
		return nil, err
	}

	mtime := ""
	if m, ok := metadata["mtime"]; ok && m != "null" {
		mtime = m
	}

	var nodeID, spaceID, parentID, dir, nodeName string
	var spaceOwner *user.UserId

	if nodeExists {
		nodeID = existing.GetId().GetOpaqueId()
		spaceID = existing.GetId().GetSpaceId()
		parentID = existing.GetParentId().GetOpaqueId()
		dir = filepath.Dir(existing.GetPath())
		nodeName = existing.GetName()
		spaceOwner = existing.GetOwner()

		// For overwrites the existing bytes will be freed on commit, so net required
		// space is uploadLength - existing.Size. Skip for size-deferred uploads.
		if uploadLength >= 0 {
			spaceRef := &provider.Reference{ResourceId: existing.GetId()}
			if _, _, remaining, qErr := c.fs.GetQuota(ctx, spaceRef); qErr == nil {
				existingSize := existing.GetSize()
				netRequired := uint64(uploadLength)
				if existingSize < netRequired {
					netRequired -= existingSize
				} else {
					netRequired = 0
				}
				if remaining < netRequired {
					return nil, errtypes.InsufficientStorage("quota exceeded")
				}
			}
		}
	} else {
		if uploadLength > 0 {
			if _, _, remaining, qErr := c.fs.GetQuota(ctx, ref); qErr == nil && remaining < uint64(uploadLength) {
				return nil, errtypes.InsufficientStorage("quota exceeded")
			}
		}

		result, tfErr := c.fs.TouchFile(ctx, ref, false, mtime)
		if tfErr != nil {
			return nil, tfErr
		}
		nodeID = result.ResourceID.GetOpaqueId()
		spaceID = result.SpaceID
		spaceOwner = result.SpaceOwner
		// Derive dir and name from the ref path (ref must carry a path for new files).
		dir = filepath.Dir(ref.GetPath())
		nodeName = filepath.Base(ref.GetPath())
		parentRef := &provider.Reference{
			ResourceId: ref.ResourceId,
			Path:       filepath.Dir(ref.GetPath()),
		}
		if parentInfo, pErr := c.fs.GetMD(ctx, parentRef, []string{}, []string{}); pErr == nil {
			parentID = parentInfo.GetId().GetOpaqueId()
		}
	}

	if nodeName == "" {
		return nil, errtypes.BadRequest("coordinator: missing filename in ref")
	}
	if dir == "" {
		return nil, errtypes.BadRequest("coordinator: could not determine upload directory")
	}

	session := c.store.New(ctx)
	session.SetMetadata("filename", nodeName)
	session.SetStorageValue("NodeName", nodeName)
	session.SetMetadata("dir", dir)
	session.SetStorageValue("Dir", dir)
	session.SetStorageValue("NodeId", nodeID)
	session.SetStorageValue("SpaceRoot", spaceID)
	if nodeExists {
		session.SetStorageValue("NodeExists", "true")
	}
	session.SetStorageValue("NodeParentId", parentID)
	if spaceOwner != nil {
		session.SetStorageValue("SpaceOwnerOrManager", spaceOwner.GetOpaqueId())
		session.SetStorageValue("SpaceOwnerIdp", spaceOwner.GetIdp())
		session.SetStorageValue("SpaceOwnerType", utils.UserTypeToString(spaceOwner.GetType()))
	}

	usr := ctxpkg.ContextMustGetUser(ctx)
	session.SetExecutant(usr)

	lockID, _ := ctxpkg.ContextGetLockID(ctx)
	session.SetMetadata("lockid", lockID)

	iid, _ := ctxpkg.ContextGetInitiator(ctx)
	session.SetMetadata("initiatorid", iid)

	session.SetSize(uploadLength)

	var mtimeSet bool
	if metadata != nil {
		session.SetMetadata("providerID", metadata["providerID"])
		if v, ok := metadata["mtime"]; ok && v != "null" {
			session.SetMetadata("mtime", v)
			mtimeSet = true
		}
		if v, ok := metadata["expires"]; ok && v != "null" {
			session.SetMetadata("expires", v)
		}
		if _, ok := metadata["sizedeferred"]; ok {
			session.SetSizeIsDeferred(true)
		}
		if checksum, ok := metadata["checksum"]; ok {
			parts := strings.SplitN(checksum, " ", 2)
			if len(parts) != 2 {
				return nil, errtypes.BadRequest("invalid checksum format. must be '[algorithm] [checksum]'")
			}
			switch parts[0] {
			case "sha1", "md5", "adler32":
				session.SetMetadata("checksum", checksum)
			default:
				return nil, errtypes.BadRequest("unsupported checksum algorithm: " + parts[0])
			}
		}
		if v := metadata["if-match"]; v != "" {
			session.SetMetadata("if-match", v)
		}
		if v := metadata["if-none-match"]; v != "" {
			session.SetMetadata("if-none-match", v)
		}
		if v := metadata["if-unmodified-since"]; v != "" {
			session.SetMetadata("if-unmodified-since", v)
		}
	}

	if !mtimeSet {
		session.SetMetadata("mtime", utils.TimeToOCMtime(time.Now()))
	}

	if err := session.TouchBin(); err != nil {
		if !nodeExists {
			_, _ = c.fs.Delete(ctx, ref)
		}
		return nil, fmt.Errorf("coordinator: could not create bin file: %w", err)
	}
	if err := session.Persist(ctx); err != nil {
		session.Cleanup(true, false)
		if !nodeExists {
			_, _ = c.fs.Delete(ctx, ref)
		}
		return nil, fmt.Errorf("coordinator: could not persist session: %w", err)
	}

	sessionRef := session.Reference()
	if err := c.fs.MarkProcessing(ctx, &sessionRef, true, session.ID()); err != nil {
		session.Cleanup(true, true)
		if !nodeExists {
			_, _ = c.fs.Delete(ctx, ref)
		}
		return nil, fmt.Errorf("coordinator: could not mark processing: %w", err)
	}

	metrics.UploadSessionsInitiated.Inc()

	if uploadLength == 0 {
		// Zero-length uploads complete immediately without postprocessing.
		commitRef := session.Reference()
		if _, err := c.fs.CommitUpload(ctx, &commitRef, storage.UploadSource{
			Body:      io.NopCloser(bytes.NewReader(nil)),
			Length:    0,
			Metadata:  session.Metadata(),
			Checksums: session.Checksums(),
		}); err != nil {
			c.rollback(ctx, session)
			return nil, fmt.Errorf("coordinator: zero-length CommitUpload: %w", err)
		}
		_ = c.fs.MarkProcessing(ctx, &commitRef, false, session.ID())
		session.Cleanup(true, true)
		metrics.UploadSessionsFinalized.Inc()
		return map[string]string{
			"simple": session.ID(),
			"tus":    session.ID(),
		}, nil
	}

	return map[string]string{
		"simple": session.ID(),
		"tus":    session.ID(),
	}, nil
}

// Upload handles the simple (single-PUT) upload path so the coordinator owns
// the complete upload lifecycle regardless of the datatx protocol used.
// simple.go calls fs.Upload(); when fs is a *Coordinator this method intercepts.
func (c *coordinator)Upload(ctx context.Context, req storage.UploadRequest, uff storage.UploadFinishedFunc) (*provider.ResourceInfo, error) {
	id := strings.TrimPrefix(req.Ref.GetPath(), "/")
	session, err := c.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = session.Context(ctx)

	size, err := session.WriteChunk(ctx, 0, req.Body)
	if err != nil {
		return nil, err
	}
	if size != req.Length {
		return nil, errtypes.PartialContent(req.Ref.String())
	}

	if err := checksumAndFinish(ctx, session); err != nil {
		c.rollback(ctx, session)
		return nil, err
	}
	if err := session.Persist(ctx); err != nil {
		c.rollback(ctx, session)
		return nil, err
	}

	metrics.UploadProcessing.Inc()
	metrics.UploadSessionsBytesReceived.Inc()

	if !c.async {
		if err := c.commitSync(ctx, session); err != nil {
			return nil, err
		}
	} else {
		s, err := session.URL(ctx)
		if err != nil {
			c.rollback(ctx, session)
			return nil, err
		}
		if err := events.Publish(ctx, c.pub, events.BytesReceived{
			UploadID:   session.ID(),
			URL:        s,
			SpaceOwner: session.SpaceOwner(),
			ExecutingUser: &user.User{
				Id: &user.UserId{
					Type:     session.Executant().Type,
					Idp:      session.Executant().Idp,
					OpaqueId: session.Executant().OpaqueId,
				},
			},
			ResourceID: &provider.ResourceId{
				StorageId: session.ProviderID(),
				SpaceId:   session.SpaceID(),
				OpaqueId:  session.NodeID(),
			},
			Filename:          session.Filename(),
			Filesize:          uint64(session.Size()),
			ImpersonatingUser: impersonatingUser(ctx),
		}); err != nil {
			c.rollback(ctx, session)
			return nil, err
		}
	}

	executant := session.Executant()
	uploadRef := &provider.Reference{
		ResourceId: &provider.ResourceId{
			StorageId: session.ProviderID(),
			SpaceId:   session.SpaceID(),
			OpaqueId:  session.SpaceID(),
		},
		Path: utils.MakeRelativePath(filepath.Join(session.Dir(), session.Filename())),
	}
	if uff != nil {
		uff(session.SpaceOwner(), &executant, uploadRef)
	}

	return &provider.ResourceInfo{
		Id: &provider.ResourceId{
			StorageId: session.ProviderID(),
			SpaceId:   session.SpaceID(),
			OpaqueId:  session.NodeID(),
		},
		Name: session.Filename(),
	}, nil
}

// ListUploadSessions returns upload sessions matching the given filter.
func (c *coordinator)ListUploadSessions(ctx context.Context, filter storage.UploadSessionFilter) ([]storage.UploadSession, error) {
	if filter.ID != nil && *filter.ID != "" {
		session, err := c.store.Get(ctx, *filter.ID)
		if err != nil {
			return nil, err
		}
		return []storage.UploadSession{session}, nil
	}
	sessions, err := c.store.List(ctx)
	if err != nil {
		return nil, err
	}
	result := []storage.UploadSession{}
	now := time.Now()
	for _, s := range sessions {
		if filter.ID != nil && *filter.ID != "" && s.ID() != *filter.ID {
			continue
		}
		if filter.Processing != nil && *filter.Processing != s.IsProcessing() {
			continue
		}
		if filter.Expired != nil {
			if *filter.Expired {
				if now.Before(s.Expires()) {
					continue
				}
			} else {
				if now.After(s.Expires()) {
					continue
				}
			}
		}
		if filter.HasVirus != nil {
			sr, _ := s.ScanData()
			infected := sr != ""
			if *filter.HasVirus != infected {
				continue
			}
		}
		result = append(result, s)
	}
	return result, nil
}

