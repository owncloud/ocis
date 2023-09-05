package bleve

import (
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	bleveQuery "github.com/blevesearch/bleve/v2/search/query"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

var _fields = map[string]string{
	"rootid":   "RootID",
	"path":     "Path",
	"id":       "ID",
	"name":     "Name",
	"size":     "Size",
	"mtime":    "Mtime",
	"mimetype": "MimeType",
	"type":     "Type",
	"tag":      "Tags",
	"tags":     "Tags",
	"content":  "Content",
	"hidden":   "Hidden",
}

// Compiler represents a KQL query search string to the bleve query formatter.
type Compiler struct{}

// Compile implements the query formatter which converts the KQL query search string to the bleve query.
func (c Compiler) Compile(givenAst *ast.Ast) (bleveQuery.Query, error) {
	q, err := compile(givenAst)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func compile(a *ast.Ast) (bleveQuery.Query, error) {
	q, _ := walk(0, a.Nodes)
	switch q.(type) {
	case *bleveQuery.ConjunctionQuery, *bleveQuery.DisjunctionQuery:
		return q, nil
	}
	return bleve.NewConjunctionQuery(q), nil
}

func walk(offset int, nodes []ast.Node) (bleveQuery.Query, int) {
	var prev, next bleveQuery.Query
	var operator *ast.OperatorNode
	var isGroup bool
	for i := offset; i < len(nodes); i++ {
		switch n := nodes[i].(type) {
		case *ast.StringNode:
			k := getField(n.Key)
			v := strings.ReplaceAll(n.Value, " ", `\ `)

			if k != "Hidden" {
				v = strings.ToLower(v)
			}

			q := bleveQuery.NewQueryStringQuery(k + ":" + v)
			if prev == nil {
				prev = q
			} else {
				next = q
			}
		case *ast.DateTimeNode:
			k := getField(n.Key)
			// fixMe: should be bleveQuery.NewDateRangeQuery or bleveQuery.NewDateRangeInclusiveQuery!?
			q := bleveQuery.NewQueryStringQuery(k + ":" + n.Operator.Value + "\"" + n.Value.Format(time.RFC3339Nano) + "\"")
			if prev == nil {
				prev = q
			} else {
				next = q
			}
		case *ast.BooleanNode:
			q := bleveQuery.NewQueryStringQuery(getField(n.Key) + fmt.Sprintf(":%v", n.Value))
			if prev == nil {
				prev = q
			} else {
				next = q
			}
		case *ast.GroupNode:
			if n.Key != "" {
				n = normalizeGroupingProperty(n)
			}
			q, _ := walk(0, n.Nodes)
			if prev == nil {
				prev = q
				isGroup = true
			} else {
				next = q
			}
		case *ast.OperatorNode:
			if n.Value == kql.BoolAND || n.Value == kql.BoolOR {
				operator = n
			} else if n.Value == kql.BoolNOT {
				next, offset = nextNode(i+1, nodes)
				q := bleve.NewBooleanQuery()
				q.AddMustNot(next)
				next = q
			}
		}
		if prev != nil && next != nil && operator != nil {
			prev = mapBinary(operator, prev, next, isGroup)
			isGroup = false
			operator = nil
			next = nil
		}
		if i < offset {
			i = offset
		}
	}
	return prev, offset
}

func nextNode(offset int, nodes []ast.Node) (bleveQuery.Query, int) {
	if n, ok := nodes[offset].(*ast.GroupNode); ok {
		gq, _ := walk(0, n.Nodes)
		return gq, offset + 1
	}
	if n, ok := nodes[offset].(*ast.OperatorNode); ok {
		if n.Value == kql.BoolNOT {
			return walk(offset, nodes)
		}
	}
	one := nodes[:offset+1]
	return walk(offset, one)
}

func mapBinary(operator *ast.OperatorNode, ln, rn bleveQuery.Query, leftIsGroup bool) bleveQuery.Query {
	if operator.Value == kql.BoolAND {
		if left, ok := ln.(*bleveQuery.ConjunctionQuery); ok {
			left.AddQuery(rn)
			return left
		}
		if left, ok := ln.(*bleveQuery.DisjunctionQuery); ok && !leftIsGroup {
			last := left.Disjuncts[len(left.Disjuncts)-1]
			rn = bleveQuery.NewConjunctionQuery([]bleveQuery.Query{
				last,
				rn,
			})
			dj := bleveQuery.NewDisjunctionQuery(left.Disjuncts[:len(left.Disjuncts)-1])
			dj.AddQuery(rn)
			return dj
		}
		return bleveQuery.NewConjunctionQuery([]bleveQuery.Query{
			ln,
			rn,
		})
	}
	if operator.Value == kql.BoolOR {
		if left, ok := ln.(*bleveQuery.DisjunctionQuery); ok {
			left.AddQuery(rn)
			return left
		}
		return bleveQuery.NewDisjunctionQuery([]bleveQuery.Query{
			ln,
			rn,
		})
	}
	return bleveQuery.NewConjunctionQuery([]bleveQuery.Query{
		ln,
		rn,
	})
}

func getField(name string) string {
	if name == "" {
		return "Name"
	}
	if _, ok := _fields[strings.ToLower(name)]; ok {
		return _fields[strings.ToLower(name)]
	}
	return name
}

func normalizeGroupingProperty(group *ast.GroupNode) *ast.GroupNode {
	for _, n := range group.Nodes {
		if onode, ok := n.(*ast.StringNode); ok {
			onode.Key = group.Key
		}
	}
	return group
}
