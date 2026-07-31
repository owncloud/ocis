package types

import (
	"testing"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/stretchr/testify/require"
)

// TestShareUpdateAction ensures every share/link update maps to a non-empty,
// countable audit action, including the case where reva only sets the
// UpdateMask (the single Updated field is deprecated and no longer populated).
func TestShareUpdateAction(t *testing.T) {
	tests := []struct {
		name    string
		updated string
		mask    []string
		want    string
	}{
		{"share permissions", "permissions", nil, ActionSharePermissionUpdated},
		{"link permissions", "TYPE_PERMISSIONS", nil, ActionSharePermissionUpdated},
		{"link password", "TYPE_PASSWORD", nil, ActionSharePasswordUpdated},
		{"share expiration", "expiration", nil, ActionShareExpirationUpdated},
		{"share displayname", "displayname", nil, ActionShareDisplayNameUpdated},
		{"deprecated empty, mask permissions", "", []string{"permissions"}, ActionSharePermissionUpdated},
		{"deprecated empty, mask expiration", "", []string{"expiration"}, ActionShareExpirationUpdated},
		{"empty updated and empty mask", "", nil, ActionShareUpdated},
		{"unknown field", "somethingelse", nil, ActionShareUpdated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shareUpdateAction(tt.updated, tt.mask)
			require.Equal(t, tt.want, got)
			require.NotEmpty(t, got, "audit action must never be empty")
		})
	}
}

// TestShareUpdatedDerivesActionFromUpdateMask reproduces the reported bug:
// reva populates UpdateMask (not the deprecated Updated field), so the audit
// event must still carry a specific action and a meaningful field in the
// message rather than an empty action and "updated field ''".
func TestShareUpdatedDerivesActionFromUpdateMask(t *testing.T) {
	ev := events.ShareUpdated{
		Sharer:        &user.UserId{OpaqueId: "sharer"},
		ShareID:       &collaboration.ShareId{OpaqueId: "shareid"},
		ItemID:        &provider.ResourceId{OpaqueId: "itemid"},
		GranteeUserID: &user.UserId{OpaqueId: "grantee"},
		Permissions:   &collaboration.SharePermissions{},
		MTime:         &types.Timestamp{Seconds: 1},
		// reva sets the mask; the deprecated Updated field stays empty
		UpdateMask: []string{"permissions"},
	}

	out := ShareUpdated(ev)
	require.Equal(t, ActionSharePermissionUpdated, out.Action)
	require.NotEmpty(t, out.Action)
	require.Contains(t, out.Message, "permissions")
}

// TestReceivedShareUpdatedStates ensures the received-share state changes map to
// non-empty actions, including the rejected case (CS3 uses SHARE_STATE_REJECTED,
// which previously did not match "SHARE_STATE_DECLINED") and a generic fallback.
func TestReceivedShareUpdatedStates(t *testing.T) {
	build := func(state string) events.ReceivedShareUpdated {
		return events.ReceivedShareUpdated{
			Sharer:        &user.UserId{OpaqueId: "sharer"},
			ShareID:       &collaboration.ShareId{OpaqueId: "shareid"},
			ItemID:        &provider.ResourceId{OpaqueId: "itemid"},
			GranteeUserID: &user.UserId{OpaqueId: "grantee"},
			MTime:         &types.Timestamp{Seconds: 1},
			State:         state,
		}
	}

	accepted := ReceivedShareUpdated(build("SHARE_STATE_ACCEPTED"))
	require.Equal(t, ActionShareAccepted, accepted.Action)

	rejected := ReceivedShareUpdated(build("SHARE_STATE_REJECTED"))
	require.Equal(t, ActionShareDeclined, rejected.Action)
	require.NotEmpty(t, rejected.Message)

	pending := ReceivedShareUpdated(build("SHARE_STATE_PENDING"))
	require.Equal(t, ActionShareStateChanged, pending.Action)
	require.NotEmpty(t, pending.Message)
}
