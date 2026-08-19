// Copyright 2018-2021 CERN
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

package decomposedfs

import (
	"context"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/google/uuid"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/lockedfile"
	tusd "github.com/tus/tusd/v2/pkg/handler"

	"github.com/owncloud/reva/v2/pkg/appctx"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/rhttp/datatx/metrics"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/utils/chunking"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/upload"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"
)

// Upload uploads data to the given resource
// TODO(OCISDEV-901): remove Upload once all drivers are migrated to CommitUpload and the coordinator (OCISDEV-900) is in place.
func (fs *Decomposedfs) Upload(ctx context.Context, req storage.UploadRequest, uff storage.UploadFinishedFunc) (*provider.ResourceInfo, error) {
	_, span := tracer.Start(ctx, "Upload")
	defer span.End()
	up, err := fs.GetUpload(ctx, req.Ref.GetPath())
	if err != nil {
		return &provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error retrieving upload")
	}

	session := up.(*upload.OcisSession)

	ctx = session.Context(ctx)

	if session.Chunk() != "" { // check chunking v1
		p, assembledFile, err := fs.chunkHandler.WriteChunk(session.Chunk(), req.Body)
		if err != nil {
			return &provider.ResourceInfo{}, err
		}
		if p == "" {
			if err = session.Terminate(ctx); err != nil {
				return &provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error removing auxiliary files")
			}
			return &provider.ResourceInfo{}, errtypes.PartialContent(req.Ref.String())
		}
		fd, err := os.Open(assembledFile)
		if err != nil {
			return &provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error opening assembled file")
		}
		defer fd.Close()
		defer os.RemoveAll(assembledFile)
		req.Body = fd

		size, err := session.WriteChunk(ctx, 0, req.Body)
		if err != nil {
			return &provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error writing to binary file")
		}
		session.SetSize(size)
	} else {
		size, err := session.WriteChunk(ctx, 0, req.Body)
		if err != nil {
			return &provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error writing to binary file")
		}
		if size != req.Length {
			return &provider.ResourceInfo{}, errtypes.PartialContent("Decomposedfs: unexpected end of stream")
		}
	}

	if err := session.FinishUploadDecomposed(ctx); err != nil {
		return &provider.ResourceInfo{}, err
	}

	if uff != nil {
		uploadRef := &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: session.ProviderID(),
				SpaceId:   session.SpaceID(),
				OpaqueId:  session.SpaceID(),
			},
			Path: utils.MakeRelativePath(filepath.Join(session.Dir(), session.Filename())),
		}
		executant := session.Executant()
		uff(session.SpaceOwner(), &executant, uploadRef)
	}

	ri := &provider.ResourceInfo{
		// fill with at least fileid, mtime and etag
		Id: &provider.ResourceId{
			StorageId: session.ProviderID(),
			SpaceId:   session.SpaceID(),
			OpaqueId:  session.NodeID(),
		},
	}

	// add etag to metadata
	ri.Etag, _ = node.CalculateEtag(session.NodeID(), session.MTime())

	if !session.MTime().IsZero() {
		ri.Mtime = utils.TimeToTS(session.MTime())
	}

	return ri, nil
}

