package types

import "fmt"

// short identifiers for audit actions
const (
	ActionShareCreated            = "share_created"
	ActionSharePermissionUpdated  = "share_permission_updated"
	ActionShareDisplayNameUpdated = "share_name_updated"
	ActionSharePasswordUpdated    = "share_password_updated"
	ActionShareExpirationUpdated  = "share_expiration_updated"
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
	return fmt.Sprintf("user '%s' modified field '%s' of public link '%s'", sharer, fieldUpdated, shareid)
}
