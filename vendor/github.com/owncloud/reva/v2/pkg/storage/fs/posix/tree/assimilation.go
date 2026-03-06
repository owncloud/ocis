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
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/storage/fs/posix/lookup"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/owncloud/reva/v2/pkg/utils"
)

type ScanDebouncer struct {
	after      time.Duration
	f          func(item scanItem)
	pending    sync.Map
	inProgress sync.Map

	mutex sync.Mutex
}

type EventAction int

const (
	ActionCreate EventAction = iota
	ActionUpdate
	ActionMove
	ActionDelete
	ActionMoveFrom
)

type queueItem struct {
	item  scanItem
	timer *time.Timer
}

const dirtyFlag = "user.ocis.dirty"

// NewScanDebouncer returns a new SpaceDebouncer instance
func NewScanDebouncer(d time.Duration, f func(item scanItem)) *ScanDebouncer {
	return &ScanDebouncer{
		after:      d,
		f:          f,
		pending:    sync.Map{},
		inProgress: sync.Map{},
	}
}

// Debounce restarts the debounce timer for the given space
func (d *ScanDebouncer) Debounce(item scanItem) {
	if d.after == 0 {
		d.f(item)
		return
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	path := item.Path
	force := item.ForceRescan
	recurse := item.Recurse
	if i, ok := d.pending.Load(item.Path); ok {
		queueItem := i.(*queueItem)
		force = force || queueItem.item.ForceRescan
		recurse = recurse || queueItem.item.Recurse
		queueItem.timer.Stop()
	}

	d.pending.Store(item.Path, &queueItem{
		item: item,
		timer: time.AfterFunc(d.after, func() {
			if _, ok := d.inProgress.Load(path); ok {
				// Reschedule this run for when the previous run has finished
				d.mutex.Lock()
				if i, ok := d.pending.Load(path); ok {
					i.(*queueItem).timer.Reset(d.after)
				}

				d.mutex.Unlock()
				return
			}

			d.pending.Delete(path)
			d.inProgress.Store(path, true)
			defer d.inProgress.Delete(path)
			d.f(scanItem{
				Path:        path,
				ForceRescan: force,
				Recurse:     recurse,
			})
		}),
	})
}

// InProgress returns true if the given path is currently being processed
func (d *ScanDebouncer) InProgress(path string) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if _, ok := d.pending.Load(path); ok {
		return true
	}

	_, ok := d.inProgress.Load(path)
	return ok
}

func (t *Tree) workScanQueue() {
	for i := 0; i < t.options.MaxConcurrency; i++ {
		go func() {
			for {
				item := <-t.scanQueue

				err := t.assimilate(item)
				if err != nil {
					log.Error().Err(err).Str("path", item.Path).Msg("failed to assimilate item")
					continue
				}

				if item.Recurse {
					err = t.WarmupIDCache(item.Path, true, false)
					if err != nil {
						log.Error().Err(err).Str("path", item.Path).Msg("failed to warmup id cache")
					}
				}
			}
		}()
	}
}

