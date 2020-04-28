package store

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"os"
	"path"
)

const folderNameBundles = "bundles"
const folderNameSettings = "settings"

// Builds the folder path for storing settings bundles
func buildFolderPathBundles(mountPath string) string {
	folderPath := path.Join(mountPath, folderNameBundles)
	ensureFolderExists(folderPath)
	return folderPath
}

// Builds a unique file name from the given bundle
func buildFilePathFromBundle(mountPath string, bundle *proto.SettingsBundle) string {
	return buildFilePathFromBundleArgs(mountPath, bundle.Extension, bundle.Key)
}

// Builds a unique file name from the given params
func buildFilePathFromBundleArgs(mountPath string, extension string, key string) string {
	extensionFolder := path.Join(mountPath, folderNameBundles, extension)
	if _, err := os.Stat(extensionFolder); os.IsNotExist(err) {
		_ = os.MkdirAll(extensionFolder, 0700)
	}
	return path.Join(extensionFolder, key + ".json")
}

// Checks if the given path is an existing folder and creates one if not existing
func ensureFolderExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, 0700)
	}
}
