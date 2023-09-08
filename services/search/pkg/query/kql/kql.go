// Package kql provides the ability to work with kql queries.
package kql

import (
	"errors"
	"strings"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

// The operator node value definition
const (
	// BoolAND connect two nodes with "AND"
	BoolAND = "AND"
	// BoolOR connect two nodes with "OR"
	BoolOR = "OR"
	// BoolNOT connect two nodes with "NOT"
	BoolNOT = "NOT"
)

// Builder implements kql Builder interface
type Builder struct{}

// Build creates an ast.Ast based on a kql query
func (b Builder) Build(q string) (*ast.Ast, error) {
	f, err := Parse("", []byte(q))
	if err != nil {
		var list errList
		errors.As(err, &list)

		for _, listError := range list {
			var parserError *parserError
			switch {
			case errors.As(listError, &parserError):
				return nil, listError
			}
		}
	}

	return f.(*ast.Ast), nil
}

// incorporateNode connects a leading node with the rest
func incorporateNode(headNode ast.Node, tailNodes ...ast.Node) *ast.OperatorNode {
	switch headNode.(type) {
	case *ast.OperatorNode:
		return nil
	}

	var nearestNeighborNode ast.Node
	var nearestNeighborOperators []*ast.OperatorNode

l:
	for _, tailNode := range tailNodes {
		switch node := tailNode.(type) {
		case *ast.OperatorNode:
			nearestNeighborOperators = append(nearestNeighborOperators, node)
		default:
			nearestNeighborNode = node
			break l
		}
	}

	if nearestNeighborNode == nil {
		return nil
	}

	headKey := strings.ToLower(nodeKey(headNode))
	neighborKey := strings.ToLower(nodeKey(nearestNeighborNode))

	connection := &ast.OperatorNode{
		Base:  &ast.Base{Loc: &ast.Location{Source: &[]string{"implicitly operator"}[0]}},
		Value: BoolAND,
	}

	// if the current node and the neighbor node have the same key
	// the connection is of type OR, same applies if no keys are in place
	//
	//		"" == ""
	//
	// spec: same
	//		author:"John Smith" author:"Jane Smith"
	//		author:"John Smith" OR author:"Jane Smith"
	if headKey == neighborKey {
		connection.Value = BoolOR
	}

	// decisions based on nearest neighbor node
	switch nearestNeighborNode.(type) {
	// nearest neighbor node type could change the default case
	// docs says, if the next value node:
	//
	//		is a group AND has no key
	//
	// even if the current node has none too, which normal leads to SAME KEY OR
	//
	// 		it should be an AND edge
	//
	// spec: same
	// 		cat (dog OR fox)
	// 		cat AND (dog OR fox)
	//
	// note:
	// 		sounds contradictory to me
	case *ast.GroupNode:
		if headKey == "" && neighborKey == "" {
			connection.Value = BoolAND
		}
	}

	// decisions based on nearest neighbor operators
	for i, node := range nearestNeighborOperators {
		// consider direct neighbor operator only
		if i == 0 {
			// no connection is necessary here because an `AND` or `OR` edge is already present
			// exit
			for _, skipValue := range []string{BoolOR, BoolAND} {
				if node.Value == skipValue {
					return nil
				}
			}

			// if neighbor node negotiates, AND edge is needed
			//
			// spec: same
			// 		cat -dog
			// 		cat AND NOT dog
			if node.Value == BoolNOT {
				connection.Value = BoolAND
			}
		}
	}

	return connection
}

// nodeKey tries to return a node key
func nodeKey(n ast.Node) string {
	switch node := n.(type) {
	case *ast.StringNode:
		return node.Key
	case *ast.DateTimeNode:
		return node.Key
	case *ast.BooleanNode:
		return node.Key
	case *ast.GroupNode:
		return node.Key
	default:
		return ""
	}
}
