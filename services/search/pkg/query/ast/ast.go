// Package ast provides available ast nodes.
package ast

// Node represents abstract syntax tree node
type Node interface {
	Location() *Location
}

// Position represents a specific location in the source
type Position struct {
	Line   int
	Column int
}

// Location represents the location of a node in the AST
type Location struct {
	Start  Position `json:"start"`
	End    Position `json:"end"`
	Source *string  `json:"source,omitempty"`
}

// Base contains shared node attributes
// each node should inherit from this
type Base struct {
	Loc *Location
}

// Location is the source location of the Node
func (b *Base) Location() *Location { return b.Loc }

// Ast represents the query - node structure as abstract syntax tree
type Ast struct {
	*Base
	Nodes []Node `json:"body"`
}

// StringNode represents a string value
type StringNode struct {
	*Base
	Key   string
	Value string
}

// BooleanNode represents a bool value
type BooleanNode struct {
	*Base
	Key   string
	Value bool
}

// OperatorNode represents an operator value like
// AND, OR, NOT, =, <= ... and so on
type OperatorNode struct {
	*Base
	Value string
}

// GroupNode represents a collection of many grouped nodes
type GroupNode struct {
	*Base
	Key   string
	Nodes []Node
}
