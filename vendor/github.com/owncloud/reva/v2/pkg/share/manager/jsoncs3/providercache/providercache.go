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

package providercache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	backoff "github.com/cenkalti/backoff/v5"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/appctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/cas"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/keylock"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/mtimesyncedcache"
	"github.com/owncloud/reva/v2/pkg/storage/utils/metadata"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/proto"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/providercache")
}

// Cache holds share information structured by provider and space
type Cache struct {
	lockMap keylock.Map

	Providers mtimesyncedcache.Map[string, *Spaces]

	storage metadata.Storage
	ttl     time.Duration
}

// Spaces holds the share information for provider
type Spaces struct {
	Spaces mtimesyncedcache.Map[string, *Shares]
}

// Shares holds the share information of one space
type Shares struct {
	Shares map[string]*collaboration.Share

	Etag string
}

// UnmarshalJSON overrides the default unmarshaling
// Shares are tricky to unmarshal because they contain an interface (Grantee) which makes the json Unmarshal bail out
// To work around that problem we unmarshal into json.RawMessage in a first step and then try to manually unmarshal
// into the specific types in a second step.
func (s *Shares) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Shares map[string]json.RawMessage
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.Shares = make(map[string]*collaboration.Share, len(tmp.Shares))
	for id, genericShare := range tmp.Shares {
		userShare := &collaboration.Share{
			Grantee: &provider.Grantee{Id: &provider.Grantee_UserId{}},
		}
		err = json.Unmarshal(genericShare, userShare) // is this a user share?
		if err == nil && userShare.Grantee.Type == provider.GranteeType_GRANTEE_TYPE_USER {
			s.Shares[id] = userShare
			continue
		}

		groupShare := &collaboration.Share{
			Grantee: &provider.Grantee{Id: &provider.Grantee_GroupId{}},
		}
		err = json.Unmarshal(genericShare, groupShare) // is this a group share?
		if err == nil && groupShare.Grantee.Type == provider.GranteeType_GRANTEE_TYPE_GROUP {
			s.Shares[id] = groupShare
			continue
		}

		invalidShare := &collaboration.Share{}
		err = json.Unmarshal(genericShare, invalidShare) // invalid
		if err == nil {
			s.Shares[id] = invalidShare
			continue
		}

		return err
	}

	return nil
}

// LockSpace locks the cache for a given space and returns an unlock function.
// It gives up with ctx.Err() once ctx is done, so abandoned requests don't stay
// parked behind a slow holder; the returned unlock func is nil in that case.
func (c *Cache) LockSpace(ctx context.Context, spaceID string) (func(), error) {
	return c.lockMap.Lock(ctx, spaceID)
}

// New returns a new Cache instance
func New(s metadata.Storage, ttl time.Duration) Cache {
	return Cache{
		Providers: mtimesyncedcache.Map[string, *Spaces]{},
		storage:   s,
		ttl:       ttl,
	}
}

func (c *Cache) isSpaceCached(storageID, spaceID string) bool {
	spaces, ok := c.Providers.Load(storageID)
	if !ok {
		return false
	}
	_, ok = spaces.Spaces.Load(spaceID)
	return ok
}

// hasAll reports if the space is in memory and has all the given shares.
func (c *Cache) hasAll(storageID, spaceID string, shareIDs ...string) bool {
	spaces, ok := c.Providers.Load(storageID)
	if !ok {
		return false
	}
	space, ok := spaces.Spaces.Load(spaceID)
	if !ok {
		return false
	}
	for _, shareID := range shareIDs {
		if _, ok := space.Shares[shareID]; !ok {
			return false
		}
	}
	return true
}

// dropOnFailedUpload drops the space from memory after a failed upload that
// is not a CAS conflict. The memory state is then different from the file,
// and the etag did not change, so a sync gets NotModified and keeps it.
// After the drop, the next call reads the file. A CAS conflict syncs anyway.
// Call it with the space lock held.
func (c *Cache) dropOnFailedUpload(err error, storageID, spaceID string) {
	if err == nil || cas.IsConflict(err) {
		return
	}
	if spaces, ok := c.Providers.Load(storageID); ok {
		spaces.Spaces.Delete(spaceID)
	}
}

