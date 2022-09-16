package query

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

type kv[K any, V any] struct {
	k K
	v V
}

// ios AST https://github.com/owncloud/ios-app/pull/933
func TestLexer(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		input string
		exp   []kv[Token, string]
	}{
		{
			input: "engineering",
			exp: []kv[Token, string]{
				{k: TValue, v: "engineering"},
				{k: TEof},
			},
		},
		{
			input: "engineering demos",
			exp: []kv[Token, string]{
				{k: TValue, v: "engineering"},
				{k: TValue, v: "demos"},
				{k: TEof},
			},
		},
		{
			input: "\"engineering demos\"",
			exp: []kv[Token, string]{
				{k: TQuotationMark},
				{k: TValue, v: "engineering"},
				{k: TValue, v: "demos"},
				{k: TQuotationMark},
				{k: TEof},
			},
		},
		{
			input: "\"engineering \"demos",
			exp: []kv[Token, string]{
				{k: TQuotationMark},
				{k: TValue, v: "engineering"},
				{k: TQuotationMark},
				{k: TValue, v: "demos"},
				{k: TEof},
			},
		},
		{
			input: "type:pdf",
			exp: []kv[Token, string]{
				{k: TField, v: "type"},
				{k: TValue, v: "pdf"},
				{k: TEof},
			},
		},
		{
			input: "before:2021",
			exp: []kv[Token, string]{
				{k: TField, v: "before"},
				{k: TValue, v: "2021"},
				{k: TEof},
			},
		},
		{
			input: "before:2021-02",
			exp: []kv[Token, string]{
				{k: TField, v: "before"},
				{k: TValue, v: "2021-02"},
				{k: TEof},
			},
		},
		{
			input: "before:2021-02-03",
			exp: []kv[Token, string]{
				{k: TField, v: "before"},
				{k: TValue, v: "2021-02-03"},
				{k: TEof},
			},
		},
		{
			input: "after:2020",
			exp: []kv[Token, string]{
				{k: TField, v: "after"},
				{k: TValue, v: "2020"},
				{k: TEof},
			},
		},
		{
			input: "after:2020-02",
			exp: []kv[Token, string]{
				{k: TField, v: "after"},
				{k: TValue, v: "2020-02"},
				{k: TEof},
			},
		},
		{
			input: "after:2020-02-03",
			exp: []kv[Token, string]{
				{k: TField, v: "after"},
				{k: TValue, v: "2020-02-03"},
				{k: TEof},
			},
		},
		{
			input: "on:2020-02-03",
			exp: []kv[Token, string]{
				{k: TField, v: "on"},
				{k: TValue, v: "2020-02-03"},
				{k: TEof},
			},
		},
		{
			input: "on:2020-02-03,2020-02-05",
			exp: []kv[Token, string]{
				{k: TField, v: "on"},
				{k: TValue, v: "2020-02-03,2020-02-05"},
				{k: TEof},
			},
		},
		{
			input: "smaller:200mb",
			exp: []kv[Token, string]{
				{k: TField, v: "smaller"},
				{k: TValue, v: "200mb"},
				{k: TEof},
			},
		},
		{
			input: "greater:1gb",
			exp: []kv[Token, string]{
				{k: TField, v: "greater"},
				{k: TValue, v: "1gb"},
				{k: TEof},
			},
		},
		{
			input: "owner:cscherm",
			exp: []kv[Token, string]{
				{k: TField, v: "owner"},
				{k: TValue, v: "cscherm"},
				{k: TEof},
			},
		},
		{
			input: ":file",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "file"},
				{k: TEof},
			},
		},
		{
			input: ":folder",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "folder"},
				{k: TEof},
			},
		},
		{
			input: ":image",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "image"},
				{k: TEof},
			},
		},
		{
			input: ":video",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "video"},
				{k: TEof},
			},
		},
		{
			input: ":year",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "year"},
				{k: TEof},
			},
		},
		{
			input: ":month",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "month"},
				{k: TEof},
			},
		},
		{
			input: ":week",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "week"},
				{k: TEof},
			},
		},
		{
			input: ":today",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "today"},
				{k: TEof},
			},
		},
		{
			input: ":d",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "d"},
				{k: TEof},
			},
		},
		{
			input: ":5d",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "5d"},
				{k: TEof},
			},
		},
		{
			input: ":2w",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "2w"},
				{k: TEof},
			},
		},
		{
			input: ":2w",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "2w"},
				{k: TEof},
			},
		},
		{
			input: ":1m",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "1m"},
				{k: TEof},
			},
		},
		{
			input: ":m",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "m"},
				{k: TEof},
			},
		},
		{
			input: ":1y",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "1y"},
				{k: TEof},
			},
		},
		{
			input: ":y",
			exp: []kv[Token, string]{
				{k: TShortcut, v: "y"},
				{k: TEof},
			},
		},
		{
			input: "-:image",
			exp: []kv[Token, string]{
				{k: TNegation},
				{k: TShortcut, v: "image"},
				{k: TEof},
			},
		},
		{
			input: "\"engineering demos\" \"engineering \"demos type:pdf,mov on:2020-02-03,2020-02-05 :image :pdf -:image",
			exp: []kv[Token, string]{
				{k: TQuotationMark},
				{k: TValue, v: "engineering"},
				{k: TValue, v: "demos"},
				{k: TQuotationMark},
				{k: TQuotationMark},
				{k: TValue, v: "engineering"},
				{k: TQuotationMark},
				{k: TValue, v: "demos"},
				{k: TField, v: "type"},
				{k: TValue, v: "pdf,mov"},
				{k: TField, v: "on"},
				{k: TValue, v: "2020-02-03,2020-02-05"},
				{k: TShortcut, v: "image"},
				{k: TShortcut, v: "pdf"},
				{k: TNegation},
				{k: TShortcut, v: "image"},
				{k: TEof},
			},
		},
		{
			input: ": a b: c :::d:::\"e \"f\"\"\"a\"",
			exp: []kv[Token, string]{
				{k: TUnknown, v: ":"},
				{k: TValue, v: "a"},
				{k: TField, v: "b"},
				{k: TValue, v: "c"},
				{k: TUnknown, v: ":"},
				{k: TUnknown, v: ":"},
				{k: TField, v: "d"},
				{k: TUnknown, v: ":"},
				{k: TUnknown, v: ":"},
				{k: TQuotationMark},
				{k: TValue, v: "e"},
				{k: TQuotationMark},
				{k: TValue, v: "f"},
				{k: TQuotationMark},
				{k: TQuotationMark},
				{k: TQuotationMark},
				{k: TValue, v: "a"},
				{k: TQuotationMark},
				{k: TEof},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			l := NewLexer(strings.NewReader(c.input))
			for _, exp := range c.exp {
				tok, lit := l.Scan()
				g.Expect(tok).To(Equal(exp.k))
				g.Expect(lit).To(Equal(exp.v))
			}
		})
	}
}
