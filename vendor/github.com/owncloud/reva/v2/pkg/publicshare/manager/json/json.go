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

package json

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/mitchellh/mapstructure"
	"github.com/owncloud/reva/v2/pkg/appctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/publicshare"
	"github.com/owncloud/reva/v2/pkg/publicshare/manager/json/persistence"
	"github.com/owncloud/reva/v2/pkg/publicshare/manager/json/persistence/cs3"
	"github.com/owncloud/reva/v2/pkg/publicshare/manager/json/persistence/file"
	"github.com/owncloud/reva/v2/pkg/publicshare/manager/json/persistence/memory"
	"github.com/owncloud/reva/v2/pkg/publicshare/manager/registry"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/storage/utils/metadata"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// shareField and passwordField are the keys used to store a share's encoded
// proto and its bcrypt-hashed password in a db entry.
const (
	shareField    = "share"
	passwordField = "password"
)

func init() {
	registry.Register("json", NewFile)
	registry.Register("jsoncs3", NewCS3)
	registry.Register("jsonmemory", NewMemory)
}

// dbEntryString safely reads the string value stored under key in a db entry.
// Entries are read back from persistence as interface{} (map[string]interface{}
// with string values), so a corrupted or legacy record can hold an unexpected
// type; this reports that via ok instead of panicking.
func dbEntryString(v interface{}, key string) (string, bool) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return "", false
	}
	s, ok := m[key].(string)
	return s, ok
}

// NewFile returns a new filesystem public shares manager.
func NewFile(c map[string]interface{}) (publicshare.Manager, error) {
	conf := &fileConfig{}
	if err := mapstructure.Decode(c, conf); err != nil {
		return nil, err
	}

	conf.init()
	if conf.File == "" {
		conf.File = "/var/tmp/reva/publicshares"
	}

	p := file.New(conf.File)
	return New(conf.GatewayAddr, conf.SharePasswordHashCost, conf.JanitorRunInterval, conf.EnableExpiredSharesCleanup, p)
}

// NewMemory returns a new in-memory public shares manager.
func NewMemory(c map[string]interface{}) (publicshare.Manager, error) {
	conf := &commonConfig{}
	if err := mapstructure.Decode(c, conf); err != nil {
		return nil, err
	}

	conf.init()
	p := memory.New()

	return New(conf.GatewayAddr, conf.SharePasswordHashCost, conf.JanitorRunInterval, conf.EnableExpiredSharesCleanup, p)
}

// NewCS3 returns a new cs3 public shares manager.
func NewCS3(c map[string]interface{}) (publicshare.Manager, error) {
	conf := &cs3Config{}
	if err := mapstructure.Decode(c, conf); err != nil {
		return nil, err
	}

	conf.init()

	s, err := metadata.NewCS3Storage(conf.ProviderAddr, conf.ProviderAddr, conf.ServiceUserID, conf.ServiceUserIdp, conf.MachineAuthAPIKey)
	if err != nil {
		return nil, err
	}
	p := cs3.New(s)

	return New(conf.GatewayAddr, conf.SharePasswordHashCost, conf.JanitorRunInterval, conf.EnableExpiredSharesCleanup, p)
}

// defaultStatConcurrency bounds how many Stat RPCs ListPublicShares may have
// in flight at once while checking ListGrants on distinct foreign resources.
// 5 mirrors the default concurrency the decomposedfs share manager clamps to
// for the same kind of bounded fan-out (see
// pkg/storage/utils/decomposedfs/options/options.go): enough to make a dent
// in a large batch of distinct resources without opening so many concurrent
// Stat RPCs that the gateway itself becomes the bottleneck.
const defaultStatConcurrency = 5

