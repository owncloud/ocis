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

package tree

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
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
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/lockedfile"
	"github.com/rs/zerolog/log"
	"go-micro.dev/v4/store"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/cs3org/reva/pkg/storage/utils/decomposedfs/tree")
}

//go:generate make --no-print-directory -C ../../../../.. mockery NAME=Blobstore

// Blobstore defines an interface for storing blobs in a blobstore
type Blobstore interface {
	Upload(node *node.Node, source string) error
	Download(node *node.Node) (io.ReadCloser, error)
	Delete(node *node.Node) error
}

// PathLookup defines the interface for the lookup component
type PathLookup interface {
	NodeFromResource(ctx context.Context, ref *provider.Reference) (*node.Node, error)
	NodeFromID(ctx context.Context, id *provider.ResourceId) (n *node.Node, err error)

	InternalRoot() string
	InternalPath(spaceID, nodeID string) string
	Path(ctx context.Context, n *node.Node, hasPermission node.PermissionFunc) (path string, err error)
	MetadataBackend() metadata.Backend
	ReadBlobSizeAttr(ctx context.Context, path string) (int64, error)
	ReadBlobIDAttr(ctx context.Context, path string) (string, error)
	TypeFromPath(ctx context.Context, path string) provider.ResourceType
}

// Tree manages a hierarchical tree
type Tree struct {
	lookup    PathLookup
	blobstore Blobstore

	options *options.Options

	idCache store.Store
}

// PermissionCheckFunc defined a function used to check resource permissions
type PermissionCheckFunc func(rp *provider.ResourcePermissions) bool

// New returns a new instance of Tree
func New(lu PathLookup, bs Blobstore, o *options.Options, cache store.Store) *Tree {
	return &Tree{
		lookup:    lu,
		blobstore: bs,
		options:   o,
		idCache:   cache,
	}
}

// Setup prepares the tree structure
func (t *Tree) Setup() error {
	// create data paths for internal layout
	dataPaths := []string{
		filepath.Join(t.options.Root, "spaces"),
		// notes contain symlinks from nodes/<u-u-i-d>/uploads/<uploadid> to ../../uploads/<uploadid>
		// better to keep uploads on a fast / volatile storage before a workflow finally moves them to the nodes dir
		filepath.Join(t.options.Root, "uploads"),
	}
	for _, v := range dataPaths {
		err := os.MkdirAll(v, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetMD returns the metadata of a node in the tree
func (t *Tree) GetMD(ctx context.Context, n *node.Node) (os.FileInfo, error) {
	md, err := os.Stat(n.InternalPath())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errtypes.NotFound(n.ID)
		}
		return nil, errors.Wrap(err, "tree: error stating "+n.ID)
	}

	return md, nil
}

// TouchFile creates a new empty file
func (t *Tree) TouchFile(ctx context.Context, n *node.Node, markprocessing bool, mtime string) error {
	if n.Exists {
		if markprocessing {
			return n.SetXattr(ctx, prefixes.StatusPrefix, []byte(node.ProcessingStatus))
		}

		return errtypes.AlreadyExists(n.ID)
	}

	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	n.SetType(provider.ResourceType_RESOURCE_TYPE_FILE)

	nodePath := n.InternalPath()
	if err := os.MkdirAll(filepath.Dir(nodePath), 0700); err != nil {
		return errors.Wrap(err, "Decomposedfs: error creating node")
	}
	_, err := os.Create(nodePath)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: error creating node")
	}

	attributes := n.NodeMetadata(ctx)
	if markprocessing {
		attributes[prefixes.StatusPrefix] = []byte(node.ProcessingStatus)
	}
	if mtime != "" {
		if err := n.SetMtimeString(mtime); err != nil {
			return errors.Wrap(err, "Decomposedfs: could not set mtime")
		}
	}
	err = n.SetXattrsWithContext(ctx, attributes, true)
	if err != nil {
		return err
	}

	// link child name to parent if it is new
	childNameLink := filepath.Join(n.ParentPath(), n.Name)
	var link string
	link, err = os.Readlink(childNameLink)
	if err == nil && link != "../"+n.ID {
		if err = os.Remove(childNameLink); err != nil {
			return errors.Wrap(err, "Decomposedfs: could not remove symlink child entry")
		}
	}
	if errors.Is(err, iofs.ErrNotExist) || link != "../"+n.ID {
		relativeNodePath := filepath.Join("../../../../../", lookup.Pathify(n.ID, 4, 2))
		if err = os.Symlink(relativeNodePath, childNameLink); err != nil {
			return errors.Wrap(err, "Decomposedfs: could not symlink child entry")
		}
	}

	return t.Propagate(ctx, n, 0)
}

