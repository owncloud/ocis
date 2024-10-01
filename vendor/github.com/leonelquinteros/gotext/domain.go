package gotext

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/leonelquinteros/gotext/plurals"
)

// Domain has all the common functions for dealing with a gettext domain
// it's initialized with a GettextFile (which represents either a Po or Mo file)
type Domain struct {
	Headers HeaderMap

	// Language header
	Language string

	// Plural-Forms header
	PluralForms string

	// Preserve comments at head of PO for round-trip
	headerComments []string

	// Parsed Plural-Forms header values
	nplurals    int
	plural      string
	pluralforms plurals.Expression

	// Storage
	translations        map[string]*Translation
	contextTranslations map[string]map[string]*Translation
	pluralTranslations  map[string]*Translation

	// Sync Mutex
	trMutex     sync.RWMutex
	pluralMutex sync.RWMutex

	// Parsing buffers
	trBuffer  *Translation
	ctxBuffer string
	refBuffer string

	customPluralResolver func(int) int
}

// Preserve MIMEHeader behaviour, without the canonicalisation
type HeaderMap map[string][]string

func (m HeaderMap) Add(key, value string) {
	m[key] = append(m[key], value)
}
func (m HeaderMap) Del(key string) {
	delete(m, key)
}
func (m HeaderMap) Get(key string) string {
	if m == nil {
		return ""
	}
	v := m[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}
func (m HeaderMap) Set(key, value string) {
	m[key] = []string{value}
}
func (m HeaderMap) Values(key string) []string {
	if m == nil {
		return nil
	}
	return m[key]
}

func NewDomain() *Domain {
	domain := new(Domain)

	domain.Headers = make(HeaderMap)
	domain.headerComments = make([]string, 0)
	domain.translations = make(map[string]*Translation)
	domain.contextTranslations = make(map[string]map[string]*Translation)
	domain.pluralTranslations = make(map[string]*Translation)

	return domain
}

func (do *Domain) SetPluralResolver(f func(int) int) {
	do.customPluralResolver = f
}

func (do *Domain) pluralForm(n int) int {
	// do we really need locking here? not sure how this plurals.Expression works, so sticking with it for now
	do.pluralMutex.RLock()
	defer do.pluralMutex.RUnlock()

	// Failure fallback
	if do.pluralforms == nil {
		if do.customPluralResolver != nil {
			return do.customPluralResolver(n)
		}

		/* Use the Germanic plural rule.  */
		if n == 1 {
			return 0
		}
		return 1
	}
	return do.pluralforms.Eval(uint32(n))
}

// parseHeaders retrieves data from previously parsed headers. it's called by both Mo and Po when parsing
func (do *Domain) parseHeaders() {
	raw := ""
	if _, ok := do.translations[raw]; ok {
		raw = do.translations[raw].Get()
	}

	// textproto.ReadMIMEHeader() forces keys through CanonicalMIMEHeaderKey(); must read header manually to have one-to-one round-trip of keys
	languageKey := "Language"
	pluralFormsKey := "Plural-Forms"

	rawLines := strings.Split(raw, "\n")
	for _, line := range rawLines {
		if len(line) == 0 {
			continue
		}

		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			continue
		}

		key := line[:colonIdx]
		lowerKey := strings.ToLower(key)
		if lowerKey == strings.ToLower(languageKey) {
			languageKey = key
		} else if lowerKey == strings.ToLower(pluralFormsKey) {
			pluralFormsKey = key
		}

		value := strings.TrimSpace(line[colonIdx+1:])
		do.Headers.Add(key, value)
	}

	// Get/save needed headers
	do.Language = do.Headers.Get(languageKey)
	do.PluralForms = do.Headers.Get(pluralFormsKey)

	// Parse Plural-Forms formula
	if do.PluralForms == "" {
		return
	}

	// Split plural form header value
	pfs := strings.Split(do.PluralForms, ";")

	// Parse values
	for _, i := range pfs {
		vs := strings.SplitN(i, "=", 2)
		if len(vs) != 2 {
			continue
		}

		switch strings.TrimSpace(vs[0]) {
		case "nplurals":
			do.nplurals, _ = strconv.Atoi(vs[1])

		case "plural":
			do.plural = vs[1]

			if expr, err := plurals.Compile(do.plural); err == nil {
				do.pluralforms = expr
			}

		}
	}
}