// InitiateUpload returns upload ids corresponding to different protocols it supports
// TODO(OCISDEV-901): remove InitiateUpload once all drivers are migrated to CommitUpload and the coordinator (OCISDEV-900) is in place.
func (fs *Decomposedfs) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	_, span := tracer.Start(ctx, "InitiateUpload")
	defer span.End()
	log := appctx.GetLogger(ctx)
	log.Debug().Interface("ref", ref).Msg("decomposedfs:InitiateUpload:start")

	// remember the path from the reference
	refpath := ref.GetPath()
	var chunk *chunking.ChunkBLOBInfo
	var err error
	if chunking.IsChunked(refpath) { // check chunking v1
		chunk, err = chunking.GetChunkBLOBInfo(refpath)
		if err != nil {
			return nil, errtypes.BadRequest(err.Error())
		}
		ref.Path = chunk.Path
	}
	n, err := fs.lu.NodeFromResource(ctx, ref)
	switch err.(type) {
	case nil:
		// ok
	case errtypes.IsNotFound:
		return nil, errtypes.PreconditionFailed(err.Error())
	default:
		return nil, err
	}

	// permissions are checked in NewUpload below

	relative, err := fs.lu.Path(ctx, n, node.NoCheck)
	// TODO why do we need the path here?
	// jfd: it is used later when emitting the UploadReady event ...
	// AAAND refPath might be . when accessing with an id / relative reference ... which causes NodeName to become . But then dir will also always be .
	// That is why we still have to read the path here: so that the event we emit contains a relative reference with a path relative to the space root. WTF
	if err != nil {
		return nil, err
	}

	lockID, _ := ctxpkg.ContextGetLockID(ctx)

	session := fs.sessionStore.New(ctx)
	session.SetMetadata("filename", n.Name)
	session.SetStorageValue("NodeName", n.Name)
	if chunk != nil {
		session.SetStorageValue("Chunk", filepath.Base(refpath))
	}
	session.SetMetadata("dir", filepath.Dir(relative))
	session.SetStorageValue("Dir", filepath.Dir(relative))
	session.SetMetadata("lockid", lockID)

	session.SetSize(uploadLength)
	session.SetStorageValue("SpaceRoot", n.SpaceRoot.ID)                                     // TODO SpaceRoot -> SpaceID
	session.SetStorageValue("SpaceOwnerOrManager", n.SpaceOwnerOrManager(ctx).GetOpaqueId()) // TODO needed for what?

	spaceGID, ok := ctx.Value(CtxKeySpaceGID).(uint32)
	if ok {
		session.SetStorageValue("SpaceGid", fmt.Sprintf("%d", spaceGID))
	}

	iid, _ := ctxpkg.ContextGetInitiator(ctx)
	session.SetMetadata("initiatorid", iid)

	if metadata != nil {
		session.SetMetadata("providerID", metadata["providerID"])
		if mtime, ok := metadata["mtime"]; ok {
			if mtime != "null" {
				session.SetMetadata("mtime", metadata["mtime"])
			}
		}
		if expiration, ok := metadata["expires"]; ok {
			if expiration != "null" {
				session.SetMetadata("expires", metadata["expires"])
			}
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

		// only check preconditions if they are not empty // TODO or is this a bad request?
		if metadata["if-match"] != "" {
			session.SetMetadata("if-match", metadata["if-match"])
		}
		if metadata["if-none-match"] != "" {
			session.SetMetadata("if-none-match", metadata["if-none-match"])
		}
		if metadata["if-unmodified-since"] != "" {
			session.SetMetadata("if-unmodified-since", metadata["if-unmodified-since"])
		}
	}

	if session.MTime().IsZero() {
		session.SetMetadata("mtime", utils.TimeToOCMtime(time.Now()))
	}

	log.Debug().Str("uploadid", session.ID()).Str("spaceid", n.SpaceID).Str("nodeid", n.ID).Interface("metadata", metadata).Msg("Decomposedfs: resolved filename")

	_, err = node.CheckQuota(ctx, n.SpaceRoot, n.Exists, uint64(n.Blobsize), uint64(session.Size()))
	if err != nil {
		return nil, err
	}

	if session.Filename() == "" {
		return nil, errors.New("Decomposedfs: missing filename in metadata")
	}
	if session.Dir() == "" {
		return nil, errors.New("Decomposedfs: missing dir in metadata")
	}

	// the parent owner will become the new owner
	parent, perr := n.Parent(ctx)
	if perr != nil {
		return nil, errors.Wrap(perr, "Decomposedfs: error getting parent "+n.ParentID)
	}

	// check permissions
	var (
		checkNode *node.Node
		path      string
	)
	if n.Exists {
		// check permissions of file to be overwritten
		checkNode = n
		path, _ = storagespace.FormatReference(&provider.Reference{ResourceId: &provider.ResourceId{
			SpaceId:  checkNode.SpaceID,
			OpaqueId: checkNode.ID,
		}})
	} else {
		// check permissions of parent
		checkNode = parent
		path, _ = storagespace.FormatReference(&provider.Reference{ResourceId: &provider.ResourceId{
			SpaceId:  checkNode.SpaceID,
			OpaqueId: checkNode.ID,
		}, Path: n.Name})
	}
	rp, err := fs.p.AssemblePermissions(ctx, checkNode)
	switch {
	case err != nil:
		return nil, err
	case !rp.InitiateFileUpload:
		return nil, errtypes.PermissionDenied(path)
	}

	// are we trying to overwriting a folder with a file?
	if n.Exists && n.IsDir(ctx) {
		return nil, errtypes.PreconditionFailed("resource is not a file")
	}

	// check lock
	if err := n.CheckLock(ctx); err != nil {
		return nil, err
	}

	usr := ctxpkg.ContextMustGetUser(ctx)

	// fill future node info
	if n.Exists {
		if session.HeaderIfNoneMatch() == "*" {
			return nil, errtypes.Aborted(fmt.Sprintf("parent %s already has a child %s, id %s", n.ParentID, n.Name, n.ID))
		}
		session.SetStorageValue("NodeId", n.ID)
		session.SetStorageValue("NodeExists", "true")
	} else {
		session.SetStorageValue("NodeId", uuid.New().String())
	}
	session.SetStorageValue("NodeParentId", n.ParentID)
	session.SetExecutant(usr)
	session.SetStorageValue("LogLevel", log.GetLevel().String())

	log.Debug().Interface("session", session).Msg("Decomposedfs: built session info")

	err = fs.um.RunInBaseScope(func() error {
		// Create binary file in the upload folder with no content
		// It will be used when determining the current offset of an upload
		err := session.TouchBin()
		if err != nil {
			return err
		}

		return session.Persist(ctx)
	})
	if err != nil {
		return nil, err
	}
	metrics.UploadSessionsInitiated.Inc()

	if uploadLength == 0 {
		// Directly finish this upload
		err = session.FinishUploadDecomposed(ctx)
		if err != nil {
			return nil, err
		}
	}

	log.Debug().Str("uploadid", session.ID()).Msg("decomposedfs:InitiateUpload:complete")
	return map[string]string{
		"simple": session.ID(),
		"tus":    session.ID(),
	}, nil
}

