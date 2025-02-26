/*
Package gotext implements GNU gettext utilities.

For quick/simple translations you can use the package level functions directly.

	    import (
		    "fmt"
		    "github.com/leonelquinteros/gotext"
	    )

	    func main() {
	        // Configure package
	        gotext.Configure("/path/to/locales/root/dir", "en_UK", "domain-name")

	        // Translate text from default domain
	        fmt.Println(gotext.Get("My text on 'domain-name' domain"))

	        // Translate text from a different domain without reconfigure
	        fmt.Println(gotext.GetD("domain2", "Another text on a different domain"))
	    }
*/
package gotext

import (
	"encoding/gob"
	"strings"
	"sync"
)

// Global environment variables
type config struct {
	sync.RWMutex

	// Path to library directory where all locale directories and Translation files are.
	library string

	// Default domain to look at when no domain is specified. Used by package level functions.
	domain string

	// Language set.
	languages []string

	// Storage for package level methods
	locales []*Locale
}

var globalConfig *config

var FallbackLocale = "en_US"

func init() {
	// Init default configuration
	globalConfig = &config{
		domain:    "default",
		languages: []string{FallbackLocale},
		library:   "/usr/local/share/locale",
		locales:   nil,
	}

	// Register Translator types for gob encoding
	gob.Register(TranslatorEncoding{})
}

// loadLocales creates a new Locale object for every language (specified using Configure)
// at package level based on the configuration of global configuration .
// It is called when trying to use Get or GetD methods.
func loadLocales(rebuildCache bool) {
	globalConfig.Lock()

	if globalConfig.locales == nil || rebuildCache {
		var locales []*Locale
		for _, language := range globalConfig.languages {
			locales = append(locales, NewLocale(globalConfig.library, language))
		}
		globalConfig.locales = locales
	}

	for _, locale := range globalConfig.locales {
		if _, ok := locale.Domains[globalConfig.domain]; !ok || rebuildCache {
			locale.AddDomain(globalConfig.domain)
		}
		locale.SetDomain(globalConfig.domain)
	}

	globalConfig.Unlock()
}

// GetDomain is the domain getter for the package configuration
func GetDomain() string {
	var dom string
	globalConfig.RLock()
	if globalConfig.locales != nil {
		// All locales have the same domain
		dom = globalConfig.locales[0].GetDomain()
	}
	if dom == "" {
		dom = globalConfig.domain
	}
	globalConfig.RUnlock()

	return dom
}

// SetDomain sets the name for the domain to be used at package level.
// It reloads the corresponding Translation file.
func SetDomain(dom string) {
	globalConfig.Lock()
	globalConfig.domain = dom
	if globalConfig.locales != nil {
		for _, locale := range globalConfig.locales {
			locale.SetDomain(dom)
		}
	}
	globalConfig.Unlock()

	loadLocales(true)
}

// GetLanguage returns the language gotext will translate into.
// If multiple languages have been supplied, the first one will be returned.
// If no language has been supplied, the fallback will be returned.
func GetLanguage() string {
	languages := GetLanguages()
	if len(languages) == 0 {
		return FallbackLocale
	}
	return languages[0]
}

// GetLanguages returns all languages that have been supplied.
func GetLanguages() []string {
	globalConfig.RLock()
	defer globalConfig.RUnlock()
	return globalConfig.languages
}

// SetLanguage sets the language code (or colon separated language codes) to be used at package level.
// It reloads the corresponding Translation file.
func SetLanguage(lang string) {
	globalConfig.Lock()
	var languages []string
	for _, language := range strings.Split(lang, ":") {
		languages = append(languages, SimplifiedLocale(language))
	}
	globalConfig.languages = languages
	globalConfig.Unlock()

	loadLocales(true)
}

// GetLibrary is the library getter for the package configuration
func GetLibrary() string {
	globalConfig.RLock()
	lib := globalConfig.library
	globalConfig.RUnlock()

	return lib
}

// SetLibrary sets the root path for the locale directories and files to be used at package level.
// It reloads the corresponding Translation file.
func SetLibrary(lib string) {
	globalConfig.Lock()
	globalConfig.library = lib
	globalConfig.Unlock()

	loadLocales(true)
}

