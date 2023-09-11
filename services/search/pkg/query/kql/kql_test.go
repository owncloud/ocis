package kql_test

import (
	"testing"

	tAssert "github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

func TestNewAST(t *testing.T) {
	tests := []struct {
		name          string
		givenQuery    string
		expectedError error
	}{
		{
			name:       "success",
			givenQuery: "foo:bar",
		},
		{
			name:       "error",
			givenQuery: kql.BoolAND,
			expectedError: kql.StartsWithBinaryOperatorError{
				Node: &ast.OperatorNode{Value: kql.BoolAND},
			},
		},
	}

	assert := tAssert.New(t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := kql.Builder{}.Build(tt.givenQuery)

			if tt.expectedError != nil {
				if tt.expectedError.Error() != "" {
					assert.Equal(err.Error(), tt.expectedError.Error())
				} else {
					assert.NotNil(err)
				}

				assert.Nil(got)

				return
			}

			assert.Nil(err)
			assert.NotNil(got)
		})
	}
}
