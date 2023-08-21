package test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

func DiffAst(x, y interface{}, opts ...cmp.Option) string {
	return cmp.Diff(
		x,
		y,
		append(
			opts,
			cmpopts.IgnoreFields(ast.Ast{}, "Base"),
			cmpopts.IgnoreFields(ast.Word{}, "Base"),
			cmpopts.IgnoreFields(ast.Phrase{}, "Base"),
			cmpopts.IgnoreFields(ast.StringProperty{}, "Base"),
			cmpopts.IgnoreFields(ast.BooleanOperator{}, "Base"),
			cmpopts.IgnoreFields(ast.Group{}, "Base"),
			cmpopts.IgnoreFields(ast.KeyGroup{}, "Base"),
		)...,
	)
}
