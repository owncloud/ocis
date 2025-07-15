/*
 * Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com. All rights reserved.
 * Licensed under the MIT License. See LICENSE file in the project root for full license information.
 */

package gotext

// Translation is the struct for the Translations parsed via Po or Mo files and all coming parsers
type Translation struct {
	ID       string
	PluralID string
	Trs      map[int]string
	Refs     []string

	dirty bool
}

// NewTranslation returns the Translation object and initialized it.
func NewTranslation() *Translation {
	return &Translation{
		Trs: make(map[int]string),
	}
}

// NewTranslationWithRefs returns the Translation object and initialized it with references.
func NewTranslationWithRefs(refs []string) *Translation {
	return &Translation{
		Trs:  make(map[int]string),
		Refs: refs,
	}
}

// IsStale returns whether the translation is stale or not
func (t *Translation) IsStale() bool {
	return !t.dirty
}

// SetRefs sets the references of the translation
func (t *Translation) SetRefs(refs []string) {
	t.Refs = refs
	t.dirty = true
}

// Set sets the string of the translation
func (t *Translation) Set(str string) {
	t.Trs[0] = str
	t.dirty = true
}

// Get returns the string of the translation
func (t *Translation) Get() string {
	// Look for Translation index 0
	if _, ok := t.Trs[0]; ok {
		if t.Trs[0] != "" {
			return t.Trs[0]
		}
	}

	// Return untranslated id by default
	return t.ID
}

// SetN sets the string of the plural translation
func (t *Translation) SetN(n int, str string) {
	t.Trs[n] = str
	t.dirty = true
}

// GetN returns the string of the plural translation
func (t *Translation) GetN(n int) string {
	// Look for Translation index
	if _, ok := t.Trs[n]; ok {
		if t.Trs[n] != "" {
			return t.Trs[n]
		}
	}

	// Return untranslated singular if corresponding
	if n == 0 {
		return t.ID
	}

	// Return untranslated plural by default
	return t.PluralID
}

// IsTranslated reports whether a string is translated
func (t *Translation) IsTranslated() bool {
	tr, ok := t.Trs[0]
	return tr != "" && ok
}

// IsTranslatedN reports whether a plural string is translated
func (t *Translation) IsTranslatedN(n int) bool {
	tr, ok := t.Trs[n]
	return tr != "" && ok
}
