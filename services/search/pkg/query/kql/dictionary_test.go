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
	// fixMe: ref1
	// `title:"Advan* Search""`,
	// fixMe: ref2
	// `title:"*anced Search""`,
	`author:"John Smith" OR author:"Jane Smith"`,
	`author:"John Smith" AND filetype:docx`,
	`author:("John Smith" "Jane Smith")`,
	// fixMe: ref4
	// `(DepartmentId:* OR RelatedHubSites:*) AND contentclass:sts_site NOT IsHubSite:true`,
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
					&ast.Word{Value: "federated"},
					&ast.Word{Value: "search"},
					&ast.Word{Value: "federat*"},
					&ast.Word{Value: "search"},
					&ast.Word{Value: "search"},
					&ast.Word{Value: "fed*"},
					&ast.StringProperty{Key: "author", Value: "John Smith"},
					&ast.StringProperty{Key: "filetype", Value: "docx"},
					&ast.StringProperty{Key: "filename", Value: "budget.xlsx"},
					&ast.Word{Value: "author"},
					&ast.Phrase{Value: "John Smith"},
					&ast.Word{Value: "author"},
					&ast.Phrase{Value: "John Smith"},
					&ast.Word{Value: "author"},
					&ast.Phrase{Value: "John Smith"},
					&ast.Word{Value: "author"},
					&ast.Phrase{Value: "John Smith"},
					&ast.Word{Value: "author"},
					&ast.Phrase{Value: "John Smith"},
					&ast.StringProperty{Key: "author", Value: "Shakespear"},
					&ast.StringProperty{Key: "author", Value: "Paul"},
					&ast.StringProperty{Key: "author", Value: "Shakesp*"},
					&ast.StringProperty{Key: "title", Value: "Advanced Search"},
					&ast.StringProperty{Key: "title", Value: "Advanced Sear*"},
					// fixMe: ref1
					// &ast.StringProperty{Key: "title", Value: "Advan Search"},
					// fixMe: ref2
					// &ast.StringProperty{Key: "title", Value: "anced Search"},
					&ast.StringProperty{Key: "author", Value: "John Smith"},
					&ast.BooleanOperator{Value: "OR"},
					&ast.StringProperty{Key: "author", Value: "Jane Smith"},
					&ast.StringProperty{Key: "author", Value: "John Smith"},
					&ast.BooleanOperator{Value: "AND"},
					&ast.StringProperty{Key: "filetype", Value: "docx"},
					// fixMe: ref3
					&ast.KeyGroup{
						Key: "author",
						Nodes: []ast.Node{
							&ast.Phrase{Value: "John Smith"},
							&ast.Phrase{Value: "Jane Smith"},
						},
					},
					// fixMe: ref4
					//&ast.Group{Nodes: []ast.Node{
					//	&ast.StringProperty{Key: "DepartmentId", Value: "*"},
					//	&ast.BooleanOperator{Value: "OR"},
					//	&ast.StringProperty{Key: "RelatedHubSites", Value: "*"},
					//}},
					//&ast.BooleanOperator{Value: "AND"},
					//&ast.StringProperty{Key: "contentclass", Value: "sts_site"},
					//&ast.BooleanOperator{Value: "NOT"},
					//&ast.YesNoQuery{Key: "IsHubSite", Value: true},
					//`author:"John Smith" (filetype:docx title:"Advanced Search")`,
					&ast.StringProperty{Key: "author", Value: "John Smith"},
					&ast.Group{Nodes: []ast.Node{
						&ast.StringProperty{Key: "filetype", Value: "docx"},
						&ast.StringProperty{Key: "title", Value: "Advanced Search"},
					}},
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
					&ast.Group{Nodes: []ast.Node{
						&ast.StringProperty{Key: "name", Value: "moby di*"},
						&ast.BooleanOperator{Value: "OR"},
						&ast.StringProperty{Key: "tag", Value: "bestseller"},
					}},
					&ast.BooleanOperator{Value: "AND"},
					&ast.StringProperty{Key: "tag", Value: "book"},
					&ast.BooleanOperator{Value: "NOT"},
					&ast.StringProperty{Key: "tag", Value: "read"},
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
					&ast.KeyGroup{
						Key: "author",
						Nodes: []ast.Node{
							&ast.Phrase{Value: "John Smith"},
							&ast.Word{Value: "Jane"},
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
