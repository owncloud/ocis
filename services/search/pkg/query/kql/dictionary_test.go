package kql_test

import (
	"strings"
	"testing"
	"time"

	tAssert "github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast/test"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

var timeMustParse = func(t *testing.T, ts string) time.Time {
	tp, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		t.Fatalf("time.Parse(...) error = %v", err)
	}

	return tp
}

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
	`author:("John Smith" OR "Jane Smith")`,
	`(DepartmentId:* OR RelatedHubSites:*) AND contentclass:sts_site NOT IsHubSite:false`,
	`author:"John Smith" (filetype:docx title:"Advanced Search")`,
}

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		givenQuery    []string
		expectedAst   *ast.Ast
		expectedError error
	}{
		{
			name:       "FullDictionary",
			givenQuery: FullDictionary,
			expectedAst: &ast.Ast{
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
					&ast.StringNode{Key: "author", Value: "Jane Smith"},
					&ast.OperatorNode{Value: "OR"},
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.StringNode{Key: "filetype", Value: "docx"},
					&ast.OperatorNode{Value: "AND"},
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: "AND"},
							&ast.StringNode{Value: "Jane Smith"},
						},
					},
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: "OR"},
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
		},
		{
			name: "Group",
			givenQuery: []string{
				`(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`,
				`author:("John Smith" Jane)`,
				`author:("John Smith" OR Jane)`,
			},
			expectedAst: &ast.Ast{
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
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.StringNode{Value: "Jane"},
						},
					},
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
		},
		{
			name: "KeyGroup or key conjunction",
			givenQuery: []string{
				`author:("John Smith" Jane) author:"Jack" AND author:"Oggy"`,
			},
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.StringNode{Value: "Jane"},
						},
					},
					&ast.StringNode{Key: "author", Value: "Jack"},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "author", Value: "Oggy"},
				},
			},
		},
		{
			name: "KeyGroup",
			givenQuery: []string{
				`author:("John Smith" OR Jane)`,
			},
			expectedAst: &ast.Ast{
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
		},
		{
			name: "not and not",
			givenQuery: []string{
				`NOT "John Smith" NOT Jane`,
			},
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.OperatorNode{Value: "NOT"},
					&ast.StringNode{Value: "John Smith"},
					&ast.OperatorNode{Value: "NOT"},
					&ast.StringNode{Value: "Jane"},
				},
			},
		},
		{
			name: "not or not and not",
			givenQuery: []string{
				`NOT author:"John Smith" NOT author:"Jane Smith" NOT tag:sifi`,
			},
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.OperatorNode{Value: "NOT"},
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.OperatorNode{Value: "NOT"},
					&ast.StringNode{Key: "author", Value: "Jane Smith"},
					&ast.OperatorNode{Value: "NOT"},
					&ast.StringNode{Key: "tag", Value: "sifi"},
				},
			},
		},
		{
			name: "misc",
			givenQuery: []string{
				`scope:"<uuid>/new folder/subfolder" file`,
			},
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{
						Key:   "scope",
						Value: "<uuid>/new folder/subfolder",
					},
					&ast.StringNode{
						Value: "file",
					},
				},
			},
		},
		{
			name: "unicode",
			givenQuery: []string{
				`	😂 "*😀 😁*" name:😂💁👌🎍😍 name:😂💁👌 😍`,
			},
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{
						Value: "😂",
					},
					&ast.StringNode{
						Value: "*😀 😁*",
					},
					&ast.StringNode{
						Key:   "name",
						Value: "😂💁👌🎍😍",
					},
					&ast.StringNode{
						Key:   "name",
						Value: "😂💁👌",
					},
					&ast.StringNode{
						Value: "😍",
					},
				},
			},
		},
		{
			name: "DateTimeRestrictionNode",
			givenQuery: []string{
				`Mtime:"2023-09-05T08:42:11.23554+02:00"`,
				`Mtime:2023-09-05T08:42:11.23554+02:00`,
				`Mtime="2023-09-05T08:42:11.23554+02:00"`,
				`Mtime=2023-09-05T08:42:11.23554+02:00`,
				`Mtime<"2023-09-05T08:42:11.23554+02:00"`,
				`Mtime<2023-09-05T08:42:11.23554+02:00`,
				`Mtime<="2023-09-05T08:42:11.23554+02:00"`,
				`Mtime<=2023-09-05T08:42:11.23554+02:00`,
				`Mtime>"2023-09-05T08:42:11.23554+02:00"`,
				`Mtime>2023-09-05T08:42:11.23554+02:00`,
				`Mtime>="2023-09-05T08:42:11.23554+02:00"`,
				`Mtime>=2023-09-05T08:42:11.23554+02:00`,
			},
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: ":"},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: ":"},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: "="},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: "="},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: "<"},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: "<"},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: "<="},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: "<="},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: ">"},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: ">"},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: ">="},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
					&ast.DateTimeNode{
						Key:      "Mtime",
						Operator: &ast.OperatorNode{Value: ">="},
						Value:    timeMustParse(t, "2023-09-05T08:42:11.23554+02:00"),
					},
				},
			},
		},
		{
			name: "id",
			givenQuery: []string{
				`id:b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c`,
				`ID:b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c`,
			},
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{
						Key:   "id",
						Value: "b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c",
					},
					&ast.StringNode{
						Key:   "ID",
						Value: "b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c",
					},
				},
			},
		},
	}

	assert := tAssert.New(t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			q := strings.Join(tt.givenQuery, " ")

			parsedAST, err := kql.Parse("", []byte(q))

			if tt.expectedError != nil {
				assert.Equal(err, tt.expectedError)
				assert.Nil(parsedAST)

				return
			}

			normalizedNodes, err := kql.NormalizeNodes(tt.expectedAst.Nodes)
			if err != nil {
				t.Fatalf("NormalizeNodes() error = %v", err)
			}
			tt.expectedAst.Nodes = normalizedNodes

			if diff := test.DiffAst(tt.expectedAst, parsedAST); diff != "" {
				t.Fatalf("AST mismatch \nquery: '%s' \n(-want +got): %s", q, diff)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := kql.Parse("", []byte(strings.Join(FullDictionary, " "))); err != nil {
			b.Fatal(err)
		}
	}
}
