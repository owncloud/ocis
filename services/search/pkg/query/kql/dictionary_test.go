package kql_test

import (
	"strings"
	"testing"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast/test"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

var FullDictionaryKql = []string{
	`author: "John Smith"`,                         // Phrase, Phrase
	`author :"John Smith"`,                         // Phrase, Phrase
	`author : "John Smith"`,                        // Phrase, Phrase
	`tags:foo AND tag:bar`,                         // TagQuery, Operator, TagQuery
	`name:book.pdf`,                                // NameQuery
	`content:letter.docx`,                          // ContentQuery
	`name:book.pdf (content:letter.docx tags:foo)`, // NameQuery, GROUP |> ContentQuery, TagQuery <|
}

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		got  []string
		want *ast.Ast
		err  bool
	}{
		{
			name: "FullDictionaryKql",
			got:  FullDictionaryKql,
			want: test.FullDictionaryAst,
			err:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := kql.NewAST(strings.NewReader(strings.Join(tt.got, " ")))
			if (err != nil) != tt.err {
				t.Fatalf("NewAST() error = %v, wantErr %v", err, tt.err)
			}

			if tt.err {
				return
			}
			if diff := test.DiffAst(
				tt.want, got); diff != "" {
				t.Fatalf("AST mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := kql.NewAST(strings.NewReader(strings.Join(FullDictionaryKql, " "))); err != nil {
			b.Fatal(err)
		}
	}
}
