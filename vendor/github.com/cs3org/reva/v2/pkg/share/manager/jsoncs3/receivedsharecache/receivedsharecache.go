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
	"path"
	"path/filepath"
	"time"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
	"github.com/cs3org/reva/v2/pkg/utils"
)

// Cache stores the list of received shares and their states
// It functions as an in-memory cache with a persistence layer
// The storage is sharded by user
type Cache struct {
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
	}
}

// Add adds a new entry to the cache
func (c *Cache) Add(ctx context.Context, userID, spaceID string, rs *collaboration.ReceivedShare) error {
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

	return c.Persist(ctx, userID)
}

// Get returns one entry from the cache
func (c *Cache) Get(userID, spaceID, shareID string) *State {
	if c.ReceivedSpaces[userID] == nil || c.ReceivedSpaces[userID].Spaces[spaceID] == nil {
		return nil
	}
	return c.ReceivedSpaces[userID].Spaces[spaceID].States[shareID]
}

// Sync updates the in-memory data with the data from the storage if it is outdated
func (c *Cache) Sync(ctx context.Context, userID string) error {
	log := appctx.GetLogger(ctx).With().Str("userID", userID).Logger()
	log.Debug().Msg("Syncing received share cache...")

	var mtime time.Time
	if c.ReceivedSpaces[userID] != nil {
		if time.Now().Before(c.ReceivedSpaces[userID].nextSync) {
			log.Debug().Msg("Skipping received share cache sync, it was just recently synced...")
			return nil
		}
		c.ReceivedSpaces[userID].nextSync = time.Now().Add(c.ttl)

		mtime = c.ReceivedSpaces[userID].Mtime
	} else {
		mtime = time.Time{} // Set zero time so that data from storage always takes precedence
	}

	jsonPath := userJSONPath(userID)
	info, err := c.storage.Stat(ctx, jsonPath)
	if err != nil {
		if _, ok := err.(errtypes.NotFound); ok {
			return nil // Nothing to sync against
		}
		log.Error().Err(err).Msg("Failed to stat the received share")
		return err
	}
	// check mtime of /users/{userid}/created.json
	if utils.TSToTime(info.Mtime).After(mtime) {
		log.Debug().Msg("Updating received share cache...")
		//  - update cached list of created shares for the user in memory if changed
		createdBlob, err := c.storage.SimpleDownload(ctx, jsonPath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to download the received share")
			return err
		}
		newSpaces := &Spaces{}
		err = json.Unmarshal(createdBlob, newSpaces)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal the received share")
			return err
		}
		newSpaces.Mtime = utils.TSToTime(info.Mtime)
		c.ReceivedSpaces[userID] = newSpaces
	}
	log.Debug().Msg("Received share cache is up to date")
	return nil
}

// Persist persists the data for one user to the storage
func (c *Cache) Persist(ctx context.Context, userID string) error {
	if c.ReceivedSpaces[userID] == nil {
		return nil
	}

	oldMtime := c.ReceivedSpaces[userID].Mtime
	c.ReceivedSpaces[userID].Mtime = time.Now()

	createdBytes, err := json.Marshal(c.ReceivedSpaces[userID])
	if err != nil {
		c.ReceivedSpaces[userID].Mtime = oldMtime
		return err
	}
	jsonPath := userJSONPath(userID)
	if err := c.storage.MakeDirIfNotExist(ctx, path.Dir(jsonPath)); err != nil {
		c.ReceivedSpaces[userID].Mtime = oldMtime
		return err
	}

	if err = c.storage.Upload(ctx, metadata.UploadRequest{
		Path:              jsonPath,
		Content:           createdBytes,
		IfUnmodifiedSince: c.ReceivedSpaces[userID].Mtime,
	}); err != nil {
		c.ReceivedSpaces[userID].Mtime = oldMtime
		return err
	}
	return nil
}

func userJSONPath(userID string) string {
	return filepath.Join("/users", userID, "received.json")
}