// New returns a new public share manager instance
func New(gwAddr string, pwHashCost, janitorRunInterval int, enableCleanup bool, p persistence.Persistence) (publicshare.Manager, error) {
	janitorCtx, janitorCancel := context.WithCancel(context.Background())
	m := &manager{
		gatewayAddr:                gwAddr,
		mutex:                      &sync.RWMutex{},
		passwordHashCost:           pwHashCost,
		janitorRunInterval:         janitorRunInterval,
		enableExpiredSharesCleanup: enableCleanup,
		persistence:                p,
		maxConcurrency:             defaultStatConcurrency,
		janitorCtx:                 janitorCtx,
		janitorCancel:              janitorCancel,
		janitorDone:                make(chan struct{}),
	}

	go m.startJanitorRun()
	return m, nil
}

type commonConfig struct {
	GatewayAddr                string `mapstructure:"gateway_addr"`
	SharePasswordHashCost      int    `mapstructure:"password_hash_cost"`
	JanitorRunInterval         int    `mapstructure:"janitor_run_interval"`
	EnableExpiredSharesCleanup bool   `mapstructure:"enable_expired_shares_cleanup"`
}

type fileConfig struct {
	commonConfig `mapstructure:",squash"`

	File string `mapstructure:"file"`
}

type cs3Config struct {
	commonConfig `mapstructure:",squash"`

	ProviderAddr      string `mapstructure:"provider_addr"`
	ServiceUserID     string `mapstructure:"service_user_id"`
	ServiceUserIdp    string `mapstructure:"service_user_idp"`
	MachineAuthAPIKey string `mapstructure:"machine_auth_apikey"`
}

func (c *commonConfig) init() {
	if c.SharePasswordHashCost == 0 {
		c.SharePasswordHashCost = 11
	}
	if c.JanitorRunInterval == 0 {
		c.JanitorRunInterval = 3600 // 1 hour
	}
}

type manager struct {
	gatewayAddr string
	mutex       *sync.RWMutex
	persistence persistence.Persistence

	passwordHashCost           int
	janitorRunInterval         int
	enableExpiredSharesCleanup bool

	// maxConcurrency bounds how many Stat RPCs ListPublicShares may have in
	// flight at once while checking ListGrants on distinct foreign resources.
	maxConcurrency int

	// janitorCtx/janitorCancel let Close ask startJanitorRun's goroutine to
	// stop - including cancelling a cleanupExpiredShares run it may be in
	// the middle of - and janitorDone lets Close wait until that goroutine
	// has actually returned, instead of racing it on shutdown.
	janitorCtx    context.Context
	janitorCancel context.CancelFunc
	janitorDone   chan struct{}
}

var _ publicshare.ClosableManager = (*manager)(nil)

// init is called at the top of every public method to lazily initialize the
// persistence layer. It must not take m.mutex: persistence.Init is already
// idempotent and self-synchronized (it returns immediately once the
// persistence layer reports itself initialized), so wrapping it in the
// manager's write lock added nothing but contention - and because a pending
// sync.RWMutex writer blocks new readers, that contention serialized every
// concurrent call through this single point regardless of whether it needed
// a read or a write lock.
func (m *manager) init(ctx context.Context) error {
	return m.persistence.Init(ctx)
}

func (m *manager) startJanitorRun() {
	defer close(m.janitorDone)

	if !m.enableExpiredSharesCleanup {
		return
	}

	ticker := time.NewTicker(time.Duration(m.janitorRunInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.janitorCtx.Done():
			return
		case <-ticker.C:
			if err := m.cleanupExpiredShares(); err != nil {
				log.Err(err).Msg("publicShareJSONManager: error cleaning up expired shares")
			}
		}
	}
}