// Drops any translations stored that have not been Set*() since 'po'
// was initialised
func (do *Domain) DropStaleTranslations() {
	do.trMutex.Lock()
	do.pluralMutex.Lock()
	defer do.trMutex.Unlock()
	defer do.pluralMutex.Unlock()

	for name, ctx := range do.contextTranslations {
		for id, trans := range ctx {
			if trans.IsStale() {
				delete(ctx, id)
			}
		}
		if len(ctx) == 0 {
			delete(do.contextTranslations, name)
		}
	}

	for id, trans := range do.translations {
		if trans.IsStale() {
			delete(do.translations, id)
		}
	}
}

// Set source references for a given translation
func (do *Domain) SetRefs(str string, refs []string) {
	do.trMutex.Lock()
	do.pluralMutex.Lock()
	defer do.trMutex.Unlock()
	defer do.pluralMutex.Unlock()

	if trans, ok := do.translations[str]; ok {
		trans.Refs = refs
	} else {
		trans = NewTranslation()
		trans.ID = str
		trans.SetRefs(refs)
		do.translations[str] = trans
	}
}

// Get source references for a given translation
func (do *Domain) GetRefs(str string) []string {
	// Sync read
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations != nil {
		if trans, ok := do.translations[str]; ok {
			return trans.Refs
		}
	}
	return nil
}

// Set the translation of a given string
func (do *Domain) Set(id, str string) {
	do.trMutex.Lock()
	do.pluralMutex.Lock()
	defer do.trMutex.Unlock()
	defer do.pluralMutex.Unlock()

	if trans, ok := do.translations[id]; ok {
		trans.Set(str)
	} else {
		trans = NewTranslation()
		trans.ID = id
		trans.Set(str)
		do.translations[id] = trans
	}
}

func (do *Domain) Get(str string, vars ...interface{}) string {
	// Sync read
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations != nil {
		if _, ok := do.translations[str]; ok {
			return Printf(do.translations[str].Get(), vars...)
		}
	}

	// Return the same we received by default
	return Printf(str, vars...)
}

func (do *Domain) Append(b []byte, str string, vars ...interface{}) []byte {
	// Sync read
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations != nil {
		if _, ok := do.translations[str]; ok {
			return Appendf(b, do.translations[str].Get(), vars...)
		}
	}

	// Return the same we received by default
	return Appendf(b, str, vars...)
}

// Set the (N)th plural form for the given string
func (do *Domain) SetN(id, plural string, n int, str string) {
	// Get plural form _before_ lock down
	pluralForm := do.pluralForm(n)

	do.trMutex.Lock()
	do.pluralMutex.Lock()
	defer do.trMutex.Unlock()
	defer do.pluralMutex.Unlock()

	if trans, ok := do.translations[id]; ok {
		trans.SetN(pluralForm, str)
	} else {
		trans = NewTranslation()
		trans.ID = id
		trans.PluralID = plural
		trans.SetN(pluralForm, str)
		do.translations[id] = trans
	}
}

// GetN retrieves the (N)th plural form of Translation for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) GetN(str, plural string, n int, vars ...interface{}) string {
	// Sync read
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations != nil {
		if _, ok := do.translations[str]; ok {
			return Printf(do.translations[str].GetN(do.pluralForm(n)), vars...)
		}
	}

	// Parse plural forms to distinguish between plural and singular
	if do.pluralForm(n) == 0 {
		return Printf(str, vars...)
	}
	return Printf(plural, vars...)
}

// GetN retrieves the (N)th plural form of Translation for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) AppendN(b []byte, str, plural string, n int, vars ...interface{}) []byte {
	// Sync read
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations != nil {
		if _, ok := do.translations[str]; ok {
			return Appendf(b, do.translations[str].GetN(do.pluralForm(n)), vars...)
		}
	}

	// Parse plural forms to distinguish between plural and singular
	if do.pluralForm(n) == 0 {
		return Appendf(b, str, vars...)
	}
	return Appendf(b, plural, vars...)
}

