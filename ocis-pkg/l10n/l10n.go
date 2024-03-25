// package l10n holds translation mechanics that are used by user facing services (notifications, userlog, graph)
package l10n

import (
	"io/fs"
	"os"

	"github.com/leonelquinteros/gotext"
)

// Translator is able to translate strings
type Translator struct {
	fs            fs.FS
	defaultLocale string
	domain        string
}

// NewTranslator creates a Translator with library path and language code and load default domain
func NewTranslator(defaultLocale string, domain string, fsys fs.FS) Translator {
	return Translator{
		fs:            fsys,
		defaultLocale: defaultLocale,
		domain:        domain,
	}
}

// NewTranslatorFromCommonConfig creates a new Translator from legacy config
func NewTranslatorFromCommonConfig(defaultLocale string, domain string, path string, fsys fs.FS, fsSubPath string) Translator {
	var filesystem fs.FS
	if path == "" {
		filesystem, _ = fs.Sub(fsys, fsSubPath)
	} else { // use custom path instead
		filesystem = os.DirFS(path)
	}
	return NewTranslator(defaultLocale, domain, filesystem)
}

// Translate translates a string to the locale
func (t Translator) Translate(str, locale string) string {
	return t.Locale(locale).Get(str)
}

// Locale returns the gotext.Locale, use `.Get` method to translate strings
func (t Translator) Locale(locale string) *gotext.Locale {
	l := gotext.NewLocaleFS(locale, t.fs)
	l.AddDomain(t.domain) // make domain configurable only if needed
	if locale != "en" && len(l.GetTranslations()) == 0 {
		l = gotext.NewLocaleFS(t.defaultLocale, t.fs)
		l.AddDomain(t.domain) // make domain configurable only if needed
	}
	return l
}