func GetLocales() []*Locale {
	globalConfig.RLock()
	defer globalConfig.RUnlock()
	return globalConfig.locales
}

// GetStorage is the locale storage getter for the package configuration.
//
// Deprecated: Storage has been renamed to Locale for consistency, use GetLocales instead.
func GetStorage() *Locale {
	locales := GetLocales()
	if len(locales) > 0 {
		return locales[0]
	}
	return nil
}

// SetLocales allows for overriding the global Locale objects with ones built manually with
// NewLocale(). This makes it possible to attach custom Domain objects from in-memory po/mo.
// The library, language and domain of the first Locale will set the default global configuration.
func SetLocales(locales []*Locale) {
	globalConfig.Lock()
	defer globalConfig.Unlock()

	globalConfig.locales = locales
	globalConfig.library = locales[0].path
	globalConfig.domain = locales[0].defaultDomain

	var languages []string
	for _, locale := range locales {
		languages = append(languages, locale.lang)
	}
	globalConfig.languages = languages
}

// SetStorage allows overriding the global Locale object with one built manually with NewLocale().
//
// Deprecated: Storage has been renamed to Locale for consistency, use SetLocales instead.
func SetStorage(locale *Locale) {
	SetLocales([]*Locale{locale})
}

// Configure sets all configuration variables to be used at package level and reloads the corresponding Translation file.
// It receives the library path, language code and domain name.
// This function is recommended to be used when changing more than one setting,
// as using each setter will introduce a I/O overhead because the Translation file will be loaded after each set.
func Configure(lib, lang, dom string) {
	globalConfig.Lock()
	globalConfig.library = lib
	var languages []string
	for _, language := range strings.Split(lang, ":") {
		languages = append(languages, SimplifiedLocale(language))
	}
	globalConfig.languages = languages
	globalConfig.domain = dom
	globalConfig.Unlock()

	loadLocales(true)
}

// Get uses the default domain globally set to return the corresponding Translation of a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func Get(str string, vars ...interface{}) string {
	return GetD(GetDomain(), str, vars...)
}

// GetN retrieves the (N)th plural form of Translation for the given string in the default domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetN(str, plural string, n int, vars ...interface{}) string {
	return GetND(GetDomain(), str, plural, n, vars...)
}

// GetD returns the corresponding Translation in the given domain for a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetD(dom, str string, vars ...interface{}) string {
	// Try to load default package Locales
	loadLocales(false)

	globalConfig.RLock()
	defer globalConfig.RUnlock()

	var tr string
	for i, locale := range globalConfig.locales {
		if _, ok := locale.Domains[dom]; !ok {
			locale.AddDomain(dom)
		}
		if !locale.IsTranslatedD(dom, str) && i < (len(globalConfig.locales)-1) {
			continue
		}
		tr = locale.GetD(dom, str, vars...)
		break
	}
	return tr
}

// GetND retrieves the (N)th plural form of Translation in the given domain for a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetND(dom, str, plural string, n int, vars ...interface{}) string {
	// Try to load default package Locales
	loadLocales(false)

	globalConfig.RLock()
	defer globalConfig.RUnlock()

	var tr string
	for i, locale := range globalConfig.locales {
		if _, ok := locale.Domains[dom]; !ok {
			locale.AddDomain(dom)
		}
		if !locale.IsTranslatedND(dom, str, n) && i < (len(globalConfig.locales)-1) {
			continue
		}
		tr = locale.GetND(dom, str, plural, n, vars...)
		break
	}
	return tr
}

// GetC uses the default domain globally set to return the corresponding Translation of the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetC(str, ctx string, vars ...interface{}) string {
	return GetDC(GetDomain(), str, ctx, vars...)
}

// GetNC retrieves the (N)th plural form of Translation for the given string in the given context in the default domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetNC(str, plural string, n int, ctx string, vars ...interface{}) string {
	return GetNDC(GetDomain(), str, plural, n, ctx, vars...)
}

