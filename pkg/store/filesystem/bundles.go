// Package store implements the go-micro store interface
package store

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
	"io/ioutil"
	"path"
)

// ListByExtension returns all bundles in the mountPath folder belonging to the given extension
func (s Store) ListByExtension(extension string) ([]*proto.SettingsBundle, error) {
	bundlesFolder := s.buildFolderPathBundles()
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
func (s Store) ReadBundle(extension string, bundleKey string) (*proto.SettingsBundle, error) {
	if len(extension) < 1 || len(bundleKey) < 1 {
		s.Logger.Error().Msg("extension and bundleKey cannot be empty")
		return nil, gstatus.Error(codes.InvalidArgument, "Missing a required identifier attribute")
	}

	filePath := s.buildFilePathFromBundleArgs(extension, bundleKey)
	record := proto.SettingsBundle{}
	if err := s.parseRecordFromFile(&record, filePath); err != nil {
		return nil, err
	}

	s.Logger.Debug().Msgf("read contents from file: %v", filePath)
	return &record, nil
}

// Write writes the given record into a file within the mountPath
func (s Store) WriteBundle(record *proto.SettingsBundle) (*proto.SettingsBundle, error) {
	if len(record.Extension) < 1 || len(record.BundleKey) < 1 {
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
