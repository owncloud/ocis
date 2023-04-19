package godata

import (
	"context"
	"strings"
)

const (
	ASC  = "asc"
	DESC = "desc"
)

type OrderByItem struct {
	Field *Token            // The raw value of the orderby field or expression.
	Tree  *GoDataExpression // The orderby expression parsed as a tree.
	Order string            // Ascending or descending order.
}

func ParseOrderByString(ctx context.Context, orderby string) (*GoDataOrderByQuery, error) {
	return GlobalExpressionParser.ParseOrderByString(ctx, orderby)
}

// The value of the $orderby System Query option contains a comma-separated
// list of expressions whose primitive result values are used to sort the items.
// The service MUST order by the specified property in ascending order.
// 4.01 services MUST support case-insensitive values for asc and desc.
func (p *ExpressionParser) ParseOrderByString(ctx context.Context, orderby string) (*GoDataOrderByQuery, error) {
	items := strings.Split(orderby, ",")

	result := make([]*OrderByItem, 0)

	for _, v := range items {
		v = strings.TrimSpace(v)

		cfg, hasComplianceConfig := ctx.Value(odataCompliance).(OdataComplianceConfig)
		if !hasComplianceConfig {
			// Strict ODATA compliance by default.
			cfg = ComplianceStrict
		}

		if len(v) == 0 && cfg&ComplianceIgnoreInvalidComma == 0 {
			return nil, BadRequestError("Extra comma in $orderby.")
		}

		var order string
		vLower := strings.ToLower(v)
		if strings.HasSuffix(vLower, " "+ASC) {
			order = ASC
		} else if strings.HasSuffix(vLower, " "+DESC) {
			order = DESC
		}
		if order == "" {
			order = ASC // default order
		} else {
			v = v[:len(v)-len(order)]
			v = strings.TrimSpace(v)
		}

		if tree, err := p.ParseExpressionString(ctx, v); err != nil {
			switch e := err.(type) {
			case *GoDataError:
				return nil, &GoDataError{
					ResponseCode: e.ResponseCode,
					Message:      "Invalid $orderby query option",
					Cause:        e,
				}
			default:
				return nil, &GoDataError{
					ResponseCode: 500,
					Message:      "Invalid $orderby query option",
					Cause:        e,
				}
			}
		} else {
			result = append(result, &OrderByItem{
				Field: &Token{Value: unescapeUtfEncoding(v)},
				Tree:  tree,
				Order: order,
			})

		}
	}

	return &GoDataOrderByQuery{result, orderby}, nil
}

func SemanticizeOrderByQuery(orderby *GoDataOrderByQuery, service *GoDataService, entity *GoDataEntityType) error {
	if orderby == nil {
		return nil
	}

	for _, item := range orderby.OrderByItems {
		if prop, ok := service.PropertyLookup[entity][item.Field.Value]; ok {
			item.Field.SemanticType = SemanticTypeProperty
			item.Field.SemanticReference = prop
		} else {
			return BadRequestError("No property " + item.Field.Value + " for entity " + entity.Name)
		}
	}

	return nil
}
