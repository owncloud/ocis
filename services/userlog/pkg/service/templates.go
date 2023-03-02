package service

import "text/template"

// the available templates
var (
	SpaceDisabled        = "space-disabled"
	SpaceDisabledSubject = "Space disabled"
	SpaceDisabledMessage = "{{ .username }} disabled Space {{ .spacename }}"

	SpaceShared        = "space-shared"
	SpaceSharedSubject = "Space shared"
	SpaceSharedMessage = "{{ .username }} shared Space {{ .spacename }} with you"
)

// NotificationTemplate is the data structure for the notifications
type NotificationTemplate struct {
	Subject *template.Template
	Message *template.Template
}

// rendered templates
var (
	_templates = map[string]NotificationTemplate{
		SpaceDisabled: {
			Subject: template.Must(template.New("").Parse(SpaceDisabledSubject)),
			Message: template.Must(template.New("").Parse(SpaceDisabledMessage)),
		},
		SpaceShared: {
			Subject: template.Must(template.New("").Parse(SpaceSharedSubject)),
			Message: template.Must(template.New("").Parse(SpaceSharedMessage)),
		},
	}
)
