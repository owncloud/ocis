// Package store implements the go-micro store interface
package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	olog "github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/owncloud/ocis-settings/pkg/settings"
)

var (
	// StoreName is the default name for the settings store
	StoreName     = "ocis-settings-store"
	managerName   = "filesystem"
	emptyKeyError = "key cannot be empty"
)

// Store interacts with the filesystem to manage settings information
type Store struct {
	mountPath string
	Logger    olog.Logger
}

// New creates a new store
func New(cfg *config.Config) settings.Manager {
	s := Store{}

	dest := path.Join(cfg.Storage.RootMountPath, StoreName)
	if _, err := os.Stat(dest); err != nil {
		s.Logger.Info().Msgf("creating container on %v", dest)
		err := os.MkdirAll(dest, 0700)
		if err != nil {
			s.Logger.Err(err).Msgf("providing container on %v", dest)
		}
	}

	s.mountPath = dest
	return &s
}

// ListByExtension returns all bundles in the mountPath folder belonging to the given extension
func (s Store) ListByExtension(extension string) ([]*proto.SettingsBundle, error) {
	bundlesFolder := buildFolderPathBundles(s.mountPath)
	extensionFolders, err := ioutil.ReadDir(bundlesFolder)
	if err != nil {
		s.Logger.Err(err).Msgf("error reading %v", bundlesFolder)
		return nil, err
	}

	s.Logger.Info().Msgf("listing bundles by extension %v", extension)
	var records []*proto.SettingsBundle
	for _, extensionFolder := range extensionFolders {
		extensionPath := path.Join(bundlesFolder, extensionFolder.Name())
		bundleFiles, err := ioutil.ReadDir(extensionPath)
		if err == nil {
			for _, bundleFile := range bundleFiles {
				record := proto.SettingsBundle{}
				err = s.parseRecordFromFile(&record, path.Join(extensionPath, bundleFile.Name()))
				if err == nil && (len(extension) == 0 || extension == record.Extension) {
					records = append(records, &record)
				}
			}
		} else {
			s.Logger.Err(err).Msgf("error reading %v", extensionPath)
		}
	}

	return records, nil
}

// Read tries to find a bundle by the given extension and key within the mountPath
func (s Store) Read(extension string, key string) (*proto.SettingsBundle, error) {
	if len(extension) < 1 || len(key) < 1 {
		s.Logger.Error().Msg("extension and key cannot be empty")
		return nil, fmt.Errorf(emptyKeyError)
	}

	filePath := buildFilePathFromBundleArgs(s.mountPath, extension, key)
	record := proto.SettingsBundle{}
	if err := s.parseRecordFromFile(&record, filePath); err != nil {
		return nil, err
	}

	s.Logger.Debug().Msgf("read contents from file: %v", filePath)
	return &record, nil
}

// Write writes the given record into a file within the mountPath
func (s Store) Write(record *proto.SettingsBundle) (*proto.SettingsBundle, error) {
	if len(record.Extension) < 1 || len(record.Key) < 1 {
		s.Logger.Error().Msg("extension and key cannot be empty")
		return nil, fmt.Errorf(emptyKeyError)
	}

	filePath := buildFilePathFromBundle(s.mountPath, record)
	if err := s.writeRecordToFile(record, filePath); err != nil {
		return nil, err
	}

	s.Logger.Debug().Msgf("request contents written to file: %v", filePath)
	return record, nil
}

func init() {
	settings.Registry[managerName] = New
}
