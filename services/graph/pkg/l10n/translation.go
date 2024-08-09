package l10n

import (
	"embed"

	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
)

var (
	//go:embed locale
	_localeFS embed.FS
)

const (
	// subfolder where the translation files are stored
	_localeSubPath = "locale"

	// domain of the graph service (transifex)
	_domain = "graph"
)

// Translate translates a string based on the locale and default locale
func Translate(content, locale, defaultLocale string) string {
	t := l10n.NewTranslatorFromCommonConfig(defaultLocale, _domain, "", _localeFS, _localeSubPath)
	return t.Translate(content, locale)
}

// TranslateEntity returns a function that translates a struct or slice based on the locale
func TranslateEntity(locale, defaultLocale string, entity any, opts ...l10n.TranslateOption) error {
	t := l10n.NewTranslatorFromCommonConfig(defaultLocale, _domain, "", _localeFS, _localeSubPath)
	return t.TranslateEntity(locale, entity, opts...)
}
