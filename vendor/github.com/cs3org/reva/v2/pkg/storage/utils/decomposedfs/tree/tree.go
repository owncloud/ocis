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
	iofs "io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/logger"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/xattrs"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

//go:generate make --no-print-directory -C ../../../../.. mockery NAME=Blobstore

// Blobstore defines an interface for storing blobs in a blobstore
type Blobstore interface {
	Upload(node *node.Node, reader io.Reader) error
	Download(node *node.Node) (io.ReadCloser, error)
	Delete(node *node.Node) error
}

// PathLookup defines the interface for the lookup component
type PathLookup interface {
	NodeFromResource(ctx context.Context, ref *provider.Reference) (*node.Node, error)
	NodeFromID(ctx context.Context, id *provider.ResourceId) (n *node.Node, err error)

	InternalRoot() string
	InternalPath(spaceID, nodeID string) string
	Path(ctx context.Context, n *node.Node) (path string, err error)
	ShareFolder() string
}

// Tree manages a hierarchical tree
type Tree struct {
	lookup    PathLookup
	blobstore Blobstore

	root               string
	treeSizeAccounting bool
	treeTimeAccounting bool
}

// PermissionCheckFunc defined a function used to check resource permissions
type PermissionCheckFunc func(rp *provider.ResourcePermissions) bool

// New returns a new instance of Tree
func New(root string, tta bool, tsa bool, lu PathLookup, bs Blobstore) *Tree {
	return &Tree{
		lookup:             lu,
		blobstore:          bs,
		root:               root,
		treeTimeAccounting: tta,
		treeSizeAccounting: tsa,
	}
}

