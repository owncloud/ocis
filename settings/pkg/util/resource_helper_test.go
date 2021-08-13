package util

import (
	"testing"

	"github.com/owncloud/ocis/settings/pkg/proto/v0"
	"gotest.tools/v3/assert"
)

func TestIsResourceMatched(t *testing.T) {
	scenarios := []struct {
		name       string
		definition *proto.Resource
		example    *proto.Resource
		matched    bool
	}{
		{
			"same resource types without ids match",
			&proto.Resource{
				Type: proto.Resource_TYPE_SYSTEM,
			},
			&proto.Resource{
				Type: proto.Resource_TYPE_SYSTEM,
			},
			true,
		},
		{
			"different resource types without ids don't match",
			&proto.Resource{
				Type: proto.Resource_TYPE_SYSTEM,
			},
			&proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			false,
		},
		{
			"same resource types with different ids don't match",
			&proto.Resource{
				Type: proto.Resource_TYPE_USER,
				Id:   "einstein",
			},
			&proto.Resource{
				Type: proto.Resource_TYPE_USER,
				Id:   "marie",
			},
			false,
		},
		{
			"same resource types with same ids match",
			&proto.Resource{
				Type: proto.Resource_TYPE_USER,
				Id:   "einstein",
			},
			&proto.Resource{
				Type: proto.Resource_TYPE_USER,
				Id:   "einstein",
			},
			true,
		},
		{
			"same resource types with definition = ALL and without id in example is a match",
			&proto.Resource{
				Type: proto.Resource_TYPE_USER,
				Id:   ResourceIDAll,
			},
			&proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			true,
		},
		{
			"same resource types with definition.id = ALL and with some id in example is a match",
			&proto.Resource{
				Type: proto.Resource_TYPE_USER,
				Id:   ResourceIDAll,
			},
			&proto.Resource{
				Type: proto.Resource_TYPE_USER,
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
