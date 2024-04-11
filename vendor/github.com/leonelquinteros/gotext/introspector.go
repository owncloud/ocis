package gotext

// IsTranslatedIntrospector is able to determine whether a given string is translated.
// Examples of this introspector are Po and Mo, which are specific to their domain.
// Locale holds multiple domains and also implements IsTranslatedDomainIntrospector.
type IsTranslatedIntrospector interface {
	IsTranslated(str string) bool
	IsTranslatedN(str string, n int) bool
	IsTranslatedC(str, ctx string) bool
	IsTranslatedNC(str string, n int, ctx string) bool
}

// IsTranslatedDomainIntrospector is able to determine whether a given string is translated.
// Example of this introspector is Locale, which holds multiple domains.
// Simpler objects that are domain-specific, like Po or Mo, implement IsTranslatedIntrospector.
type IsTranslatedDomainIntrospector interface {
	IsTranslated(str string) bool
	IsTranslatedN(str string, n int) bool
	IsTranslatedD(dom, str string) bool
	IsTranslatedND(dom, str string, n int) bool
	IsTranslatedC(str, ctx string) bool
	IsTranslatedNC(str string, n int, ctx string) bool
	IsTranslatedDC(dom, str, ctx string) bool
	IsTranslatedNDC(dom, str string, n int, ctx string) bool
}