// Setup prepares the tree structure
func (t *Tree) Setup() error {
	// create data paths for internal layout
	dataPaths := []string{
		filepath.Join(t.root, "spaces"),
		// notes contain symlinks from nodes/<u-u-i-d>/uploads/<uploadid> to ../../uploads/<uploadid>
		// better to keep uploads on a fast / volatile storage before a workflow finally moves them to the nodes dir
		filepath.Join(t.root, "uploads"),
	}
	for _, v := range dataPaths {
		err := os.MkdirAll(v, 0700)
		if err != nil {
			return err
		}
	}

	// create spaces folder and iterate over existing nodes to populate it
	nodesPath := filepath.Join(t.root, "nodes")
	fi, err := os.Stat(nodesPath)
	if err == nil && fi.IsDir() {

		f, err := os.Open(nodesPath)
		if err != nil {
			return err
		}
		nodes, err := f.Readdir(0)
		if err != nil {
			return err
		}

		for _, node := range nodes {
			nodePath := filepath.Join(nodesPath, node.Name())

			if isRootNode(nodePath) {
				if err := t.moveNode(node.Name(), node.Name()); err != nil {
					logger.New().Error().Err(err).
						Str("space", node.Name()).
						Msg("could not move space")
					continue
				}
				t.linkSpace("personal", node.Name())
			}
		}
		// TODO delete nodesPath if empty

	}

	return nil
}
func (t *Tree) moveNode(spaceID, nodeID string) error {
	dirPath := filepath.Join(t.root, "nodes", nodeID)
	f, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	children, err := f.Readdir(0)
	if err != nil {
		return err
	}
	for _, child := range children {
		old := filepath.Join(t.root, "nodes", child.Name())
		new := filepath.Join(t.root, "spaces", lookup.Pathify(spaceID, 1, 2), "nodes", lookup.Pathify(child.Name(), 4, 2))
		if err := os.Rename(old, new); err != nil {
			logger.New().Error().Err(err).
				Str("space", spaceID).
				Str("nodes", child.Name()).
				Str("oldpath", old).
				Str("newpath", new).
				Msg("could not rename node")
		}
		if child.IsDir() {
			if err := t.moveNode(spaceID, child.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

// linkSpace creates a new symbolic link for a space with the given type st, and node id
func (t *Tree) linkSpace(spaceType, spaceID string) {
	spaceTypesPath := filepath.Join(t.root, "spacetypes", spaceType, spaceID)
	expectedTarget := "../../spaces/" + lookup.Pathify(spaceID, 1, 2) + "/nodes/" + lookup.Pathify(spaceID, 4, 2)
	linkTarget, err := os.Readlink(spaceTypesPath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Symlink(expectedTarget, spaceTypesPath)
		if err != nil {
			logger.New().Error().Err(err).
				Str("space_type", spaceType).
				Str("space", spaceID).
				Msg("could not create symlink")
		}
	} else {
		if err != nil {
			logger.New().Error().Err(err).
				Str("space_type", spaceType).
				Str("space", spaceID).
				Msg("could not read symlink")
		}
		if linkTarget != expectedTarget {
			logger.New().Warn().
				Str("space_type", spaceType).
				Str("space", spaceID).
				Str("expected", expectedTarget).
				Str("actual", linkTarget).
				Msg("expected a different link target")
		}
	}
}

// isRootNode checks if a node is a space root
func isRootNode(nodePath string) bool {
	attr, err := xattrs.Get(nodePath, xattrs.ParentidAttr)
	return err == nil && attr == node.RootID
}

/*
func isSharedNode(nodePath string) bool {
	if attrs, err := xattr.List(nodePath); err == nil {
		for i := range attrs {
			if strings.HasPrefix(attrs[i], xattrs.GrantPrefix) {
				return true
			}
		}
	}
	return false
}
*/

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
func (t *Tree) TouchFile(ctx context.Context, n *node.Node) error {
	if n.Exists {
		return errtypes.AlreadyExists(n.ID)
	}

	if n.ID == "" {
		n.ID = uuid.New().String()
	}

	nodePath := n.InternalPath()
	if err := os.MkdirAll(filepath.Dir(nodePath), 0700); err != nil {
		return errors.Wrap(err, "Decomposedfs: error creating node")
	}
	_, err := os.Create(nodePath)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: error creating node")
	}

	err = n.WriteAllNodeMetadata()
	if err != nil {
		return err
	}

	// link child name to parent if it is new
	childNameLink := filepath.Join(n.ParentInternalPath(), n.Name)
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

	return t.Propagate(ctx, n)
}

// CreateDir creates a new directory entry in the tree
func (t *Tree) CreateDir(ctx context.Context, n *node.Node) (err error) {

	if n.Exists {
		return errtypes.AlreadyExists(n.ID) // path?
	}

	// create a directory node
	if n.ID == "" {
		n.ID = uuid.New().String()
	}

	err = t.createNode(n)
	if err != nil {
		return
	}

	// make child appear in listings
	relativeNodePath := filepath.Join("../../../../../", lookup.Pathify(n.ID, 4, 2))
	err = os.Symlink(relativeNodePath, filepath.Join(n.ParentInternalPath(), n.Name))
	if err != nil {
		// no better way to check unfortunately
		if !strings.Contains(err.Error(), "file exists") {
			return
		}

		// try to remove the node
		e := os.RemoveAll(n.InternalPath())
		if e != nil {
			appctx.GetLogger(ctx).Debug().Err(e).Msg("cannot delete node")
		}
		return errtypes.AlreadyExists(err.Error())
	}
	return t.Propagate(ctx, n)
}

// Move replaces the target with the source
func (t *Tree) Move(ctx context.Context, oldNode *node.Node, newNode *node.Node) (err error) {
	// if target exists delete it without trashing it
	if newNode.Exists {
		// TODO make sure all children are deleted
		if err := os.RemoveAll(newNode.InternalPath()); err != nil {
			return errors.Wrap(err, "Decomposedfs: Move: error deleting target node "+newNode.ID)
		}
	}

	// Always target the old node ID for xattr updates.
	// The new node id is empty if the target does not exist
	// and we need to overwrite the new one when overwriting an existing path.
	tgtPath := oldNode.InternalPath()

	// are we just renaming (parent stays the same)?
	if oldNode.ParentID == newNode.ParentID {

		// parentPath := t.lookup.InternalPath(oldNode.SpaceID, oldNode.ParentID)
		parentPath := oldNode.ParentInternalPath()

		// rename child
		err = os.Rename(
			filepath.Join(parentPath, oldNode.Name),
			filepath.Join(parentPath, newNode.Name),
		)
		if err != nil {
			return errors.Wrap(err, "Decomposedfs: could not rename child")
		}

		// update name attribute
		if err := xattrs.Set(tgtPath, xattrs.NameAttr, newNode.Name); err != nil {
			return errors.Wrap(err, "Decomposedfs: could not set name attribute")
		}

		return t.Propagate(ctx, newNode)
	}

	// we are moving the node to a new parent, any target has been removed
	// bring old node to the new parent

	// rename child
	err = os.Rename(
		filepath.Join(oldNode.ParentInternalPath(), oldNode.Name),
		filepath.Join(newNode.ParentInternalPath(), newNode.Name),
	)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: could not move child")
	}

	// update target parentid and name
	if err := xattrs.Set(tgtPath, xattrs.ParentidAttr, newNode.ParentID); err != nil {
		return errors.Wrap(err, "Decomposedfs: could not set parentid attribute")
	}
	if err := xattrs.Set(tgtPath, xattrs.NameAttr, newNode.Name); err != nil {
		return errors.Wrap(err, "Decomposedfs: could not set name attribute")
	}

	// TODO inefficient because we might update several nodes twice, only propagate unchanged nodes?
	// collect in a list, then only stat each node once
	// also do this in a go routine ... webdav should check the etag async

	err = t.Propagate(ctx, oldNode)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: Move: could not propagate old node")
	}
	err = t.Propagate(ctx, newNode)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: Move: could not propagate new node")
	}
	return nil
}

