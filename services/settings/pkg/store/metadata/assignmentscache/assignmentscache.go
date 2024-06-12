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

package assignmentscache

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/mtimesyncedcache"
	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// name is the Tracer name used to identify this instrumentation library.
const tracerName = "assignmentscache"

// Cache caches the list of roleassignments for roles
// It functions as an in-memory cache with a persistence layer
// The storage is sharded by roleid
type Cache struct {
	lockMap sync.Map

	RoleAssignments mtimesyncedcache.Map[string, *RoleAssignmentCache]

	storage   metadata.Storage
	namespace string
	filename  string
}

// RoleAssignmentCache holds the assignments for one role
type RoleAssignmentCache struct {
	RoleAssignments map[string]*RoleAssignment `json:"roleassignments"`
	Etag            string                     `json:"etag"`
}

// RoleAssignment holds the unique list of assignments ids for a role
type RoleAssignment struct {
	AssignmentID string `json:"assignmentid"`
}

func (c *Cache) lockRole(roleID string) func() {
	v, _ := c.lockMap.LoadOrStore(roleID, &sync.Mutex{})
	lock := v.(*sync.Mutex)

	lock.Lock()
	return func() { lock.Unlock() }
}

// New returns a new Cache instance
func New(s metadata.Storage, namespace, filename string) Cache {
	return Cache{
		RoleAssignments: mtimesyncedcache.Map[string, *RoleAssignmentCache]{},
		storage:         s,
		namespace:       namespace,
		filename:        filename,
		lockMap:         sync.Map{},
	}
}

// Add adds a role assignment to the cache
func (c *Cache) Add(ctx context.Context, roleID string, assignment *settingsmsg.UserRoleAssignment) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Grab lock")
	unlock := c.lockRole(roleID)
	span.End()
	span.SetAttributes(attribute.String("roleid", roleID))
	defer unlock()
	log := appctx.GetLogger(ctx).With().
		Str("hostname", os.Getenv("HOSTNAME")).
		Str("roleID", roleID).
		Str("assignmentID", assignment.GetId()).Logger()

	if _, ok := c.RoleAssignments.Load(roleID); !ok {
		err := c.syncWithLock(ctx, roleID)
		if err != nil {
			return err
		}
	}

	ctx, span = appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Add")
	defer span.End()
	span.SetAttributes(attribute.String("roleid", roleID), attribute.String("assignmentid", assignment.GetId()))

	persistFunc := func() error {
		c.initializeIfNeeded(roleID, assignment.GetAccountUuid())

		us, _ := c.RoleAssignments.Load(roleID)
		us.RoleAssignments[assignment.GetAccountUuid()] = &RoleAssignment{
			AssignmentID: assignment.GetId(),
		}

		return c.Persist(ctx, roleID)
	}

	var err error
	for retries := 100; retries > 0; retries-- {
		err = persistFunc()
		switch err.(type) {
		case nil:
			span.SetStatus(codes.Ok, "")
			return nil
		case errtypes.Aborted:
			log.Debug().Msg("aborted when persisting added assignemt: etag changed. retrying...")
			// this is the expected status code from the server when the if-match etag check fails
			// continue with sync below
		case errtypes.PreconditionFailed:
			log.Debug().Msg("precondition failed when persisting added assignemt: etag changed. retrying...")
			// actually, this is the wrong status code and we treat it like errtypes.Aborted because of inconsistencies on the server side
			// continue with sync below
		case errtypes.AlreadyExists:
			log.Debug().Msg("already exists when persisting added assignemt. retrying...")
			// CS3 uses an already exists error instead of precondition failed when using an If-None-Match=* header / IfExists flag in the InitiateFileUpload call.
			// Thas happens when the cache thinks there is no file.
			// continue with sync below
		default:
			span.SetStatus(codes.Error, "persisting added assignemt failed. giving up: "+err.Error())
			log.Error().Err(err).Msg("persisting added assignemt failed")
			return err
		}
		if err := c.syncWithLock(ctx, roleID); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error().Err(err).Msg("persisting added assignment failed. giving up.")
			return err
		}
	}
	return err
}