// CreateDir creates a new directory entry in the tree
func (t *Tree) CreateDir(ctx context.Context, n *node.Node) (err error) {
	ctx, span := tracer.Start(ctx, "CreateDir")
	defer span.End()
	if n.Exists {
		return errtypes.AlreadyExists(n.ID) // path?
	}

	// create a directory node
	n.SetType(provider.ResourceType_RESOURCE_TYPE_CONTAINER)
	if n.ID == "" {
		n.ID = uuid.New().String()
	}

	err = t.createDirNode(ctx, n)
	if err != nil {
		return
	}

	// make child appear in listings
	relativeNodePath := filepath.Join("../../../../../", lookup.Pathify(n.ID, 4, 2))
	ctx, subspan := tracer.Start(ctx, "os.Symlink")
	err = os.Symlink(relativeNodePath, filepath.Join(n.ParentPath(), n.Name))
	subspan.End()
	if err != nil {
		// no better way to check unfortunately
		if !strings.Contains(err.Error(), "file exists") {
			return
		}

		// try to remove the node
		ctx, subspan = tracer.Start(ctx, "os.RemoveAll")
		e := os.RemoveAll(n.InternalPath())
		subspan.End()
		if e != nil {
			appctx.GetLogger(ctx).Debug().Err(e).Msg("cannot delete node")
		}
		return errtypes.AlreadyExists(err.Error())
	}
	return t.Propagate(ctx, n, 0)
}

// Move replaces the target with the source
func (t *Tree) Move(ctx context.Context, oldNode *node.Node, newNode *node.Node) (err error) {
	if oldNode.SpaceID != newNode.SpaceID {
		// WebDAV RFC https://www.rfc-editor.org/rfc/rfc4918#section-9.9.4 says to use
		// > 502 (Bad Gateway) - This may occur when the destination is on another
		// > server and the destination server refuses to accept the resource.
		// > This could also occur when the destination is on another sub-section
		// > of the same server namespace.
		// but we only have a not supported error
		return errtypes.NotSupported("cannot move across spaces")
	}
	// if target exists delete it without trashing it
	if newNode.Exists {
		// TODO make sure all children are deleted
		if err := os.RemoveAll(newNode.InternalPath()); err != nil {
			return errors.Wrap(err, "Decomposedfs: Move: error deleting target node "+newNode.ID)
		}
	}

	// remove cache entry in any case to avoid inconsistencies
	defer func() { _ = t.idCache.Delete(filepath.Join(oldNode.ParentPath(), oldNode.Name)) }()

	// Always target the old node ID for xattr updates.
	// The new node id is empty if the target does not exist
	// and we need to overwrite the new one when overwriting an existing path.
	// are we just renaming (parent stays the same)?
	if oldNode.ParentID == newNode.ParentID {

		// parentPath := t.lookup.InternalPath(oldNode.SpaceID, oldNode.ParentID)
		parentPath := oldNode.ParentPath()

		// rename child
		err = os.Rename(
			filepath.Join(parentPath, oldNode.Name),
			filepath.Join(parentPath, newNode.Name),
		)
		if err != nil {
			return errors.Wrap(err, "Decomposedfs: could not rename child")
		}

		// update name attribute
		if err := oldNode.SetXattrString(ctx, prefixes.NameAttr, newNode.Name); err != nil {
			return errors.Wrap(err, "Decomposedfs: could not set name attribute")
		}

		return t.Propagate(ctx, newNode, 0)
	}

	// we are moving the node to a new parent, any target has been removed
	// bring old node to the new parent

	// rename child
	err = os.Rename(
		filepath.Join(oldNode.ParentPath(), oldNode.Name),
		filepath.Join(newNode.ParentPath(), newNode.Name),
	)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: could not move child")
	}

	// update target parentid and name
	attribs := node.Attributes{}
	attribs.SetString(prefixes.ParentidAttr, newNode.ParentID)
	attribs.SetString(prefixes.NameAttr, newNode.Name)
	if err := oldNode.SetXattrsWithContext(ctx, attribs, true); err != nil {
		return errors.Wrap(err, "Decomposedfs: could not update old node attributes")
	}

	// the size diff is the current treesize or blobsize of the old/source node
	var sizeDiff int64
	if oldNode.IsDir(ctx) {
		treeSize, err := oldNode.GetTreeSize(ctx)
		if err != nil {
			return err
		}
		sizeDiff = int64(treeSize)
	} else {
		sizeDiff = oldNode.Blobsize
	}

	// TODO inefficient because we might update several nodes twice, only propagate unchanged nodes?
	// collect in a list, then only stat each node once
	// also do this in a go routine ... webdav should check the etag async

	err = t.Propagate(ctx, oldNode, -sizeDiff)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: Move: could not propagate old node")
	}
	err = t.Propagate(ctx, newNode, sizeDiff)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: Move: could not propagate new node")
	}
	return nil
}