func readChildNodeFromLink(path string) (string, error) {
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
	dir := n.InternalPath()
	f, err := os.Open(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errtypes.NotFound(dir)
		}
		return nil, errors.Wrap(err, "tree: error listing "+dir)
	}
	defer f.Close()

	names, err := f.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	nodes := []*node.Node{}
	for i := range names {
		nodeID, err := readChildNodeFromLink(filepath.Join(dir, names[i]))
		if err != nil {
			// TODO log
			continue
		}

		child, err := node.ReadNode(ctx, t.lookup, n.SpaceID, nodeID, false)
		if err != nil {
			// TODO log
			continue
		}
		if child.SpaceRoot == nil {
			child.SpaceRoot = n.SpaceRoot
		}
		nodes = append(nodes, child)
	}
	return nodes, nil
}

// Delete deletes a node in the tree by moving it to the trash
func (t *Tree) Delete(ctx context.Context, n *node.Node) (err error) {
	deletingSharedResource := ctx.Value(appctx.DeletingSharedResource)

	if deletingSharedResource != nil && deletingSharedResource.(bool) {
		src := filepath.Join(n.ParentInternalPath(), n.Name)
		return os.Remove(src)
	}

	// get the original path
	origin, err := t.lookup.Path(ctx, n)
	if err != nil {
		return
	}

	// set origin location in metadata
	nodePath := n.InternalPath()
	if err := n.SetMetadata(xattrs.TrashOriginAttr, origin); err != nil {
		return err
	}

	deletionTime := time.Now().UTC().Format(time.RFC3339Nano)

	// Prepare the trash
	trashLink := filepath.Join(t.root, "spaces", lookup.Pathify(n.SpaceRoot.ID, 1, 2), "trash", lookup.Pathify(n.ID, 4, 2))
	if err := os.MkdirAll(filepath.Dir(trashLink), 0700); err != nil {
		// Roll back changes
		_ = n.RemoveMetadata(xattrs.TrashOriginAttr)
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
		_ = n.RemoveMetadata(xattrs.TrashOriginAttr)
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
		_ = n.RemoveMetadata(xattrs.TrashOriginAttr)
		return
	}

	// Remove lock file if it exists
	_ = os.Remove(n.LockFilePath())

	// finally remove the entry from the parent dir
	src := filepath.Join(n.ParentInternalPath(), n.Name)
	err = os.Remove(src)
	if err != nil {
		// To roll back changes
		// TODO revert the rename
		// TODO remove symlink
		// Roll back changes
		_ = n.RemoveMetadata(xattrs.TrashOriginAttr)
		return
	}

	return t.Propagate(ctx, n)
}

