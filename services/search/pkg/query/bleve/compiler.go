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
		case *ast.TextPropertyRestriction:
			s = append(s, textQuery(n))
			continue
		case *ast.Phrase, *ast.Word:
			// fixMe:
			// how should bleve treat an phrase or word
			continue
		case *ast.Group:
			// fixMe:
			// how should bleve treat an group
			// hint, recursion
			continue
		case *ast.BooleanOperator:
			// fixMe:
			// how should bleve treat an boolean operator
			continue
		}
	}

	_, err := io.WriteString(w, strings.Join(s, " "))

	return err
}

func textQuery(n *ast.TextPropertyRestriction) string {
	switch n.Key {
	case _tagKey:
		return fmt.Sprintf("%s:%s", _tagKey, n.Value)
	case _nameKey:
		return fmt.Sprintf("%s:%s", _nameKey, n.Value)
	case _contentKey:
		return fmt.Sprintf("%s:%s", _contentKey, n.Value)
	}

	return ""
}