// Scan scans the given path and updates the id chache
func (t *Tree) Scan(path string, action EventAction, isDir bool) error {
	// cases:
	switch action {
	case ActionCreate:
		t.log.Debug().Str("path", path).Bool("isDir", isDir).Msg("scanning path (ActionCreate)")
		if !isDir {
			// 1. New file (could be emitted as part of a new directory)
			//	 -> assimilate file
			//   -> scan parent directory recursively to update tree size and catch nodes that weren't covered by an event
			if !t.scanDebouncer.InProgress(filepath.Dir(path)) {
				t.scanDebouncer.Debounce(scanItem{
					Path:        path,
					ForceRescan: false,
				})
			}
			if err := t.setDirty(filepath.Dir(path), true); err != nil {
				return err
			}
			t.scanDebouncer.Debounce(scanItem{
				Path:        filepath.Dir(path),
				ForceRescan: true,
				Recurse:     true,
			})
		} else {
			// 2. New directory
			//  -> scan directory
			if err := t.setDirty(path, true); err != nil {
				return err
			}
			t.scanDebouncer.Debounce(scanItem{
				Path:        path,
				ForceRescan: true,
				Recurse:     true,
			})
		}

	case ActionUpdate:
		t.log.Debug().Str("path", path).Bool("isDir", isDir).Msg("scanning path (ActionUpdate)")
		// 3. Updated file
		//   -> update file unless parent directory is being rescanned
		if !t.scanDebouncer.InProgress(filepath.Dir(path)) {
			t.scanDebouncer.Debounce(scanItem{
				Path:        path,
				ForceRescan: true,
			})
		}

	case ActionMove:
		t.log.Debug().Str("path", path).Bool("isDir", isDir).Msg("scanning path (ActionMove)")
		// 4. Moved file
		//   -> update file
		// 5. Moved directory
		//   -> update directory and all children
		t.scanDebouncer.Debounce(scanItem{
			Path:        path,
			ForceRescan: isDir,
			Recurse:     isDir,
		})

	case ActionMoveFrom:
		t.log.Debug().Str("path", path).Bool("isDir", isDir).Msg("scanning path (ActionMoveFrom)")
		// 6. file/directory moved out of the watched directory
		//   -> update directory
		if err := t.setDirty(filepath.Dir(path), true); err != nil {
			return err
		}

		go func() { _ = t.WarmupIDCache(filepath.Dir(path), false, true) }()

	case ActionDelete:
		t.log.Debug().Str("path", path).Bool("isDir", isDir).Msg("handling deleted item")

		// 7. Deleted file or directory
		//   -> update parent and all children

		err := t.HandleFileDelete(path)
		if err != nil {
			return err
		}

		t.scanDebouncer.Debounce(scanItem{
			Path:        filepath.Dir(path),
			ForceRescan: true,
			Recurse:     true,
		})
	}

	return nil
}

func (t *Tree) HandleFileDelete(path string) error {
	// purge metadata
	if err := t.lookup.(*lookup.Lookup).IDCache.DeleteByPath(context.Background(), path); err != nil {
		t.log.Error().Err(err).Str("path", path).Msg("could not delete id cache entry by path")
	}
	if err := t.lookup.MetadataBackend().Purge(context.Background(), path); err != nil {
		t.log.Error().Err(err).Str("path", path).Msg("could not purge metadata")
	}

	// send event
	owner, spaceID, nodeID, parentID, err := t.getOwnerAndIDs(filepath.Dir(path))
	if err != nil {
		return err
	}

	t.PublishEvent(events.ItemTrashed{
		Owner:     owner,
		Executant: owner,
		Ref: &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: t.options.MountID,
				SpaceId:   spaceID,
				OpaqueId:  parentID,
			},
			Path: filepath.Base(path),
		},
		ID: &provider.ResourceId{
			StorageId: t.options.MountID,
			SpaceId:   spaceID,
			OpaqueId:  nodeID,
		},
		Timestamp: utils.TSNow(),
	})

	return nil
}

func (t *Tree) getOwnerAndIDs(path string) (*userv1beta1.UserId, string, string, string, error) {
	lu := t.lookup.(*lookup.Lookup)

	spaceID, nodeID, err := lu.IDsForPath(context.Background(), path)
	if err != nil {
		return nil, "", "", "", err
	}

	attrs, err := t.lookup.MetadataBackend().All(context.Background(), path)
	if err != nil {
		return nil, "", "", "", err
	}

	parentID := string(attrs[prefixes.ParentidAttr])

	spacePath, ok := lu.GetCachedID(context.Background(), spaceID, spaceID)
	if !ok {
		return nil, "", "", "", fmt.Errorf("could not find space root for path %s", path)
	}

	spaceAttrs, err := t.lookup.MetadataBackend().All(context.Background(), spacePath)
	if err != nil {
		return nil, "", "", "", err
	}

	owner := &userv1beta1.UserId{
		Idp:      string(spaceAttrs[prefixes.OwnerIDPAttr]),
		OpaqueId: string(spaceAttrs[prefixes.OwnerIDAttr]),
	}

	return owner, nodeID, spaceID, parentID, nil
}