// Add adds a share to the cache
func (c *Cache) Add(ctx context.Context, storageID, spaceID, shareID string, share *collaboration.Share) error {
	ctx, span := tracer.Start(ctx, "Add")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID), attribute.String("cs3.shareid", shareID))

	switch {
	case storageID == "":
		return fmt.Errorf("missing storage id")
	case spaceID == "":
		return fmt.Errorf("missing space id")
	case shareID == "":
		return fmt.Errorf("missing share id")
	}

	unlock, err := c.LockSpace(ctx, spaceID)
	if err != nil {
		return err
	}
	defer unlock()
	span.AddEvent("got lock")

	if !c.isSpaceCached(storageID, spaceID) {
		err = c.syncWithRetry(ctx, storageID, spaceID)
		if err != nil {
			return err
		}
	}

	log := appctx.GetLogger(ctx).With().
		Str("hostname", os.Getenv("HOSTNAME")).
		Str("storageID", storageID).
		Str("spaceID", spaceID).
		Str("shareID", shareID).Logger()

	persistFunc := func() error {

		spaces, _ := c.Providers.Load(storageID)
		space, _ := spaces.Spaces.Load(spaceID)

		log.Info().Interface("shares", maps.Keys(space.Shares)).Str("New share", shareID).Msg("Adding share to space")
		space.Shares[shareID] = share

		err := c.Persist(ctx, storageID, spaceID)
		// With the share in memory, this pod lists a share that is not in the file.
		c.dropOnFailedUpload(err, storageID, spaceID)
		return err
	}

	err = cas.RetryPersist(ctx, 100, persistFunc,
		func() error { return c.resyncAfterConflict(ctx, storageID, spaceID) },
		func(err error) { log.Debug().Err(err).Msg("CAS conflict persisting added provider share, retrying...") },
		func(err error) {
			span.SetStatus(codes.Error, fmt.Sprintf("persisting added provider share failed. giving up: %s", err.Error()))
			log.Error().Err(err).Msg("persisting added provider share failed")
		},
		func(err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error().Err(err).Msg("persisting added provider share failed. giving up.")
		},
	)
	if err == nil {
		span.SetStatus(codes.Ok, "")
	}
	return err
}

// Remove removes a share from the cache
func (c *Cache) Remove(ctx context.Context, storageID, spaceID, shareID string) error {
	_, err := c.RemoveMany(ctx, storageID, spaceID, []string{shareID})
	return err
}

// RemoveMany removes the given shares of one space with one upload.
// It does not upload when none of the shares is in the space.
// It returns the IDs that this call removed. A share that another caller
// removed first is not in the result. The CAS upload makes sure that only one
// caller removes a share, also across pods.
func (c *Cache) RemoveMany(ctx context.Context, storageID, spaceID string, shareIDs []string) ([]string, error) {
	ctx, span := tracer.Start(ctx, "Remove")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID), attribute.Int("cs3.sharecount", len(shareIDs)))

	unlock, err := c.LockSpace(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	defer unlock()
	span.AddEvent("got lock")

	// When a share is not in memory, sync first. A cold or stale cache does
	// not show if the share is in the file. For example, a listing that got a
	// temporary NotFound stores an empty space. Without the sync, the removal
	// does not upload and the share stays in the file. When all shares are in
	// memory, the CAS write in Persist detects stale state.
	if !c.hasAll(storageID, spaceID, shareIDs...) {
		err := c.syncWithRetry(ctx, storageID, spaceID)
		if err != nil {
			return nil, err
		}
	}

	var removed []string
	persistFunc := func() error {
		// A CAS conflict syncs the space and calls this function again.
		removed = nil
		spaces, ok := c.Providers.Load(storageID)
		if !ok {
			return nil
		}
		space, ok := spaces.Spaces.Load(spaceID)
		if !ok {
			return nil
		}
		for _, shareID := range shareIDs {
			if _, ok := space.Shares[shareID]; ok {
				delete(space.Shares, shareID)
				removed = append(removed, shareID)
			}
		}
		if len(removed) == 0 {
			return nil // nothing to remove, do not upload
		}

		err := c.Persist(ctx, storageID, spaceID)
		// Without the shares in memory, a retry finds nothing to remove and does not upload.
		c.dropOnFailedUpload(err, storageID, spaceID)
		return err
	}

	log := appctx.GetLogger(ctx).With().
		Str("hostname", os.Getenv("HOSTNAME")).
		Str("storageID", storageID).
		Str("spaceID", spaceID).
		Int("shareCount", len(shareIDs)).Logger()

	err = cas.RetryPersist(ctx, 100, persistFunc,
		func() error { return c.resyncAfterConflict(ctx, storageID, spaceID) },
		func(err error) {
			log.Debug().Err(err).Msg("CAS conflict persisting removed provider share, retrying...")
		},
		func(err error) {
			span.SetStatus(codes.Error, fmt.Sprintf("persisting removed provider share failed. giving up: %s", err.Error()))
			log.Error().Err(err).Msg("persisting removed provider share failed")
		},
		func(err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error().Err(err).Msg("persisting removed provider share failed. giving up.")
		},
	)
	if err == nil {
		span.SetStatus(codes.Ok, "")
		return removed, nil
	}
	return nil, err
}

