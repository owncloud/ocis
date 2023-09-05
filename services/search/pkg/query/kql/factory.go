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

func buildAST(n interface{}, text []byte, pos position) (*ast.Ast, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	nodes, err := toNodes[ast.Node](n)
	if err != nil {
		return nil, err
	}

	normalizedNodes, err := NormalizeNodes(nodes)
	if err != nil {
		return nil, err
	}

	return &ast.Ast{
		Base:  b,
		Nodes: normalizedNodes,
	}, nil
}

func buildNodes(e interface{}) ([]ast.Node, error) {
	maybeNodesGroups := toIfaceSlice(e)

	nodes := make([]ast.Node, len(maybeNodesGroups))
	for i, maybeNodesGroup := range maybeNodesGroups {
		node, err := toNode[ast.Node](toIfaceSlice(maybeNodesGroup)[1])
		if err != nil {
			return nil, err
		}

		nodes[i] = node
	}

	return nodes, nil
}

func buildStringNode(k, v interface{}, text []byte, pos position) (*ast.StringNode, error) {
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

func buildDateTimeNode(k, o, v interface{}, text []byte, pos position) (*ast.DateTimeNode, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	operator, err := toNode[*ast.OperatorNode](o)
	if err != nil {
		return nil, err
	}

	key, err := toString(k)
	if err != nil {
		return nil, err
	}

	value, err := toTime(v)
	if err != nil {
		return nil, err
	}

	return &ast.DateTimeNode{
		Base:     b,
		Key:      key,
		Operator: operator,
		Value:    value,
	}, nil
}

func buildBooleanNode(k, v interface{}, text []byte, pos position) (*ast.BooleanNode, error) {
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

func buildOperatorNode(text []byte, pos position) (*ast.OperatorNode, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	value, err := toString(text)
	if err != nil {
		return nil, err
	}

	return &ast.OperatorNode{
		Base:  b,
		Value: value,
	}, nil
}

func buildGroupNode(k, n interface{}, text []byte, pos position) (*ast.GroupNode, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	key, _ := toString(k)

	nodes, err := toNodes[ast.Node](n)
	if err != nil {
		return nil, err
	}

	return &ast.GroupNode{
		Base:  b,
		Key:   key,
		Nodes: nodes,
	}, nil
}
