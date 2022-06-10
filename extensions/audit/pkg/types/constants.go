package types

import (
	"fmt"
	"strings"

	"github.com/cs3org/reva/v2/pkg/events"
)

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
	ActionContainerCreated    = "container_create"
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

	// Users
	ActionUserCreated        = "user_created"
	ActionUserDeleted        = "user_deleted"
	ActionUserFeatureChanged = "user_feature_changed"

	// Groups
	ActionGroupCreated       = "group_created"
	ActionGroupDeleted       = "group_deleted"
	ActionGroupMemberAdded   = "group_member_added"
	ActionGroupMemberRemoved = "group_member_removed"
)

// MessageShareCreated returns the human readable string that describes the action
func MessageShareCreated(sharer, item, grantee string) string {
	return fmt.Sprintf("user '%s' shared file '%s' with '%s'", sharer, item, grantee)
}

// MessageLinkCreated returns the human readable string that describes the action
func MessageLinkCreated(sharer, item, shareid string) string {
	return fmt.Sprintf("user '%s' created a public link to file '%s' with id '%s'", sharer, item, shareid)
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

// MessageContainerCreated returns the human readable string that describes the action
func MessageContainerCreated(item string) string {
	return fmt.Sprintf("Folder '%s' was created", item)
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

// MessageUserCreated returns the human readable string that describes the action
func MessageUserCreated(userID string) string {
	return fmt.Sprintf("User '%s' was created", userID)
}

// MessageUserDeleted returns the human readable string that describes the action
func MessageUserDeleted(userID string) string {
	return fmt.Sprintf("User '%s' was deleted", userID)
}

// MessageUserFeatureChanged returns the human readable string that describes the action
func MessageUserFeatureChanged(userID string, features []events.UserFeature) string {
	// Result is: "User %username%'s feature changed: %featurename%=%featurevalue% %featurename%=%featurevalue%"
	var sb strings.Builder
	sb.WriteString("User ")
	sb.WriteString(userID)
	sb.WriteString("'s feature changed: ")
	for _, f := range features {
		sb.WriteString(f.Name)
		sb.WriteRune('=')
		sb.WriteString(f.Value)
		sb.WriteRune(' ')
	}
	return sb.String()
}

// MessageGroupCreated returns the human readable string that describes the action
func MessageGroupCreated(groupID string) string {
	return fmt.Sprintf("Group '%s' was created", groupID)
}

// MessageGroupDeleted returns the human readable string that describes the action
func MessageGroupDeleted(groupID string) string {
	return fmt.Sprintf("Group '%s' was deleted", groupID)
}

// MessageGroupMemberAdded returns the human readable string that describes the action
func MessageGroupMemberAdded(userID, groupID string) string {
	return fmt.Sprintf("User '%s' was added to group '%s'", userID, groupID)
}

// MessageGroupMemberRemoved returns the human readable string that describes the action
func MessageGroupMemberRemoved(userID, groupID string) string {
	return fmt.Sprintf("User '%s' was removed from group '%s'", userID, groupID)
}
