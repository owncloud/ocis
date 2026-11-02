/*
 * Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com. All rights reserved.
 * Licensed under the MIT License. See LICENSE file in the project root for full license information.
 */

package gotext

import (
	"bytes"
	"encoding/binary"
	"io/fs"
)

const (
	// MoMagicLittleEndian encoding
	MoMagicLittleEndian = 0x950412de
	// MoMagicBigEndian encoding
	MoMagicBigEndian = 0xde120495

	// EotSeparator msgctxt and msgid separator
	EotSeparator = "\x04"
	// NulSeparator msgid and msgstr separator
	NulSeparator = "\x00"
)

/*
Mo parses the content of any MO file and provides all the Translation functions needed.
It's the base object used by all package methods.
And it's safe for concurrent use by multiple goroutines by using the sync package for locking.

Example:

	import (
		"fmt"
		"github.com/leonelquinteros/gotext"
	)

	func main() {
		// Create mo object
		mo := gotext.NewMo()

		// Parse .mo file
		mo.ParseFile("/path/to/po/file/translations.mo")

		// Get Translation
		fmt.Println(mo.Get("Translate this"))
	}
*/
type Mo struct {
	// these three public members are for backwards compatibility. they are just set to the value in the domain
	Headers     HeaderMap
	Language    string
	PluralForms string
	domain      *Domain
	fs          fs.FS
}

// NewMo should always be used to instantiate a new Mo object
func NewMo() *Mo {
	mo := new(Mo)
	mo.domain = NewDomain()

	return mo
}

// NewMoFS works like NewMO but adds an optional fs.FS
func NewMoFS(filesystem fs.FS) *Mo {
	mo := NewMo()
	mo.fs = filesystem
	return mo
}

// GetDomain returns the domain object
func (mo *Mo) GetDomain() *Domain {
	return mo.domain
}

// Get returns the translation for the given string
func (mo *Mo) Get(str string, vars ...any) string {
	return mo.domain.Get(str, vars...)
}

// Append a translation string into the domain
func (mo *Mo) Append(b []byte, str string, vars ...any) []byte {
	return mo.domain.Append(b, str, vars...)
}

// GetN returns the translation for the given string and plural form
func (mo *Mo) GetN(str, plural string, n int, vars ...any) string {
	return mo.domain.GetN(str, plural, n, vars...)
}

// AppendN appends a translation string for the given plural form into the domain
func (mo *Mo) AppendN(b []byte, str, plural string, n int, vars ...any) []byte {
	return mo.domain.AppendN(b, str, plural, n, vars...)
}

// GetC returns the translation for the given string and context
func (mo *Mo) GetC(str, ctx string, vars ...any) string {
	return mo.domain.GetC(str, ctx, vars...)
}

// AppendC appends a translation string for the given context into the domain
func (mo *Mo) AppendC(b []byte, str, ctx string, vars ...any) []byte {
	return mo.domain.AppendC(b, str, ctx, vars...)
}

// GetNC returns the translation for the given string, plural form and context
func (mo *Mo) GetNC(str, plural string, n int, ctx string, vars ...any) string {
	return mo.domain.GetNC(str, plural, n, ctx, vars...)
}

// AppendNC appends a translation string for the given plural form and context into the domain
func (mo *Mo) AppendNC(b []byte, str, plural string, n int, ctx string, vars ...any) []byte {
	return mo.domain.AppendNC(b, str, plural, n, ctx, vars...)
}

// IsTranslated checks if the given string is translated
func (mo *Mo) IsTranslated(str string) bool {
	return mo.domain.IsTranslated(str)
}

// IsTranslatedN checks if the given string is translated with plural form
func (mo *Mo) IsTranslatedN(str string, n int) bool {
	return mo.domain.IsTranslatedN(str, n)
}

// IsTranslatedC checks if the given string is translated with context
func (mo *Mo) IsTranslatedC(str, ctx string) bool {
	return mo.domain.IsTranslatedC(str, ctx)
}

// IsTranslatedNC checks if the given string is translated with plural form and context
func (mo *Mo) IsTranslatedNC(str string, n int, ctx string) bool {
	return mo.domain.IsTranslatedNC(str, n, ctx)
}

// MarshalBinary marshals the Mo object into a binary format
func (mo *Mo) MarshalBinary() ([]byte, error) {
	return mo.domain.MarshalBinary()
}

// UnmarshalBinary unmarshals the Mo object from a binary format
func (mo *Mo) UnmarshalBinary(data []byte) error {
	return mo.domain.UnmarshalBinary(data)
}

// ParseFile loads the translations specified in the provided file, in the GNU gettext .mo format
func (mo *Mo) ParseFile(f string) {
	data, err := getFileData(f, mo.fs)
	if err != nil {
		return
	}

	mo.Parse(data)
}

