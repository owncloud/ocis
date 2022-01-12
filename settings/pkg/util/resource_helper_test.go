package util

import (
	"testing"

	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v1"
	"gotest.tools/v3/assert"
)

func TestIsResourceMatched(t *testing.T) {
	scenarios := []struct {
		name       string
		definition *settingsmsg.Resource
		example    *settingsmsg.Resource
		matched    bool
	}{
		{
			"same resource types without ids match",
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_SYSTEM,
			},
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_SYSTEM,
			},
			true,
		},
		{
			"different resource types without ids don't match",
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_SYSTEM,
			},
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
			},
			false,
		},
		{
			"same resource types with different ids don't match",
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
				Id:   "einstein",
			},
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
				Id:   "marie",
			},
			false,
		},
		{
			"same resource types with same ids match",
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
				Id:   "einstein",
			},
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
				Id:   "einstein",
			},
			true,
		},
		{
			"same resource types with definition = ALL and without id in example is a match",
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
				Id:   ResourceIDAll,
			},
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
			},
			true,
		},
		{
			"same resource types with definition.id = ALL and with some id in example is a match",
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
				Id:   ResourceIDAll,
			},
			&settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
				Id:   "einstein",
			},
			true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			assert.Equal(t, scenario.matched, IsResourceMatched(scenario.definition, scenario.example))
		})
	}
}
