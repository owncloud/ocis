package email

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
)

var (
	//go:embed l10n/locale
	_translationFS embed.FS
	_domain        = "notifications"
)

// NewTextTemplate replace the body message template placeholders with the translated template
func NewTextTemplate(mt MessageTemplate, locale, defaultLocale string, translationPath string, vars map[string]string) (MessageTemplate, error) {
	var err error
	t := l10n.NewTranslatorFromCommonConfig(defaultLocale, _domain, translationPath, _translationFS, "l10n/locale").Locale(locale)
	mt.Subject, err = composeMessage(t.Get(mt.Subject), vars)
	if err != nil {
		return mt, err
	}
	mt.Greeting, err = composeMessage(t.Get(mt.Greeting), vars)
	if err != nil {
		return mt, err
	}
	mt.MessageBody, err = composeMessage(t.Get(mt.MessageBody), vars)
	if err != nil {
		return mt, err
	}
	mt.CallToAction, err = composeMessage(t.Get(mt.CallToAction), vars)
	if err != nil {
		return mt, err
	}
	return mt, nil
}

// NewHTMLTemplate replace the body message template placeholders with the translated template
func NewHTMLTemplate(mt MessageTemplate, locale, defaultLocale string, translationPath string, vars map[string]string) (MessageTemplate, error) {
	var err error
	t := l10n.NewTranslatorFromCommonConfig(defaultLocale, _domain, translationPath, _translationFS, "l10n/locale").Locale(locale)
	mt.Subject, err = composeMessage(t.Get(mt.Subject), vars)
	if err != nil {
		return mt, err
	}
	mt.Greeting, err = composeMessage(newlineToBr(t.Get(mt.Greeting)), vars)
	if err != nil {
		return mt, err
	}
	mt.MessageBody, err = composeMessage(newlineToBr(t.Get(mt.MessageBody)), vars)
	if err != nil {
		return mt, err
	}
	mt.CallToAction, err = composeMessage(callToActionToHTML(t.Get(mt.CallToAction)), vars)
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
	s = strings.TrimSuffix(s, "{ShareLink}")
	return s + `<a href="{ShareLink}">{ShareLink}</a>`
}