// MarkProcessing toggles a processing flag on the resource.
func (fs *Decomposedfs) MarkProcessing(ctx context.Context, ref *provider.Reference, processing bool, sessionID string) error {
	n, err := fs.lu.NodeFromResource(ctx, ref)
	if err != nil {
		return err
	}
	if !n.Exists {
		return errtypes.NotFound(ref.String())
	}

	// Early lock, so MarkProcessing is atomic.
	f, err := lockedfile.OpenFile(fs.lu.MetadataBackend().LockfilePath(n.InternalPath()), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			appctx.GetLogger(ctx).Error().Err(cerr).Str("nodeid", n.ID).Msg("could not close mark-processing lock")
		}
	}()

	// Evict the node's in-process xattr cache so IsProcessing reads from disk while we hold the lock.
	n.ResetXattrsCache()

	if !processing {
		if !n.IsProcessing(ctx) {
			return nil
		}
		id, _ := n.ProcessingID(ctx)
		if id != sessionID {
			return nil // owned by a different session, do not clear
		}
		return n.RemoveXattr(ctx, prefixes.StatusPrefix, false)
	}

	return n.SetXattrsWithContext(ctx, node.Attributes{
		prefixes.StatusPrefix: []byte(node.ProcessingStatus + sessionID),
	}, false) // acquireLock=false, because outer lock already held
}