// Close stops the janitor and waits for any in-flight cleanupExpiredShares
// run to finish, bounded by ctx. It satisfies publicshare.ClosableManager.
func (m *manager) Close(ctx context.Context) error {
	m.janitorCancel()

	select {
	case <-m.janitorDone:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Dump exports public shares to channels (e.g. during migration)
func (m *manager) Dump(ctx context.Context, shareChan chan<- *publicshare.WithPassword) error {
	log := appctx.GetLogger(ctx)

	if err := m.init(ctx); err != nil {
		return err
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	db, err := m.persistence.Read(ctx)
	if err != nil {
		return err
	}

	for _, v := range db {
		var local publicshare.WithPassword
		if share, ok := dbEntryString(v, shareField); ok {
			if err := utils.UnmarshalJSONToProtoV1([]byte(share), &local.PublicShare); err != nil {
				log.Error().Err(err).Msg("error unmarshalling share")
			}
		} else {
			log.Error().Msg("error reading share entry: missing or invalid \"share\" field")
		}
		// password is optional: shares without password protection have no
		// password field, so a missing/invalid entry just means "no password".
		local.Password, _ = dbEntryString(v, passwordField)
		shareChan <- &local
	}

	return nil
}

// Load imports public shares and received shares from channels (e.g. during migration)
func (m *manager) Load(ctx context.Context, shareChan <-chan *publicshare.WithPassword) error {
	if err := m.init(ctx); err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	db, err := m.persistence.Read(ctx)
	if err != nil {
		return err
	}
	dbCopy := persistence.Copy(db)

	for ps := range shareChan {
		encShare, err := utils.MarshalProtoV1ToJSON(&ps.PublicShare)
		if err != nil {
			return err
		}

		dbCopy[ps.PublicShare.Id.GetOpaqueId()] = map[string]interface{}{
			shareField:    string(encShare),
			passwordField: ps.Password,
		}
	}
	return m.persistence.Write(ctx, dbCopy)
}

// CreatePublicShare adds a new entry to manager.shares
func (m *manager) CreatePublicShare(ctx context.Context, u *user.User, rInfo *provider.ResourceInfo, g *link.Grant) (*link.PublicShare, error) {
	if rInfo.GetId() == nil || rInfo.GetId().GetStorageId() == "" {
		return nil, errtypes.BadRequest("resource id is required to create a public share")
	}

	id := &link.PublicShareId{
		OpaqueId: utils.RandString(15),
	}

	tkn := utils.RandString(15)
	now := time.Now().UnixNano()

	displayName, ok := rInfo.ArbitraryMetadata.Metadata["name"]
	if !ok {
		displayName = tkn
	}

	quicklink, _ := strconv.ParseBool(rInfo.ArbitraryMetadata.Metadata["quicklink"])

	var passwordProtected bool
	password := g.Password
	if len(password) > 0 {
		h, err := bcrypt.GenerateFromPassword([]byte(password), m.passwordHashCost)
		if err != nil {
			return nil, errors.Wrap(err, "could not hash share password")
		}
		password = string(h)
		passwordProtected = true
	}

	createdAt := &typespb.Timestamp{
		Seconds: uint64(now / int64(time.Second)),
		Nanos:   uint32(now % int64(time.Second)),
	}

	s := &link.PublicShare{
		Id:                id,
		Owner:             rInfo.GetOwner(),
		Creator:           u.Id,
		ResourceId:        rInfo.Id,
		Token:             tkn,
		Permissions:       g.Permissions,
		Ctime:             createdAt,
		Mtime:             createdAt,
		PasswordProtected: passwordProtected,
		Expiration:        g.Expiration,
		DisplayName:       displayName,
		Quicklink:         quicklink,
	}

	ps := &publicShare{
		Password: password,
	}
	proto.Merge(&ps.PublicShare, s)

	encShare, err := utils.MarshalProtoV1ToJSON(&ps.PublicShare)
	if err != nil {
		return nil, err
	}

	if err := m.init(ctx); err != nil {
		return nil, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	db, err := m.persistence.Read(ctx)
	if err != nil {
		return nil, err
	}
	dbCopy := persistence.Copy(db)

	if _, ok := dbCopy[s.Id.GetOpaqueId()]; !ok {
		dbCopy[s.Id.GetOpaqueId()] = map[string]interface{}{
			shareField:    string(encShare),
			passwordField: ps.Password,
		}
	} else {
		return nil, errors.New("key already exists")
	}

	err = m.persistence.Write(ctx, dbCopy)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// UpdatePublicShare updates the public share
func (m *manager) UpdatePublicShare(ctx context.Context, u *user.User, req *link.UpdatePublicShareRequest) (*link.PublicShare, error) {
	log := appctx.GetLogger(ctx)
	share, err := m.GetPublicShare(ctx, u, req.Ref, false)
	if err != nil {
		return nil, errors.New("ref does not exist")
	}

	now := time.Now().UnixNano()
	var newPasswordEncoded string
	passwordChanged := false

	switch req.GetUpdate().GetType() {
	case link.UpdatePublicShareRequest_Update_TYPE_DISPLAYNAME:
		log.Debug().Str("json", "update display name").Msgf("from: `%v` to `%v`", share.DisplayName, req.Update.GetDisplayName())
		share.DisplayName = req.Update.GetDisplayName()
	case link.UpdatePublicShareRequest_Update_TYPE_PERMISSIONS:
		old, _ := json.Marshal(share.Permissions)
		new, _ := json.Marshal(req.Update.GetGrant().Permissions)

		if req.GetUpdate().GetGrant().GetPassword() != "" {
			passwordChanged = true
			h, err := bcrypt.GenerateFromPassword([]byte(req.Update.GetGrant().Password), m.passwordHashCost)
			if err != nil {
				return nil, errors.Wrap(err, "could not hash share password")
			}
			newPasswordEncoded = string(h)
			share.PasswordProtected = true
		}

		log.Debug().Str("json", "update grants").Msgf("from: `%v`\nto\n`%v`", old, new)
		share.Permissions = req.Update.GetGrant().GetPermissions()
	case link.UpdatePublicShareRequest_Update_TYPE_EXPIRATION:
		old, _ := json.Marshal(share.Expiration)
		new, _ := json.Marshal(req.Update.GetGrant().Expiration)
		log.Debug().Str("json", "update expiration").Msgf("from: `%v`\nto\n`%v`", old, new)
		share.Expiration = req.Update.GetGrant().Expiration
	case link.UpdatePublicShareRequest_Update_TYPE_PASSWORD:
		passwordChanged = true
		if req.Update.GetGrant().Password == "" {
			share.PasswordProtected = false
			newPasswordEncoded = ""
		} else {
			h, err := bcrypt.GenerateFromPassword([]byte(req.Update.GetGrant().Password), m.passwordHashCost)
			if err != nil {
				return nil, errors.Wrap(err, "could not hash share password")
			}
			newPasswordEncoded = string(h)
			share.PasswordProtected = true
		}
	default:
		return nil, fmt.Errorf("invalid update type: %v", req.GetUpdate().GetType())
	}

	share.Mtime = &typespb.Timestamp{
		Seconds: uint64(now / int64(time.Second)),
		Nanos:   uint32(now % int64(time.Second)),
	}

	if err := m.init(ctx); err != nil {
		return nil, err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	db, err := m.persistence.Read(ctx)
	if err != nil {
		return nil, err
	}
	dbCopy := persistence.Copy(db)

	encShare, err := utils.MarshalProtoV1ToJSON(share)
	if err != nil {
		return nil, err
	}

	data, ok := dbCopy[share.Id.OpaqueId].(map[string]interface{})
	if !ok {
		data = map[string]interface{}{}
	}

	if ok && passwordChanged {
		data[passwordField] = newPasswordEncoded
	}
	data[shareField] = string(encShare)

	dbCopy[share.Id.OpaqueId] = data

	err = m.persistence.Write(ctx, dbCopy)
	if err != nil {
		return nil, err
	}

	return share, nil
}

// GetPublicShare gets a public share either by ID or Token.
func (m *manager) GetPublicShare(ctx context.Context, u *user.User, ref *link.PublicShareReference, sign bool) (*link.PublicShare, error) {
	if err := m.init(ctx); err != nil {
		return nil, err
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if ref.GetToken() != "" {
		ps, pw, err := m.getByToken(ctx, ref.GetToken())
		if err != nil {
			return nil, errtypes.NotFound("no shares found by token")
		}
		if ps.PasswordProtected && sign {
			err := publicshare.AddSignature(ps, pw)
			if err != nil {
				return nil, err
			}
		}
		return ps, nil
	}

	db, err := m.persistence.Read(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range db {
		share, ok := dbEntryString(v, shareField)
		if !ok {
			continue
		}
		passDB, _ := dbEntryString(v, passwordField)

		var ps link.PublicShare
		if err := utils.UnmarshalJSONToProtoV1([]byte(share), &ps); err != nil {
			return nil, err
		}

		if ref.GetId().GetOpaqueId() == ps.Id.OpaqueId {
			if publicshare.IsExpired(&ps) {
				// actual deletion is left to the janitor (cleanupExpiredShares)
				return nil, errtypes.NotFound("no shares found by id:" + ref.GetId().String())
			}
			if ps.PasswordProtected && sign {
				err := publicshare.AddSignature(&ps, passDB)
				if err != nil {
					return nil, err
				}
			}
			return &ps, nil
		}

	}
	return nil, errtypes.NotFound("no shares found by id:" + ref.GetId().String())
}

// ListPublicShares retrieves all the shares on the manager that are valid.
//
// Visibility of a foreign share (one not created by the calling user) is
// decided by a per-resource Stat, exactly as it always was: ListGrants on the
// share's resource is the OR of every ACE from that resource up to the space
// root (see assemblePermissions in
// pkg/storage/utils/decomposedfs/node/permissions.go, which also
// short-circuits on deny grants), so it cannot be derived from any
// precomputed set of space or resource ids without risking either false
// negatives or a privilege escalation (OCISDEV-861). What this method bounds
// is the *cost* of that check: pass 1 collects the set of distinct resources
// referenced by foreign shares, and pass 2 stats each of them at most once,
// concurrently, within a time budget derived from the caller's own context
// deadline (see statBudgetContext). N links on M distinct resources thus
// costs at most M stats, not N, and the whole call is bounded by whatever
// deadline the caller supplied, regardless of how large M is. If the caller
// supplies no deadline, the stat fan-out is unbounded, matching pre-existing
// behaviour.
func (m *manager) ListPublicShares(ctx context.Context, u *user.User, filters []*link.ListPublicSharesRequest_Filter, sign bool) ([]*link.PublicShare, error) {
	if err := m.init(ctx); err != nil {
		return nil, err
	}

	m.mutex.RLock()

	// Ranging over db below happens after the lock is released, which is safe
	// because we never mutate it: writers copy before they mutate, and the
	// persistence layer publishes a new map instead of changing this one.
	db, err := m.persistence.Read(ctx)
	if err != nil {
		m.mutex.RUnlock()
		return nil, err
	}

	m.mutex.RUnlock()

	log := appctx.GetLogger(ctx)

	// Pass 1 (in-memory, no RPCs): decode every persisted share once, handle
	// expiry and filters exactly as before, and split the survivors into
	// shares the caller created (no permission check needed) and foreign
	// shares (which do need one). While doing so, collect the set of
	// distinct resources the foreign shares point at, keyed by
	// storagespace.FormatResourceID, so pass 2 can stat each of them exactly
	// once.
	ownShares := make([]*publicShare, 0)
	foreignShares := make([]*publicShare, 0)
	foreignResourceIDs := make(map[string]*provider.ResourceId)

	for _, v := range db {
		var local publicShare
		share, ok := dbEntryString(v, shareField)
		if !ok {
			log.Warn().Interface("entry", v).Msg("ListPublicShares: skipping entry with missing or invalid \"share\" field")
			continue
		}
		if err := utils.UnmarshalJSONToProtoV1([]byte(share), &local.PublicShare); err != nil {
			return nil, err
		}

		if publicshare.IsExpired(&local.PublicShare) {
			// actual deletion is left to the janitor (cleanupExpiredShares)
			continue
		}

		if !publicshare.MatchesFilters(&local.PublicShare, filters) {
			continue
		}

		if local.ResourceId == nil {
			log.Warn().
				Str("share_id", local.PublicShare.GetId().GetOpaqueId()).
				Str("share_token", local.Token).
				Msg("ListPublicShares: skipping share with nil resource_id")
			continue
		}

		if publicshare.IsCreatedByUser(&local.PublicShare, u) {
			ownShares = append(ownShares, &local)
			continue
		}

		foreignShares = append(foreignShares, &local)
		foreignResourceIDs[storagespace.FormatResourceID(local.ResourceId)] = local.ResourceId
	}

	// Pass 2 (bounded RPCs): stat each distinct foreign resource once,
	// concurrently, within a time budget. A caller who created every share
	// (or has no foreign shares surviving the filters) issues no RPC at all.
	var permitted map[string]bool
	if len(foreignResourceIDs) > 0 {
		client, err := pool.GetGatewayServiceClient(m.gatewayAddr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to list shares")
		}
		permitted = m.statForeignResources(ctx, u, client, foreignResourceIDs)
	}

	shares := make([]*link.PublicShare, 0, len(ownShares)+len(foreignShares))
	for _, local := range ownShares {
		if local.PublicShare.PasswordProtected && sign {
			if err := publicshare.AddSignature(&local.PublicShare, local.Password); err != nil {
				return nil, err
			}
		}
		shares = append(shares, &local.PublicShare)
	}
	for _, local := range foreignShares {
		// Any resource whose permission was never determined (e.g. because
		// the time budget ran out) is absent here and therefore excluded:
		// fail closed, never include a share whose permission is unknown.
		if !permitted[storagespace.FormatResourceID(local.ResourceId)] {
			continue
		}
		if local.PublicShare.PasswordProtected && sign {
			if err := publicshare.AddSignature(&local.PublicShare, local.Password); err != nil {
				return nil, err
			}
		}
		shares = append(shares, &local.PublicShare)
	}
	return shares, nil
}

// statForeignResources stats each of the given distinct resources at most
// once, using a bounded pool of at most m.maxConcurrency concurrent workers,
// and returns a map of storagespace.FormatResourceID -> whether the calling
// user may list grants on that resource.
//
// The whole operation is bounded by a time budget derived from the caller's
// own context deadline, minus a small margin (see statBudgetContext): if the
// budget runs out before every resource has been stated, statting stops, a
// single warning is logged naming how many resources were skipped, and the
// partial map is returned rather than letting the caller block until its own
// deadline cancels the whole request with code = Canceled (OCISDEV-861). If
// the caller supplies no deadline, no budget is imposed. Any resource not
// present in the returned map was never decided and must be treated as not
// permitted by the caller.
func (m *manager) statForeignResources(ctx context.Context, u *user.User, client gateway.GatewayAPIClient, resourceIDs map[string]*provider.ResourceId) map[string]bool {
	log := appctx.GetLogger(ctx)

	statCtx, cancel := m.statBudgetContext(ctx)
	defer cancel()

	results := newStatResults()

	numWorkers := m.maxConcurrency
	if numWorkers > len(resourceIDs) {
		numWorkers = len(resourceIDs)
	}
	if numWorkers < 1 {
		numWorkers = 1
	}

	type job struct {
		rid *provider.ResourceId
	}
	jobs := make(chan job)

	g, gctx := errgroup.WithContext(statCtx)

	// Distribute work. Stop feeding jobs once the budget runs out so workers
	// drain and exit instead of blocking forever on a full channel.
	g.Go(func() error {
		defer close(jobs)
		for _, rid := range resourceIDs {
			select {
			case jobs <- job{rid}:
			case <-gctx.Done():
				return nil
			}
		}
		return nil
	})

	// Spawn workers that concurrently work the queue, bounded by
	// numWorkers <= m.maxConcurrency concurrent Stat RPCs in flight.
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for j := range jobs {
				if gctx.Err() != nil {
					// Budget exhausted: stop statting. A resource left
					// undecided is simply absent from the returned map, so
					// the caller treats it as not permitted.
					continue
				}
				m.userCanListGrants(statCtx, client, results, j.rid)
			}
			return nil
		})
	}
	_ = g.Wait()

	result := results.snapshot()
	if skipped := len(resourceIDs) - len(result); skipped > 0 {
		log.Warn().
			Str("user_id", u.GetId().GetOpaqueId()).
			Int("resources_checked", len(result)).
			Int("resources_skipped", skipped).
			Msg("ListPublicShares: stat time budget exhausted before every resource could be checked, returned list may be incomplete")
	}
	return result
}

// statBudgetContext derives a child context bounding how long
// statForeignResources may spend statting resources. The budget is derived
// solely from the caller's own context deadline, minus a small margin so the
// rest of the request (decoding, filtering, signing) still has time to run
// before the caller's deadline fires: this method never imposes a bound of
// its own. If the incoming context has no deadline, the returned context has
// none either, and the stat fan-out is unbounded - the same as before this
// budget existed.
func (m *manager) statBudgetContext(ctx context.Context) (context.Context, context.CancelFunc) {
	const margin = 200 * time.Millisecond

	deadline, ok := ctx.Deadline()
	if !ok {
		return ctx, func() {}
	}

	budget := time.Until(deadline) - margin
	if budget < 0 {
		budget = 0
	}
	return context.WithTimeout(ctx, budget)
}

// statResults collects ListGrants answers for resources, keyed by
// storagespace.FormatResourceID. It is safe for concurrent use by the bounded
// worker pool in statForeignResources.
type statResults struct {
	mu   sync.Mutex
	data map[string]bool
}

func newStatResults() *statResults {
	return &statResults{data: make(map[string]bool)}
}

func (r *statResults) set(key string, allowed bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = allowed
}

// snapshot returns a copy of the results collected so far. Call only once no
// more writers are running (e.g. after an errgroup.Wait), or take a copy
// under the same lock discipline as set.
func (r *statResults) snapshot() map[string]bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]bool, len(r.data))
	for k, v := range r.data {
		out[k] = v
	}
	return out
}

// userCanListGrants reports whether the current user may list grants on the
// given resource and records the answer in results. The resource IDs are
// already deduplicated by the caller (see the foreignResourceIDs map built
// in ListPublicShares) before statForeignResources ever runs, so each
// resource is stated at most once and there is nothing to look up here
// beforehand.
func (m *manager) userCanListGrants(ctx context.Context, client gateway.GatewayAPIClient, results *statResults, rid *provider.ResourceId) bool {
	log := appctx.GetLogger(ctx)
	key := storagespace.FormatResourceID(rid)

	sRes, err := client.Stat(ctx, &provider.StatRequest{
		Ref:       &provider.Reference{ResourceId: rid},
		FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"permissions"}},
	})
	switch {
	case err != nil:
		log.Error().Err(err).Interface("resource_id", rid).Msg("ListShares: an error occurred during stat on the resource")
		results.set(key, false)
		return false
	case sRes.Status.Code == rpc.Code_CODE_NOT_FOUND:
		log.Debug().Str("message", sRes.Status.Message).Interface("status", sRes.Status).Interface("resource_id", rid).Msg("ListShares: Resource not found")
		results.set(key, false)
		return false
	case sRes.Status.Code != rpc.Code_CODE_OK:
		log.Error().Str("message", sRes.Status.Message).Interface("status", sRes.Status).Interface("resource_id", rid).Msg("ListShares: could not stat resource")
		results.set(key, false)
		return false
	}

	allowed := sRes.GetInfo().GetPermissionSet().GetListGrants()
	results.set(key, allowed)
	return allowed
}