// GetDC returns the corresponding Translation in the given domain for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetDC(dom, str, ctx string, vars ...interface{}) string {
	// Try to load default package Locales
	loadLocales(false)

	globalConfig.RLock()
	defer globalConfig.RUnlock()

	var tr string
	for i, locale := range globalConfig.locales {
		if !locale.IsTranslatedDC(dom, str, ctx) && i < (len(globalConfig.locales)-1) {
			continue
		}
		tr = locale.GetDC(dom, str, ctx, vars...)
		break
	}
	return tr
}

// GetNDC retrieves the (N)th plural form of Translation in the given domain for a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func GetNDC(dom, str, plural string, n int, ctx string, vars ...interface{}) string {
	// Try to load default package Locales
	loadLocales(false)

	// Return Translation
	globalConfig.RLock()
	defer globalConfig.RUnlock()

	var tr string
	for i, locale := range globalConfig.locales {
		if !locale.IsTranslatedNDC(dom, str, n, ctx) && i < (len(globalConfig.locales)-1) {
			continue
		}
		tr = locale.GetNDC(dom, str, plural, n, ctx, vars...)
		break
	}
	return tr
}

// IsTranslated reports whether a string is translated in given languages.
// When the langs argument is omitted, the output of GetLanguages is used.
func IsTranslated(str string, langs ...string) bool {
	return IsTranslatedND(GetDomain(), str, 1, langs...)
}

// IsTranslatedN reports whether a plural string is translated in given languages.
// When the langs argument is omitted, the output of GetLanguages is used.
func IsTranslatedN(str string, n int, langs ...string) bool {
	return IsTranslatedND(GetDomain(), str, n, langs...)
}

// IsTranslatedD reports whether a domain string is translated in given languages.
// When the langs argument is omitted, the output of GetLanguages is used.
func IsTranslatedD(dom, str string, langs ...string) bool {
	return IsTranslatedND(dom, str, 1, langs...)
}

// IsTranslatedND reports whether a plural domain string is translated in any of given languages.
// When the langs argument is omitted, the output of GetLanguages is used.
func IsTranslatedND(dom, str string, n int, langs ...string) bool {
	if len(langs) == 0 {
		langs = GetLanguages()
	}

	loadLocales(false)

	globalConfig.RLock()
	defer globalConfig.RUnlock()

	for _, lang := range langs {
		lang = SimplifiedLocale(lang)

		for _, supportedLocale := range globalConfig.locales {
			if lang != supportedLocale.GetActualLanguage(dom) {
				continue
			}
			return supportedLocale.IsTranslatedND(dom, str, n)
		}
	}
	return false
}

// IsTranslatedC reports whether a context string is translated in given languages.
// When the langs argument is omitted, the output of GetLanguages is used.
func IsTranslatedC(str, ctx string, langs ...string) bool {
	return IsTranslatedNDC(GetDomain(), str, 1, ctx, langs...)
}

// IsTranslatedNC reports whether a plural context string is translated in given languages.
// When the langs argument is omitted, the output of GetLanguages is used.
func IsTranslatedNC(str string, n int, ctx string, langs ...string) bool {
	return IsTranslatedNDC(GetDomain(), str, n, ctx, langs...)
}

// IsTranslatedDC reports whether a domain context string is translated in given languages.
// When the langs argument is omitted, the output of GetLanguages is used.
func IsTranslatedDC(dom, str, ctx string, langs ...string) bool {
	return IsTranslatedNDC(dom, str, 0, ctx, langs...)
}

// IsTranslatedNDC reports whether a plural domain context string is translated in any of given languages.
// When the langs argument is omitted, the output of GetLanguages is used.
func IsTranslatedNDC(dom, str string, n int, ctx string, langs ...string) bool {
	if len(langs) == 0 {
		langs = GetLanguages()
	}

	loadLocales(false)

	globalConfig.RLock()
	defer globalConfig.RUnlock()

	for _, lang := range langs {
		lang = SimplifiedLocale(lang)

		for _, locale := range globalConfig.locales {
			if lang != locale.GetActualLanguage(dom) {
				continue
			}
			return locale.IsTranslatedNDC(dom, str, n, ctx)
		}
	}
	return false
}
