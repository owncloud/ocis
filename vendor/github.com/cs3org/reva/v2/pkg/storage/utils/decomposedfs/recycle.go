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
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storagespace"
)

// Recycle items are stored inside the node folder and start with the uuid of the deleted node.
// The `.T.` indicates it is a trash item and what follows is the timestamp of the deletion.
// The deleted file is kept in the same location/dir as the original node. This prevents deletes
// from triggering cross storage moves when the trash is accidentally stored on another partition,
// because the admin mounted a different partition there.
// For an efficient listing of deleted nodes the ocis storage driver maintains a 'trash' folder
// with symlinks to trash files for every storagespace.

// ListRecycle returns the list of available recycle items
// ref -> the space (= resourceid), key -> deleted node id, relativePath = relative to key
func (fs *Decomposedfs) ListRecycle(ctx context.Context, ref *provider.Reference, key, relativePath string) ([]*provider.RecycleItem, error) {

	if ref == nil || ref.ResourceId == nil || ref.ResourceId.OpaqueId == "" {
		return nil, errtypes.BadRequest("spaceid required")
	}
	if key == "" && relativePath != "" {
		return nil, errtypes.BadRequest("key is required when navigating with a path")
	}
	spaceID := ref.ResourceId.OpaqueId

	sublog := appctx.GetLogger(ctx).With().Str("spaceid", spaceID).Str("key", key).Str("relative_path", relativePath).Logger()

	// check permissions
	trashnode, err := fs.lu.NodeFromSpaceID(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	rp, err := fs.p.AssembleTrashPermissions(ctx, trashnode)
	switch {
	case err != nil:
		return nil, err
	case !rp.ListRecycle:
		if rp.Stat {
			return nil, errtypes.PermissionDenied(key)
		}
		return nil, errtypes.NotFound(key)
	}

	if key == "" && relativePath == "" {
		return fs.listTrashRoot(ctx, spaceID)
	}

	// build a list of trash items relative to the given trash root and path
	items := make([]*provider.RecycleItem, 0)

	trashRootPath := filepath.Join(fs.getRecycleRoot(spaceID), lookup.Pathify(key, 4, 2))
	originalPath, _, timeSuffix, err := readTrashLink(trashRootPath)
	if err != nil {
		sublog.Error().Err(err).Str("trashRoot", trashRootPath).Msg("error reading trash link")
		return nil, err
	}

	origin := ""
	attrs, err := fs.lu.MetadataBackend().All(ctx, originalPath)
	if err != nil {
		return items, err
	}
	// lookup origin path in extended attributes
	if attrBytes, ok := attrs[prefixes.TrashOriginAttr]; ok {
		origin = string(attrBytes)
	} else {
		sublog.Error().Err(err).Str("spaceid", spaceID).Msg("could not read origin path, skipping")
		return nil, err
	}

	// all deleted items have the same deletion time
	var deletionTime *types.Timestamp
	if parsed, err := time.Parse(time.RFC3339Nano, timeSuffix); err == nil {
		deletionTime = &types.Timestamp{
			Seconds: uint64(parsed.Unix()),
			// TODO nanos
		}
	} else {
		sublog.Error().Err(err).Msg("could not parse time format, ignoring")
	}

	var size int64
	if relativePath == "" {
		// this is the case when we want to directly list a file in the trashbin
		nodeType := fs.lu.TypeFromPath(ctx, originalPath)
		switch nodeType {
		case provider.ResourceType_RESOURCE_TYPE_FILE:
			size, err = fs.lu.ReadBlobSizeAttr(ctx, originalPath)
			if err != nil {
				return items, err
			}
		case provider.ResourceType_RESOURCE_TYPE_CONTAINER:
			size, err = fs.lu.MetadataBackend().GetInt64(ctx, originalPath, prefixes.TreesizeAttr)
			if err != nil {
				return items, err
			}
		}
		item := &provider.RecycleItem{
			Type:         fs.lu.TypeFromPath(ctx, originalPath),
			Size:         uint64(size),
			Key:          filepath.Join(key, relativePath),
			DeletionTime: deletionTime,
			Ref: &provider.Reference{
				Path: filepath.Join(origin, relativePath),
			},
		}
		items = append(items, item)
		return items, err
	}

	// we have to read the names and stat the path to follow the symlinks
	childrenPath := filepath.Join(originalPath, relativePath)
	childrenDir, err := os.Open(childrenPath)
	if err != nil {
		return nil, err
	}

	names, err := childrenDir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		resolvedChildPath, err := filepath.EvalSymlinks(filepath.Join(childrenPath, name))
		if err != nil {
			sublog.Error().Err(err).Str("name", name).Msg("could not resolve symlink, skipping")
			continue
		}

		// reset size
		size = 0

		nodeType := fs.lu.TypeFromPath(ctx, resolvedChildPath)
		switch nodeType {
		case provider.ResourceType_RESOURCE_TYPE_FILE:
			size, err = fs.lu.ReadBlobSizeAttr(ctx, resolvedChildPath)
			if err != nil {
				sublog.Error().Err(err).Str("name", name).Msg("invalid blob size, skipping")
				continue
			}
		case provider.ResourceType_RESOURCE_TYPE_CONTAINER:
			size, err = fs.lu.MetadataBackend().GetInt64(ctx, resolvedChildPath, prefixes.TreesizeAttr)
			if err != nil {
				sublog.Error().Err(err).Str("name", name).Msg("invalid tree size, skipping")
				continue
			}
		case provider.ResourceType_RESOURCE_TYPE_INVALID:
			sublog.Error().Err(err).Str("name", name).Str("resolvedChildPath", resolvedChildPath).Msg("invalid node type, skipping")
			continue
		}

		item := &provider.RecycleItem{
			Type:         nodeType,
			Size:         uint64(size),
			Key:          filepath.Join(key, relativePath, name),
			DeletionTime: deletionTime,
			Ref: &provider.Reference{
				Path: filepath.Join(origin, relativePath, name),
			},
		}
		items = append(items, item)
	}
	return items, nil
}

