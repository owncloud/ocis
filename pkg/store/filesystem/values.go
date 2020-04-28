// Package store implements the go-micro store interface
package store

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
	"os"
)

// Read tries to find a value by the given identifier attributes within the mountPath
func (s Store) ReadValue(accountUuid string, extension string, bundleKey string, settingKey string) (*proto.SettingsValue, error) {
	if len(accountUuid) < 1 || len(extension) < 1 || len(bundleKey) < 1 || len(settingKey) < 1 {
		s.Logger.Error().Msg("account, extension, bundle and setting are required")
		return nil, gstatus.Errorf(codes.InvalidArgument, "Missing a required identifier attribute")
	}

	filePath := s.buildFilePathFromValueArgs(accountUuid, extension, bundleKey)
	values, err := s.readValuesMapFromFile(filePath)
	if err != nil {
		return nil, err
	}
	if value := values.Values[settingKey]; value != nil {
		return value, nil
	}
	// TODO: we want to return sensible defaults here, when the value was not found
	return nil, gstatus.Error(codes.NotFound, "SettingsValue not set")
}

// Write writes the given SettingsValue into a file within the mountPath
func (s Store) WriteValue(value *proto.SettingsValue) (*proto.SettingsValue, error) {
	if len(value.AccountUuid) < 1 || len(value.Extension) < 1 || len(value.BundleKey) < 1 || len(value.SettingKey) < 1 {
		s.Logger.Error().Msg("account, extension, bundle and setting are required")
		return nil, gstatus.Errorf(codes.InvalidArgument, "Missing a required identifier attribute")
	}

	filePath := s.buildFilePathFromValue(value)
	values, err := s.readValuesMapFromFile(filePath)
	if err != nil {
		return nil, err
	}
	values.Values[value.SettingKey] = value
	if err := s.writeRecordToFile(values, filePath); err != nil {
		return nil, err
	}
	return value, nil
}

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
