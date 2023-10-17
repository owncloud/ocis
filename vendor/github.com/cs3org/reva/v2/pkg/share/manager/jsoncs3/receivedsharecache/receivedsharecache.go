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
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
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
	Spaces map[string]*Space

	etag string
}

// Space holds the received shares of one user in one space
type Space struct {
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
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Grab lock")
	unlock := c.lockUser(userID)
	span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))
	defer unlock()

	if c.ReceivedSpaces[userID] == nil {
		err := c.syncWithLock(ctx, userID)
		if err != nil {
			return err
		}
	}

	ctx, span = appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Add")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID), attribute.String("cs3.spaceid", spaceID))

	persistFunc := func() error {
		c.initializeIfNeeded(userID, spaceID)

		receivedSpace := c.ReceivedSpaces[userID].Spaces[spaceID]
		if receivedSpace.States == nil {
			receivedSpace.States = map[string]*State{}
		}
		receivedSpace.States[rs.Share.Id.GetOpaqueId()] = &State{
			State:      rs.State,
			MountPoint: rs.MountPoint,
		}

		return c.persist(ctx, userID)
	}

	log := appctx.GetLogger(ctx).With().
		Str("hostname", os.Getenv("HOSTNAME")).
		Str("userID", userID).
		Str("spaceID", spaceID).Logger()

	var err error
	for retries := 100; retries > 0; retries-- {
		err = persistFunc()
		switch err.(type) {
		case nil:
			span.SetStatus(codes.Ok, "")
			return nil
		case errtypes.Aborted:
			log.Debug().Msg("aborted when persisting added received share: etag changed. retrying...")
			// this is the expected status code from the server when the if-match etag check fails
			// continue with sync below
		case errtypes.PreconditionFailed:
			log.Debug().Msg("precondition failed when persisting added received share: etag changed. retrying...")
			// actually, this is the wrong status code and we treat it like errtypes.Aborted because of inconsistencies on the server side
			// continue with sync below
		default:
			span.SetStatus(codes.Error, fmt.Sprintf("persisting added received share failed. giving up: %s", err.Error()))
			log.Error().Err(err).Msg("persisting added received share failed")
			return err
		}
		if err := c.syncWithLock(ctx, userID); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error().Err(err).Msg("persisting added received share failed. giving up.")
			return err
		}
	}
	return err
}

// Get returns one entry from the cache
func (c *Cache) Get(ctx context.Context, userID, spaceID, shareID string) (*State, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Grab lock")
	unlock := c.lockUser(userID)
	span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))
	defer unlock()

	err := c.syncWithLock(ctx, userID)
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
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Grab lock")
	unlock := c.lockUser(userID)
	span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))
	defer unlock()

	return c.syncWithLock(ctx, userID)
}

func (c *Cache) syncWithLock(ctx context.Context, userID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Sync")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))

	log := appctx.GetLogger(ctx).With().Str("userID", userID).Logger()

	c.initializeIfNeeded(userID, "")

	jsonPath := userJSONPath(userID)
	span.AddEvent("updating cache")
	//  - update cached list of created shares for the user in memory if changed
	dlres, err := c.storage.Download(ctx, metadata.DownloadRequest{
		Path:        jsonPath,
		IfNoneMatch: []string{c.ReceivedSpaces[userID].etag},
	})
	switch err.(type) {
	case nil:
		span.AddEvent("updating local cache")
	case errtypes.NotFound:
		span.SetStatus(codes.Ok, "")
		return nil
	case errtypes.NotModified:
		span.SetStatus(codes.Ok, "")
		return nil
	default:
		span.SetStatus(codes.Error, fmt.Sprintf("Failed to download the received share: %s", err.Error()))
		log.Error().Err(err).Msg("Failed to download the received share")
		return err
	}

	newSpaces := &Spaces{}
	err = json.Unmarshal(dlres.Content, newSpaces)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("Failed to unmarshal the received share: %s", err.Error()))
		log.Error().Err(err).Msg("Failed to unmarshal the received share")
		return err
	}
	newSpaces.etag = dlres.Etag

	c.ReceivedSpaces[userID] = newSpaces
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

	createdBytes, err := json.Marshal(c.ReceivedSpaces[userID])
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	jsonPath := userJSONPath(userID)
	if err := c.storage.MakeDirIfNotExist(ctx, path.Dir(jsonPath)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	ur := metadata.UploadRequest{
		Path:        jsonPath,
		Content:     createdBytes,
		IfMatchEtag: c.ReceivedSpaces[userID].etag,
	}
	// when there is no etag in memory make sure the file has not been created on the server, see https://www.rfc-editor.org/rfc/rfc9110#field.if-match
	// > If the field value is "*", the condition is false if the origin server has a current representation for the target resource.
	if c.ReceivedSpaces[userID].etag == "" {
		ur.IfNoneMatch = []string{"*"}
	}

	_, err = c.storage.Upload(ctx, ur)
	if err != nil {
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

func (c *Cache) initializeIfNeeded(userID, spaceID string) {
	if c.ReceivedSpaces[userID] == nil {
		c.ReceivedSpaces[userID] = &Spaces{
			Spaces: map[string]*Space{},
		}
	}
	if spaceID != "" && c.ReceivedSpaces[userID].Spaces[spaceID] == nil {
		c.ReceivedSpaces[userID].Spaces[spaceID] = &Space{}
	}
}
