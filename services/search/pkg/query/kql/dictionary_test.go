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
	// fixMe: ref3
	// `author:("John Smith" "Jane Smith")`,
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
					&ast.TextPropertyRestriction{Key: "author", Value: "John Smith"},
					&ast.TextPropertyRestriction{Key: "filetype", Value: "docx"},
					&ast.TextPropertyRestriction{Key: "filename", Value: "budget.xlsx"},
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
					&ast.TextPropertyRestriction{Key: "author", Value: "Shakespear"},
					&ast.TextPropertyRestriction{Key: "author", Value: "Paul"},
					&ast.TextPropertyRestriction{Key: "author", Value: "Shakesp*"},
					&ast.TextPropertyRestriction{Key: "title", Value: "Advanced Search"},
					&ast.TextPropertyRestriction{Key: "title", Value: "Advanced Sear*"},
					// fixMe: ref1
					// &ast.TextPropertyRestriction{Key: "title", Value: "Advan Search"},
					// fixMe: ref2
					// &ast.TextPropertyRestriction{Key: "title", Value: "anced Search"},
					&ast.TextPropertyRestriction{Key: "author", Value: "John Smith"},
					&ast.BooleanOperator{Value: "OR"},
					&ast.TextPropertyRestriction{Key: "author", Value: "Jane Smith"},
					&ast.TextPropertyRestriction{Key: "author", Value: "John Smith"},
					&ast.BooleanOperator{Value: "AND"},
					&ast.TextPropertyRestriction{Key: "filetype", Value: "docx"},
					// fixMe: ref3
					//&ast.Group{Nodes: []ast.Node{
					//	&ast.TextPropertyRestriction{Key: "author", Value: "John Smith"},
					//	&ast.TextPropertyRestriction{Key: "author", Value: "Jane Smith"},
					//}},
					// fixMe: ref4
					//&ast.Group{Nodes: []ast.Node{
					//	&ast.TextPropertyRestriction{Key: "DepartmentId", Value: "*"},
					//	&ast.BooleanOperator{Value: "OR"},
					//	&ast.TextPropertyRestriction{Key: "RelatedHubSites", Value: "*"},
					//}},
					//&ast.BooleanOperator{Value: "AND"},
					//&ast.TextPropertyRestriction{Key: "contentclass", Value: "sts_site"},
					//&ast.BooleanOperator{Value: "NOT"},
					//&ast.YesNoQuery{Key: "IsHubSite", Value: true},
					//`author:"John Smith" (filetype:docx title:"Advanced Search")`,
					&ast.TextPropertyRestriction{Key: "author", Value: "John Smith"},
					&ast.Group{Nodes: []ast.Node{
						&ast.TextPropertyRestriction{Key: "filetype", Value: "docx"},
						&ast.TextPropertyRestriction{Key: "title", Value: "Advanced Search"},
					}},
				},
			},
			err: false,
		},
		{
			name: "Conjunctive normal form",
			got: []string{
				`(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`,
			},
			want: &ast.Ast{
				Nodes: []ast.Node{
					&ast.Group{Nodes: []ast.Node{
						&ast.TextPropertyRestriction{Key: "name", Value: "moby di*"},
						&ast.BooleanOperator{Value: "OR"},
						&ast.TextPropertyRestriction{Key: "tag", Value: "bestseller"},
					}},
					&ast.BooleanOperator{Value: "AND"},
					&ast.TextPropertyRestriction{Key: "tag", Value: "book"},
					&ast.BooleanOperator{Value: "NOT"},
					&ast.TextPropertyRestriction{Key: "tag", Value: "read"},
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
