package store

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
)

// Builds a unique file name from the given bundle
func buildFileNameFromBundle(bundle *proto.SettingsBundle) string {
	return buildFileNameFromBundleArgs(bundle.Extension, bundle.Key)
}

// Builds a unique file name from the given params
func buildFileNameFromBundleArgs(extension string, key string) string {
	return extension + "__" + key + ".json"
}