// Remove removes an assignment from the roles cache
func (c *Cache) Remove(ctx context.Context, roleID, accountID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Grab lock")
	unlock := c.lockRole(roleID)
	span.End()
	span.SetAttributes(attribute.String("roleid", roleID))
	defer unlock()

	if _, ok := c.RoleAssignments.Load(roleID); ok {
		err := c.syncWithLock(ctx, roleID)
		if err != nil {
			return err
		}
	}

	ctx, span = appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Remove")
	defer span.End()
	span.SetAttributes(attribute.String("roleid", roleID), attribute.String("userid", accountID))

	persistFunc := func() error {
		us, loaded := c.RoleAssignments.LoadOrStore(roleID, &RoleAssignmentCache{
			RoleAssignments: map[string]*RoleAssignment{},
		})

		if loaded {
			// remove user id
			delete(us.RoleAssignments, accountID)
		}

		return c.Persist(ctx, roleID)
	}

	log := appctx.GetLogger(ctx).With().
		Str("hostname", os.Getenv("HOSTNAME")).
		Str("roleID", roleID).
		Str("accountID", accountID).Logger()

	var err error
	for retries := 100; retries > 0; retries-- {
		err = persistFunc()
		switch err.(type) {
		case nil:
			span.SetStatus(codes.Ok, "")
			return nil
		case errtypes.Aborted:
			log.Debug().Msg("aborted when persisting removed assignment: etag changed. retrying...")
			// this is the expected status code from the server when the if-match etag check fails
			// continue with sync below
		case errtypes.PreconditionFailed:
			log.Debug().Msg("precondition failed when persisting removed assignment: etag changed. retrying...")
			// actually, this is the wrong status code and we treat it like errtypes.Aborted because of inconsistencies on the server side
			// continue with sync below
		default:
			span.SetStatus(codes.Error, "persisting removed assignment failed. giving up: "+err.Error())
			log.Error().Err(err).Msg("persisting removed assignment failed")
			return err
		}
		if err := c.syncWithLock(ctx, roleID); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
	}

	return err
}

// List return the list of assignments for the given role
func (c *Cache) List(ctx context.Context, roleID string) (map[string]RoleAssignment, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Grab lock")
	unlock := c.lockRole(roleID)
	span.End()
	span.SetAttributes(attribute.String("roleid", roleID))
	defer unlock()
	if err := c.syncWithLock(ctx, roleID); err != nil {
		return nil, err
	}

	r := map[string]RoleAssignment{}
	us, ok := c.RoleAssignments.Load(roleID)
	if !ok {
		return r, nil
	}

	for roleid, cached := range us.RoleAssignments {
		r[roleid] = *cached
	}
	return r, nil
}

func (c *Cache) syncWithLock(ctx context.Context, roleID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Sync")
	defer span.End()
	span.SetAttributes(attribute.String("roleid", roleID))

	log := appctx.GetLogger(ctx).With().Str("roleID", roleID).Logger()

	c.initializeIfNeeded(roleID, "")

	assignmentsCachePath := c.assignmentsForRolePath(roleID)
	span.AddEvent("updating cache")
	//  - update cached list of assignments for the role in memory if changed
	dlreq := metadata.DownloadRequest{
		Path: assignmentsCachePath,
	}
	if us, ok := c.RoleAssignments.Load(roleID); ok && us.Etag != "" {
		dlreq.IfNoneMatch = []string{us.Etag}
	}

	dlres, err := c.storage.Download(ctx, dlreq)
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
		span.SetStatus(codes.Error, "Failed to download the assignment cache: "+err.Error())
		log.Error().Err(err).Msg("Failed to download the assignment cache")
		return err
	}

	assignmentCache := &RoleAssignmentCache{}
	err = json.Unmarshal(dlres.Content, assignmentCache)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to unmarshal the assignment cache: "+err.Error())
		log.Error().Err(err).Msg("Failed to unmarshal the assignment cache")
		return err
	}
	assignmentCache.Etag = dlres.Etag

	c.RoleAssignments.Store(roleID, assignmentCache)
	span.SetStatus(codes.Ok, "")
	return nil
}

// Persist persists the data for one role to the storage
func (c *Cache) Persist(ctx context.Context, roleID string) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Persist")
	defer span.End()
	span.SetAttributes(attribute.String("roleid", roleID))

	ra, ok := c.RoleAssignments.Load(roleID)
	if !ok {
		span.SetStatus(codes.Ok, "no role assignments")
		return nil
	}
	createdBytes, err := json.Marshal(ra)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	jsonPath := c.assignmentsForRolePath(roleID)
	if err := c.storage.MakeDirIfNotExist(ctx, path.Dir(jsonPath)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	ur := metadata.UploadRequest{
		Path:        jsonPath,
		Content:     createdBytes,
		IfMatchEtag: ra.Etag,
	}
	// when there is no etag in memory make sure the file has not been created on the server, see https://www.rfc-editor.org/rfc/rfc9110#field.if-match
	// > If the field value is "*", the condition is false if the origin server has a current representation for the target resource.
	if ra.Etag == "" {
		ur.IfNoneMatch = []string{"*"}
	}

	res, err := c.storage.Upload(ctx, ur)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	ra.Etag = res.Etag

	span.SetStatus(codes.Ok, "")
	return nil
}

func (c *Cache) assignmentsForRolePath(roleid string) string {
	return filepath.Join("/", c.namespace, roleid, c.filename)
}

func (c *Cache) initializeIfNeeded(roleID, accountID string) {
	us, _ := c.RoleAssignments.LoadOrStore(roleID, &RoleAssignmentCache{
		RoleAssignments: map[string]*RoleAssignment{},
	})
	if accountID != "" && us.RoleAssignments[accountID] == nil {
		us.RoleAssignments[accountID] = &RoleAssignment{}
	}
}
