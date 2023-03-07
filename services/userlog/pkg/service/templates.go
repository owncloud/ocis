package service

// the available templates
var (
	SpaceShared = NotificationTemplate{
		Subject: "Space shared",
		Message: "{{ .username }} added you to Space {{ .spacename }}",
	}

	SpaceUnshared = NotificationTemplate{
		Subject: "Removed from Space",
		Message: "{{ .username }} removed you from Space {{ .spacename }}",
	}

	SpaceDisabled = NotificationTemplate{
		Subject: "Space disabled",
		Message: "{{ .username }} disabled Space {{ .spacename }}",
	}

	SpaceDeleted = NotificationTemplate{
		Subject: "Space deleted",
		Message: "{{ .username }} deleted Space {{ .spacename }}",
	}

	SpaceMembershipExpired = NotificationTemplate{
		Subject: "Membership expired",
		Message: "Access to Space {{ .spacename }} lost",
	}

	ShareCreated = NotificationTemplate{
		Subject: "Resource shared",
		Message: "{{ .username }} shared {{ .resourcename }} with you",
	}

	ShareRemoved = NotificationTemplate{
		Subject: "Resource unshared",
		Message: "{{ .username }} unshared {{ .resourcename }} with you",
	}

	ShareExpired = NotificationTemplate{
		Subject: "Share expired",
		Message: "Access to {{ .resourcename }} expired",
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
