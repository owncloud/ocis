package types

import (
	"fmt"
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

// LinkCreated converts a ShareCreated Event to an AuditEventShareCreated
func LinkCreated(ev events.LinkCreated) AuditEventShareCreated {
	uid := ev.Sharer.OpaqueId
	with, typ := "", "link"
	base := BasicAuditEvent(uid, formatTime(ev.CTime), MessageLinkCreated(uid, ev.ItemID.OpaqueId, ev.ShareID.OpaqueId), ActionShareCreated)
	return AuditEventShareCreated{
		AuditEventSharing: SharingAuditEvent("", ev.ItemID.OpaqueId, uid, base),
		ShareOwner:        uid,
		ShareWith:         with,
		ShareType:         typ,
		ExpirationDate:    formatTime(ev.Expiration),
		SharePass:         ev.PasswordProtected,
		Permissions:       ev.Permissions.String(),
		ShareToken:        ev.Token,

		// NOTE: those values are not in the event and can therefore not be filled at the moment
		ItemType: "",
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

// LinkUpdated converts a LinkUpdated event to an AuditEventShareUpdated
func LinkUpdated(ev events.LinkUpdated) AuditEventShareUpdated {
	uid := ev.Sharer.OpaqueId
	with, typ := "", "link"
	base := BasicAuditEvent(uid, formatTime(ev.CTime), MessageLinkUpdated(uid, ev.ShareID.OpaqueId, ev.FieldUpdated), updateType(ev.FieldUpdated))
	return AuditEventShareUpdated{
		AuditEventSharing: SharingAuditEvent(ev.ShareID.GetOpaqueId(), ev.ItemID.OpaqueId, uid, base),
		ShareOwner:        uid,
		ShareWith:         with,
		ShareType:         typ,
		Permissions:       ev.Permissions.Permissions.String(),
		ExpirationDate:    formatTime(ev.Expiration),
		SharePass:         ev.PasswordProtected,
		ShareToken:        ev.Token,

		// NOTE: those values are not in the event and can therefore not be filled at the moment
		ItemType: "",
	}
}

// ShareRemoved converts a ShareRemoved event to an AuditEventShareRemoved
func ShareRemoved(ev events.ShareRemoved) AuditEventShareRemoved {
	sid, uid, iid, with, typ := "", "", "", "", ""
	if ev.ShareID != nil {
		sid = ev.ShareID.GetOpaqueId()
	}

	if ev.ShareKey != nil {
		uid = ev.ShareKey.GetOwner().GetOpaqueId()
		iid = ev.ShareKey.GetResourceId().GetOpaqueId()
		with, typ = extractGrantee(ev.ShareKey.GetGrantee().GetUserId(), ev.ShareKey.GetGrantee().GetGroupId())
	}
	base := BasicAuditEvent(uid, "", MessageShareRemoved(uid, sid, iid), ActionShareRemoved)
	return AuditEventShareRemoved{
		AuditEventSharing: SharingAuditEvent(sid, iid, uid, base),
		ShareWith:         with,
		ShareType:         typ,

		// NOTE: those values are not in the event and can therefore not be filled at the moment
		ItemType: "",
	}
}

// LinkRemoved converts a LinkRemoved event to an AuditEventShareRemoved
func LinkRemoved(ev events.LinkRemoved) AuditEventShareRemoved {
	uid, sid, typ := "", "", "link"
	if ev.ShareID != nil {
		sid = ev.ShareID.GetOpaqueId()
	} else {
		sid = ev.ShareToken
	}

	base := BasicAuditEvent(uid, "", MessageLinkRemoved(sid), ActionShareRemoved)
	return AuditEventShareRemoved{
		AuditEventSharing: SharingAuditEvent(sid, "", uid, base),
		ShareWith:         "",
		ShareType:         typ,

		// NOTE: those values are not in the event and can therefore not be filled at the moment
		ItemType: "",
	}
}

// ReceivedShareUpdated converts a ReceivedShareUpdated event to an AuditEventReceivedShareUpdated
func ReceivedShareUpdated(ev events.ReceivedShareUpdated) AuditEventReceivedShareUpdated {
	uid := ev.Sharer.GetOpaqueId()
	sid := ev.ShareID.GetOpaqueId()
	with, typ := extractGrantee(ev.GranteeUserID, ev.GranteeGroupID)
	itemID := ev.ItemID.GetOpaqueId()

	msg, utype := "", ""
	switch ev.State {
	case "SHARE_STATE_ACCEPTED":
		msg = MessageShareAccepted(with, sid, uid)
		utype = ActionShareAccepted
	case "SHARE_STATE_DECLINED":
		msg = MessageShareDeclined(with, sid, uid)
		utype = ActionShareDeclined
	}
	base := BasicAuditEvent(with, formatTime(ev.MTime), msg, utype)
	return AuditEventReceivedShareUpdated{
		AuditEventSharing: SharingAuditEvent(sid, itemID, uid, base),
		ShareType:         typ,
		ShareWith:         with,

		// NOTE: those values are not in the event and can therefore not be filled at the moment
		ItemType: "",
	}
}

// LinkAccessed converts a LinkAccessed event to an AuditEventLinkAccessed
func LinkAccessed(ev events.LinkAccessed) AuditEventLinkAccessed {
	uid := ev.Sharer.OpaqueId
	base := BasicAuditEvent(uid, formatTime(ev.CTime), MessageLinkAccessed(ev.ShareID.GetOpaqueId(), true), ActionLinkAccessed)
	return AuditEventLinkAccessed{
		AuditEventSharing: SharingAuditEvent(ev.ShareID.GetOpaqueId(), ev.ItemID.OpaqueId, uid, base),
		ShareToken:        ev.Token,
		Success:           true,

		// NOTE: those values are not in the event and can therefore not be filled at the moment
		ItemType: "",
	}
}

// LinkAccessFailed converts a LinkAccessFailed event to an AuditEventLinkAccessed
func LinkAccessFailed(ev events.LinkAccessFailed) AuditEventLinkAccessed {
	base := BasicAuditEvent("", "", MessageLinkAccessed(ev.ShareID.GetOpaqueId(), false), ActionLinkAccessed)
	return AuditEventLinkAccessed{
		AuditEventSharing: SharingAuditEvent(ev.ShareID.GetOpaqueId(), "", "", base),
		ShareToken:        ev.Token,
		Success:           false,

		// NOTE: those values are not in the event and can therefore not be filled at the moment
		ItemType: "",
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
	if t == nil {
		return ""
	}
	return time.Unix(int64(t.Seconds), int64(t.Nanos)).UTC().Format(time.RFC3339)
}

func updateType(u string) string {
	switch {
	case u == "permissions":
		return ActionSharePermissionUpdated
	case u == "displayname":
		return ActionShareDisplayNameUpdated
	case u == "TYPE_PERMISSIONS":
		return ActionSharePermissionUpdated
	case u == "TYPE_DISPLAYNAME":
		return ActionShareDisplayNameUpdated
	case u == "TYPE_PASSWORD":
		return ActionSharePasswordUpdated
	case u == "TYPE_EXPIRATION":
		return ActionShareExpirationUpdated
	default:
		fmt.Println("Unknown update type", u)
		return ""
	}
}
