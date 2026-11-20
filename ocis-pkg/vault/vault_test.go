package vault

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/reva/v2/pkg/permission"
	"github.com/owncloud/reva/v2/pkg/utils"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestIsVaultSpace(t *testing.T) {
	tests := []struct {
		name    string
		spaceID string
		want    bool
	}{
		{"vault space", utils.VaultStorageProviderID + "$space-id", true},
		{"regular space", "storage-id$space-id", false},
		{"unparsable space id fails closed", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVaultSpace(tt.spaceID); got != tt.want {
				t.Errorf("IsVaultSpace(%q) = %v, want %v", tt.spaceID, got, tt.want)
			}
		})
	}
}

func TestHasVaultMode(t *testing.T) {
	tests := []struct {
		name    string
		res     *permissions.CheckPermissionResponse
		err     error
		want    bool
		wantErr bool
	}{
		{"ok grants", &permissions.CheckPermissionResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}}, nil, true, false},
		{"permission denied denies", &permissions.CheckPermissionResponse{Status: &rpc.Status{Code: rpc.Code_CODE_PERMISSION_DENIED}}, nil, false, false},
		{"internal status fails closed", &permissions.CheckPermissionResponse{Status: &rpc.Status{Code: rpc.Code_CODE_INTERNAL}}, nil, false, true},
		{"transport error fails closed", nil, errors.New("unavailable"), false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gwc := &cs3mocks.GatewayAPIClient{}
			gwc.On("CheckPermission", mock.Anything, mock.MatchedBy(func(req *permissions.CheckPermissionRequest) bool {
				return req.GetPermission() == permission.VaultMode && req.GetSubjectRef().GetUserId().GetOpaqueId() == "alice"
			})).Return(tt.res, tt.err)

			got, err := HasVaultMode(context.Background(), gwc, &user.UserId{OpaqueId: "alice"})
			if (err != nil) != tt.wantErr {
				t.Fatalf("HasVaultMode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("HasVaultMode() = %v, want %v", got, tt.want)
			}
			gwc.AssertExpectations(t)
		})
	}
}

func TestFilterVaultModeUsers(t *testing.T) {
	// alice and every "ok-" user hold the permission, bob does not, and the check for carol fails
	gwc := &cs3mocks.GatewayAPIClient{}
	gwc.On("CheckPermission", mock.Anything, mock.Anything).Return(func(_ context.Context, r *permissions.CheckPermissionRequest, _ ...grpc.CallOption) (*permissions.CheckPermissionResponse, error) {
		switch id := r.GetSubjectRef().GetUserId().GetOpaqueId(); {
		case id == "carol":
			return nil, errors.New("unavailable")
		case id == "bob":
			return &permissions.CheckPermissionResponse{Status: &rpc.Status{Code: rpc.Code_CODE_PERMISSION_DENIED}}, nil
		default:
			return &permissions.CheckPermissionResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}}, nil
		}
	})

	ids := []*user.UserId{{OpaqueId: "alice"}, {OpaqueId: "bob"}, {OpaqueId: "carol"}}
	for i := range 50 {
		ids = append(ids, &user.UserId{OpaqueId: fmt.Sprintf("ok-%02d", i)})
	}

	var failed []string
	got := FilterVaultModeUsers(context.Background(), gwc, ids, func(id *user.UserId, err error) {
		failed = append(failed, id.GetOpaqueId())
	})

	want := []string{"alice"}
	for i := range 50 {
		want = append(want, fmt.Sprintf("ok-%02d", i))
	}
	gotIDs := make([]string, 0, len(got))
	for _, id := range got {
		gotIDs = append(gotIDs, id.GetOpaqueId())
	}
	if !slices.Equal(gotIDs, want) {
		t.Errorf("FilterVaultModeUsers() = %v, want %v", gotIDs, want)
	}
	if !slices.Equal(failed, []string{"carol"}) {
		t.Errorf("onError called for %v, want [carol]", failed)
	}
	gwc.AssertNumberOfCalls(t, "CheckPermission", len(ids))
}
