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

// Package upload provides the driver-agnostic upload coordinator.
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
	"github.com/google/uuid"

	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/rhttp/datatx/metrics"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/utils/chunking"
	"github.com/owncloud/reva/v2/pkg/utils"
)

// Coordinator owns the upload lifecycle.
type Coordinator interface {
	// InitiateUpload returns a list of protocols with urls that can be used to append bytes to a new upload session.
	InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error)
}

// coordinator is the concrete implementation of Coordinator.
type coordinator struct {
	fs    storage.FS
	store SessionStore
}

// NewCoordinator constructs a coordinator backed by the given storage driver
// and session store. The store must use an on-disk session format the driver's
// data path can read (the decomposedfs family: ocis/s3ng/posix).
func NewCoordinator(fs storage.FS, store SessionStore) *coordinator {
	return &coordinator{fs: fs, store: store}
}

// InitiateUpload returns a list of protocols with urls that can be used to append bytes to a new upload session.
func (c *coordinator) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	return c.initiateUpload(ctx, ref, uploadLength, metadata)
}

// initiateUpload is the driver-agnostic port of decomposedfs.InitiateUpload.
//
// Known open divergences from main, tracked as findings and NOT yet resolved:
//   - B2: permission-gated GetMD hides deny-granted files → late 409 instead of 403.
//   - B6/B7: spaceOwner manager-fallback and posix scoping (SpaceGid, RunInBaseScope).
func (c *coordinator) initiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	var chunkName string
	if chunking.IsChunked(ref.GetPath()) { // check legacy chunking v1
		var rerr error
		ref, chunkName, rerr = rewriteChunkedRef(ref)
		if rerr != nil {
			return nil, rerr
		}
	}

	// nodeExists=false is overloaded: genuinely absent, or exists-but-hidden by a deny-grant.
	//
	// TODO(OCISDEV-900): permission-gated GetMD hides a deny-granted
	// file as NotFound, so we take the new-file branch (200) and fail late at
	// finish with 409 instead of main's up-front 403 — an existence oracle plus a
	// wasted upload. Accepted for now; a clean fix needs a permission-free resolve.
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
		// GetQuota is permission-gated: roles that can upload but lack GetQuota (Uploader, share) error here, so we fail open and let finish enforce.
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
		// GetMD returns only the basename for relative (id-based) refs, so
		// filepath.Dir would yield "." here. Reconstruct the space-relative
		// path via the public FS interface — mirrors main's fs.lu.Path.
		// Best-effort: on error keep the basename rather than failing an
		// upload main would allow.
		relPath := existing.GetPath()
		if utils.IsRelativeReference(ref) {
			if full, pErr := c.fs.GetPathByID(ctx, existing.GetId()); pErr == nil {
				relPath = full
			}
		}
		dir = filepath.Dir(relPath)
		nodeName = existing.GetName()
		// TODO(OCISDEV-900, finding B6): main uses SpaceOwnerOrManager (falls back to a
		// manager when owner is nil/SPACE_OWNER, e.g. project drives). GetOwner() has no
		// such fallback, and the new-file branch never sets spaceOwner at all.
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
			// GetMD collapses "dir missing" and "dir hidden (no access)" both into NotFound.
			// Walk up: if any ancestor is visible, the dir is genuinely missing → PreconditionFailed.
			// If nothing is visible up to the root, the caller has no access → PermissionDenied.
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
	} else {
		//todo not sure if this is correct
		// mint the future node id for the new file (main: upload.go:308)
		session.SetStorageValue("NodeId", uuid.New().String())
	}
	session.SetStorageValue("NodeParentId", parentID)
	if spaceOwner != nil {
		session.SetStorageValue("SpaceOwnerOrManager", spaceOwner.GetOpaqueId())
		session.SetStorageValue("SpaceOwnerIdp", spaceOwner.GetIdp())
		session.SetStorageValue("SpaceOwnerType", utils.UserTypeToString(spaceOwner.GetType()))
	}

	// TODO(OCISDEV-900, finding B7): main copies CtxKeySpaceGID into the session
	// (upload.go:188) to drive posix uid/gid scoping at commit. That key lives in the
	// decomposedfs package; reading it here would make the driver-agnostic coordinator
	// depend on a concrete driver. posix-only concern (unset on ocis/s3ng). Deferred.

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

	// TODO(OCISDEV-900, finding B7): main wraps TouchBin+Persist in fs.um.RunInBaseScope
	// (upload.go:316) so the .bin/.info files get correct posix ownership. That usermapper
	// lives in decomposedfs; the driver-agnostic coordinator can't reach it. posix-only
	// (no-op on ocis/s3ng). Same root cause as SpaceGid; deferred.
	if err := session.TouchBin(); err != nil {
		return nil, fmt.Errorf("coordinator: could not create bin file: %w", err)
	}
	if err := session.Persist(ctx); err != nil {
		session.Cleanup(true, true)
		return nil, fmt.Errorf("coordinator: could not persist session: %w", err)
	}

	metrics.UploadSessionsInitiated.Inc()

	if uploadLength == 0 {
		// zero-length uploads have no bytes to append, so finish immediately (main: upload.go:333)
		if err := c.finishUpload(ctx, session); err != nil {
			return nil, err
		}
	}

	return map[string]string{
		"simple": session.ID(),
		"tus":    session.ID(),
	}, nil
}

// finishUpload lands a fully-received upload: create the node (new files), verify
// checksums, then commit the staged bytes. Zero-length uploads always finish here
// synchronously; the async postprocessing path is not ported yet.
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

	return c.finishSync(ctx, session)
}

// touchAndMark creates the node for new files (via the public TouchFile, since
// CommitUpload requires an existing node) and marks it as processing. TouchFile
// mints the real node id, so we overwrite the id minted at initiate.
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

// finishSync commits the staged bytes inline, then unmarks processing and cleans up.
func (c *coordinator) finishSync(ctx context.Context, session Session) error {
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
	session.Cleanup(true, true)
	metrics.UploadSessionsFinalized.Inc()
	return nil
}

// rollback unmarks processing, cleans up session files, and deletes the node if
// this upload created it (NodeExists=false at initiation).
func (c *coordinator) rollback(ctx context.Context, session Session) {
	ref := session.Reference()
	_ = c.fs.MarkProcessing(ctx, &ref, false, session.ID())
	session.Cleanup(true, true)
	if !session.NodeExists() {
		_, _ = c.fs.Delete(ctx, &ref)
	}
}

// verifyAndStoreChecksums computes checksums over the staged binary, validates any
// client-supplied checksum, and stores the results on the session for CommitUpload.
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
			return checkErr
		}
	}
	session.SetChecksums(sha1h.Sum(nil), md5h.Sum(nil), adler32h.Sum(nil))
	return nil
}

// rewriteChunkedRef parses a legacy chunking-v1 path, returning a reference to the
// real target file plus the original chunk name.
func rewriteChunkedRef(ref *provider.Reference) (*provider.Reference, string, error) {
	ci, err := chunking.GetChunkBLOBInfo(ref.GetPath())
	if err != nil {
		return nil, "", errtypes.BadRequest(err.Error())
	}
	return &provider.Reference{ResourceId: ref.ResourceId, Path: ci.Path}, filepath.Base(ref.GetPath()), nil
}
