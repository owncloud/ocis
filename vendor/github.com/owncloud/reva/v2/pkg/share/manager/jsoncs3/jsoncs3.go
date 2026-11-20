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

package jsoncs3

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"sync"
	"time"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/owncloud/reva/v2/pkg/appctx"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/owncloud/reva/v2/pkg/logger"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/share"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/providercache"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/receivedsharecache"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/sharecache"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/shareid"
	"github.com/owncloud/reva/v2/pkg/share/manager/registry"
	"github.com/owncloud/reva/v2/pkg/storage/utils/metadata" // nolint:staticcheck // we need the legacy package to convert V1 to V2 messages
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

/*
  The sharded json driver splits the json file per storage space. Similar to fileids shareids are prefixed with the spaceid for easier lookup.
  In addition to the space json the share manager keeps lists for users and groups to cache their lists of created and received shares
  and to hold the state of received shares.

  FAQ
  Q: Why not split shares by user and have a list per user?
  A: While shares are created by users, they are persisted as grants on a file.
     If we persist shares by their creator/owner they would vanish if a user is deprovisioned: shares
	 in project spaces could not be managed collaboratively.
	 By splitting by space, we are in fact not only splitting by user, but more granular, per space.


  File structure in the jsoncs3 space:

  /storages/{storageid}/{spaceid.json} 	// contains the share information of all shares in that space
  /users/{userid}/created.json			// points to the spaces the user created shares in, including the list of shares
  /users/{userid}/received.json			// holds the accepted/pending state and mount point of received shares for users
  /groups/{groupid}/received.json		// points to the spaces the group has received shares in including the list of shares

  Example:
  	├── groups
  	│	└── group1
  	│		└── received.json
  	├── storages
  	│	└── storageid
  	│		└── spaceid.json
  	└── users
   		├── admin
 		│	└── created.json
 		└── einstein
 			└── received.json

  Whenever a share is created, the share manager has to
  1. update the /storages/{storageid}/{spaceid}.json file,
  2. create /users/{userid}/created.json if it doesn't exist yet and add the space/share
  3. create /users/{userid}/received.json or /groups/{groupid}/received.json if it doesn exist yet and add the space/share

  When updating shares /storages/{storageid}/{spaceid}.json is updated accordingly. The etag is used to invalidate in-memory caches:
  - TODO the upload is tried with an if-unmodified-since header
  - TODO when if fails, the {spaceid}.json file is downloaded, the changes are reapplied and the upload is retried with the new etag

  When updating received shares the mountpoint and state are updated in /users/{userid}/received.json (for both user and group shares).

  When reading the list of received shares the /users/{userid}/received.json file and the /groups/{groupid}/received.json files are statted.
  - if the etag changed we download the file to update the local cache

  When reading the list of created shares the /users/{userid}/created.json file is statted
  - if the etag changed we download the file to update the local cache
*/

// TODO implement a channel based aggregation of sharing requests: every in memory cache should read as many share updates to a space that are available and update them all in one go
// whenever a persist operation fails we check if we can read more shares from the channel

// name is the Tracer name used to identify this instrumentation library.
const tracerName = "jsoncs3"

// rollbackTimeout limits the rollback of a failed Share. The rollback does not
// use the cancellation of the request, because a cancelled request is a
// common cause of the failure. It is a variable, so the tests can make it short.
var rollbackTimeout = 30 * time.Second

// defaultStaleShareGracePeriod is the default of Manager.StaleShareGracePeriod.
// An older version writes the grantee entry and the share at the same time.
// During a rolling upgrade, a listing can find the grantee entry before the
// share is on storage. The grace period must be longer than such a write.
const defaultStaleShareGracePeriod = 5 * time.Minute

// defaultStaleForgetAfter is the default of Manager.staleForgetAfter. It must
// be longer than the grace period. It must also be longer than the usual time
// between two listings that find the same stale share on one pod. Otherwise
// the pod starts the grace period again and does not remove the entry. One
// session of an active user lists the received shares many times. A stale
// entry that stays costs only a debug line in each listing, because the
// listing does not upload in the grace period.
const defaultStaleForgetAfter = time.Hour

// staleSweepsPerForget sets how often a pod looks for the stale shares to
// forget. Each sweep goes through all entries under staleMu. A pod sweeps at
// most once in each staleForgetAfter/staleSweepsPerForget, so a listing that
// finds many new stale shares does not run many sweeps. A pod forgets an entry
// at most one sweep interval later than staleForgetAfter.
const staleSweepsPerForget = 60

func init() {
	registry.Register("jsoncs3", NewDefault)
}

var (
	_registeredEvents = []events.Unmarshaller{
		events.SpaceDeleted{},
	}
)

type config struct {
	GatewayAddr       string       `mapstructure:"gateway_addr"`
	MaxConcurrency    int          `mapstructure:"max_concurrency"`
	ProviderAddr      string       `mapstructure:"provider_addr"`
	ServiceUserID     string       `mapstructure:"service_user_id"`
	ServiceUserIdp    string       `mapstructure:"service_user_idp"`
	MachineAuthAPIKey string       `mapstructure:"machine_auth_apikey"`
	CacheTTL          int          `mapstructure:"ttl"`
	Events            EventOptions `mapstructure:"events"`
}

// EventOptions are the configurable options for events
type EventOptions struct {
	Endpoint             string `mapstructure:"natsaddress"`
	Cluster              string `mapstructure:"natsclusterid"`
	TLSInsecure          bool   `mapstructure:"tlsinsecure"`
	TLSRootCACertificate string `mapstructure:"tlsrootcacertificate"`
	EnableTLS            bool   `mapstructure:"enabletls"`
	AuthUsername         string `mapstructure:"authusername"`
	AuthPassword         string `mapstructure:"authpassword"`
	CredsFile            string `mapstructure:"credsfile"`
	Creds                string `mapstructure:"creds"`
}

// Manager implements a share manager using a cs3 storage backend with local caching
type Manager struct {
	sync.RWMutex

	Cache              providercache.Cache      // holds all shares, sharded by provider id and space id
	CreatedCache       sharecache.Cache         // holds the list of shares a user has created, sharded by user id
	GroupReceivedCache sharecache.Cache         // holds the list of shares a group has access to, sharded by group id
	UserReceivedStates receivedsharecache.Cache // holds the state of shares a user has received, sharded by user id

	storage   metadata.Storage
	SpaceRoot *provider.ResourceId

	initialized bool

	MaxConcurrency int

	// StaleShareGracePeriod is the time that a pod waits before it removes a
	// stale grantee entry. Zero removes the entry at the first listing.
	StaleShareGracePeriod time.Duration

	// staleSince holds the times when this pod found each stale share.
	staleMu    sync.Mutex
	staleSince map[string]staleSeen
	// staleForgetAfter is the time after which this pod forgets a stale share
	// that no listing found again. It limits the size of staleSince.
	staleForgetAfter time.Duration
	// staleLastSweep is the time of the last sweep of staleSince.
	staleLastSweep time.Time
	// staleClock gives the time for the grace period. Nil means time.Now.
	staleClock func() time.Time

	gatewaySelector pool.Selectable[gatewayv1beta1.GatewayAPIClient]
	eventStream     events.Stream
}

