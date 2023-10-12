package email

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/owncloud/ocis/v2/services/notifications/pkg/email/l10n"
)

// NewTextTemplate replace the body message template placeholders with the translated template
func NewTextTemplate(mt MessageTemplate, locale, defaultLocale string, translationPath string, vars map[string]string) (MessageTemplate, error) {
	var err error
	t := l10n.NewTranslator(locale, defaultLocale, translationPath)
	mt.Subject, err = composeMessage(t.Translate(mt.Subject), vars)
	if err != nil {
		return mt, err
	}
	mt.Greeting, err = composeMessage(t.Translate(mt.Greeting), vars)
	if err != nil {
		return mt, err
	}
	mt.MessageBody, err = composeMessage(t.Translate(mt.MessageBody), vars)
	if err != nil {
		return mt, err
	}
	mt.CallToAction, err = composeMessage(t.Translate(mt.CallToAction), vars)
	if err != nil {
		return mt, err
	}
	return mt, nil
}

// NewHTMLTemplate replace the body message template placeholders with the translated template
func NewHTMLTemplate(mt MessageTemplate, locale, defaultLocale string, translationPath string, vars map[string]string) (MessageTemplate, error) {
	var err error
	t := l10n.NewTranslator(locale, defaultLocale, translationPath)
	mt.Subject, err = composeMessage(t.Translate(mt.Subject), vars)
	if err != nil {
		return mt, err
	}
	mt.Greeting, err = composeMessage(newlineToBr(t.Translate(mt.Greeting)), vars)
	if err != nil {
		return mt, err
	}
	mt.MessageBody, err = composeMessage(newlineToBr(t.Translate(mt.MessageBody)), vars)
	if err != nil {
		return mt, err
	}
	mt.CallToAction, err = composeMessage(callToActionToHTML(t.Translate(mt.CallToAction)), vars)
	if err != nil {
		return mt, err
	}
	return mt, nil
}

// composeMessage renders the message based on template
func composeMessage(tmpl string, vars map[string]string) (string, error) {
	tpl, err := template.New("").Parse(replacePlaceholders(tmpl))
	if err != nil {
		return "", err
	}
	var writer bytes.Buffer
	if err := tpl.Execute(&writer, vars); err != nil {
		return "", err
	}
	return writer.String(), nil
}

func replacePlaceholders(raw string) string {
	for o, n := range _placeholders {
		raw = strings.ReplaceAll(raw, o, n)
	}
	return raw
}

func newlineToBr(s string) string {
	return strings.Replace(s, "\n", "<br>", -1)
}

func callToActionToHTML(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}
	s = strings.TrimSuffix(s, "{{ .ShareLink }}")
	return `<a href="{{ .ShareLink }}">` + s + `</a>`
}
