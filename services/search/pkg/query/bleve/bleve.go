// Package bleve provides the ability to work with bleve queries.
package bleve

import (
	bQuery "github.com/blevesearch/bleve/v2/search/query"

	"github.com/owncloud/ocis/v2/services/search/pkg/query"
)

// Creator is combines a Builder and a Compiler which is used to Create the query.
type Creator[T any] struct {
	builder  query.Builder
	compiler query.Compiler[T]
}

// Create implements the Creator interface
func (c Creator[T]) Create(qs string) (T, error) {
	var t T
	builderAst, err := c.builder.Build(qs)
	if err != nil {
		return t, err
	}

	t, err = c.compiler.Compile(builderAst)
	if err != nil {
		return t, err
	}

	return t, nil
}

// LegacyCreator exposes an ocis legacy bleve query creator.
var LegacyCreator = Creator[bQuery.Query]{LegacyBuilder{}, LegacyCompiler{}}
