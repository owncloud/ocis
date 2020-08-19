// Package store implements the go-micro store interface
package store

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/gofrs/uuid"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
)

var m = &sync.RWMutex{}

// ListBundles returns all bundles in the dataPath folder that match the given type.
func (s Store) ListBundles(bundleType proto.Bundle_Type) ([]*proto.Bundle, error) {
	// FIXME: list requests should be ran against a cache, not FS
	m.RLock()
	defer m.RUnlock()

	bundlesFolder := s.buildFolderPathForBundles(false)
	bundleFiles, err := ioutil.ReadDir(bundlesFolder)
	if err != nil {
		return []*proto.Bundle{}, nil
	}

	records := make([]*proto.Bundle, 0, len(bundleFiles))
	for _, bundleFile := range bundleFiles {
		record := proto.Bundle{}
		err = s.parseRecordFromFile(&record, filepath.Join(bundlesFolder, bundleFile.Name()))
		if err != nil {
			s.Logger.Warn().Msgf("error reading %v", bundleFile)
			continue
		}
		if record.Type != bundleType {
			continue
		}
		records = append(records, &record)
	}

	return records, nil
}

// ReadBundle tries to find a bundle by the given id within the dataPath.
func (s Store) ReadBundle(bundleID string) (*proto.Bundle, error) {
	// FIXME: locking should happen on the file here, not globally.
	m.RLock()
	defer m.RUnlock()

	filePath := s.buildFilePathForBundle(bundleID, false)
	record := proto.Bundle{}
	if err := s.parseRecordFromFile(&record, filePath); err != nil {
		return nil, err
	}

	s.Logger.Debug().Msgf("read contents from file: %v", filePath)
	return &record, nil
}

// ReadSetting tries to find a setting by the given id within the dataPath.
func (s Store) ReadSetting(settingID string) (*proto.Setting, error) {
	// FIXME: locking should happen on the file here, not globally.
	m.RLock()
	defer m.RUnlock()

	bundles, err := s.ListBundles(proto.Bundle_TYPE_DEFAULT)
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
	return nil, merrors.NotFound(settingID, fmt.Sprintf("could not read setting: %v", settingID))
}

// WriteBundle writes the given record into a file within the dataPath.
func (s Store) WriteBundle(record *proto.Bundle) (*proto.Bundle, error) {
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
