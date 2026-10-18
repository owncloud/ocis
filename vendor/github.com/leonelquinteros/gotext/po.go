/*
 * Copyright (c) 2018 DeineAgentur UG https://www.deineagentur.com. All rights reserved.
 * Licensed under the MIT License. See LICENSE file in the project root for full license information.
 */

package gotext

import (
	"io/fs"
	"strconv"
	"strings"
)

/*
Po parses the content of any PO file and provides all the Translation functions needed.
It's the base object used by all package methods.
And it's safe for concurrent use by multiple goroutines by using the sync package for locking.

Example:

	import (
		"fmt"
		"github.com/leonelquinteros/gotext"
	)

	func main() {
		// Create po object
		po := gotext.NewPo()

		// Parse .po file
		po.ParseFile("/path/to/po/file/translations.po")

		// Get Translation
		fmt.Println(po.Get("Translate this"))
	}
*/
type Po struct {
	// these three public members are for backwards compatibility. they are just set to the value in the domain
	Headers     HeaderMap
	Language    string
	PluralForms string

	domain *Domain
	fs     fs.FS

	parseBufferHasID      bool
	parseBufferHasPlural  bool
	parseBufferHasContext bool
}

type parseState int

const (
	head parseState = iota
	msgCtxt
	msgID
	msgIDPlural
	msgStr
	msgStrInvalid
)

// NewPo should always be used to instantiate a new Po object
func NewPo() *Po {
	po := new(Po)
	po.domain = NewDomain()

	return po
}

// NewPoFS works like NewPO but adds an optional fs.FS
func NewPoFS(filesystem fs.FS) *Po {
	po := NewPo()
	po.fs = filesystem
	return po
}

// GetDomain returns the domain object
func (po *Po) GetDomain() *Domain {
	return po.domain
}

// Convenience interfaces
// ---------------------------------------------------------------

// DropStaleTranslations removes all translations that are not referenced in the current domain
func (po *Po) DropStaleTranslations() {
	po.domain.DropStaleTranslations()
}

// SetRefs sets the references for a given translation
func (po *Po) SetRefs(str string, refs []string) {
	po.domain.SetRefs(str, refs)
}

// GetRefs returns the references for a given translation
func (po *Po) GetRefs(str string) []string {
	return po.domain.GetRefs(str)
}

// SetPluralResolver sets the plural resolver function
func (po *Po) SetPluralResolver(f func(int) int) {
	po.domain.SetPluralResolver(f)
}

// Set translation
func (po *Po) Set(id, str string) {
	po.domain.Set(id, str)
}

// Get translation
func (po *Po) Get(str string, vars ...any) string {
	return po.domain.Get(str, vars...)
}

// Append translation
func (po *Po) Append(b []byte, str string, vars ...any) []byte {
	return po.domain.Append(b, str, vars...)
}

// SetN sets the plural translation
func (po *Po) SetN(id, plural string, n int, str string) {
	po.domain.SetN(id, plural, n, str)
}

// GetN gets the plural translation
func (po *Po) GetN(str, plural string, n int, vars ...any) string {
	return po.domain.GetN(str, plural, n, vars...)
}

// AppendN appends the plural translation
func (po *Po) AppendN(b []byte, str, plural string, n int, vars ...any) []byte {
	return po.domain.AppendN(b, str, plural, n, vars...)
}

// SetC sets the translation for a given context
func (po *Po) SetC(id, ctx, str string) {
	po.domain.SetC(id, ctx, str)
}

// GetC gets the translation for a given context
func (po *Po) GetC(str, ctx string, vars ...any) string {
	return po.domain.GetC(str, ctx, vars...)
}

// AppendC appends the translation for a given context
func (po *Po) AppendC(b []byte, str, ctx string, vars ...any) []byte {
	return po.domain.AppendC(b, str, ctx, vars...)
}

// SetNC sets the plural translation for a given context
func (po *Po) SetNC(id, plural, ctx string, n int, str string) {
	po.domain.SetNC(id, plural, ctx, n, str)
}

// GetNC gets the plural translation for a given context
func (po *Po) GetNC(str, plural string, n int, ctx string, vars ...any) string {
	return po.domain.GetNC(str, plural, n, ctx, vars...)
}