// Set the translation for the given string in the given context
func (do *Domain) SetC(id, ctx, str string) {
	do.trMutex.Lock()
	do.pluralMutex.Lock()
	defer do.trMutex.Unlock()
	defer do.pluralMutex.Unlock()

	if context, ok := do.contextTranslations[ctx]; ok {
		if trans, hasTrans := context[id]; hasTrans {
			trans.Set(str)
		} else {
			trans = NewTranslation()
			trans.ID = id
			trans.Set(str)
			context[id] = trans
		}
	} else {
		trans := NewTranslation()
		trans.ID = id
		trans.Set(str)
		do.contextTranslations[ctx] = map[string]*Translation{
			id: trans,
		}
	}
}

// GetC retrieves the corresponding Translation for a given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) GetC(str, ctx string, vars ...interface{}) string {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.contextTranslations != nil {
		if _, ok := do.contextTranslations[ctx]; ok {
			if do.contextTranslations[ctx] != nil {
				if _, ok := do.contextTranslations[ctx][str]; ok {
					return Printf(do.contextTranslations[ctx][str].Get(), vars...)
				}
			}
		}
	}

	// Return the string we received by default
	return Printf(str, vars...)
}

// AppendC retrieves the corresponding Translation for a given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) AppendC(b []byte, str, ctx string, vars ...interface{}) []byte {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.contextTranslations != nil {
		if _, ok := do.contextTranslations[ctx]; ok {
			if do.contextTranslations[ctx] != nil {
				if _, ok := do.contextTranslations[ctx][str]; ok {
					return Appendf(b, do.contextTranslations[ctx][str].Get(), vars...)
				}
			}
		}
	}

	// Return the string we received by default
	return Appendf(b, str, vars...)
}

// Set the (N)th plural form for the given string in the given context
func (do *Domain) SetNC(id, plural, ctx string, n int, str string) {
	// Get plural form _before_ lock down
	pluralForm := do.pluralForm(n)

	do.trMutex.Lock()
	do.pluralMutex.Lock()
	defer do.trMutex.Unlock()
	defer do.pluralMutex.Unlock()

	if context, ok := do.contextTranslations[ctx]; ok {
		if trans, hasTrans := context[id]; hasTrans {
			trans.SetN(pluralForm, str)
		} else {
			trans = NewTranslation()
			trans.ID = id
			trans.SetN(pluralForm, str)
			context[id] = trans
		}
	} else {
		trans := NewTranslation()
		trans.ID = id
		trans.SetN(pluralForm, str)
		do.contextTranslations[ctx] = map[string]*Translation{
			id: trans,
		}
	}
}

// GetNC retrieves the (N)th plural form of Translation for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) GetNC(str, plural string, n int, ctx string, vars ...interface{}) string {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.contextTranslations != nil {
		if _, ok := do.contextTranslations[ctx]; ok {
			if do.contextTranslations[ctx] != nil {
				if _, ok := do.contextTranslations[ctx][str]; ok {
					return Printf(do.contextTranslations[ctx][str].GetN(do.pluralForm(n)), vars...)
				}
			}
		}
	}

	if n == 1 {
		return Printf(str, vars...)
	}
	return Printf(plural, vars...)
}

// AppendNC retrieves the (N)th plural form of Translation for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) AppendNC(b []byte, str, plural string, n int, ctx string, vars ...interface{}) []byte {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.contextTranslations != nil {
		if _, ok := do.contextTranslations[ctx]; ok {
			if do.contextTranslations[ctx] != nil {
				if _, ok := do.contextTranslations[ctx][str]; ok {
					return Appendf(b, do.contextTranslations[ctx][str].GetN(do.pluralForm(n)), vars...)
				}
			}
		}
	}

	if n == 1 {
		return Appendf(b, str, vars...)
	}
	return Appendf(b, plural, vars...)
}

// IsTranslated reports whether a string is translated
func (do *Domain) IsTranslated(str string) bool {
	return do.IsTranslatedN(str, 1)
}

// IsTranslatedN reports whether a plural string is translated
func (do *Domain) IsTranslatedN(str string, n int) bool {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations == nil {
		return false
	}
	tr, ok := do.translations[str]
	if !ok {
		return false
	}
	return tr.IsTranslatedN(do.pluralForm(n))
}

