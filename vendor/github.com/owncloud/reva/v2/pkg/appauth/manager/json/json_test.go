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

package json

import (
	"context"
	encjson "encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	apppb "github.com/cs3org/go-cs3apis/cs3/auth/applications/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"golang.org/x/crypto/bcrypt"
)

func newTestManager(t *testing.T) (*jsonManager, string) {
	t.Helper()
	dir := t.TempDir()
	file := filepath.Join(dir, "appauth.json")
	mgr, err := New(map[string]interface{}{
		"file":               file,
		"token_strength":     16,
		"password_hash_cost": 4, // low cost for fast tests
	})
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	return mgr.(*jsonManager), file
}

func testCtx(uid string) context.Context {
	user := &userpb.User{
		Id: &userpb.UserId{
			OpaqueId: uid,
			Idp:      "test",
		},
	}
	return ctxpkg.ContextSetUser(context.Background(), user)
}

func TestGetAppPassword_SkipsExpiredTokens(t *testing.T) {
	mgr, _ := newTestManager(t)
	ctx := testCtx("user1")
	userID := ctxpkg.ContextMustGetUser(ctx).GetId()

	// Create an expired token directly in the map.
	expiredPassword := "expired-secret"
	hash, err := bcrypt.GenerateFromPassword([]byte(expiredPassword), 4)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	mgr.passwords[userID.String()] = map[string]*apppb.AppPassword{
		string(hash): {
			Password: string(hash),
			User:     userID,
			Expiration: &typespb.Timestamp{
				Seconds: uint64(time.Now().Add(-1 * time.Hour).Unix()),
			},
			Ctime: now(),
			Utime: now(),
		},
	}

	// Should not find the expired token.
	_, err = mgr.GetAppPassword(ctx, userID, expiredPassword)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestGetAppPassword_ValidTokenWorks(t *testing.T) {
	mgr, _ := newTestManager(t)
	ctx := testCtx("user1")

	// Generate a real token.
	expiration := &typespb.Timestamp{
		Seconds: uint64(time.Now().Add(1 * time.Hour).Unix()),
	}
	appPass, err := mgr.GenerateAppPassword(ctx, nil, "test-token", expiration)
	if err != nil {
		t.Fatalf("GenerateAppPassword: %v", err)
	}

	userID := ctxpkg.ContextMustGetUser(ctx).GetId()
	result, err := mgr.GetAppPassword(ctx, userID, appPass.Password)
	if err != nil {
		t.Fatalf("GetAppPassword: %v", err)
	}
	if result.Label != "test-token" {
		t.Errorf("expected label 'test-token', got %q", result.Label)
	}
}

func TestGetAppPassword_NoExpirationNeverExpires(t *testing.T) {
	mgr, _ := newTestManager(t)
	ctx := testCtx("user1")

	// Token with nil expiration should never expire.
	appPass, err := mgr.GenerateAppPassword(ctx, nil, "no-expiry", nil)
	if err != nil {
		t.Fatalf("GenerateAppPassword: %v", err)
	}

	userID := ctxpkg.ContextMustGetUser(ctx).GetId()
	result, err := mgr.GetAppPassword(ctx, userID, appPass.Password)
	if err != nil {
		t.Fatalf("GetAppPassword should succeed for non-expiring token: %v", err)
	}
	if result.Label != "no-expiry" {
		t.Errorf("expected label 'no-expiry', got %q", result.Label)
	}
}

func TestGetAppPassword_ZeroExpirationNeverExpires(t *testing.T) {
	mgr, _ := newTestManager(t)
	ctx := testCtx("user1")

	// Token with Seconds == 0 should never expire.
	appPass, err := mgr.GenerateAppPassword(ctx, nil, "zero-expiry", &typespb.Timestamp{Seconds: 0})
	if err != nil {
		t.Fatalf("GenerateAppPassword: %v", err)
	}

	userID := ctxpkg.ContextMustGetUser(ctx).GetId()
	result, err := mgr.GetAppPassword(ctx, userID, appPass.Password)
	if err != nil {
		t.Fatalf("GetAppPassword should succeed for zero-expiry token: %v", err)
	}
	if result.Label != "zero-expiry" {
		t.Errorf("expected label 'zero-expiry', got %q", result.Label)
	}
}

func TestPurgeExpiredTokensOnLoad(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "appauth.json")

	// Write a JSON file with one expired and one valid token.
	validPassword := "valid-secret"
	expiredPassword := "expired-secret"
	validHash, _ := bcrypt.GenerateFromPassword([]byte(validPassword), 4)
	expiredHash, _ := bcrypt.GenerateFromPassword([]byte(expiredPassword), 4)

	userKey := "idp:opaqueid"
	passwords := map[string]map[string]*apppb.AppPassword{
		userKey: {
			string(validHash): {
				Password: string(validHash),
				Label:    "valid",
				Ctime:    now(),
				Utime:    now(),
				// No expiration — should survive purge.
			},
			string(expiredHash): {
				Password: string(expiredHash),
				Label:    "expired",
				Expiration: &typespb.Timestamp{
					Seconds: uint64(time.Now().Add(-1 * time.Hour).Unix()),
				},
				Ctime: now(),
				Utime: now(),
			},
		},
	}

	data, _ := marshalPasswords(passwords)
	if err := os.WriteFile(file, data, 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	mgr, err := New(map[string]interface{}{
		"file":               file,
		"token_strength":     16,
		"password_hash_cost": 4,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	jm := mgr.(*jsonManager)
	tokens := jm.passwords[userKey]
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token after purge, got %d", len(tokens))
	}
	for _, pw := range tokens {
		if pw.Label != "valid" {
			t.Errorf("expected remaining token to be 'valid', got %q", pw.Label)
		}
	}
}

func TestGenerateAppPassword_PurgesExpired(t *testing.T) {
	mgr, _ := newTestManager(t)
	ctx := testCtx("user1")
	userID := ctxpkg.ContextMustGetUser(ctx).GetId()

	// Insert an expired token directly.
	expiredHash, _ := bcrypt.GenerateFromPassword([]byte("old-secret"), 4)
	mgr.passwords[userID.String()] = map[string]*apppb.AppPassword{
		string(expiredHash): {
			Password: string(expiredHash),
			Label:    "expired",
			Expiration: &typespb.Timestamp{
				Seconds: uint64(time.Now().Add(-1 * time.Hour).Unix()),
			},
			Ctime: now(),
			Utime: now(),
		},
	}

	// Generate a new token — should purge the expired one.
	expiration := &typespb.Timestamp{
		Seconds: uint64(time.Now().Add(1 * time.Hour).Unix()),
	}
	_, err := mgr.GenerateAppPassword(ctx, nil, "new-token", expiration)
	if err != nil {
		t.Fatalf("GenerateAppPassword: %v", err)
	}

	tokens := mgr.passwords[userID.String()]
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token (expired purged), got %d", len(tokens))
	}
	for _, pw := range tokens {
		if pw.Label != "new-token" {
			t.Errorf("expected remaining token to be 'new-token', got %q", pw.Label)
		}
	}
}

func TestIsExpired(t *testing.T) {
	nowSec := uint64(time.Now().Unix())

	tests := []struct {
		name     string
		pw       *apppb.AppPassword
		expected bool
	}{
		{
			name:     "nil expiration",
			pw:       &apppb.AppPassword{},
			expected: false,
		},
		{
			name:     "zero seconds",
			pw:       &apppb.AppPassword{Expiration: &typespb.Timestamp{Seconds: 0}},
			expected: false,
		},
		{
			name:     "future expiration",
			pw:       &apppb.AppPassword{Expiration: &typespb.Timestamp{Seconds: nowSec + 3600}},
			expected: false,
		},
		{
			name:     "past expiration",
			pw:       &apppb.AppPassword{Expiration: &typespb.Timestamp{Seconds: nowSec - 3600}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExpired(tt.pw, nowSec); got != tt.expected {
				t.Errorf("isExpired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// marshalPasswords is a test helper to serialize the password map to JSON.
func marshalPasswords(passwords map[string]map[string]*apppb.AppPassword) ([]byte, error) {
	return encjson.Marshal(passwords)
}
