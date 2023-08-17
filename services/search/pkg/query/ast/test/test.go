package test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

var FullDictionaryAst = &ast.Ast{
	Nodes: []ast.Node{
		&ast.Phrase{Value: "author"},
		&ast.Phrase{Value: "John Smith"},
		&ast.Phrase{Value: "author"},
		&ast.Phrase{Value: "John Smith"},
		&ast.Phrase{Value: "author"},
		&ast.Phrase{Value: "John Smith"},
		&ast.TagQuery{Value: "foo"},
		&ast.Operator{Value: "AND"},
		&ast.TagQuery{Value: "bar"},
		&ast.NameQuery{Value: "book.pdf"},
		&ast.ContentQuery{Value: "letter.docx"},
	},
}

func DiffAst(x, y interface{}, opts ...cmp.Option) string {
	return cmp.Diff(
		x,
		y,
		append(
			opts,
			cmpopts.IgnoreFields(ast.Ast{}, "Base"),
			cmpopts.IgnoreFields(ast.Phrase{}, "Base"),
			cmpopts.IgnoreFields(ast.TagQuery{}, "Base"),
			cmpopts.IgnoreFields(ast.NameQuery{}, "Base"),
			cmpopts.IgnoreFields(ast.ContentQuery{}, "Base"),
			cmpopts.IgnoreFields(ast.Operator{}, "Base"),
		)...,
	)
}
