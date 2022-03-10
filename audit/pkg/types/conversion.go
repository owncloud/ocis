package types

import (
	"fmt"
	"time"

	"github.com/cs3org/reva/v2/pkg/events"
)

// actions
const (
	actionShareCreated = "file_shared"
)

// messages
const (
	messageShareCreated = "user '%s' shared file '%s' with '%s'"
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
func SharingAuditEvent(fileid string, uid string, base AuditEvent) AuditEventSharing {
	return AuditEventSharing{
		AuditEvent: base,
		FileID:     fileid,
		Owner:      uid,

		// NOTE: those values are not in the events and can therefore not be filled at the moment
		ShareID: "",
		Path:    "",
	}
}

// ShareCreated converts a ShareCreated Event to an AuditEventShareCreated
func ShareCreated(ev events.ShareCreated) AuditEventShareCreated {
	with := ""
	typ := ""
	if ev.GranteeUserID != nil && ev.GranteeUserID.OpaqueId != "" {
		with = ev.GranteeUserID.OpaqueId
		typ = "user"
	} else if ev.GranteeGroupID != nil && ev.GranteeGroupID.OpaqueId != "" {
		with = ev.GranteeGroupID.OpaqueId
		typ = "group"
	}
	uid := ev.Sharer.OpaqueId
	t := time.Unix(int64(ev.CTime.Seconds), int64(ev.CTime.Nanos)).Format(time.RFC3339)
	base := BasicAuditEvent(uid, t, fmt.Sprintf(messageShareCreated, uid, ev.ItemID.OpaqueId, with), actionShareCreated)
	return AuditEventShareCreated{
		AuditEventSharing: SharingAuditEvent(ev.ItemID.OpaqueId, uid, base),
		ShareOwner:        uid,
		ShareWith:         with,
		ShareType:         typ,

		// NOTE: those values are not in the events and can therefore not be filled at the moment
		ItemType:       "",
		ExpirationDate: "",
		SharePass:      false,
		Permissions:    "",
		ShareToken:     "",
	}
}
