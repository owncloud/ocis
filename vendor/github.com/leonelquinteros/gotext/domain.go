package gotext

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"maps"
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

// HeaderMap preserves MIMEHeader behaviour, without the canonicalisation
type HeaderMap map[string][]string

// Add key/value pair to HeaderMap
func (m HeaderMap) Add(key, value string) {
	m[key] = append(m[key], value)
}

// Del key from HeaderMap
func (m HeaderMap) Del(key string) {
	delete(m, key)
}

// Get value for key from HeaderMap
func (m HeaderMap) Get(key string) string {
	v := m[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Set key/value pair in HeaderMap
func (m HeaderMap) Set(key, value string) {
	m[key] = []string{value}
}

// Values returns all values for a given key from HeaderMap
func (m HeaderMap) Values(key string) []string {
	return m[key]
}

// NewDomain creates a new Domain instance
func NewDomain() *Domain {
	domain := new(Domain)

	domain.Headers = make(HeaderMap)
	domain.headerComments = make([]string, 0)
	domain.translations = make(map[string]*Translation)
	domain.contextTranslations = make(map[string]map[string]*Translation)
	domain.pluralTranslations = make(map[string]*Translation)

	return domain
}

func cloneRefs(refs []string) []string {
	if refs == nil {
		return nil
	}

	owned := make([]string, len(refs))
	copy(owned, refs)
	return owned
}

func cloneHeaderMap(headers HeaderMap) HeaderMap {
	if headers == nil {
		return nil
	}

	owned := make(HeaderMap, len(headers))
	for key, values := range headers {
		owned[key] = cloneRefs(values)
	}
	return owned
}

func cloneTranslation(trans *Translation) *Translation {
	if trans == nil {
		return nil
	}

	owned := NewTranslation()
	owned.ID = trans.ID
	owned.PluralID = trans.PluralID
	owned.dirty = trans.dirty
	owned.Refs = cloneRefs(trans.Refs)

	maps.Copy(owned.Trs, trans.Trs)

	return owned
}

func cloneTranslations(translations map[string]*Translation) map[string]*Translation {
	if translations == nil {
		return nil
	}

	owned := make(map[string]*Translation, len(translations))
	for id, trans := range translations {
		owned[id] = cloneTranslation(trans)
	}
	return owned
}

func cloneContextTranslations(contexts map[string]map[string]*Translation) map[string]map[string]*Translation {
	if contexts == nil {
		return nil
	}

	owned := make(map[string]map[string]*Translation, len(contexts))
	for context, translations := range contexts {
		owned[context] = cloneTranslations(translations)
	}
	return owned
}

func (do *Domain) hasTranslation(id string) bool {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	return do.translations[id] != nil
}

func (do *Domain) hasContextTranslation(id, context string) bool {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	return do.contextTranslations[context][id] != nil
}

// SetPluralResolver sets a custom plural resolver function
func (do *Domain) SetPluralResolver(f func(int) int) {
	do.pluralMutex.Lock()
	defer do.pluralMutex.Unlock()

	do.customPluralResolver = f
}

func (do *Domain) pluralForm(n int) int {
	// Snapshot plural state before calling an expression or custom resolver.
	do.pluralMutex.RLock()
	pluralforms := do.pluralforms
	customPluralResolver := do.customPluralResolver
	do.pluralMutex.RUnlock()

	if pluralforms != nil {
		return pluralforms.Eval(uint32(n))
	}

	// Failure fallback
	if customPluralResolver != nil {
		return customPluralResolver(n)
	}

	/* Use the Germanic plural rule.  */
	if n == 1 {
		return 0
	}
	return 1
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

	for line := range strings.SplitSeq(raw, "\n") {
		if len(line) == 0 {
			continue
		}

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		lowerKey := strings.ToLower(key)
		if lowerKey == strings.ToLower(languageKey) {
			languageKey = key
		} else if lowerKey == strings.ToLower(pluralFormsKey) {
			pluralFormsKey = key
		}

		value = strings.TrimSpace(value)
		do.Headers.Add(key, value)
	}

	// Get/save needed headers
	do.Language = do.Headers.Get(languageKey)
	do.PluralForms = do.Headers.Get(pluralFormsKey)

	// Parse Plural-Forms formula
	if do.PluralForms == "" {
		return
	}

	for i := range strings.SplitSeq(do.PluralForms, ";") {
		key, value, ok := strings.Cut(i, "=")
		if !ok {
			continue
		}

		switch strings.TrimSpace(key) {
		case "nplurals":
			do.nplurals, _ = strconv.Atoi(value)

		case "plural":
			do.plural = value

			if expr, err := plurals.Compile(do.plural); err == nil {
				do.pluralforms = expr
			}

		}
	}
}

// DropStaleTranslations drops any translations stored that have not been Set*()
// since 'po' was initialised
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

// SetRefs set source references for a given translation
func (do *Domain) SetRefs(str string, refs []string) {
	ownedRefs := cloneRefs(refs)

	do.trMutex.Lock()
	do.pluralMutex.Lock()
	defer do.trMutex.Unlock()
	defer do.pluralMutex.Unlock()

	if trans, ok := do.translations[str]; ok && trans != nil {
		trans.SetRefs(ownedRefs)
	} else {
		trans := NewTranslation()
		trans.ID = str
		trans.SetRefs(ownedRefs)
		do.translations[str] = trans
	}
}

// GetRefs get source references for a given translation
func (do *Domain) GetRefs(str string) []string {
	// Sync read
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations != nil {
		if trans, ok := do.translations[str]; ok && trans != nil {
			return cloneRefs(trans.Refs)
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

// Get retrieves the Translation for the given string.
func (do *Domain) Get(str string, vars ...any) string {
	// Sync read
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations != nil {
		if _, ok := do.translations[str]; ok {
			return FormatString(do.translations[str].Get(), vars...)
		}
	}

	// Return the same we received by default
	return FormatString(str, vars...)
}

// Append retrieves the Translation for the given string.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) Append(b []byte, str string, vars ...any) []byte {
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

// SetN sets the (N)th plural form for the given string
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
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) GetN(str, plural string, n int, vars ...any) string {
	// Get plural form before acquiring the translation lock so a custom
	// resolver is never called while a project lock is held.
	pluralForm := do.pluralForm(n)

	// Sync read
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations != nil {
		if _, ok := do.translations[str]; ok {
			return FormatString(do.translations[str].GetN(pluralForm), vars...)
		}
	}

	// Parse plural forms to distinguish between plural and singular
	if pluralForm == 0 {
		return FormatString(str, vars...)
	}
	return FormatString(plural, vars...)
}

// AppendN adds the (N)th plural form of Translation for the given string.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) AppendN(b []byte, str, plural string, n int, vars ...any) []byte {
	// Get plural form before acquiring the translation lock so a custom
	// resolver is never called while a project lock is held.
	pluralForm := do.pluralForm(n)

	// Sync read
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations != nil {
		if _, ok := do.translations[str]; ok {
			return Appendf(b, do.translations[str].GetN(pluralForm), vars...)
		}
	}

	// Parse plural forms to distinguish between plural and singular
	if pluralForm == 0 {
		return Appendf(b, str, vars...)
	}
	return Appendf(b, plural, vars...)
}

// SetC sets the translation for the given string in the given context
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
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) GetC(str, ctx string, vars ...any) string {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.contextTranslations != nil {
		if _, ok := do.contextTranslations[ctx]; ok {
			if do.contextTranslations[ctx] != nil {
				if _, ok := do.contextTranslations[ctx][str]; ok {
					return FormatString(do.contextTranslations[ctx][str].Get(), vars...)
				}
			}
		}
	}

	// Return the string we received by default
	return FormatString(str, vars...)
}

// AppendC retrieves the corresponding Translation for a given string in the given context.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) AppendC(b []byte, str, ctx string, vars ...any) []byte {
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

// SetNC sets the (N)th plural form for the given string in the given context
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
			trans.PluralID = plural
			trans.SetN(pluralForm, str)
			context[id] = trans
		}
	} else {
		trans := NewTranslation()
		trans.ID = id
		trans.PluralID = plural
		trans.SetN(pluralForm, str)
		do.contextTranslations[ctx] = map[string]*Translation{
			id: trans,
		}
	}
}

