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
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	rtypes "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
)

var testCases = []struct {
	Alias           string
	SystemEvent     interface{}
	CheckAuditEvent func(*testing.T, []byte)
}{
	{
		Alias: "ShareCreated",
		SystemEvent: events.ShareCreated{
			Sharer: &user.UserId{
				OpaqueId: "sharing-userid",
				Idp:      "idp",
			},
			GranteeUserID: &user.UserId{
				OpaqueId: "beshared-userid",
				Idp:      "idp",
			},
			GranteeGroupID: &group.GroupId{},
			Sharee:         &provider.Grantee{},
			ItemID: &provider.ResourceId{
				StorageId: "storage-1",
				OpaqueId:  "itemid-1",
			},
			CTime: &rtypes.Timestamp{
				Seconds: 0,
				Nanos:   0,
			},
		},
		CheckAuditEvent: func(t *testing.T, b []byte) {
			ev := types.AuditEventShareCreated{}
			err := json.Unmarshal(b, &ev)
			require.NoError(t, err)

			require.Equal(t, ev.User, "sharing-userid")
			require.Equal(t, ev.ShareWith, "beshared-userid")
			require.Equal(t, ev.FileID, "itemid-1")
			require.Equal(t, ev.Time, "1970-01-01T01:00:00+01:00")

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

	for _, tc := range testCases {
		inch <- tc.SystemEvent
		tc.CheckAuditEvent(t, <-outch)
	}
}