// NewDefault returns a new manager instance with default dependencies
func NewDefault(m map[string]interface{}) (share.Manager, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error creating a new manager")
		return nil, err
	}

	s, err := metadata.NewCS3Storage(c.ProviderAddr, c.ProviderAddr, c.ServiceUserID, c.ServiceUserIdp, c.MachineAuthAPIKey)
	if err != nil {
		return nil, err
	}

	gatewaySelector, err := pool.GatewaySelector(c.GatewayAddr)
	if err != nil {
		return nil, err
	}

	var es events.Stream
	if c.Events.Endpoint != "" {
		es, err = stream.NatsFromConfig("jsoncs3-share-manager", false, stream.NatsConfig(c.Events))
		if err != nil {
			return nil, err
		}
	}

	return New(s, gatewaySelector, c.CacheTTL, es, c.MaxConcurrency)
}

// New returns a new manager instance.
func New(s metadata.Storage, gatewaySelector pool.Selectable[gatewayv1beta1.GatewayAPIClient], ttlSeconds int, es events.Stream, maxconcurrency int) (*Manager, error) {
	ttl := time.Duration(ttlSeconds) * time.Second

	m := &Manager{
		Cache:              providercache.New(s, ttl),
		CreatedCache:       sharecache.New(s, "users", "created.json", ttl),
		UserReceivedStates: receivedsharecache.New(s, ttl),
		GroupReceivedCache: sharecache.New(s, "groups", "received.json", ttl),
		storage:            s,
		gatewaySelector:    gatewaySelector,
		eventStream:        es,
		MaxConcurrency:     maxconcurrency,

		StaleShareGracePeriod: defaultStaleShareGracePeriod,
		staleSince:            map[string]staleSeen{},
		staleForgetAfter:      defaultStaleForgetAfter,
	}

	// listen for events
	if m.eventStream != nil {
		ch, err := events.Consume(m.eventStream, "jsoncs3sharemanager", _registeredEvents...)
		if err != nil {
			appctx.GetLogger(context.Background()).Error().Err(err).Msg("error consuming events")
		}
		go m.ProcessEvents(ch)
	}

	return m, nil
}

func (m *Manager) initialize(ctx context.Context) error {
	_, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "initialize")
	defer span.End()
	if m.initialized {
		span.SetStatus(codes.Ok, "already initialized")
		return nil
	}

	m.Lock()
	defer m.Unlock()

	if m.initialized { // check if initialization happened while grabbing the lock
		span.SetStatus(codes.Ok, "initialized while grabbing lock")
		return nil
	}

	ctx = context.Background()
	err := m.storage.Init(ctx, "jsoncs3-share-manager-metadata")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err = m.storage.MakeDirIfNotExist(ctx, "storages")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	err = m.storage.MakeDirIfNotExist(ctx, "users")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	err = m.storage.MakeDirIfNotExist(ctx, "groups")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	m.initialized = true
	span.SetStatus(codes.Ok, "initialized")
	return nil
}

func (m *Manager) ProcessEvents(ch <-chan events.Event) {
	log := logger.New()
	for event := range ch {
		// Put a logger into the context. Without it, the code below logs to a
		// disabled logger. Log only warnings and errors, because this logger
		// does not use the configured log level.
		ctxLog := log.Level(zerolog.WarnLevel)
		ctx := ctxLog.WithContext(context.Background())

		if err := m.initialize(ctx); err != nil {
			log.Error().Err(err).Msg("error initializing manager")
		}

		if ev, ok := event.Event.(events.SpaceDeleted); ok {
			log.Debug().Msgf("space deleted event: %v", ev)
			go func() {
				if err := m.purgeSpace(ctx, ev.ID); err != nil {
					log.Error().Err(err).Str("spaceid", ev.ID.GetOpaqueId()).Msg("could not remove all shares of the deleted space")
				}
			}()
		}
	}
}