// Update changes one share and uploads the space. It reads the share under
// the space lock and gives a copy to the update function. A CAS conflict syncs
// the space and calls the update function again with the new share.
// The cache keeps the old share until the upload succeeds.
// Update returns NotFound when the space does not have the share. Then it does
// not upload, so it does not write a share that a parallel caller removed.
func (c *Cache) Update(ctx context.Context, storageID, spaceID, shareID string, update func(*collaboration.Share)) (*collaboration.Share, error) {
	ctx, span := tracer.Start(ctx, "Update")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID), attribute.String("cs3.shareid", shareID))

	switch {
	case storageID == "":
		return nil, fmt.Errorf("missing storage id")
	case spaceID == "":
		return nil, fmt.Errorf("missing space id")
	case shareID == "":
		return nil, fmt.Errorf("missing share id")
	}

	unlock, err := c.LockSpace(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	defer unlock()
	span.AddEvent("got lock")

	// When the share is not in memory, sync first, as in RemoveMany. Else
	// Update returns NotFound for a share that is in the file.
	if !c.hasAll(storageID, spaceID, shareID) {
		if err := c.syncWithRetry(ctx, storageID, spaceID); err != nil {
			return nil, err
		}
	}

	var updated *collaboration.Share
	persistFunc := func() error {
		// A CAS conflict syncs the space and calls this function again.
		updated = nil
		spaces, ok := c.Providers.Load(storageID)
		if !ok {
			return nil
		}
		space, ok := spaces.Spaces.Load(spaceID)
		if !ok {
			return nil
		}
		old, ok := space.Shares[shareID]
		if !ok || old == nil {
			return nil // the share is gone, do not upload
		}
		s := proto.Clone(old).(*collaboration.Share)
		update(s)
		space.Shares[shareID] = s

		err := c.Persist(ctx, storageID, spaceID)
		if err != nil {
			// Put the old share back, so that memory agrees with the file.
			// Persist does not change the etag when the upload fails.
			space.Shares[shareID] = old
			return err
		}
		updated = s
		return nil
	}

	log := appctx.GetLogger(ctx).With().
		Str("hostname", os.Getenv("HOSTNAME")).
		Str("storageID", storageID).
		Str("spaceID", spaceID).
		Str("shareID", shareID).Logger()

	err = cas.RetryPersist(ctx, 100, persistFunc,
		func() error { return c.resyncAfterConflict(ctx, storageID, spaceID) },
		func(err error) {
			log.Debug().Err(err).Msg("CAS conflict persisting updated provider share, retrying...")
		},
		func(err error) {
			span.SetStatus(codes.Error, fmt.Sprintf("persisting updated provider share failed. giving up: %s", err.Error()))
			log.Error().Err(err).Msg("persisting updated provider share failed")
		},
		func(err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error().Err(err).Msg("persisting updated provider share failed. giving up.")
		},
	)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errtypes.NotFound(shareID)
	}
	span.SetStatus(codes.Ok, "")
	return updated, nil
}

// Get returns one entry from the cache
func (c *Cache) Get(ctx context.Context, storageID, spaceID, shareID string, skipSync bool) (*collaboration.Share, error) {
	ctx, span := tracer.Start(ctx, "Get")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID), attribute.String("cs3.shareid", shareID))

	unlock, err := c.LockSpace(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	defer unlock()
	span.AddEvent("got lock")

	if !skipSync {
		// sync cache, maybe our data is outdated
		err := c.syncWithRetry(ctx, storageID, spaceID)
		if err != nil {
			return nil, err
		}
	}

	spaces, ok := c.Providers.Load(storageID)
	if !ok {
		return nil, nil
	}
	space, ok := spaces.Spaces.Load(spaceID)
	if !ok {
		return nil, nil
	}
	return space.Shares[shareID], nil
}