// CommitUpload writes the staged bytes from source to the resource at ref.
// sessionID is used to identify the correct blob slot prepared before postprocessing.
// Caller owns source.Body and must close it after CommitUpload returns.
func (fs *Decomposedfs) CommitUpload(ctx context.Context, ref *provider.Reference, sessionID string, source storage.UploadSource) error {
	if source.Body == nil {
		return errtypes.BadRequest("Decomposedfs: source body is nil")
	}
	if sessionID == "" {
		return errtypes.BadRequest("Decomposedfs: sessionID is empty")
	}

	n, err := fs.lu.NodeFromResource(ctx, ref)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("CommitUpload: unexpected NodeFromResource failure")
		return errtypes.InternalError("CommitUpload: node lookup failed unexpectedly")
	}
	if !n.Exists {
		return errtypes.NotFound(ref.String())
	}

	blobNode := node.New(n.SpaceID, n.ID, "", "", source.Length, sessionID,
		provider.ResourceType_RESOURCE_TYPE_FILE, nil, fs.lu)

	if err := fs.tp.WriteBlobFromReader(blobNode, source.Body, source.Length); err != nil {
		if derr := fs.tp.DeleteBlob(blobNode); derr != nil {
			appctx.GetLogger(ctx).Error().Err(derr).Str("nodeid", n.ID).Str("blobid", sessionID).Msg("could not clean up orphaned blob after failed write")
		}
		return errors.Wrap(err, "Decomposedfs: failed to write blob")
	}

	// on the node, not the session: the session is deleted when the upload finishes
	if !source.ScanDate.IsZero() {
		if err := n.SetScanData(ctx, source.ScanResult, source.ScanDate); err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("nodeid", n.ID).Msg("could not set scan results")
		}
	}

	now := time.Now()
	if p, err := n.Parent(ctx); err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Str("nodeid", n.ID).Msg("could not read parent for etag propagation")
	} else {
		_ = p.SetTMTime(ctx, &now)
		if err := fs.tp.Propagate(ctx, p, 0); err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("nodeid", n.ID).Msg("could not propagate etag change")
		}
	}

	return nil
}