// RestoreRecycleItemFunc returns a node and a function to restore it from the trash.
func (t *Tree) RestoreRecycleItemFunc(ctx context.Context, spaceid, key, trashPath string, targetNode *node.Node) (*node.Node, *node.Node, func() error, error) {
	recycleNode, trashItem, deletedNodePath, origin, err := t.readRecycleItem(ctx, spaceid, key, trashPath)
	if err != nil {
		return nil, nil, nil, err
	}

	targetRef := &provider.Reference{
		ResourceId: &provider.ResourceId{StorageId: spaceid, OpaqueId: spaceid},
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

	parent, err := targetNode.Parent()
	if err != nil {
		return nil, nil, nil, err
	}

	fn := func() error {
		if targetNode.Exists {
			return errtypes.AlreadyExists("origin already exists")
		}

		// add the entry for the parent dir
		err = os.Symlink("../../../../../"+lookup.Pathify(recycleNode.ID, 4, 2), filepath.Join(targetNode.ParentInternalPath(), targetNode.Name))
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
		}

		targetNode.Exists = true
		// update name attribute
		if err := recycleNode.SetMetadata(xattrs.NameAttr, targetNode.Name); err != nil {
			return errors.Wrap(err, "Decomposedfs: could not set name attribute")
		}

		// set ParentidAttr to restorePath's node parent id
		if trashPath != "" {
			if err := recycleNode.SetMetadata(xattrs.ParentidAttr, targetNode.ParentID); err != nil {
				return errors.Wrap(err, "Decomposedfs: could not set name attribute")
			}
		}

		// delete item link in trash
		if err = os.Remove(trashItem); err != nil {
			log.Error().Err(err).Str("trashItem", trashItem).Msg("error deleting trash item")
		}
		return t.Propagate(ctx, targetNode)
	}
	return recycleNode, parent, fn, nil
}

// PurgeRecycleItemFunc returns a node and a function to purge it from the trash
func (t *Tree) PurgeRecycleItemFunc(ctx context.Context, spaceid, key string, path string) (*node.Node, func() error, error) {
	rn, trashItem, deletedNodePath, _, err := t.readRecycleItem(ctx, spaceid, key, path)
	if err != nil {
		return nil, nil, err
	}

	// only the root node is trashed, the rest is still in normal file system
	children, err := os.ReadDir(deletedNodePath)
	var nodes []*node.Node
	for _, c := range children {
		n, _, _, _, err := t.readRecycleItem(ctx, spaceid, key, filepath.Join(path, c.Name()))
		if err != nil {
			return nil, nil, err
		}
		nodes, err = appendChildren(ctx, n, nodes)
		if err != nil {
			return nil, nil, err
		}
	}

	fn := func() error {
		if err := t.removeNode(deletedNodePath, rn); err != nil {
			return err
		}

		// delete item link in trash
		if err = os.Remove(trashItem); err != nil {
			log.Error().Err(err).Str("trashItem", trashItem).Msg("error deleting trash item")
			return err
		}

		// delete children
		for i := len(nodes) - 1; i >= 0; i-- {
			n := nodes[i]
			if err := t.removeNode(n.InternalPath(), n); err != nil {
				return err
			}

		}

		return nil
	}

	return rn, fn, nil
}