// IsTranslatedC reports whether a context string is translated
func (do *Domain) IsTranslatedC(str, ctx string) bool {
	return do.IsTranslatedNC(str, 1, ctx)
}

// IsTranslatedNC reports whether a plural context string is translated
func (do *Domain) IsTranslatedNC(str string, n int, ctx string) bool {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.contextTranslations == nil {
		return false
	}
	translations, ok := do.contextTranslations[ctx]
	if !ok {
		return false
	}
	tr, ok := translations[str]
	if !ok {
		return false
	}
	return tr.IsTranslatedN(do.pluralForm(n))
}

// GetTranslations returns a copy of every translation in the domain. It does not support contexts.
func (do *Domain) GetTranslations() map[string]*Translation {
	all := make(map[string]*Translation, len(do.translations))

	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	for msgID, trans := range do.translations {
		newTrans := NewTranslation()
		newTrans.ID = trans.ID
		newTrans.PluralID = trans.PluralID
		newTrans.dirty = trans.dirty
		if len(trans.Refs) > 0 {
			newTrans.Refs = make([]string, len(trans.Refs))
			copy(newTrans.Refs, trans.Refs)
		}
		for k, v := range trans.Trs {
			newTrans.Trs[k] = v
		}
		all[msgID] = newTrans
	}

	return all
}

type SourceReference struct {
	path    string
	line    int
	context string
	trans   *Translation
}

func extractPathAndLine(ref string) (string, int) {
	var path string
	var line int
	colonIdx := strings.IndexRune(ref, ':')
	if colonIdx >= 0 {
		path = ref[:colonIdx]
		line, _ = strconv.Atoi(ref[colonIdx+1:])
	} else {
		path = ref
		line = 0
	}
	return path, line
}

// MarshalText implements encoding.TextMarshaler interface
// Assists round-trip of POT/PO content
func (do *Domain) MarshalText() ([]byte, error) {
	var buf bytes.Buffer
	if len(do.headerComments) > 0 {
		buf.WriteString(strings.Join(do.headerComments, "\n"))
		buf.WriteByte(byte('\n'))
	}
	buf.WriteString("msgid \"\"\nmsgstr \"\"")

	// Standard order consistent with xgettext
	headerOrder := map[string]int{
		"project-id-version":        0,
		"report-msgid-bugs-to":      1,
		"pot-creation-date":         2,
		"po-revision-date":          3,
		"last-translator":           4,
		"language-team":             5,
		"language":                  6,
		"mime-version":              7,
		"content-type":              9,
		"content-transfer-encoding": 10,
		"plural-forms":              11,
	}

	headerKeys := make([]string, 0, len(do.Headers))

	for k, _ := range do.Headers {
		headerKeys = append(headerKeys, k)
	}

	sort.Slice(headerKeys, func(i, j int) bool {
		var iOrder int
		var jOrder int
		var ok bool
		if iOrder, ok = headerOrder[strings.ToLower(headerKeys[i])]; !ok {
			iOrder = 8
		}

		if jOrder, ok = headerOrder[strings.ToLower(headerKeys[j])]; !ok {
			jOrder = 8
		}

		if iOrder < jOrder {
			return true
		}
		if iOrder > jOrder {
			return false
		}
		return headerKeys[i] < headerKeys[j]
	})

	for _, k := range headerKeys {
		// Access Headers map directly so as not to canonicalise
		v := do.Headers[k]

		for _, value := range v {
			buf.WriteString("\n\"" + k + ": " + value + "\\n\"")
		}
	}

	// Just as with headers, output translations in consistent order (to minimise diffs between round-trips), with (first) source reference taking priority, followed by context and finally ID
	references := make([]SourceReference, 0)
	for name, ctx := range do.contextTranslations {
		for id, trans := range ctx {
			if id == "" {
				continue
			}
			if len(trans.Refs) > 0 {
				path, line := extractPathAndLine(trans.Refs[0])
				references = append(references, SourceReference{
					path,
					line,
					name,
					trans,
				})
			} else {
				references = append(references, SourceReference{
					"",
					0,
					name,
					trans,
				})
			}
		}
	}

	for id, trans := range do.translations {
		if id == "" {
			continue
		}

		if len(trans.Refs) > 0 {
			path, line := extractPathAndLine(trans.Refs[0])
			references = append(references, SourceReference{
				path,
				line,
				"",
				trans,
			})
		} else {
			references = append(references, SourceReference{
				"",
				0,
				"",
				trans,
			})
		}
	}

	sort.Slice(references, func(i, j int) bool {
		if references[i].path < references[j].path {
			return true
		}
		if references[i].path > references[j].path {
			return false
		}
		if references[i].line < references[j].line {
			return true
		}
		if references[i].line > references[j].line {
			return false
		}

		if references[i].context < references[j].context {
			return true
		}
		if references[i].context > references[j].context {
			return false
		}
		return references[i].trans.ID < references[j].trans.ID
	})

	for _, ref := range references {
		trans := ref.trans
		if len(trans.Refs) > 0 {
			buf.WriteString("\n\n#: " + strings.Join(trans.Refs, " "))
		} else {
			buf.WriteByte(byte('\n'))
		}

		if ref.context == "" {
			buf.WriteString("\nmsgid \"" + EscapeSpecialCharacters(trans.ID) + "\"")
		} else {
			buf.WriteString("\nmsgctxt \"" + EscapeSpecialCharacters(ref.context) + "\"\nmsgid \"" + EscapeSpecialCharacters(trans.ID) + "\"")
		}

		if trans.PluralID == "" {
			buf.WriteString("\nmsgstr \"" + EscapeSpecialCharacters(trans.Trs[0]) + "\"")
		} else {
			buf.WriteString("\nmsgid_plural \"" + trans.PluralID + "\"")
			for i, tr := range trans.Trs {
				buf.WriteString("\nmsgstr[" + EscapeSpecialCharacters(strconv.Itoa(i)) + "] \"" + tr + "\"")
			}
		}
	}

	return buf.Bytes(), nil
}

