package email

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/pkg/errors"

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
	// the translated fields below are translation lookup keys, not printf format
	// strings; go vet's non-constant format string check can't trace through a
	// func value, so assign Get to a local before calling it to avoid a false positive.
	get := t.Get
	mt.Subject, err = composeMessage(get(mt.Subject), vars)
	if err != nil {
		return mt, err
	}
	mt.Greeting, err = composeMessage(get(mt.Greeting), vars)
	if err != nil {
		return mt, err
	}
	mt.MessageBody, err = composeMessage(get(mt.MessageBody), vars)
	if err != nil {
		return mt, err
	}

	if mt.CallToAction != "" {
		// Some templates have an empty call-to-action. We don't want to translate
		// an empty key and get an unexpected message.
		mt.CallToAction, err = composeMessage(get(mt.CallToAction), vars)
		if err != nil {
			return mt, err
		}
	}
	return mt, nil
}

// NewHTMLTemplate replace the body message template placeholders with the translated template
func NewHTMLTemplate(mt MessageTemplate, locale, defaultLocale string, translationPath string, vars map[string]string) (MessageTemplate, error) {
	var err error
	t := l10n.NewTranslatorFromCommonConfig(defaultLocale, _domain, translationPath, _translationFS, "l10n/locale").Locale(locale)
	// the translated fields below are translation lookup keys, not printf format
	// strings; go vet's non-constant format string check can't trace through a
	// func value, so assign Get to a local before calling it to avoid a false positive.
	get := t.Get
	mt.Subject, err = composeMessage(get(mt.Subject), vars)
	if err != nil {
		return mt, err
	}
	mt.Greeting, err = composeMessage(newlineToBr(get(mt.Greeting)), vars)
	if err != nil {
		return mt, err
	}
	mt.MessageBody, err = composeMessage(newlineToBr(get(mt.MessageBody)), vars)
	if err != nil {
		return mt, err
	}
	if mt.CallToAction != "" {
		// Some templates have an empty call-to-action. We don't want to translate
		// an empty key and get an unexpected message.
		mt.CallToAction, err = composeMessage(callToActionToHTML(get(mt.CallToAction)), vars)
		if err != nil {
			return mt, err
		}
	}
	return mt, nil
}

// NewGroupedTextTemplate replace the body message template placeholders with the translated template
func NewGroupedTextTemplate(gmt GroupedMessageTemplate, vars map[string]string, locale, defaultLocale string, translationPath string, mts []MessageTemplate, mtsVars []map[string]string) (GroupedMessageTemplate, error) {
	if len(mts) != len(mtsVars) {
		return gmt, errors.New("number of templates does not match number of variables")
	}

	var err error
	t := l10n.NewTranslatorFromCommonConfig(defaultLocale, _domain, translationPath, _translationFS, "l10n/locale").Locale(locale)
	// the translated fields below are translation lookup keys, not printf format
	// strings; go vet's non-constant format string check can't trace through a
	// func value, so assign Get to a local before calling it to avoid a false positive.
	get := t.Get
	gmt.Subject, err = composeMessage(get(gmt.Subject), vars)
	if err != nil {
		return gmt, err
	}
	gmt.Greeting, err = composeMessage(get(gmt.Greeting), vars)
	if err != nil {
		return gmt, err
	}

	bodyParts := make([]string, 0, len(mtsVars))
	for i, mt := range mts {
		bodyPart, err := composeMessage(get(mt.MessageBody), mtsVars[i])
		if err != nil {
			return gmt, err
		}
		bodyParts = append(bodyParts, bodyPart)
	}
	gmt.MessageBody = strings.Join(bodyParts, "\n\n\n")
	return gmt, nil
}

// NewGroupedHTMLTemplate replace the body message template placeholders with the translated template
func NewGroupedHTMLTemplate(gmt GroupedMessageTemplate, vars map[string]string, locale, defaultLocale string, translationPath string, mts []MessageTemplate, mtsVars []map[string]string) (GroupedMessageTemplate, error) {
	if len(mts) != len(mtsVars) {
		return gmt, errors.New("number of templates does not match number of variables")
	}

	var err error
	t := l10n.NewTranslatorFromCommonConfig(defaultLocale, _domain, translationPath, _translationFS, "l10n/locale").Locale(locale)
	// the translated fields below are translation lookup keys, not printf format
	// strings; go vet's non-constant format string check can't trace through a
	// func value, so assign Get to a local before calling it to avoid a false positive.
	get := t.Get
	gmt.Subject, err = composeMessage(get(gmt.Subject), vars)
	if err != nil {
		return gmt, err
	}
	gmt.Greeting, err = composeMessage(newlineToBr(get(gmt.Greeting)), vars)
	if err != nil {
		return gmt, err
	}

	bodyParts := make([]string, 0, len(mtsVars))
	for i, mt := range mts {
		bodyPart, err := composeMessage(get(mt.MessageBody), mtsVars[i])
		if err != nil {
			return gmt, err
		}
		bodyParts = append(bodyParts, bodyPart)
	}
	gmt.MessageBody = newlineToBr(strings.Join(bodyParts, "<br><br><br>"))

	return gmt, nil
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
