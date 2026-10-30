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

package ocdav_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncloud/reva/v2/internal/http/services/owncloud/ocdav"
	"github.com/owncloud/reva/v2/internal/http/services/owncloud/ocdav/config"
	"github.com/owncloud/reva/v2/pkg/rhttp/global"
	"github.com/owncloud/reva/v2/pkg/storage/favorite/memory"
)

func newTestSvc(t *testing.T, enableEmptyLegacyWebdavResponse bool) global.Service {
	t.Helper()
	fm, err := memory.New(nil)
	if err != nil {
		t.Fatalf("memory favorite manager: %v", err)
	}
	cfg := &config.Config{
		EnableEmptyLegacyWebdavResponse: enableEmptyLegacyWebdavResponse,
	}
	cfg.Init()
	svc, err := ocdav.NewWith(cfg, fm, nil, nil, nil)
	if err != nil {
		t.Fatalf("ocdav.NewWith: %v", err)
	}
	return svc
}

func TestEmptyLegacyWebdavResponse_PropfindRoot_ReturnsMinimal207(t *testing.T) {
	svc := newTestSvc(t, true)

	for _, path := range []string{"/remote.php/webdav/", "/remote.php/webdav"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(ocdav.MethodPropfind, path, http.NoBody)
			rr := httptest.NewRecorder()
			svc.Handler().ServeHTTP(rr, req)

			if rr.Code != http.StatusMultiStatus {
				t.Errorf("expected 207, got %d", rr.Code)
			}
			if ct := rr.Header().Get("Content-Type"); ct != "application/xml; charset=utf-8" {
				t.Errorf("unexpected Content-Type: %q", ct)
			}
			if rr.Body.Len() == 0 {
				t.Error("expected non-empty XML body")
			}
		})
	}
}

func TestEmptyLegacyWebdavResponse_Disabled_PassesThrough(t *testing.T) {
	svc := newTestSvc(t, false)

	// When disabled, the request passes through to the normal handler.
	// Without a valid auth context / gateway it will fail, but it must NOT return 207.
	req := httptest.NewRequest(ocdav.MethodPropfind, "/remote.php/webdav/", http.NoBody)
	rr := httptest.NewRecorder()
	svc.Handler().ServeHTTP(rr, req)

	if rr.Code == http.StatusMultiStatus {
		t.Error("feature is disabled but got 207 — should have passed through")
	}
}

func TestEmptyLegacyWebdavResponse_NonRootPath_PassesThrough(t *testing.T) {
	svc := newTestSvc(t, true)

	// Sub-paths must NOT return the empty response even when enabled.
	req := httptest.NewRequest(ocdav.MethodPropfind, "/remote.php/webdav/Documents/", http.NoBody)
	rr := httptest.NewRecorder()
	svc.Handler().ServeHTTP(rr, req)

	if rr.Code == http.StatusMultiStatus {
		t.Error("sub-path should not return empty 207")
	}
}

func TestEmptyLegacyWebdavResponse_NonPropfindMethod_PassesThrough(t *testing.T) {
	svc := newTestSvc(t, true)

	// Only PROPFIND on the root is intercepted; other methods must pass through.
	req := httptest.NewRequest(http.MethodGet, "/remote.php/webdav/", http.NoBody)
	rr := httptest.NewRecorder()
	svc.Handler().ServeHTTP(rr, req)

	if rr.Code == http.StatusMultiStatus {
		t.Error("GET on webdav root should not return 207")
	}
}

func TestEmptyLegacyWebdavResponse_WithoutRemotePHP_PassesThrough(t *testing.T) {
	svc := newTestSvc(t, true)

	// /webdav/ without the remote.php prefix must NOT be intercepted.
	req := httptest.NewRequest(ocdav.MethodPropfind, "/webdav/", http.NoBody)
	rr := httptest.NewRecorder()
	svc.Handler().ServeHTTP(rr, req)

	if rr.Code == http.StatusMultiStatus {
		t.Error("/webdav/ without remote.php should not return empty 207")
	}
}
