package vault

import (
	"context"
	"errors"
	"testing"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/reva/v2/pkg/permission"
	"github.com/owncloud/reva/v2/pkg/utils"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
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
