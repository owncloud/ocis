package store

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"os"
	"path"
)

const folderNameBundles = "bundles"
const folderNameValues = "values"

// Builds the folder path for storing settings bundles
func (s Store) buildFolderPathBundles() string {
	folderPath := path.Join(s.mountPath, folderNameBundles)
	s.ensureFolderExists(folderPath)
	return folderPath
}

// Builds a unique file name from the given settings bundle
func (s Store) buildFilePathFromBundle(bundle *proto.SettingsBundle) string {
	return s.buildFilePathFromBundleArgs(bundle.Extension, bundle.BundleKey)
}

// Builds a unique file name from the given params
func (s Store) buildFilePathFromBundleArgs(extension string, bundleKey string) string {
	extensionFolder := path.Join(s.mountPath, folderNameBundles, extension)
	s.ensureFolderExists(extensionFolder)
	return path.Join(extensionFolder, bundleKey+".json")
}

// Builds the folder path for storing settings values
func (s Store) buildFolderPathValues() string {
	folderPath := path.Join(s.mountPath, folderNameValues)
	s.ensureFolderExists(folderPath)
	return folderPath
}

// Builds a unique file name from the given settings value
func (s Store) buildFilePathFromValue(value *proto.SettingsValue) string {
	return s.buildFilePathFromValueArgs(value.AccountUuid, value.Extension, value.BundleKey)
}

// Builds a unique file name from the given params
func (s Store) buildFilePathFromValueArgs(accountUuid string, extension string, bundleKey string) string {
	extensionFolder := path.Join(s.mountPath, folderNameValues, accountUuid, extension)
	s.ensureFolderExists(extensionFolder)
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