// Share creates a new share
func (m *Manager) Share(ctx context.Context, md *provider.ResourceInfo, g *collaboration.ShareGrant) (*collaboration.Share, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Share")
	defer span.End()
	if err := m.initialize(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	user := ctxpkg.ContextMustGetUser(ctx)
	ts := utils.TSNow()

	// do not allow share to myself or the owner if share is for a user
	// TODO: should this not already be caught at the gw level?
	if g.Grantee.Type == provider.GranteeType_GRANTEE_TYPE_USER &&
		(utils.UserEqual(g.Grantee.GetUserId(), user.Id) || utils.UserEqual(g.Grantee.GetUserId(), md.Owner)) {
		err := errtypes.BadRequest("jsoncs3: owner/creator and grantee are the same")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// check if share already exists.
	key := &collaboration.ShareKey{
		// Owner:      md.Owner, owner no longer matters as it belongs to the space
		ResourceId: md.Id,
		Grantee:    g.Grantee,
	}

	_, err := m.getByKey(ctx, key)
	if err == nil {
		// share already exists
		err := errtypes.AlreadyExists(key.String())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	shareID := shareid.Encode(md.GetId().GetStorageId(), md.GetId().GetSpaceId(), uuid.NewString())
	s := &collaboration.Share{
		Id: &collaboration.ShareId{
			OpaqueId: shareID,
		},
		ResourceId:  md.Id,
		Permissions: g.Permissions,
		Grantee:     g.Grantee,
		Expiration:  g.Expiration,
		Owner:       md.Owner,
		Creator:     user.Id,
		Ctime:       ts,
		Mtime:       ts,
	}

	// Write the share first, then the entry for the grantee. When the share
	// cannot be stored, no grantee entry exists without a share.
	// A listing that sees the grantee entry then always finds the share.
	// When a write fails, remove what was written. Otherwise the share stays
	// without a grantee entry, and a new Share fails with AlreadyExists.
	spaceID := md.Id.StorageId + shareid.IDDelimiter + md.Id.SpaceId
	granteeWritten := false
	rollback := func() {
		log := appctx.GetLogger(ctx)
		// Each removal gets its own timeout. A slow removal must not use the
		// time of the next one.
		undo := func(what string, remove func(ctx context.Context) error) {
			ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), rollbackTimeout)
			defer cancel()
			if err := remove(ctx); err != nil {
				log.Error().Err(err).Str("shareid", shareID).Msg("could not roll back the " + what)
			}
		}
		// A failed grantee write can still reach the storage, for example when
		// the response times out. Remove the grantee entry first, so that no
		// grantee entry stays without a share.
		if granteeWritten {
			undo("grantee entry", func(ctx context.Context) error {
				switch g.Grantee.Type {
				case provider.GranteeType_GRANTEE_TYPE_USER:
					return m.UserReceivedStates.Remove(ctx, g.Grantee.GetUserId().GetOpaqueId(), spaceID, shareID)
				case provider.GranteeType_GRANTEE_TYPE_GROUP:
					return m.GroupReceivedCache.Remove(ctx, g.Grantee.GetGroupId().GetOpaqueId(), shareID)
				}
				return nil
			})
		}
		// The two removals do not depend on each other. Run them at the same
		// time, so that a failed Share does not wait for the sum of their timeouts.
		_ = runAll(
			func() error {
				undo("share in the provider cache", func(ctx context.Context) error {
					return m.Cache.Remove(ctx, md.Id.StorageId, md.Id.SpaceId, shareID)
				})
				return nil
			},
			func() error {
				undo("share in the created cache", func(ctx context.Context) error {
					return m.CreatedCache.Remove(ctx, s.GetCreator().GetOpaqueId(), shareID)
				})
				return nil
			},
		)
	}

	err = runAll(
		func() error { return m.Cache.Add(ctx, md.Id.StorageId, md.Id.SpaceId, shareID, s) },
		func() error { return m.CreatedCache.Add(ctx, s.GetCreator().GetOpaqueId(), shareID) },
	)
	if err == nil {
		granteeWritten = true
		// set flag for grantee to have access to share
		switch g.Grantee.Type {
		case provider.GranteeType_GRANTEE_TYPE_USER:
			userid := g.Grantee.GetUserId().GetOpaqueId()

			rs := &collaboration.ReceivedShare{
				Share: s,
				State: collaboration.ShareState_SHARE_STATE_PENDING,
			}
			err = m.UserReceivedStates.Add(ctx, userid, spaceID, rs)
		case provider.GranteeType_GRANTEE_TYPE_GROUP:
			groupid := g.Grantee.GetGroupId().GetOpaqueId()
			err = m.GroupReceivedCache.Add(ctx, groupid, shareID)
		}
	}
	if err != nil {
		rollback()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "")

	return s, nil
}

// getByID must be called in a lock-controlled block.
func (m *Manager) getByID(ctx context.Context, id *collaboration.ShareId) (*collaboration.Share, error) {
	storageID, spaceID, _ := shareid.Decode(id.OpaqueId)

	share, err := m.Cache.Get(ctx, storageID, spaceID, id.OpaqueId, false)
	if err != nil {
		return nil, err
	}
	if share == nil {
		return nil, errtypes.NotFound(id.String())
	}
	return share, nil
}

// getByKey must be called in a lock-controlled block.
func (m *Manager) getByKey(ctx context.Context, key *collaboration.ShareKey) (*collaboration.Share, error) {
	spaceShares, err := m.Cache.ListSpace(ctx, key.ResourceId.StorageId, key.ResourceId.SpaceId)
	if err != nil {
		return nil, err
	}
	for _, share := range spaceShares.Shares {
		if utils.GranteeEqual(key.Grantee, share.Grantee) && utils.ResourceIDEqual(share.ResourceId, key.ResourceId) {
			return share, nil
		}
	}
	return nil, errtypes.NotFound(key.String())
}

// get must be called in a lock-controlled block.
func (m *Manager) get(ctx context.Context, ref *collaboration.ShareReference) (s *collaboration.Share, err error) {
	switch {
	case ref.GetId() != nil:
		s, err = m.getByID(ctx, ref.GetId())
	case ref.GetKey() != nil:
		s, err = m.getByKey(ctx, ref.GetKey())
	default:
		err = errtypes.NotFound(ref.String())
	}
	return
}

// GetShare gets the information for a share by the given ref.
func (m *Manager) GetShare(ctx context.Context, ref *collaboration.ShareReference) (*collaboration.Share, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "GetShare")
	defer span.End()
	if err := m.initialize(ctx); err != nil {
		return nil, err
	}

	s, err := m.get(ctx, ref)
	if err != nil {
		return nil, err
	}
	if share.IsExpired(s) {
		m.expireShares(ctx, []*collaboration.Share{s})
	}
	// check if we are the creator or the grantee
	// TODO allow manager to get shares in a space created by other users
	user := ctxpkg.ContextMustGetUser(ctx)
	if share.IsCreatedByUser(s, user) || share.IsGrantedToUser(s, user) {
		return s, nil
	}

	req := &provider.StatRequest{
		Ref: &provider.Reference{ResourceId: s.ResourceId},
		FieldMask: &fieldmaskpb.FieldMask{
			Paths: []string{"permissions"},
		},
	}
	client, err := m.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}
	res, err := client.Stat(ctx, req)
	if err == nil &&
		res.Status.Code == rpcv1beta1.Code_CODE_OK &&
		res.Info.PermissionSet.ListGrants {
		return s, nil
	}

	// we return not found to not disclose information
	return nil, errtypes.NotFound(ref.String())
}

// Unshare deletes a share
func (m *Manager) Unshare(ctx context.Context, ref *collaboration.ShareReference) error {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "Unshare")
	defer span.End()

	if err := m.initialize(ctx); err != nil {
		return err
	}

	user := ctxpkg.ContextMustGetUser(ctx)

	s, err := m.get(ctx, ref)
	if err != nil {
		return err
	}
	// TODO allow manager to unshare shares in a space created by other users
	if !share.IsCreatedByUser(s, user) {
		// TODO why not permission denied?
		return errtypes.NotFound(ref.String())
	}

	return m.removeShare(ctx, s, false)
}

// UpdateShare updates the mode of the given share.
func (m *Manager) UpdateShare(ctx context.Context, ref *collaboration.ShareReference, p *collaboration.SharePermissions, updated *collaboration.Share, fieldMask *field_mask.FieldMask) (*collaboration.Share, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "UpdateShare")
	defer span.End()

	if err := m.initialize(ctx); err != nil {
		return nil, err
	}

	var toUpdate *collaboration.Share

	if ref != nil {
		var err error
		toUpdate, err = m.get(ctx, ref)
		if err != nil {
			return nil, err
		}
	} else if updated != nil {
		var err error
		toUpdate, err = m.getByID(ctx, updated.Id)
		if err != nil {
			return nil, err
		}
	}

	for _, fp := range fieldMask.GetPaths() {
		switch fp {
		case "permissions", "expiration":
		default:
			return nil, errtypes.NotSupported("updating " + fp + " is not supported")
		}
	}

	user := ctxpkg.ContextMustGetUser(ctx)
	if !share.IsCreatedByUser(toUpdate, user) {
		req := &provider.StatRequest{
			Ref: &provider.Reference{ResourceId: toUpdate.ResourceId},
			FieldMask: &fieldmaskpb.FieldMask{
				Paths: []string{"permissions"},
			},
		}
		client, err := m.gatewaySelector.Next()
		if err != nil {
			return nil, err
		}
		res, err := client.Stat(ctx, req)
		if err != nil ||
			res.Status.Code != rpcv1beta1.Code_CODE_OK ||
			!res.Info.PermissionSet.UpdateGrant {
			return nil, errtypes.NotFound(ref.String())
		}
	}

	// toUpdate is the share in the cache, and other requests can read it. Do not
	// change it. Update reads the share again under the space lock, changes a
	// copy and uploads the space. A share that a parallel request removed is
	// not written again.
	return m.Cache.Update(ctx, toUpdate.GetResourceId().GetStorageId(), toUpdate.GetResourceId().GetSpaceId(), toUpdate.GetId().GetOpaqueId(), func(s *collaboration.Share) {
		for _, fp := range fieldMask.GetPaths() {
			switch fp {
			case "permissions":
				s.Permissions = updated.GetPermissions()
			case "expiration":
				s.Expiration = updated.GetExpiration()
			}
		}
		if p != nil {
			s.Permissions = p
		}
		s.Mtime = utils.TSNow()
	})
}

