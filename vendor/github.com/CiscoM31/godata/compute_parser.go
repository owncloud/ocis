package godata

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// The $compute query option must have a value which is a comma separated list of <expression> as <dynamic property name>
// See https://docs.oasis-open.org/odata/odata/v4.01/os/part2-url-conventions/odata-v4.01-os-part2-url-conventions.html#sec_SystemQueryOptioncompute
const computeAsSeparator = " as "

// Dynamic property names are restricted to case-insensitive a-z and the path separator /.
var computeFieldRegex = regexp.MustCompile("^[a-zA-Z/]+$")

type ComputeItem struct {
	Tree  *ParseNode // The compute expression parsed as a tree.
	Field string     // The name of the computed dynamic property.
}

// GlobalAllTokenParser is a Tokenizer which matches all tokens and ignores none. It differs from the
// GlobalExpressionTokenizer which ignores whitespace tokens.
var GlobalAllTokenParser *Tokenizer

func init() {
	t := NewExpressionParser().tokenizer
	t.TokenMatchers = append(t.IgnoreMatchers, t.TokenMatchers...)
	t.IgnoreMatchers = nil
	GlobalAllTokenParser = t
}

func ParseComputeString(ctx context.Context, compute string) (*GoDataComputeQuery, error) {
	items, err := SplitComputeItems(compute)
	if err != nil {
		return nil, err
	}

	result := make([]*ComputeItem, 0)
	fields := map[string]struct{}{}

	for _, v := range items {
		v = strings.TrimSpace(v)
		parts := strings.Split(v, computeAsSeparator)
		if len(parts) != 2 {
			return nil, &GoDataError{
				ResponseCode: 400,
				Message:      "Invalid $compute query option",
			}
		}
		field := strings.TrimSpace(parts[1])
		if !computeFieldRegex.MatchString(field) {
			return nil, &GoDataError{
				ResponseCode: 400,
				Message:      "Invalid $compute query option",
			}
		}

		if tree, err := GlobalExpressionParser.ParseExpressionString(ctx, parts[0]); err != nil {
			switch e := err.(type) {
			case *GoDataError:
				return nil, &GoDataError{
					ResponseCode: e.ResponseCode,
					Message:      fmt.Sprintf("Invalid $compute query option, %s", e.Message),
					Cause:        e,
				}
			default:
				return nil, &GoDataError{
					ResponseCode: 500,
					Message:      "Invalid $compute query option",
					Cause:        e,
				}
			}
		} else {
			if tree == nil {
				return nil, &GoDataError{
					ResponseCode: 500,
					Message:      "Invalid $compute query option",
				}
			}

			if _, ok := fields[field]; ok {
				return nil, &GoDataError{
					ResponseCode: 400,
					Message:      "Invalid $compute query option",
				}
			}

			fields[field] = struct{}{}

			result = append(result, &ComputeItem{
				Tree:  tree.Tree,
				Field: field,
			})
		}
	}

	return &GoDataComputeQuery{result, compute}, nil
}

// SplitComputeItems splits the input string based on the comma delimiter. It does so with awareness as to
// which commas delimit $compute items and which ones are an inline part of the item, such as a separator
// for function arguments.
//
// For example the input "someFunc(one,two) as three, 1 add 2 as four" results in the
// output ["someFunc(one,two) as three", "1 add 2 as four"]
func SplitComputeItems(in string) ([]string, error) {

	var ret []string

	tokens, err := GlobalAllTokenParser.Tokenize(context.Background(), in)
	if err != nil {
		return nil, err
	}

	item := strings.Builder{}
	parenGauge := 0

	for _, v := range tokens {
		switch v.Type {
		case ExpressionTokenOpenParen:
			parenGauge++
		case ExpressionTokenCloseParen:
			if parenGauge == 0 {
				return nil, errors.New("unmatched parentheses")
			}
			parenGauge--
		case ExpressionTokenComma:
			if parenGauge == 0 {
				ret = append(ret, item.String())
				item.Reset()
				continue
			}
		}

		item.WriteString(v.Value)
	}

	if parenGauge != 0 {
		return nil, errors.New("unmatched parentheses")
	}

	if item.Len() > 0 {
		ret = append(ret, item.String())
	}

	return ret, nil
}
