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

/*
   Sharing
*/

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

// AuditEventReceivedShareUpdated is the event logged when a share is accepted or declined
type AuditEventReceivedShareUpdated struct {
	AuditEventSharing
	ItemType  string // file or folder
	ShareType string // group user or link
	ShareWith string // The UID or GID of the share recipient.
}

// AuditEventLinkAccessed is the event logged when a link is accessed
type AuditEventLinkAccessed struct {
	AuditEventSharing
	ShareToken string // The share token.
	Success    bool   // If the request was successful.
	ItemType   string // file or folder
}

/*
   Files
*/

// AuditEventFiles is the basic audit event for files
type AuditEventFiles struct {
	AuditEvent

	Path   string // The full path to the create file.
	Owner  string // The UID of the owner of the file.
	FileID string // The newly created files identifier.
}

// AuditEventFileCreated is the event logged when a file is created
type AuditEventFileCreated struct {
	AuditEventFiles
}

// AuditEventFileRead is the event logged when a file is read (aka downloaded)
type AuditEventFileRead struct {
	AuditEventFiles
}

// AuditEventFileUpdated is the event logged when a file is updated
// TODO: How to differentiate between new uploads and new version uploads?
// FIXME: implement
type AuditEventFileUpdated struct {
	AuditEventFiles
}

// AuditEventFileDeleted is the event logged when a file is deleted (aka trashed)
type AuditEventFileDeleted struct {
	AuditEventFiles
}

// AuditEventFileCopied is the event logged when a file is copied
// TODO: copy is a download&upload for now. How to know it was a copy?
// FIXME: implement
type AuditEventFileCopied struct {
	AuditEventFiles
}

// AuditEventFileRenamed is the event logged when a file is renamed (moved)
type AuditEventFileRenamed struct {
	AuditEventFiles
}

// AuditEventFilePurged is the event logged when a file is purged (deleted from trashbin)
type AuditEventFilePurged struct {
	AuditEventFiles
}

// AuditEventFileRestored is the event logged when a file is restored (from trashbin)
type AuditEventFileRestored struct {
	AuditEventFiles
}

// AuditEventFileVersionRestored is the event logged when a file version is restored
type AuditEventFileVersionRestored struct {
	AuditEventFiles
}

// AuditEventFileVersionDeleted is the event logged when a file version is deleted
// TODO: is this even possible?
type AuditEventFileVersionDeleted struct {
	AuditEventFiles
}