// All returns all entries in the storage
func (c *Cache) All(ctx context.Context) (*mtimesyncedcache.Map[string, *Spaces], error) {
	ctx, span := tracer.Start(ctx, "All")
	defer span.End()

	providers, err := c.storage.ListDir(ctx, "/storages")
	if err != nil {
		return nil, err
	}
	for _, provider := range providers {
		storageID := provider.Name
		spaces, err := c.storage.ListDir(ctx, path.Join("/storages", storageID))
		if err != nil {
			return nil, err
		}
		for _, space := range spaces {
			spaceID := strings.TrimSuffix(space.Name, ".json")

			unlock, err := c.LockSpace(ctx, spaceID)
			if err != nil {
				return nil, err
			}
			span.AddEvent("got lock for space " + spaceID)
			if err := c.syncWithRetry(ctx, storageID, spaceID); err != nil {
				unlock()
				return nil, err
			}
			unlock()
		}
	}

	return &c.Providers, nil
}

// ListSpace returns the list of shares in a given space
func (c *Cache) ListSpace(ctx context.Context, storageID, spaceID string) (*Shares, error) {
	ctx, span := tracer.Start(ctx, "ListSpace")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID))

	unlock, err := c.LockSpace(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	defer unlock()
	span.AddEvent("got lock")

	// sync cache, maybe our data is outdated
	err = c.syncWithRetry(ctx, storageID, spaceID)
	if err != nil {
		return nil, err
	}

	spaces, ok := c.Providers.Load(storageID)
	if !ok {
		return &Shares{}, nil
	}

	space, ok := spaces.Spaces.Load(spaceID)
	if !ok {
		return &Shares{}, nil
	}

	shares := &Shares{
		Shares: maps.Clone(space.Shares),
		Etag:   space.Etag,
	}
	return shares, nil
}

// SpaceFileMissing reports if a Stat of the space file returns NotFound.
// Use it as a second proof after a download got NotFound. It returns false
// for all other results, also for errors.
func (c *Cache) SpaceFileMissing(ctx context.Context, storageID, spaceID string) bool {
	_, err := c.storage.Stat(ctx, spaceJSONPath(storageID, spaceID))
	_, notFound := err.(errtypes.NotFound)
	return notFound
}

// Persist persists the data of one space
func (c *Cache) Persist(ctx context.Context, storageID, spaceID string) error {
	ctx, span := tracer.Start(ctx, "Persist")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID))

	spaces, ok := c.Providers.Load(storageID)
	if !ok {
		span.AddEvent("nothing to persist")
		span.SetStatus(codes.Ok, "")
		return nil
	}
	space, ok := spaces.Spaces.Load(spaceID)
	if !ok {
		span.AddEvent("nothing to persist")
		span.SetStatus(codes.Ok, "")
		return nil
	}
	span.SetAttributes(attribute.String("BeforeEtag", space.Etag))
	log := appctx.GetLogger(ctx).With().Str("storageID", storageID).Str("spaceID", spaceID).Logger()
	log = log.With().Str("BeforeEtag", space.Etag).Logger()

	createdBytes, err := json.Marshal(space)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	jsonPath := spaceJSONPath(storageID, spaceID)
	if err := c.storage.MakeDirIfNotExist(ctx, path.Dir(jsonPath)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetAttributes(attribute.String("etag", space.Etag))

	ur := metadata.UploadRequest{
		Path:        jsonPath,
		Content:     createdBytes,
		IfMatchEtag: space.Etag,
	}
	// when there is no etag in memory make sure the file has not been created on the server, see https://www.rfc-editor.org/rfc/rfc9110#field.if-match
	// > If the field value is "*", the condition is false if the origin server has a current representation for the target resource.
	if space.Etag == "" {
		ur.IfNoneMatch = []string{"*"}
	}

	res, err := c.storage.Upload(ctx, ur)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Debug().Err(err).Msg("persisting provider cache failed")
		return err
	}
	space.Etag = res.Etag

	span.SetStatus(codes.Ok, "")
	shares := []string{}
	for _, s := range space.Shares {
		shares = append(shares, s.GetId().GetOpaqueId())
	}
	log.Debug().Str("AfterEtag", space.Etag).Interface("Shares", shares).Msg("persisted provider cache")
	return nil
}

