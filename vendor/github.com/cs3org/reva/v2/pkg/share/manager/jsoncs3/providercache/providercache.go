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
const tracerName = "providercache"

// Cache holds share information structured by provider and space
type Cache struct {
	Providers map[string]*Spaces

	storage metadata.Storage
	ttl     time.Duration
}

// Spaces holds the share information for provider
type Spaces struct {
	Spaces map[string]*Shares
}

// Shares holds the share information of one space
type Shares struct {
	Shares   map[string]*collaboration.Share
	Mtime    time.Time
	nextSync time.Time
}

// UnmarshalJSON overrides the default unmarshaling
// Shares are tricky to unmarshal because they contain an interface (Grantee) which makes the json Unmarshal bail out
// To work around that problem we unmarshal into json.RawMessage in a first step and then try to manually unmarshal
// into the specific types in a second step.
func (s *Shares) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Shares map[string]json.RawMessage
		Mtime  time.Time
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.Mtime = tmp.Mtime
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
		err = json.Unmarshal(genericShare, groupShare) // try to unmarshal to a group share if the user share unmarshalling failed
		if err != nil {
			return err
		}
		s.Shares[id] = groupShare
	}

	return nil
}

// New returns a new Cache instance
func New(s metadata.Storage, ttl time.Duration) Cache {
	return Cache{
		Providers: map[string]*Spaces{},
		storage:   s,
		ttl:       ttl,
	}
}

// Add adds a share to the cache
func (c *Cache) Add(ctx context.Context, storageID, spaceID, shareID string, share *collaboration.Share) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Add")
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
	c.initializeIfNeeded(storageID, spaceID)
	c.Providers[storageID].Spaces[spaceID].Shares[shareID] = share

	return c.Persist(ctx, storageID, spaceID)
}

// Remove removes a share from the cache
func (c *Cache) Remove(ctx context.Context, storageID, spaceID, shareID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Remove")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID), attribute.String("cs3.shareid", shareID))

	if c.Providers[storageID] == nil ||
		c.Providers[storageID].Spaces[spaceID] == nil {
		return nil
	}
	delete(c.Providers[storageID].Spaces[spaceID].Shares, shareID)

	return c.Persist(ctx, storageID, spaceID)
}

// Get returns one entry from the cache
func (c *Cache) Get(storageID, spaceID, shareID string) *collaboration.Share {
	if c.Providers[storageID] == nil ||
		c.Providers[storageID].Spaces[spaceID] == nil {
		return nil
	}
	return c.Providers[storageID].Spaces[spaceID].Shares[shareID]
}

// ListSpace returns the list of shares in a given space
func (c *Cache) ListSpace(storageID, spaceID string) *Shares {
	if c.Providers[storageID] == nil || c.Providers[storageID].Spaces[spaceID] == nil {
		return &Shares{}
	}
	return c.Providers[storageID].Spaces[spaceID]
}

// PersistWithTime persists the data of one space if it has not been modified since the given mtime
func (c *Cache) PersistWithTime(ctx context.Context, storageID, spaceID string, mtime time.Time) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "PersistWithTime")
	defer span.End()
	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID))

	if c.Providers[storageID] == nil || c.Providers[storageID].Spaces[spaceID] == nil {
		return nil
	}

	oldMtime := c.Providers[storageID].Spaces[spaceID].Mtime
	c.Providers[storageID].Spaces[spaceID].Mtime = mtime

	// FIXME there is a race when between this time now and the below Uploed another process also updates the file -> we need a lock
	createdBytes, err := json.Marshal(c.Providers[storageID].Spaces[spaceID])
	if err != nil {
		c.Providers[storageID].Spaces[spaceID].Mtime = oldMtime
		return err
	}
	jsonPath := spaceJSONPath(storageID, spaceID)
	if err := c.storage.MakeDirIfNotExist(ctx, path.Dir(jsonPath)); err != nil {
		c.Providers[storageID].Spaces[spaceID].Mtime = oldMtime
		return err
	}

	if err = c.storage.Upload(ctx, metadata.UploadRequest{
		Path:              jsonPath,
		Content:           createdBytes,
		IfUnmodifiedSince: c.Providers[storageID].Spaces[spaceID].Mtime,
	}); err != nil {
		c.Providers[storageID].Spaces[spaceID].Mtime = oldMtime
		return err
	}
	return nil
}

