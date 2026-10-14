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

	backoff "github.com/cenkalti/backoff/v5"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/appctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/cas"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/mtimesyncedcache"
	"github.com/owncloud/reva/v2/pkg/storage/utils/metadata"
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

	ReceivedSpaces mtimesyncedcache.Map[string, *Spaces]

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
	Hidden     bool
}

// New returns a new Cache instance
func New(s metadata.Storage, ttl time.Duration) Cache {
	return Cache{
		ReceivedSpaces: mtimesyncedcache.Map[string, *Spaces]{},
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
	span.SetAttributes(attribute.String("cs3.userid", userID))
	span.End()
	defer unlock()

	if _, ok := c.ReceivedSpaces.Load(userID); !ok {
		err := c.syncWithRetry(ctx, userID)
		if err != nil {
			return err
		}
	}

	ctx, span = appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Add")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID), attribute.String("cs3.spaceid", spaceID))

	err := c.retryPersist(ctx, userID, spaceID, func() error {
		c.initializeIfNeeded(userID, spaceID)

		rss, _ := c.ReceivedSpaces.Load(userID)
		receivedSpace := rss.Spaces[spaceID]
		if receivedSpace.States == nil {
			receivedSpace.States = map[string]*State{}
		}
		receivedSpace.States[rs.Share.Id.GetOpaqueId()] = &State{
			State:      rs.State,
			MountPoint: rs.MountPoint,
			Hidden:     rs.Hidden,
		}

		return c.persist(ctx, userID)
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	return err
}

// Get returns one entry from the cache
func (c *Cache) Get(ctx context.Context, userID, spaceID, shareID string) (*State, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Grab lock")
	unlock := c.lockUser(userID)
	span.SetAttributes(attribute.String("cs3.userid", userID))
	span.End()
	defer unlock()

	err := c.syncWithRetry(ctx, userID)
	if err != nil {
		return nil, err
	}
	rss, ok := c.ReceivedSpaces.Load(userID)
	if !ok || rss.Spaces[spaceID] == nil {
		return nil, nil
	}
	return rss.Spaces[spaceID].States[shareID], nil
}

// Remove removes an entry from the cache
func (c *Cache) Remove(ctx context.Context, userID, spaceID, shareID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Grab lock")
	unlock := c.lockUser(userID)
	span.SetAttributes(attribute.String("cs3.userid", userID))
	span.End()
	defer unlock()

	ctx, span = appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Remove")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID), attribute.String("cs3.spaceid", spaceID))

	err := c.retryPersist(ctx, userID, spaceID, func() error {
		c.initializeIfNeeded(userID, spaceID)

		rss, _ := c.ReceivedSpaces.Load(userID)
		receivedSpace := rss.Spaces[spaceID]
		if receivedSpace.States == nil {
			receivedSpace.States = map[string]*State{}
		}
		delete(receivedSpace.States, shareID)
		if len(receivedSpace.States) == 0 {
			delete(rss.Spaces, spaceID)
		}

		return c.persist(ctx, userID)
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	return err
}

// List returns a list of received shares for a given user
// The return list is guaranteed to be thread-safe
func (c *Cache) List(ctx context.Context, userID string) (map[string]*Space, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "List")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))

	unlock := c.lockUser(userID)
	defer unlock()

	err := c.syncWithRetry(ctx, userID)
	if err != nil {
		return nil, err
	}

	spaces := map[string]*Space{}
	rss, _ := c.ReceivedSpaces.Load(userID)
	for spaceID, space := range rss.Spaces {
		spaceCopy := &Space{
			States: map[string]*State{},
		}
		for shareID, state := range space.States {
			spaceCopy.States[shareID] = &State{
				State:      state.State,
				MountPoint: state.MountPoint,
				Hidden:     state.Hidden,
			}
		}
		spaces[spaceID] = spaceCopy
	}
	return spaces, nil
}

func (c *Cache) retryPersist(ctx context.Context, userID, spaceID string, persistFunc func() error) error {
	log := appctx.GetLogger(ctx).With().
		Str("hostname", os.Getenv("HOSTNAME")).
		Str("userID", userID).
		Str("spaceID", spaceID).Logger()

	bo := cas.NewBackoff()

	const maxPersistAttempts = 20
	// independent budget, not derived from maxPersistAttempts: a flaky resync must not starve the persist budget.
	const maxResyncAttemptsPerFailure = 20

	var err error
	needsResync := false
	persistAttempts := 0
	resyncAttempts := 0
	for persistAttempts < maxPersistAttempts {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if needsResync {
			// a previous persist attempt failed and the re-read to pick up fresh state
			// was itself transient; keep retrying the re-read instead of hammering
			// persistFunc again with the same stale in-memory rss/etag
			serr := c.resyncAfterConflict(ctx, userID)
			if serr == nil {
				needsResync = false
				resyncAttempts = 0
				continue // fresh state is in memory; retry persistFunc next attempt, no need to wait
			}
			// keep err at its last real (persist or resync) failure; never let a
			// budget-exhausting resync success clobber it into a false nil return
			err = serr
			if !cas.IsSyncTransient(serr) {
				log.Error().Err(serr).Msg("lost update: re-read failed, aborting")
				return serr
			}
			resyncAttempts++
			if resyncAttempts >= maxResyncAttemptsPerFailure {
				log.Error().Err(serr).Msg("lost update: resync would not stabilize, giving up")
				return serr
			}
			log.Warn().Err(serr).Msg("lost update: re-read before retry")
		} else {
			persistAttempts++
			err = persistFunc()
			switch {
			case err == nil:
				return nil
			case cas.IsTransientGRPCStatus(err):
				log.Debug().Err(err).Msg("persist failed: transient gRPC error, retrying")
			case cas.IsConflict(err):
				log.Debug().Err(err).Msg("CAS conflict persisting received share, retrying")
			default:
				log.Error().Err(err).Msg("persisting received share failed, giving up")
				return err
			}
			needsResync = true
		}

		timer := time.NewTimer(bo.NextBackOff())
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}
	return err
}

// syncWithRetry retries sync's cold-start read on transient errors, since
// Add/Get/List call it directly with no other retry wrapper around it.
func (c *Cache) syncWithRetry(ctx context.Context, userID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		err := c.sync(ctx, userID)
		if err != nil && !cas.IsSyncTransient(err) {
			return struct{}{}, backoff.Permanent(err)
		}
		return struct{}{}, err
	}, backoff.WithBackOff(cas.NewBackoff()), backoff.WithMaxTries(20))
	return err
}