func (t *Tree) findSpaceId(path string) (string, node.Attributes, error) {
	// find the space id, scope by the according user
	spaceCandidate := path
	spaceAttrs := node.Attributes{}
	for strings.HasPrefix(spaceCandidate, t.options.Root) {
		spaceAttrs, err := t.lookup.MetadataBackend().All(context.Background(), spaceCandidate)
		spaceID := spaceAttrs[prefixes.SpaceIDAttr]
		if err == nil && len(spaceID) > 0 {
			if t.options.UseSpaceGroups {
				// set the uid and gid for the space
				fi, err := os.Stat(spaceCandidate)
				if err != nil {
					return "", spaceAttrs, err
				}
				sys := fi.Sys().(*syscall.Stat_t)
				gid := int(sys.Gid)
				_, err = t.userMapper.ScopeUserByIds(-1, gid)
				if err != nil {
					return "", spaceAttrs, err
				}
			}

			return string(spaceID), spaceAttrs, nil
		}
		spaceCandidate = filepath.Dir(spaceCandidate)
	}
	return "", spaceAttrs, fmt.Errorf("could not find space for path %s", path)
}

func (t *Tree) assimilate(item scanItem) error {
	var id []byte
	var err error

	// First find the space id
	spaceID, spaceAttrs, err := t.findSpaceId(item.Path)
	if err != nil {
		return err
	}

	// lock the file for assimilation
	unlock, err := t.lookup.MetadataBackend().Lock(item.Path)
	if err != nil {
		return errors.Wrap(err, "failed to lock item for assimilation")
	}
	defer func() {
		_ = unlock()
	}()

	user := &userv1beta1.UserId{
		Idp:      string(spaceAttrs[prefixes.OwnerIDPAttr]),
		OpaqueId: string(spaceAttrs[prefixes.OwnerIDAttr]),
	}

	// check for the id attribute again after grabbing the lock, maybe the file was assimilated/created by us in the meantime
	id, err = t.lookup.MetadataBackend().Get(context.Background(), item.Path, prefixes.IDAttr)
	if err == nil {
		previousPath, ok := t.lookup.(*lookup.Lookup).GetCachedID(context.Background(), spaceID, string(id))
		previousParentID, _ := t.lookup.MetadataBackend().Get(context.Background(), item.Path, prefixes.ParentidAttr)

		// was it moved or copied/restored with a clashing id?
		if ok && len(previousParentID) > 0 && previousPath != item.Path {
			_, err := os.Stat(previousPath)
			if err == nil {
				// this id clashes with an existing item -> clear metadata and re-assimilate
				t.log.Debug().Str("path", item.Path).Msg("ID clash detected, purging metadata and re-assimilating")

				if err := t.lookup.MetadataBackend().Purge(context.Background(), item.Path); err != nil {
					t.log.Error().Err(err).Str("path", item.Path).Msg("could not purge metadata")
				}
				go func() {
					if err := t.assimilate(scanItem{Path: item.Path, ForceRescan: true}); err != nil {
						t.log.Error().Err(err).Str("path", item.Path).Msg("could not re-assimilate")
					}
				}()
			} else {
				// this is a move
				t.log.Debug().Str("path", item.Path).Msg("move detected")

				if err := t.lookup.(*lookup.Lookup).CacheID(context.Background(), spaceID, string(id), item.Path); err != nil {
					t.log.Error().Err(err).Str("spaceID", spaceID).Str("id", string(id)).Str("path", item.Path).Msg("could not cache id")
				}
				_, err := t.updateFile(item.Path, string(id), spaceID)
				if err != nil {
					return err
				}

				// purge original metadata. Only delete the path entry using DeletePath(reverse lookup), not the whole entry pair.
				if err := t.lookup.(*lookup.Lookup).IDCache.DeletePath(context.Background(), previousPath); err != nil {
					t.log.Error().Err(err).Str("path", previousPath).Msg("could not delete id cache entry by path")
				}
				if err := t.lookup.MetadataBackend().Purge(context.Background(), previousPath); err != nil {
					t.log.Error().Err(err).Str("path", previousPath).Msg("could not purge metadata")
				}

				fi, err := os.Stat(item.Path)
				if err != nil {
					return err
				}
				if fi.IsDir() {
					// if it was moved and it is a directory we need to propagate the move
					go func() {
						if err := t.WarmupIDCache(item.Path, false, true); err != nil {
							t.log.Error().Err(err).Str("path", item.Path).Msg("could not warmup id cache")
						}
					}()
				}

				parentID, err := t.lookup.MetadataBackend().Get(context.Background(), item.Path, prefixes.ParentidAttr)
				if err == nil && len(parentID) > 0 {
					ref := &provider.Reference{
						ResourceId: &provider.ResourceId{
							StorageId: t.options.MountID,
							SpaceId:   spaceID,
							OpaqueId:  string(parentID),
						},
						Path: filepath.Base(item.Path),
					}
					oldRef := &provider.Reference{
						ResourceId: &provider.ResourceId{
							StorageId: t.options.MountID,
							SpaceId:   spaceID,
							OpaqueId:  string(previousParentID),
						},
						Path: filepath.Base(previousPath),
					}
					t.PublishEvent(events.ItemMoved{
						SpaceOwner:   user,
						Executant:    user,
						Owner:        user,
						Ref:          ref,
						OldReference: oldRef,
						Timestamp:    utils.TSNow(),
					})
				}
			}
		} else {
			// This item had already been assimilated in the past. Update the path
			t.log.Debug().Str("path", item.Path).Msg("updating cached path")
			if err := t.lookup.(*lookup.Lookup).CacheID(context.Background(), spaceID, string(id), item.Path); err != nil {
				t.log.Error().Err(err).Str("spaceID", spaceID).Str("id", string(id)).Str("path", item.Path).Msg("could not cache id")
			}

			_, err := t.updateFile(item.Path, string(id), spaceID)
			if err != nil {
				return err
			}
		}
	} else {
		t.log.Debug().Str("path", item.Path).Msg("new item detected")
		// assimilate new file
		newId := uuid.New().String()
		fi, err := t.updateFile(item.Path, newId, spaceID)
		if err != nil {
			return err
		}

		ref := &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: t.options.MountID,
				SpaceId:   spaceID,
				OpaqueId:  newId,
			},
		}
		if fi.IsDir() {
			t.PublishEvent(events.ContainerCreated{
				SpaceOwner: user,
				Executant:  user,
				Owner:      user,
				Ref:        ref,
				Timestamp:  utils.TSNow(),
			})
		} else {
			if fi.Size() == 0 {
				t.PublishEvent(events.FileTouched{
					SpaceOwner: user,
					Executant:  user,
					Ref:        ref,
					Timestamp:  utils.TSNow(),
				})
			} else {
				t.PublishEvent(events.UploadReady{
					SpaceOwner: user,
					FileRef:    ref,
					ResourceID: ref.ResourceId,
					Timestamp:  utils.TSNow(),
				})
			}
		}
	}
	return nil
}

