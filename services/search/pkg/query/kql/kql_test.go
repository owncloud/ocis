package kql_test

import (
	"testing"

	tAssert "github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

func TestNewAST(t *testing.T) {
	tests := []struct {
		name        string
		givenQuery  string
		shouldError bool
	}{
		{
			name:       "success",
			givenQuery: "foo:bar",
		},
		{
			name:        "error",
			givenQuery:  "AND",
			shouldError: true,
		},
	}

	assert := tAssert.New(t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := kql.Builder{}.Build(tt.givenQuery)

			if tt.shouldError {
				assert.NotNil(err)
				assert.Nil(got)
			} else {
				assert.Nil(err)
				assert.NotNil(got)
			}
		})
	}
}