// GetNC retrieves the (N)th plural form of Translation for the given string in the given context.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) GetNC(str, plural string, n int, ctx string, vars ...any) string {
	if do.hasContextTranslation(str, ctx) {
		// Get plural form before acquiring the translation lock so a custom
		// resolver is never called while a project lock is held.
		pluralForm := do.pluralForm(n)

		do.trMutex.RLock()
		if translations, ok := do.contextTranslations[ctx]; ok && translations != nil {
			if trans, ok := translations[str]; ok && trans != nil {
				translated := trans.GetN(pluralForm)
				do.trMutex.RUnlock()
				return FormatString(translated, vars...)
			}
		}
		do.trMutex.RUnlock()
	}

	if n == 1 {
		return FormatString(str, vars...)
	}
	return FormatString(plural, vars...)
}

// AppendNC retrieves the (N)th plural form of Translation for the given string in the given context.
// Supports optional parameters (vars... any) to be inserted on the formatted string using the fmt.Printf syntax.
func (do *Domain) AppendNC(b []byte, str, plural string, n int, ctx string, vars ...any) []byte {
	if do.hasContextTranslation(str, ctx) {
		// Get plural form before acquiring the translation lock so a custom
		// resolver is never called while a project lock is held.
		pluralForm := do.pluralForm(n)

		do.trMutex.RLock()
		if translations, ok := do.contextTranslations[ctx]; ok && translations != nil {
			if trans, ok := translations[str]; ok && trans != nil {
				translated := trans.GetN(pluralForm)
				do.trMutex.RUnlock()
				return Appendf(b, translated, vars...)
			}
		}
		do.trMutex.RUnlock()
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
	if !do.hasTranslation(str) {
		return false
	}

	// Get plural form before acquiring the translation lock so a custom
	// resolver is never called while a project lock is held.
	pluralForm := do.pluralForm(n)

	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.translations == nil {
		return false
	}
	tr, ok := do.translations[str]
	if !ok || tr == nil {
		return false
	}
	return tr.IsTranslatedN(pluralForm)
}

// IsTranslatedC reports whether a context string is translated
func (do *Domain) IsTranslatedC(str, ctx string) bool {
	return do.IsTranslatedNC(str, 1, ctx)
}

// IsTranslatedNC reports whether a plural context string is translated
func (do *Domain) IsTranslatedNC(str string, n int, ctx string) bool {
	if !do.hasContextTranslation(str, ctx) {
		return false
	}

	// Get plural form before acquiring the translation lock so a custom
	// resolver is never called while a project lock is held.
	pluralForm := do.pluralForm(n)

	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	if do.contextTranslations == nil {
		return false
	}
	translations, ok := do.contextTranslations[ctx]
	if !ok || translations == nil {
		return false
	}
	tr, ok := translations[str]
	if !ok || tr == nil {
		return false
	}
	return tr.IsTranslatedN(pluralForm)
}

// GetTranslations returns a copy of every translation in the domain. It does not support contexts.
func (do *Domain) GetTranslations() map[string]*Translation {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	all := make(map[string]*Translation, len(do.translations))
	for msgID, trans := range do.translations {
		all[msgID] = cloneTranslation(trans)
	}

	return all
}

// GetCtxTranslations returns a copy of every translation in the domain with context
func (do *Domain) GetCtxTranslations() map[string]map[string]*Translation {
	do.trMutex.RLock()
	defer do.trMutex.RUnlock()

	all := make(map[string]map[string]*Translation, len(do.contextTranslations))
	for ctx, translations := range do.contextTranslations {
		for msgID, trans := range translations {
			if all[ctx] == nil {
				all[ctx] = make(map[string]*Translation)
			}

			all[ctx][msgID] = cloneTranslation(trans)
		}
	}

	return all
}

// SourceReference is a struct to hold source reference information
type SourceReference struct {
	path       string
	line       int
	context    string
	hasContext bool
	trans      *Translation
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
	do.trMutex.RLock()
	do.pluralMutex.RLock()
	defer do.trMutex.RUnlock()
	defer do.pluralMutex.RUnlock()

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

	for k := range do.Headers {
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
			encoded := EscapeSpecialCharacters(k + ": " + value + "\n")
			buf.WriteString("\n\"" + encoded + "\"")
		}
	}

	// Just as with headers, output translations in consistent order (to minimise diffs between round-trips), with (first) source reference taking priority, followed by context and finally ID
	references := make([]SourceReference, 0)
	for name, ctx := range do.contextTranslations {
		for _, trans := range ctx {
			if len(trans.Refs) > 0 {
				path, line := extractPathAndLine(trans.Refs[0])
				references = append(references, SourceReference{
					path:       path,
					line:       line,
					context:    name,
					hasContext: true,
					trans:      trans,
				})
			} else {
				references = append(references, SourceReference{
					context:    name,
					hasContext: true,
					trans:      trans,
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
				path:  path,
				line:  line,
				trans: trans,
			})
		} else {
			references = append(references, SourceReference{
				trans: trans,
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
		if references[i].hasContext != references[j].hasContext {
			return !references[i].hasContext
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

		if !ref.hasContext {
			buf.WriteString("\nmsgid \"" + EscapeSpecialCharacters(trans.ID) + "\"")
		} else {
			buf.WriteString("\nmsgctxt \"" + EscapeSpecialCharacters(ref.context) + "\"\nmsgid \"" + EscapeSpecialCharacters(trans.ID) + "\"")
		}

		if trans.PluralID == "" {
			buf.WriteString("\nmsgstr \"" + EscapeSpecialCharacters(trans.Trs[0]) + "\"")
		} else {
			buf.WriteString("\nmsgid_plural \"" + EscapeSpecialCharacters(trans.PluralID) + "\"")

			pluralIndexes := make([]int, 0, len(trans.Trs))
			for i := range trans.Trs {
				pluralIndexes = append(pluralIndexes, i)
			}
			sort.Ints(pluralIndexes)

			for _, i := range pluralIndexes {
				buf.WriteString("\nmsgstr[" + strconv.Itoa(i) + "] \"" + EscapeSpecialCharacters(trans.Trs[i]) + "\"")
			}
		}
	}

	return buf.Bytes(), nil
}

// EscapeSpecialCharacters escapes special characters in a string
func EscapeSpecialCharacters(s string) string {
	var escaped strings.Builder
	escaped.Grow(len(s))

	for i := range len(s) {
		switch s[i] {
		case '\\':
			escaped.WriteByte('\\')
			escaped.WriteByte('\\')
		case '"':
			escaped.WriteByte('\\')
			escaped.WriteByte('"')
		default:
			escaped.WriteByte(s[i])
		}
	}

	s = escaped.String()
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

func (do *Domain) translatorEncodingSnapshot() *TranslatorEncoding {
	do.trMutex.RLock()
	do.pluralMutex.RLock()
	defer do.trMutex.RUnlock()
	defer do.pluralMutex.RUnlock()

	return &TranslatorEncoding{
		Headers:      cloneHeaderMap(do.Headers),
		Language:     do.Language,
		PluralForms:  do.PluralForms,
		Nplurals:     do.nplurals,
		Plural:       do.plural,
		Translations: cloneTranslations(do.translations),
		Contexts:     cloneContextTranslations(do.contextTranslations),
	}
}

// MarshalBinary implements encoding.BinaryMarshaler interface
func (do *Domain) MarshalBinary() ([]byte, error) {
	obj := do.translatorEncodingSnapshot()

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
	if err := decoder.Decode(obj); err != nil {
		return err
	}

	headers := cloneHeaderMap(obj.Headers)
	if headers == nil {
		headers = make(HeaderMap)
	}
	translations := cloneTranslations(obj.Translations)
	if translations == nil {
		translations = make(map[string]*Translation)
	}
	contextTranslations := cloneContextTranslations(obj.Contexts)
	if contextTranslations == nil {
		contextTranslations = make(map[string]map[string]*Translation)
	}

	for id, trans := range translations {
		if trans == nil {
			trans = NewTranslation()
			trans.ID = id
			translations[id] = trans
		}
	}
	for context, translationsForContext := range contextTranslations {
		if translationsForContext == nil {
			contextTranslations[context] = make(map[string]*Translation)
			continue
		}
		for id, trans := range translationsForContext {
			if trans == nil {
				trans = NewTranslation()
				trans.ID = id
				translationsForContext[id] = trans
			}
		}
	}

	var pluralforms plurals.Expression
	if expr, err := plurals.Compile(obj.Plural); err == nil {
		pluralforms = expr
	}

	do.trMutex.Lock()
	do.pluralMutex.Lock()

	do.Headers = headers
	do.Language = obj.Language
	do.PluralForms = obj.PluralForms
	do.nplurals = obj.Nplurals
	do.plural = obj.Plural
	do.pluralforms = pluralforms
	do.translations = translations
	do.contextTranslations = contextTranslations
	do.pluralTranslations = make(map[string]*Translation)
	do.headerComments = make([]string, 0)
	do.trBuffer = nil
	do.ctxBuffer = ""
	do.refBuffer = ""

	do.pluralMutex.Unlock()
	do.trMutex.Unlock()

	return nil
}