func (t *Tree) updateFile(path, id, spaceID string) (fs.FileInfo, error) {
	retries := 1
	parentID := ""
assimilate:
	if id != spaceID {
		// read parent
		parentAttribs, err := t.lookup.MetadataBackend().All(context.Background(), filepath.Dir(path))
		if err != nil {
			return nil, fmt.Errorf("failed to read parent item attributes")
		}

		if len(parentAttribs) == 0 || len(parentAttribs[prefixes.IDAttr]) == 0 {
			if retries == 0 {
				return nil, fmt.Errorf("got empty parent attribs even after assimilating")
			}

			// assimilate parent first
			err = t.assimilate(scanItem{Path: filepath.Dir(path), ForceRescan: false})
			if err != nil {
				return nil, err
			}

			// retry
			retries--
			goto assimilate
		}
		parentID = string(parentAttribs[prefixes.IDAttr])
	}

	// assimilate file
	fi, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to stat item")
	}

	attrs, err := t.lookup.MetadataBackend().All(context.Background(), path)
	if err != nil && !metadata.IsAttrUnset(err) {
		return nil, errors.Wrap(err, "failed to get item attribs")
	}
	previousAttribs := node.Attributes(attrs)

	attributes := node.Attributes{
		prefixes.IDAttr:   []byte(id),
		prefixes.NameAttr: []byte(filepath.Base(path)),
	}
	if len(parentID) > 0 {
		attributes[prefixes.ParentidAttr] = []byte(parentID)
	}

	sha1h, md5h, adler32h, err := node.CalculateChecksums(context.Background(), path)
	if err == nil {
		attributes[prefixes.ChecksumPrefix+"sha1"] = sha1h.Sum(nil)
		attributes[prefixes.ChecksumPrefix+"md5"] = md5h.Sum(nil)
		attributes[prefixes.ChecksumPrefix+"adler32"] = adler32h.Sum(nil)
	}

	if fi.IsDir() {
		attributes.SetInt64(prefixes.TypeAttr, int64(provider.ResourceType_RESOURCE_TYPE_CONTAINER))
		attributes.SetInt64(prefixes.TreesizeAttr, 0)
		if previousAttribs != nil && previousAttribs[prefixes.TreesizeAttr] != nil {
			attributes[prefixes.TreesizeAttr] = previousAttribs[prefixes.TreesizeAttr]
		}
		attributes[prefixes.PropagationAttr] = []byte("1")
	} else {
		attributes.SetInt64(prefixes.TypeAttr, int64(provider.ResourceType_RESOURCE_TYPE_FILE))
	}

	n := node.New(spaceID, id, parentID, filepath.Base(path), fi.Size(), "", provider.ResourceType_RESOURCE_TYPE_FILE, nil, t.lookup)
	n.SpaceRoot = &node.Node{SpaceID: spaceID, ID: spaceID}
	err = t.Propagate(context.Background(), n, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to propagate")
	}

	t.log.Debug().Str("path", path).Interface("attributes", attributes).Msg("setting attributes")
	err = t.lookup.MetadataBackend().SetMultiple(context.Background(), path, attributes, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set attributes")
	}

	if err := t.lookup.(*lookup.Lookup).CacheID(context.Background(), spaceID, id, path); err != nil {
		t.log.Error().Err(err).Str("spaceID", spaceID).Str("id", id).Str("path", path).Msg("could not cache id")
	}

	return fi, nil
}

