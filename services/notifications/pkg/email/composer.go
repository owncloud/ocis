package email

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/leonelquinteros/gotext"
)

var (
	//go:embed l10n/locale
	_translationFS embed.FS
	_domain        = "notifications"
)

// NewTextTemplate replace the body message template placeholders with the translated template
func NewTextTemplate(mt MessageTemplate, locale string, translationPath string, vars map[string]interface{}) (MessageTemplate, error) {
	var err error
	mt.Subject, err = ComposeMessage(mt.Subject, locale, translationPath, vars)
	if err != nil {
		return mt, err
	}
	mt.Greeting, err = ComposeMessage(mt.Greeting, locale, translationPath, vars)
	if err != nil {
		return mt, err
	}
	mt.MessageBody, err = ComposeMessage(mt.MessageBody, locale, translationPath, vars)
	if err != nil {
		return mt, err
	}
	mt.CallToAction, err = ComposeMessage(mt.CallToAction, locale, translationPath, vars)
	if err != nil {
		return mt, err
	}
	return mt, nil
}

// NewHTMLTemplate replace the body message template placeholders with the translated template
func NewHTMLTemplate(mt MessageTemplate, locale string, translationPath string, vars map[string]interface{}) (MessageTemplate, error) {
	var err error
	mt.Subject, err = ComposeMessage(mt.Subject, locale, translationPath, vars)
	if err != nil {
		return mt, err
	}
	mt.Greeting, err = ComposeMessage(newlineToBr(mt.Greeting), locale, translationPath, vars)
	if err != nil {
		return mt, err
	}
	mt.MessageBody, err = ComposeMessage(newlineToBr(mt.MessageBody), locale, translationPath, vars)
	if err != nil {
		return mt, err
	}
	mt.CallToAction, err = ComposeMessage(callToActionToHTML(mt.CallToAction), locale, translationPath, vars)
	if err != nil {
		return mt, err
	}
	return mt, nil
}

// ComposeMessage renders the message based on template
func ComposeMessage(template, locale string, path string, vars map[string]interface{}) (string, error) {
	raw := loadTemplate(template, locale, path)
	return executeRaw(replacePlaceholders(raw), vars)
}

func loadTemplate(template, locale string, path string) string {
	// Create Locale with library path and language code and load default domain
	var l *gotext.Locale
	if path == "" {
		filesystem, _ := fs.Sub(_translationFS, "l10n/locale")
		l = gotext.NewLocaleFS(locale, filesystem)
	} else { // use custom path instead
		l = gotext.NewLocale(path, locale)
	}
	l.AddDomain(_domain) // make domain configurable only if needed
	return l.Get(template)
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
