package query

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

// Lexer is responsible for lexing the query.
type Lexer struct {
	r *bufio.Reader
}

// NewLexer creates a new Lexer
func NewLexer(r io.Reader) Lexer {
	return Lexer{r: bufio.NewReader(r)}
}

// Scan reads a query section and returns a Token and literal
func (l *Lexer) Scan() (Token, string) {

	for {
		r := l.read()

		if r == eof {
			return T_EOF, ""
		}

		if unicode.IsSpace(r) {
			continue
		}

		if r == '"' {
			return T_QUOTATION_MARK, ""
		}

		if r == '-' {
			return T_NEGATION, ""
		}

		if r != ':' {
			l.unread()
			return l.scanUnknown(T_VALUE)
		}

		if r == ':' && (unicode.IsLetter(l.peek(1)) || unicode.IsNumber(l.peek(1))) {
			return l.scanUnknown(T_SHORTCUT)
		}

		return T_UNKNOWN, string(r)
	}
}

func (l *Lexer) scanUnknown(t Token) (Token, string) {
	var buf bytes.Buffer

	for {
		r := l.read()

		if r == eof || unicode.IsSpace(r) {
			break
		}

		if r == '"' {
			l.unread()
			break
		}

		if r == ':' {
			return T_FIELD, buf.String()
		}

		buf.WriteRune(r)
	}

	return t, buf.String()
}

func (l *Lexer) peek(n int) rune {
	b, _ := l.r.Peek(n)
	r, _ := utf8.DecodeRune(b)
	return r
}

func (l *Lexer) read() rune {
	r, _, err := l.r.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

func (l *Lexer) unread() {
	_ = l.r.UnreadRune()
}
