// Copyright 2018-2022 CERN
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
	"encoding/json"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/lockedfile"
	tusd "github.com/tus/tusd/pkg/handler"
)

var _idRegexp = regexp.MustCompile(".*/([^/]+).info")

// PermissionsChecker defines an interface for checking permissions on a Node
type PermissionsChecker interface {
	AssemblePermissions(ctx context.Context, n *node.Node) (ap provider.ResourcePermissions, err error)
}

// OcisStore manages upload sessions
type OcisStore struct {
	lu      node.PathLookup
	tp      Tree
	root    string
	pub     events.Publisher
	async   bool
	tknopts options.TokenOptions
}

// NewSessionStore returns a new OcisStore
func NewSessionStore(lu node.PathLookup, tp Tree, root string, pub events.Publisher, async bool, tknopts options.TokenOptions) *OcisStore {
	return &OcisStore{
		lu:      lu,
		tp:      tp,
		root:    root,
		pub:     pub,
		async:   async,
		tknopts: tknopts,
	}
}

// New returns a new upload session
func (store OcisStore) New(ctx context.Context) *OcisSession {
	return &OcisSession{
		store: store,
		info: tusd.FileInfo{
			ID: uuid.New().String(),
			Storage: map[string]string{
				"Type": "OCISStore",
			},
			MetaData: tusd.MetaData{},
		},
	}
}

// List lists all upload sessions
func (store OcisStore) List(ctx context.Context) ([]*OcisSession, error) {
	uploads := []*OcisSession{}
	infoFiles, err := filepath.Glob(filepath.Join(store.root, "uploads", "*.info"))
	if err != nil {
		return nil, err
	}

	for _, info := range infoFiles {
		id := strings.TrimSuffix(filepath.Base(info), filepath.Ext(info))
		progress, err := store.Get(ctx, id)
		if err != nil {
			appctx.GetLogger(ctx).Error().Interface("path", info).Msg("Decomposedfs: could not getUploadSession")
			continue
		}

		uploads = append(uploads, progress)
	}
	return uploads, nil
}

// Get returns the upload session for the given upload id
func (store OcisStore) Get(ctx context.Context, id string) (*OcisSession, error) {
	sessionPath := sessionPath(store.root, id)
	match := _idRegexp.FindStringSubmatch(sessionPath)
	if match == nil || len(match) < 2 {
		return nil, fmt.Errorf("invalid upload path")
	}

	session := OcisSession{
		store: store,
		info:  tusd.FileInfo{},
	}
	lock, err := lockedfile.Open(sessionPath + ".lock")
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = tusd.ErrNotFound
		}
		return nil, err
	}
	defer lock.Close()
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = tusd.ErrNotFound
		}
		return nil, err
	}
	lock.Close() // release lock asap

	if err := json.Unmarshal(data, &session.info); err != nil {
		return nil, err
	}

	stat, err := os.Stat(session.binPath())
	if err != nil {
		if os.IsNotExist(err) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = tusd.ErrNotFound
		}
		return nil, err
	}

	session.info.Offset = stat.Size()

	return &session, nil
}

// Session is the interface used by the Cleanup call
type Session interface {
	ID() string
	Node(ctx context.Context) (*node.Node, error)
	Context(ctx context.Context) context.Context
	Cleanup(revertNodeMetadata, cleanBin, cleanInfo bool)
}

// Cleanup cleans upload metadata, binary data and processing status as necessary
func (store OcisStore) Cleanup(ctx context.Context, session Session, revertNodeMetadata, keepUpload, unmarkPostprocessing bool) {
	ctx, span := tracer.Start(session.Context(ctx), "Cleanup")
	defer span.End()
	session.Cleanup(revertNodeMetadata, !keepUpload, !keepUpload)

	// unset processing status
	if unmarkPostprocessing {
		n, err := session.Node(ctx)
		if err != nil {
			appctx.GetLogger(ctx).Info().Str("session", session.ID()).Err(err).Msg("could not read node")
			return
		}
		// FIXME: after cleanup the node might already be deleted ...
		if n != nil { // node can be nil when there was an error before it was created (eg. checksum-mismatch)
			if err := n.UnmarkProcessing(ctx, session.ID()); err != nil {
				appctx.GetLogger(ctx).Info().Str("path", n.InternalPath()).Err(err).Msg("unmarking processing failed")
			}
		}
	}
}