// ListShares returns the shares created by the user
func (m *Manager) ListShares(ctx context.Context, filters []*collaboration.Filter) ([]*collaboration.Share, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "ListShares")
	defer span.End()

	if err := m.initialize(ctx); err != nil {
		return nil, err
	}

	user := ctxpkg.ContextMustGetUser(ctx)

	if len(share.FilterFiltersByType(filters, collaboration.Filter_TYPE_RESOURCE_ID)) > 0 {
		return m.listSharesByIDs(ctx, user, filters)
	}

	return m.listCreatedShares(ctx, user, filters)
}

func (m *Manager) listSharesByIDs(ctx context.Context, user *userv1beta1.User, filters []*collaboration.Filter) ([]*collaboration.Share, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "listSharesByIDs")
	defer span.End()
	sublog := appctx.GetLogger(ctx).With().Str("userid", user.GetId().GetOpaqueId()).Str("useridp", user.GetId().GetIdp()).Str("driver", "jsoncs3").Str("handler", "listSharesByIDs").Logger()

	providerSpaces := make(map[string]map[string]struct{})
	for _, f := range share.FilterFiltersByType(filters, collaboration.Filter_TYPE_RESOURCE_ID) {
		storageID := f.GetResourceId().GetStorageId()
		spaceID := f.GetResourceId().GetSpaceId()
		if providerSpaces[storageID] == nil {
			providerSpaces[storageID] = make(map[string]struct{})
		}
		providerSpaces[storageID][spaceID] = struct{}{}
	}

	statCache := make(map[string]struct{})
	var ss []*collaboration.Share
	// remove the expired shares of all spaces together, with one upload per document
	var expired []*collaboration.Share
	for providerID, spaces := range providerSpaces {
		for spaceID := range spaces {
			shares, err := m.Cache.ListSpace(ctx, providerID, spaceID)
			if err != nil {
				m.expireShares(ctx, expired)
				return nil, err
			}

			for _, s := range shares.Shares {
				resourceID := s.GetResourceId()
				// A new logger for each share. Do not add the fields to sublog, or
				// each log line carries the fields of all earlier shares.
				log := sublog.With().Str("storageid", resourceID.GetStorageId()).Str("spaceid", resourceID.GetSpaceId()).Str("opaqueid", resourceID.GetOpaqueId()).Logger()
				if share.IsExpired(s) {
					expired = append(expired, s)
					continue
				}
				if !share.MatchesFilters(s, filters) {
					continue
				}

				if !(share.IsCreatedByUser(s, user) || share.IsGrantedToUser(s, user)) {
					key := storagespace.FormatResourceID(resourceID)
					if _, hit := statCache[key]; !hit {
						req := &provider.StatRequest{
							Ref: &provider.Reference{ResourceId: resourceID},
							FieldMask: &fieldmaskpb.FieldMask{
								Paths: []string{"permissions"},
							},
						}
						client, err := m.gatewaySelector.Next()
						if err != nil {
							log.Error().Err(err).Msg("failed to select next gateway client")
							continue
						}
						res, err := client.Stat(ctx, req)
						if err != nil {
							log.Error().Err(err).Msg("failed to make stat call")
							continue
						}
						if res.Status.Code != rpcv1beta1.Code_CODE_OK {
							log.Debug().Str("code", res.GetStatus().GetCode().String()).Msg(res.GetStatus().GetMessage())
							continue
						}
						if !res.Info.PermissionSet.ListGrants {
							log.Debug().Msg("user has no list grants permission")
							continue
						}
						log.Debug().Msg("listing share for non participating user")
						statCache[key] = struct{}{}
					}
				}

				ss = append(ss, s)
			}
		}
	}
	m.expireShares(ctx, expired)
	span.SetStatus(codes.Ok, "")
	return ss, nil
}

func (m *Manager) listCreatedShares(ctx context.Context, user *userv1beta1.User, filters []*collaboration.Filter) ([]*collaboration.Share, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "listCreatedShares")
	defer span.End()
	sublog := appctx.GetLogger(ctx).With().Str("userid", user.GetId().GetOpaqueId()).Str("useridp", user.GetId().GetIdp()).Str("driver", "jsoncs3").Str("handler", "listCreatedShares").Logger()

	list, err := m.CreatedCache.List(ctx, user.Id.OpaqueId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	numWorkers := m.MaxConcurrency
	if numWorkers == 0 || len(list) < numWorkers {
		numWorkers = len(list)
	}

	type w struct {
		ssid string
		ids  sharecache.SpaceShareIDs
	}
	work := make(chan w)
	results := make(chan *collaboration.Share)
	// The workers collect the expired shares here. The expiration runs after
	// all workers stop, with one upload per document for all spaces.
	var (
		expiredMu     sync.Mutex
		expiredShares []*collaboration.Share
	)

	// gctx is cancelled when g.Wait returns. Do not use it after the wait.
	g, gctx := errgroup.WithContext(ctx)

	// Distribute work
	g.Go(func() error {
		defer close(work)
		for ssid, ids := range list {
			select {
			case work <- w{ssid, ids}:
			case <-gctx.Done():
				return gctx.Err()
			}
		}
		return nil
	})
	// Spawn workers that'll concurrently work the queue
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for w := range work {
				storageID, spaceID, _ := shareid.Decode(w.ssid)
				// fetch all shares from space with one request
				spaceShares, err := m.Cache.ListSpace(gctx, storageID, spaceID)
				if err != nil {
					sublog.Error().Err(err).
						Str("storageid", storageID).
						Str("spaceid", spaceID).
						Msg("failed to list shares in space")
					continue
				}
				for shareID := range w.ids.IDs {
					// Read the share from the snapshot. A parallel sync can reset
					// the live cache, and then the listing misses valid shares.
					s := spaceShares.Shares[shareID]
					if s == nil {
						continue
					}
					if share.IsExpired(s) {
						expiredMu.Lock()
						expiredShares = append(expiredShares, s)
						expiredMu.Unlock()
						continue
					}
					if utils.UserEqual(user.GetId(), s.GetCreator()) {
						if share.MatchesFilters(s, filters) {
							select {
							case results <- s:
							case <-gctx.Done():
								return gctx.Err()
							}
						}
					}
				}
			}
			return nil
		})
	}

	// Wait for things to settle down, then close the channel
	go func() {
		_ = g.Wait() // error is checked later
		close(results)
	}()

	ss := []*collaboration.Share{}
	for s := range results {
		ss = append(ss, s)
	}

	err = g.Wait()
	m.expireShares(ctx, expiredShares)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "")
	return ss, nil
}