func (t *Tree) removeNode(path string, n *node.Node) error {
	// delete the actual node
	if err := utils.RemoveItem(path); err != nil {
		log.Error().Err(err).Str("path", path).Msg("error node")
		return err
	}

	// delete blob from blobstore
	if n.BlobID != "" {
		if err := t.DeleteBlob(n); err != nil {
			log.Error().Err(err).Str("blobID", n.BlobID).Msg("error deleting nodes blob")
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
		bID, err := node.ReadBlobIDAttr(rev)
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
func (t *Tree) Propagate(ctx context.Context, n *node.Node) (err error) {
	sublog := appctx.GetLogger(ctx).With().Interface("node", n).Logger()
	if !t.treeTimeAccounting && !t.treeSizeAccounting {
		// no propagation enabled
		sublog.Debug().Msg("propagation disabled")
		return
	}

	// is propagation enabled for the parent node?
	root := n.SpaceRoot

	// use a sync time and don't rely on the mtime of the current node, as the stat might not change when a rename happened too quickly
	sTime := time.Now().UTC()

	// we loop until we reach the root
	for err == nil && n.ID != root.ID {
		sublog.Debug().Msg("propagating")

		if n, err = n.Parent(); err != nil {
			break
		}

		sublog = sublog.With().Interface("node", n).Logger()

		// TODO none, sync and async?
		if !n.HasPropagation() {
			sublog.Debug().Str("attr", xattrs.PropagationAttr).Msg("propagation attribute not set or unreadable, not propagating")
			// if the attribute is not set treat it as false / none / no propagation
			return nil
		}

		if t.treeTimeAccounting {
			// update the parent tree time if it is older than the nodes mtime
			updateSyncTime := false

			var tmTime time.Time
			tmTime, err = n.GetTMTime()
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
				if err = n.SetTMTime(&sTime); err != nil {
					sublog.Error().Err(err).Time("tmtime", sTime).Msg("could not update tmtime of parent node")
				} else {
					sublog.Debug().Time("tmtime", sTime).Msg("updated tmtime of parent node")
				}
			}

			if err := n.UnsetTempEtag(); err != nil {
				sublog.Error().Err(err).Msg("could not remove temporary etag attribute")
			}
		}

		// size accounting
		if t.treeSizeAccounting {
			// update the treesize if it differs from the current size
			updateTreeSize := false

			var treeSize, calculatedTreeSize uint64
			calculatedTreeSize, err = calculateTreeSize(ctx, n.InternalPath())
			if err != nil {
				continue
			}

			treeSize, err = n.GetTreeSize()
			switch {
			case err != nil:
				// missing attribute, or invalid format, overwrite
				sublog.Debug().Err(err).Msg("could not read treesize attribute, overwriting")
				updateTreeSize = true
			case treeSize != calculatedTreeSize:
				sublog.Debug().
					Uint64("treesize", treeSize).
					Uint64("calculatedTreeSize", calculatedTreeSize).
					Msg("parent treesize is different then calculated treesize, updating")
				updateTreeSize = true
			default:
				sublog.Debug().
					Uint64("treesize", treeSize).
					Uint64("calculatedTreeSize", calculatedTreeSize).
					Msg("parent size matches calculated size, not updating")
			}

			if updateTreeSize {
				// update the tree time of the parent node
				if err = n.SetTreeSize(calculatedTreeSize); err != nil {
					sublog.Error().Err(err).Uint64("calculatedTreeSize", calculatedTreeSize).Msg("could not update treesize of parent node")
				} else {
					sublog.Debug().Uint64("calculatedTreeSize", calculatedTreeSize).Msg("updated treesize of parent node")
				}
			}
		}
	}
	if err != nil {
		sublog.Error().Err(err).Msg("error propagating")
		return
	}
	return
}

func calculateTreeSize(ctx context.Context, nodePath string) (uint64, error) {
	var size uint64

	f, err := os.Open(nodePath)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Str("nodepath", nodePath).Msg("could not open dir")
		return 0, err
	}
	defer f.Close()

	names, err := f.Readdirnames(0)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Str("nodepath", nodePath).Msg("could not read dirnames")
		return 0, err
	}
	for i := range names {
		cPath := filepath.Join(nodePath, names[i])
		info, err := os.Stat(cPath)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("childpath", cPath).Msg("could not stat child entry")
			continue // continue after an error
		}
		if !info.IsDir() {
			blobSize, err := node.ReadBlobSizeAttr(cPath)
			if err != nil {
				appctx.GetLogger(ctx).Error().Err(err).Str("childpath", cPath).Msg("could not read blobSize xattr")
				continue // continue after an error
			}
			size += uint64(blobSize)
		} else {
			// read from attr
			var b string
			// xattrs.Get will follow the symlink
			if b, err = xattrs.Get(cPath, xattrs.TreesizeAttr); err != nil {
				// TODO recursively descend and recalculate treesize
				continue // continue after an error
			}
			csize, err := strconv.ParseUint(b, 10, 64)
			if err != nil {
				// TODO recursively descend and recalculate treesize
				continue // continue after an error
			}
			size += csize
		}
	}
	return size, err

}

// WriteBlob writes a blob to the blobstore
func (t *Tree) WriteBlob(node *node.Node, reader io.Reader) error {
	return t.blobstore.Upload(node, reader)
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
	if node.BlobID == "" {
		return fmt.Errorf("could not delete blob, node with empty blob id was given")
	}

	return t.blobstore.Delete(node)
}

// TODO check if node exists?
func (t *Tree) createNode(n *node.Node) (err error) {
	// create a directory node
	nodePath := n.InternalPath()
	if err = os.MkdirAll(nodePath, 0700); err != nil {
		return errors.Wrap(err, "Decomposedfs: error creating node")
	}

	return n.WriteAllNodeMetadata()
}

// readTrashLink returns nodeID and timestamp
func readTrashLink(path string) (string, string, error) {
	link, err := os.Readlink(path)
	if err != nil {
		return "", "", err
	}
	// ../../../../../nodes/e5/6c/75/a8/-d235-4cbb-8b4e-48b6fd0f2094.T.2022-02-16T14:38:11.769917408Z
	// TODO use filepath.Separator to support windows
	link = strings.ReplaceAll(link, "/", "")
	// ..........nodese56c75a8-d235-4cbb-8b4e-48b6fd0f2094.T.2022-02-16T14:38:11.769917408Z
	if link[0:15] != "..........nodes" || link[51:54] != ".T." {
		return "", "", errtypes.InternalError("malformed trash link")
	}
	return link[15:51], link[54:], nil
}

