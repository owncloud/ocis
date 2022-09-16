package query

// Token maps to token type
type Token int

const (
	// TEof end of file token type
	TEof Token = iota
	// TUnknown unknown token type
	TUnknown
	// TNegation negation token type
	TNegation
	// TQuotationMark quotation-mark token type, e.g. not of type - "-"
	TQuotationMark
	// TField field token type, e.g. - '"'
	TField
	// TValue value token type, e.g. - "content:"
	TValue
	// TShortcut shortcut token type, e.g. all images - ":image"
	TShortcut
)

var eof = rune(0)

var tokens = map[Token]string{
	TEof:           "EOF",
	TUnknown:       "UNKNOWN",
	TNegation:      "NEGATION",
	TQuotationMark: "QUOTATION_MARK",
	TField:         "FIELD",
	TValue:         "VALUE",
	TShortcut:      "SHORTCUT",
}

// String returns the Token human-readable.
func (t Token) String() string {
	return tokens[t]
}