func (m *manager) cleanupExpiredShares() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// Another goroutine may hold the lock for a while; once we get it, skip the work if Close already asked us to stop.
	if err := m.janitorCtx.Err(); err != nil {
		return err
	}
	// Deriving from m.janitorCtx (not context.Background()) means Close
	// cancels an in-flight run immediately instead of leaving it to run out its full timeout.
	ctx, cancel := context.WithTimeout(m.janitorCtx, 60*time.Second)
	defer cancel()

	if err := m.init(ctx); err != nil {
		return err
	}

	read, err := m.persistence.Read(ctx)
	if err != nil {
		return err
	}
	db := persistence.Copy(read)

	var changed bool
	for id, v := range db {
		share, ok := dbEntryString(v, shareField)
		if !ok {
			continue
		}

		var ps link.PublicShare
		if err := utils.UnmarshalJSONToProtoV1([]byte(share), &ps); err != nil {
			continue
		}

		if publicshare.IsExpired(&ps) {
			delete(db, id)
			changed = true
		}
	}

	if !changed {
		return nil
	}

	return m.persistence.Write(ctx, db)
}

// RevokePublicShare undocumented.
func (m *manager) RevokePublicShare(ctx context.Context, _ *user.User, ref *link.PublicShareReference) error {
	if err := m.init(ctx); err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.revokePublicShare(ctx, ref)
}

