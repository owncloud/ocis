package kql

import (
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

func source(text []byte) *string {
	str := string(text)
	return &str
}

func base(text []byte, pos position) *ast.Base {
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
			Source: source(text),
		},
	}
}

func root(elements interface{}, text []byte, pos position) (*ast.Ast, error) {
	return &ast.Ast{
		Base:  base(text, pos),
		Nodes: elements.([]ast.Node),
	}, nil
}

func nodes(head, tails interface{}) ([]ast.Node, error) {
	elems := []ast.Node{head.(ast.Node)}
	for _, tail := range toIfaceSlice(tails) {
		elem := toIfaceSlice(tail)[1]
		elems = append(elems, elem.(ast.Node))
	}
	return elems, nil
}

func textPropertyRestriction(k, v interface{}, text []byte, pos position) (*ast.TextPropertyRestriction, error) {
	key, err := toString(k)
	if err != nil {
		return nil, err
	}

	value, err := toString(v)
	if err != nil {
		return nil, err
	}

	return &ast.TextPropertyRestriction{
		Base:  base(text, pos),
		Key:   key,
		Value: value,
	}, nil
}

func phrase(v interface{}, text []byte, pos position) (*ast.Phrase, error) {
	value, err := toString(v)
	if err != nil {
		return nil, err
	}

	return &ast.Phrase{
		Base:  base(text, pos),
		Value: value,
	}, nil
}

func word(v interface{}, text []byte, pos position) (*ast.Word, error) {
	value, err := toString(v)
	if err != nil {
		return nil, err
	}

	return &ast.Word{
		Base:  base(text, pos),
		Value: value,
	}, nil
}

func booleanOperator(text []byte, pos position) (*ast.BooleanOperator, error) {
	return &ast.BooleanOperator{
		Base:  base(text, pos),
		Value: string(text),
	}, nil
}

func group(elements interface{}, text []byte, pos position) (*ast.Group, error) {
	return &ast.Group{
		Base:  base(text, pos),
		Nodes: elements.([]ast.Node),
	}, nil
}