// AppendNC appends the plural translation for a given context
func (po *Po) AppendNC(b []byte, str, plural string, n int, ctx string, vars ...any) []byte {
	return po.domain.AppendNC(b, str, plural, n, ctx, vars...)
}

// IsTranslated checks if the given string is translated
func (po *Po) IsTranslated(str string) bool {
	return po.domain.IsTranslated(str)
}

// IsTranslatedN checks if the given string is translated with plural form
func (po *Po) IsTranslatedN(str string, n int) bool {
	return po.domain.IsTranslatedN(str, n)
}

// IsTranslatedC checks if the given string is translated with context
func (po *Po) IsTranslatedC(str, ctx string) bool {
	return po.domain.IsTranslatedC(str, ctx)
}

// IsTranslatedNC checks if the given string is translated with plural form and context
func (po *Po) IsTranslatedNC(str string, n int, ctx string) bool {
	return po.domain.IsTranslatedNC(str, n, ctx)
}

// MarshalText marshals the Po object to text
func (po *Po) MarshalText() ([]byte, error) {
	return po.domain.MarshalText()
}

// MarshalBinary marshals the Po object to binary
func (po *Po) MarshalBinary() ([]byte, error) {
	return po.domain.MarshalBinary()
}

// UnmarshalBinary unmarshals the Po object from binary
func (po *Po) UnmarshalBinary(data []byte) error {
	return po.domain.UnmarshalBinary(data)
}

// ParseFile loads the translations from a file
func (po *Po) ParseFile(f string) {
	data, err := getFileData(f, po.fs)
	if err != nil {
		return
	}

	po.Parse(data)
}

// Parse loads the translations specified in the provided byte slice (buf)
func (po *Po) Parse(buf []byte) {
	if po.domain == nil {
		panic("NewPo() was not used to instantiate this object")
	}

	// Lock while parsing
	po.domain.trMutex.Lock()
	po.domain.pluralMutex.Lock()
	defer po.domain.trMutex.Unlock()
	defer po.domain.pluralMutex.Unlock()

	// Init buffer
	po.domain.trBuffer = NewTranslation()
	po.domain.ctxBuffer = ""
	po.domain.refBuffer = ""
	po.parseBufferHasID = false
	po.parseBufferHasPlural = false
	po.parseBufferHasContext = false

	state := head
	activeMsgStrIndex := 0
	for l := range strings.Lines(string(buf)) {
		// Trim spaces
		l = strings.TrimSpace(l)

		// Skip invalid lines
		if !po.isValidLine(l) {
			po.parseComment(l, state)
			continue
		}

		// Multi line strings and headers
		if strings.HasPrefix(l, "\"") {
			if !po.parseString(l, state, activeMsgStrIndex) && state != head {
				state = msgStrInvalid
			}
			continue
		}

		// Buffer context and continue
		if po.hasKeywordBoundary(l, "msgctxt") {
			if po.parseContext(l) {
				state = msgCtxt
			} else {
				state = msgStrInvalid
			}
			continue
		}

		// Check for plural form
		if po.hasKeywordBoundary(l, "msgid_plural") {
			if state != msgID || !po.parsePluralID(l) {
				state = msgStrInvalid
			} else {
				state = msgIDPlural
			}
			continue
		}

		// Buffer msgid and continue
		if po.hasKeywordBoundary(l, "msgid") {
			if po.parseID(l) {
				state = msgID
			} else {
				state = msgStrInvalid
			}
			continue
		}

		// Save Translation
		if po.hasKeywordBoundary(l, "msgstr") {
			if state != msgID && state != msgIDPlural && state != msgStr {
				activeMsgStrIndex = 0
				state = msgStrInvalid
				continue
			}

			var ok bool
			activeMsgStrIndex, ok = po.parseMessage(l)
			if ok {
				state = msgStr
			} else {
				activeMsgStrIndex = 0
				state = msgStrInvalid
			}
			continue
		}
	}

	// Save last Translation buffer, but do not synthesize an entry after an
	// invalid directive. Preserve the historical empty-input entry.
	if po.parseBufferHasID {
		po.saveBuffer()
	} else if len(buf) == 0 {
		po.domain.translations[""] = po.domain.trBuffer
	}

	// Parse headers
	po.domain.parseHeaders()

	// set values on this struct
	// this is for backwards compatibility
	po.Language = po.domain.Language
	po.PluralForms = po.domain.PluralForms
	po.Headers = po.domain.Headers
}

