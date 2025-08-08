package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	cs3permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	"github.com/go-chi/chi/v5"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockServiceHandler implements all required methods for testing
type MockServiceHandler struct{}

// Bundle service methods
func (m *MockServiceHandler) SaveBundle(ctx context.Context, req *settingssvc.SaveBundleRequest, resp *settingssvc.SaveBundleResponse) error {
	return nil
}

func (m *MockServiceHandler) GetBundle(ctx context.Context, req *settingssvc.GetBundleRequest, resp *settingssvc.GetBundleResponse) error {
	return nil
}

func (m *MockServiceHandler) ListBundles(ctx context.Context, req *settingssvc.ListBundlesRequest, resp *settingssvc.ListBundlesResponse) error {
	return nil
}

func (m *MockServiceHandler) AddSettingToBundle(ctx context.Context, req *settingssvc.AddSettingToBundleRequest, resp *settingssvc.AddSettingToBundleResponse) error {
	return nil
}

func (m *MockServiceHandler) RemoveSettingFromBundle(ctx context.Context, req *settingssvc.RemoveSettingFromBundleRequest, resp *emptypb.Empty) error {
	return nil
}

// Value service methods
func (m *MockServiceHandler) SaveValue(ctx context.Context, req *settingssvc.SaveValueRequest, resp *settingssvc.SaveValueResponse) error {
	return nil
}

func (m *MockServiceHandler) GetValue(ctx context.Context, req *settingssvc.GetValueRequest, resp *settingssvc.GetValueResponse) error {
	return nil
}

func (m *MockServiceHandler) ListValues(ctx context.Context, req *settingssvc.ListValuesRequest, resp *settingssvc.ListValuesResponse) error {
	return nil
}

func (m *MockServiceHandler) GetValueByUniqueIdentifiers(ctx context.Context, req *settingssvc.GetValueByUniqueIdentifiersRequest, resp *settingssvc.GetValueResponse) error {
	return nil
}

// Role service methods
func (m *MockServiceHandler) ListRoles(ctx context.Context, req *settingssvc.ListBundlesRequest, resp *settingssvc.ListBundlesResponse) error {
	return nil
}

func (m *MockServiceHandler) ListRoleAssignments(ctx context.Context, req *settingssvc.ListRoleAssignmentsRequest, resp *settingssvc.ListRoleAssignmentsResponse) error {
	return nil
}

func (m *MockServiceHandler) ListRoleAssignmentsFiltered(ctx context.Context, req *settingssvc.ListRoleAssignmentsFilteredRequest, resp *settingssvc.ListRoleAssignmentsResponse) error {
	return nil
}

func (m *MockServiceHandler) AssignRoleToUser(ctx context.Context, req *settingssvc.AssignRoleToUserRequest, resp *settingssvc.AssignRoleToUserResponse) error {
	return nil
}

func (m *MockServiceHandler) RemoveRoleFromUser(ctx context.Context, req *settingssvc.RemoveRoleFromUserRequest, resp *emptypb.Empty) error {
	return nil
}

// Permission service methods
func (m *MockServiceHandler) ListPermissions(ctx context.Context, req *settingssvc.ListPermissionsRequest, resp *settingssvc.ListPermissionsResponse) error {
	return nil
}

func (m *MockServiceHandler) ListPermissionsByResource(ctx context.Context, req *settingssvc.ListPermissionsByResourceRequest, resp *settingssvc.ListPermissionsByResourceResponse) error {
	return nil
}

func (m *MockServiceHandler) GetPermissionByID(ctx context.Context, req *settingssvc.GetPermissionByIDRequest, resp *settingssvc.GetPermissionByIDResponse) error {
	return nil
}

// CS3 permissions methods
func (m *MockServiceHandler) CheckPermission(ctx context.Context, req *cs3permissions.CheckPermissionRequest) (*cs3permissions.CheckPermissionResponse, error) {
	return &cs3permissions.CheckPermissionResponse{}, nil
}

