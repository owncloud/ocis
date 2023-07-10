package service

// Template marks the string as a translatable template
func Template(s string) string { return s }

// the available templates
var (
	VirusFound = NotificationTemplate{
		Subject: Template("Virus found"),
		Message: Template("Virus found in {resource}. Upload not possible. Virus: {virus}"),
	}

	PoliciesEnforced = NotificationTemplate{
		Subject: Template("Policies enforced"),
		Message: Template("File {resource} was deleted because it violates the policies"),
	}

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

	PlatformDeprovision = NotificationTemplate{
		Subject: Template("Instance will be shut down and deprovisioned"),
		Message: Template("Attention! The instance will be shut down and deprovisioned on {date}. Download all your data before that date as no access past that date is possible."),
	}
)

// holds the information to turn the raw template into a parseable go template
var _placeholders = map[string]string{
	"{user}":     "{{ .username }}",
	"{space}":    "{{ .spacename }}",
	"{resource}": "{{ .resourcename }}",
	"{virus}":    "{{ .virusdescription }}",
	"{date}":     "{{ .date }}",
}

// NotificationTemplate is the data structure for the notifications
type NotificationTemplate struct {
	Subject string
	Message string
}
