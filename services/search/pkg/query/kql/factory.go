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
		Nodes: nodes,
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
