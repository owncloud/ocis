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

// FallbackLocale is the default language to be used when no language is set.
var FallbackLocale = "en_US"

func init() {
	// Init default configuration
	globalConfig = &config{
		domain:    "default",
		languages: []string{FallbackLocale},
		library:   "/usr/local/share/locale",
		locales:   make([]*Locale, 0),
	}

	// Register Translator types for gob encoding
	gob.Register(TranslatorEncoding{})
}

// loadLocales ensures package-level locales and their configured domain are loaded.
func loadLocales(rebuildCache bool) {
	globalConfig.Lock()

	if globalConfig.locales == nil || rebuildCache {
		locales := make([]*Locale, 0, len(globalConfig.languages))
		for _, language := range globalConfig.languages {
			locales = append(locales, NewLocale(globalConfig.library, language))
		}
		globalConfig.locales = locales
	}

	locales := append([]*Locale(nil), globalConfig.locales...)
	domain := globalConfig.domain
	globalConfig.Unlock()

	for _, locale := range locales {
		if locale == nil {
			continue
		}
		if rebuildCache || !locale.hasDomain(domain) {
			locale.AddDomain(domain)
		}
		locale.SetDomain(domain)
	}
}

// GetDomain is the domain getter for the package configuration
func GetDomain() string {
	globalConfig.RLock()
	domain := globalConfig.domain
	var locale *Locale
	for _, candidate := range globalConfig.locales {
		if candidate != nil {
			locale = candidate
			break
		}
	}
	globalConfig.RUnlock()

	if locale != nil {
		if localeDomain := locale.GetDomain(); localeDomain != "" {
			return localeDomain
		}
	}
	return domain
}

// SetDomain sets the name for the domain to be used at package level.
// It reloads the corresponding Translation file.
func SetDomain(dom string) {
	globalConfig.Lock()
	globalConfig.domain = dom
	for _, locale := range globalConfig.locales {
		if locale != nil {
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
	languages := make([]string, len(globalConfig.languages))
	copy(languages, globalConfig.languages)
	globalConfig.RUnlock()
	return languages
}

// SetLanguage sets the language code (or colon separated language codes) to be used at package level.
// It reloads the corresponding Translation file.
func SetLanguage(lang string) {
	globalConfig.Lock()
	var languages []string
	for language := range strings.SplitSeq(lang, ":") {
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

// GetLocales returns the locales that have been set for the package configuration.
func GetLocales() []*Locale {
	globalConfig.RLock()
	locales := make([]*Locale, len(globalConfig.locales))
	copy(locales, globalConfig.locales)
	globalConfig.RUnlock()
	return locales
}

// GetStorage is the locale storage getter for the package configuration.
//
// Deprecated: Storage has been renamed to Locale for consistency, use GetLocales instead.
func GetStorage() *Locale {
	for _, locale := range GetLocales() {
		if locale != nil {
			return locale
		}
	}
	return nil
}

// SetLocales allows for overriding the global Locale objects with ones built manually with
// NewLocale(). This makes it possible to attach custom Domain objects from in-memory po/mo.
// The library, language and domain of the first Locale will set the default global configuration.
func SetLocales(locales []*Locale) {
	ownedLocales := make([]*Locale, len(locales))
	copy(ownedLocales, locales)

	globalConfig.Lock()
	globalConfig.locales = ownedLocales

	languages := make([]string, 0, len(ownedLocales))
	var first *Locale
	for _, locale := range ownedLocales {
		if locale == nil {
			continue
		}

		locale.RLock()
		if first == nil {
			first = locale
			globalConfig.library = locale.path
			globalConfig.domain = locale.defaultDomain
		}
		languages = append(languages, locale.lang)
		locale.RUnlock()
	}
	globalConfig.languages = languages
	globalConfig.Unlock()
}

// SetStorage allows overriding the global Locale object with one built manually with NewLocale().
//
// Deprecated: Storage has been renamed to Locale for consistency, use SetLocales instead.
func SetStorage(locale *Locale) {
	if locale == nil {
		SetLocales(nil)
		return
	}
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
	for language := range strings.SplitSeq(lang, ":") {
		languages = append(languages, SimplifiedLocale(language))
	}
	globalConfig.languages = languages
	globalConfig.domain = dom
	globalConfig.Unlock()

	loadLocales(true)
}

// Get uses the default domain globally set to return the corresponding Translation of a given string.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func Get(str string, vars ...any) string {
	return GetD(GetDomain(), str, vars...)
}

// GetN retrieves the (N)th plural form of Translation for the given string in the default domain.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func GetN(str, plural string, n int, vars ...any) string {
	return GetND(GetDomain(), str, plural, n, vars...)
}

// GetD returns the corresponding Translation in the given domain for a given string.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func GetD(dom, str string, vars ...any) string {
	// Try to load default package Locales
	loadLocales(false)

	locales := GetLocales()
	var fallback *Locale
	for _, locale := range locales {
		if locale == nil {
			continue
		}
		fallback = locale
		if !locale.hasDomain(dom) {
			locale.AddDomain(dom)
		}
		if locale.IsTranslatedD(dom, str) {
			return locale.GetD(dom, str, vars...)
		}
	}
	if fallback != nil {
		return fallback.GetD(dom, str, vars...)
	}
	return FormatString(str, vars...)
}

// GetND retrieves the (N)th plural form of Translation in the given domain for a given string.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func GetND(dom, str, plural string, n int, vars ...any) string {
	// Try to load default package Locales
	loadLocales(false)

	locales := GetLocales()
	var fallback *Locale
	for _, locale := range locales {
		if locale == nil {
			continue
		}
		fallback = locale
		if !locale.hasDomain(dom) {
			locale.AddDomain(dom)
		}
		if locale.IsTranslatedND(dom, str, n) {
			return locale.GetND(dom, str, plural, n, vars...)
		}
	}
	if fallback != nil {
		return fallback.GetND(dom, str, plural, n, vars...)
	}
	if n == 1 {
		return FormatString(str, vars...)
	}
	return FormatString(plural, vars...)
}