// ListReceivedShares returns the list of shares the user has access to.
func (m *Manager) ListReceivedShares(ctx context.Context, filters []*collaboration.Filter, forUser *userv1beta1.UserId) ([]*collaboration.ReceivedShare, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "ListReceivedShares")
	defer span.End()
	sublog := appctx.GetLogger(ctx).With().Str("driver", "jsoncs3").Str("handler", "ListReceivedShares").Logger()

	if err := m.initialize(ctx); err != nil {
		return nil, err
	}

	user := ctxpkg.ContextMustGetUser(ctx)
	if user.GetId().GetType() == userv1beta1.UserType_USER_TYPE_SERVICE {
		client, err := m.gatewaySelector.Next()
		if err != nil {
			return nil, err
		}
		u, err := utils.GetUser(forUser, client)
		if err != nil {
			return nil, errtypes.BadRequest("user not found")
		}
		user = u
	}

	ssids := map[string]*receivedsharecache.Space{}
	// groupsOfShare holds the groups that give each share ID, for the cleanup of
	// stale group shares. The workers below only read it.
	groupsOfShare := map[string][]string{}

	// first collect all spaceids the user has access to as a group member
	for _, group := range user.Groups {
		list, err := m.GroupReceivedCache.List(ctx, group)
		if err != nil {
			continue // ignore error, cache will be updated on next read
		}
		for ssid, spaceShareIDs := range list {
			// add a pending entry, the state will be updated
			// when reading the received shares below if they have already been accepted or denied
			var rs *receivedsharecache.Space
			var ok bool
			if rs, ok = ssids[ssid]; !ok {
				rs = &receivedsharecache.Space{
					States: make(map[string]*receivedsharecache.State, len(spaceShareIDs.IDs)),
				}
				ssids[ssid] = rs
			}

			for shareid := range spaceShareIDs.IDs {
				rs.States[shareid] = &receivedsharecache.State{
					State: collaboration.ShareState_SHARE_STATE_PENDING,
				}
				groupsOfShare[shareid] = append(groupsOfShare[shareid], group)
			}
		}
	}

	// add all spaces the user has receved shares for, this includes mount points and share state for groups
	spaces, err := m.UserReceivedStates.List(ctx, user.Id.OpaqueId)
	if err != nil {
		return nil, err
	}
	for ssid, rspace := range spaces {
		if rs, ok := ssids[ssid]; ok {
			for shareid, state := range rspace.States {
				// overwrite state
				rs.States[shareid] = state
			}
		} else {
			ssids[ssid] = rspace
		}
	}

	numWorkers := m.MaxConcurrency
	if numWorkers == 0 || len(ssids) < numWorkers {
		numWorkers = len(ssids)
	}

	type w struct {
		ssid   string
		rspace *receivedsharecache.Space
	}
	work := make(chan w)
	results := make(chan *collaboration.ReceivedShare)
	// The workers collect the expired shares here. The expiration runs after
	// all workers stop, with one upload per document for all spaces.
	var (
		expiredMu     sync.Mutex
		expiredShares []*collaboration.Share
	)

	// gctx is cancelled when g.Wait returns. Do not use it after the wait.
	g, gctx := errgroup.WithContext(ctx)

	// Distribute work
	g.Go(func() error {
		defer close(work)
		for ssid, rspace := range ssids {
			select {
			case work <- w{ssid, rspace}:
			case <-gctx.Done():
				return gctx.Err()
			}
		}
		return nil
	})

	// Spawn workers that'll concurrently work the queue
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for w := range work {
				storageID, spaceID, _ := shareid.Decode(w.ssid)
				sublogr := sublog.With().Str("storageid", storageID).Str("spaceid", spaceID).Logger()
				// fetch all shares from space with one request
				spaceShares, err := m.Cache.ListSpace(gctx, storageID, spaceID)
				if err != nil {
					sublogr.Error().Err(err).
						Msg("failed to list shares in space")
					continue
				}
				// An empty etag means that the download of the space file got NotFound.
				// That can be temporary, so it is no proof that a share is gone. A Stat
				// that also gets NotFound is the proof. A failed first Share in a space
				// leaves a grantee entry without a space file.
				staleIsProven := spaceShares.Etag != ""
				spaceFileChecked := false
				// remove the stale shares of the space together, with one upload per document
				var stale []string
				staleOfGroup := map[string][]string{}
				// cleaned holds the IDs of all stale shares that this worker removes.
				var cleaned []string
				for shareID, state := range w.rspace.States {
					// Read the share from the same snapshot as the etag. A parallel
					// sync can reset the live cache, and then a valid share looks stale.
					s := spaceShares.Shares[shareID]
					if s == nil && !staleIsProven && !spaceFileChecked {
						spaceFileChecked = true
						staleIsProven = m.Cache.SpaceFileMissing(gctx, storageID, spaceID)
					}
					if s == nil && !staleIsProven {
						sublogr.Debug().Str("shareid", shareID).Msg("share not found, but the space file can exist. skipping the cleanup")
						continue
					}
					if s == nil && !m.staleGracePeriodOver(shareID) {
						sublogr.Debug().Str("shareid", shareID).Msg("share not found. waiting for the grace period before the cleanup")
						continue
					}
					if s == nil {
						sublogr.Warn().Str("shareid", shareID).Msg("share not found. cleaning up")
						cleaned = append(cleaned, shareID)
						// Remove the ID from received.json only when the file has it. Otherwise
						// the user cache downloads the file again to find nothing to remove.
						if us := spaces[w.ssid]; us != nil {
							if _, inUserFile := us.States[shareID]; inUserFile {
								stale = append(stale, shareID)
							}
						}
						// A stale group share is in the group cache, not in received.json.
						for _, group := range groupsOfShare[shareID] {
							staleOfGroup[group] = append(staleOfGroup[group], shareID)
						}
						continue
					}
					m.forgetStale(shareID)
					if share.IsExpired(s) {
						expiredMu.Lock()
						expiredShares = append(expiredShares, s)
						expiredMu.Unlock()
						continue
					}

					if share.IsGrantedToUser(s, user) {
						if share.MatchesFiltersWithState(s, state.State, filters) {
							rs := &collaboration.ReceivedShare{
								Share:      s,
								State:      state.State,
								MountPoint: state.MountPoint,
								Hidden:     state.Hidden,
							}
							select {
							case results <- rs:
							case <-gctx.Done():
								return gctx.Err()
							}
						}
					}
				}
				// failed holds the stale shares that a removal did not remove.
				failed := map[string]bool{}
				if len(stale) > 0 {
					if err := m.UserReceivedStates.RemoveMany(gctx, user.Id.OpaqueId, w.ssid, stale); err != nil {
						sublogr.Error().Err(err).Strs("shareids", stale).Msg("could not remove stale shares from the user received states")
						for _, id := range stale {
							failed[id] = true
						}
					}
				}
				for group, ids := range staleOfGroup {
					if err := m.GroupReceivedCache.RemoveMany(gctx, group, ids); err != nil {
						sublogr.Error().Err(err).Strs("shareids", ids).Str("groupid", group).Msg("could not remove stale shares from the group cache")
						for _, id := range ids {
							failed[id] = true
						}
					}
				}
				// Forget a stale share only after all removals of it. After a failed
				// removal, the grace period is still over, and the next listing
				// tries the removal again at once.
				for _, id := range cleaned {
					if !failed[id] {
						m.forgetStale(id)
					}
				}
			}
			return nil
		})
	}

	// Wait for things to settle down, then close the channel
	go func() {
		_ = g.Wait() // error is checked later
		close(results)
	}()

	rss := []*collaboration.ReceivedShare{}
	for rs := range results {
		rss = append(rss, rs)
	}

	err = g.Wait()
	m.expireShares(ctx, expiredShares)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "")
	return rss, nil
}