func readChildNodeFromLink(ctx context.Context, path string) (string, error) {
	_, span := tracer.Start(ctx, "readChildNodeFromLink")
	defer span.End()
	link, err := os.Readlink(path)
	if err != nil {
		return "", err
	}
	nodeID := strings.TrimLeft(link, "/.")
	nodeID = strings.ReplaceAll(nodeID, "/", "")
	return nodeID, nil
}

// ListFolder lists the content of a folder node
func (t *Tree) ListFolder(ctx context.Context, n *node.Node) ([]*node.Node, error) {
	ctx, span := tracer.Start(ctx, "ListFolder")
	defer span.End()
	dir := n.InternalPath()

	_, subspan := tracer.Start(ctx, "os.Open")
	f, err := os.Open(dir)
	subspan.End()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errtypes.NotFound(dir)
		}
		return nil, errors.Wrap(err, "tree: error listing "+dir)
	}
	defer f.Close()

	_, subspan = tracer.Start(ctx, "f.Readdirnames")
	names, err := f.Readdirnames(0)
	subspan.End()
	if err != nil {
		return nil, err
	}

	numWorkers := t.options.MaxConcurrency
	if len(names) < numWorkers {
		numWorkers = len(names)
	}
	work := make(chan string)
	results := make(chan *node.Node)

	g, ctx := errgroup.WithContext(ctx)

	// Distribute work
	g.Go(func() error {
		defer close(work)
		for _, name := range names {
			select {
			case work <- name:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	// Spawn workers that'll concurrently work the queue
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for name := range work {
				path := filepath.Join(dir, name)
				nodeID := getNodeIDFromCache(ctx, path, t.idCache)
				if nodeID == "" {
					nodeID, err = readChildNodeFromLink(ctx, path)
					if err != nil {
						return err
					}
					err = storeNodeIDInCache(ctx, path, nodeID, t.idCache)
					if err != nil {
						return err
					}
				}

				child, err := node.ReadNode(ctx, t.lookup, n.SpaceID, nodeID, false, n.SpaceRoot, true)
				if err != nil {
					return err
				}

				// prevent listing denied resources
				if !child.IsDenied(ctx) {
					if child.SpaceRoot == nil {
						child.SpaceRoot = n.SpaceRoot
					}
					select {
					case results <- child:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			}
			return nil
		})
	}
	// Wait for things to settle down, then close results chan
	go func() {
		_ = g.Wait() // error is checked later
		close(results)
	}()

	retNodes := []*node.Node{}
	for n := range results {
		retNodes = append(retNodes, n)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return retNodes, nil
}

// Delete deletes a node in the tree by moving it to the trash
func (t *Tree) Delete(ctx context.Context, n *node.Node) (err error) {
	path := filepath.Join(n.ParentPath(), n.Name)
	// remove entry from cache immediately to avoid inconsistencies
	defer func() { _ = t.idCache.Delete(path) }()

	deletingSharedResource := ctx.Value(appctx.DeletingSharedResource)

	if deletingSharedResource != nil && deletingSharedResource.(bool) {
		src := filepath.Join(n.ParentPath(), n.Name)
		return os.Remove(src)
	}

	// get the original path
	origin, err := t.lookup.Path(ctx, n, node.NoCheck)
	if err != nil {
		return
	}

	// set origin location in metadata
	nodePath := n.InternalPath()
	if err := n.SetXattrString(ctx, prefixes.TrashOriginAttr, origin); err != nil {
		return err
	}

	var sizeDiff int64
	if n.IsDir(ctx) {
		treesize, err := n.GetTreeSize(ctx)
		if err != nil {
			return err // TODO calculate treesize if it is not set
		}
		sizeDiff = -int64(treesize)
	} else {
		sizeDiff = -n.Blobsize
	}

	deletionTime := time.Now().UTC().Format(time.RFC3339Nano)

	// Prepare the trash
	trashLink := filepath.Join(t.options.Root, "spaces", lookup.Pathify(n.SpaceRoot.ID, 1, 2), "trash", lookup.Pathify(n.ID, 4, 2))
	if err := os.MkdirAll(filepath.Dir(trashLink), 0700); err != nil {
		// Roll back changes
		_ = n.RemoveXattr(ctx, prefixes.TrashOriginAttr)
		return err
	}

	// FIXME can we just move the node into the trash dir? instead of adding another symlink and appending a trash timestamp?
	// can we just use the mtime as the trash time?
	// TODO store a trashed by userid

	// first make node appear in the space trash
	// parent id and name are stored as extended attributes in the node itself
	err = os.Symlink("../../../../../nodes/"+lookup.Pathify(n.ID, 4, 2)+node.TrashIDDelimiter+deletionTime, trashLink)
	if err != nil {
		// Roll back changes
		_ = n.RemoveXattr(ctx, prefixes.TrashOriginAttr)
		return
	}

	// at this point we have a symlink pointing to a non existing destination, which is fine

	// rename the trashed node so it is not picked up when traversing up the tree and matches the symlink
	trashPath := nodePath + node.TrashIDDelimiter + deletionTime
	err = os.Rename(nodePath, trashPath)
	if err != nil {
		// To roll back changes
		// TODO remove symlink
		// Roll back changes
		_ = n.RemoveXattr(ctx, prefixes.TrashOriginAttr)
		return
	}
	err = t.lookup.MetadataBackend().Rename(nodePath, trashPath)
	if err != nil {
		_ = n.RemoveXattr(ctx, prefixes.TrashOriginAttr)
		_ = os.Rename(trashPath, nodePath)
		return
	}

	// Remove lock file if it exists
	_ = os.Remove(n.LockFilePath())

	// finally remove the entry from the parent dir
	if err = os.Remove(path); err != nil {
		// To roll back changes
		// TODO revert the rename
		// TODO remove symlink
		// Roll back changes
		_ = n.RemoveXattr(ctx, prefixes.TrashOriginAttr)
		return
	}

	return t.Propagate(ctx, n, sizeDiff)
}

// RestoreRecycleItemFunc returns a node and a function to restore it from the trash.
func (t *Tree) RestoreRecycleItemFunc(ctx context.Context, spaceid, key, trashPath string, targetNode *node.Node) (*node.Node, *node.Node, func() error, error) {
	recycleNode, trashItem, deletedNodePath, origin, err := t.readRecycleItem(ctx, spaceid, key, trashPath)
	if err != nil {
		return nil, nil, nil, err
	}

	targetRef := &provider.Reference{
		ResourceId: &provider.ResourceId{SpaceId: spaceid, OpaqueId: spaceid},
		Path:       utils.MakeRelativePath(origin),
	}

	if targetNode == nil {
		targetNode, err = t.lookup.NodeFromResource(ctx, targetRef)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if err := targetNode.CheckLock(ctx); err != nil {
		return nil, nil, nil, err
	}

	parent, err := targetNode.Parent(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	fn := func() error {
		if targetNode.Exists {
			return errtypes.AlreadyExists("origin already exists")
		}

		// add the entry for the parent dir
		err = os.Symlink("../../../../../"+lookup.Pathify(recycleNode.ID, 4, 2), filepath.Join(targetNode.ParentPath(), targetNode.Name))
		if err != nil {
			return err
		}

		// rename to node only name, so it is picked up by id
		nodePath := recycleNode.InternalPath()

		// attempt to rename only if we're not in a subfolder
		if deletedNodePath != nodePath {
			err = os.Rename(deletedNodePath, nodePath)
			if err != nil {
				return err
			}
			err = t.lookup.MetadataBackend().Rename(deletedNodePath, nodePath)
			if err != nil {
				return err
			}
		}

		targetNode.Exists = true

		attrs := node.Attributes{}
		attrs.SetString(prefixes.NameAttr, targetNode.Name)
		if trashPath != "" {
			// set ParentidAttr to restorePath's node parent id
			attrs.SetString(prefixes.ParentidAttr, targetNode.ParentID)
		}

		if err = recycleNode.SetXattrsWithContext(ctx, attrs, true); err != nil {
			return errors.Wrap(err, "Decomposedfs: could not update recycle node")
		}

		// delete item link in trash
		deletePath := trashItem
		if trashPath != "" && trashPath != "/" {
			resolvedTrashRoot, err := filepath.EvalSymlinks(trashItem)
			if err != nil {
				return errors.Wrap(err, "Decomposedfs: could not resolve trash root")
			}
			deletePath = filepath.Join(resolvedTrashRoot, trashPath)
		}
		if err = os.Remove(deletePath); err != nil {
			log.Error().Err(err).Str("trashItem", trashItem).Msg("error deleting trash item")
		}

		var sizeDiff int64
		if recycleNode.IsDir(ctx) {
			treeSize, err := recycleNode.GetTreeSize(ctx)
			if err != nil {
				return err
			}
			sizeDiff = int64(treeSize)
		} else {
			sizeDiff = recycleNode.Blobsize
		}
		return t.Propagate(ctx, targetNode, sizeDiff)
	}
	return recycleNode, parent, fn, nil
}

// PurgeRecycleItemFunc returns a node and a function to purge it from the trash
func (t *Tree) PurgeRecycleItemFunc(ctx context.Context, spaceid, key string, path string) (*node.Node, func() error, error) {
	rn, trashItem, deletedNodePath, _, err := t.readRecycleItem(ctx, spaceid, key, path)
	if err != nil {
		return nil, nil, err
	}

	fn := func() error {
		if err := t.removeNode(ctx, deletedNodePath, rn); err != nil {
			return err
		}

		// delete item link in trash
		deletePath := trashItem
		if path != "" && path != "/" {
			resolvedTrashRoot, err := filepath.EvalSymlinks(trashItem)
			if err != nil {
				return errors.Wrap(err, "Decomposedfs: could not resolve trash root")
			}
			deletePath = filepath.Join(resolvedTrashRoot, path)
		}
		if err = os.Remove(deletePath); err != nil {
			log.Error().Err(err).Str("deletePath", deletePath).Msg("error deleting trash item")
			return err
		}

		return nil
	}

	return rn, fn, nil
}

func (t *Tree) removeNode(ctx context.Context, path string, n *node.Node) error {
	// delete the actual node
	if err := utils.RemoveItem(path); err != nil {
		log.Error().Err(err).Str("path", path).Msg("error purging node")
		return err
	}

	if err := t.lookup.MetadataBackend().Purge(path); err != nil {
		log.Error().Err(err).Str("path", t.lookup.MetadataBackend().MetadataPath(path)).Msg("error purging node metadata")
		return err
	}

	// delete blob from blobstore
	if n.BlobID != "" {
		if err := t.DeleteBlob(n); err != nil {
			log.Error().Err(err).Str("blobID", n.BlobID).Msg("error purging nodes blob")
			return err
		}
	}

	// delete revisions
	revs, err := filepath.Glob(n.InternalPath() + node.RevisionIDDelimiter + "*")
	if err != nil {
		log.Error().Err(err).Str("path", n.InternalPath()+node.RevisionIDDelimiter+"*").Msg("glob failed badly")
		return err
	}
	for _, rev := range revs {
		if t.lookup.MetadataBackend().IsMetaFile(rev) {
			continue
		}

		bID, err := t.lookup.ReadBlobIDAttr(ctx, rev)
		if err != nil {
			log.Error().Err(err).Str("revision", rev).Msg("error reading blobid attribute")
			return err
		}

		if err := utils.RemoveItem(rev); err != nil {
			log.Error().Err(err).Str("revision", rev).Msg("error removing revision node")
			return err
		}

		if bID != "" {
			if err := t.DeleteBlob(&node.Node{SpaceID: n.SpaceID, BlobID: bID}); err != nil {
				log.Error().Err(err).Str("revision", rev).Str("blobID", bID).Msg("error removing revision node blob")
				return err
			}
		}

	}

	return nil
}

// Propagate propagates changes to the root of the tree
func (t *Tree) Propagate(ctx context.Context, n *node.Node, sizeDiff int64) (err error) {
	ctx, span := tracer.Start(ctx, "Propagate")
	defer span.End()
	sublog := appctx.GetLogger(ctx).With().
		Str("method", "tree.Propagate").
		Str("spaceid", n.SpaceID).
		Str("nodeid", n.ID).
		Int64("sizeDiff", sizeDiff).
		Logger()

	if !t.options.TreeTimeAccounting && (!t.options.TreeSizeAccounting || sizeDiff == 0) {
		// no propagation enabled
		sublog.Debug().Msg("propagation disabled or nothing to propagate")
		return
	}

	// is propagation enabled for the parent node?
	root := n.SpaceRoot

	// use a sync time and don't rely on the mtime of the current node, as the stat might not change when a rename happened too quickly
	sTime := time.Now().UTC()

	// we loop until we reach the root
	for err == nil && n.ID != root.ID {
		sublog.Debug().Msg("propagating")

		attrs := node.Attributes{}

		var f *lockedfile.File
		// lock parent before reading treesize or tree time

		_, subspan := tracer.Start(ctx, "lockedfile.OpenFile")
		parentFilename := t.lookup.MetadataBackend().LockfilePath(n.ParentPath())
		f, err = lockedfile.OpenFile(parentFilename, os.O_RDWR|os.O_CREATE, 0600)
		subspan.End()
		if err != nil {
			sublog.Error().Err(err).
				Str("parent filename", parentFilename).
				Msg("Propagation failed. Could not open metadata for parent with lock.")
			return err
		}
		// always log error if closing node fails
		defer func() {
			// ignore already closed error
			cerr := f.Close()
			if err == nil && cerr != nil && !errors.Is(cerr, os.ErrClosed) {
				err = cerr // only overwrite err with en error from close if the former was nil
			}
		}()

		if n, err = n.ParentWithReader(ctx, f); err != nil {
			sublog.Error().Err(err).
				Msg("Propagation failed. Could not read parent node.")
			return err
		}

		// TODO none, sync and async?
		if !n.HasPropagation(ctx) {
			sublog.Debug().Str("attr", prefixes.PropagationAttr).Msg("propagation attribute not set or unreadable, not propagating")
			// if the attribute is not set treat it as false / none / no propagation
			return nil
		}

		sublog = sublog.With().Str("spaceid", n.SpaceID).Str("nodeid", n.ID).Logger()

		if t.options.TreeTimeAccounting {
			// update the parent tree time if it is older than the nodes mtime
			updateSyncTime := false

			var tmTime time.Time
			tmTime, err = n.GetTMTime(ctx)
			switch {
			case err != nil:
				// missing attribute, or invalid format, overwrite
				sublog.Debug().Err(err).
					Msg("could not read tmtime attribute, overwriting")
				updateSyncTime = true
			case tmTime.Before(sTime):
				sublog.Debug().
					Time("tmtime", tmTime).
					Time("stime", sTime).
					Msg("parent tmtime is older than node mtime, updating")
				updateSyncTime = true
			default:
				sublog.Debug().
					Time("tmtime", tmTime).
					Time("stime", sTime).
					Dur("delta", sTime.Sub(tmTime)).
					Msg("parent tmtime is younger than node mtime, not updating")
			}

			if updateSyncTime {
				// update the tree time of the parent node
				attrs.SetString(prefixes.TreeMTimeAttr, sTime.UTC().Format(time.RFC3339Nano))
			}

			attrs.SetString(prefixes.TmpEtagAttr, "")
		}

		// size accounting
		if t.options.TreeSizeAccounting && sizeDiff != 0 {
			var newSize uint64

			// read treesize
			treeSize, err := n.GetTreeSize(ctx)
			switch {
			case metadata.IsAttrUnset(err):
				// fallback to calculating the treesize
				sublog.Warn().Msg("treesize attribute unset, falling back to calculating the treesize")
				newSize, err = t.calculateTreeSize(ctx, n.InternalPath())
				if err != nil {
					return err
				}
			case err != nil:
				sublog.Error().Err(err).
					Msg("Faild to propagate treesize change. Error when reading the treesize attribute from parent")
				return err
			case sizeDiff > 0:
				newSize = treeSize + uint64(sizeDiff)
			case uint64(-sizeDiff) > treeSize:
				// The sizeDiff is larger than the current treesize. Which would result in
				// a negative new treesize. Something must have gone wrong with the accounting.
				// Reset the current treesize to 0.
				sublog.Error().Uint64("treeSize", treeSize).Int64("sizeDiff", sizeDiff).
					Msg("Error when updating treesize of parent node. Updated treesize < 0. Reestting to 0")
				newSize = 0
			default:
				newSize = treeSize - uint64(-sizeDiff)
			}

			// update the tree size of the node
			attrs.SetString(prefixes.TreesizeAttr, strconv.FormatUint(newSize, 10))
			sublog.Debug().Uint64("newSize", newSize).Msg("updated treesize of parent node")
		}

		if err = n.SetXattrsWithContext(ctx, attrs, false); err != nil {
			sublog.Error().Err(err).Msg("Failed to update extend attributes of parent node")
			return err
		}

		// Release node lock early, ignore already closed error
		_, subspan = tracer.Start(ctx, "f.Close")
		cerr := f.Close()
		subspan.End()
		if cerr != nil && !errors.Is(cerr, os.ErrClosed) {
			sublog.Error().Err(cerr).Msg("Failed to close parent node and release lock")
			return cerr
		}
	}
	if err != nil {
		sublog.Error().Err(err).Msg("error propagating")
		return
	}
	return
}

func (t *Tree) calculateTreeSize(ctx context.Context, childrenPath string) (uint64, error) {
	ctx, span := tracer.Start(ctx, "calculateTreeSize")
	defer span.End()
	var size uint64

	f, err := os.Open(childrenPath)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Str("childrenPath", childrenPath).Msg("could not open dir")
		return 0, err
	}
	defer f.Close()

	names, err := f.Readdirnames(0)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Str("childrenPath", childrenPath).Msg("could not read dirnames")
		return 0, err
	}
	for i := range names {
		cPath := filepath.Join(childrenPath, names[i])
		resolvedPath, err := filepath.EvalSymlinks(cPath)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("childpath", cPath).Msg("could not resolve child entry symlink")
			continue // continue after an error
		}

		// raw read of the attributes for performance reasons
		attribs, err := t.lookup.MetadataBackend().All(ctx, resolvedPath)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("childpath", cPath).Msg("could not read attributes of child entry")
			continue // continue after an error
		}
		sizeAttr := ""
		if string(attribs[prefixes.TypeAttr]) == strconv.FormatUint(uint64(provider.ResourceType_RESOURCE_TYPE_FILE), 10) {
			sizeAttr = string(attribs[prefixes.BlobsizeAttr])
		} else {
			sizeAttr = string(attribs[prefixes.TreesizeAttr])
		}
		csize, err := strconv.ParseInt(sizeAttr, 10, 64)
		if err != nil {
			return 0, errors.Wrapf(err, "invalid blobsize xattr format")
		}
		size += uint64(csize)
	}
	return size, err
}

// WriteBlob writes a blob to the blobstore
func (t *Tree) WriteBlob(node *node.Node, source string) error {
	return t.blobstore.Upload(node, source)
}

// ReadBlob reads a blob from the blobstore
func (t *Tree) ReadBlob(node *node.Node) (io.ReadCloser, error) {
	if node.BlobID == "" {
		// there is no blob yet - we are dealing with a 0 byte file
		return io.NopCloser(bytes.NewReader([]byte{})), nil
	}
	return t.blobstore.Download(node)
}

// DeleteBlob deletes a blob from the blobstore
func (t *Tree) DeleteBlob(node *node.Node) error {
	if node == nil {
		return fmt.Errorf("could not delete blob, nil node was given")
	}
	if node.BlobID == "" {
		return fmt.Errorf("could not delete blob, node with empty blob id was given")
	}

	return t.blobstore.Delete(node)
}

// TODO check if node exists?
func (t *Tree) createDirNode(ctx context.Context, n *node.Node) (err error) {
	ctx, span := tracer.Start(ctx, "createDirNode")
	defer span.End()
	// create a directory node
	nodePath := n.InternalPath()
	if err := os.MkdirAll(nodePath, 0700); err != nil {
		return errors.Wrap(err, "Decomposedfs: error creating node")
	}

	attributes := n.NodeMetadata(ctx)
	attributes[prefixes.TreesizeAttr] = []byte("0") // initialize as empty, TODO why bother? if it is not set we could treat it as 0?
	if t.options.TreeTimeAccounting || t.options.TreeSizeAccounting {
		attributes[prefixes.PropagationAttr] = []byte("1") // mark the node for propagation
	}
	return n.SetXattrsWithContext(ctx, attributes, true)
}

var nodeIDRegep = regexp.MustCompile(`.*/nodes/([^.]*).*`)

// TODO refactor the returned params into Node properties? would make all the path transformations go away...
func (t *Tree) readRecycleItem(ctx context.Context, spaceID, key, path string) (recycleNode *node.Node, trashItem string, deletedNodePath string, origin string, err error) {
	if key == "" {
		return nil, "", "", "", errtypes.InternalError("key is empty")
	}

	backend := t.lookup.MetadataBackend()
	var nodeID string

	trashItem = filepath.Join(t.lookup.InternalRoot(), "spaces", lookup.Pathify(spaceID, 1, 2), "trash", lookup.Pathify(key, 4, 2))
	resolvedTrashItem, err := filepath.EvalSymlinks(trashItem)
	if err != nil {
		return
	}
	deletedNodePath, err = filepath.EvalSymlinks(filepath.Join(resolvedTrashItem, path))
	if err != nil {
		return
	}
	nodeID = nodeIDRegep.ReplaceAllString(deletedNodePath, "$1")
	nodeID = strings.ReplaceAll(nodeID, "/", "")

	recycleNode = node.New(spaceID, nodeID, "", "", 0, "", provider.ResourceType_RESOURCE_TYPE_INVALID, nil, t.lookup)
	recycleNode.SpaceRoot, err = node.ReadNode(ctx, t.lookup, spaceID, spaceID, false, nil, false)
	if err != nil {
		return
	}
	recycleNode.SetType(t.lookup.TypeFromPath(ctx, deletedNodePath))

	var attrBytes []byte
	if recycleNode.Type(ctx) == provider.ResourceType_RESOURCE_TYPE_FILE {
		// lookup blobID in extended attributes
		if attrBytes, err = backend.Get(ctx, deletedNodePath, prefixes.BlobIDAttr); err == nil {
			recycleNode.BlobID = string(attrBytes)
		} else {
			return
		}

		// lookup blobSize in extended attributes
		if recycleNode.Blobsize, err = backend.GetInt64(ctx, deletedNodePath, prefixes.BlobsizeAttr); err != nil {
			return
		}
	}

	// lookup parent id in extended attributes
	if attrBytes, err = backend.Get(ctx, deletedNodePath, prefixes.ParentidAttr); err == nil {
		recycleNode.ParentID = string(attrBytes)
	} else {
		return
	}

	// lookup name in extended attributes
	if attrBytes, err = backend.Get(ctx, deletedNodePath, prefixes.NameAttr); err == nil {
		recycleNode.Name = string(attrBytes)
	} else {
		return
	}

	// get origin node, is relative to space root
	origin = "/"

	// lookup origin path in extended attributes
	if attrBytes, err = backend.Get(ctx, resolvedTrashItem, prefixes.TrashOriginAttr); err == nil {
		origin = filepath.Join(string(attrBytes), path)
	} else {
		log.Error().Err(err).Str("trashItem", trashItem).Str("deletedNodePath", deletedNodePath).Msg("could not read origin path, restoring to /")
	}

	return
}

func getNodeIDFromCache(ctx context.Context, path string, cache store.Store) string {
	_, span := tracer.Start(ctx, "getNodeIDFromCache")
	defer span.End()
	recs, err := cache.Read(path)
	if err == nil && len(recs) > 0 {
		return string(recs[0].Value)
	}
	return ""
}

func storeNodeIDInCache(ctx context.Context, path string, nodeID string, cache store.Store) error {
	_, span := tracer.Start(ctx, "storeNodeIDInCache")
	defer span.End()
	return cache.Write(&store.Record{
		Key:   path,
		Value: []byte(nodeID),
	})
}