// CreateNodeForUpload will create the target node for the Upload
// TODO move this to the node package as NodeFromUpload?
// should we in InitiateUpload create the node first? and then the upload?
func (store OcisStore) CreateNodeForUpload(session *OcisSession, initAttrs node.Attributes) (*node.Node, error) {
	ctx, span := tracer.Start(session.Context(context.Background()), "CreateNodeForUpload")
	defer span.End()
	n := node.New(
		session.SpaceID(),
		session.NodeID(),
		session.NodeParentID(),
		session.Filename(),
		session.Size(),
		session.ID(),
		provider.ResourceType_RESOURCE_TYPE_FILE,
		nil,
		store.lu,
	)
	var err error
	n.SpaceRoot, err = node.ReadNode(ctx, store.lu, session.SpaceID(), session.SpaceID(), false, nil, false)
	if err != nil {
		return nil, err
	}

	// check lock
	if err := n.CheckLock(ctx); err != nil {
		return nil, err
	}

	var f *lockedfile.File
	if session.NodeExists() { // TODO this is wrong. The node should be created when the upload starts, the revisions should be created independently of the node
		// we do not need to propagate a change when a node is created, only when the upload is ready.
		// that still creates problems for desktop clients because if another change causes propagation it will detects an empty file
		// so the first upload has to point to the first revision with the expected size. The file cannot be downloaded, but it can be overwritten (which will create a new revision and make the node reflect the latest revision)
		// any finished postprocessing will not affect the node metadata.
		// *thinking* but then initializing an upload will lock the file until the upload has finished. That sucks.
		// so we have to check if the node has been created meanwhile (well, only in case the upload does not know the nodeid ... or the NodeExists array that is checked by session.NodeExists())
		// FIXME look at the disk again to see if the file has been created in between, or just try initializing a new node and do the update existing node as a fallback. <- the latter!

		f, err = store.updateExistingNode(ctx, session, n, session.SpaceID(), uint64(session.Size()))
		if f != nil {
			appctx.GetLogger(ctx).Info().Str("lockfile", f.Name()).Interface("err", err).Msg("got lock file from updateExistingNode")
		}
	} else {
		f, err = store.initNewNode(ctx, session, n, uint64(session.Size()))
		if f != nil {
			appctx.GetLogger(ctx).Info().Str("lockfile", f.Name()).Interface("err", err).Msg("got lock file from initNewNode")
		}
	}
	defer func() {
		if f == nil {
			return
		}
		if err := f.Close(); err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("nodeid", n.ID).Str("parentid", n.ParentID).Msg("could not close lock")
		}
	}()
	if err != nil {
		return nil, err
	}

	mtime := time.Now()
	if !session.MTime().IsZero() {
		// overwrite mtime if requested
		mtime = session.MTime()
	}

	// overwrite technical information
	initAttrs.SetString(prefixes.MTimeAttr, mtime.UTC().Format(time.RFC3339Nano))
	initAttrs.SetInt64(prefixes.TypeAttr, int64(provider.ResourceType_RESOURCE_TYPE_FILE))
	initAttrs.SetString(prefixes.ParentidAttr, n.ParentID)
	initAttrs.SetString(prefixes.NameAttr, n.Name)
	initAttrs.SetString(prefixes.BlobIDAttr, n.BlobID)
	initAttrs.SetInt64(prefixes.BlobsizeAttr, n.Blobsize)
	initAttrs.SetString(prefixes.StatusPrefix, node.ProcessingStatus+session.ID())

	// update node metadata with new blobid etc
	err = n.SetXattrsWithContext(ctx, initAttrs, false)
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: could not write metadata")
	}

	if err := session.Persist(ctx); err != nil {
		return nil, err
	}

	return n, nil
}

