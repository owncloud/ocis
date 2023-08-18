package bleve_test

import (
	"strings"
	"testing"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast/test"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/bleve"
)

var FullDictionaryBleve = []string{
	`tag:foo`,
	`tag:bar`,
	`name:book.pdf`,
	`content:letter.docx`,
	`name:book.pdf`,
}

func TestCompile(t *testing.T) {
	tests := []struct {
		name string
		got  *ast.Ast
		want []string
		err  bool
	}{
		{
			name: "FullDictionaryBleve",
			got:  test.FullDictionaryAst,
			want: FullDictionaryBleve,
			err:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var got strings.Builder
			err := bleve.Compile(&got, tt.got)
			if (err != nil) != tt.err {
				t.Fatalf("Compile() error = %v, wantErr %v", err, tt.err)
			}

			if tt.err {
				return
			}

			if got.String() != strings.Join(tt.want, " ") {
				t.Fatalf("Compile mismatch \ngot: `%s` \nwant: `%s`", got.String(), strings.Join(tt.want, " "))
			}
		})
	}
}

func BenchmarkCompile(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if err := bleve.Compile(&strings.Builder{}, test.FullDictionaryAst); err != nil {
			b.Fatal(err)
		}
	}
}