// Persist persists the data of one space
func (c *Cache) Persist(ctx context.Context, storageID, spaceID string) error {
	return c.PersistWithTime(ctx, storageID, spaceID, time.Now())
}

// Sync updates the in-memory data with the data from the storage if it is outdated
func (c *Cache) Sync(ctx context.Context, storageID, spaceID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Sync")
	defer span.End()

	span.SetAttributes(attribute.String("cs3.storageid", storageID), attribute.String("cs3.spaceid", spaceID))

	log := appctx.GetLogger(ctx).With().Str("storageID", storageID).Str("spaceID", spaceID).Logger()

	var mtime time.Time
	if c.Providers[storageID] != nil && c.Providers[storageID].Spaces[spaceID] != nil {
		mtime = c.Providers[storageID].Spaces[spaceID].Mtime

		if time.Now().Before(c.Providers[storageID].Spaces[spaceID].nextSync) {
			span.AddEvent("skip sync")
			span.SetStatus(codes.Ok, "")
			return nil
		}
		c.Providers[storageID].Spaces[spaceID].nextSync = time.Now().Add(c.ttl)
	} else {
		mtime = time.Time{} // Set zero time so that data from storage always takes precedence
	}

	jsonPath := spaceJSONPath(storageID, spaceID)
	info, err := c.storage.Stat(ctx, jsonPath)
	if err != nil {
		if _, ok := err.(errtypes.NotFound); ok {
			span.AddEvent("no file")
			span.SetStatus(codes.Ok, "")
			return nil // Nothing to sync against
		}
		if _, ok := err.(*os.PathError); ok {
			span.AddEvent("no dir")
			span.SetStatus(codes.Ok, "")
			return nil // Nothing to sync against
		}
		span.SetStatus(codes.Error, fmt.Sprintf("Failed to stat the provider cache: %s", err.Error()))
		log.Error().Err(err).Msg("Failed to stat the provider cache")
		return err
	}
	// check mtime of /users/{userid}/created.json
	if utils.TSToTime(info.Mtime).After(mtime) {
		span.AddEvent("updating cache")
		//  - update cached list of created shares for the user in memory if changed
		createdBlob, err := c.storage.SimpleDownload(ctx, jsonPath)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("Failed to download the provider cache: %s", err.Error()))
			log.Error().Err(err).Msg("Failed to download the provider cache")
			return err
		}
		newShares := &Shares{}
		err = json.Unmarshal(createdBlob, newShares)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("Failed to unmarshal the provider cache: %s", err.Error()))
			log.Error().Err(err).Msg("Failed to unmarshal the provider cache")
			return err
		}
		newShares.Mtime = utils.TSToTime(info.Mtime)
		c.initializeIfNeeded(storageID, spaceID)
		c.Providers[storageID].Spaces[spaceID] = newShares
	}
	span.SetStatus(codes.Ok, "")
	return nil
}

func (c *Cache) initializeIfNeeded(storageID, spaceID string) {
	if c.Providers[storageID] == nil {
		c.Providers[storageID] = &Spaces{
			Spaces: map[string]*Shares{},
		}
	}
	if c.Providers[storageID].Spaces[spaceID] == nil {
		c.Providers[storageID].Spaces[spaceID] = &Shares{
			Shares: map[string]*collaboration.Share{},
		}
	}
}

func spaceJSONPath(storageID, spaceID string) string {
	return filepath.Join("/storages", storageID, spaceID+".json")
}
