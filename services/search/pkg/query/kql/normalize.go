package kql

import (
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

var implicitOperatorNodeSource = "implicitly operator"
var operatorNodeAnd = ast.OperatorNode{Base: &ast.Base{Loc: &ast.Location{Source: &implicitOperatorNodeSource}}, Value: BoolAND}
var operatorNodeOr = ast.OperatorNode{Base: &ast.Base{Loc: &ast.Location{Source: &implicitOperatorNodeSource}}, Value: BoolOR}

// NormalizeNodes Populate the implicit logical operators in the ast
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
func NormalizeNodes(nodes []ast.Node) ([]ast.Node, error) {
	res := make([]ast.Node, 0, len(nodes))
	var currentNode ast.Node
	var prevKey, currentKey *string
	var operator *ast.OperatorNode
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
		case *ast.BooleanNode:
			if prevKey == nil {
				prevKey = &n.Key
				res = append(res, node)
				continue
			}
			currentNode = n
			currentKey = &n.Key
		case *ast.GroupNode:
			var err error
			n.Nodes, err = NormalizeNodes(n.Nodes)
			if err != nil {
				return nil, err
			}
			if prevKey == nil {
				prevKey = &n.Key
				res = append(res, n)
				continue
			}
			currentNode = n
			currentKey = &n.Key
		case *ast.OperatorNode:
			if n.Value == BoolNOT {
				if prevKey == nil {
					res = append(res, n)
				} else {
					operator = n
				}
			} else {
				if prevKey == nil {
					return nil, &StartsWithBinaryOperatorError{Op: n.Value}
				}
				prevKey = nil
				res = append(res, node)
			}
		default:
			prevKey = nil
			res = append(res, node)
		}
		if prevKey != nil && currentKey != nil {
			if *prevKey == *currentKey && *prevKey != "" {
				res = append(res, &operatorNodeOr)
			} else {
				res = append(res, &operatorNodeAnd)
			}
			if operator != nil {
				res = append(res, operator)
				operator = nil
			}
			res = append(res, currentNode)

			prevKey = currentKey
			currentNode = nil
			currentKey = nil
			continue
		}
	}

	return trimOrphan(res), nil
}

func trimOrphan(nodes []ast.Node) []ast.Node {
	offset := len(nodes)
	for i := len(nodes) - 1; i >= 0; i-- {
		if _, ok := nodes[i].(*ast.OperatorNode); ok {
			offset--
		} else {
			break
		}
	}
	return nodes[:offset]
}
