package service

import "text/template"

// the available templates
var (
	SpaceShared        = "space-shared"
	SpaceSharedSubject = "Space shared"
	SpaceSharedMessage = "{{ .username }} shared Space {{ .spacename }} with you"

	SpaceUnshared        = "space-unshared"
	SpaceUnsharedSubject = "Removed from Space"
	SpaceUnsharedMessage = "{{ .username }} removed you from Space {{ .spacename }}"

	SpaceDisabled        = "space-disabled"
	SpaceDisabledSubject = "Space disabled"
	SpaceDisabledMessage = "{{ .username }} disabled Space {{ .spacename }}"

	SpaceDeleted        = "space-deleted"
	SpaceDeletedSubject = "Space deleted"
	SpaceDeletedMessage = "{{ .username }} deleted Space {{ .spacename }}"

	SpaceMembershipExpired        = "space-membership-expired"
	SpaceMembershipExpiredSubject = "Membership expired"
	SpaceMembershipExpiredMessage = "Access to Space {{ .spacename }} lost"

	ShareCreated        = "item-shared"
	ShareCreatedSubject = "Resource shared"
	ShareCreatedMessage = "{{ .username }} shared {{ .itemname }} with you"

	ShareRemoved        = "item-unshared"
	ShareRemovedSubject = "Resource unshared"
	ShareRemovedMessage = "{{ .username }} unshared {{ .itemname }} with you"

	ShareExpired        = "share-expired"
	ShareExpiredSubject = "Share expired"
	ShareExpiredMessage = "Access to {{ .resourcename }} expired"
)

// rendered templates
var (
	_templates = map[string]NotificationTemplate{
		SpaceShared:            notiTmpl(SpaceSharedSubject, SpaceSharedMessage),
		SpaceUnshared:          notiTmpl(SpaceUnsharedSubject, SpaceUnsharedMessage),
		SpaceDisabled:          notiTmpl(SpaceDisabledSubject, SpaceDisabledMessage),
		SpaceDeleted:           notiTmpl(SpaceDeletedSubject, SpaceDeletedMessage),
		SpaceMembershipExpired: notiTmpl(SpaceMembershipExpiredSubject, SpaceMembershipExpiredMessage),
		ShareCreated:           notiTmpl(ShareCreatedSubject, ShareCreatedMessage),
		ShareRemoved:           notiTmpl(ShareRemovedSubject, ShareRemovedMessage),
		ShareExpired:           notiTmpl(ShareExpiredSubject, ShareExpiredMessage),
	}
)

// NotificationTemplate is the data structure for the notifications
type NotificationTemplate struct {
	Subject *template.Template
	Message *template.Template
}

func notiTmpl(subjectname string, messagename string) NotificationTemplate {
	return NotificationTemplate{
		Subject: template.Must(template.New("").Parse(subjectname)),
		Message: template.Must(template.New("").Parse(messagename)),
	}
}
