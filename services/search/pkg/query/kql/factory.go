package kql

import (
	"strings"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

func base(text []byte, pos position) (*ast.Base, error) {
	source, err := toString(text)
	if err != nil {
		return nil, err
	}

	return &ast.Base{
		Loc: &ast.Location{
			Start: ast.Position{
				Line:   pos.line,
				Column: pos.col,
			},
			End: ast.Position{
				Line:   pos.line,
				Column: pos.col + len(text),
			},
			Source: &source,
		},
	}, nil
}

func root(n interface{}, text []byte, pos position) (*ast.Ast, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	nodes, err := toNodes(n)
	if err != nil {
		return nil, err
	}

	return &ast.Ast{
		Base:  b,
		Nodes: normalize(nodes),
	}, nil
}

func nodes(head, t interface{}) ([]ast.Node, error) {
	node, err := toNode(head)
	if err != nil {
		return nil, err
	}

	tails := toIfaceSlice(t)
	for i, tail := range tails {
		tails[i] = toIfaceSlice(tail)[1]
	}

	nodes, err := toNodes(tails)
	if err != nil {
		return nil, err
	}

	return append(append([]ast.Node{}, node), nodes...), nil
}

func stringNode(k, v interface{}, text []byte, pos position) (*ast.StringNode, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	key, err := toString(k)
	if err != nil {
		return nil, err
	}

	value, err := toString(v)
	if err != nil {
		return nil, err
	}

	return &ast.StringNode{
		Base:  b,
		Key:   key,
		Value: value,
	}, nil
}

func booleanNode(k, v interface{}, text []byte, pos position) (*ast.BooleanNode, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	key, err := toString(k)
	if err != nil {
		return nil, err
	}

	value, err := toString(v)
	if err != nil {
		return nil, err
	}

	return &ast.BooleanNode{
		Base:  b,
		Key:   key,
		Value: strings.ToLower(value) == "true",
	}, nil
}

func operatorNode(text []byte, pos position) (*ast.OperatorNode, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	return &ast.OperatorNode{
		Base:  b,
		Value: string(text),
	}, nil
}

func groupNode(k, n interface{}, text []byte, pos position) (*ast.GroupNode, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	key, _ := toString(k)

	nodes, err := toNodes(n)
	if err != nil {
		return nil, err
	}

	return &ast.GroupNode{
		Base:  b,
		Key:   key,
		Nodes: nodes,
	}, nil
}

var source = "implicitly operator"
var operatorNodeAnd = ast.OperatorNode{Base: &ast.Base{Loc: &ast.Location{Source: &source}}, Value: BoolAND}
var operatorNodeOr = ast.OperatorNode{Base: &ast.Base{Loc: &ast.Location{Source: &source}}, Value: BoolOR}

// normalize Populate the implicit logical operators in the ast
//
// https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference#constructing-free-text-queries-using-kql
// If there are multiple free-text expressions without any operators in between them, the query behavior is the same as using the AND operator.
// "John Smith" "Jane Smith"
// This functionally is the same as using the OR Boolean operator, as follows:
// "John Smith" AND "Jane Smith"
//
// https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference#using-multiple-property-restrictions-within-a-kql-query
// When you use multiple instances of the same property restriction, matches are based on the union of the property restrictions in the KQL query.
// author:"John Smith" author:"Jane Smith"
// This functionally is the same as using the OR Boolean operator, as follows:
// author:"John Smith" OR author:"Jane Smith"
//
// When you use different property restrictions, matches are based on an intersection of the property restrictions in the KQL query, as follows:
// author:"John Smith" filetype:docx
// This is the same as using the AND Boolean operator, as follows:
// author:"John Smith" AND filetype:docx
//
// https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference#grouping-property-restrictions-within-a-kql-query
// author:("John Smith" "Jane Smith")
// This is the same as using the AND Boolean operator, as follows:
// author:"John Smith" AND author:"Jane Smith"
func normalize(nodes []ast.Node) []ast.Node {
	res := make([]ast.Node, 0, len(nodes))
	var currentNode ast.Node
	var prevKey, currentKey *string
	for _, node := range nodes {
		switch n := node.(type) {
		case *ast.StringNode:
			if prevKey == nil {
				prevKey = &n.Key
				res = append(res, node)
				continue
			}
			currentNode = n
			currentKey = &n.Key
		case *ast.GroupNode:
			n.Nodes = normalize(n.Nodes)
			if prevKey == nil {
				prevKey = &n.Key
				res = append(res, n)
				continue
			}
			currentNode = n
			currentKey = &n.Key
		default:
			prevKey = nil
			res = append(res, node)
		}
		if prevKey != nil && currentKey != nil {
			if *prevKey == *currentKey && *prevKey != "" {
				res = append(res, &operatorNodeOr, currentNode)
			} else {
				res = append(res, &operatorNodeAnd, currentNode)
			}
			currentNode = nil
			currentKey = nil
			prevKey = nil
			continue
		}
	}
	return res
}
