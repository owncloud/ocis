package email

// Template marks the string as a translatable template
func Template(s string) string { return s }

// the available templates
var (
	// Shares
	ShareCreated = MessageTemplate{
		bodyTemplate: "shares/shareCreated.email.body.tmpl",
		// ShareCreated email template, Subject field
		Subject: Template(`{ShareSharer} shared '{ShareFolder}' with you`),
		// ShareCreated email template, Greeting field
		Greeting: Template(`Hello {ShareGrantee}`),
		// ShareCreated email template, Message field
		MessageBody: Template(`{ShareSharer} has shared "{ShareFolder}" with you.`),
		// ShareCreated email template, CallToAction field
		CallToAction: Template(`Click here to view it: {ShareLink}`),
	}

	ShareExpired = MessageTemplate{
		bodyTemplate: "shares/shareExpired.email.body.tmpl",
		// ShareExpired email template, Subject field
		Subject: Template(`Share to '{ShareFolder}' expired at {ExpiredAt}`),
		// ShareExpired email template, Greeting field
		Greeting: Template(`Hello {ShareGrantee},`),
		// ShareExpired email template, Message field
		MessageBody: Template(`Your share to {ShareFolder} has expired at {ExpiredAt}

Even though this share has been revoked you still might have access through other shares and/or space memberships.`),
	}

	// Spaces templates
	SharedSpace = MessageTemplate{
		bodyTemplate: "spaces/sharedSpace.email.body.tmpl",
		// SharedSpace email template, Subject field
		Subject: Template("{SpaceSharer} invited you to join {SpaceName}"),
		// SharedSpace email template, Greeting field
		Greeting: Template(`Hello {SpaceGrantee},`),
		// SharedSpace email template, Message field
		MessageBody: Template(`{SpaceSharer} has invited you to join "{SpaceName}".`),
		// SharedSpace email template, CallToAction field
		CallToAction: Template(`Click here to view it: {ShareLink}`),
	}

	UnsharedSpace = MessageTemplate{
		bodyTemplate: "spaces/unsharedSpace.email.body.tmpl",
		// UnsharedSpace email template, Subject field
		Subject: Template(`{SpaceSharer} removed you from {SpaceName}`),
		// UnsharedSpace email template, Greeting field
		Greeting: Template(`Hello {SpaceGrantee},`),
		// UnsharedSpace email template, Message field
		MessageBody: Template(`{SpaceSharer} has removed you from "{SpaceName}".

You might still have access through your other groups or direct membership.`),
		// UnsharedSpace email template, CallToAction field
		CallToAction: Template(`Click here to check it: {ShareLink}`),
	}

	MembershipExpired = MessageTemplate{
		bodyTemplate: "spaces/membershipExpired.email.body.tmpl",
		// MembershipExpired email template, Subject field
		Subject: Template(`Membership of '{SpaceName}' expired at {ExpiredAt}`),
		// MembershipExpired email template, Greeting field
		Greeting: Template(`Hello {SpaceGrantee},`),
		// MembershipExpired email template, Message field
		MessageBody: Template(`Your membership of space {SpaceName} has expired at {ExpiredAt}

Even though this membership has expired you still might have access through other shares and/or space memberships`),
	}
)

// holds the information to turn the raw template into a parseable go template
var _placeholders = map[string]string{
	"{ShareSharer}":  "{{ .ShareSharer }}",
	"{ShareFolder}":  "{{ .ShareFolder }}",
	"{ShareGrantee}": "{{ .ShareGrantee }}",
	"{ShareLink}":    "{{ .ShareLink }}",
	"{SpaceName}":    "{{ .SpaceName }}",
	"{SpaceGrantee}": "{{ .SpaceGrantee }}",
	"{SpaceSharer}":  "{{ .SpaceSharer }}",
	"{ExpiredAt}":    "{{ .ExpiredAt }}",
}

// MessageTemplate is the data structure for the email
type MessageTemplate struct {
	// bodyTemplate represent the path to .tmpl file
	bodyTemplate string
	Subject      string
	Greeting     string
	MessageBody  string
	CallToAction string
}