// PrepareUpload finalizes node metadata after bytes are received, before postprocessing.
// CommitUpload is called after postprocessing completes.
func (fs *Decomposedfs) PrepareUpload(ctx context.Context, ref *provider.Reference, sessionID string, info storage.UploadInfo) (*storage.PrepareUploadResult, error) {
	ctx, span := tracer.Start(ctx, "PrepareUpload")
	defer span.End()

	n, err := fs.lu.NodeFromResource(ctx, ref)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("PrepareUpload: unexpected NodeFromResource failure")
		return nil, errtypes.InternalError("PrepareUpload: node lookup failed unexpectedly")
	}
	if !n.Exists {
		return nil, errtypes.NotFound(ref.String())
	}
	n.SpaceRoot, err = node.ReadNode(ctx, fs.lu, n.SpaceID, n.SpaceID, false, nil, false)
	if err != nil {
		return nil, err
	}

	if err := n.CheckLock(ctx); err != nil {
		return nil, err
	}

	// scope to space owner GID for posix deployments; no-op with NullMapper
	if spaceGID, ok := ctx.Value(CtxKeySpaceGID).(uint32); ok {
		unscope, err := fs.um.ScopeUserByIds(-1, int(spaceGID))
		if err != nil {
			return nil, errors.Wrap(err, "failed to scope user")
		}
		if unscope != nil {
			defer func() { _ = unscope() }()
		}
	}

	targetPath := n.InternalPath()
	f, err := lockedfile.OpenFile(fs.lu.MetadataBackend().LockfilePath(targetPath), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	unlock := func() error { return f.Close() }
	defer func() {
		if err := unlock(); err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("nodeid", n.ID).Msg("could not close lock")
		}
	}()

	var (
		sizeDiff       int64
		versionCreated bool
		versionPath    string
		oldAttrs       node.Attributes
		oldMtime       time.Time
		committed      bool
	)

	defer func() {
		if committed {
			return
		}
		if versionCreated {
			if err := os.Remove(versionPath); err != nil && !errors.Is(err, os.ErrNotExist) {
				appctx.GetLogger(ctx).Error().Err(err).Str("versionpath", versionPath).Msg("could not remove version file during rollback")
			}
		}
		if info.NodeExisted && oldAttrs != nil {
			// mtime goes in the same batch: we still hold the metadata lock and
			// SetMTime would deadlock trying to retake it
			if err := fs.lu.TimeManager().OverrideMtime(ctx, n, &oldAttrs, oldMtime); err != nil {
				appctx.GetLogger(ctx).Error().Err(err).Str("nodeid", n.ID).Msg("could not restore node mtime during rollback")
			}
			if err := n.SetXattrsWithContext(ctx, oldAttrs, false); err != nil {
				appctx.GetLogger(ctx).Error().Err(err).Str("nodeid", n.ID).Msg("could not restore node xattrs during rollback")
			}
		}
	}()

	var (
		old         *node.Node
		overwrite   bool
		oldBlobsize int64
	)
	if info.NodeExisted {
		old, err = node.ReadNode(ctx, fs.lu, n.SpaceID, n.ID, false, nil, false)
		if err != nil {
			return nil, errors.Wrap(err, "PrepareUpload: failed to read existing node")
		}
		overwrite = old.BlobID != ""
		oldBlobsize = old.Blobsize
	}

	// also for new files: the coordinator's check fails open without GetQuota
	// permission, and CheckQuota guards disk space too
	if _, err := node.CheckQuota(ctx, n.SpaceRoot, overwrite, uint64(oldBlobsize), uint64(info.Size)); err != nil {
		return nil, err
	}

	if info.NodeExisted {
		oldMtime, err = old.GetMTime(ctx)
		if err != nil {
			return nil, err
		}
		oldEtag, err := node.CalculateEtag(old.ID, oldMtime)
		if err != nil {
			return nil, err
		}

		if info.IfMatch != "" && info.IfMatch != oldEtag {
			return nil, errtypes.Aborted("etag mismatch")
		}
		if info.IfNoneMatch != "" {
			if info.IfNoneMatch == "*" {
				return nil, errtypes.Aborted("etag mismatch, resource exists")
			}
			for _, tag := range strings.Split(info.IfNoneMatch, ",") {
				if tag == oldEtag {
					return nil, errtypes.Aborted("etag mismatch")
				}
			}
		}
		if !info.IfUnmodifiedSince.IsZero() && oldMtime.After(info.IfUnmodifiedSince) {
			return nil, errtypes.Aborted("if-unmodified-since mismatch")
		}

		// capture full node xattrs for rollback before any write
		oldAttrs, err = fs.lu.MetadataBackend().All(ctx, targetPath)
		if err != nil {
			return nil, err
		}

		if !fs.o.DisableVersioning {
			versionPath = fs.lu.InternalPath(n.SpaceID, n.ID+node.RevisionIDDelimiter+oldMtime.UTC().Format(time.RFC3339Nano))
			revFile, err := os.OpenFile(versionPath, os.O_CREATE|os.O_EXCL, 0600)
			if err != nil {
				if !errors.Is(err, os.ErrExist) {
					return nil, err
				}
				// revision with this mtime already exists; verify blobs match then reclaim the slot
				if err := validateChecksums(ctx, fs.lu, n, versionPath); err != nil {
					return nil, err
				}
				blobID, _, err := fs.lu.ReadBlobIDAndSizeAttr(ctx, versionPath, nil)
				if err != nil {
					return nil, err
				}
				if err := fs.tp.DeleteBlob(&node.Node{BlobID: blobID, SpaceID: n.SpaceID}); err != nil {
					return nil, err
				}
				f2, err := os.Create(versionPath)
				if err != nil {
					return nil, err
				}
				f2.Close()
			} else {
				revFile.Close()
			}

			if err := fs.lu.CopyMetadataWithSourceLock(ctx, targetPath, versionPath, func(attributeName string, value []byte) ([]byte, bool) {
				return value, strings.HasPrefix(attributeName, prefixes.ChecksumPrefix) ||
					attributeName == prefixes.TypeAttr ||
					attributeName == prefixes.BlobIDAttr ||
					attributeName == prefixes.BlobsizeAttr ||
					attributeName == prefixes.MTimeAttr
			}, f, true); err != nil {
				return nil, err
			}
			if err := os.Chtimes(versionPath, oldMtime, oldMtime); err != nil {
				return nil, errtypes.InternalError(fmt.Sprintf("failed to set mtime on version node: %s", err))
			}
			versionCreated = true
		}

		sizeDiff = info.Size - old.Blobsize
	} else {
		if c, ok := fs.lu.(node.IDCacher); ok {
			if err := c.CacheID(ctx, n.SpaceID, n.ID, filepath.Join(n.ParentPath(), n.Name)); err != nil {
				appctx.GetLogger(ctx).Error().Err(err).Msg("failed to cache id")
			}
		}
		sizeDiff = info.Size
	}

	attrs := node.Attributes{}
	attrs.SetString(prefixes.IDAttr, n.ID)
	attrs.SetInt64(prefixes.TypeAttr, int64(provider.ResourceType_RESOURCE_TYPE_FILE))
	attrs.SetString(prefixes.ParentidAttr, n.ParentID)
	attrs.SetString(prefixes.NameAttr, n.Name)
	attrs.SetString(prefixes.BlobIDAttr, sessionID)
	attrs.SetInt64(prefixes.BlobsizeAttr, info.Size)
	attrs[prefixes.ChecksumPrefix+"sha1"] = info.Checksums.SHA1
	attrs[prefixes.ChecksumPrefix+"md5"] = info.Checksums.MD5
	attrs[prefixes.ChecksumPrefix+"adler32"] = info.Checksums.Adler32

	mtime := time.Now()
	if !info.MTime.IsZero() {
		mtime = info.MTime
	}
	if err := fs.lu.TimeManager().OverrideMtime(ctx, n, &attrs, mtime); err != nil {
		return nil, errors.Wrap(err, "failed to set mtime")
	}

	if err := n.SetXattrsWithContext(ctx, attrs, false); err != nil {
		return nil, errors.Wrap(err, "could not write metadata")
	}

	if err := fs.tp.Propagate(ctx, n, sizeDiff); err != nil {
		return nil, errors.Wrap(err, "could not propagate size change")
	}
	committed = true

	return &storage.PrepareUploadResult{VersionCreated: versionCreated, SizeDiff: sizeDiff}, nil
}

