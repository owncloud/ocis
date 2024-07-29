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
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/utils"
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
)

type queueItem struct {
	item  scanItem
	timer *time.Timer
}

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
					err = t.WarmupIDCache(item.Path, true)
					if err != nil {
						log.Error().Err(err).Str("path", item.Path).Msg("failed to warmup id cache")
					}
				}
			}
		}()
	}
}

// Scan scans the given path and updates the id chache
func (t *Tree) Scan(path string, action EventAction, isDir bool, recurse bool) error {
	// cases:
	switch action {
	case ActionCreate:
		if !isDir {
			// 1. New file (could be emitted as part of a new directory)
			//	 -> assimilate file
			//   -> scan parent directory recursively
			if !t.scanDebouncer.InProgress(filepath.Dir(path)) {
				t.scanDebouncer.Debounce(scanItem{
					Path:        path,
					ForceRescan: false,
				})
			}
			t.scanDebouncer.Debounce(scanItem{
				Path:        filepath.Dir(path),
				ForceRescan: true,
				Recurse:     true,
			})
		} else {
			// 2. New directory
			//  -> scan directory
			t.scanDebouncer.Debounce(scanItem{
				Path:        path,
				ForceRescan: true,
				Recurse:     true,
			})
		}

	case ActionUpdate:
		// 3. Updated file
		//   -> update file unless parent directory is being rescanned
		if !t.scanDebouncer.InProgress(filepath.Dir(path)) {
			t.scanDebouncer.Debounce(scanItem{
				Path:        path,
				ForceRescan: true,
			})
		}

	case ActionMove:
		// 4. Moved file
		//   -> update file
		// 5. Moved directory
		//   -> update directory and all children
		t.scanDebouncer.Debounce(scanItem{
			Path:        path,
			ForceRescan: isDir,
			Recurse:     isDir,
		})

	case ActionDelete:
		_ = t.HandleFileDelete(path)
	}

	return nil
}

func (t *Tree) HandleFileDelete(path string) error {
	// purge metadata
	_ = t.lookup.(*lookup.Lookup).IDCache.DeleteByPath(context.Background(), path)
	_ = t.lookup.MetadataBackend().Purge(path)

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

		// This item had already been assimilated in the past. Update the path
		_ = t.lookup.(*lookup.Lookup).CacheID(context.Background(), spaceID, string(id), item.Path)

		previousParentID, _ := t.lookup.MetadataBackend().Get(context.Background(), item.Path, prefixes.ParentidAttr)

		fi, err := t.updateFile(item.Path, string(id), spaceID)
		if err != nil {
			return err
		}

		// was it moved?
		if ok && len(previousParentID) > 0 && previousPath != item.Path {
			// purge original metadata. Only delete the path entry using DeletePath(reverse lookup), not the whole entry pair.
			_ = t.lookup.(*lookup.Lookup).IDCache.DeletePath(context.Background(), previousPath)
			_ = t.lookup.MetadataBackend().Purge(previousPath)

			if fi.IsDir() {
				// if it was moved and it is a directory we need to propagate the move
				go func() { _ = t.WarmupIDCache(item.Path, false) }()
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
			// }
		}
	} else {
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
	prevMtime, err := previousAttribs.Time(prefixes.MTimeAttr)
	if err != nil || prevMtime.Before(fi.ModTime()) {
		attributes[prefixes.MTimeAttr] = []byte(fi.ModTime().Format(time.RFC3339Nano))
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
		attributes.SetString(prefixes.BlobIDAttr, id)
		attributes.SetInt64(prefixes.BlobsizeAttr, fi.Size())

		// propagate the change
		sizeDiff := fi.Size()
		if previousAttribs != nil && previousAttribs[prefixes.BlobsizeAttr] != nil {
			oldSize, err := attributes.Int64(prefixes.BlobsizeAttr)
			if err == nil {
				sizeDiff -= oldSize
			}
		}

		n := node.New(spaceID, id, parentID, filepath.Base(path), fi.Size(), "", provider.ResourceType_RESOURCE_TYPE_FILE, nil, t.lookup)
		n.SpaceRoot = &node.Node{SpaceID: spaceID, ID: spaceID}
		err = t.Propagate(context.Background(), n, sizeDiff)
		if err != nil {
			return nil, errors.Wrap(err, "failed to propagate")
		}
	}
	err = t.lookup.MetadataBackend().SetMultiple(context.Background(), path, attributes, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set attributes")
	}

	_ = t.lookup.(*lookup.Lookup).CacheID(context.Background(), spaceID, id, path)

	return fi, nil
}

// WarmupIDCache warms up the id cache
func (t *Tree) WarmupIDCache(root string, assimilate bool) error {
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
		if err != nil {
			return err
		}

		// skip lock files
		if isLockFile(path) {
			return nil
		}

		// calculate tree sizes
		if !info.IsDir() {
			dir := filepath.Dir(path)
			sizes[dir] += info.Size()
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
				_ = t.lookup.(*lookup.Lookup).CacheID(context.Background(), string(spaceID), string(id), path)
			}
		} else if assimilate {
			_ = t.assimilate(scanItem{Path: path, ForceRescan: true})
		}
		return nil
	})
	if err != nil {
		return err
	}

	for dir, size := range sizes {
		_ = t.lookup.MetadataBackend().Set(context.Background(), dir, prefixes.TreesizeAttr, []byte(fmt.Sprintf("%d", size)))
	}
	return nil
}
