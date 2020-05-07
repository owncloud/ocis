// Package store implements the go-micro store interface
package store

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
	"io/ioutil"
	"path"
)

// ListBundles returns all bundles in the mountPath folder belonging to the given extension
func (s Store) ListBundles(identifier *proto.Identifier) ([]*proto.SettingsBundle, error) {
	bundlesFolder := s.buildFolderPathBundles()
	extensionFolders, err := ioutil.ReadDir(bundlesFolder)
	if err != nil {
		s.Logger.Err(err).Msgf("error reading %v", bundlesFolder)
		return nil, err
	}

	if len(identifier.Extension) < 1 {
		s.Logger.Info().Msg("listing all bundles")
	} else {
		s.Logger.Info().Msgf("listing bundles by extension %v", identifier.Extension)
	}
	var records []*proto.SettingsBundle
	for _, extensionFolder := range extensionFolders {
		extensionPath := path.Join(bundlesFolder, extensionFolder.Name())
		bundleFiles, err := ioutil.ReadDir(extensionPath)
		if err == nil {
			for _, bundleFile := range bundleFiles {
				record := proto.SettingsBundle{}
				bundlePath := path.Join(extensionPath, bundleFile.Name())
				err = s.parseRecordFromFile(&record, bundlePath)
				if err != nil {
					s.Logger.Warn().Msgf("error reading %v", bundlePath)
					continue
				}
				if len(identifier.Extension) == 0 || identifier.Extension == record.Identifier.Extension {
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
func (s Store) ReadBundle(identifier *proto.Identifier) (*proto.SettingsBundle, error) {
	if len(identifier.Extension) < 1 || len(identifier.BundleKey) < 1 {
		s.Logger.Error().Msg("extension and bundleKey cannot be empty")
		return nil, gstatus.Error(codes.InvalidArgument, "Missing a required identifier attribute")
	}

	filePath := s.buildFilePathFromBundleArgs(identifier.Extension, identifier.BundleKey)
	record := proto.SettingsBundle{}
	if err := s.parseRecordFromFile(&record, filePath); err != nil {
		return nil, err
	}

	s.Logger.Debug().Msgf("read contents from file: %v", filePath)
	return &record, nil
}

// Write writes the given record into a file within the mountPath
func (s Store) WriteBundle(record *proto.SettingsBundle) (*proto.SettingsBundle, error) {
	if len(record.Identifier.Extension) < 1 || len(record.Identifier.BundleKey) < 1 {
		s.Logger.Error().Msg("extension and bundleKey cannot be empty")
		return nil, gstatus.Error(codes.InvalidArgument, "Missing a required identifier attribute")
	}

	filePath := s.buildFilePathFromBundle(record)
	if err := s.writeRecordToFile(record, filePath); err != nil {
		return nil, err
	}

	s.Logger.Debug().Msgf("request contents written to file: %v", filePath)
	return record, nil
}
