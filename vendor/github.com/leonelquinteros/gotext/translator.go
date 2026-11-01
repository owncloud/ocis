/*
 * Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com. All rights reserved.
 * Licensed under the MIT License. See LICENSE file in the project root for full license information.
 */

package gotext

import (
	"errors"
	"io/fs"
	"os"

	"github.com/leonelquinteros/gotext/plurals"
)

// Translator interface is used by Locale and Po objects.
// It contains all methods needed to parse translation sources and obtain corresponding translations.
// It includes encoding.BinaryMarshaler and encoding.BinaryUnmarshaler methods for Locale serialization.
type Translator interface {
	ParseFile(f string)
	Parse(buf []byte)
	Get(str string, vars ...any) string
	GetN(str, plural string, n int, vars ...any) string
	GetC(str, ctx string, vars ...any) string
	GetNC(str, plural string, n int, ctx string, vars ...any) string

	MarshalBinary() ([]byte, error)
	UnmarshalBinary([]byte) error
	GetDomain() *Domain
}

// AppendTranslator interface is used by Locale and Po objects.
// It contains all methods needed to parse translation sources and to append entries to the object.
type AppendTranslator interface {
	Translator
	Append(b []byte, str string, vars ...any) []byte
	AppendN(b []byte, str, plural string, n int, vars ...any) []byte
	AppendC(b []byte, str, ctx string, vars ...any) []byte
	AppendNC(b []byte, str, plural string, n int, ctx string, vars ...any) []byte
}

// TranslatorEncoding is used as intermediary storage to encode Translator objects to Gob.
type TranslatorEncoding struct {
	// Headers storage
	Headers HeaderMap

	// Language header
	Language string

	// Plural-Forms header
	PluralForms string

	// Parsed Plural-Forms header values
	Nplurals int
	Plural   string

	// Storage
	Translations map[string]*Translation
	Contexts     map[string]map[string]*Translation
}

// GetTranslator is used to recover a Translator object after unmarshalling the TranslatorEncoding object.
// Internally uses a Po object as it should be switchable with Mo objects without problem.
// External Translator implementations should be able to serialize into a TranslatorEncoding object in order to
// deserialize into a Po-compatible object.
func (te *TranslatorEncoding) GetTranslator() Translator {
	po := NewPo()
	if te == nil {
		return po
	}

	headers := te.Headers
	if headers == nil {
		headers = make(HeaderMap)
	}
	translations := te.Translations
	if translations == nil {
		translations = make(map[string]*Translation)
	}
	for id, translation := range translations {
		if translation == nil {
			translation = NewTranslation()
			translation.ID = id
			translations[id] = translation
		}
	}
	contexts := te.Contexts
	if contexts == nil {
		contexts = make(map[string]map[string]*Translation)
	}
	for context, translationsForContext := range contexts {
		if translationsForContext == nil {
			contexts[context] = make(map[string]*Translation)
			continue
		}
		for id, translation := range translationsForContext {
			if translation == nil {
				translation = NewTranslation()
				translation.ID = id
				translationsForContext[id] = translation
			}
		}
	}

	po.domain.Headers = headers
	po.domain.Language = te.Language
	po.domain.PluralForms = te.PluralForms
	po.domain.nplurals = te.Nplurals
	po.domain.plural = te.Plural
	po.domain.translations = translations
	po.domain.contextTranslations = contexts

	if expr, err := plurals.Compile(te.Plural); err == nil {
		po.domain.pluralforms = expr
	}

	po.Headers = po.domain.Headers
	po.Language = po.domain.Language
	po.PluralForms = po.domain.PluralForms

	return po
}

// getFileData reads a file and returns the byte slice after doing some basic sanity checking
func getFileData(f string, filesystem fs.FS) ([]byte, error) {
	if filesystem != nil {
		return fs.ReadFile(filesystem, f)
	}

	// Check if file exists
	info, err := os.Stat(f)
	if err != nil {
		return nil, err
	}

	// Check that isn't a directory
	if info.IsDir() {
		return nil, errors.New("cannot parse a directory")
	}

	return os.ReadFile(f)
}
