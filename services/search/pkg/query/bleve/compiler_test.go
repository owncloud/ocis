package bleve

import (
	"testing"

	"github.com/blevesearch/bleve/v2/search/query"
	tAssert "github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

func Test_compile(t *testing.T) {
	tests := []struct {
		name    string
		args    *ast.Ast
		want    query.Query
		wantErr bool
	}{
		{
			name: `federated`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "federated"},
				},
			},
			want: query.NewConjunctionQuery([]query.Query{
				query.NewQueryStringQuery(`Name:federated`),
			}),
			wantErr: false,
		},
		{
			name: `"John Smith"`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "John Smith"},
				},
			},
			want: query.NewConjunctionQuery([]query.Query{
				query.NewQueryStringQuery(`Name:John\ Smith`),
			}),
			wantErr: false,
		},
		{
			name: `"John Smith" Jane`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Key: "name", Value: "John Smith"},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "name", Value: "Jane"},
				},
			},
			want: query.NewConjunctionQuery([]query.Query{
				query.NewQueryStringQuery(`Name:John\ Smith`),
				query.NewQueryStringQuery(`Name:Jane`),
			}),
			wantErr: false,
		},
		{
			name: `tag:bestseller tag:book`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Key: "tag", Value: "bestseller"},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "tag", Value: "book"},
				},
			},
			want: query.NewConjunctionQuery([]query.Query{
				query.NewQueryStringQuery(`Tags:bestseller`),
				query.NewQueryStringQuery(`Tags:book`),
			}),
			wantErr: false,
		},
		{
			name: `name:"moby di*" OR tag:bestseller AND tag:book`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Key: "name", Value: "moby di*"},
					&ast.OperatorNode{Value: "OR"},
					&ast.StringNode{Key: "tag", Value: "bestseller"},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "tag", Value: "book"},
				},
			},
			want: query.NewDisjunctionQuery([]query.Query{
				query.NewQueryStringQuery(`Name:moby\ di*`),
				query.NewConjunctionQuery([]query.Query{
					query.NewQueryStringQuery(`Tags:bestseller`),
					query.NewQueryStringQuery(`Tags:book`),
				}),
			}),
			wantErr: false,
		},
		{
			name: `a AND b OR c`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "a"},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Value: "b"},
					&ast.OperatorNode{Value: "OR"},
					&ast.StringNode{Value: "c"},
				},
			},
			want: query.NewDisjunctionQuery([]query.Query{
				query.NewConjunctionQuery([]query.Query{
					query.NewQueryStringQuery(`Name:a`),
					query.NewQueryStringQuery(`Name:b`),
				}),
				query.NewQueryStringQuery(`Name:c`),
			}),
			wantErr: false,
		},
		{
			name: `(name:"moby di*" OR tag:bestseller) AND tag:book`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.GroupNode{Nodes: []ast.Node{
						&ast.StringNode{Key: "name", Value: "moby di*"},
						&ast.OperatorNode{Value: "OR"},
						&ast.StringNode{Key: "tag", Value: "bestseller"},
					}},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "tag", Value: "book"},
				},
			},
			want: query.NewConjunctionQuery([]query.Query{
				query.NewDisjunctionQuery([]query.Query{
					query.NewQueryStringQuery(`Name:moby\ di*`),
					query.NewQueryStringQuery(`Tags:bestseller`),
				}),
				query.NewQueryStringQuery(`Tags:book`),
			}),
			wantErr: false,
		},
		{
			name: `(name:"moby di*" OR tag:bestseller) AND tag:book AND NOT tag:read`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.GroupNode{Nodes: []ast.Node{
						&ast.StringNode{Key: "name", Value: "moby di*"},
						&ast.OperatorNode{Value: "OR"},
						&ast.StringNode{Key: "tag", Value: "bestseller"},
					}},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "tag", Value: "book"},
					&ast.OperatorNode{Value: "AND"},
					&ast.OperatorNode{Value: "NOT"},
					&ast.StringNode{Key: "tag", Value: "read"},
				},
			},
			want: query.NewConjunctionQuery([]query.Query{
				query.NewDisjunctionQuery([]query.Query{
					query.NewQueryStringQuery(`Name:moby\ di*`),
					query.NewQueryStringQuery(`Tags:bestseller`),
				}),
				query.NewQueryStringQuery(`Tags:book`),
				query.NewBooleanQuery(nil, nil, []query.Query{query.NewQueryStringQuery(`Tags:read`)}),
			}),
			wantErr: false,
		},
		{
			name: `author:("John Smith" Jane)`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: "AND"},
							&ast.StringNode{Value: "Jane"},
						},
					},
				},
			},
			want: query.NewConjunctionQuery([]query.Query{
				query.NewQueryStringQuery(`author:John\ Smith`),
				query.NewQueryStringQuery(`author:Jane`),
			}),
			wantErr: false,
		},
		{
			name: `author:("John Smith" Jane) AND tag:bestseller`,
			args: &ast.Ast{
				Nodes: []ast.Node{
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: "AND"},
							&ast.StringNode{Value: "Jane"},
						},
					},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "tag", Value: "bestseller"},
				},
			},
			want: query.NewConjunctionQuery([]query.Query{
				query.NewQueryStringQuery(`author:John\ Smith`),
				query.NewQueryStringQuery(`author:Jane`),
				query.NewQueryStringQuery(`Tags:bestseller`),
			}),
			wantErr: false,
		},
	}

	assert := tAssert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compile(tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(tt.want, got)
		})
	}
}