// readTrashLink returns path, nodeID and timestamp
func readTrashLink(path string) (string, string, string, error) {
	link, err := os.Readlink(path)
	if err != nil {
		return "", "", "", err
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", "", "", err
	}
	// ../../../../../nodes/e5/6c/75/a8/-d235-4cbb-8b4e-48b6fd0f2094.T.2022-02-16T14:38:11.769917408Z
	// TODO use filepath.Separator to support windows
	link = strings.ReplaceAll(link, "/", "")
	// ..........nodese56c75a8-d235-4cbb-8b4e-48b6fd0f2094.T.2022-02-16T14:38:11.769917408Z
	if link[0:15] != "..........nodes" || link[51:54] != node.TrashIDDelimiter {
		return "", "", "", errtypes.InternalError("malformed trash link")
	}
	return resolved, link[15:51], link[54:], nil
}

func (fs *Decomposedfs) listTrashRoot(ctx context.Context, spaceID string) ([]*provider.RecycleItem, error) {
	log := appctx.GetLogger(ctx)
	trashRoot := fs.getRecycleRoot(spaceID)
	items := []*provider.RecycleItem{}
	subTrees, err := filepath.Glob(trashRoot + "/*")
	if err != nil {
		return nil, err
	}

	numWorkers := fs.o.MaxConcurrency
	if len(subTrees) < numWorkers {
		numWorkers = len(subTrees)
	}

	work := make(chan string, len(subTrees))
	results := make(chan *provider.RecycleItem, len(subTrees))

	g, ctx := errgroup.WithContext(ctx)

	// Distribute work
	g.Go(func() error {
		defer close(work)
		for _, itemPath := range subTrees {
			select {
			case work <- itemPath:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	// Spawn workers that'll concurrently work the queue
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for subTree := range work {
				matches, err := filepath.Glob(subTree + "/*/*/*/*")
				if err != nil {
					return err
				}

				for _, itemPath := range matches {
					// TODO can we encode this in the path instead of reading the link?
					nodePath, nodeID, timeSuffix, err := readTrashLink(itemPath)
					if err != nil {
						log.Error().Err(err).Str("trashRoot", trashRoot).Str("item", itemPath).Msg("error reading trash link, skipping")
						continue
					}

					md, err := os.Stat(nodePath)
					if err != nil {
						log.Error().Err(err).Str("trashRoot", trashRoot).Str("item", itemPath).Str("node_path", nodePath).Msg("could not stat trash item, skipping")
						continue
					}

					attrs, err := fs.lu.MetadataBackend().All(ctx, nodePath)
					if err != nil {
						log.Error().Err(err).Str("trashRoot", trashRoot).Str("item", itemPath).Str("node_path", nodePath).Msg("could not get extended attributes, skipping")
						continue
					}

					nodeType := fs.lu.TypeFromPath(ctx, nodePath)
					if nodeType == provider.ResourceType_RESOURCE_TYPE_INVALID {
						log.Error().Err(err).Str("trashRoot", trashRoot).Str("item", itemPath).Str("node_path", nodePath).Msg("invalid node type, skipping")
						continue
					}

					item := &provider.RecycleItem{
						Type: nodeType,
						Size: uint64(md.Size()),
						Key:  nodeID,
					}
					if deletionTime, err := time.Parse(time.RFC3339Nano, timeSuffix); err == nil {
						item.DeletionTime = &types.Timestamp{
							Seconds: uint64(deletionTime.Unix()),
							// TODO nanos
						}
					} else {
						log.Error().Err(err).Str("trashRoot", trashRoot).Str("item", itemPath).Str("spaceid", spaceID).Str("nodeid", nodeID).Str("dtime", timeSuffix).Msg("could not parse time format, ignoring")
					}

					// lookup origin path in extended attributes
					if attr, ok := attrs[prefixes.TrashOriginAttr]; ok {
						item.Ref = &provider.Reference{Path: string(attr)}
					} else {
						log.Error().Str("trashRoot", trashRoot).Str("item", itemPath).Str("spaceid", spaceID).Str("nodeid", nodeID).Str("dtime", timeSuffix).Msg("could not read origin path")
					}

					select {
					case results <- item:
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

	// Collect results
	for ri := range results {
		items = append(items, ri)
	}
	return items, nil
}

// RestoreRecycleItem restores the specified item
func (fs *Decomposedfs) RestoreRecycleItem(ctx context.Context, ref *provider.Reference, key, relativePath string, restoreRef *provider.Reference) error {
	if ref == nil {
		return errtypes.BadRequest("missing reference, needs a space id")
	}

	var targetNode *node.Node
	if restoreRef != nil {
		tn, err := fs.lu.NodeFromResource(ctx, restoreRef)
		if err != nil {
			return err
		}

		targetNode = tn
	}

	rn, parent, restoreFunc, err := fs.tp.RestoreRecycleItemFunc(ctx, ref.ResourceId.SpaceId, key, relativePath, targetNode)
	if err != nil {
		return err
	}

	// check permissions of deleted node
	rp, err := fs.p.AssembleTrashPermissions(ctx, rn)
	switch {
	case err != nil:
		return err
	case !rp.RestoreRecycleItem:
		if rp.Stat {
			return errtypes.PermissionDenied(key)
		}
		return errtypes.NotFound(key)
	}

	// Set space owner in context
	storagespace.ContextSendSpaceOwnerID(ctx, rn.SpaceOwnerOrManager(ctx))

	// check we can write to the parent of the restore reference
	pp, err := fs.p.AssemblePermissions(ctx, parent)
	switch {
	case err != nil:
		return err
	case !pp.InitiateFileUpload:
		// share receiver cannot restore to a shared resource to which she does not have write permissions.
		if rp.Stat {
			return errtypes.PermissionDenied(key)
		}
		return errtypes.NotFound(key)
	}

	// Run the restore func
	return restoreFunc()
}

// PurgeRecycleItem purges the specified item, all its children and all their revisions
func (fs *Decomposedfs) PurgeRecycleItem(ctx context.Context, ref *provider.Reference, key, relativePath string) error {
	if ref == nil {
		return errtypes.BadRequest("missing reference, needs a space id")
	}

	rn, purgeFunc, err := fs.tp.PurgeRecycleItemFunc(ctx, ref.ResourceId.OpaqueId, key, relativePath)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return errtypes.NotFound(key)
		}
		return err
	}

	// check permissions of deleted node
	rp, err := fs.p.AssembleTrashPermissions(ctx, rn)
	switch {
	case err != nil:
		return err
	case !rp.PurgeRecycle:
		if rp.Stat {
			return errtypes.PermissionDenied(key)
		}
		return errtypes.NotFound(key)
	}

	// Run the purge func
	return purgeFunc()
}

// EmptyRecycle empties the trash
func (fs *Decomposedfs) EmptyRecycle(ctx context.Context, ref *provider.Reference) error {
	if ref == nil || ref.ResourceId == nil || ref.ResourceId.OpaqueId == "" {
		return errtypes.BadRequest("spaceid must be set")
	}

	items, err := fs.ListRecycle(ctx, ref, "", "")
	if err != nil {
		return err
	}

	for _, i := range items {
		if err := fs.PurgeRecycleItem(ctx, ref, i.Key, ""); err != nil {
			return err
		}
	}
	// TODO what permission should we check? we could check the root node of the user? or the owner permissions on his home root node?
	// The current impl will wipe your own trash. or when no user provided the trash of 'root'
	return os.RemoveAll(fs.getRecycleRoot(ref.ResourceId.SpaceId))
}

func (fs *Decomposedfs) getRecycleRoot(spaceID string) string {
	return filepath.Join(fs.getSpaceRoot(spaceID), "trash")
}
