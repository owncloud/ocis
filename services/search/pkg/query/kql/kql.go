package kql

import (
	"io"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

func NewAST(r io.Reader, opts ...Option) (*ast.Ast, error) {
	f, err := ParseReader("", r, opts...)
	if err != nil {
		return nil, err
	}
	return f.(*ast.Ast), nil
}
