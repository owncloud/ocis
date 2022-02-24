// Package store implements the go-micro store interface
package store

import (
	"errors"
	"sync"

	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
)

var m = &sync.RWMutex{}

// ListBundles returns all bundles in the dataPath folder that match the given type.
func (s Store) ListBundles(bundleType settingsmsg.Bundle_Type, bundleIDs []string) ([]*settingsmsg.Bundle, error) {
	return nil, errors.New("not implemented")
}

// ReadBundle tries to find a bundle by the given id within the dataPath.
func (s Store) ReadBundle(bundleID string) (*settingsmsg.Bundle, error) {
	return nil, errors.New("not implemented")
}

// ReadSetting tries to find a setting by the given id within the dataPath.
func (s Store) ReadSetting(settingID string) (*settingsmsg.Setting, error) {
	return nil, errors.New("not implemented")
}

// WriteBundle writes the given record into a file within the dataPath.
func (s Store) WriteBundle(record *settingsmsg.Bundle) (*settingsmsg.Bundle, error) {
	return nil, errors.New("not implemented")
}

// AddSettingToBundle adds the given setting to the bundle with the given bundleID.
func (s Store) AddSettingToBundle(bundleID string, setting *settingsmsg.Setting) (*settingsmsg.Setting, error) {
	return nil, errors.New("not implemented")
}

// RemoveSettingFromBundle removes the setting from the bundle with the given ids.
func (s Store) RemoveSettingFromBundle(bundleID string, settingID string) error {
	return errors.New("not implemented")
}
