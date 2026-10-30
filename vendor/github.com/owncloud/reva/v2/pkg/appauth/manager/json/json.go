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
	"io"

	"os"
	"strings"
	"sync"
	"time"

	apppb "github.com/cs3org/go-cs3apis/cs3/auth/applications/v1beta1"
	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/mitchellh/mapstructure"
	"github.com/owncloud/reva/v2/pkg/appauth"
	"github.com/owncloud/reva/v2/pkg/appauth/manager/registry"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

func init() {
	registry.Register("json", New)
}

type config struct {
	File                    string `mapstructure:"file"`
	TokenStrength           int    `mapstructure:"token_strength"`
	PasswordHashCost        int    `mapstructure:"password_hash_cost"`
	KeepExpiredTokensOnLoad bool   `mapstructure:"keep_expired_tokens_on_load"`
}

type jsonManager struct {
	sync.RWMutex
	config *config
	// map[userid][password]AppPassword
	passwords map[string]map[string]*apppb.AppPassword
}

// New returns a new mgr.
func New(m map[string]interface{}) (appauth.Manager, error) {
	c, err := parseConfig(m)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a new manager")
	}

	c.init()

	// load or create file
	manager, err := loadOrCreate(c.File)
	if err != nil {
		return nil, errors.Wrap(err, "error loading the file containing the application passwords")
	}

	manager.config = c

	if len(manager.passwords) == 0 {
		return manager, nil
	}
	// Both steps run before the manager is shared, so no lock is needed.
	// Re-key buckets written by older versions under UserId.String(), then
	// purge expired tokens so they don't accumulate over time.
	changed := manager.migrateLegacyKeys()
	
	if !c.KeepExpiredTokensOnLoad {
		manager.purgeExpiredTokens()
		changed = true
	}
	if changed {
		// A failed save is not fatal: the in-memory state is correct and
		// the next start repeats the migration.
		_ = manager.save()
	}

	return manager, nil
}

// userKey returns the key of the bucket that holds the app passwords of a user.
// It does not use UserId.String(): that is the protobuf text format, and
// protobuf-go varies its whitespace between binaries, so the same user gets a
// different key after a rebuild and loses access to the stored app passwords.
func userKey(id *userpb.UserId) string {
	return id.GetIdp() + "|" + id.GetOpaqueId()
}

// migrateLegacyKeys re-keys the buckets that older versions stored under
// UserId.String(). It parses each such key back into a user id and moves the
// bucket under userKey. Buckets of the same user that differ only in whitespace
// are merged. It reports whether it changed anything.
// Must be called before the manager is shared (no lock needed).
func (mgr *jsonManager) migrateLegacyKeys() bool {
	migrated := make(map[string]map[string]*apppb.AppPassword, len(mgr.passwords))
	changed := false
	// Keys that are already stable are copied first, so that on a conflict
	// their tokens win over the tokens of a legacy bucket.
	var legacy []string
	for key, tokens := range mgr.passwords {
		if isLegacyKey(key) {
			legacy = append(legacy, key)
			continue
		}
		migrated[key] = tokens
	}
	for _, key := range legacy {
		id := &userpb.UserId{}
		if err := prototext.Unmarshal([]byte(key), id); err != nil {
			// Not a user id in text format, keep the bucket as it is.
			migrated[key] = mgr.passwords[key]
			continue
		}
		newKey := userKey(id)
		if migrated[newKey] == nil {
			migrated[newKey] = make(map[string]*apppb.AppPassword, len(mgr.passwords[key]))
		}
		for hash, pw := range mgr.passwords[key] {
			if _, exists := migrated[newKey][hash]; !exists {
				migrated[newKey][hash] = pw
			}
		}
		changed = true
	}
	if changed {
		mgr.passwords = migrated
	}
	return changed
}

// isLegacyKey reports whether key looks like the text form of a user id.
func isLegacyKey(key string) bool {
	return strings.Contains(key, "opaque_id:") || strings.Contains(key, "idp:")
}