// revokePublicShare doesn't have a lock inside, ensure a lock before call
func (m *manager) revokePublicShare(ctx context.Context, ref *link.PublicShareReference) error {
	read, err := m.persistence.Read(ctx)
	if err != nil {
		return err
	}
	db := persistence.Copy(read)

	switch {
	case ref.GetId() != nil && ref.GetId().OpaqueId != "":
		if _, ok := db[ref.GetId().OpaqueId]; ok {
			delete(db, ref.GetId().OpaqueId)
		} else {
			return errors.New("reference does not exist")
		}
	case ref.GetToken() != "":
		share, _, err := m.getByToken(ctx, ref.GetToken())
		if err != nil {
			return err
		}
		delete(db, share.Id.OpaqueId)
	default:
		return errors.New("reference does not exist")
	}

	return m.persistence.Write(ctx, db)
}

// getByToken doesn't have a lock inside, ensure a lock before call
func (m *manager) getByToken(ctx context.Context, token string) (*link.PublicShare, string, error) {
	db, err := m.persistence.Read(ctx)
	if err != nil {
		return nil, "", err
	}

	for _, v := range db {
		share, ok := dbEntryString(v, shareField)
		if !ok {
			continue
		}

		var local link.PublicShare
		if err := utils.UnmarshalJSONToProtoV1([]byte(share), &local); err != nil {
			return nil, "", err
		}

		if local.Token == token {
			passDB, _ := dbEntryString(v, passwordField)
			return &local, passDB, nil
		}
	}

	return nil, "", fmt.Errorf("share with token: `%v` not found", token)
}

