package kql_test

import (
	"strings"
	"testing"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast/test"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

var FullDictionary = []string{
	`federated search`,
	`federat* search`,
	`search fed*`,
	`author:"John Smith"`,
	`filetype:docx`,
	`filename:budget.xlsx`,
	`author: "John Smith"`,
	`author :"John Smith"`,
	`author : "John Smith"`,
	`author "John Smith"`,
	`author "John Smith"`,
	`author:Shakespear`,
	`author:Paul`,
	`author:Shakesp*`,
	`title:"Advanced Search"`,
	`title:"Advanced Sear*"`,
	`title:"Advan* Search"`,
	`title:"*anced Search"`,
	`author:"John Smith" OR author:"Jane Smith"`,
	`author:"John Smith" AND filetype:docx`,
	`author:("John Smith" "Jane Smith")`,
	`(DepartmentId:* OR RelatedHubSites:*) AND contentclass:sts_site NOT IsHubSite:false`,
	`author:"John Smith" (filetype:docx title:"Advanced Search")`,
}

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		got  []string
		want *ast.Ast
		err  bool
	}{
		{
			name: "FullDictionary",
			got:  FullDictionary,
			want: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "federated"},
					&ast.StringNode{Value: "search"},
					&ast.StringNode{Value: "federat*"},
					&ast.StringNode{Value: "search"},
					&ast.StringNode{Value: "search"},
					&ast.StringNode{Value: "fed*"},
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.StringNode{Key: "filetype", Value: "docx"},
					&ast.StringNode{Key: "filename", Value: "budget.xlsx"},
					&ast.StringNode{Value: "author"},
					&ast.StringNode{Value: "John Smith"},
					&ast.StringNode{Value: "author"},
					&ast.StringNode{Value: "John Smith"},
					&ast.StringNode{Value: "author"},
					&ast.StringNode{Value: "John Smith"},
					&ast.StringNode{Value: "author"},
					&ast.StringNode{Value: "John Smith"},
					&ast.StringNode{Value: "author"},
					&ast.StringNode{Value: "John Smith"},
					&ast.StringNode{Key: "author", Value: "Shakespear"},
					&ast.StringNode{Key: "author", Value: "Paul"},
					&ast.StringNode{Key: "author", Value: "Shakesp*"},
					&ast.StringNode{Key: "title", Value: "Advanced Search"},
					&ast.StringNode{Key: "title", Value: "Advanced Sear*"},
					&ast.StringNode{Key: "title", Value: "Advan* Search"},
					&ast.StringNode{Key: "title", Value: "*anced Search"},
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.OperatorNode{Value: "OR"},
					&ast.StringNode{Key: "author", Value: "Jane Smith"},
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "filetype", Value: "docx"},
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.StringNode{Value: "Jane Smith"},
						},
					},
					&ast.GroupNode{
						Nodes: []ast.Node{
							&ast.StringNode{Key: "DepartmentId", Value: "*"},
							&ast.OperatorNode{Value: "OR"},
							&ast.StringNode{Key: "RelatedHubSites", Value: "*"},
						},
					},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "contentclass", Value: "sts_site"},
					&ast.OperatorNode{Value: "NOT"},
					&ast.BooleanNode{Key: "IsHubSite", Value: false},
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.GroupNode{
						Nodes: []ast.Node{
							&ast.StringNode{Key: "filetype", Value: "docx"},
							&ast.StringNode{Key: "title", Value: "Advanced Search"},
						},
					},
				},
			},
			err: false,
		},
		{
			name: "Group",
			got: []string{
				`(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`,
			},
			want: &ast.Ast{
				Nodes: []ast.Node{
					&ast.GroupNode{
						Nodes: []ast.Node{
							&ast.StringNode{Key: "name", Value: "moby di*"},
							&ast.OperatorNode{Value: "OR"},
							&ast.StringNode{Key: "tag", Value: "bestseller"},
						},
					},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "tag", Value: "book"},
					&ast.OperatorNode{Value: "NOT"},
					&ast.StringNode{Key: "tag", Value: "read"},
				},
			},
			err: false,
		},
		{
			name: "KeyGroup",
			got: []string{
				`author:("John Smith" Jane)`,
			},
			want: &ast.Ast{
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
			err: false,
		},
		{
			name: "KeyGroup or key",
			got: []string{
				`author:("John Smith" Jane) author:"Jack" AND author:"Oggy"`,
			},
			want: &ast.Ast{
				Nodes: []ast.Node{
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: "AND"},
							&ast.StringNode{Value: "Jane"},
						},
					},
					&ast.OperatorNode{Value: "OR"},
					&ast.StringNode{Key: "author", Value: "Jack"},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "author", Value: "Oggy"},
				},
			},
			err: false,
		},
		{
			name: "KeyGroup",
			got: []string{
				`author:("John Smith" OR Jane)"`,
			},
			want: &ast.Ast{
				Nodes: []ast.Node{
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: "OR"},
							&ast.StringNode{Value: "Jane"},
						},
					},
				},
			},
			err: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			q := strings.Join(tt.got, " ")
			got, err := kql.NewAST(strings.NewReader(q))
			if (err != nil) != tt.err {
				t.Fatalf("NewAST() error = %v, wantErr %v", err, tt.err)
			}

			if tt.err {
				return
			}
			if diff := test.DiffAst(
				tt.want, got); diff != "" {
				t.Fatalf("AST mismatch \nquery: '%s' \n(-want +got): %s", q, diff)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := kql.NewAST(strings.NewReader(strings.Join(FullDictionary, " "))); err != nil {
			b.Fatal(err)
		}
	}
}