// readTrashChildLink returns nodeID
func readTrashChildLink(path string) (string, error) {
	link, err := os.Readlink(path)
	if err != nil {
		return "", err
	}
	// ../../../../../e5/6c/75/a8/-d235-4cbb-8b4e-48b6fd0f2094
	// TODO use filepath.Separator to support windows
	link = strings.ReplaceAll(link, "/", "")
	// ..........e56c75a8-d235-4cbb-8b4e-48b6fd0f2094
	if link[0:10] != ".........." {
		return "", errtypes.InternalError("malformed trash child link")
	}
	return link[10:], nil
}

// TODO refactor the returned params into Node properties? would make all the path transformations go away...
func (t *Tree) readRecycleItem(ctx context.Context, spaceID, key, path string) (recycleNode *node.Node, trashItem string, deletedNodePath string, origin string, err error) {
	if key == "" {
		return nil, "", "", "", errtypes.InternalError("key is empty")
	}

	var nodeID, timeSuffix string

	trashItem = filepath.Join(t.lookup.InternalRoot(), "spaces", lookup.Pathify(spaceID, 1, 2), "trash", lookup.Pathify(key, 4, 2), path)
	if path == "" || path == "/" {
		nodeID, timeSuffix, err = readTrashLink(trashItem)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("trashItem", trashItem).Msg("error reading trash link")
			return
		}
		deletedNodePath = filepath.Join(t.lookup.InternalPath(spaceID, nodeID) + node.TrashIDDelimiter + timeSuffix)
	} else {
		// children of a trashed node are in the nodes folder
		nodeID, err = readTrashChildLink(trashItem)
		if err != nil {
			appctx.GetLogger(ctx).Error().Err(err).Str("trashItem", trashItem).Msg("error reading trash child link")
			return
		}
		deletedNodePath = t.lookup.InternalPath(spaceID, nodeID)
	}

	recycleNode = node.New(spaceID, nodeID, "", "", 0, "", nil, t.lookup)
	recycleNode.SpaceRoot, err = node.ReadNode(ctx, t.lookup, spaceID, spaceID, false)
	if err != nil {
		return
	}

	var attrStr string
	// lookup blobID in extended attributes
	if attrStr, err = xattrs.Get(deletedNodePath, xattrs.BlobIDAttr); err == nil {
		recycleNode.BlobID = attrStr
	} else {
		return
	}

	// lookup parent id in extended attributes
	if attrStr, err = xattrs.Get(deletedNodePath, xattrs.ParentidAttr); err == nil {
		recycleNode.ParentID = attrStr
	} else {
		return
	}

	// lookup name in extended attributes
	if attrStr, err = xattrs.Get(deletedNodePath, xattrs.NameAttr); err == nil {
		recycleNode.Name = attrStr
	} else {
		return
	}

	// get origin node, is relative to space root
	origin = "/"

	trashRootItemPath := filepath.Join(t.lookup.InternalRoot(), "spaces", lookup.Pathify(spaceID, 1, 2), "trash", lookup.Pathify(key, 4, 2))
	// lookup origin path in extended attributes
	if attrStr, err = xattrs.Get(trashRootItemPath, xattrs.TrashOriginAttr); err == nil {
		origin = filepath.Join(attrStr, path)
	} else {
		log.Error().Err(err).Str("trashItem", trashItem).Str("deletedNodePath", deletedNodePath).Msg("could not read origin path, restoring to /")
	}

	return
}

// appendChildren appends `n` and all its children to `nodes`
func appendChildren(ctx context.Context, n *node.Node, nodes []*node.Node) ([]*node.Node, error) {
	nodes = append(nodes, n)

	children, err := os.ReadDir(n.InternalPath())
	if err != nil {
		// TODO: How to differentiate folders from files?
		return nodes, nil
	}

	for _, c := range children {
		cn, err := n.Child(ctx, c.Name())
		if err != nil {
			// continue?
			return nil, err
		}
		nodes, err = appendChildren(ctx, cn, nodes)
		if err != nil {
			// continue?
			return nil, err
		}
	}

	return nodes, nil
}
