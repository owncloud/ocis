package kql

import (
	"fmt"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

// StartsWithBinaryOperatorError records an error and the operation that caused it.
type StartsWithBinaryOperatorError struct {
	Node *ast.OperatorNode
}

func (e StartsWithBinaryOperatorError) Error() string {
	return "the expression can't begin from a binary operator: '" + e.Node.Value + "'"
}

// NamedGroupInvalidNodesError records an error and the operation that caused it.
type NamedGroupInvalidNodesError struct {
	Node ast.Node
}

func (e NamedGroupInvalidNodesError) Error() string {
	return fmt.Errorf(
		"'%T' - '%v' - '%v' is not valid",
		e.Node,
		ast.NodeKey(e.Node),
		ast.NodeValue(e.Node),
	).Error()
}

// UnsupportedTimeRange records an error and the value that caused it.
type UnsupportedTimeRange struct {
	Value interface{}
}

func (e UnsupportedTimeRange) Error() string {
	return fmt.Errorf(
		"unable to convert '%v' to a time range",
		e.Value,
	).Error()
}
