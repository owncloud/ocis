// Package store implements the go-micro store interface
package store

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
	"os"
	"path"
	"path/filepath"
)

// Tries to find a value by the given identifier attributes within the mountPath
func (s Store) ReadValue(identifier *proto.Identifier) (*proto.SettingsValue, error) {
	if len(identifier.AccountUuid) < 1 || len(identifier.Extension) < 1 || len(identifier.BundleKey) < 1 || len(identifier.SettingKey) < 1 {
		s.Logger.Error().Msg("account-uuid, extension, bundle and setting are required")
		return nil, gstatus.Errorf(codes.InvalidArgument, "Missing a required identifier attribute")
	}

	filePath := s.buildFilePathFromValueArgs(identifier.AccountUuid, identifier.Extension, identifier.BundleKey)
	values, err := s.readValuesMapFromFile(filePath)
	if err != nil {
		return nil, err
	}
	if value := values.Values[identifier.SettingKey]; value != nil {
		return value, nil
	}
	// TODO: we want to return sensible defaults here, when the value was not found
	return nil, gstatus.Error(codes.NotFound, "SettingsValue not set")
}

// Writes the given SettingsValue into a file within the mountPath
func (s Store) WriteValue(value *proto.SettingsValue) (*proto.SettingsValue, error) {
	if len(value.Identifier.AccountUuid) < 1 || len(value.Identifier.Extension) < 1 || len(value.Identifier.BundleKey) < 1 || len(value.Identifier.SettingKey) < 1 {
		s.Logger.Error().Msg("all identifier keys are required")
		return nil, gstatus.Errorf(codes.InvalidArgument, "Missing a required identifier attribute")
	}

	filePath := s.buildFilePathFromValue(value)
	values, err := s.readValuesMapFromFile(filePath)
	if err != nil {
		return nil, err
	}
	values.Values[value.Identifier.SettingKey] = value
	if err := s.writeRecordToFile(values, filePath); err != nil {
		return nil, err
	}
	return value, nil
}

// Reads all values within the scope of the given identifier
func (s Store) ListValues(identifier *proto.Identifier) ([]*proto.SettingsValue, error) {
	if len(identifier.AccountUuid) < 1 {
		s.Logger.Error().Msg("account-uuid is required")
		return nil, gstatus.Errorf(codes.InvalidArgument, "Missing a required identifier attribute")
	}

	accountFolderPath := path.Join(s.mountPath, folderNameValues, identifier.AccountUuid)
	var values []*proto.SettingsValue
	if _, err := os.Stat(accountFolderPath); err != nil {
		return values, nil
	}

	// TODO: might be a good idea to do this non-hierarchical. i.e. allowing all fragments in the identifier being set or not.
	// depending on the set values in the identifier arg, collect all SettingValues files for the account
	var valueFilePaths []string
	if len(identifier.Extension) < 1 {
		if err := filepath.Walk(accountFolderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			valueFilePaths = append(valueFilePaths, path)
			return nil
		}); err != nil {
			s.Logger.Err(err).Msgf("error reading %v", accountFolderPath)
			return nil, err
		}
	} else if len(identifier.BundleKey) < 1 {
		extensionPath := path.Join(accountFolderPath, identifier.Extension)
		if err := filepath.Walk(extensionPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			valueFilePaths = append(valueFilePaths, path)
			return nil
		}); err != nil {
			s.Logger.Err(err).Msgf("error reading %v", extensionPath)
			return nil, err
		}
	} else {
		bundlePath := path.Join(accountFolderPath, identifier.Extension, identifier.BundleKey+".json")
		valueFilePaths = append(valueFilePaths, bundlePath)
	}

	// parse the SettingValues from the collected files
	for _, filePath := range valueFilePaths {
		bundleValues, err := s.readValuesMapFromFile(filePath)
		if err != nil {
			s.Logger.Err(err).Msgf("error reading %v", filePath)
		} else {
			for _, value := range bundleValues.Values {
				values = append(values, value)
			}
		}
	}
	return values, nil
}

// Reads SettingsValues as map from the given file or returns an empty map if the file doesn't exist.
func (s Store) readValuesMapFromFile(filePath string) (*proto.SettingsValues, error) {
	values := &proto.SettingsValues{}
	err := s.parseRecordFromFile(values, filePath)
	if err != nil {
		if os.IsNotExist(err) {
			values.Values = map[string]*proto.SettingsValue{}
		} else {
			return nil, err
		}
	}
	return values, nil
}
