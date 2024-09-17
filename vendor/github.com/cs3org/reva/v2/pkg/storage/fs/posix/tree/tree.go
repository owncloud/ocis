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
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-micro.dev/v4/store"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/logger"
	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/options"
	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/trashbin"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/tree/propagator"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/usermapper"
	"github.com/cs3org/reva/v2/pkg/utils"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/cs3org/reva/pkg/storage/utils/decomposedfs/tree")
}

// Blobstore defines an interface for storing blobs in a blobstore
type Blobstore interface {
	Upload(node *node.Node, source string) error
	Download(node *node.Node) (io.ReadCloser, error)
	Delete(node *node.Node) error
}

type Watcher interface {
	Watch(path string)
}

type scanItem struct {
	Path        string
	ForceRescan bool
	Recurse     bool
}

// Tree manages a hierarchical tree
type Tree struct {
	lookup     node.PathLookup
	blobstore  Blobstore
	trashbin   *trashbin.Trashbin
	propagator propagator.Propagator

	options *options.Options

	userMapper    usermapper.Mapper
	idCache       store.Store
	watcher       Watcher
	scanQueue     chan scanItem
	scanDebouncer *ScanDebouncer

	es  events.Stream
	log *zerolog.Logger
}

// PermissionCheckFunc defined a function used to check resource permissions
type PermissionCheckFunc func(rp *provider.ResourcePermissions) bool

// New returns a new instance of Tree
func New(lu node.PathLookup, bs Blobstore, um usermapper.Mapper, trashbin *trashbin.Trashbin, o *options.Options, es events.Stream, cache store.Store) (*Tree, error) {
	log := logger.New()
	scanQueue := make(chan scanItem)
	t := &Tree{
		lookup:     lu,
		blobstore:  bs,
		userMapper: um,
		trashbin:   trashbin,
		options:    o,
		idCache:    cache,
		propagator: propagator.New(lu, &o.Options),
		scanQueue:  scanQueue,
		scanDebouncer: NewScanDebouncer(o.ScanDebounceDelay, func(item scanItem) {
			scanQueue <- item
		}),
		es:  es,
		log: log,
	}

	watchPath := o.WatchPath
	var err error
	switch o.WatchType {
	case "gpfswatchfolder":
		t.watcher, err = NewGpfsWatchFolderWatcher(t, strings.Split(o.WatchFolderKafkaBrokers, ","))
		if err != nil {
			return nil, err
		}
	case "gpfsfileauditlogging":
		t.watcher, err = NewGpfsFileAuditLoggingWatcher(t, o.WatchPath)
		if err != nil {
			return nil, err
		}
	default:
		t.watcher = NewInotifyWatcher(t)
		watchPath = o.Root
	}

	// Start watching for fs events and put them into the queue
	if o.WatchFS {
		go t.watcher.Watch(watchPath)
		go t.workScanQueue()
		go func() {
			_ = t.WarmupIDCache(o.Root, true)
		}()
	}

	return t, nil
}

func (t *Tree) PublishEvent(ev interface{}) {
	if t.es == nil {
		return
	}

	if err := events.Publish(context.Background(), t.es, ev); err != nil {
		t.log.Error().Err(err).Interface("event", ev).Msg("failed to publish event")
	}
}

