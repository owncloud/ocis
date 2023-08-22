package bleve

import (
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
		case *ast.StringNode:
			s = append(s, stringNode(n))
		case *ast.BooleanNode:
			// how should bleve treat an BooleanNode
		case *ast.GroupNode:
			// fixMe:
			// how should bleve treat an GroupNode
			// hint, recursion
		case *ast.OperatorNode:
			// fixMe:
			// how should bleve treat an OperatorNode
		}
	}

	_, err := io.WriteString(w, strings.Join(s, " "))

	return err
}

func stringNode(n *ast.StringNode) string {
	switch n.Key {
	case _tagKey:
		return _tagKey + ":" + n.Value
	case _nameKey:
		return _nameKey + ":" + n.Value
	case _contentKey:
		return _contentKey + ":" + n.Value
	}

	return ""
}