// sync pulls the authoritative state from storage; caller must hold the user
// lock. Treats NotFound as a legitimately empty user (cold-start only).
func (c *Cache) sync(ctx context.Context, userID string) error {
	return c.doSync(ctx, userID, true)
}

// Unlike sync, a NotFound here must not reset the cache -- a pending mutation
// is about to overwrite it, dropping a sibling writer's data.
func (c *Cache) resyncAfterConflict(ctx context.Context, userID string) error {
	return c.doSync(ctx, userID, false)
}

func (c *Cache) doSync(ctx context.Context, userID string, resetOnNotFound bool) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Sync")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))

	log := appctx.GetLogger(ctx).With().Str("userID", userID).Logger()

	c.initializeIfNeeded(userID, "")

	jsonPath := userJSONPath(userID)
	span.AddEvent("updating cache")
	//  - update cached list of created shares for the user in memory if changed
	rss, _ := c.ReceivedSpaces.Load(userID)
	dlreq := metadata.DownloadRequest{
		Path:        jsonPath,
		IfNoneMatch: []string{rss.etag},
	}
	dlres, err := c.storage.Download(ctx, dlreq)

	switch err.(type) {
	case nil:
		span.AddEvent("updating local cache")
	case errtypes.NotFound:
		if reset, trashed := cas.DecideNotFoundReset(ctx, c.storage, jsonPath, resetOnNotFound, rss.etag != "", log); !reset {
			span.SetStatus(codes.Error, err.Error())
			if trashed {
				log.Error().Err(err).Msg("lost update: backing file was recently trashed")
			} else {
				log.Error().Err(err).Msg("lost update: backing file disappeared mid-retry")
			}
			return err
		}
		c.ReceivedSpaces.Store(userID, &Spaces{Spaces: map[string]*Space{}})
		span.SetStatus(codes.Ok, "")
		return nil
	case errtypes.NotModified:
		span.SetStatus(codes.Ok, "")
		return nil
	default:
		span.SetStatus(codes.Error, err.Error())
		if cas.IsSyncTransient(err) {
			log.Warn().Err(err).Msg("lost update: re-read transient error")
		} else {
			log.Error().Err(err).Msg("lost update: re-read failed")
		}
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

	c.ReceivedSpaces.Store(userID, newSpaces)
	span.SetStatus(codes.Ok, "")
	return nil
}

// persist persists the data for one user to the storage
func (c *Cache) persist(ctx context.Context, userID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Persist")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.userid", userID))

	rss, ok := c.ReceivedSpaces.Load(userID)
	if !ok {
		span.SetStatus(codes.Ok, "no received shares")
		return nil
	}

	createdBytes, err := json.Marshal(rss)
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
		IfMatchEtag: rss.etag,
	}
	// when there is no etag in memory make sure the file has not been created on the server, see https://www.rfc-editor.org/rfc/rfc9110#field.if-match
	// > If the field value is "*", the condition is false if the origin server has a current representation for the target resource.
	if rss.etag == "" {
		ur.IfNoneMatch = []string{"*"}
	}

	res, err := c.storage.Upload(ctx, ur)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	rss.etag = res.Etag
	span.SetStatus(codes.Ok, "")
	return nil
}

func userJSONPath(userID string) string {
	return filepath.Join("/users", userID, "received.json")
}

func (c *Cache) initializeIfNeeded(userID, spaceID string) {
	rss, _ := c.ReceivedSpaces.LoadOrStore(userID, &Spaces{Spaces: map[string]*Space{}})
	if spaceID != "" && rss.Spaces[spaceID] == nil {
		rss.Spaces[spaceID] = &Space{}
		c.ReceivedSpaces.Store(userID, rss)
	}
}
