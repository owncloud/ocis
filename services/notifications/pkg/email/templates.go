package email

import "github.com/owncloud/ocis/v2/ocis-pkg/l10n"

// the available templates
var (
	// Shares
	ShareCreated = MessageTemplate{
		textTemplate: "templates/text/email.text.tmpl",
		htmlTemplate: "templates/html/email.html.tmpl",
		// ShareCreated email template, Subject field (resolves directly)
		Subject: l10n.Template(`{ShareSharer} shared '{ShareFolder}' with you`),
		// ShareCreated email template, resolves via {{ .Greeting }}
		Greeting: l10n.Template(`Hello {ShareGrantee}`),
		// ShareCreated email template, resolves via {{ .MessageBody }}
		MessageBody: l10n.Template(`{ShareSharer} has shared "{ShareFolder}" with you.`),
		// ShareCreated email template, resolves via {{ .CallToAction }}
		CallToAction: l10n.Template(`Click here to view it: {ShareLink}`),
	}

	ShareExpired = MessageTemplate{
		textTemplate: "templates/text/email.text.tmpl",
		htmlTemplate: "templates/html/email.html.tmpl",
		// ShareExpired email template, Subject field (resolves directly)
		Subject: l10n.Template(`Share to '{ShareFolder}' expired at {ExpiredAt}`),
		// ShareExpired email template, resolves via {{ .Greeting }}
		Greeting: l10n.Template(`Hello {ShareGrantee},`),
		// ShareExpired email template, resolves via {{ .MessageBody }}
		MessageBody: l10n.Template(`Your share to {ShareFolder} has expired at {ExpiredAt}

Even though this share has been revoked you still might have access through other shares and/or space memberships.`),
	}

	// Spaces templates
	SharedSpace = MessageTemplate{
		textTemplate: "templates/text/email.text.tmpl",
		htmlTemplate: "templates/html/email.html.tmpl",
		// SharedSpace email template, Subject field (resolves directly)
		Subject: l10n.Template("{SpaceSharer} invited you to join {SpaceName}"),
		// SharedSpace email template, resolves via {{ .Greeting }}
		Greeting: l10n.Template(`Hello {SpaceGrantee},`),
		// SharedSpace email template, resolves via {{ .MessageBody }}
		MessageBody: l10n.Template(`{SpaceSharer} has invited you to join "{SpaceName}".`),
		// SharedSpace email template, resolves via {{ .CallToAction }}
		CallToAction: l10n.Template(`Click here to view it: {ShareLink}`),
	}

	UnsharedSpace = MessageTemplate{
		textTemplate: "templates/text/email.text.tmpl",
		htmlTemplate: "templates/html/email.html.tmpl",
		// UnsharedSpace email template, Subject field (resolves directly)
		Subject: l10n.Template(`{SpaceSharer} removed you from {SpaceName}`),
		// UnsharedSpace email template, resolves via {{ .Greeting }}
		Greeting: l10n.Template(`Hello {SpaceGrantee},`),
		// UnsharedSpace email template, resolves via {{ .MessageBody }}
		MessageBody: l10n.Template(`{SpaceSharer} has removed you from "{SpaceName}".

You might still have access through your other groups or direct membership.`),
		// UnsharedSpace email template, resolves via {{ .CallToAction }}
		CallToAction: l10n.Template(`Click here to check it: {ShareLink}`),
	}

	MembershipExpired = MessageTemplate{
		textTemplate: "templates/text/email.text.tmpl",
		htmlTemplate: "templates/html/email.html.tmpl",
		// MembershipExpired email template, Subject field (resolves directly)
		Subject: l10n.Template(`Membership of '{SpaceName}' expired at {ExpiredAt}`),
		// MembershipExpired email template, resolves via {{ .Greeting }}
		Greeting: l10n.Template(`Hello {SpaceGrantee},`),
		// MembershipExpired email template, resolves via {{ .MessageBody }}
		MessageBody: l10n.Template(`Your membership of space {SpaceName} has expired at {ExpiredAt}

Even though this membership has expired you still might have access through other shares and/or space memberships`),
	}

	ScienceMeshInviteTokenGenerated = MessageTemplate{
		textTemplate: "templates/text/email.text.tmpl",
		htmlTemplate: "templates/html/email.html.tmpl",
		// ScienceMeshInviteTokenGenerated email template, Subject field (resolves directly)
		Subject: l10n.Template(`ScienceMesh: {InitiatorName} wants to collaborate with you`),
		// ScienceMeshInviteTokenGenerated email template, resolves via {{ .Greeting }}
		Greeting: l10n.Template(`Hi,`),
		// ScienceMeshInviteTokenGenerated email template, resolves via {{ .MessageBody }}
		MessageBody: l10n.Template(`{ShareSharer} ({ShareSharerMail}) wants to start sharing OCM resources with you.
{{if .ShareLink }}To accept the invite, please visit the following URL:
{ShareLink}

Alternatively, you can{{else}}
Please{{end}} visit your mesh provider and use the following details:
  Token: {Token}
  ProviderDomain: {ProviderDomain}`),
	}
)

// holds the information to turn the raw template into a parseable go template
var _placeholders = map[string]string{
	"{ShareSharer}":     "{{ .ShareSharer }}",
	"{ShareFolder}":     "{{ .ShareFolder }}",
	"{ShareGrantee}":    "{{ .ShareGrantee }}",
	"{ShareLink}":       "{{ .ShareLink }}",
	"{SpaceName}":       "{{ .SpaceName }}",
	"{SpaceGrantee}":    "{{ .SpaceGrantee }}",
	"{SpaceSharer}":     "{{ .SpaceSharer }}",
	"{ExpiredAt}":       "{{ .ExpiredAt }}",
	"{ShareSharerMail}": "{{ .ShareSharerMail }}",
	"{ProviderDomain}":  "{{ .ProviderDomain }}",
	"{Token}":           "{{ .Token }}",
}

// MessageTemplate is the data structure for the email
type MessageTemplate struct {
	// textTemplate represent the path to text plain .tmpl file
	textTemplate string
	// htmlTemplate represent the path to html .tmpl file
	htmlTemplate string
	// The fields below represent the placeholders for the translatable templates
	Subject      string
	Greeting     string
	MessageBody  string
	CallToAction string
}