// RollbackUpload reverts the node state written by PrepareUpload after a failed or aborted
// postprocessing run. It restores the previous revision (or purges the node if versioning is
// disabled and no prior version exists) and reverts the optimistic size propagation.
func (fs *Decomposedfs) RollbackUpload(ctx context.Context, ref *provider.Reference, sessionID string, info storage.RollbackInfo) error {
	n, err := fs.lu.NodeFromResource(ctx, ref)
	if err != nil {
		// The node metadata is unreadable, so the upload can never finish and its
		// quota would stay consumed forever.
		return fs.rollbackOrphaned(ctx, ref, sessionID, info, err)
	}
	if !n.Exists {
		return nil // nothing was written yet
	}
	n.SpaceRoot, err = node.ReadNode(ctx, fs.lu, n.SpaceID, n.SpaceID, false, nil, false)
	if err != nil {
		return err
	}

	curProcessingID, err := n.ProcessingID(ctx)
	if err != nil {
		return fmt.Errorf("RollbackUpload: could not read processing ID: %w", err)
	}
	if curProcessingID != sessionID {
		return nil
	}

	if info.NodeExisted {
		if err := n.RevertCurrentRevision(ctx, false); err != nil {
			return err
		}
	} else {
		// the upload created this node, so undo it. Purge, not trash: the file never
		// became visible, and Purge needs no Delete permission
		if err := n.Purge(ctx); err != nil {
			return fmt.Errorf("RollbackUpload: could not purge node: %w", err)
		}
	}

	if info.SizeDiff != 0 {
		if err := fs.tp.Propagate(ctx, n, -info.SizeDiff); err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Msg("RollbackUpload: could not revert propagate")
		}
	}
	return nil
}

