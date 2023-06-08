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

package sharecache

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"time"

	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/share/manager/jsoncs3/shareid"
	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
	"github.com/cs3org/reva/v2/pkg/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// name is the Tracer name used to identify this instrumentation library.
const tracerName = "sharecache"

// Cache caches the list of share ids for users/groups
// It functions as an in-memory cache with a persistence layer
// The storage is sharded by user/group
type Cache struct {
	UserShares map[string]*UserShareCache

	storage   metadata.Storage
	namespace string
	filename  string
	ttl       time.Duration
}

// UserShareCache holds the space/share map for one user
type UserShareCache struct {
	Mtime      time.Time
	UserShares map[string]*SpaceShareIDs

	nextSync time.Time
}

// SpaceShareIDs holds the unique list of share ids for a space
type SpaceShareIDs struct {
	Mtime time.Time
	IDs   map[string]struct{}
}

// New returns a new Cache instance
func New(s metadata.Storage, namespace, filename string, ttl time.Duration) Cache {
	return Cache{
		UserShares: map[string]*UserShareCache{},
		storage:    s,
		namespace:  namespace,
		filename:   filename,
		ttl:        ttl,
	}
}

// Add adds a share to the cache
func (c *Cache) Add(ctx context.Context, userid, shareID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Add")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userid), attribute.String("cs3.shareid", shareID))

	storageid, spaceid, _ := shareid.Decode(shareID)
	ssid := storageid + shareid.IDDelimiter + spaceid

	now := time.Now()
	if c.UserShares[userid] == nil {
		c.UserShares[userid] = &UserShareCache{
			UserShares: map[string]*SpaceShareIDs{},
		}
	}
	if c.UserShares[userid].UserShares[ssid] == nil {
		c.UserShares[userid].UserShares[ssid] = &SpaceShareIDs{
			IDs: map[string]struct{}{},
		}
	}
	// add share id
	c.UserShares[userid].Mtime = now
	c.UserShares[userid].UserShares[ssid].Mtime = now
	c.UserShares[userid].UserShares[ssid].IDs[shareID] = struct{}{}
	return c.Persist(ctx, userid)
}

// Remove removes a share for the given user
func (c *Cache) Remove(ctx context.Context, userid, shareID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Remove")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userid), attribute.String("cs3.shareid", shareID))

	storageid, spaceid, _ := shareid.Decode(shareID)
	ssid := storageid + shareid.IDDelimiter + spaceid

	now := time.Now()
	if c.UserShares[userid] == nil {
		c.UserShares[userid] = &UserShareCache{
			Mtime:      now,
			UserShares: map[string]*SpaceShareIDs{},
		}
	}

	if c.UserShares[userid].UserShares[ssid] != nil {
		// remove share id
		c.UserShares[userid].Mtime = now
		c.UserShares[userid].UserShares[ssid].Mtime = now
		delete(c.UserShares[userid].UserShares[ssid].IDs, shareID)
	}

	return c.Persist(ctx, userid)
}

// List return the list of spaces/shares for the given user/group
func (c *Cache) List(userid string) map[string]SpaceShareIDs {
	r := map[string]SpaceShareIDs{}
	if c.UserShares[userid] == nil {
		return r
	}

	for ssid, cached := range c.UserShares[userid].UserShares {
		r[ssid] = SpaceShareIDs{
			Mtime: cached.Mtime,
			IDs:   cached.IDs,
		}
	}
	return r
}

// Sync updates the in-memory data with the data from the storage if it is outdated
func (c *Cache) Sync(ctx context.Context, userID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Sync")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))

	log := appctx.GetLogger(ctx).With().Str("userID", userID).Logger()

	var mtime time.Time
	//  - do we have a cached list of created shares for the user in memory?
	if usc := c.UserShares[userID]; usc != nil {
		if time.Now().Before(c.UserShares[userID].nextSync) {
			span.AddEvent("skip sync")
			span.SetStatus(codes.Ok, "")
			return nil
		}
		c.UserShares[userID].nextSync = time.Now().Add(c.ttl)

		mtime = usc.Mtime
		//    - y: set If-Modified-Since header to only download if it changed
	} else {
		mtime = time.Time{} // Set zero time so that data from storage always takes precedence
	}

	userCreatedPath := c.userCreatedPath(userID)
	info, err := c.storage.Stat(ctx, userCreatedPath)
	if err != nil {
		if _, ok := err.(errtypes.NotFound); ok {
			span.AddEvent("no file")
			span.SetStatus(codes.Ok, "")
			return nil // Nothing to sync against
		}
		span.SetStatus(codes.Error, fmt.Sprintf("Failed to stat the share cache: %s", err.Error()))
		log.Error().Err(err).Msg("Failed to stat the share cache")
		return err
	}
	// check mtime of /users/{userid}/created.json
	if utils.TSToTime(info.Mtime).After(mtime) {
		span.AddEvent("updating cache")
		//  - update cached list of created shares for the user in memory if changed
		createdBlob, err := c.storage.SimpleDownload(ctx, userCreatedPath)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("Failed to download the share cache: %s", err.Error()))
			log.Error().Err(err).Msg("Failed to download the share cache")
			return err
		}
		newShareCache := &UserShareCache{}
		err = json.Unmarshal(createdBlob, newShareCache)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("Failed to unmarshal the share cache: %s", err.Error()))
			log.Error().Err(err).Msg("Failed to unmarshal the share cache")
			return err
		}
		newShareCache.Mtime = utils.TSToTime(info.Mtime)
		c.UserShares[userID] = newShareCache
	}
	span.SetStatus(codes.Ok, "")
	return nil
}

// Persist persists the data for one user/group to the storage
func (c *Cache) Persist(ctx context.Context, userid string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Persist")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userid))

	oldMtime := c.UserShares[userid].Mtime
	c.UserShares[userid].Mtime = time.Now()

	createdBytes, err := json.Marshal(c.UserShares[userid])
	if err != nil {
		c.UserShares[userid].Mtime = oldMtime
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	jsonPath := c.userCreatedPath(userid)
	if err := c.storage.MakeDirIfNotExist(ctx, path.Dir(jsonPath)); err != nil {
		c.UserShares[userid].Mtime = oldMtime
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err = c.storage.Upload(ctx, metadata.UploadRequest{
		Path:              jsonPath,
		Content:           createdBytes,
		IfUnmodifiedSince: c.UserShares[userid].Mtime,
	}); err != nil {
		c.UserShares[userid].Mtime = oldMtime
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetStatus(codes.Ok, "")
	return nil
}

func (c *Cache) userCreatedPath(userid string) string {
	return filepath.Join("/", c.namespace, userid, c.filename)
}
