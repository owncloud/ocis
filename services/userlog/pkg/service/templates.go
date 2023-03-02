package service

import "text/template"

// the available templates
var (
	SpaceDisabled        = "space-disabled"
	SpaceDisabledSubject = "space disabled"
	SpaceDisabledMessage = "{{ .username }} disabled space {{ .spacename }}"
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
			Subject: template.Must(template.New(SpaceDisabled).Parse(SpaceDisabledSubject)),
			Message: template.Must(template.New("").Parse(SpaceDisabledMessage)),
		},
	}
)
