package kql_test

import (
	"testing"
	"time"

	tAssert "github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast/test"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

var now = time.Now()

func TestNormalizeNodes(t *testing.T) {
	tests := []struct {
		name          string
		givenNodes    []ast.Node
		expectedNodes []ast.Node
		fixme         bool
		expectedError error
	}{
		{
			name: "start with binary operator",
			givenNodes: []ast.Node{
				&ast.OperatorNode{Value: "OR"},
			},
			expectedError: &kql.StartsWithBinaryOperatorError{Op: "OR"},
		},
		{
			name: "same key implicit OR",
			givenNodes: []ast.Node{
				&ast.StringNode{Key: "author", Value: "John Smith"},
				&ast.StringNode{Key: "author", Value: "Jane Smith"},
			},
			expectedNodes: []ast.Node{
				&ast.StringNode{Key: "author", Value: "John Smith"},
				&ast.OperatorNode{Value: "OR"},
				&ast.StringNode{Key: "author", Value: "Jane Smith"},
			},
		},
		{
			name: "no key implicit AND",
			givenNodes: []ast.Node{
				&ast.StringNode{Value: "John Smith"},
				&ast.StringNode{Value: "Jane Smith"},
			},
			expectedNodes: []ast.Node{
				&ast.StringNode{Value: "John Smith"},
				&ast.OperatorNode{Value: "AND"},
				&ast.StringNode{Value: "Jane Smith"},
			},
		},
		{
			name: "same key explicit AND",
			givenNodes: []ast.Node{
				&ast.StringNode{Key: "author", Value: "John Smith"},
				&ast.OperatorNode{Value: "AND"},
				&ast.StringNode{Key: "author", Value: "Jane Smith"},
			},
			expectedNodes: []ast.Node{
				&ast.StringNode{Key: "author", Value: "John Smith"},
				&ast.OperatorNode{Value: "AND"},
				&ast.StringNode{Key: "author", Value: "Jane Smith"},
			},
		},
		{
			name: "key-group implicit AND",
			// https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference#grouping-property-restrictions-within-a-kql-query
			fixme: true,
			givenNodes: []ast.Node{
				&ast.GroupNode{Key: "author", Nodes: []ast.Node{
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.StringNode{Key: "author", Value: "Jane Smith"},
				}},
			},
			expectedNodes: []ast.Node{
				&ast.GroupNode{Key: "author", Nodes: []ast.Node{
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.OperatorNode{Value: "AND"},
					&ast.StringNode{Key: "author", Value: "Jane Smith"},
				}},
			},
		},
		{
			name: "different key implicit AND",
			givenNodes: []ast.Node{
				&ast.StringNode{Key: "author", Value: "John Smith"},
				&ast.StringNode{Key: "filetype", Value: "docx"},
				&ast.DateTimeNode{Key: "mtime", Operator: &ast.OperatorNode{Value: "="}, Value: now},
			},
			expectedNodes: []ast.Node{
				&ast.StringNode{Key: "author", Value: "John Smith"},
				&ast.OperatorNode{Value: "AND"},
				&ast.StringNode{Key: "filetype", Value: "docx"},
				&ast.OperatorNode{Value: "AND"},
				&ast.DateTimeNode{Key: "mtime", Operator: &ast.OperatorNode{Value: "="}, Value: now},
			},
		},
	}

	assert := tAssert.New(t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.fixme {
				t.Skip("not implemented")
			}

			normalizedNodes, err := kql.NormalizeNodes(tt.givenNodes)

			if tt.expectedError != nil {
				assert.Equal(err, tt.expectedError)
				assert.Nil(normalizedNodes)

				return
			}

			if diff := test.DiffAst(tt.expectedNodes, normalizedNodes); diff != "" {
				t.Fatalf("Nodes mismatch (-want +got): %s", diff)
			}
		})
	}
}