// GetPublicShareByToken gets a public share by its opaque token.
func (m *manager) GetPublicShareByToken(ctx context.Context, token string, auth *link.PublicShareAuthentication, sign bool) (*link.PublicShare, error) {
	if err := m.init(ctx); err != nil {
		return nil, err
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	db, err := m.persistence.Read(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range db {
		share, ok := dbEntryString(v, shareField)
		if !ok {
			continue
		}

		passDB, _ := dbEntryString(v, passwordField)
		var local link.PublicShare
		if err := utils.UnmarshalJSONToProtoV1([]byte(share), &local); err != nil {
			return nil, err
		}

		if local.Token == token {
			if publicshare.IsExpired(&local) {
				// actual deletion is left to the janitor (cleanupExpiredShares)
				break
			}

			if local.PasswordProtected {
				if publicshare.Authenticate(&local, passDB, auth) {
					if sign {
						err := publicshare.AddSignature(&local, passDB)
						if err != nil {
							return nil, err
						}
					}
					return &local, nil
				}

				return nil, errtypes.InvalidCredentials("json: invalid password")
			}
			return &local, nil
		}
	}

	return nil, errtypes.NotFound(fmt.Sprintf("share with token: `%v` not found", token))
}

type publicShare struct {
	link.PublicShare
	Password string `json:"password"`
}
