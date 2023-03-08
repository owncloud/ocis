package service

// Template marks the string as a translatable template
func Template(s string) string { return s }

// the available templates
var (
	SpaceShared = NotificationTemplate{
		Subject: Template("Space shared"),
		Message: Template("{{ .username }} added you to Space {{ .spacename }}"),
	}

	SpaceUnshared = NotificationTemplate{
		Subject: Template("Removed from Space"),
		Message: Template("{{ .username }} removed you from Space {{ .spacename }}"),
	}

	SpaceDisabled = NotificationTemplate{
		Subject: Template("Space disabled"),
		Message: Template("{{ .username }} disabled Space {{ .spacename }}"),
	}

	SpaceDeleted = NotificationTemplate{
		Subject: Template("Space deleted"),
		Message: Template("{{ .username }} deleted Space {{ .spacename }}"),
	}

	SpaceMembershipExpired = NotificationTemplate{
		Subject: Template("Membership expired"),
		Message: Template("Access to Space {{ .spacename }} lost"),
	}

	ShareCreated = NotificationTemplate{
		Subject: Template("Resource shared"),
		Message: Template("{{ .username }} shared {{ .resourcename }} with you"),
	}

	ShareRemoved = NotificationTemplate{
		Subject: Template("Resource unshared"),
		Message: Template("{{ .username }} unshared {{ .resourcename }} with you"),
	}

	ShareExpired = NotificationTemplate{
		Subject: Template("Share expired"),
		Message: Template("Access to {{ .resourcename }} expired"),
	}
)

// holds the information to link the raw template to the details
var _placeholders = map[string]interface{}{
	"username":     "{user}",
	"spacename":    "{space}",
	"resourcename": "{resource}",
}

// NotificationTemplate is the data structure for the notifications
type NotificationTemplate struct {
	Subject string
	Message string
}