// saveBuffer takes the context and Translation buffers
// and saves it on the translations collection
func (po *Po) saveBuffer() {
	if !po.parseBufferHasID {
		return
	}

	// Store a completed plural ID after all of its continuation lines have
	// been decoded.
	if po.parseBufferHasPlural {
		po.domain.pluralTranslations[po.domain.trBuffer.PluralID] = po.domain.trBuffer
	}

	if po.parseBufferHasContext {
		if _, ok := po.domain.contextTranslations[po.domain.ctxBuffer]; !ok {
			po.domain.contextTranslations[po.domain.ctxBuffer] = make(map[string]*Translation)
		}
		po.domain.contextTranslations[po.domain.ctxBuffer][po.domain.trBuffer.ID] = po.domain.trBuffer
		// Context applies only to the record it prefixes, even when its ID is empty.
		po.domain.ctxBuffer = ""
		po.parseBufferHasContext = false
	} else {
		po.domain.translations[po.domain.trBuffer.ID] = po.domain.trBuffer
	}

	po.parseBufferHasID = false
	po.parseBufferHasPlural = false
	// Flush Translation buffer
	if po.domain.refBuffer == "" {
		po.domain.trBuffer = NewTranslation()
	} else {
		po.domain.trBuffer = NewTranslationWithRefs(strings.Split(po.domain.refBuffer, " "))
	}
}

// Either preserves comments before the first "msgid", for later round-trip.
// Or preserves source references for a given translation.
func (po *Po) parseComment(l string, state parseState) {
	if len(l) > 0 && l[0] == '#' {
		if state == head {
			po.domain.headerComments = append(po.domain.headerComments, l)
		} else if len(l) > 1 {
			switch l[1] {
			case ':':
				if len(l) > 2 {
					po.domain.refBuffer = strings.TrimSpace(l[2:])
				}
			}
		}
	}
}

// parseContext takes a line starting with "msgctxt",
// saves the current Translation buffer and creates a new context.
func (po *Po) parseContext(l string) bool {
	value, ok := po.parseQuotedDirective(l, "msgctxt")
	if !ok {
		return false
	}

	// Save current Translation buffer.
	po.saveBuffer()

	// Buffer context, including an explicitly empty context.
	po.domain.ctxBuffer = value
	po.parseBufferHasContext = true
	return true
}

// parseID takes a line starting with "msgid",
// saves the current Translation and creates a new msgid buffer.
func (po *Po) parseID(l string) bool {
	value, ok := po.parseQuotedDirective(l, "msgid")
	if !ok {
		return false
	}

	// Save current Translation buffer.
	po.saveBuffer()

	// Set id
	po.domain.trBuffer.ID = value
	po.parseBufferHasID = true
	return true
}

// parsePluralID saves the plural id buffer from a line starting with "msgid_plural"
func (po *Po) parsePluralID(l string) bool {
	value, ok := po.parseQuotedDirective(l, "msgid_plural")
	if !ok {
		return false
	}

	po.domain.trBuffer.PluralID = value
	po.parseBufferHasPlural = true
	return true
}

// parseMessage takes a line starting with "msgstr" and saves it into the current buffer.
func (po *Po) parseMessage(l string) (int, bool) {
	if !po.hasKeywordBoundary(l, "msgstr") {
		return 0, false
	}
	l = strings.TrimSpace(l[len("msgstr"):])
	// Check for indexed Translation forms
	if strings.HasPrefix(l, "[") {
		idx := strings.IndexByte(l, ']')
		if idx <= 1 || !isASCIIPluralIndex(l[1:idx]) {
			// Skip wrong index formatting
			return 0, false
		}

		// Parse index
		i, err := strconv.Atoi(l[1:idx])
		if err != nil {
			// Skip wrong index formatting
			return 0, false
		}

		// Parse Translation string
		clean, ok := decodePOQuotedString(strings.TrimSpace(l[idx+1:]))
		if !ok {
			return 0, false
		}
		po.domain.trBuffer.Trs[i] = clean

		return i, true
	}

	// Save single Translation form under 0 index
	clean, ok := decodePOQuotedString(l)
	if !ok {
		return 0, false
	}
	po.domain.trBuffer.Trs[0] = clean
	return 0, true
}

