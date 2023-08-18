package bleve

import (
	"fmt"
	"io"
	"strings"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

const (
	_tagKey     = "tag"
	_nameKey    = "name"
	_contentKey = "content"
)

func Compile(w io.Writer, a *ast.Ast) error {
	var s []string

	for _, node := range a.Nodes {
		switch n := node.(type) {
		case *ast.TagQuery:
			s = append(s, fmt.Sprintf("%s:%s", _tagKey, n.Value))
			continue
		case *ast.NameQuery:
			s = append(s, fmt.Sprintf("%s:%s", _nameKey, n.Value))
			continue
		case *ast.ContentQuery:
			s = append(s, fmt.Sprintf("%s:%s", _contentKey, n.Value))
			continue
		case *ast.Operator:
			// fixMe:
			// how should bleve treat an operator
			continue
		case *ast.Phrase:
			// fixMe:
			// how should bleve treat an phrase
			continue
		case *ast.Group:
			// fixMe:
			// how should bleve treat an group
			// hint, recursion
			continue
		}
	}

	_, err := io.WriteString(w, strings.Join(s, " "))

	return err
}