// rollbackOrphaned rolls back an upload whose target node can no longer be read,
// which happens when the node file outlives its metadata: an ancestor trashed
// mid-upload leaves a node whose ReadNode fails on the missing parent id. The ids
// recorded on the session are enough to walk up the tree, so the quota is released
// and the unreachable node removed rather than left consuming space forever.
//
// lookupErr is what made the node unreadable, returned when there is nothing on
// the session to fall back to.
func (fs *Decomposedfs) rollbackOrphaned(ctx context.Context, ref *provider.Reference, sessionID string, info storage.RollbackInfo, lookupErr error) error {
	if info.NodeID == "" || info.ParentID == "" {
		return fmt.Errorf("RollbackUpload: node lookup failed: %w", lookupErr)
	}
	spaceID := ref.GetResourceId().GetSpaceId()
	if spaceID == "" {
		return fmt.Errorf("RollbackUpload: node lookup failed: %w", lookupErr)
	}
	log := appctx.GetLogger(ctx)
	log.Info().Err(lookupErr).Str("sessionid", sessionID).Str("nodeid", info.NodeID).
		Msg("node unreadable, rolling back orphaned upload")

	n := node.New(spaceID, info.NodeID, info.ParentID, info.Filename, info.Size, sessionID,
		provider.ResourceType_RESOURCE_TYPE_FILE, nil, fs.lu)
	spaceRoot, err := node.ReadNode(ctx, fs.lu, spaceID, spaceID, false, nil, false)
	if err != nil {
		return fmt.Errorf("RollbackUpload: space root lookup failed: %w", err)
	}
	n.SpaceRoot = spaceRoot

	if info.SizeDiff != 0 {
		// Stop before removing anything: the caller keeps the session on an error,
		// so this stays retryable instead of leaking the quota silently.
		if err := fs.tp.Propagate(ctx, n, -info.SizeDiff); err != nil {
			return fmt.Errorf("RollbackUpload: could not revert propagate: %w", err)
		}
	}

	// Nothing can resolve this node any more, so remove it and its metadata.
	// Already gone is fine: the node may never have been created.
	nodePath := n.InternalPath()
	if err := utils.RemoveItem(nodePath); err != nil && !errors.Is(err, iofs.ErrNotExist) {
		log.Error().Err(err).Str("nodepath", nodePath).Msg("RollbackUpload: removing orphaned node failed")
	}
	if err := fs.lu.MetadataBackend().Purge(ctx, nodePath); err != nil && !errors.Is(err, iofs.ErrNotExist) {
		log.Error().Err(err).Str("nodepath", nodePath).Msg("RollbackUpload: purging orphaned node metadata failed")
	}
	return nil
}

func validateChecksums(ctx context.Context, lu node.PathLookup, n *node.Node, versionPath string) error {
	for _, t := range []string{"md5", "sha1", "adler32"} {
		key := prefixes.ChecksumPrefix + t
		checksum, err := n.Xattr(ctx, key)
		if err != nil {
			return err
		}
		revisionChecksum, err := lu.MetadataBackend().Get(ctx, versionPath, key)
		if err != nil {
			return err
		}
		if string(checksum) == "" || string(revisionChecksum) == "" {
			return errors.New("checksum not found")
		}
		if string(checksum) != string(revisionChecksum) {
			return errors.New("checksum mismatch")
		}
	}
	return nil
}

