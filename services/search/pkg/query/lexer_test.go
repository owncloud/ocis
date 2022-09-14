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
				{k: T_VALUE, v: "engineering"},
				{k: T_EOF},
			},
		},
		{
			input: "engineering demos",
			exp: []kv[Token, string]{
				{k: T_VALUE, v: "engineering"},
				{k: T_VALUE, v: "demos"},
				{k: T_EOF},
			},
		},
		{
			input: "\"engineering demos\"",
			exp: []kv[Token, string]{
				{k: T_QUOTATION_MARK},
				{k: T_VALUE, v: "engineering"},
				{k: T_VALUE, v: "demos"},
				{k: T_QUOTATION_MARK},
				{k: T_EOF},
			},
		},
		{
			input: "\"engineering \"demos",
			exp: []kv[Token, string]{
				{k: T_QUOTATION_MARK},
				{k: T_VALUE, v: "engineering"},
				{k: T_QUOTATION_MARK},
				{k: T_VALUE, v: "demos"},
				{k: T_EOF},
			},
		},
		{
			input: "type:pdf",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "type"},
				{k: T_VALUE, v: "pdf"},
				{k: T_EOF},
			},
		},
		{
			input: "before:2021",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "before"},
				{k: T_VALUE, v: "2021"},
				{k: T_EOF},
			},
		},
		{
			input: "before:2021-02",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "before"},
				{k: T_VALUE, v: "2021-02"},
				{k: T_EOF},
			},
		},
		{
			input: "before:2021-02-03",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "before"},
				{k: T_VALUE, v: "2021-02-03"},
				{k: T_EOF},
			},
		},
		{
			input: "after:2020",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "after"},
				{k: T_VALUE, v: "2020"},
				{k: T_EOF},
			},
		},
		{
			input: "after:2020-02",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "after"},
				{k: T_VALUE, v: "2020-02"},
				{k: T_EOF},
			},
		},
		{
			input: "after:2020-02-03",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "after"},
				{k: T_VALUE, v: "2020-02-03"},
				{k: T_EOF},
			},
		},
		{
			input: "on:2020-02-03",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "on"},
				{k: T_VALUE, v: "2020-02-03"},
				{k: T_EOF},
			},
		},
		{
			input: "on:2020-02-03,2020-02-05",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "on"},
				{k: T_VALUE, v: "2020-02-03,2020-02-05"},
				{k: T_EOF},
			},
		},
		{
			input: "smaller:200mb",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "smaller"},
				{k: T_VALUE, v: "200mb"},
				{k: T_EOF},
			},
		},
		{
			input: "greater:1gb",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "greater"},
				{k: T_VALUE, v: "1gb"},
				{k: T_EOF},
			},
		},
		{
			input: "owner:cscherm",
			exp: []kv[Token, string]{
				{k: T_FIELD, v: "owner"},
				{k: T_VALUE, v: "cscherm"},
				{k: T_EOF},
			},
		},
		{
			input: ":file",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "file"},
				{k: T_EOF},
			},
		},
		{
			input: ":folder",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "folder"},
				{k: T_EOF},
			},
		},
		{
			input: ":image",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "image"},
				{k: T_EOF},
			},
		},
		{
			input: ":video",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "video"},
				{k: T_EOF},
			},
		},
		{
			input: ":year",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "year"},
				{k: T_EOF},
			},
		},
		{
			input: ":month",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "month"},
				{k: T_EOF},
			},
		},
		{
			input: ":week",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "week"},
				{k: T_EOF},
			},
		},
		{
			input: ":today",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "today"},
				{k: T_EOF},
			},
		},
		{
			input: ":d",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "d"},
				{k: T_EOF},
			},
		},
		{
			input: ":5d",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "5d"},
				{k: T_EOF},
			},
		},
		{
			input: ":2w",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "2w"},
				{k: T_EOF},
			},
		},
		{
			input: ":2w",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "2w"},
				{k: T_EOF},
			},
		},
		{
			input: ":1m",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "1m"},
				{k: T_EOF},
			},
		},
		{
			input: ":m",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "m"},
				{k: T_EOF},
			},
		},
		{
			input: ":1y",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "1y"},
				{k: T_EOF},
			},
		},
		{
			input: ":y",
			exp: []kv[Token, string]{
				{k: T_SHORTCUT, v: "y"},
				{k: T_EOF},
			},
		},
		{
			input: "-:image",
			exp: []kv[Token, string]{
				{k: T_NEGATION},
				{k: T_SHORTCUT, v: "image"},
				{k: T_EOF},
			},
		},
		{
			input: "\"engineering demos\" \"engineering \"demos type:pdf,mov on:2020-02-03,2020-02-05 :image :pdf -:image",
			exp: []kv[Token, string]{
				{k: T_QUOTATION_MARK},
				{k: T_VALUE, v: "engineering"},
				{k: T_VALUE, v: "demos"},
				{k: T_QUOTATION_MARK},
				{k: T_QUOTATION_MARK},
				{k: T_VALUE, v: "engineering"},
				{k: T_QUOTATION_MARK},
				{k: T_VALUE, v: "demos"},
				{k: T_FIELD, v: "type"},
				{k: T_VALUE, v: "pdf,mov"},
				{k: T_FIELD, v: "on"},
				{k: T_VALUE, v: "2020-02-03,2020-02-05"},
				{k: T_SHORTCUT, v: "image"},
				{k: T_SHORTCUT, v: "pdf"},
				{k: T_NEGATION},
				{k: T_SHORTCUT, v: "image"},
				{k: T_EOF},
			},
		},
		{
			input: ": a b: c :::d:::\"e \"f\"\"\"a\"",
			exp: []kv[Token, string]{
				{k: T_UNKNOWN, v: ":"},
				{k: T_VALUE, v: "a"},
				{k: T_FIELD, v: "b"},
				{k: T_VALUE, v: "c"},
				{k: T_UNKNOWN, v: ":"},
				{k: T_UNKNOWN, v: ":"},
				{k: T_FIELD, v: "d"},
				{k: T_UNKNOWN, v: ":"},
				{k: T_UNKNOWN, v: ":"},
				{k: T_QUOTATION_MARK},
				{k: T_VALUE, v: "e"},
				{k: T_QUOTATION_MARK},
				{k: T_VALUE, v: "f"},
				{k: T_QUOTATION_MARK},
				{k: T_QUOTATION_MARK},
				{k: T_QUOTATION_MARK},
				{k: T_VALUE, v: "a"},
				{k: T_QUOTATION_MARK},
				{k: T_EOF},
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
