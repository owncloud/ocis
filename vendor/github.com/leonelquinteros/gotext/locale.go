/*
 * Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com. All rights reserved.
 * Licensed under the MIT License. See LICENSE file in the project root for full license information.
 */

package gotext

import (
	"bytes"
	"encoding/gob"
	"io/fs"
	"os"
	"path"
	"sync"
)

/*
Locale wraps the entire i18n collection for a single language (locale)
It's used by the package functions, but it can also be used independently to handle
multiple languages at the same time by working with this object.

Example:

	    import (
		"encoding/gob"
		"bytes"
		    "fmt"
		    "github.com/leonelquinteros/gotext"
	    )

	    func main() {
	        // Create Locale with library path and language code
	        l := gotext.NewLocale("/path/to/i18n/dir", "en_US")

	        // Load domain '/path/to/i18n/dir/en_US/LC_MESSAGES/default.{po,mo}'
	        l.AddDomain("default")

	        // Translate text from default domain
	        fmt.Println(l.Get("Translate this"))

	        // Load different domain ('/path/to/i18n/dir/en_US/LC_MESSAGES/extras.{po,mo}')
	        l.AddDomain("extras")

	        // Translate text from domain
	        fmt.Println(l.GetD("extras", "Translate this"))
	    }
*/
type Locale struct {
	// Path to locale files.
	path string

	// Language for this Locale
	lang string

	// List of available Domains for this locale.
	Domains map[string]Translator

	// First AddDomain is default Domain
	defaultDomain string

	// Sync Mutex
	sync.RWMutex

	// optional fs to use
	fs fs.FS
}

// NewLocale creates and initializes a new Locale object for a given language.
// It receives a path for the i18n .po/.mo files directory (p) and a language code to use (l).
func NewLocale(p, l string) *Locale {
	return &Locale{
		path:    p,
		lang:    SimplifiedLocale(l),
		Domains: make(map[string]Translator),
	}
}

// NewLocaleFS returns a Locale working with a fs.FS
func NewLocaleFS(l string, filesystem fs.FS) *Locale {
	loc := NewLocale("", l)
	loc.fs = filesystem
	return loc
}

// NewLocaleFSWithPath returns a Locale working with a fs.FS on a p path folder.
func NewLocaleFSWithPath(l string, filesystem fs.FS, p string) *Locale {
	loc := NewLocale("", l)
	loc.fs = filesystem
	loc.path = p
	return loc
}

func (l *Locale) findExt(dom, ext string) string {
	filename := path.Join(l.path, l.lang, "LC_MESSAGES", dom+"."+ext)
	if l.fileExists(filename) {
		return filename
	}

	if len(l.lang) > 2 {
		filename = path.Join(l.path, l.lang[:2], "LC_MESSAGES", dom+"."+ext)
		if l.fileExists(filename) {
			return filename
		}
	}

	filename = path.Join(l.path, l.lang, dom+"."+ext)
	if l.fileExists(filename) {
		return filename
	}

	if len(l.lang) > 2 {
		filename = path.Join(l.path, l.lang[:2], dom+"."+ext)
		if l.fileExists(filename) {
			return filename
		}
	}

	return ""
}

// GetActualLanguage inspects the filesystem and decides whether to strip
// a CC part of the ll_CC locale string.
func (l *Locale) GetActualLanguage(dom string) string {
	extensions := []string{"mo", "po"}
	var fp string
	for _, ext := range extensions {
		// 'll' (or 'll_CC') exists, and it was specified as-is
		fp = path.Join(l.path, l.lang, "LC_MESSAGES", dom+"."+ext)
		if l.fileExists(fp) {
			return l.lang
		}
		// 'll' exists, but 'll_CC' was specified
		if len(l.lang) > 2 {
			fp = path.Join(l.path, l.lang[:2], "LC_MESSAGES", dom+"."+ext)
			if l.fileExists(fp) {
				return l.lang[:2]
			}
		}
		// 'll' (or 'll_CC') exists outside of LC_category, and it was specified as-is
		fp = path.Join(l.path, l.lang, dom+"."+ext)
		if l.fileExists(fp) {
			return l.lang
		}
		// 'll' exists outside of LC_category, but 'll_CC' was specified
		if len(l.lang) > 2 {
			fp = path.Join(l.path, l.lang[:2], dom+"."+ext)
			if l.fileExists(fp) {
				return l.lang[:2]
			}
		}
	}
	return ""
}