// UseIn tells the tus upload middleware which extensions it supports.
func (fs *Decomposedfs) UseIn(composer *tusd.StoreComposer) {
	composer.UseCore(fs)
	composer.UseTerminater(fs)
	composer.UseConcater(fs)
	composer.UseLengthDeferrer(fs)
}

// To implement the core tus.io protocol as specified in https://tus.io/protocols/resumable-upload.html#core-protocol
// - the storage needs to implement NewUpload and GetUpload
// - the upload needs to implement the tusd.Upload interface: WriteChunk, GetInfo, GetReader and FinishUpload

// NewUpload returns a new tus Upload instance
func (fs *Decomposedfs) NewUpload(ctx context.Context, info tusd.FileInfo) (tusd.Upload, error) {
	return nil, fmt.Errorf("not implemented, use InitiateUpload on the CS3 API to start a new upload")
}

// GetUpload returns the Upload for the given upload id
func (fs *Decomposedfs) GetUpload(ctx context.Context, id string) (tusd.Upload, error) {
	var ul tusd.Upload
	var err error
	_ = fs.um.RunInBaseScope(func() error {
		ul, err = fs.sessionStore.Get(ctx, id)
		return nil
	})
	return ul, err
}

// ListUploadSessions returns the upload sessions for the given filter
func (fs *Decomposedfs) ListUploadSessions(ctx context.Context, filter storage.UploadSessionFilter) ([]storage.UploadSession, error) {
	var sessions []*upload.OcisSession
	if filter.ID != nil && *filter.ID != "" {
		session, err := fs.sessionStore.Get(ctx, *filter.ID)
		if err != nil {
			return nil, err
		}
		sessions = []*upload.OcisSession{session}
	} else {
		var err error
		sessions, err = fs.sessionStore.List(ctx)
		if err != nil {
			return nil, err
		}
	}
	filteredSessions := []storage.UploadSession{}
	now := time.Now()
	for _, session := range sessions {
		if filter.Processing != nil && *filter.Processing != session.IsProcessing() {
			continue
		}
		if filter.Expired != nil {
			if *filter.Expired {
				if now.Before(session.Expires()) {
					continue
				}
			} else {
				if now.After(session.Expires()) {
					continue
				}
			}
		}
		if filter.HasVirus != nil {
			sr, _ := session.ScanData()
			infected := sr != ""
			if *filter.HasVirus != infected {
				continue
			}
		}
		// evaluated last: unlike the other filters this reads the node metadata
		// from disk, so it is only done for sessions that passed all other filters
		if filter.Orphaned != nil && *filter.Orphaned != session.IsOrphaned(ctx) {
			continue
		}
		filteredSessions = append(filteredSessions, session)
	}
	return filteredSessions, nil
}

// AsTerminatableUpload returns a TerminatableUpload
// To implement the termination extension as specified in https://tus.io/protocols/resumable-upload.html#termination
// the storage needs to implement AsTerminatableUpload
func (fs *Decomposedfs) AsTerminatableUpload(up tusd.Upload) tusd.TerminatableUpload {
	return up.(*upload.OcisSession)
}

// AsLengthDeclarableUpload returns a LengthDeclarableUpload
// To implement the creation-defer-length extension as specified in https://tus.io/protocols/resumable-upload.html#creation
// the storage needs to implement AsLengthDeclarableUpload
func (fs *Decomposedfs) AsLengthDeclarableUpload(up tusd.Upload) tusd.LengthDeclarableUpload {
	return up.(*upload.OcisSession)
}

// AsConcatableUpload returns a ConcatableUpload
// To implement the concatenation extension as specified in https://tus.io/protocols/resumable-upload.html#concatenation
// the storage needs to implement AsConcatableUpload
func (fs *Decomposedfs) AsConcatableUpload(up tusd.Upload) tusd.ConcatableUpload {
	return up.(*upload.OcisSession)
}