// TestCurrentProtocGenMicrowebImplementation captures the current protoc-gen-microweb implementation
// for future reference after it gets replaced with manual handlers
func TestCurrentProtocGenMicrowebImplementation(t *testing.T) {
	mockHandler := &MockServiceHandler{}
	r := chi.NewRouter()

	// Test that manual HTTP handlers are registered without panicking
	assert.NotPanics(t, func() {
		settingssvc.RegisterBundleServiceWeb(r, mockHandler)
		settingssvc.RegisterValueServiceWeb(r, mockHandler)
		settingssvc.RegisterRoleServiceWeb(r, mockHandler)
		settingssvc.RegisterPermissionServiceWeb(r, mockHandler)
	})

	// Test all endpoints with minimal JSON - these are the current protoc-gen-microweb routes
	testCases := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		// Bundle endpoints (current protoc-gen-microweb routes)
		{name: "bundle-save", path: "/api/v0/settings/bundle-save", expectedStatus: http.StatusCreated},
		{name: "bundle-get", path: "/api/v0/settings/bundle-get", expectedStatus: http.StatusCreated},
		{name: "bundles-list", path: "/api/v0/settings/bundles-list", expectedStatus: http.StatusCreated},
		{name: "bundles-add-setting", path: "/api/v0/settings/bundles-add-setting", expectedStatus: http.StatusCreated},
		{name: "bundles-remove-setting", path: "/api/v0/settings/bundles-remove-setting", expectedStatus: http.StatusNoContent},

		// Value endpoints (current protoc-gen-microweb routes)
		{name: "values-save", path: "/api/v0/settings/values-save", expectedStatus: http.StatusCreated},
		{name: "values-get", path: "/api/v0/settings/values-get", expectedStatus: http.StatusCreated},
		{name: "values-list", path: "/api/v0/settings/values-list", expectedStatus: http.StatusCreated},
		{name: "values-get-by-unique-identifiers", path: "/api/v0/settings/values-get-by-unique-identifiers", expectedStatus: http.StatusCreated},

		// Role endpoints (current protoc-gen-microweb routes)
		{name: "roles-list", path: "/api/v0/settings/roles-list", expectedStatus: http.StatusCreated},
		{name: "assignments-list", path: "/api/v0/settings/assignments-list", expectedStatus: http.StatusCreated},
		{name: "assignments-list-filtered", path: "/api/v0/settings/assignments-list-filtered", expectedStatus: http.StatusCreated},
		{name: "assignments-add", path: "/api/v0/settings/assignments-add", expectedStatus: http.StatusCreated},
		{name: "assignments-remove", path: "/api/v0/settings/assignments-remove", expectedStatus: http.StatusNoContent},

		// Permission endpoints (current protoc-gen-microweb routes)
		{name: "permissions-list", path: "/api/v0/settings/permissions-list", expectedStatus: http.StatusCreated},
		{name: "permissions-list-by-resource", path: "/api/v0/settings/permissions-list-by-resource", expectedStatus: http.StatusCreated},
		{name: "permissions-get-by-id", path: "/api/v0/settings/permissions-get-by-id", expectedStatus: http.StatusCreated},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test with empty JSON object
			req := httptest.NewRequest("POST", tc.path, bytes.NewBufferString("{}"))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Should not panic
			assert.NotPanics(t, func() {
				r.ServeHTTP(w, req)
			})

			// Should return expected status code
			assert.Equal(t, tc.expectedStatus, w.Code, "Expected status %d for %s, got %d", tc.expectedStatus, tc.path, w.Code)
		})
	}
}

func TestEdgeCaseRequestPayloads(t *testing.T) {
	mockHandler := &MockServiceHandler{}
	r := chi.NewRouter()

	// Register manual HTTP handlers
	settingssvc.RegisterBundleServiceWeb(r, mockHandler)
	settingssvc.RegisterValueServiceWeb(r, mockHandler)
	settingssvc.RegisterRoleServiceWeb(r, mockHandler)
	settingssvc.RegisterPermissionServiceWeb(r, mockHandler)

	testCases := []struct {
		name           string
		path           string
		payload        string
		expectedStatus int
		description    string
	}{
		{
			name:           "invalid-json",
			path:           "/api/v0/settings/bundle-save",
			payload:        "{invalid json}",
			expectedStatus: http.StatusPreconditionFailed,
			description:    "Should return 412 PreconditionFailed for invalid JSON",
		},
		{
			name:           "empty-json-object",
			path:           "/api/v0/settings/bundle-save",
			payload:        "{}",
			expectedStatus: http.StatusCreated,
			description:    "Should handle empty JSON object without panicking",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tc.path, bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Should not panic
			assert.NotPanics(t, func() {
				r.ServeHTTP(w, req)
			}, "Should not panic for %s", tc.description)

			// Always show response for debugging
			t.Logf("Status: %d, Response body: %s", w.Code, w.Body.String())

			// Should return expected status code
			assert.Equal(t, tc.expectedStatus, w.Code, "Expected status %d for %s, got %d", tc.expectedStatus, tc.description, w.Code)
		})
	}
}

func TestProtobufJSONFormatHandling(t *testing.T) {
	mockHandler := &MockServiceHandler{}
	r := chi.NewRouter()

	// Register manual HTTP handlers
	settingssvc.RegisterBundleServiceWeb(r, mockHandler)
	settingssvc.RegisterValueServiceWeb(r, mockHandler)
	settingssvc.RegisterRoleServiceWeb(r, mockHandler)
	settingssvc.RegisterPermissionServiceWeb(r, mockHandler)

	// Test with new protobuf JSON format
	testCases := []struct {
		name           string
		path           string
		payload        string
		expectedStatus int
		description    string
	}{
		{
			name:           "new-format-bundle-save",
			path:           "/api/v0/settings/bundle-save",
			payload:        `{"bundle":{"id":"test-bundle","name":"Test Bundle","type":"TYPE_DEFAULT","extension":"test","display_name":"Test Bundle","resource":{"type":"TYPE_SYSTEM"},"settings":[]}}`,
			expectedStatus: http.StatusCreated,
			description:    "Should handle new protobuf JSON format correctly",
		},
		{
			name:           "new-format-value-save",
			path:           "/api/v0/settings/values-save",
			payload:        `{"value":{"id":"test-value","bundle_id":"test-bundle","setting_id":"test-setting","account_uuid":"test-account","string_value":"test"}}`,
			expectedStatus: http.StatusCreated,
			description:    "Should handle new protobuf JSON format for values correctly",
		},
		{
			name:           "new-format-permission",
			path:           "/api/v0/settings/permissions-list",
			payload:        `{"account_uuid":"test-account"}`,
			expectedStatus: http.StatusCreated,
			description:    "Should handle new protobuf JSON format for permissions correctly",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tc.path, bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Should not panic
			assert.NotPanics(t, func() {
				r.ServeHTTP(w, req)
			}, "Should not panic for %s", tc.description)

			// Should return expected status code
			assert.Equal(t, tc.expectedStatus, w.Code, "Expected status %d for %s, got %d", tc.expectedStatus, tc.description, w.Code)
		})
	}
}
