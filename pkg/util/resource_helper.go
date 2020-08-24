package util

import "github.com/owncloud/ocis-settings/pkg/proto/v0"

const (
	ResourceIdAll = "all"
)

// IsResourceMatched checks if the `example` resource is an exact match or a subset of `definition`
func IsResourceMatched(definition, example *proto.Resource) bool {
	if definition.Type != example.Type {
		return false
	}
	return definition.Id == ResourceIdAll || definition.Id == example.Id
}
