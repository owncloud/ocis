package bleve_test

import (
	"strings"
	"testing"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/bleve"
)

var FullAst = &ast.Ast{
	Nodes: []ast.Node{
		&ast.TextPropertyRestriction{Key: "tag", Value: "foo"},
		&ast.TextPropertyRestriction{Key: "tag", Value: "bar"},
		&ast.TextPropertyRestriction{Key: "name", Value: "book.pdf"},
		&ast.TextPropertyRestriction{Key: "content", Value: "ahab"},
	},
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
			got:  FullAst,
			want: []string{
				`tag:foo`,
				`tag:bar`,
				`name:book.pdf`,
				`content:ahab`,
			},
			err: false,
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
		if err := bleve.Compile(&strings.Builder{}, FullAst); err != nil {
			b.Fatal(err)
		}
	}
}
