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
	"context"
	"fmt"
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

	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rhttp/datatx/metrics"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/utils/chunking"
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

// rewriteChunkedRef strips the chunk suffix from ref.Path and returns the chunk basename.
// Only called when chunking.IsChunked(ref.GetPath()) is true.
func rewriteChunkedRef(ref *provider.Reference) (*provider.Reference, string, error) {
	ci, err := chunking.GetChunkBLOBInfo(ref.GetPath())
	if err != nil {
		return nil, "", errtypes.BadRequest(err.Error())
	}
	return &provider.Reference{ResourceId: ref.ResourceId, Path: ci.Path}, filepath.Base(ref.GetPath()), nil
}

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

// finishSync commits the upload inline, without postprocessing.
// Used when async=false or the upload is empty (size==0).
func (c *coordinator) finishSync(ctx context.Context, session Session) error {
	ref := session.Reference()
	f, err := os.Open(session.BinPath())
	if err != nil {
		c.rollback(ctx, session)
		return err
	}
	cs := session.Checksums()
	if _, err := c.fs.CommitUpload(ctx, &ref, storage.UploadSource{
		Body:      f,
		Length:    session.Size(),
		Metadata:  session.Metadata(),
		Checksums: cs,
	}); err != nil {
		c.rollback(ctx, session)
		return err
	}
	_ = c.fs.MarkProcessing(ctx, &ref, false, session.ID())
	session.Cleanup(true, true)
	metrics.UploadSessionsFinalized.Inc()
	return nil
}

