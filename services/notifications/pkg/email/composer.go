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

// ComposeMessage renders the message based on template
func ComposeMessage(template, locale string, path string) string {
	raw := loadTemplate(template, locale, path)
	return replacePlaceholders(raw)
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