func (l *Locale) fileExists(filename string) bool {
	if l.fs != nil {
		_, err := fs.Stat(l.fs, filename)
		return err == nil
	}
	_, err := os.Stat(filename)
	return err == nil
}

// AddDomain creates a new domain for a given locale object and initializes the Po object.
// If the domain exists, it gets reloaded.
func (l *Locale) AddDomain(dom string) {
	var poObj Translator

	file := l.findExt(dom, "po")
	if file != "" {
		poObj = NewPoFS(l.fs)
		// Parse file.
		poObj.ParseFile(file)
	} else {
		file = l.findExt(dom, "mo")
		if file != "" {
			poObj = NewMoFS(l.fs)
			// Parse file.
			poObj.ParseFile(file)
		} else {
			// fallback return if no file found with
			return
		}
	}

	// Save new domain
	l.Lock()

	if l.Domains == nil {
		l.Domains = make(map[string]Translator)
	}
	if l.defaultDomain == "" {
		l.defaultDomain = dom
	}
	l.Domains[dom] = poObj

	// Unlock "Save new domain"
	l.Unlock()
}

// AddTranslator takes a domain name and a Translator object to make it available in the Locale object.
func (l *Locale) AddTranslator(dom string, tr Translator) {
	l.Lock()

	if l.Domains == nil {
		l.Domains = make(map[string]Translator)
	}
	if l.defaultDomain == "" {
		l.defaultDomain = dom
	}
	l.Domains[dom] = tr

	l.Unlock()
}

// GetDomain is the domain getter for Locale configuration
func (l *Locale) GetDomain() string {
	l.RLock()
	dom := l.defaultDomain
	l.RUnlock()
	return dom
}

// SetDomain sets the name for the domain to be used.
func (l *Locale) SetDomain(dom string) {
	l.Lock()
	l.defaultDomain = dom
	l.Unlock()
}

// GetLanguage is the lang getter for Locale configuration
func (l *Locale) GetLanguage() string {
	l.RLock()
	lang := l.lang
	l.RUnlock()
	return lang
}

// Get uses a domain "default" to return the corresponding Translation of a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) Get(str string, vars ...interface{}) string {
	return l.GetD(l.GetDomain(), str, vars...)
}

// GetN retrieves the (N)th plural form of Translation for the given string in the "default" domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetN(str, plural string, n int, vars ...interface{}) string {
	return l.GetND(l.GetDomain(), str, plural, n, vars...)
}

// GetD returns the corresponding Translation in the given domain for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetD(dom, str string, vars ...interface{}) string {
	// Sync read
	l.RLock()
	defer l.RUnlock()

	if l.Domains != nil {
		if _, ok := l.Domains[dom]; ok {
			if l.Domains[dom] != nil {
				return l.Domains[dom].Get(str, vars...)
			}
		}
	}

	return Printf(str, vars...)
}

// GetND retrieves the (N)th plural form of Translation in the given domain for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetND(dom, str, plural string, n int, vars ...interface{}) string {
	// Sync read
	l.RLock()
	defer l.RUnlock()

	if l.Domains != nil {
		if _, ok := l.Domains[dom]; ok {
			if l.Domains[dom] != nil {
				return l.Domains[dom].GetN(str, plural, n, vars...)
			}
		}
	}

	// Use western default rule (plural > 1) to handle missing domain default result.
	if n == 1 {
		return Printf(str, vars...)
	}
	return Printf(plural, vars...)
}

// GetC uses a domain "default" to return the corresponding Translation of the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetC(str, ctx string, vars ...interface{}) string {
	return l.GetDC(l.GetDomain(), str, ctx, vars...)
}

// GetNC retrieves the (N)th plural form of Translation for the given string in the given context in the "default" domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetNC(str, plural string, n int, ctx string, vars ...interface{}) string {
	return l.GetNDC(l.GetDomain(), str, plural, n, ctx, vars...)
}

// GetDC returns the corresponding Translation in the given domain for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetDC(dom, str, ctx string, vars ...interface{}) string {
	// Sync read
	l.RLock()
	defer l.RUnlock()

	if l.Domains != nil {
		if _, ok := l.Domains[dom]; ok {
			if l.Domains[dom] != nil {
				return l.Domains[dom].GetC(str, ctx, vars...)
			}
		}
	}

	return Printf(str, vars...)
}

