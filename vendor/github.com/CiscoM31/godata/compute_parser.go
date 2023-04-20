package godata

import (
	"context"
	"regexp"
	"strings"
)

// The $compute query option must have a value which is a comma separated list of <expression> as <dynamic property name>
// See https://docs.oasis-open.org/odata/odata/v4.01/os/part2-url-conventions/odata-v4.01-os-part2-url-conventions.html#sec_SystemQueryOptioncompute
const computeAsSeparator = " as "

// Dynamic property names are restricted to case-insensitive a-z
var computeFieldRegex = regexp.MustCompile("^[a-zA-Z]+$")

type ComputeItem struct {
	Tree  *ParseNode // The compute expression parsed as a tree.
	Field string     // The name of the computed dynamic property.
}

func ParseComputeString(ctx context.Context, compute string) (*GoDataComputeQuery, error) {
	items := strings.Split(compute, ",")

	result := make([]*ComputeItem, 0)

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
					Message:      "Invalid $compute query option",
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
			result = append(result, &ComputeItem{
				Tree:  tree.Tree,
				Field: field,
			})
		}
	}

	return &GoDataComputeQuery{result, compute}, nil
}