// Parse loads the translations specified in the provided byte slice, in the GNU gettext .mo format
func (mo *Mo) Parse(buf []byte) {
	// Lock while parsing
	mo.domain.trMutex.Lock()
	mo.domain.pluralMutex.Lock()
	defer mo.domain.trMutex.Unlock()
	defer mo.domain.pluralMutex.Unlock()

	if len(buf) < 28 {
		return
	}

	var bo binary.ByteOrder
	switch binary.LittleEndian.Uint32(buf[:4]) {
	case MoMagicLittleEndian:
		bo = binary.LittleEndian
	case MoMagicBigEndian:
		bo = binary.BigEndian
	default:
		return
	}

	var header struct {
		MajorVersion uint16
		MinorVersion uint16
		MsgIDCount   uint32
		MsgIDOffset  uint32
		MsgStrOffset uint32
	}

	bufLen := uint64(len(buf))
	header.MajorVersion = bo.Uint16(buf[4:6])
	header.MinorVersion = bo.Uint16(buf[6:8])
	header.MsgIDCount = bo.Uint32(buf[8:12])
	header.MsgIDOffset = bo.Uint32(buf[12:16])
	header.MsgStrOffset = bo.Uint32(buf[16:20])

	if header.MajorVersion != 0 && header.MajorVersion != 1 {
		return
	}
	if header.MinorVersion != 0 && header.MinorVersion != 1 {
		return
	}

	directorySize := uint64(header.MsgIDCount) * 8
	if uint64(header.MsgIDOffset)+directorySize > bufLen ||
		uint64(header.MsgStrOffset)+directorySize > bufLen {
		return
	}

	for i := uint32(0); i < header.MsgIDCount; i++ {
		msgIDEntry := int(uint64(header.MsgIDOffset) + uint64(i)*8)
		msgIDLen := bo.Uint32(buf[msgIDEntry:])
		msgIDStart := bo.Uint32(buf[msgIDEntry+4:])
		if uint64(msgIDStart) > bufLen ||
			uint64(msgIDLen) > bufLen-uint64(msgIDStart) {
			return
		}

		msgStrEntry := int(uint64(header.MsgStrOffset) + uint64(i)*8)
		msgStrLen := bo.Uint32(buf[msgStrEntry:])
		msgStrStart := bo.Uint32(buf[msgStrEntry+4:])
		if uint64(msgStrStart) > bufLen ||
			uint64(msgStrLen) > bufLen-uint64(msgStrStart) {
			return
		}
	}

	for i := uint32(0); i < header.MsgIDCount; i++ {
		msgIDEntry := int(uint64(header.MsgIDOffset) + uint64(i)*8)
		msgIDLen := bo.Uint32(buf[msgIDEntry:])
		msgIDStart := bo.Uint32(buf[msgIDEntry+4:])
		msgIDEnd := int(uint64(msgIDStart) + uint64(msgIDLen))

		msgStrEntry := int(uint64(header.MsgStrOffset) + uint64(i)*8)
		msgStrLen := bo.Uint32(buf[msgStrEntry:])
		msgStrStart := bo.Uint32(buf[msgStrEntry+4:])
		msgStrEnd := int(uint64(msgStrStart) + uint64(msgStrLen))

		mo.addTranslation(buf[int(msgIDStart):msgIDEnd], buf[int(msgStrStart):msgStrEnd])
	}

	// Parse headers
	mo.domain.parseHeaders()

	// set values on this struct
	// this is for backwards compatibility
	mo.Language = mo.domain.Language
	mo.PluralForms = mo.domain.PluralForms
	mo.Headers = mo.domain.Headers
}

func (mo *Mo) addTranslation(msgid, msgstr []byte) {
	translation := NewTranslation()
	var msgctxt []byte
	var msgidPlural []byte
	hasContext := false

	if eotIndex := bytes.IndexByte(msgid, EotSeparator[0]); eotIndex >= 0 {
		hasContext = true
		msgctxt = msgid[:eotIndex]
		msgid = msgid[eotIndex+len(EotSeparator):]
	}

	if nulIndex := bytes.IndexByte(msgid, NulSeparator[0]); nulIndex >= 0 {
		msgidPlural = msgid[nulIndex+1:]
		msgid = msgid[:nulIndex]
	}
	translation.ID = string(msgid)

	if len(msgidPlural) > 0 {
		translation.PluralID = string(msgidPlural)
	}

	i := 0
	for s := range bytes.SplitSeq(msgstr, []byte(NulSeparator)) {
		translation.Trs[i] = string(s)
		i++
	}

	if hasContext {
		// With context, including an explicitly empty context...
		if _, ok := mo.domain.contextTranslations[string(msgctxt)]; !ok {
			mo.domain.contextTranslations[string(msgctxt)] = make(map[string]*Translation)
		}
		mo.domain.contextTranslations[string(msgctxt)][translation.ID] = translation
	} else {
		mo.domain.translations[translation.ID] = translation
	}
}
