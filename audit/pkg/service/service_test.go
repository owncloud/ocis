package svc

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/audit/pkg/types"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/test-go/testify/require"

	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	rtypes "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
)

var testCases = []struct {
	Alias           string
	SystemEvent     interface{}
	CheckAuditEvent func(*testing.T, []byte)
}{
	{
		Alias: "ShareCreated - user",
		SystemEvent: events.ShareCreated{
			Sharer:         userID("sharing-userid"),
			GranteeUserID:  userID("beshared-userid"),
			GranteeGroupID: nil,
			ItemID:         resourceID("storage-1", "itemid-1"),
			CTime:          timestamp(0),
		},
		CheckAuditEvent: func(t *testing.T, b []byte) {
			ev := types.AuditEventShareCreated{}
			require.NoError(t, json.Unmarshal(b, &ev))

			// AuditEvent fields
			checkBaseAuditEvent(t, ev.AuditEvent, "sharing-userid", "1970-01-01T01:00:00+01:00", "user 'sharing-userid' shared file 'itemid-1' with 'beshared-userid'", "file_shared")
			// AuditEventSharing fields
			checkSharingAuditEvent(t, ev.AuditEventSharing, "itemid-1", "sharing-userid", "")
			// AuditEventShareCreated fields
			require.Equal(t, "", ev.ItemType)
			require.Equal(t, "", ev.ExpirationDate)
			require.Equal(t, false, ev.SharePass)
			//require.Equal(t, "stat:true ", ev.Permissions) // TODO: BUG! Should work
			require.Equal(t, "user", ev.ShareType)
			require.Equal(t, "beshared-userid", ev.ShareWith)
			require.Equal(t, "sharing-userid", ev.ShareOwner)
			require.Equal(t, "", ev.ShareToken)
		},
	}, {
		Alias: "ShareCreated - group",
		SystemEvent: events.ShareCreated{
			Sharer:         userID("sharing-userid"),
			GranteeUserID:  nil,
			GranteeGroupID: groupID("beshared-groupid"),
			ItemID:         resourceID("storage-1", "itemid-1"),
			CTime:          timestamp(10e8),
		},
		CheckAuditEvent: func(t *testing.T, b []byte) {
			ev := types.AuditEventShareCreated{}
			require.NoError(t, json.Unmarshal(b, &ev))

			// AuditEvent fields
			checkBaseAuditEvent(t, ev.AuditEvent, "sharing-userid", "2001-09-09T03:46:40+02:00", "user 'sharing-userid' shared file 'itemid-1' with 'beshared-groupid'", "file_shared")
			// AuditEventSharing fields
			checkSharingAuditEvent(t, ev.AuditEventSharing, "itemid-1", "sharing-userid", "")
			// AuditEventShareCreated fields
			require.Equal(t, "", ev.ItemType)
			require.Equal(t, "", ev.ExpirationDate)
			require.Equal(t, false, ev.SharePass)
			//require.Equal(t, "stat:true ", ev.Permissions) // TODO: BUG! Should work
			require.Equal(t, "group", ev.ShareType)
			require.Equal(t, "beshared-groupid", ev.ShareWith)
			require.Equal(t, "sharing-userid", ev.ShareOwner)
			require.Equal(t, "", ev.ShareToken)

		},
	}, {
		Alias: "ShareUpdated",
		SystemEvent: events.ShareUpdated{
			ShareID:        shareID("shareid"),
			Sharer:         userID("sharing-userid"),
			GranteeUserID:  nil,
			GranteeGroupID: groupID("beshared-groupid"),
			ItemID:         resourceID("storage-1", "itemid-1"),
			Permissions:    sharePermissions("stat", "get_quota"),
			MTime:          timestamp(10e8),
			Updated:        "permissions",
		},
		CheckAuditEvent: func(t *testing.T, b []byte) {
			ev := types.AuditEventShareUpdated{}
			require.NoError(t, json.Unmarshal(b, &ev))

			// AuditEvent fields
			checkBaseAuditEvent(t, ev.AuditEvent, "sharing-userid", "2001-09-09T03:46:40+02:00", "user 'sharing-userid' updated field 'permissions' of share 'shareid'", "share_permission_updated")
			// AuditEventSharing fields
			checkSharingAuditEvent(t, ev.AuditEventSharing, "itemid-1", "sharing-userid", "shareid")
			// AuditEventShareUpdated fields
			require.Equal(t, "", ev.ItemType)       // not implemented atm
			require.Equal(t, "", ev.ExpirationDate) // no expiration for shares
			require.Equal(t, false, ev.SharePass)
			require.Equal(t, "get_quota:true stat:true ", ev.Permissions)
			require.Equal(t, "group", ev.ShareType)
			require.Equal(t, "beshared-groupid", ev.ShareWith)
			require.Equal(t, "sharing-userid", ev.ShareOwner)
			require.Equal(t, "", ev.ShareToken) // token not filled for shares
		},
	}, {
		Alias: "LinkUpdated - permissions",
		SystemEvent: events.LinkUpdated{
			ShareID:           linkID("shareid"),
			Sharer:            userID("sharing-userid"),
			ItemID:            resourceID("storage-1", "itemid-1"),
			Permissions:       linkPermissions("stat"),
			CTime:             timestamp(10e8),
			DisplayName:       "link",
			Expiration:        timestamp(10e8 + 10e5),
			PasswordProtected: true,
			Token:             "token-123",
			FieldUpdated:      "permissions",
		},
		CheckAuditEvent: func(t *testing.T, b []byte) {
			ev := types.AuditEventShareUpdated{}
			require.NoError(t, json.Unmarshal(b, &ev))

			// AuditEvent fields
			checkBaseAuditEvent(t, ev.AuditEvent, "sharing-userid", "2001-09-09T03:46:40+02:00", "user 'sharing-userid' updated field 'permissions' of public link 'shareid'", "share_permission_updated")
			// AuditEventSharing fields
			checkSharingAuditEvent(t, ev.AuditEventSharing, "itemid-1", "sharing-userid", "shareid")
			// AuditEventShareUpdated fields
			require.Equal(t, "", ev.ItemType) // not implemented atm
			require.Equal(t, "2001-09-20T17:33:20+02:00", ev.ExpirationDate)
			require.Equal(t, true, ev.SharePass)
			require.Equal(t, "stat:true ", ev.Permissions)
			require.Equal(t, "link", ev.ShareType)
			require.Equal(t, "", ev.ShareWith) // not filled on links
			require.Equal(t, "sharing-userid", ev.ShareOwner)
			require.Equal(t, "token-123", ev.ShareToken)
		},
	}, {
		Alias: "ShareRemoved",
		SystemEvent: events.ShareRemoved{
			ShareID:  shareID("shareid"),
			ShareKey: nil,
		},
		CheckAuditEvent: func(t *testing.T, b []byte) {
			ev := types.AuditEventShareRemoved{}
			require.NoError(t, json.Unmarshal(b, &ev))

			// AuditEvent fields
			checkBaseAuditEvent(t, ev.AuditEvent, "", "", "share id:'shareid' uid:'' item-id:'' was removed", "file_unshared")
			// AuditEventSharing fields
			checkSharingAuditEvent(t, ev.AuditEventSharing, "", "", "shareid")
			// AuditEventShareUpdated fields
			require.Equal(t, "", ev.ItemType) // not implemented atm
			require.Equal(t, "", ev.ShareType)
			require.Equal(t, "", ev.ShareWith) // not filled on links
		},
	}, {
		Alias: "LinkRemoved - id",
		SystemEvent: events.LinkRemoved{
			ShareID:    linkID("shareid"),
			ShareToken: "",
		},
		CheckAuditEvent: func(t *testing.T, b []byte) {
			ev := types.AuditEventShareRemoved{}
			require.NoError(t, json.Unmarshal(b, &ev))

			// AuditEvent fields
			checkBaseAuditEvent(t, ev.AuditEvent, "", "", "public link id:'shareid' was removed", "file_unshared")
			// AuditEventSharing fields
			checkSharingAuditEvent(t, ev.AuditEventSharing, "", "", "shareid")
			// AuditEventShareUpdated fields
			require.Equal(t, "", ev.ItemType) // not implemented atm
			require.Equal(t, "link", ev.ShareType)
			require.Equal(t, "", ev.ShareWith) // not filled on links
		},
	}, {
		Alias: "LinkRemoved - token",
		SystemEvent: events.LinkRemoved{
			ShareID:    nil,
			ShareToken: "token-123",
		},
		CheckAuditEvent: func(t *testing.T, b []byte) {
			ev := types.AuditEventShareRemoved{}
			require.NoError(t, json.Unmarshal(b, &ev))

			// AuditEvent fields
			checkBaseAuditEvent(t, ev.AuditEvent, "", "", "public link id:'token-123' was removed", "file_unshared")
			// AuditEventSharing fields
			checkSharingAuditEvent(t, ev.AuditEventSharing, "", "", "token-123")
			// AuditEventShareUpdated fields
			require.Equal(t, "", ev.ItemType) // not implemented atm
			require.Equal(t, "link", ev.ShareType)
			require.Equal(t, "", ev.ShareWith) // not filled on links
		},
	}, {
		Alias: "Share accepted",
		SystemEvent: events.ReceivedShareUpdated{
			ShareID:        shareID("shareid"),
			ItemID:         resourceID("storageid-1", "itemid-1"),
			Permissions:    sharePermissions("get_quota"),
			GranteeUserID:  userID("beshared-userid"),
			GranteeGroupID: nil,
			Sharer:         userID("sharing-userid"),
			MTime:          timestamp(10e8),
			State:          "SHARE_STATE_ACCEPTED",
		},
		CheckAuditEvent: func(t *testing.T, b []byte) {
			ev := types.AuditEventReceivedShareUpdated{}
			require.NoError(t, json.Unmarshal(b, &ev))

			// AuditEvent fields
			checkBaseAuditEvent(t, ev.AuditEvent, "", "", "public link id:'token-123' was removed", "file_unshared")
			// AuditEventSharing fields
			checkSharingAuditEvent(t, ev.AuditEventSharing, "", "", "token-123")
			// AuditEventShareUpdated fields
			require.Equal(t, "", ev.ItemType) // not implemented atm
			require.Equal(t, "link", ev.ShareType)
			require.Equal(t, "", ev.ShareWith) // not filled on links
		},
	},
}