func (c *config) init() {
	if c.File == "" {
		c.File = "/var/tmp/reva/appauth.json"
	}
	if c.TokenStrength == 0 {
		c.TokenStrength = 16
	}
	if c.PasswordHashCost == 0 {
		c.PasswordHashCost = 11
	}
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}
	return c, nil
}

func loadOrCreate(file string) (*jsonManager, error) {
	stat, err := os.Stat(file)
	if os.IsNotExist(err) || stat.Size() == 0 {
		if err = os.WriteFile(file, []byte("{}"), 0644); err != nil {
			return nil, errors.Wrapf(err, "error creating the file %s", file)
		}
	}

	fd, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening the file %s", file)
	}
	defer fd.Close()

	data, err := io.ReadAll(fd)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading the file %s", file)
	}

	m := &jsonManager{}
	if err = json.Unmarshal(data, &m.passwords); err != nil {
		return nil, errors.Wrapf(err, "error parsing the file %s", file)
	}

	if m.passwords == nil {
		m.passwords = make(map[string]map[string]*apppb.AppPassword)
	}

	return m, nil
}

func (mgr *jsonManager) GenerateAppPassword(ctx context.Context, scope map[string]*authpb.Scope, label string, expiration *typespb.Timestamp) (*apppb.AppPassword, error) {
	// Reject already-expired tokens outright.
	if expiration != nil && expiration.Seconds != 0 && expiration.Seconds < uint64(time.Now().Unix()) {
		return nil, errors.New("cannot create an already-expired app password")
	}

	token, err := password.Generate(mgr.config.TokenStrength, mgr.config.TokenStrength/2, 0, false, false)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new token")
	}
	tokenHashed, err := bcrypt.GenerateFromPassword([]byte(token), mgr.config.PasswordHashCost)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new token")
	}
	userID := ctxpkg.ContextMustGetUser(ctx).GetId()
	uKey := userKey(userID)
	ctime := now()

	password := string(tokenHashed)
	appPass := &apppb.AppPassword{
		Password:   password,
		TokenScope: scope,
		Label:      label,
		Expiration: expiration,
		Ctime:      ctime,
		Utime:      ctime,
		User:       userID,
	}
	mgr.Lock()
	defer mgr.Unlock()

	mgr.purgeExpiredUserTokens(uKey)

	// check if user has some previous password
	if _, ok := mgr.passwords[uKey]; !ok {
		mgr.passwords[uKey] = make(map[string]*apppb.AppPassword)
	}

	mgr.passwords[uKey][password] = appPass

	err = mgr.save()
	if err != nil {
		return nil, errors.Wrap(err, "error saving new token")
	}

	clonedAppPass := proto.Clone(appPass).(*apppb.AppPassword)
	clonedAppPass.Password = token
	return clonedAppPass, nil
}

func (mgr *jsonManager) ListAppPasswords(ctx context.Context) ([]*apppb.AppPassword, error) {
	uKey := userKey(ctxpkg.ContextMustGetUser(ctx).GetId())
	mgr.RLock()
	defer mgr.RUnlock()
	appPasswords := []*apppb.AppPassword{}
	for _, pw := range mgr.passwords[uKey] {
		appPasswords = append(appPasswords, pw)
	}
	return appPasswords, nil
}

func (mgr *jsonManager) InvalidateAppPassword(ctx context.Context, password string) error {
	uKey := userKey(ctxpkg.ContextMustGetUser(ctx).GetId())
	mgr.Lock()
	defer mgr.Unlock()

	// see if user has a list of passwords
	appPasswords, ok := mgr.passwords[uKey]
	if !ok || len(appPasswords) == 0 {
		return errtypes.NotFound("password not found")
	}

	if _, ok := appPasswords[password]; !ok {
		return errtypes.NotFound("password not found")
	}
	delete(mgr.passwords[uKey], password)

	// if user has 0 passwords, delete user key from state map
	if len(mgr.passwords[uKey]) == 0 {
		delete(mgr.passwords, uKey)
	}

	return mgr.save()
}

