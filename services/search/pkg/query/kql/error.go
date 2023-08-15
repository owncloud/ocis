package kql

// StartsWithBinaryOperatorError records an error and the operation that caused it.
type StartsWithBinaryOperatorError struct {
	Op string
}

func (e *StartsWithBinaryOperatorError) Error() string {
	return "the expression can't begin from a binary operator: '" + e.Op + "'"
}
