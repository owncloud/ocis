package query

// Token maps to token type
type Token int

const (
	T_EOF Token = iota
	T_UNKNOWN
	T_NEGATION
	T_QUOTATION_MARK
	T_FIELD
	T_VALUE
	T_SHORTCUT
)

var eof = rune(0)

var tokens = map[Token]string{
	T_EOF:            "EOF",
	T_UNKNOWN:        "UNKNOWN",
	T_NEGATION:       "NEGATION",
	T_QUOTATION_MARK: "QUOTATION_MARK",
	T_FIELD:          "FIELD",
	T_VALUE:          "VALUE",
	T_SHORTCUT:       "SHORTCUT",
}

// String returns the Token human-readable.
func (t Token) String() string {
	return tokens[t]
}
