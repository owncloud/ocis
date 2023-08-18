package ast

type Node interface {
	Location() *Location
}

type Position struct {
	Line   int
	Column int
}

type Location struct {
	Start  Position `json:"start"`
	End    Position `json:"end"`
	Source *string  `json:"source,omitempty"`
}

type Base struct {
	Loc *Location
}

func (b *Base) Location() *Location { return b.Loc }

type Ast struct {
	*Base
	Nodes []Node `json:"body"`
}

type TextPropertyRestriction struct {
	*Base
	Key   string
	Value string
}

type Word struct {
	*Base
	Value string
}

type Phrase struct {
	*Base
	Value string
}

type BooleanOperator struct {
	*Base
	Value string
}

type Group struct {
	*Base
	Nodes []Node `json:"body"`
}
