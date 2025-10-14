package godata

import (
	"context"
	"fmt"
	"strconv"
)

type ExpandTokenType int

func (e ExpandTokenType) Value() int {
	return (int)(e)
}

const (
	ExpandTokenOpenParen ExpandTokenType = iota
	ExpandTokenCloseParen
	ExpandTokenNav
	ExpandTokenComma
	ExpandTokenSemicolon
	ExpandTokenEquals
	ExpandTokenLiteral
)

var GlobalExpandTokenizer = ExpandTokenizer()

// Represents an item to expand in an OData query. Tracks the path of the entity
// to expand and also the filter, levels, and reference options, etc.
type ExpandItem struct {
	Path    []*Token
	Filter  *GoDataFilterQuery
	At      *GoDataFilterQuery
	Search  *GoDataSearchQuery
	OrderBy *GoDataOrderByQuery
	Skip    *GoDataSkipQuery
	Top     *GoDataTopQuery
	Select  *GoDataSelectQuery
	Compute *GoDataComputeQuery
	Expand  *GoDataExpandQuery
	Levels  int
}

func ExpandTokenizer() *Tokenizer {
	t := Tokenizer{}
	t.Add("^\\(", ExpandTokenOpenParen)
	t.Add("^\\)", ExpandTokenCloseParen)
	t.Add("^/", ExpandTokenNav)
	t.Add("^,", ExpandTokenComma)
	t.Add("^;", ExpandTokenSemicolon)
	t.Add("^=", ExpandTokenEquals)
	t.Add("^[a-zA-Z0-9_\\'\\.:\\$ \\*]+", ExpandTokenLiteral)

	return &t
}

func ParseExpandString(ctx context.Context, expand string) (*GoDataExpandQuery, error) {
	tokens, err := GlobalExpandTokenizer.Tokenize(ctx, expand)

	if err != nil {
		return nil, err
	}

	stack := tokenStack{}
	queue := tokenQueue{}
	items := make([]*ExpandItem, 0)

	for len(tokens) > 0 {
		token := tokens[0]
		tokens = tokens[1:]

		if token.Value == "(" {
			queue.Enqueue(token)
			stack.Push(token)
		} else if token.Value == ")" {
			queue.Enqueue(token)
			stack.Pop()
		} else if token.Value == "," {
			if stack.Empty() {
				// no paren on the stack, parse this item and start a new queue
				item, err := ParseExpandItem(ctx, queue)
				if err != nil {
					return nil, err
				}
				items = append(items, item)
				queue = tokenQueue{}
			} else {
				// this comma is inside a nested expression, keep it in the queue
				queue.Enqueue(token)
			}
		} else {
			queue.Enqueue(token)
		}
	}

	if !stack.Empty() {
		return nil, BadRequestError("Mismatched parentheses in expand clause.")
	}

	item, err := ParseExpandItem(ctx, queue)
	if err != nil {
		return nil, err
	}
	items = append(items, item)

	return &GoDataExpandQuery{ExpandItems: items}, nil
}

func ParseExpandItem(ctx context.Context, input tokenQueue) (*ExpandItem, error) {

	item := &ExpandItem{}
	item.Path = []*Token{}

	stack := &tokenStack{}
	queue := &tokenQueue{}

	for !input.Empty() {
		token := input.Dequeue()
		if token.Value == "(" {
			if !stack.Empty() {
				// this is a nested slash, it belongs on the queue
				queue.Enqueue(token)
			} else {
				// top level slash means we're done parsing the path
				item.Path = append(item.Path, queue.Dequeue())
			}
			stack.Push(token)
		} else if token.Value == ")" {
			stack.Pop()
			if !stack.Empty() {
				// this is a nested slash, it belongs on the queue
				queue.Enqueue(token)
			} else {
				// top level slash means we're done parsing the options
				err := ParseExpandOption(ctx, queue, item)
				if err != nil {
					return nil, err
				}
				// reset the queue
				queue = &tokenQueue{}
			}
		} else if token.Value == "/" && stack.Empty() {
			if queue.Empty() {
				// Disallow extra leading and intermediate slash, like /Product and Product//Info
				return nil, BadRequestError("Empty path segment in expand clause.")
			}
			if input.Empty() {
				// Disallow extra trailing slash, like Product/
				return nil, BadRequestError("Empty path segment in expand clause.")
			}
			// at root level, slashes separate path segments
			item.Path = append(item.Path, queue.Dequeue())
		} else if token.Value == ";" && stack.Size == 1 {
			// semicolons only split expand options at the first level
			err := ParseExpandOption(ctx, queue, item)
			if err != nil {
				return nil, err
			}
			// reset the queue
			queue = &tokenQueue{}
		} else {
			queue.Enqueue(token)
		}
	}

	if !stack.Empty() {
		return nil, BadRequestError("Mismatched parentheses in expand clause.")
	}

	if !queue.Empty() {
		item.Path = append(item.Path, queue.Dequeue())
	}

	cfg, hasComplianceConfig := ctx.Value(odataCompliance).(OdataComplianceConfig)
	if !hasComplianceConfig {
		// Strict ODATA compliance by default.
		cfg = ComplianceStrict
	}

	if len(item.Path) == 0 && cfg&ComplianceIgnoreInvalidComma == 0 {
		return nil, BadRequestError("Extra comma in $expand.")
	}

	return item, nil
}

