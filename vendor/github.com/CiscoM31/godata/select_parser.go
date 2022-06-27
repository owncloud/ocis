package godata

import (
	"context"
	"errors"
	"strings"
)

type SelectItem struct {
	Segments []*Token
}

func ParseSelectString(ctx context.Context, sel string) (*GoDataSelectQuery, error) {
	items := strings.Split(sel, ",")

	result := []*SelectItem{}

	for _, item := range items {
		segments := []*Token{}
		for _, val := range strings.Split(item, "/") {
			segments = append(segments, &Token{Value: val})
		}
		result = append(result, &SelectItem{segments})
	}

	return &GoDataSelectQuery{result, sel}, nil
}

func SemanticizeSelectQuery(sel *GoDataSelectQuery, service *GoDataService, entity *GoDataEntityType) error {
	if sel == nil {
		return nil
	}

	newItems := []*SelectItem{}

	// replace wildcards with every property of the entity
	for _, item := range sel.SelectItems {
		// TODO: allow multiple path segments
		if len(item.Segments) > 1 {
			return NotImplementedError("Multiple path segments in select clauses are not yet supported.")
		}

		if item.Segments[0].Value == "*" {
			for _, prop := range service.PropertyLookup[entity] {
				newItems = append(newItems, &SelectItem{[]*Token{{Value: prop.Name}}})
			}
		} else {
			newItems = append(newItems, item)
		}
	}

	sel.SelectItems = newItems

	for _, item := range sel.SelectItems {
		if prop, ok := service.PropertyLookup[entity][item.Segments[0].Value]; ok {
			item.Segments[0].SemanticType = SemanticTypeProperty
			item.Segments[0].SemanticReference = prop
		} else {
			return errors.New("Entity " + entity.Name + " has no property " + item.Segments[0].Value)
		}
	}

	return nil
}
