// Package query provides functions to work with the different search query flavours.
package query

import (
	"io"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/bleve"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

// KqlToBleveQuery translates a kql query into a bleve query
func KqlToBleveQuery(r io.Reader, w io.Writer, opts ...kql.Option) error {
	a, err := kql.NewAST(r, opts...)
	if err != nil {
		return err
	}

	if err := bleve.Compile(w, a); err != nil {
		return err
	}

	return nil
}