func ParseExpandOption(ctx context.Context, queue *tokenQueue, item *ExpandItem) error {
	head := queue.Dequeue().Value
	if queue.Head == nil {
		return BadRequestError("Invalid expand clause.")
	}
	queue.Dequeue() // drop the '=' from the front of the queue
	body := queue.GetValue()

	cfg, hasComplianceConfig := ctx.Value(odataCompliance).(OdataComplianceConfig)
	if !hasComplianceConfig {
		// Strict ODATA compliance by default.
		cfg = ComplianceStrict
	}

	if cfg == ComplianceStrict {
		// Enforce that only supported keywords are specified in expand.
		// The $levels keyword supported within expand is checked explicitly in addition to
		// keywords listed in supportedOdataKeywords[] which are permitted within expand and
		// at the top level of the odata query.
		if _, ok := supportedOdataKeywords[head]; !ok && head != "$levels" {
			return BadRequestError(fmt.Sprintf("Unsupported item '%s' in expand clause.", head))
		}
	}

	if head == "$filter" {
		filter, err := ParseFilterString(ctx, body)
		if err == nil {
			item.Filter = filter
		} else {
			return err
		}
	}

	if head == "at" {
		at, err := ParseFilterString(ctx, body)
		if err == nil {
			item.At = at
		} else {
			return err
		}
	}

	if head == "$search" {
		search, err := ParseSearchString(ctx, body)
		if err == nil {
			item.Search = search
		} else {
			return err
		}
	}

	if head == "$orderby" {
		orderby, err := ParseOrderByString(ctx, body)
		if err == nil {
			item.OrderBy = orderby
		} else {
			return err
		}
	}

	if head == "$skip" {
		skip, err := ParseSkipString(ctx, body)
		if err == nil {
			item.Skip = skip
		} else {
			return err
		}
	}

	if head == "$top" {
		top, err := ParseTopString(ctx, body)
		if err == nil {
			item.Top = top
		} else {
			return err
		}
	}

	if head == "$select" {
		sel, err := ParseSelectString(ctx, body)
		if err == nil {
			item.Select = sel
		} else {
			return err
		}
	}

	if head == "$compute" {
		comp, err := ParseComputeString(ctx, body)
		if err == nil {
			item.Compute = comp
		} else {
			return err
		}
	}

	if head == "$expand" {
		expand, err := ParseExpandString(ctx, body)
		if err == nil {
			item.Expand = expand
		} else {
			return err
		}
	}

	if head == "$levels" {
		i, err := strconv.Atoi(body)
		if err != nil {
			return err
		}
		item.Levels = i
	}

	return nil
}

func SemanticizeExpandQuery(
	expand *GoDataExpandQuery,
	service *GoDataService,
	entity *GoDataEntityType,
) error {

	if expand == nil {
		return nil
	}

	// Replace $levels with a nested expand clause
	for _, item := range expand.ExpandItems {
		if item.Levels > 0 {
			if item.Expand == nil {
				item.Expand = &GoDataExpandQuery{[]*ExpandItem{}}
			}
			// Future recursive calls to SemanticizeExpandQuery() will build out
			// this expand tree completely
			item.Expand.ExpandItems = append(
				item.Expand.ExpandItems,
				&ExpandItem{
					Path:   item.Path,
					Levels: item.Levels - 1,
				},
			)
			item.Levels = 0
		}
	}

	// we're gonna rebuild the items list, replacing wildcards where possible
	// TODO: can we save the garbage collector some heartache?
	newItems := []*ExpandItem{}

	for _, item := range expand.ExpandItems {
		if item.Path[0].Value == "*" {
			// replace wildcard with a copy of every navigation property
			for _, navProp := range service.NavigationPropertyLookup[entity] {
				path := []*Token{{Value: navProp.Name, Type: ExpandTokenLiteral}}
				newItem := &ExpandItem{
					Path:   append(path, item.Path[1:]...),
					Levels: item.Levels,
					Expand: item.Expand,
				}
				newItems = append(newItems, newItem)
			}
			// TODO: check for duplicates?
		} else {
			newItems = append(newItems, item)
		}
	}

	expand.ExpandItems = newItems

	for _, item := range expand.ExpandItems {
		err := semanticizeExpandItem(item, service, entity)
		if err != nil {
			return err
		}
	}

	return nil
}

func semanticizeExpandItem(
	item *ExpandItem,
	service *GoDataService,
	entity *GoDataEntityType,
) error {

	// TODO: allow multiple path segments in expand clause
	// TODO: handle $ref
	if len(item.Path) > 1 {
		return NotImplementedError("Multiple path segments not currently supported in expand clauses.")
	}

	navProps := service.NavigationPropertyLookup[entity]
	target := item.Path[len(item.Path)-1]
	if prop, ok := navProps[target.Value]; ok {
		target.SemanticType = SemanticTypeEntity
		entityType, err := service.LookupEntityType(prop.Type)
		if err != nil {
			return err
		}
		target.SemanticReference = entityType

		err = SemanticizeFilterQuery(item.Filter, service, entityType)
		if err != nil {
			return err
		}
		err = SemanticizeExpandQuery(item.Expand, service, entityType)
		if err != nil {
			return err
		}
		err = SemanticizeSelectQuery(item.Select, service, entityType)
		if err != nil {
			return err
		}
		err = SemanticizeOrderByQuery(item.OrderBy, service, entityType)
		if err != nil {
			return err
		}

	} else {
		return BadRequestError("Entity type " + entity.Name + " has no navigational property " + target.Value)
	}

	return nil
}