// convert must be called in a lock-controlled block.
func (m *Manager) convert(ctx context.Context, userID string, s *collaboration.Share) *collaboration.ReceivedShare {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "convert")
	defer span.End()

	rs := &collaboration.ReceivedShare{
		Share: s,
		State: collaboration.ShareState_SHARE_STATE_PENDING,
	}

	storageID, spaceID, _ := shareid.Decode(s.Id.OpaqueId)

	state, err := m.UserReceivedStates.Get(ctx, userID, storageID+shareid.IDDelimiter+spaceID, s.Id.GetOpaqueId())
	if err == nil && state != nil {
		rs.State = state.State
		rs.MountPoint = state.MountPoint
		rs.Hidden = state.Hidden
	}
	return rs
}

// GetReceivedShare returns the information for a received share.
func (m *Manager) GetReceivedShare(ctx context.Context, ref *collaboration.ShareReference) (*collaboration.ReceivedShare, error) {
	if err := m.initialize(ctx); err != nil {
		return nil, err
	}

	return m.getReceived(ctx, ref)
}

func (m *Manager) getReceived(ctx context.Context, ref *collaboration.ShareReference) (*collaboration.ReceivedShare, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "getReceived")
	defer span.End()

	s, err := m.get(ctx, ref)
	if err != nil {
		return nil, err
	}
	user := ctxpkg.ContextMustGetUser(ctx)
	if user.GetId().GetType() != userv1beta1.UserType_USER_TYPE_SERVICE && !share.IsGrantedToUser(s, user) {
		return nil, errtypes.NotFound(ref.String())
	}
	if share.IsExpired(s) {
		m.expireShares(ctx, []*collaboration.Share{s})
	}
	return m.convert(ctx, user.Id.GetOpaqueId(), s), nil
}

// UpdateReceivedShare updates the received share with share state.
func (m *Manager) UpdateReceivedShare(ctx context.Context, receivedShare *collaboration.ReceivedShare, fieldMask *field_mask.FieldMask, forUser *userv1beta1.UserId) (*collaboration.ReceivedShare, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "UpdateReceivedShare")
	defer span.End()

	if err := m.initialize(ctx); err != nil {
		return nil, err
	}

	rs, err := m.getReceived(ctx, &collaboration.ShareReference{Spec: &collaboration.ShareReference_Id{Id: receivedShare.Share.Id}})
	if err != nil {
		return nil, err
	}

	for i := range fieldMask.Paths {
		switch fieldMask.Paths[i] {
		case "state":
			rs.State = receivedShare.State
		case "mount_point":
			rs.MountPoint = receivedShare.MountPoint
		case "hidden":
			rs.Hidden = receivedShare.Hidden
		default:
			return nil, errtypes.NotSupported("updating " + fieldMask.Paths[i] + " is not supported")
		}
	}

	// write back
	u := ctxpkg.ContextMustGetUser(ctx)
	uid := u.GetId().GetOpaqueId()
	if u.GetId().GetType() == userv1beta1.UserType_USER_TYPE_SERVICE {
		uid = forUser.GetOpaqueId()
	}

	err = m.UserReceivedStates.Add(ctx, uid, rs.Share.ResourceId.StorageId+shareid.IDDelimiter+rs.Share.ResourceId.SpaceId, rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func shareIsRoutable(share *collaboration.Share) bool {
	return strings.Contains(share.Id.OpaqueId, shareid.IDDelimiter)
}

func updateShareID(share *collaboration.Share) {
	share.Id.OpaqueId = shareid.Encode(share.ResourceId.StorageId, share.ResourceId.SpaceId, share.Id.OpaqueId)
}

// Load imports shares and received shares from channels (e.g. during migration)
func (m *Manager) Load(ctx context.Context, shareChan <-chan *collaboration.Share, receivedShareChan <-chan share.ReceivedShareWithUser) error {
	log := appctx.GetLogger(ctx)
	if err := m.initialize(ctx); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for s := range shareChan {
			if s == nil {
				continue
			}
			if !shareIsRoutable(s) {
				updateShareID(s)
			}
			if err := m.Cache.Add(context.Background(), s.GetResourceId().GetStorageId(), s.GetResourceId().GetSpaceId(), s.Id.OpaqueId, s); err != nil {
				log.Error().Err(err).Interface("share", s).Msg("error persisting share")
			} else {
				log.Debug().Str("storageid", s.GetResourceId().GetStorageId()).Str("spaceid", s.GetResourceId().GetSpaceId()).Str("shareid", s.Id.OpaqueId).Msg("imported share")
			}
			if err := m.CreatedCache.Add(ctx, s.GetCreator().GetOpaqueId(), s.Id.OpaqueId); err != nil {
				log.Error().Err(err).Interface("share", s).Msg("error persisting created cache")
			} else {
				log.Debug().Str("creatorid", s.GetCreator().GetOpaqueId()).Str("shareid", s.Id.OpaqueId).Msg("updated created cache")
			}
		}
		wg.Done()
	}()
	go func() {
		for s := range receivedShareChan {
			if s.ReceivedShare != nil {
				if !shareIsRoutable(s.ReceivedShare.GetShare()) {
					updateShareID(s.ReceivedShare.GetShare())
				}
				switch s.ReceivedShare.Share.Grantee.Type {
				case provider.GranteeType_GRANTEE_TYPE_USER:
					if err := m.UserReceivedStates.Add(context.Background(), s.ReceivedShare.GetShare().GetGrantee().GetUserId().GetOpaqueId(), s.ReceivedShare.GetShare().GetResourceId().GetSpaceId(), s.ReceivedShare); err != nil {
						log.Error().Err(err).Interface("received share", s).Msg("error persisting received share for user")
					} else {
						log.Debug().Str("userid", s.ReceivedShare.GetShare().GetGrantee().GetUserId().GetOpaqueId()).Str("spaceid", s.ReceivedShare.GetShare().GetResourceId().GetSpaceId()).Str("shareid", s.ReceivedShare.GetShare().Id.OpaqueId).Msg("updated received share userdata")
					}
				case provider.GranteeType_GRANTEE_TYPE_GROUP:
					if err := m.GroupReceivedCache.Add(context.Background(), s.ReceivedShare.GetShare().GetGrantee().GetGroupId().GetOpaqueId(), s.ReceivedShare.GetShare().GetId().GetOpaqueId()); err != nil {
						log.Error().Err(err).Interface("received share", s).Msg("error persisting received share to group cache")
					} else {
						log.Debug().Str("groupid", s.ReceivedShare.GetShare().GetGrantee().GetGroupId().GetOpaqueId()).Str("shareid", s.ReceivedShare.GetShare().Id.OpaqueId).Msg("updated received share group cache")
					}
				}
			}
		}
		wg.Done()
	}()
	wg.Wait()

	return nil
}

