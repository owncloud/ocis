// Package l10n implements utility for translation the text templates.
//
// The l10n package use transifex translation for text templates.
package l10n

import (
	"embed"
	"io/fs"

	"github.com/leonelquinteros/gotext"
)

var (
	//go:embed locale
	_translationFS embed.FS
	_domain        = "notifications"
)

// Translator is the interface to the translation
type Translator interface {
	Translate(str string) string
}

type translator struct {
	l *gotext.Locale
}

// NewTranslator Create Translator with library path and language code and load default domain
func NewTranslator(local string, path string) Translator {
	var l *gotext.Locale
	if path == "" {
		filesystem, _ := fs.Sub(_translationFS, "locale")
		l = gotext.NewLocaleFS(local, filesystem)
	} else { // use custom path instead
		l = gotext.NewLocale(path, local)
	}
	l.AddDomain(_domain) // make domain configurable only if needed
	return &translator{l: l}
}

func (t *translator) Translate(str string) string {
	return t.l.Get(str)
}