// triggerPostprocessing publishes BytesReceived to start async postprocessing.
func (c *coordinator) triggerPostprocessing(ctx context.Context, session Session) error {
	s, err := session.URL(ctx)
	if err != nil {
		c.rollback(ctx, session)
		return err
	}
	executingUser, _ := ctxpkg.ContextGetUser(ctx)
	if err := events.Publish(ctx, c.pub, events.BytesReceived{
		UploadID:      session.ID(),
		URL:           s,
		SpaceOwner:    session.SpaceOwner(),
		ExecutingUser: executingUser,
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
		return err
	}
	return nil
}

// finishUpload is called after all bytes are received (TUS FinishUpload and simple PUT).
// It creates the node, validates checksums, then either commits inline or triggers postprocessing.
func (c *coordinator) finishUpload(ctx context.Context, session Session) error {
	if err := c.touchAndMark(ctx, session); err != nil {
		return err
	}
	if err := verifyAndStoreChecksums(ctx, session); err != nil {
		c.rollback(ctx, session)
		return err
	}
	if err := session.Persist(ctx); err != nil {
		c.rollback(ctx, session)
		return err
	}

	metrics.UploadProcessing.Inc()
	metrics.UploadSessionsBytesReceived.Inc()

	if !c.async || session.Size() == 0 {
		return c.finishSync(ctx, session)
	}
	return c.triggerPostprocessing(ctx, session)
}

// verifyAndStoreChecksums computes checksums from the staged binary, validates against
// any client-supplied checksum, and stores them on the session.
func verifyAndStoreChecksums(ctx context.Context, session Session) error {
	sha1h, md5h, adler32h, err := calculateChecksums(ctx, session.BinPath())
	if err != nil {
		return err
	}
	info, err := session.GetInfo(ctx)
	if err != nil {
		return err
	}
	if checksum := info.MetaData["checksum"]; checksum != "" {
		parts := strings.SplitN(checksum, " ", 2)
		if len(parts) != 2 {
			return errtypes.BadRequest("invalid checksum format. must be '[algorithm] [checksum]'")
		}
		var checkErr error
		switch parts[0] {
		case "sha1":
			checkErr = checkHash(parts[1], sha1h)
		case "md5":
			checkErr = checkHash(parts[1], md5h)
		case "adler32":
			checkErr = checkHash(parts[1], adler32h)
		default:
			checkErr = errtypes.BadRequest("unsupported checksum algorithm: " + parts[0])
		}
		if checkErr != nil {
			session.Cleanup(true, true)
			return checkErr
		}
	}
	session.SetChecksums(sha1h.Sum(nil), md5h.Sum(nil), adler32h.Sum(nil))
	return nil
}

// touchAndMark creates the node (new files only) and marks it as processing.
// Called from FinishUpload and Upload after all bytes have been received.
func (c *coordinator) touchAndMark(ctx context.Context, session Session) error {
	if !session.NodeExists() {
		pathRef := &provider.Reference{
			ResourceId: &provider.ResourceId{
				SpaceId:  session.SpaceID(),
				OpaqueId: session.NodeParentID(),
			},
			Path: session.Filename(),
		}
		result, err := c.fs.TouchFile(ctx, pathRef, false, session.Metadata()["mtime"])
		if err != nil {
			session.Cleanup(true, true)
			if _, ok := err.(errtypes.IsNotFound); ok {
				return errtypes.PreconditionFailed(err.Error())
			}
			return err
		}
		session.SetStorageValue("NodeId", result.ResourceID.GetOpaqueId())
		session.SetStorageValue("SpaceRoot", result.SpaceID)
		if result.SpaceOwner != nil {
			session.SetStorageValue("SpaceOwnerOrManager", result.SpaceOwner.GetOpaqueId())
			session.SetStorageValue("SpaceOwnerIdp", result.SpaceOwner.GetIdp())
			session.SetStorageValue("SpaceOwnerType", utils.UserTypeToString(result.SpaceOwner.GetType()))
		}
	}
	nodeRef := session.Reference()
	if err := c.fs.MarkProcessing(ctx, &nodeRef, true, session.ID()); err != nil {
		session.Cleanup(true, true)
		if !session.NodeExists() {
			_, _ = c.fs.Delete(ctx, &nodeRef)
		}
		return err
	}
	return session.Persist(ctx)
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
	fs           storage.FS
	store        SessionStore
	pub          events.Publisher
	async        bool
	mountID      string
	numConc      int
	conGroup     string
	log          *zerolog.Logger
	chunkHandler *chunking.ChunkHandler // nil when legacy chunking v1 is not needed
}

// NewCoordinator constructs a Coordinator. Call Start to begin consuming events.
// async=true requires a non-nil pub.
// chunkFolder enables legacy chunking v1 support; pass "" to disable it.
func NewCoordinator(
	fs storage.FS,
	store SessionStore,
	pub events.Publisher,
	async bool,
	mountID string,
	consumerGroup string,
	numConsumers int,
	log *zerolog.Logger,
	chunkFolder string,
) (Coordinator, error) {
	if async && pub == nil {
		return nil, fmt.Errorf("need event stream for async upload processing")
	}
	if numConsumers <= 0 {
		numConsumers = 1
	}
	var ch *chunking.ChunkHandler
	if chunkFolder != "" {
		ch = chunking.NewChunkHandler(chunkFolder)
	}
	return &coordinator{
		fs:           fs,
		store:        store,
		pub:          pub,
		async:        async,
		mountID:      mountID,
		numConc:      numConsumers,
		conGroup:     consumerGroup,
		log:          log,
		chunkHandler: ch,
	}, nil
}

func (c *coordinator) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	var chunkName string
	if chunking.IsChunked(ref.GetPath()) { // check legacy chunking v1
		var rerr error
		ref, chunkName, rerr = rewriteChunkedRef(ref)
		if rerr != nil {
			return nil, rerr
		}
	}

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

	var nodeID, spaceID, parentID, dir, nodeName string
	var spaceOwner *user.UserId

	// check quota
	if uploadLength >= 0 {
		spaceRef := &provider.Reference{ResourceId: &provider.ResourceId{
			StorageId: ref.GetResourceId().GetStorageId(),
			SpaceId:   ref.GetResourceId().GetSpaceId(),
		}}
		if _, _, remaining, qErr := c.fs.GetQuota(ctx, spaceRef); qErr == nil {
			var existingSize uint64
			if nodeExists {
				existingSize = existing.GetSize()
			}
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

	if nodeExists {
		nodeID = existing.GetId().GetOpaqueId()
		spaceID = existing.GetId().GetSpaceId()
		parentID = existing.GetParentId().GetOpaqueId()
		dir = filepath.Dir(ref.GetPath())
		nodeName = existing.GetName()
		spaceOwner = existing.GetOwner()

		diskLock, _ := c.fs.GetLock(ctx, ref)
		contextLockID, _ := ctxpkg.ContextGetLockID(ctx)
		if diskLock != nil {
			switch contextLockID {
			case "":
				return nil, errtypes.Locked(diskLock.LockId)
			case diskLock.LockId:
				// ok
			default:
				return nil, errtypes.Aborted("mismatching lock")
			}
		} else if contextLockID != "" {
			return nil, errtypes.Aborted("not locked")
		}
	} else {
		spaceID = ref.GetResourceId().GetSpaceId()
		dir = filepath.Dir(ref.GetPath())
		nodeName = filepath.Base(ref.GetPath())
	}

	if nodeExists {
		if !existing.GetPermissionSet().GetInitiateFileUpload() {
			return nil, errtypes.PermissionDenied(ref.GetPath())
		}
		if existing.GetType() == provider.ResourceType_RESOURCE_TYPE_CONTAINER {
			return nil, errtypes.PreconditionFailed("resource is not a file")
		}
		if metadata["if-none-match"] == "*" {
			return nil, errtypes.Aborted(fmt.Sprintf("parent %s already has a child %s, id %s", parentID, nodeName, nodeID))
		}
	} else {
		parentRef := &provider.Reference{
			ResourceId: ref.GetResourceId(),
			Path:       dir,
		}
		parentMD, pErr := c.fs.GetMD(ctx, parentRef, []string{}, []string{})
		switch pErr.(type) {
		case nil:
		case errtypes.IsNotFound:
			// RFC 4918: missing intermediate dir → 409, no permission → 404.
			// GetMD returns NotFound for both (hides resources from unauthorized callers).
			// Walk up the path: if an ancestor is visible, the dir is truly missing (409).
			// If nothing is visible up to the root, caller has no access (404).
			ancestor := dir
			permDenied := true
			for ancestor != "." && ancestor != "/" {
				ancestor = filepath.Dir(ancestor)
				ancestorRef := &provider.Reference{ResourceId: ref.GetResourceId(), Path: ancestor}
				if _, aErr := c.fs.GetMD(ctx, ancestorRef, []string{}, []string{}); aErr == nil {
					permDenied = false
					break
				}
			}
			if permDenied {
				return nil, errtypes.PermissionDenied(ref.GetPath())
			}
			return nil, errtypes.PreconditionFailed(pErr.Error())
		default:
			return nil, pErr
		}
		if !parentMD.GetPermissionSet().GetInitiateFileUpload() {
			return nil, errtypes.PermissionDenied(ref.GetPath())
		}
		parentID = parentMD.GetId().GetOpaqueId()
		spaceID = parentMD.GetId().GetSpaceId()
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
	session.SetStorageValue("SpaceRoot", spaceID)
	if nodeExists {
		session.SetStorageValue("NodeId", nodeID)
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
	if chunkName != "" { // check legacy chunking v1
		session.SetStorageValue("Chunk", chunkName)
	}

	if err := session.TouchBin(); err != nil {
		return nil, fmt.Errorf("coordinator: could not create bin file: %w", err)
	}
	if err := session.Persist(ctx); err != nil {
		session.Cleanup(true, true)
		return nil, fmt.Errorf("coordinator: could not persist session: %w", err)
	}

	metrics.UploadSessionsInitiated.Inc()

	if uploadLength == 0 {
		// Zero-length uploads complete immediately without postprocessing.
		if err := c.finishUpload(ctx, session); err != nil {
			return nil, err
		}
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
func (c *coordinator) Upload(ctx context.Context, req storage.UploadRequest, uff storage.UploadFinishedFunc) (*provider.ResourceInfo, error) {
	id := strings.TrimPrefix(req.Ref.GetPath(), "/")
	session, err := c.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	ctx = session.Context(ctx)

	if session.Chunk() != "" { // check legacy chunking v1
		assembled, assembledSize, done, err := c.chunkHandler.Assemble(session.Chunk(), req.Body)
		if err != nil {
			return nil, err
		}
		if !done {
			session.Cleanup(true, true)
			return nil, errtypes.PartialContent(req.Ref.String())
		}
		defer assembled.Close()
		req.Body, req.Length = assembled, assembledSize
		session.SetSize(assembledSize)
	}

	size, err := session.WriteChunk(ctx, 0, req.Body)
	if err != nil {
		return nil, err
	}
	if size != req.Length {
		return nil, errtypes.PartialContent(req.Ref.String())
	}

	if err := c.finishUpload(ctx, session); err != nil {
		return nil, err
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

	ri := &provider.ResourceInfo{
		Id: &provider.ResourceId{
			StorageId: session.ProviderID(),
			SpaceId:   session.SpaceID(),
			OpaqueId:  session.NodeID(),
		},
		Name: session.Filename(),
	}
	if mt, ok := session.Metadata()["mtime"]; ok && mt != "" {
		if t, err := utils.MTimeToTime(mt); err == nil {
			ri.Etag, _ = utils.CalculateEtag(session.NodeID(), t)
			ri.Mtime = utils.TimeToTS(t)
		}
	}
	return ri, nil
}

// ListUploadSessions returns upload sessions matching the given filter.
func (c *coordinator) ListUploadSessions(ctx context.Context, filter storage.UploadSessionFilter) ([]storage.UploadSession, error) {
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