func TestAuditLogging(t *testing.T) {
	log := log.NewLogger()

	inch := make(chan interface{})
	defer close(inch)

	outch := make(chan []byte)
	defer close(outch)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go StartAuditLogger(ctx, inch, log, Marshal("json", log), func(b []byte) {
		outch <- b
	})

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Alias, func(t *testing.T) {
			inch <- tc.SystemEvent
			tc.CheckAuditEvent(t, <-outch)
		})
	}
}

func checkBaseAuditEvent(t *testing.T, ev types.AuditEvent, user string, time string, message string, action string) {
	require.Equal(t, "", ev.RemoteAddr) // not implemented atm
	require.Equal(t, user, ev.User)
	require.Equal(t, "", ev.URL)       // not implemented atm
	require.Equal(t, "", ev.Method)    // not implemented atm
	require.Equal(t, "", ev.UserAgent) // not implemented atm
	require.Equal(t, time, ev.Time)
	require.Equal(t, "admin_audit", ev.App)
	require.Equal(t, message, ev.Message)
	require.Equal(t, action, ev.Action)
	require.Equal(t, false, ev.CLI) // not implemented atm
	require.Equal(t, 1, ev.Level)
}

func checkSharingAuditEvent(t *testing.T, ev types.AuditEventSharing, itemID string, owner string, shareID string) {
	require.Equal(t, itemID, ev.FileID)
	require.Equal(t, owner, ev.Owner)
	require.Equal(t, "", ev.Path) // not implemented atm
	require.Equal(t, shareID, ev.ShareID)
}

