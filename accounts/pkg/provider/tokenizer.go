package provider

import "github.com/CiscoM31/godata"

// FilterTokenizer creates a tokenizer capable of tokenizing filter statements
// TODO disable tokens we don't handle anyway
func FilterTokenizer() *godata.Tokenizer {
	t := godata.Tokenizer{}
	t.Add("^[0-9]{4,4}-[0-9]{2,2}-[0-9]{2,2}T[0-9]{2,2}:[0-9]{2,2}(:[0-9]{2,2}(.[0-9]+)?)?(Z|[+-][0-9]{2,2}:[0-9]{2,2})", godata.FilterTokenDateTime)
	t.Add("^-?[0-9]{4,4}-[0-9]{2,2}-[0-9]{2,2}", godata.FilterTokenDate)
	t.Add("^[0-9]{2,2}:[0-9]{2,2}(:[0-9]{2,2}(.[0-9]+)?)?", godata.FilterTokenTime)
	t.Add("^\\(", godata.FilterTokenOpenParen)
	t.Add("^\\)", godata.FilterTokenCloseParen)
	t.Add("^/", godata.FilterTokenNav)
	t.Add("^:", godata.FilterTokenColon)
	t.Add("^,", godata.FilterTokenComma)
	t.Add("^(geo.distance|geo.intersects|geo.length)", godata.FilterTokenFunc)
	t.Add("^(substringof|substring|length|indexof)", godata.FilterTokenFunc)
	// only change from the global tokenizer is the added ap
	t.Add("^(eq|ne|gt|ge|lt|le|and|or|not|has|in|ap)", godata.FilterTokenLogical)
	t.Add("^(add|sub|mul|divby|div|mod)", godata.FilterTokenOp)
	t.Add("^(contains|endswith|startswith|tolower|toupper|"+
		"trim|concat|year|month|day|hour|minute|second|fractionalseconds|date|"+
		"time|totaloffsetminutes|now|maxdatetime|mindatetime|totalseconds|round|"+
		"floor|ceiling|isof|cast)", godata.FilterTokenFunc)
	t.Add("^(any|all)", godata.FilterTokenLambda)
	t.Add("^null", godata.FilterTokenNull)
	t.Add("^\\$it", godata.FilterTokenIt)
	t.Add("^\\$root", godata.FilterTokenRoot)
	t.Add("^-?[0-9]+\\.[0-9]+", godata.FilterTokenFloat)
	t.Add("^-?[0-9]+", godata.FilterTokenInteger)
	t.Add("^'(''|[^'])*'", godata.FilterTokenString)
	t.Add("^(true|false)", godata.FilterTokenBoolean)
	t.Add("^@*[a-zA-Z][a-zA-Z0-9_.]*", godata.FilterTokenLiteral) // The optional '@' character is used to identify parameter aliases
	t.Ignore("^ ", godata.FilterTokenWhitespace)

	return &t
}
