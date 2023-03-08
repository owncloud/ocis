package service

// Template marks the string as a translatable template
func Template(s string) string { return s }

// the available templates
var (
	SpaceShared = NotificationTemplate{
		Subject: Template("Space shared"),
		Message: Template("{user} added you to Space {space}"),
	}

	SpaceUnshared = NotificationTemplate{
		Subject: Template("Removed from Space"),
		Message: Template("{user} removed you from Space {space}"),
	}

	SpaceDisabled = NotificationTemplate{
		Subject: Template("Space disabled"),
		Message: Template("{user} disabled Space {space}"),
	}

	SpaceDeleted = NotificationTemplate{
		Subject: Template("Space deleted"),
		Message: Template("{user} deleted Space {space}"),
	}

	SpaceMembershipExpired = NotificationTemplate{
		Subject: Template("Membership expired"),
		Message: Template("Access to Space {space} lost"),
	}

	ShareCreated = NotificationTemplate{
		Subject: Template("Resource shared"),
		Message: Template("{user} shared {resource} with you"),
	}

	ShareRemoved = NotificationTemplate{
		Subject: Template("Resource unshared"),
		Message: Template("{user} unshared {resource} with you"),
	}

	ShareExpired = NotificationTemplate{
		Subject: Template("Share expired"),
		Message: Template("Access to {resource} expired"),
	}
)

// holds the information to turn the raw template into a parseable go template
var _placeholders = map[string]string{
	"{user}":     "{{ .username }}",
	"{space}":    "{{ .spacename }}",
	"{resource}": "{{ .resourcename }}",
}

// NotificationTemplate is the data structure for the notifications
type NotificationTemplate struct {
	Subject string
	Message string
}