func shareID(id string) *collaboration.ShareId {
	return &collaboration.ShareId{
		OpaqueId: id,
	}
}

func linkID(id string) *link.PublicShareId {
	return &link.PublicShareId{
		OpaqueId: id,
	}
}

func userID(id string) *user.UserId {
	return &user.UserId{
		OpaqueId: id,
		Idp:      "idp",
	}
}

func groupID(id string) *group.GroupId {
	return &group.GroupId{
		OpaqueId: id,
		Idp:      "idp",
	}
}

func resourceID(sid, oid string) *provider.ResourceId {
	return &provider.ResourceId{
		StorageId: sid,
		OpaqueId:  oid,
	}
}

func timestamp(seconds uint64) *rtypes.Timestamp {
	return &rtypes.Timestamp{
		Seconds: seconds,
		Nanos:   0,
	}
}

func sharePermissions(perms ...string) *collaboration.SharePermissions {
	return &collaboration.SharePermissions{
		Permissions: permissions(perms...),
	}
}

func linkPermissions(perms ...string) *link.PublicSharePermissions {
	return &link.PublicSharePermissions{
		Permissions: permissions(perms...),
	}
}
func permissions(permissions ...string) *provider.ResourcePermissions {
	perms := &provider.ResourcePermissions{}

	for _, p := range permissions {
		switch p {
		case "stat":
			perms.Stat = true
		case "get_path":
			perms.GetPath = true
		case "list_container":
			perms.ListContainer = true
		case "get_quota":
			perms.GetQuota = true

		}
	}

	return perms
}