// Setup prepares the tree structure
func (t *Tree) Setup() error {
	err := os.MkdirAll(t.options.Root, 0700)
	if err != nil {
		return err
	}

	err = os.MkdirAll(t.options.UploadDirectory, 0700)
	if err != nil {
		return err
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

	parentPath := n.ParentPath()
	nodePath := filepath.Join(parentPath, n.Name)

	// lock the meta file
	unlock, err := t.lookup.MetadataBackend().Lock(nodePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = unlock()
	}()

	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	n.SetType(provider.ResourceType_RESOURCE_TYPE_FILE)

	// Set id in cache
	_ = t.lookup.(*lookup.Lookup).CacheID(context.Background(), n.SpaceID, n.ID, nodePath)

	if err := os.MkdirAll(filepath.Dir(nodePath), 0700); err != nil {
		return errors.Wrap(err, "Decomposedfs: error creating node")
	}
	_, err = os.Create(nodePath)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: error creating node")
	}

	attributes := n.NodeMetadata(ctx)
	attributes[prefixes.IDAttr] = []byte(n.ID)
	if markprocessing {
		attributes[prefixes.StatusPrefix] = []byte(node.ProcessingStatus)
	}
	if mtime != "" {
		nodeMTime, err := utils.MTimeToTime(mtime)
		if err != nil {
			return err
		}
		err = os.Chtimes(nodePath, nodeMTime, nodeMTime)
		if err != nil {
			return err
		}
	}

	err = n.SetXattrsWithContext(ctx, attributes, false)
	if err != nil {
		return err
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

	// we are moving the node to a new parent, any target has been removed
	// bring old node to the new parent

	// update target parentid and name
	attribs := node.Attributes{}
	attribs.SetString(prefixes.ParentidAttr, newNode.ParentID)
	attribs.SetString(prefixes.NameAttr, newNode.Name)
	if err := oldNode.SetXattrsWithContext(ctx, attribs, true); err != nil {
		return errors.Wrap(err, "Decomposedfs: could not update old node attributes")
	}

	// rename node
	err = os.Rename(
		filepath.Join(oldNode.ParentPath(), oldNode.Name),
		filepath.Join(newNode.ParentPath(), newNode.Name),
	)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: could not move child")
	}

	// rename the lock (if it exists)
	if _, err := os.Stat(oldNode.LockFilePath()); err == nil {
		err = os.Rename(
			filepath.Join(oldNode.ParentPath(), oldNode.Name+".lock"),
			filepath.Join(newNode.ParentPath(), newNode.Name+".lock"),
		)
		if err != nil {
			return errors.Wrap(err, "Decomposedfs: could not move lock")
		}
	}

	// update the id cache
	if newNode.ID == "" {
		newNode.ID = oldNode.ID
	}
	_ = t.lookup.(*lookup.Lookup).CacheID(ctx, newNode.SpaceID, newNode.ID, filepath.Join(newNode.ParentPath(), newNode.Name))
	// update id cache for the moved subtree.
	if oldNode.IsDir(ctx) {
		err = t.WarmupIDCache(filepath.Join(newNode.ParentPath(), newNode.Name), false)
		if err != nil {
			return err
		}
	}

	err = t.Propagate(ctx, oldNode, 0)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: Move: could not propagate old node")
	}
	err = t.Propagate(ctx, newNode, 0)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: Move: could not propagate new node")
	}
	return nil
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
			if isLockFile(name) || isTrash(name) {
				continue
			}

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
			// switch user if necessary
			spaceGID, ok := ctx.Value(decomposedfs.CtxKeySpaceGID).(uint32)
			if ok {
				unscope, err := t.userMapper.ScopeUserByIds(-1, int(spaceGID))
				if err != nil {
					return errors.Wrap(err, "failed to scope user")
				}
				defer func() { _ = unscope() }()
			}

			for name := range work {
				path := filepath.Join(dir, name)
				nodeID, err := t.lookup.MetadataBackend().Get(ctx, path, prefixes.IDAttr)
				if err != nil {
					if metadata.IsAttrUnset(err) {
						continue
					}
					return err
				}

				child, err := node.ReadNode(ctx, t.lookup, n.SpaceID, string(nodeID), false, n.SpaceRoot, true)
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
func (t *Tree) Delete(ctx context.Context, n *node.Node) error {
	path := n.InternalPath()

	if !strings.HasPrefix(path, t.options.Root) {
		return errtypes.InternalError("invalid internal path")
	}

	// remove entry from cache immediately to avoid inconsistencies
	defer func() { _ = t.idCache.Delete(path) }()

	deletingSharedResource := ctx.Value(appctx.DeletingSharedResource)

	if deletingSharedResource != nil && deletingSharedResource.(bool) {
		src := filepath.Join(n.ParentPath(), n.Name)
		return os.RemoveAll(src)
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

	// Remove lock file if it exists
	_ = os.Remove(n.LockFilePath())

	err := t.trashbin.MoveToTrash(ctx, n, path)
	if err != nil {
		return err
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

	if err := t.lookup.MetadataBackend().Purge(ctx, path); err != nil {
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
	// posixfs doesn't do revisions yet

	return nil
}

// Propagate propagates changes to the root of the tree
func (t *Tree) Propagate(ctx context.Context, n *node.Node, _ int64) (err error) {
	// We do not propagate size diffs here but rely on the assimilation to take care of the tree sizes instead
	return t.propagator.Propagate(ctx, n, 0)
}

// WriteBlob writes a blob to the blobstore
func (t *Tree) WriteBlob(node *node.Node, source string) error {
	return t.blobstore.Upload(node, source)
}

// ReadBlob reads a blob from the blobstore
func (t *Tree) ReadBlob(node *node.Node) (io.ReadCloser, error) {
	return t.blobstore.Download(node)
}

// DeleteBlob deletes a blob from the blobstore
func (t *Tree) DeleteBlob(node *node.Node) error {
	if node == nil {
		return fmt.Errorf("could not delete blob, nil node was given")
	}
	return t.blobstore.Delete(node)
}

// BuildSpaceIDIndexEntry returns the entry for the space id index
func (t *Tree) BuildSpaceIDIndexEntry(spaceID, nodeID string) string {
	return nodeID
}

// ResolveSpaceIDIndexEntry returns the node id for the space id index entry
func (t *Tree) ResolveSpaceIDIndexEntry(spaceid, entry string) (string, string, error) {
	return spaceid, entry, nil
}

// InitNewNode initializes a new node
func (t *Tree) InitNewNode(ctx context.Context, n *node.Node, fsize uint64) (metadata.UnlockFunc, error) {
	// create folder structure (if needed)
	if err := os.MkdirAll(filepath.Dir(n.InternalPath()), 0700); err != nil {
		return nil, err
	}

	// create and write lock new node metadata
	unlock, err := t.lookup.MetadataBackend().Lock(n.InternalPath())
	if err != nil {
		return nil, err
	}

	// we also need to touch the actual node file here it stores the mtime of the resource
	h, err := os.OpenFile(n.InternalPath(), os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return unlock, err
	}
	h.Close()

	if _, err := node.CheckQuota(ctx, n.SpaceRoot, false, 0, fsize); err != nil {
		return unlock, err
	}

	return unlock, nil
}

// TODO check if node exists?
func (t *Tree) createDirNode(ctx context.Context, n *node.Node) (err error) {
	ctx, span := tracer.Start(ctx, "createDirNode")
	defer span.End()

	idcache := t.lookup.(*lookup.Lookup).IDCache
	// create a directory node
	parentPath, ok := idcache.Get(ctx, n.SpaceID, n.ParentID)
	if !ok {
		return errtypes.NotFound(n.ParentID)
	}
	path := filepath.Join(parentPath, n.Name)

	// lock the meta file
	unlock, err := t.lookup.MetadataBackend().Lock(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = unlock()
	}()

	if err := os.MkdirAll(path, 0700); err != nil {
		return errors.Wrap(err, "Decomposedfs: error creating node")
	}

	_ = idcache.Set(ctx, n.SpaceID, n.ID, path)

	attributes := n.NodeMetadata(ctx)
	attributes[prefixes.IDAttr] = []byte(n.ID)
	attributes[prefixes.TreesizeAttr] = []byte("0") // initialize as empty, TODO why bother? if it is not set we could treat it as 0?
	if t.options.TreeTimeAccounting || t.options.TreeSizeAccounting {
		attributes[prefixes.PropagationAttr] = []byte("1") // mark the node for propagation
	}
	return n.SetXattrsWithContext(ctx, attributes, false)
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

func isLockFile(path string) bool {
	return strings.HasSuffix(path, ".lock") || strings.HasSuffix(path, ".flock") || strings.HasSuffix(path, ".mlock")
}

func isTrash(path string) bool {
	return strings.HasSuffix(path, ".trashinfo") || strings.HasSuffix(path, ".trashitem")
}
