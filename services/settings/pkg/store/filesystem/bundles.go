// Package store implements the go-micro store interface
package store

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/gofrs/uuid"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/errortypes"
)

var m = &sync.RWMutex{}

// ListBundles returns all bundles in the dataPath folder that match the given type.
func (s Store) ListBundles(bundleType settingsmsg.Bundle_Type, bundleIDs []string) ([]*settingsmsg.Bundle, error) {
	// FIXME: list requests should be ran against a cache, not FS
	m.RLock()
	defer m.RUnlock()

	bundlesFolder := s.buildFolderPathForBundles(false)
	bundleFiles, err := ioutil.ReadDir(bundlesFolder)
	if err != nil {
		return []*settingsmsg.Bundle{}, nil
	}

	records := make([]*settingsmsg.Bundle, 0, len(bundleFiles))
	for _, bundleFile := range bundleFiles {
		record := settingsmsg.Bundle{}
		err = s.parseRecordFromFile(&record, filepath.Join(bundlesFolder, bundleFile.Name()))
		if err != nil {
			s.Logger.Warn().Msgf("error reading %v", bundleFile)
			continue
		}
		if record.Type != bundleType {
			continue
		}
		if len(bundleIDs) > 0 && !containsStr(record.Id, bundleIDs) {
			continue
		}
		records = append(records, &record)
	}

	return records, nil
}

// containsStr checks if the strs slice contains str
func containsStr(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

// ReadBundle tries to find a bundle by the given id within the dataPath.
func (s Store) ReadBundle(bundleID string) (*settingsmsg.Bundle, error) {
	m.RLock()
	defer m.RUnlock()

	filePath := s.buildFilePathForBundle(bundleID, false)
	record := settingsmsg.Bundle{}
	if err := s.parseRecordFromFile(&record, filePath); err != nil {
		return nil, err
	}

	s.Logger.Debug().Msgf("read contents from file: %v", filePath)
	return &record, nil
}

// ReadSetting tries to find a setting by the given id within the dataPath.
func (s Store) ReadSetting(settingID string) (*settingsmsg.Setting, error) {
	m.RLock()
	defer m.RUnlock()

	bundles, err := s.ListBundles(settingsmsg.Bundle_TYPE_DEFAULT, []string{})
	if err != nil {
		return nil, err
	}
	for _, bundle := range bundles {
		for _, setting := range bundle.Settings {
			if setting.Id == settingID {
				return setting, nil
			}
		}
	}
	return nil, fmt.Errorf("could not read setting: %v", settingID)
}

// WriteBundle writes the given record into a file within the dataPath.
func (s Store) WriteBundle(record *settingsmsg.Bundle) (*settingsmsg.Bundle, error) {
	// FIXME: locking should happen on the file here, not globally.
	m.Lock()
	defer m.Unlock()
	if record.Id == "" {
		record.Id = uuid.Must(uuid.NewV4()).String()
	}
	filePath := s.buildFilePathForBundle(record.Id, true)
	if err := s.writeRecordToFile(record, filePath); err != nil {
		return nil, err
	}

	s.Logger.Debug().Msgf("request contents written to file: %v", filePath)
	return record, nil
}

// AddSettingToBundle adds the given setting to the bundle with the given bundleID.
func (s Store) AddSettingToBundle(bundleID string, setting *settingsmsg.Setting) (*settingsmsg.Setting, error) {
	bundle, err := s.ReadBundle(bundleID)
	if err != nil {
		if _, notFound := err.(errortypes.BundleNotFound); !notFound {
			return nil, err
		}
		bundle = new(settingsmsg.Bundle)
		bundle.Id = bundleID
		bundle.Type = settingsmsg.Bundle_TYPE_DEFAULT
	}
	if setting.Id == "" {
		setting.Id = uuid.Must(uuid.NewV4()).String()
	}
	setSetting(bundle, setting)
	_, err = s.WriteBundle(bundle)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

// RemoveSettingFromBundle removes the setting from the bundle with the given ids.
func (s Store) RemoveSettingFromBundle(bundleID string, settingID string) error {
	bundle, err := s.ReadBundle(bundleID)
	if err != nil {
		return nil
	}
	if ok := removeSetting(bundle, settingID); ok {
		if _, err := s.WriteBundle(bundle); err != nil {
			return err
		}
	}
	return nil
}

// indexOfSetting finds the index of the given setting within the given bundle.
// returns -1 if the setting was not found.
func indexOfSetting(bundle *settingsmsg.Bundle, settingID string) int {
	for index := range bundle.Settings {
		s := bundle.Settings[index]
		if s.Id == settingID {
			return index
		}
	}
	return -1
}

// setSetting will append or overwrite the given setting within the given bundle
func setSetting(bundle *settingsmsg.Bundle, setting *settingsmsg.Setting) {
	m.Lock()
	defer m.Unlock()
	index := indexOfSetting(bundle, setting.Id)
	if index == -1 {
		bundle.Settings = append(bundle.Settings, setting)
	} else {
		bundle.Settings[index] = setting
	}
}

// removeSetting will remove the given setting from the given bundle
func removeSetting(bundle *settingsmsg.Bundle, settingID string) bool {
	m.Lock()
	defer m.Unlock()
	index := indexOfSetting(bundle, settingID)
	if index == -1 {
		return false
	}
	bundle.Settings = append(bundle.Settings[:index], bundle.Settings[index+1:]...)
	return true
}
