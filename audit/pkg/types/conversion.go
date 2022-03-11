package types

import (
	"time"

	"github.com/cs3org/reva/v2/pkg/events"

	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
)

// BasicAuditEvent creates an AuditEvent from given values
func BasicAuditEvent(uid string, ctime string, msg string, action string) AuditEvent {
	return AuditEvent{
		User:    uid,
		Time:    ctime,
		App:     "admin_audit",
		Message: msg,
		Action:  action,
		Level:   1,

		// NOTE: those values are not in the events and can therefore not be filled at the moment
		RemoteAddr: "",
		URL:        "",
		Method:     "",
		UserAgent:  "",
		CLI:        false,
	}
}

// SharingAuditEvent creates an AuditEventSharing from given values
func SharingAuditEvent(shareid string, fileid string, uid string, base AuditEvent) AuditEventSharing {
	return AuditEventSharing{
		AuditEvent: base,
		FileID:     fileid,
		Owner:      uid,
		ShareID:    shareid,

		// NOTE: those values are not in the events and can therefore not be filled at the moment
		Path: "",
	}
}

// ShareCreated converts a ShareCreated Event to an AuditEventShareCreated
func ShareCreated(ev events.ShareCreated) AuditEventShareCreated {
	uid := ev.Sharer.OpaqueId
	with, typ := extractGrantee(ev.GranteeUserID, ev.GranteeGroupID)
	base := BasicAuditEvent(uid, formatTime(ev.CTime), MessageShareCreated(uid, ev.ItemID.OpaqueId, with), ActionShareCreated)
	return AuditEventShareCreated{
		AuditEventSharing: SharingAuditEvent("", ev.ItemID.OpaqueId, uid, base),
		ShareOwner:        uid,
		ShareWith:         with,
		ShareType:         typ,

		// NOTE: those values are not in the event and can therefore not be filled at the moment
		ItemType:       "",
		ExpirationDate: "",
		SharePass:      false,
		Permissions:    "",
		ShareToken:     "",
	}
}

// ShareUpdated converts a ShareUpdated event to an AuditEventShareUpdated
func ShareUpdated(ev events.ShareUpdated) AuditEventShareUpdated {
	uid := ev.Sharer.OpaqueId
	with, typ := extractGrantee(ev.GranteeUserID, ev.GranteeGroupID)
	base := BasicAuditEvent(uid, formatTime(ev.MTime), MessageShareUpdated(uid, ev.ShareID.OpaqueId, ev.Updated), updateType(ev.Updated))
	return AuditEventShareUpdated{
		AuditEventSharing: SharingAuditEvent(ev.ShareID.GetOpaqueId(), ev.ItemID.OpaqueId, uid, base),
		ShareOwner:        uid,
		ShareWith:         with,
		ShareType:         typ,
		Permissions:       ev.Permissions.Permissions.String(),

		// NOTE: those values are not in the event and can therefore not be filled at the moment
		ItemType:       "",
		ExpirationDate: "",
		SharePass:      false,
		ShareToken:     "",
	}
}

func extractGrantee(uid *user.UserId, gid *group.GroupId) (string, string) {
	switch {
	case uid != nil && uid.OpaqueId != "":
		return uid.OpaqueId, "user"
	case gid != nil && gid.OpaqueId != "":
		return gid.OpaqueId, "group"
	}

	return "", ""
}

func formatTime(t *types.Timestamp) string {
	return time.Unix(int64(t.Seconds), int64(t.Nanos)).Format(time.RFC3339)
}

func updateType(u string) string {
	switch {
	case u == "permissions":
		return ActionSharePermissionUpdated
	case u == "displayname":
		return ActionShareDisplayNameUpdated
	default:
		return ""
	}
}