// WarmupIDCache warms up the id cache
func (t *Tree) WarmupIDCache(root string, assimilate, onlyDirty bool) error {
	root = filepath.Clean(root)
	spaceID := []byte("")

	scopeSpace := func(spaceCandidate string) error {
		if !t.options.UseSpaceGroups {
			return nil
		}

		// set the uid and gid for the space
		fi, err := os.Stat(spaceCandidate)
		if err != nil {
			return err
		}
		sys := fi.Sys().(*syscall.Stat_t)
		gid := int(sys.Gid)
		_, err = t.userMapper.ScopeUserByIds(-1, gid)
		return err
	}

	sizes := make(map[string]int64)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// skip lock and upload files
		if isLockFile(path) {
			return nil
		}
		if isTrash(path) || t.isUpload(path) {
			return filepath.SkipDir
		}

		if err != nil {
			return err
		}

		// calculate tree sizes
		if !info.IsDir() {
			dir := path
			for dir != root {
				dir = filepath.Clean(filepath.Dir(dir))
				sizes[dir] += info.Size()
			}
		} else if onlyDirty {
			dirty, err := t.isDirty(path)
			if err != nil {
				return err
			}
			if !dirty {
				return filepath.SkipDir
			}
			sizes[path] += 0 // Make sure to set the size to 0 for empty directories
		}

		attribs, err := t.lookup.MetadataBackend().All(context.Background(), path)
		if err == nil && len(attribs[prefixes.IDAttr]) > 0 {
			nodeSpaceID := attribs[prefixes.SpaceIDAttr]
			if len(nodeSpaceID) > 0 {
				spaceID = nodeSpaceID

				err = scopeSpace(path)
				if err != nil {
					return err
				}
			} else {
				// try to find space
				spaceCandidate := path
				for strings.HasPrefix(spaceCandidate, t.options.Root) {
					spaceID, err = t.lookup.MetadataBackend().Get(context.Background(), spaceCandidate, prefixes.SpaceIDAttr)
					if err == nil {
						err = scopeSpace(path)
						if err != nil {
							return err
						}
						break
					}
					spaceCandidate = filepath.Dir(spaceCandidate)
				}
			}
			if len(spaceID) == 0 {
				return nil // no space found
			}

			id, ok := attribs[prefixes.IDAttr]
			if ok {
				// Check if the item on the previous still exists. In this case it might have been a copy with extended attributes -> set new ID
				previousPath, ok := t.lookup.(*lookup.Lookup).GetCachedID(context.Background(), string(spaceID), string(id))
				if ok && previousPath != path {
					// this id clashes with an existing id -> clear metadata and re-assimilate
					_, err := os.Stat(previousPath)
					if err == nil {
						_ = t.lookup.MetadataBackend().Purge(context.Background(), path)
						_ = t.assimilate(scanItem{Path: path, ForceRescan: true})
					}
				}
				if err := t.lookup.(*lookup.Lookup).CacheID(context.Background(), string(spaceID), string(id), path); err != nil {
					t.log.Error().Err(err).Str("spaceID", string(spaceID)).Str("id", string(id)).Str("path", path).Msg("could not cache id")
				}
			}
		} else if assimilate {
			if err := t.assimilate(scanItem{Path: path, ForceRescan: true}); err != nil {
				t.log.Error().Err(err).Str("path", path).Msg("could not assimilate item")
			}
		}
		return t.setDirty(path, false)
	})

	for dir, size := range sizes {
		if dir == root {
			// Propagate the size diff further up the tree
			if err := t.propagateSizeDiff(dir, size); err != nil {
				t.log.Error().Err(err).Str("path", dir).Msg("could not propagate size diff")
			}
		}
		if err := t.lookup.MetadataBackend().Set(context.Background(), dir, prefixes.TreesizeAttr, []byte(fmt.Sprintf("%d", size))); err != nil {
			t.log.Error().Err(err).Str("path", dir).Int64("size", size).Msg("could not set tree size")
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (t *Tree) propagateSizeDiff(dir string, size int64) error {
	// First find the space id
	spaceID, _, err := t.findSpaceId(dir)
	if err != nil {
		return err
	}
	attrs, err := t.lookup.MetadataBackend().All(context.Background(), dir)
	if err != nil {
		return err
	}
	n, err := t.lookup.NodeFromID(context.Background(), &provider.ResourceId{
		StorageId: t.options.MountID,
		SpaceId:   spaceID,
		OpaqueId:  string(attrs[prefixes.IDAttr]),
	})
	if err != nil {
		return err
	}
	oldSize, err := node.Attributes(attrs).Int64(prefixes.TreesizeAttr)
	if err != nil {
		return err
	}
	return t.Propagate(context.Background(), n, size-oldSize)
}

func (t *Tree) setDirty(path string, dirty bool) error {
	return t.lookup.MetadataBackend().Set(context.Background(), path, dirtyFlag, []byte(fmt.Sprintf("%t", dirty)))
}

func (t *Tree) isDirty(path string) (bool, error) {
	dirtyAttr, err := t.lookup.MetadataBackend().Get(context.Background(), path, dirtyFlag)
	if err != nil {
		return false, err
	}
	return string(dirtyAttr) == "true", nil
}
