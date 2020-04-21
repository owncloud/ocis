package store

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"strings"
)

// Builds a unique file name from the given bundle
func buildFileNameFromBundle(bundle *proto.SettingsBundle) string {
	return buildFileNameFromBundleArgs(bundle.Extension, bundle.Key)
}

// Builds a unique file name from the given params
func buildFileNameFromBundleArgs(extension string, key string) string {
	return extension + "__" + key + ".json"
}

// Extracts extension and key from the given fileName and builds a (minimalistic) bundle from it
func parseBundleFromFileName(fileName string) *proto.SettingsBundle {
	parts := strings.Split(strings.Replace(fileName, ".json", "", 1), "__")
	return &proto.SettingsBundle{
		Key:       parts[1],
		Extension: parts[0],
	}
}
