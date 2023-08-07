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

package receivedsharecache

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"sync"
	"time"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
	"github.com/cs3org/reva/v2/pkg/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// name is the Tracer name used to identify this instrumentation library.
const tracerName = "receivedsharecache"

// Cache stores the list of received shares and their states
// It functions as an in-memory cache with a persistence layer
// The storage is sharded by user
type Cache struct {
	lockMap sync.Map

	ReceivedSpaces map[string]*Spaces

	storage metadata.Storage
	ttl     time.Duration
}

// Spaces holds the received shares of one user per space
type Spaces struct {
	Mtime  time.Time
	Spaces map[string]*Space

	nextSync time.Time
}

// Space holds the received shares of one user in one space
type Space struct {
	Mtime  time.Time
	States map[string]*State
}

// State holds the state information of a received share
type State struct {
	State      collaboration.ShareState
	MountPoint *provider.Reference
}

// New returns a new Cache instance
func New(s metadata.Storage, ttl time.Duration) Cache {
	return Cache{
		ReceivedSpaces: map[string]*Spaces{},
		storage:        s,
		ttl:            ttl,
		lockMap:        sync.Map{},
	}
}

func (c *Cache) lockUser(userID string) func() {
	v, _ := c.lockMap.LoadOrStore(userID, &sync.Mutex{})
	lock := v.(*sync.Mutex)

	lock.Lock()
	return func() { lock.Unlock() }
}

// Add adds a new entry to the cache
func (c *Cache) Add(ctx context.Context, userID, spaceID string, rs *collaboration.ReceivedShare) error {
	unlock := c.lockUser(userID)
	defer unlock()

	if c.ReceivedSpaces[userID] == nil {
		err := c.syncWithLock(ctx, userID)
		if err != nil {
			return err
		}
	}

	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Add")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID), attribute.String("cs3.spaceid", spaceID))

	persistFunc := func() error {
		if c.ReceivedSpaces[userID] == nil {
			c.ReceivedSpaces[userID] = &Spaces{
				Spaces: map[string]*Space{},
			}
		}
		if c.ReceivedSpaces[userID].Spaces[spaceID] == nil {
			c.ReceivedSpaces[userID].Spaces[spaceID] = &Space{}
		}

		receivedSpace := c.ReceivedSpaces[userID].Spaces[spaceID]
		receivedSpace.Mtime = time.Now()
		if receivedSpace.States == nil {
			receivedSpace.States = map[string]*State{}
		}
		receivedSpace.States[rs.Share.Id.GetOpaqueId()] = &State{
			State:      rs.State,
			MountPoint: rs.MountPoint,
		}

		return c.persist(ctx, userID)
	}
	err := persistFunc()
	if _, ok := err.(errtypes.IsPreconditionFailed); ok {
		if err := c.syncWithLock(ctx, userID); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		err = persistFunc()
	}
	return err
}

// Get returns one entry from the cache
func (c *Cache) Get(ctx context.Context, userID, spaceID, shareID string) (*State, error) {
	err := c.Sync(ctx, userID)
	if err != nil {
		return nil, err
	}
	if c.ReceivedSpaces[userID] == nil || c.ReceivedSpaces[userID].Spaces[spaceID] == nil {
		return nil, nil
	}
	return c.ReceivedSpaces[userID].Spaces[spaceID].States[shareID], nil
}

// Sync updates the in-memory data with the data from the storage if it is outdated
func (c *Cache) Sync(ctx context.Context, userID string) error {
	unlock := c.lockUser(userID)
	defer unlock()

	return c.syncWithLock(ctx, userID)
}

func (c *Cache) syncWithLock(ctx context.Context, userID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Sync")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))

	log := appctx.GetLogger(ctx).With().Str("userID", userID).Logger()

	var mtime time.Time
	if c.ReceivedSpaces[userID] != nil {
		if time.Now().Before(c.ReceivedSpaces[userID].nextSync) {
			span.AddEvent("skip sync")
			span.SetStatus(codes.Ok, "")
			return nil
		}
		c.ReceivedSpaces[userID].nextSync = time.Now().Add(c.ttl)

		mtime = c.ReceivedSpaces[userID].Mtime
	} else {
		mtime = time.Time{} // Set zero time so that data from storage always takes precedence
	}

	jsonPath := userJSONPath(userID)
	info, err := c.storage.Stat(ctx, jsonPath) // TODO we only need the mtime ... use fieldmask to make the request cheaper
	if err != nil {
		if _, ok := err.(errtypes.NotFound); ok {
			span.AddEvent("no file")
			span.SetStatus(codes.Ok, "")
			return nil // Nothing to sync against
		}
		span.SetStatus(codes.Error, fmt.Sprintf("Failed to stat the received share: %s", err.Error()))
		log.Error().Err(err).Msg("Failed to stat the received share")
		return err
	}
	// check mtime of /users/{userid}/created.json
	if utils.TSToTime(info.Mtime).After(mtime) {
		span.AddEvent("updating cache")
		//  - update cached list of created shares for the user in memory if changed
		createdBlob, err := c.storage.SimpleDownload(ctx, jsonPath)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("Failed to download the received share: %s", err.Error()))
			log.Error().Err(err).Msg("Failed to download the received share")
			return err
		}
		newSpaces := &Spaces{}
		err = json.Unmarshal(createdBlob, newSpaces)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("Failed to unmarshal the received share: %s", err.Error()))
			log.Error().Err(err).Msg("Failed to unmarshal the received share")
			return err
		}
		newSpaces.Mtime = utils.TSToTime(info.Mtime)
		c.ReceivedSpaces[userID] = newSpaces
	}
	span.SetStatus(codes.Ok, "")
	return nil
}

// persist persists the data for one user to the storage
func (c *Cache) persist(ctx context.Context, userID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Persist")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))

	if c.ReceivedSpaces[userID] == nil {
		span.SetStatus(codes.Ok, "no received shares")
		return nil
	}

	oldMtime := c.ReceivedSpaces[userID].Mtime
	c.ReceivedSpaces[userID].Mtime = time.Now()

	createdBytes, err := json.Marshal(c.ReceivedSpaces[userID])
	if err != nil {
		c.ReceivedSpaces[userID].Mtime = oldMtime
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	jsonPath := userJSONPath(userID)
	if err := c.storage.MakeDirIfNotExist(ctx, path.Dir(jsonPath)); err != nil {
		c.ReceivedSpaces[userID].Mtime = oldMtime
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err = c.storage.Upload(ctx, metadata.UploadRequest{
		Path:              jsonPath,
		Content:           createdBytes,
		IfUnmodifiedSince: oldMtime,
		MTime:             c.ReceivedSpaces[userID].Mtime,
	}); err != nil {
		c.ReceivedSpaces[userID].Mtime = oldMtime
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetStatus(codes.Ok, "")
	return nil
}

func userJSONPath(userID string) string {
	return filepath.Join("/users", userID, "received.json")
}
