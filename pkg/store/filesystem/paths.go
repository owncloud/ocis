package store

import (
	"os"
	"path"

	"github.com/owncloud/ocis-settings/pkg/proto/v0"
)

const folderNameBundles = "bundles"
const folderNameValues = "values"

// Builds the folder path for storing settings bundles. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFolderPathBundles(mkdir bool) string {
	folderPath := path.Join(s.mountPath, folderNameBundles)
	if mkdir {
		s.ensureFolderExists(folderPath)
	}
	return folderPath
}

// Builds a unique file name from the given settings bundle. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFilePathFromBundle(bundle *proto.SettingsBundle, mkdir bool) string {
	return s.buildFilePathFromBundleArgs(bundle.Identifier.Extension, bundle.Identifier.BundleKey, mkdir)
}

// Builds a unique file name from the given params. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFilePathFromBundleArgs(extension string, bundleKey string, mkdir bool) string {
	extensionFolder := path.Join(s.mountPath, folderNameBundles, extension)
	if mkdir {
		s.ensureFolderExists(extensionFolder)
	}
	return path.Join(extensionFolder, bundleKey+".json")
}

// Builds a unique file name from the given settings value. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFilePathFromValue(value *proto.SettingsValue, mkdir bool) string {
	return s.buildFilePathFromValueArgs(value.Identifier.AccountUuid, value.Identifier.Extension, value.Identifier.BundleKey, mkdir)
}

// Builds a unique file name from the given params. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFilePathFromValueArgs(accountUUID string, extension string, bundleKey string, mkdir bool) string {
	extensionFolder := path.Join(s.mountPath, folderNameValues, accountUUID, extension)
	if mkdir {
		s.ensureFolderExists(extensionFolder)
	}
	return path.Join(extensionFolder, bundleKey+".json")
}

// Checks if the given path is an existing folder and creates one if not existing
func (s Store) ensureFolderExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			s.Logger.Err(err).Msgf("Error creating folder %v", path)
		}
	}
}
