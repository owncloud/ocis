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

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ScanDebouncer struct {
	after      time.Duration
	f          func(item scanItem)
	pending    map[string]*time.Timer
	inProgress sync.Map

	mutex sync.Mutex
}

// NewScanDebouncer returns a new SpaceDebouncer instance
func NewScanDebouncer(d time.Duration, f func(item scanItem)) *ScanDebouncer {
	return &ScanDebouncer{
		after:      d,
		f:          f,
		pending:    map[string]*time.Timer{},
		inProgress: sync.Map{},
	}
}

// Debounce restars the debounce timer for the given space
func (d *ScanDebouncer) Debounce(item scanItem) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	path := item.Path
	force := item.ForceRescan
	if t := d.pending[item.Path]; t != nil {
		force = force || item.ForceRescan
		t.Stop()
	}

	d.pending[item.Path] = time.AfterFunc(d.after, func() {
		if _, ok := d.inProgress.Load(path); ok {
			// Reschedule this run for when the previous run has finished
			d.mutex.Lock()
			d.pending[path].Reset(d.after)
			d.mutex.Unlock()
			return
		}

		d.inProgress.Store(path, true)
		defer d.inProgress.Delete(path)
		d.f(scanItem{
			Path:        path,
			ForceRescan: force,
		})
	})
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
			}
		}()
	}
}

// Scan scans the given path and updates the id chache
func (t *Tree) Scan(path string, forceRescan bool) error {
	t.scanDebouncer.Debounce(scanItem{
		Path:        path,
		ForceRescan: forceRescan,
	})
	return nil
}

func (t *Tree) assimilate(item scanItem) error {
	var err error
	// find the space id, scope by the according user
	spaceID := []byte("")
	spaceCandidate := item.Path
	for strings.HasPrefix(spaceCandidate, t.options.Root) {
		spaceID, err = t.lookup.MetadataBackend().Get(context.Background(), spaceCandidate, prefixes.SpaceIDAttr)
		if err == nil {
			if t.options.UseSpaceGroups {
				// set the uid and gid for the space
				fi, err := os.Stat(spaceCandidate)
				if err != nil {
					return err
				}
				sys := fi.Sys().(*syscall.Stat_t)
				gid := int(sys.Gid)
				_, err = t.userMapper.ScopeUserByIds(-1, gid)
				if err != nil {
					return err
				}
			}
			break
		}
		spaceCandidate = filepath.Dir(spaceCandidate)
	}
	if len(spaceID) == 0 {
		return fmt.Errorf("did not find space id for path")
	}

	var id []byte
	if !item.ForceRescan {
		// already assimilated?
		id, err := t.lookup.MetadataBackend().Get(context.Background(), item.Path, prefixes.IDAttr)
		if err == nil {
			_ = t.lookup.(*lookup.Lookup).CacheID(context.Background(), string(spaceID), string(id), item.Path)
			return nil
		}
	}

	// lock the file for assimilation
	unlock, err := t.lookup.MetadataBackend().Lock(item.Path)
	if err != nil {
		return errors.Wrap(err, "failed to lock item for assimilation")
	}
	defer func() {
		_ = unlock()
	}()

	// check for the id attribute again after grabbing the lock, maybe the file was assimilated/created by us in the meantime
	id, err = t.lookup.MetadataBackend().Get(context.Background(), item.Path, prefixes.IDAttr)
	if err == nil {
		_ = t.lookup.(*lookup.Lookup).CacheID(context.Background(), string(spaceID), string(id), item.Path)
		if item.ForceRescan {
			_, err = t.updateFile(item.Path, string(id), string(spaceID))
			if err != nil {
				return err
			}
		}
	} else {
		// assimilate new file
		newId := uuid.New().String()
		_, err = t.updateFile(item.Path, newId, string(spaceID))
		if err != nil {
			return err
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
