package types

import "fmt"

// short identifiers for audit actions
const (
	// Sharing
	ActionShareCreated            = "file_shared"
	ActionSharePermissionUpdated  = "share_permission_updated"
	ActionShareDisplayNameUpdated = "share_name_updated"
	ActionSharePasswordUpdated    = "share_password_updated"
	ActionShareExpirationUpdated  = "share_expiration_updated"
	ActionShareRemoved            = "file_unshared"
	ActionShareAccepted           = "share_accepted"
	ActionShareDeclined           = "share_declined"
	ActionLinkAccessed            = "public_link_accessed"

	// Files
	ActionFileCreated         = "file_create"
	ActionFileRead            = "file_read"
	ActionFileTrashed         = "file_delete"
	ActionFileRenamed         = "file_rename"
	ActionFilePurged          = "file_trash_delete"
	ActionFileRestored        = "file_trash_restore"
	ActionFileVersionRestored = "file_version_restore"

	// Spaces
	ActionSpaceCreated  = "space_created"
	ActionSpaceRenamed  = "space_renamed"
	ActionSpaceDisabled = "space_disabled"
	ActionSpaceEnabled  = "space_enabled"
	ActionSpaceDeleted  = "space_deleted"
)

// MessageShareCreated returns the human readable string that describes the action
func MessageShareCreated(sharer, item, grantee string) string {
	return fmt.Sprintf("user '%s' shared file '%s' with '%s'", sharer, item, grantee)
}

// MessageLinkCreated returns the human readable string that describes the action
func MessageLinkCreated(sharer, item, shareid string) string {
	return fmt.Sprintf("user '%s' created a public to file '%s' with id '%s'", sharer, item, shareid)
}

// MessageShareUpdated returns the human readable string that describes the action
func MessageShareUpdated(sharer, shareID, fieldUpdated string) string {
	return fmt.Sprintf("user '%s' updated field '%s' of share '%s'", sharer, fieldUpdated, shareID)
}

// MessageLinkUpdated returns the human readable string that describes the action
func MessageLinkUpdated(sharer, shareid, fieldUpdated string) string {
	return fmt.Sprintf("user '%s' updated field '%s' of public link '%s'", sharer, fieldUpdated, shareid)
}

// MessageShareRemoved returns the human readable string that describes the action
func MessageShareRemoved(sharer, shareid, itemid string) string {
	return fmt.Sprintf("share id:'%s' uid:'%s' item-id:'%s' was removed", shareid, sharer, itemid)
}

// MessageLinkRemoved returns the human readable string that describes the action
func MessageLinkRemoved(shareid string) string {
	return fmt.Sprintf("public link id:'%s' was removed", shareid)
}

// MessageShareAccepted returns the human readable string that describes the action
func MessageShareAccepted(userid, shareid, sharerid string) string {
	return fmt.Sprintf("user '%s' accepted share '%s' from user '%s'", userid, shareid, sharerid)
}

// MessageShareDeclined returns the human readable string that describes the action
func MessageShareDeclined(userid, shareid, sharerid string) string {
	return fmt.Sprintf("user '%s' declined share '%s' from user '%s'", userid, shareid, sharerid)
}

// MessageLinkAccessed returns the human readable string that describes the action
func MessageLinkAccessed(linkid string, success bool) string {
	return fmt.Sprintf("link '%s' was accessed. Success: %v", linkid, success)
}

// MessageFileCreated returns the human readable string that describes the action
func MessageFileCreated(item string) string {
	return fmt.Sprintf("File '%s' was created", item)
}

// MessageFileRead returns the human readable string that describes the action
func MessageFileRead(item string) string {
	return fmt.Sprintf("File '%s' was read", item)
}

// MessageFileTrashed returns the human readable string that describes the action
func MessageFileTrashed(item string) string {
	return fmt.Sprintf("File '%s' was trashed", item)
}

// MessageFileRenamed returns the human readable string that describes the action
func MessageFileRenamed(item, oldpath, newpath string) string {
	return fmt.Sprintf("File '%s' was moved from '%s' to '%s'", item, oldpath, newpath)
}

// MessageFilePurged returns the human readable string that describes the action
func MessageFilePurged(item string) string {
	return fmt.Sprintf("File '%s' was removed from trashbin", item)
}

// MessageFileRestored returns the human readable string that describes the action
func MessageFileRestored(item, path string) string {
	return fmt.Sprintf("File '%s' was restored from trashbin to '%s'", item, path)
}

// MessageFileVersionRestored returns the human readable string that describes the action
func MessageFileVersionRestored(item string, version string) string {
	return fmt.Sprintf("File '%s' was restored in version '%s'", item, version)
}

// MessageSpaceCreated returns the human readable string that describes the action
func MessageSpaceCreated(spaceID string, name string) string {
	return fmt.Sprintf("Space '%s' with name '%s' was created", spaceID, name)
}

// MessageSpaceRenamed returns the human readable string that describes the action
func MessageSpaceRenamed(spaceID string, name string) string {
	return fmt.Sprintf("Space '%s' was renamed to '%s'", spaceID, name)
}

// MessageSpaceDisabled returns the human readable string that describes the action
func MessageSpaceDisabled(spaceID string) string {
	return fmt.Sprintf("Space '%s' was disabled", spaceID)
}

// MessageSpaceEnabled returns the human readable string that describes the action
func MessageSpaceEnabled(spaceID string) string {
	return fmt.Sprintf("Space '%s' was (re-) enabled", spaceID)
}

// MessageSpaceDeleted returns the human readable string that describes the action
func MessageSpaceDeleted(spaceID string) string {
	return fmt.Sprintf("Space '%s' was deleted", spaceID)
}
