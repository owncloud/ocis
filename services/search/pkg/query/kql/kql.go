// Package kql provides the ability to work with kql queries.
package kql

import (
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

// Builder implements kql Builder interface
type Builder struct{}

// Build creates an ast.Ast based on a kql query
func (b Builder) Build(q string) (*ast.Ast, error) {
	f, err := Parse("", []byte(q))
	if err != nil {
		return nil, err
	}
	return f.(*ast.Ast), nil
}
