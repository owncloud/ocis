package kql

import (
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

func nodes(head, tails interface{}) ([]ast.Node, error) {
	node, err := toNode(head)
	if err != nil {
		return nil, err
	}

	var nodes []ast.Node

	for _, tail := range toIfaceSlice(tails) {
		node, err := toNode(toIfaceSlice(tail)[1])
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return append(append([]ast.Node{}, node), nodes...), nil
}

func textPropertyRestriction(k, v interface{}, text []byte, pos position) (*ast.StringProperty, error) {
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

	return &ast.StringProperty{
		Base:  b,
		Key:   key,
		Value: value,
	}, nil
}

func phrase(v interface{}, text []byte, pos position) (*ast.Phrase, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	value, err := toString(v)
	if err != nil {
		return nil, err
	}

	return &ast.Phrase{
		Base:  b,
		Value: value,
	}, nil
}

func word(v interface{}, text []byte, pos position) (*ast.Word, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	value, err := toString(v)
	if err != nil {
		return nil, err
	}

	return &ast.Word{
		Base:  b,
		Value: value,
	}, nil
}

func booleanOperator(text []byte, pos position) (*ast.BooleanOperator, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	return &ast.BooleanOperator{
		Base:  b,
		Value: string(text),
	}, nil
}

func group(n interface{}, text []byte, pos position) (*ast.Group, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	nodes, err := toNodes(n)
	if err != nil {
		return nil, err
	}

	return &ast.Group{
		Base:  b,
		Nodes: nodes,
	}, nil
}

func propertyGroup(k, n interface{}, text []byte, pos position) (*ast.KeyGroup, error) {
	b, err := base(text, pos)
	if err != nil {
		return nil, err
	}

	key, err := toString(k)
	if err != nil {
		return nil, err
	}

	var nodes []ast.Node

	for _, el := range toIfaceSlice(n) {
		node, err := toNode(el)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return &ast.KeyGroup{
		Base:  b,
		Key:   key,
		Nodes: nodes,
	}, nil
}