// GetNDC retrieves the (N)th plural form of Translation in the given domain for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetNDC(dom, str, plural string, n int, ctx string, vars ...interface{}) string {
	// Sync read
	l.RLock()
	defer l.RUnlock()

	if l.Domains != nil {
		if _, ok := l.Domains[dom]; ok {
			if l.Domains[dom] != nil {
				return l.Domains[dom].GetNC(str, plural, n, ctx, vars...)
			}
		}
	}

	// Use western default rule (plural > 1) to handle missing domain default result.
	if n == 1 {
		return Printf(str, vars...)
	}
	return Printf(plural, vars...)
}

// GetTranslations returns a copy of all translations in all domains of this locale. It does not support contexts.
func (l *Locale) GetTranslations() map[string]*Translation {
	all := make(map[string]*Translation)

	l.RLock()
	defer l.RUnlock()
	for _, translator := range l.Domains {
		for msgID, trans := range translator.GetDomain().GetTranslations() {
			all[msgID] = trans
		}
	}

	return all
}

// IsTranslated reports whether a string is translated
func (l *Locale) IsTranslated(str string) bool {
	return l.IsTranslatedND(l.GetDomain(), str, 0)
}

// IsTranslatedN reports whether a plural string is translated
func (l *Locale) IsTranslatedN(str string, n int) bool {
	return l.IsTranslatedND(l.GetDomain(), str, n)
}

// IsTranslatedD reports whether a domain string is translated
func (l *Locale) IsTranslatedD(dom, str string) bool {
	return l.IsTranslatedND(dom, str, 0)
}

// IsTranslatedND reports whether a plural domain string is translated
func (l *Locale) IsTranslatedND(dom, str string, n int) bool {
	l.RLock()
	defer l.RUnlock()

	if l.Domains == nil {
		return false
	}
	translator, ok := l.Domains[dom]
	if !ok {
		return false
	}
	return translator.GetDomain().IsTranslatedN(str, n)
}

// IsTranslatedC reports whether a context string is translated
func (l *Locale) IsTranslatedC(str, ctx string) bool {
	return l.IsTranslatedNDC(l.GetDomain(), str, 0, ctx)
}

// IsTranslatedNC reports whether a plural context string is translated
func (l *Locale) IsTranslatedNC(str string, n int, ctx string) bool {
	return l.IsTranslatedNDC(l.GetDomain(), str, n, ctx)
}

// IsTranslatedDC reports whether a domain context string is translated
func (l *Locale) IsTranslatedDC(dom, str, ctx string) bool {
	return l.IsTranslatedNDC(dom, str, 0, ctx)
}

// IsTranslatedNDC reports whether a plural domain context string is translated
func (l *Locale) IsTranslatedNDC(dom string, str string, n int, ctx string) bool {
	l.RLock()
	defer l.RUnlock()

	if l.Domains == nil {
		return false
	}
	translator, ok := l.Domains[dom]
	if !ok {
		return false
	}
	return translator.GetDomain().IsTranslatedNC(str, n, ctx)
}

// LocaleEncoding is used as intermediary storage to encode Locale objects to Gob.
type LocaleEncoding struct {
	Path          string
	Lang          string
	Domains       map[string][]byte
	DefaultDomain string
}

// MarshalBinary implements encoding BinaryMarshaler interface
func (l *Locale) MarshalBinary() ([]byte, error) {
	obj := new(LocaleEncoding)
	obj.DefaultDomain = l.defaultDomain
	obj.Domains = make(map[string][]byte)
	for k, v := range l.Domains {
		var err error
		obj.Domains[k], err = v.MarshalBinary()
		if err != nil {
			return nil, err
		}
	}
	obj.Lang = l.lang
	obj.Path = l.path

	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(obj)

	return buff.Bytes(), err
}

// UnmarshalBinary implements encoding BinaryUnmarshaler interface
func (l *Locale) UnmarshalBinary(data []byte) error {
	buff := bytes.NewBuffer(data)
	obj := new(LocaleEncoding)

	decoder := gob.NewDecoder(buff)
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}

	l.defaultDomain = obj.DefaultDomain
	l.lang = obj.Lang
	l.path = obj.Path

	// Decode Domains
	l.Domains = make(map[string]Translator)
	for k, v := range obj.Domains {
		var tr TranslatorEncoding
		buff := bytes.NewBuffer(v)
		trDecoder := gob.NewDecoder(buff)
		err := trDecoder.Decode(&tr)
		if err != nil {
			return err
		}

		l.Domains[k] = tr.GetTranslator()
	}

	return nil
}