func (mgr *jsonManager) GetAppPassword(ctx context.Context, userID *userpb.UserId, password string) (*apppb.AppPassword, error) {
	// Phase 1: find the matching token under a read lock.
	// Expired tokens are skipped before the expensive bcrypt comparison so
	// that accumulated expired tokens do not slow down authentication.
	// A read lock allows concurrent GetAppPassword calls.
	uKey := userKey(userID)
	mgr.RLock()
	appPasswords, ok := mgr.passwords[uKey]
	if !ok {
		mgr.RUnlock()
		return nil, errtypes.NotFound("password not found")
	}

	nowSec := uint64(time.Now().Unix())
	var matchedHash string
	var matchedPw *apppb.AppPassword
	for hash, pw := range appPasswords {
		if isExpired(pw, nowSec) {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err == nil {
			matchedHash = hash
			matchedPw = pw
			break
		}
	}
	mgr.RUnlock()

	if matchedPw == nil {
		return nil, errtypes.NotFound("password not found")
	}

	// Phase 2: update last-used time under a write lock.
	// Between RUnlock and Lock, the token could have been invalidated by
	// another goroutine, so re-check that it still exists.
	mgr.Lock()
	defer mgr.Unlock()

	if current, ok := mgr.passwords[uKey][matchedHash]; ok {
		current.Utime = now()
		if err := mgr.save(); err != nil {
			return nil, errors.Wrap(err, "error saving file")
		}
		return current, nil
	}

	return nil, errtypes.NotFound("password not found")
}

func now() *typespb.Timestamp {
	return &typespb.Timestamp{Seconds: uint64(time.Now().Unix())}
}

func (mgr *jsonManager) save() error {
	data, err := json.Marshal(mgr.passwords)
	if err != nil {
		return errors.Wrap(err, "error encoding json file")
	}

	if err = os.WriteFile(mgr.config.File, data, 0644); err != nil {
		return errors.Wrapf(err, "error writing to file %s", mgr.config.File)
	}

	return nil
}

// isExpired returns true if the token has a non-zero expiration time that is in the past.
func isExpired(pw *apppb.AppPassword, nowSec uint64) bool {
	return pw.Expiration != nil && pw.Expiration.Seconds != 0 && pw.Expiration.Seconds < nowSec
}

// purgeExpiredUserTokens removes expired tokens for a single user.
// Must be called while holding the write lock.
//
// Deleting map entries during a range loop is safe in Go per the language spec:
// https://go.dev/ref/spec#For_range
// "If a map entry that has not yet been reached is removed during iteration,
// the corresponding iteration value will not be produced."
func (mgr *jsonManager) purgeExpiredUserTokens(uid string) {
	tokens, ok := mgr.passwords[uid]
	if !ok {
		return
	}
	nowSec := uint64(time.Now().Unix())
	for hash, pw := range tokens {
		if isExpired(pw, nowSec) {
			delete(tokens, hash)
		}
	}
	if len(tokens) == 0 {
		delete(mgr.passwords, uid)
	}
}

// purgeExpiredTokens removes expired tokens for all users.
// Must be called before the manager is shared (no lock needed).
//
// Deleting map entries during a range loop is safe in Go per the language spec:
// https://go.dev/ref/spec#For_range
// "If a map entry that has not yet been reached is removed during iteration,
// the corresponding iteration value will not be produced."
func (mgr *jsonManager) purgeExpiredTokens() {
	nowSec := uint64(time.Now().Unix())
	for uid, tokens := range mgr.passwords {
		for hash, pw := range tokens {
			if isExpired(pw, nowSec) {
				delete(tokens, hash)
			}
		}
		if len(tokens) == 0 {
			delete(mgr.passwords, uid)
		}
	}
}
