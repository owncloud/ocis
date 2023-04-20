// Copyright 2018-2023 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package migrator

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cs3org/reva/v2/pkg/logger"
)

// Migration0002 migrates spacetypes to indexes
func (m *Migrator) Migration0002() (Result, error) {
	m.log.Info().Msg("Migrating space types indexes...")

	spaceTypesPath := filepath.Join(m.lu.InternalRoot(), "spacetypes")
	fi, err := os.Stat(spaceTypesPath)
	if err == nil && fi.IsDir() {

		f, err := os.Open(spaceTypesPath)
		if err != nil {
			return resultFailed, err
		}
		spaceTypes, err := f.Readdir(0)
		if err != nil {
			return resultFailed, err
		}

		for _, st := range spaceTypes {
			err := m.moveSpaceType(st.Name())
			if err != nil {
				logger.New().Error().Err(err).
					Str("space", st.Name()).
					Msg("could not move space")
				continue
			}
		}

		// delete spacetypespath
		d, err := os.Open(spaceTypesPath)
		if err != nil {
			logger.New().Error().Err(err).
				Str("spacetypesdir", spaceTypesPath).
				Msg("could not open spacetypesdir")
			return resultFailed, nil
		}
		defer d.Close()
		_, err = d.Readdirnames(1) // Or f.Readdir(1)
		if err == io.EOF {
			// directory is empty we can delete
			err := os.Remove(spaceTypesPath)
			if err != nil {
				logger.New().Error().Err(err).
					Str("spacetypesdir", d.Name()).
					Msg("could not delete")
			}
		} else {
			logger.New().Error().Err(err).
				Str("spacetypesdir", d.Name()).
				Msg("could not delete, not empty")
		}
	}
	return resultSucceeded, nil
}

func (m *Migrator) moveSpaceType(spaceType string) error {
	dirPath := filepath.Join(m.lu.InternalRoot(), "spacetypes", spaceType)
	f, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	children, err := f.Readdir(0)
	if err != nil {
		return err
	}
	for _, child := range children {
		old := filepath.Join(m.lu.InternalRoot(), "spacetypes", spaceType, child.Name())
		target, err := os.Readlink(old)
		if err != nil {
			logger.New().Error().Err(err).
				Str("space", spaceType).
				Str("nodes", child.Name()).
				Str("oldLink", old).
				Msg("could not read old symlink")
			continue
		}
		newDir := filepath.Join(m.lu.InternalRoot(), "indexes", "by-type", spaceType)
		if err := os.MkdirAll(newDir, 0700); err != nil {
			logger.New().Error().Err(err).
				Str("space", spaceType).
				Str("nodes", child.Name()).
				Str("targetDir", newDir).
				Msg("could not read old symlink")
		}
		newLink := filepath.Join(newDir, child.Name())
		if err := os.Symlink(filepath.Join("..", target), newLink); err != nil {
			logger.New().Error().Err(err).
				Str("space", spaceType).
				Str("nodes", child.Name()).
				Str("oldpath", old).
				Str("newpath", newLink).
				Msg("could not rename node")
			continue
		}
		if err := os.Remove(old); err != nil {
			logger.New().Error().Err(err).
				Str("space", spaceType).
				Str("nodes", child.Name()).
				Str("oldLink", old).
				Msg("could not remove old symlink")
			continue
		}
	}
	if err := os.Remove(dirPath); err != nil {
		logger.New().Error().Err(err).
			Str("space", spaceType).
			Str("dir", dirPath).
			Msg("could not remove spaces folder, folder probably not empty")
	}
	return nil
}
