package provider

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/CiscoM31/godata"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

func init() {
	// add (ap)prox filter
	godata.GlobalFilterTokenizer = FilterTokenizer()
	godata.GlobalFilterParser.DefineOperator("ap", 2, godata.OpAssociationLeft, 4, false)
}

// BuildBleveQuery converts a GoDataFilterQuery into a bleve query
func BuildBleveQuery(r *godata.GoDataFilterQuery) (query.Query, error) {
	return recursiveBuildQuery(r.Tree)
}

// Builds the filter recursively using DFS
func recursiveBuildQuery(n *godata.ParseNode) (query.Query, error) {
	if n.Token.Type == godata.FilterTokenFunc {
		switch n.Token.Value {
		case "startswith":
			if len(n.Children) != 2 {
				return nil, errors.New("startswith match must have two children")
			}
			if n.Children[0].Token.Type != godata.FilterTokenLiteral {
				return nil, errors.New("startswith expected a literal as the first param")
			}
			if n.Children[1].Token.Type != godata.FilterTokenString {
				return nil, errors.New("startswith expected a string as the second param")
			}
			q := bleve.NewTermQuery(n.Children[1].Token.Value)
			q.SetField(n.Children[0].Token.Value)
			return q, nil
		default:
			return nil, godata.NotImplementedError(n.Token.Value + " is not implemented.")
		}
	}
	if n.Token.Type == godata.FilterTokenLogical {
		switch n.Token.Value {
		case "eq":
			if len(n.Children) != 2 {
				return nil, errors.New("equality match must have two children")
			}
			if n.Children[0].Token.Type != godata.FilterTokenLiteral {
				return nil, errors.New("equality expected a literal on the lhs")
			}
			if n.Children[1].Token.Type == godata.FilterTokenString {
				// string tokens are enclosed with 'some string'
				// ' is escaped as ''
				// TODO unescape '' as '
				// http://docs.oasis-open.org/odata/odata/v4.01/cs01/part2-url-conventions/odata-v4.01-cs01-part2-url-conventions.html#sec_URLComponents
				q := bleve.NewTermQuery(n.Children[1].Token.Value[1 : len(n.Children[1].Token.Value)-1])
				q.SetField(n.Children[0].Token.Value)
				return q, nil
			} else if n.Children[1].Token.Type == godata.FilterTokenInteger {
				v, err := strconv.ParseFloat(n.Children[1].Token.Value, 64)
				if err != nil {
					return nil, err
				}
				incl := true
				q := bleve.NewNumericRangeInclusiveQuery(&v, &v, &incl, &incl)
				q.SetField(n.Children[0].Token.Value)
				return q, nil
			}
			return nil, fmt.Errorf("equality expected a string or int on the rhs, got %d", n.Children[1].Token.Type)
		case "and":
			q := query.NewConjunctionQuery([]query.Query{})
			for _, child := range n.Children {
				subQuery, err := recursiveBuildQuery(child)
				if err != nil {
					return nil, err
				}
				if subQuery != nil {
					q.AddQuery(subQuery)
				}
			}
			return q, nil
		case "or":
			q := query.NewDisjunctionQuery([]query.Query{})
			for _, child := range n.Children {
				subQuery, err := recursiveBuildQuery(child)
				if err != nil {
					return nil, err
				}
				if subQuery != nil {
					q.AddQuery(subQuery)
				}
			}
			return q, nil
		case "Not":
			if len(n.Children) != 1 {
				return nil, errors.New("not filter must have only one child")
			}
			subQuery, err := recursiveBuildQuery(n.Children[0])
			if err != nil {
				return nil, err
			}
			q := query.NewBooleanQuery(nil, nil, []query.Query{subQuery})
			return q, nil
		default:
			return nil, godata.NotImplementedError(n.Token.Value + " is not implemented.")
		}
	}

	return nil, godata.NotImplementedError(n.Token.Value + " is not implemented.")
}