func EscapeSpecialCharacters(s string) string {
	s = regexp.MustCompile(`([^\\])(")`).ReplaceAllString(s, "$1\\\"") // Escape non-escaped double quotation marks

	if strings.Count(s, "\n") == 0 {
		return s
	}

	// Handle EOL and multi-lines
	// Only one line, but finishing with \n
	if strings.Count(s, "\n") == 1 && strings.HasSuffix(s, "\n") {
		return strings.ReplaceAll(s, "\n", "\\n")
	}

	elems := strings.Split(s, "\n")
	// Skip last element for multiline which is an empty
	var shouldEndWithEOL bool
	if elems[len(elems)-1] == "" {
		elems = elems[:len(elems)-1]
		shouldEndWithEOL = true
	}
	data := []string{(`"`)}
	for i, v := range elems {
		l := fmt.Sprintf(`"%s\n"`, v)
		// Last element without EOL
		if i == len(elems)-1 && !shouldEndWithEOL {
			l = fmt.Sprintf(`"%s"`, v)
		}
		// Remove finale " to last element as the whole string will be quoted
		if i == len(elems)-1 {
			l = strings.TrimSuffix(l, `"`)
		}
		data = append(data, l)
	}
	return strings.Join(data, "\n")
}

// MarshalBinary implements encoding.BinaryMarshaler interface
func (do *Domain) MarshalBinary() ([]byte, error) {
	obj := new(TranslatorEncoding)
	obj.Headers = do.Headers
	obj.Language = do.Language
	obj.PluralForms = do.PluralForms
	obj.Nplurals = do.nplurals
	obj.Plural = do.plural
	obj.Translations = do.translations
	obj.Contexts = do.contextTranslations

	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(obj)

	return buff.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (do *Domain) UnmarshalBinary(data []byte) error {
	buff := bytes.NewBuffer(data)
	obj := new(TranslatorEncoding)

	decoder := gob.NewDecoder(buff)
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}

	do.Headers = obj.Headers
	do.Language = obj.Language
	do.PluralForms = obj.PluralForms
	do.nplurals = obj.Nplurals
	do.plural = obj.Plural
	do.translations = obj.Translations
	do.contextTranslations = obj.Contexts

	if expr, err := plurals.Compile(do.plural); err == nil {
		do.pluralforms = expr
	}

	return nil
}
