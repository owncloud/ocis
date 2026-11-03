package query

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/ast"
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

// UnsupportedTimeRangeError records an error and the value that caused it.
type UnsupportedTimeRangeError struct {
	Value interface{}
}

func (e UnsupportedTimeRangeError) Error() string {
	return fmt.Sprintf("unable to convert '%v' to a time range", e.Value)
}

// UncompilableQueryError is returned when the query can't be compiled,
// e.g. because it uses an unsupported operator.
type UncompilableQueryError struct{}

func (e UncompilableQueryError) Error() string {
	return "can not compile the query"
}

// IsValidationError reports whether err was caused by an invalid user query.
func IsValidationError(err error) bool {
	switch err.(type) {
	case *StartsWithBinaryOperatorError, *NamedGroupInvalidNodesError, *UnsupportedTimeRangeError, *UncompilableQueryError:
		return true
	}
	return false
}
