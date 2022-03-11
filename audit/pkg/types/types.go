package types

// AuditEvent is the basic audit event
type AuditEvent struct {
	RemoteAddr string // the remote client IP
	User       string // the UID of the user performing the action. Or "IP x.x.x.x.", "cron", "CLI", "unknown"
	URL        string // the process request URI
	Method     string // the HTTP request method
	UserAgent  string // the HTTP request user agent
	Time       string // the time of the event eg: 2018-05-08T08:26:00+00:00
	App        string // always 'admin_audit'
	Message    string // sentence explaining the action
	Action     string // unique action identifier eg: file_delete or public_link_created
	CLI        bool   // if the action was performed from the CLI
	Level      int    // the log level of the entry (usually 1 for audit events)
}

// AuditEventSharing is the basic audit event for shares
type AuditEventSharing struct {
	AuditEvent

	FileID  string // The file identifier for the item shared.
	Owner   string // The UID of the owner of the shared item.
	Path    string // The path to the shared item.
	ShareID string // The sharing identifier. (not available for public_link_accessed or when recipient unshares)
}

// AuditEventShareCreated is the event logged when a share is created
type AuditEventShareCreated struct {
	AuditEventSharing

	ItemType       string // file or folder
	ExpirationDate string // The text expiration date in format 'yyyy-mm-dd'
	SharePass      bool   // If the share is password protected.
	Permissions    string // The permissions string eg: "READ"
	ShareType      string // group user or link
	ShareWith      string // The UID or GID of the share recipient. (not available for public link)
	ShareOwner     string // The UID of the share owner.
	ShareToken     string // For link shares the unique token, else null
}

// AuditEventShareUpdated is the event logged when a share is updated
type AuditEventShareUpdated struct {
	AuditEventSharing

	ItemType       string // file or folder
	ExpirationDate string // The text expiration date in format 'yyyy-mm-dd'
	SharePass      bool   // If the share is password protected.
	Permissions    string // The permissions string eg: "READ"
	ShareType      string // group user or link
	ShareWith      string // The UID or GID of the share recipient. (not available for public link)
	ShareOwner     string // The UID of the share owner.
	ShareToken     string // For link shares the unique token, else null
}

// AuditEventShareRemoved is the event logged when a share is removed
type AuditEventShareRemoved struct {
	AuditEventSharing
	ItemType  string // file or folder
	ShareType string // group user or link
	ShareWith string // The UID or GID of the share recipient.
}
