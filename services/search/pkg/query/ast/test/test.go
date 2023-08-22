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
			cmpopts.IgnoreFields(ast.StringNode{}, "Base"),
			cmpopts.IgnoreFields(ast.OperatorNode{}, "Base"),
			cmpopts.IgnoreFields(ast.GroupNode{}, "Base"),
			cmpopts.IgnoreFields(ast.BooleanNode{}, "Base"),
		)...,
	)
}