// GetC uses the default domain globally set to return the corresponding Translation of the given string in the given context.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func GetC(str, ctx string, vars ...any) string {
	return GetDC(GetDomain(), str, ctx, vars...)
}

// GetNC retrieves the (N)th plural form of Translation for the given string in the given context in the default domain.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func GetNC(str, plural string, n int, ctx string, vars ...any) string {
	return GetNDC(GetDomain(), str, plural, n, ctx, vars...)
}

// GetDC returns the corresponding Translation in the given domain for the given string in the given context.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func GetDC(dom, str, ctx string, vars ...any) string {
	// Try to load default package Locales
	loadLocales(false)

	locales := GetLocales()
	var fallback *Locale
	for _, locale := range locales {
		if locale == nil {
			continue
		}
		fallback = locale
		if !locale.hasDomain(dom) {
			locale.AddDomain(dom)
		}
		if locale.IsTranslatedDC(dom, str, ctx) {
			return locale.GetDC(dom, str, ctx, vars...)
		}
	}
	if fallback != nil {
		return fallback.GetDC(dom, str, ctx, vars...)
	}
	return FormatString(str, vars...)
}

// GetNDC retrieves the (N)th plural form of Translation in the given domain for a given string.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func GetNDC(dom, str, plural string, n int, ctx string, vars ...any) string {
	// Try to load default package Locales
	loadLocales(false)

	locales := GetLocales()
	var fallback *Locale
	for _, locale := range locales {
		if locale == nil {
			continue
		}
		fallback = locale
		if !locale.hasDomain(dom) {
			locale.AddDomain(dom)
		}
		if locale.IsTranslatedNDC(dom, str, n, ctx) {
			return locale.GetNDC(dom, str, plural, n, ctx, vars...)
		}
	}
	if fallback != nil {
		return fallback.GetNDC(dom, str, plural, n, ctx, vars...)
	}
	if n == 1 {
		return FormatString(str, vars...)
	}
	return FormatString(plural, vars...)
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
	locales := GetLocales()
	for _, lang := range langs {
		lang = SimplifiedLocale(lang)

		for _, supportedLocale := range locales {
			if supportedLocale == nil {
				continue
			}
			if lang != supportedLocale.GetLanguage() && lang != supportedLocale.GetActualLanguage(dom) {
				continue
			}
			if !supportedLocale.hasDomain(dom) {
				supportedLocale.AddDomain(dom)
			}
			if supportedLocale.IsTranslatedND(dom, str, n) {
				return true
			}
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
	return IsTranslatedNDC(dom, str, 1, ctx, langs...)
}

// IsTranslatedNDC reports whether a plural domain context string is translated in any of given languages.
// When the langs argument is omitted, the output of GetLanguages is used.
func IsTranslatedNDC(dom, str string, n int, ctx string, langs ...string) bool {
	if len(langs) == 0 {
		langs = GetLanguages()
	}

	loadLocales(false)
	locales := GetLocales()
	for _, lang := range langs {
		lang = SimplifiedLocale(lang)

		for _, supportedLocale := range locales {
			if supportedLocale == nil {
				continue
			}
			if lang != supportedLocale.GetLanguage() && lang != supportedLocale.GetActualLanguage(dom) {
				continue
			}
			if !supportedLocale.hasDomain(dom) {
				supportedLocale.AddDomain(dom)
			}
			if supportedLocale.IsTranslatedNDC(dom, str, n, ctx) {
				return true
			}
		}
	}
	return false
}