// purgeSpace removes all shares of a deleted space. It continues after an
// error, because the space is gone. It returns all errors, so the caller
// can report them.
func (m *Manager) purgeSpace(ctx context.Context, id *provider.StorageSpaceId) error {
	storageID, spaceID := storagespace.SplitStorageID(id.OpaqueId)

	shares, err := m.Cache.ListSpace(ctx, storageID, spaceID)
	if err != nil {
		return fmt.Errorf("could not list the shares in space %s: %w", id.GetOpaqueId(), err)
	}

	// Remove all shares together, so each grantee and creator document gets one upload.
	all := make([]*collaboration.Share, 0, len(shares.Shares))
	for _, s := range shares.Shares {
		all = append(all, s)
	}
	var errs []error
	if _, err := m.removeShares(ctx, all, true); err != nil {
		errs = append(errs, err)
	}

	if err := m.Cache.PurgeSpace(ctx, storageID, spaceID); err != nil {
		errs = append(errs, fmt.Errorf("could not purge space %s: %w", id.GetOpaqueId(), err))
	}
	return stderrors.Join(errs...)
}

// removeShare removes a share and its entries. Set purgingSpace when the
// caller removes the space entry itself, because the space is deleted.
func (m *Manager) removeShare(ctx context.Context, s *collaboration.Share, purgingSpace bool) error {
	_, err := m.removeShares(ctx, []*collaboration.Share{s}, purgingSpace)
	return err
}