// parseString takes a well formatted string without prefix
// and creates headers or attach multi-line strings when corresponding
func (po *Po) parseString(l string, state parseState, activeMsgStrIndex int) bool {
	if state == msgStrInvalid {
		return false
	}

	clean, ok := decodePOQuotedString(l)
	if !ok {
		return false
	}

	switch state {
	case msgStr:
		// Append to active Translation form
		po.domain.trBuffer.Trs[activeMsgStrIndex] += clean

	case msgID:
		// Multiline msgid - Append to current id
		po.domain.trBuffer.ID += clean

	case msgIDPlural:
		// Multiline msgid - Append to current id
		po.domain.trBuffer.PluralID += clean

	case msgCtxt:
		// Multiline context - Append to current context
		po.domain.ctxBuffer += clean
	}

	return true
}

func (po *Po) isValidLine(l string) bool {
	if strings.HasPrefix(l, "\"") {
		return true
	}

	for _, keyword := range []string{
		"msgctxt",
		"msgid_plural",
		"msgid",
		"msgstr",
	} {
		if po.hasKeywordBoundary(l, keyword) {
			return true
		}
	}

	return false
}

func (po *Po) hasKeywordBoundary(l, keyword string) bool {
	if !strings.HasPrefix(l, keyword) {
		return false
	}
	if len(l) == len(keyword) {
		return true
	}

	next := l[len(keyword)]
	if keyword == "msgstr" && next == '[' {
		return true
	}
	return next == ' ' || next == '\t'
}

func (po *Po) parseQuotedDirective(l, keyword string) (string, bool) {
	if !po.hasKeywordBoundary(l, keyword) {
		return "", false
	}

	value, ok := decodePOQuotedString(strings.TrimSpace(l[len(keyword):]))
	if !ok {
		return "", false
	}
	return value, true
}

// decodePOQuotedString decodes a GNU PO double-quoted string.
func decodePOQuotedString(source string) (string, bool) {
	if len(source) < 2 || source[0] != '"' {
		return "", false
	}

	var decoded []byte
	for i := 1; i < len(source); {
		switch source[i] {
		case '"':
			if i != len(source)-1 {
				return "", false
			}
			if decoded == nil {
				return source[1:i], true
			}
			return string(decoded), true

		case '\\':
			if decoded == nil {
				decoded = make([]byte, 0, len(source)-2)
				decoded = append(decoded, source[1:i]...)
			}
			i++
			if i >= len(source) {
				return "", false
			}

			switch source[i] {
			case 'a':
				decoded = append(decoded, '\a')
				i++
			case 'b':
				decoded = append(decoded, '\b')
				i++
			case 'f':
				decoded = append(decoded, '\f')
				i++
			case 'n':
				decoded = append(decoded, '\n')
				i++
			case 'r':
				decoded = append(decoded, '\r')
				i++
			case 't':
				decoded = append(decoded, '\t')
				i++
			case 'v':
				decoded = append(decoded, '\v')
				i++
			case '\\':
				decoded = append(decoded, '\\')
				i++
			case '"':
				decoded = append(decoded, '"')
				i++
			case '0', '1', '2', '3', '4', '5', '6', '7':
				var value byte
				for digits := 0; i < len(source) && digits < 3; digits++ {
					if source[i] < '0' || source[i] > '7' {
						break
					}
					value = value<<3 | (source[i] - '0')
					i++
				}
				decoded = append(decoded, value)
			case 'x':
				i++
				if i >= len(source) || !isPOHexDigit(source[i]) {
					return "", false
				}
				var value byte
				for i < len(source) && isPOHexDigit(source[i]) {
					value = value<<4 | poHexValue(source[i])
					i++
				}
				decoded = append(decoded, value)
			default:
				return "", false
			}

		default:
			if decoded != nil {
				decoded = append(decoded, source[i])
			}
			i++
		}
	}

	return "", false
}

func isPOHexDigit(c byte) bool {
	return c >= '0' && c <= '9' ||
		c >= 'a' && c <= 'f' ||
		c >= 'A' && c <= 'F'
}

func poHexValue(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	default:
		return c - 'A' + 10
	}
}

func isASCIIPluralIndex(index string) bool {
	if index == "" {
		return false
	}
	for i := 0; i < len(index); i++ {
		if index[i] < '0' || index[i] > '9' {
			return false
		}
	}
	return true
}