// PurgeSpace removes a space from the cache
func (c *Cache) PurgeSpace(ctx context.Context, storageID, spaceID string) error {
	ctx, span := tracer.Start(ctx, "PurgeSpace")
	defer span.End()

	unlock, err := c.LockSpace(ctx, spaceID)
	if err != nil {
		return err
	}
	defer unlock()
	span.AddEvent("got lock")

	if !c.isSpaceCached(storageID, spaceID) {
		err := c.syncWithRetry(ctx, storageID, spaceID)
		if err != nil {
			return err
		}
	}

	spaces, ok := c.Providers.Load(storageID)
	if !ok {
		return nil
	}
	// Add does not sync a space that is in memory, and it writes to this map.
	newShares := &Shares{Shares: map[string]*collaboration.Share{}}
	if space, ok := spaces.Spaces.Load(spaceID); ok {
		newShares.Etag = space.Etag // keep the etag to allow overwriting the state on the server
	}
	spaces.Spaces.Store(spaceID, newShares)

	return c.Persist(ctx, storageID, spaceID)
}

// syncWithLock treats NotFound as a legitimately empty space (cold-start only).
func (c *Cache) syncWithLock(ctx context.Context, storageID, spaceID string) error {
	return c.doSync(ctx, storageID, spaceID, true)
}

// Unlike syncWithLock, a NotFound here must not reset the cache -- a pending
// mutation is about to overwrite it, dropping a sibling writer's data.
func (c *Cache) resyncAfterConflict(ctx context.Context, storageID, spaceID string) error {
	return c.doSync(ctx, storageID, spaceID, false)
}

func (c *Cache) doSync(ctx context.Context, storageID, spaceID string, resetOnNotFound bool) error {
	ctx, span := tracer.Start(ctx, "syncWithLock")
	defer span.End()

	c.initializeIfNeeded(storageID, spaceID)

	spaces, _ := c.Providers.Load(storageID)
	space, _ := spaces.Spaces.Load(spaceID)
	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID), attribute.String("etag", space.Etag))
	log := appctx.GetLogger(ctx).With().Str("storageID", storageID).Str("spaceID", spaceID).Str("etag", space.Etag).Str("hostname", os.Getenv("HOSTNAME")).Logger()

	dlreq := metadata.DownloadRequest{
		Path: spaceJSONPath(storageID, spaceID),
	}
	// when we know an etag, only download if it changed remotely
	if space.Etag != "" {
		dlreq.IfNoneMatch = []string{space.Etag}
	}

	dlres, err := c.storage.Download(ctx, dlreq)

	switch err.(type) {
	case nil:
		span.AddEvent("updating local cache")
	case errtypes.NotFound:
		if reset, trashed := cas.DecideNotFoundReset(ctx, c.storage, dlreq.Path, resetOnNotFound, space.Etag != "", log); !reset {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			if trashed {
				log.Error().Err(err).Msg("lost update: backing file was recently trashed")
			} else {
				log.Error().Err(err).Msg("lost update: backing file disappeared mid-retry")
			}
			return err
		}
		spaces.Spaces.Store(spaceID, &Shares{Shares: map[string]*collaboration.Share{}})
		span.SetStatus(codes.Ok, "")
		return nil
	case errtypes.NotModified:
		span.SetStatus(codes.Ok, "")
		return nil
	default:
		span.RecordError(err)
		span.SetStatus(codes.Error, "downloading provider cache failed")
		return err
	}

	span.AddEvent("updating local cache")
	newShares := &Shares{}
	err = json.Unmarshal(dlres.Content, newShares)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unmarshaling provider cache failed")
		log.Error().Err(err).Msg("unmarshaling provider cache failed")
		return err
	}
	newShares.Etag = dlres.Etag

	spaces.Spaces.Store(spaceID, newShares)
	span.SetStatus(codes.Ok, "")
	return nil
}

// syncWithRetry retries syncWithLock's cold-start read on transient errors, since callers have no other retry wrapper.
func (c *Cache) syncWithRetry(ctx context.Context, storageID, spaceID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		err := c.syncWithLock(ctx, storageID, spaceID)
		if err != nil && !cas.IsSyncTransient(err) {
			return struct{}{}, backoff.Permanent(err)
		}
		return struct{}{}, err
	}, backoff.WithBackOff(cas.NewBackoff()), backoff.WithMaxTries(20))
	return err
}

func (c *Cache) initializeIfNeeded(storageID, spaceID string) {
	spaces, _ := c.Providers.LoadOrStore(storageID, &Spaces{
		Spaces: mtimesyncedcache.Map[string, *Shares]{},
	})
	_, _ = spaces.Spaces.LoadOrStore(spaceID, &Shares{
		Shares: map[string]*collaboration.Share{},
	})
}

func spaceJSONPath(storageID, spaceID string) string {
	return filepath.Join("/storages", storageID, spaceID+".json")
}