// removeShares removes shares and their entries. It uploads each touched
// document once, also when it removes many shares from it.
// Set purgingSpace when the caller removes the space entries itself.
// It returns the IDs of the shares that this call removed from the space
// entries. A share that another caller removed first is not in the result.
// With purgingSpace, the result is empty.
func (m *Manager) removeShares(ctx context.Context, shares []*collaboration.Share, purgingSpace bool) ([]string, error) {
	if len(shares) == 0 {
		return nil, nil
	}
	ctx, span := appctx.GetTracerProvider(ctx).Tracer(tracerName).Start(ctx, "removeShares")
	defer span.End()

	type userSpace struct{ userID, ssid string }
	type space struct{ storageID, spaceID string }
	userShares := map[userSpace][]string{}
	groupShares := map[string][]string{}
	createdShares := map[string][]string{}
	spaceShares := map[space][]string{}
	for _, s := range shares {
		id := s.GetId().GetOpaqueId()
		switch s.GetGrantee().GetType() {
		case provider.GranteeType_GRANTEE_TYPE_USER:
			k := userSpace{s.GetGrantee().GetUserId().GetOpaqueId(), s.GetResourceId().GetStorageId() + shareid.IDDelimiter + s.GetResourceId().GetSpaceId()}
			userShares[k] = append(userShares[k], id)
		case provider.GranteeType_GRANTEE_TYPE_GROUP:
			gid := s.GetGrantee().GetGroupId().GetOpaqueId()
			groupShares[gid] = append(groupShares[gid], id)
		}
		createdShares[s.GetCreator().GetOpaqueId()] = append(createdShares[s.GetCreator().GetOpaqueId()], id)
		storageID, spaceID, _ := shareid.Decode(id)
		spaceShares[space{storageID, spaceID}] = append(spaceShares[space{storageID, spaceID}], id)
	}

	var granteeFns, createdFns, spaceFns []func() error
	for k, ids := range userShares {
		granteeFns = append(granteeFns, func() error {
			if err := m.UserReceivedStates.RemoveMany(ctx, k.userID, k.ssid, ids); err != nil {
				return fmt.Errorf("could not remove shares %v from the received states of user %s: %w", ids, k.userID, err)
			}
			return nil
		})
	}
	for gid, ids := range groupShares {
		granteeFns = append(granteeFns, func() error {
			if err := m.GroupReceivedCache.RemoveMany(ctx, gid, ids); err != nil {
				return fmt.Errorf("could not remove shares %v from the cache of group %s: %w", ids, gid, err)
			}
			return nil
		})
	}
	for creator, ids := range createdShares {
		createdFns = append(createdFns, func() error {
			if err := m.CreatedCache.RemoveMany(ctx, creator, ids); err != nil {
				return fmt.Errorf("could not remove shares %v from the created shares of %s: %w", ids, creator, err)
			}
			return nil
		})
	}
	var removedMu sync.Mutex
	var removed []string
	for sp, ids := range spaceShares {
		spaceFns = append(spaceFns, func() error {
			r, err := m.Cache.RemoveMany(ctx, sp.storageID, sp.spaceID, ids)
			if err != nil {
				return fmt.Errorf("could not remove shares %v from space %s: %w", ids, sp.spaceID, err)
			}
			removedMu.Lock()
			removed = append(removed, r...)
			removedMu.Unlock()
			return nil
		})
	}

	// Limit the parallel removals. Each one downloads and uploads a document,
	// and a purge of a large space can touch hundreds of documents.
	limitAll(m.removeConcurrency(), granteeFns, createdFns, spaceFns)

	var err error
	if purgingSpace {
		// The caller purges the space entries. The order of the other writes does not matter.
		err = runAll(append(granteeFns, createdFns...)...)
	} else {
		// Remove the entries for the grantees first, then the shares. Continue also
		// after an error: the gateway removes the storage grant in any case.
		// The next listing of a grantee removes a stale grantee entry.
		granteeErr := runAll(granteeFns...)
		createdErr := make(chan error, 1)
		go func() { createdErr <- runAll(createdFns...) }()
		// Only the removal of the shares decides the result. When the shares are
		// gone, a retry gets NotFound, so an error makes the caller see a failed
		// removal that it cannot repeat. A listing skips a stale created entry.
		err = runAll(spaceFns...)
		if leftover := stderrors.Join(granteeErr, <-createdErr); leftover != nil {
			appctx.GetLogger(ctx).Warn().Err(leftover).Msg("removed the shares, but could not remove all entries. the next listing skips or removes them")
		}
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	// A failed space does not undo the removal in the other spaces, so return
	// the removed IDs also with an error.
	return removed, err
}

// expireShares removes expired shares, with one upload for each touched
// document. It publishes a ShareExpired event only for the shares that it
// removed. Many listings can find the same expired share, and each event can
// send a notification.
func (m *Manager) expireShares(ctx context.Context, shares []*collaboration.Share) {
	if len(shares) == 0 {
		return
	}
	log := appctx.GetLogger(ctx)
	removedIDs, err := m.removeShares(ctx, shares, false)
	if err != nil {
		log.Error().Err(err).Int("count", len(shares)).Msg("failed to unshare expired shares")
	}
	removed := make(map[string]struct{}, len(removedIDs))
	for _, id := range removedIDs {
		removed[id] = struct{}{}
	}
	for _, s := range shares {
		if _, ok := removed[s.GetId().GetOpaqueId()]; !ok {
			continue // another caller removed it, or the removal failed
		}
		if err := events.Publish(ctx, m.eventStream, events.ShareExpired{
			ShareID:        s.GetId(),
			ShareOwner:     s.GetOwner(),
			ItemID:         s.GetResourceId(),
			ExpiredAt:      time.Unix(int64(s.GetExpiration().GetSeconds()), int64(s.GetExpiration().GetNanos())),
			GranteeUserID:  s.GetGrantee().GetUserId(),
			GranteeGroupID: s.GetGrantee().GetGroupId(),
		}); err != nil {
			log.Error().Err(err).Str("shareid", s.GetId().GetOpaqueId()).Msg("failed to publish share expired event")
		}
	}
}

// staleSeen holds the times when this pod found a stale share.
type staleSeen struct {
	first time.Time // the start of the grace period
	last  time.Time // the time to forget the share is counted from here
}

// staleGracePeriodOver reports if a listing can remove a stale share. It is
// true when this pod first found the share stale StaleShareGracePeriod ago.
// The first call starts the grace period. Each call sets the last find time
// for sweepStale. Call forgetStale only after the removal succeeds, so that a
// failed removal does not start a new grace period.
func (m *Manager) staleGracePeriodOver(shareID string) bool {
	if m.StaleShareGracePeriod <= 0 {
		return true
	}
	now := time.Now()
	if m.staleClock != nil {
		now = m.staleClock()
	}
	m.staleMu.Lock()
	defer m.staleMu.Unlock()
	if m.staleSince == nil {
		m.staleSince = map[string]staleSeen{}
	}
	seen, ok := m.staleSince[shareID]
	if !ok {
		m.sweepStale(now)
		m.staleSince[shareID] = staleSeen{first: now, last: now}
		return false
	}
	seen.last = now
	m.staleSince[shareID] = seen
	return now.Sub(seen.first) >= m.StaleShareGracePeriod
}

// sweepStale removes the entries that no listing on this pod found for
// staleForgetAfter. forgetStale cannot remove an entry when only other pods
// find the share again or remove it. Without the sweep, staleSince grows
// until the pod restarts. A sweep reads the full map, so it runs at most once
// in each staleForgetAfter/staleSweepsPerForget. The caller holds staleMu.
func (m *Manager) sweepStale(now time.Time) {
	forgetAfter := m.staleForgetAfter
	if forgetAfter <= 0 {
		forgetAfter = defaultStaleForgetAfter
	}
	if now.Sub(m.staleLastSweep) < forgetAfter/staleSweepsPerForget {
		return
	}
	m.staleLastSweep = now
	for id, s := range m.staleSince {
		if now.Sub(s.last) > forgetAfter {
			delete(m.staleSince, id)
		}
	}
}

// forgetStale stops the grace period of a share that a listing found again.
// Two short losses of a share then do not add up to one long loss.
func (m *Manager) forgetStale(shareID string) {
	m.staleMu.Lock()
	defer m.staleMu.Unlock()
	delete(m.staleSince, shareID)
}

// defaultRemoveConcurrency limits the parallel removals when MaxConcurrency is not set.
const defaultRemoveConcurrency = 10

// minRemoveConcurrency is the lowest limit for the parallel removals. Before
// the batched removal, each removal of one share ran its 3 document updates
// in parallel. A lower limit makes an unshare slower than before.
const minRemoveConcurrency = 3

func (m *Manager) removeConcurrency() int {
	if m.MaxConcurrency > 0 {
		return max(m.MaxConcurrency, minRemoveConcurrency)
	}
	return defaultRemoveConcurrency
}

// limitAll wraps the functions in place, so that at most n of them run at
// the same time. All lists share the same limit.
func limitAll(n int, lists ...[]func() error) {
	sem := make(chan struct{}, n)
	for _, fns := range lists {
		for i, fn := range fns {
			fns[i] = func() error {
				sem <- struct{}{}
				defer func() { <-sem }()
				return fn()
			}
		}
	}
}

// runAll runs the functions at the same time and returns all errors.
// It does not cancel a function when another one fails.
func runAll(fns ...func() error) error {
	errs := make([]error, len(fns))
	var wg sync.WaitGroup
	for i, fn := range fns {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = fn()
		}()
	}
	wg.Wait()
	return stderrors.Join(errs...)
}

func (m *Manager) CleanupStaleShares(ctx context.Context) {
	log := appctx.GetLogger(ctx)

	if err := m.initialize(ctx); err != nil {
		return
	}

	// list all shares
	providers, err := m.Cache.All(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error listing all shares")
		return
	}

	client, err := m.gatewaySelector.Next()
	if err != nil {
		log.Error().Err(err).Msg("could not get gateway client")
		return
	}

	// Collect the orphaned shares first, and remove them together after the
	// loops. Then each document gets one upload, and the loops do not change
	// the maps that they read.
	var orphaned []*collaboration.Share
	providers.Range(func(storage string, spaces *providercache.Spaces) bool {
		log.Info().Str("storage", storage).Interface("spaceCount", spaces.Spaces.Count()).Msg("checking storage")

		spaces.Spaces.Range(func(space string, shares *providercache.Shares) bool {
			log.Info().Str("storage", storage).Str("space", space).Interface("shareCount", len(shares.Shares)).Msg("checking space")

			for _, s := range shares.Shares {
				req := &provider.StatRequest{
					Ref: &provider.Reference{ResourceId: s.ResourceId, Path: "."},
				}
				res, err := client.Stat(ctx, req)
				if err != nil {
					log.Error().Err(err).Str("storage", storage).Str("space", space).Msg("could not stat shared resource")
					continue
				}
				if res.GetStatus().GetCode() == rpcv1beta1.Code_CODE_NOT_FOUND {
					log.Info().Str("storage", storage).Str("space", space).Str("shareid", s.GetId().GetOpaqueId()).Msg("shared resource does not exist anymore. cleaning up share")
					orphaned = append(orphaned, s)
				}
			}

			return true
		})

		return true
	})

	removed, err := m.removeShares(ctx, orphaned, false)
	if err != nil {
		log.Error().Err(err).Msg("could not remove all orphaned shares")
	}
	log.Info().Int("orphaned", len(orphaned)).Int("removed", len(removed)).Msg("cleaned up orphaned shares")
}