func (store OcisStore) initNewNode(ctx context.Context, session *OcisSession, n *node.Node, fsize uint64) (*lockedfile.File, error) {
	// create folder structure (if needed)
	if err := os.MkdirAll(filepath.Dir(n.InternalPath()), 0700); err != nil {
		return nil, err
	}

	// create and write lock new node metadata
	f, err := lockedfile.OpenFile(store.lu.MetadataBackend().LockfilePath(n.InternalPath()), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	// we also need to touch the actual node file here it stores the mtime of the resource
	h, err := os.OpenFile(n.InternalPath(), os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return f, err
	}
	h.Close()

	if _, err := node.CheckQuota(ctx, n.SpaceRoot, false, 0, fsize); err != nil {
		return f, err
	}

	// link child name to parent if it is new
	childNameLink := filepath.Join(n.ParentPath(), n.Name)
	relativeNodePath := filepath.Join("../../../../../", lookup.Pathify(n.ID, 4, 2))
	log := appctx.GetLogger(ctx).With().Str("childNameLink", childNameLink).Str("relativeNodePath", relativeNodePath).Logger()
	log.Info().Msg("initNewNode: creating symlink")

	if err = os.Symlink(relativeNodePath, childNameLink); err != nil {
		log.Info().Err(err).Msg("initNewNode: symlink failed")
		if errors.Is(err, iofs.ErrExist) {
			log.Info().Err(err).Msg("initNewNode: symlink already exists")
			return f, errtypes.AlreadyExists(n.Name)
		}
		return f, errors.Wrap(err, "Decomposedfs: could not symlink child entry")
	}
	log.Info().Msg("initNewNode: symlink created")

	// on a new file the sizeDiff is the fileSize
	session.info.MetaData["sizeDiff"] = strconv.FormatInt(int64(fsize), 10)
	return f, nil
}

func (store OcisStore) updateExistingNode(ctx context.Context, session *OcisSession, n *node.Node, spaceID string, fsize uint64) (*lockedfile.File, error) {
	targetPath := n.InternalPath()

	// write lock existing node before reading any metadata
	f, err := lockedfile.OpenFile(store.lu.MetadataBackend().LockfilePath(targetPath), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	old, _ := node.ReadNode(ctx, store.lu, spaceID, n.ID, false, nil, false)
	if _, err := node.CheckQuota(ctx, n.SpaceRoot, true, uint64(old.Blobsize), fsize); err != nil {
		return f, err
	}

	oldNodeMtime, err := old.GetMTime(ctx)
	if err != nil {
		return f, err
	}
	oldNodeEtag, err := node.CalculateEtag(old.ID, oldNodeMtime)
	if err != nil {
		return f, err
	}

	// When the if-match header was set we need to check if the
	// etag still matches before finishing the upload.
	if session.HeaderIfMatch() != "" && session.HeaderIfMatch() != oldNodeEtag {
		return f, errtypes.Aborted("etag mismatch")
	}

	// When the if-none-match header was set we need to check if any of the
	// etags matches before finishing the upload.
	if session.HeaderIfNoneMatch() != "" {
		if session.HeaderIfNoneMatch() == "*" {
			return f, errtypes.Aborted("etag mismatch, resource exists")
		}
		for _, ifNoneMatchTag := range strings.Split(session.HeaderIfNoneMatch(), ",") {
			if ifNoneMatchTag == oldNodeEtag {
				return f, errtypes.Aborted("etag mismatch")
			}
		}
	}

	// When the if-unmodified-since header was set we need to check if the
	// etag still matches before finishing the upload.
	if session.HeaderIfUnmodifiedSince() != "" {
		ifUnmodifiedSince, err := time.Parse(time.RFC3339Nano, session.HeaderIfUnmodifiedSince())
		if err != nil {
			return f, errtypes.InternalError(fmt.Sprintf("failed to parse if-unmodified-since time: %s", err))
		}

		if oldNodeMtime.After(ifUnmodifiedSince) {
			return f, errtypes.Aborted("if-unmodified-since mismatch")
		}
	}

	session.info.MetaData["versionsPath"] = session.store.lu.InternalPath(spaceID, n.ID+node.RevisionIDDelimiter+oldNodeMtime.UTC().Format(time.RFC3339Nano))
	session.info.MetaData["sizeDiff"] = strconv.FormatInt((int64(fsize) - old.Blobsize), 10)

	// create version node
	if _, err := os.Create(session.info.MetaData["versionsPath"]); err != nil {
		return f, err
	}

	// copy blob metadata to version node
	if err := store.lu.CopyMetadataWithSourceLock(ctx, targetPath, session.info.MetaData["versionsPath"], func(attributeName string, value []byte) (newValue []byte, copy bool) {
		return value, strings.HasPrefix(attributeName, prefixes.ChecksumPrefix) ||
			attributeName == prefixes.TypeAttr ||
			attributeName == prefixes.BlobIDAttr ||
			attributeName == prefixes.BlobsizeAttr ||
			attributeName == prefixes.MTimeAttr
	}, f, true); err != nil {
		return f, err
	}

	// keep mtime from previous version
	if err := os.Chtimes(session.info.MetaData["versionsPath"], oldNodeMtime, oldNodeMtime); err != nil {
		return f, errtypes.InternalError(fmt.Sprintf("failed to change mtime of version node: %s", err))
	}

	return f, nil
}
