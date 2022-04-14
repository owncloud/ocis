package util

import (
	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
)

const (
	// ResourceIDAll declares on a resource that it matches any id
	ResourceIDAll = "all"
)

// IsResourceMatched checks if the `example` resource is an exact match or a subset of `definition`
func IsResourceMatched(definition, example *settingsmsg.Resource) bool {
	if definition.Type != example.Type {
		return false
	}
	return definition.Id == ResourceIDAll || definition.Id == example.Id
}
