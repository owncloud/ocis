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
	"os"
	"path/filepath"
	"strings"

	"github.com/shamaton/msgpack/v2"
)

func init() {
	registerMigration("0004", Migration0004{})
}

type Migration0004 struct{}

// Migrate migrates the directory tree based space indexes to messagepack
func (Migration0004) Migrate(migrator *Migrator) (Result, error) {
	root := migrator.lu.InternalRoot()

	// migrate user indexes
	users, err := os.ReadDir(filepath.Join(root, "indexes", "by-user-id"))
	if err != nil {
		migrator.log.Warn().Err(err).Msg("error listing user indexes")
	}
	for _, user := range users {
		if !user.IsDir() {
			continue
		}
		id := user.Name()
		indexPath := filepath.Join(root, "indexes", "by-user-id", id+".mpk")
		dirIndexPath := filepath.Join(root, "indexes", "by-user-id", id)

		migrator.log.Info().Str("root", migrator.lu.InternalRoot()).Msg("Migrating " + indexPath + " to messagepack index format...")
		err := migrateSpaceIndex(indexPath, dirIndexPath)
		if err != nil {
			migrator.log.Error().Err(err).Str("path", dirIndexPath).Msg("error migrating index")
		}
	}

	// migrate group indexes
	groups, err := os.ReadDir(filepath.Join(root, "indexes", "by-group-id"))
	if err != nil {
		migrator.log.Warn().Err(err).Msg("error listing group indexes")
	}
	for _, group := range groups {
		if !group.IsDir() {
			continue
		}
		id := group.Name()
		indexPath := filepath.Join(root, "indexes", "by-group-id", id+".mpk")
		dirIndexPath := filepath.Join(root, "indexes", "by-group-id", id)

		migrator.log.Info().Str("root", migrator.lu.InternalRoot()).Msg("Migrating " + indexPath + " to messagepack index format...")
		err := migrateSpaceIndex(indexPath, dirIndexPath)
		if err != nil {
			migrator.log.Error().Err(err).Str("path", dirIndexPath).Msg("error migrating index")
		}
	}

	// migrate project indexes
	for _, spaceType := range []string{"personal", "project", "share"} {
		indexPath := filepath.Join(root, "indexes", "by-type", spaceType+".mpk")
		dirIndexPath := filepath.Join(root, "indexes", "by-type", spaceType)

		_, err := os.Stat(dirIndexPath)
		if err != nil {
			continue
		}

		migrator.log.Info().Str("root", migrator.lu.InternalRoot()).Msg("Migrating " + indexPath + " to messagepack index format...")
		err = migrateSpaceIndex(indexPath, dirIndexPath)
		if err != nil {
			migrator.log.Error().Err(err).Str("path", dirIndexPath).Msg("error migrating index")
		}
	}

	migrator.log.Info().Msg("done.")
	return stateSucceeded, nil
}

func migrateSpaceIndex(indexPath, dirIndexPath string) error {
	links := map[string][]byte{}
	m, err := filepath.Glob(dirIndexPath + "/*")
	if err != nil {
		return err
	}
	for _, match := range m {
		link, err := os.Readlink(match)
		if err != nil {
			continue
		}
		links[filepath.Base(match)] = []byte(link)
	}

	// rewrite index as file
	d, err := msgpack.Marshal(links)
	if err != nil {
		return err
	}
	err = os.WriteFile(indexPath, d, 0600)
	if err != nil {
		return err
	}
	return os.RemoveAll(dirIndexPath)
}

// Rollback migrates the directory messagepack indexes to symlinks
func (Migration0004) Rollback(m *Migrator) (Result, error) {
	root := m.lu.InternalRoot()

	// migrate user indexes
	users, err := filepath.Glob(filepath.Join(root, "indexes", "by-user-id", "*.mpk"))
	if err != nil {
		m.log.Warn().Err(err).Msg("error listing user indexes")
	}
	for _, indexPath := range users {
		dirIndexPath := strings.TrimSuffix(indexPath, ".mpk")

		m.log.Info().Str("root", m.lu.InternalRoot()).Msg("Migrating " + indexPath + " to symlinks index format...")
		err := downSpaceIndex(indexPath, dirIndexPath)
		if err != nil {
			m.log.Error().Err(err).Str("path", dirIndexPath).Msg("error migrating index")
		}
	}

	// migrate group indexes
	groups, err := filepath.Glob(filepath.Join(root, "indexes", "by-group-id", "*.mpk"))
	if err != nil {
		m.log.Warn().Err(err).Msg("error listing group indexes")
	}
	for _, indexPath := range groups {
		dirIndexPath := strings.TrimSuffix(indexPath, ".mpk")

		m.log.Info().Str("root", m.lu.InternalRoot()).Msg("Migrating " + indexPath + " to symlinks index format...")
		err := downSpaceIndex(indexPath, dirIndexPath)
		if err != nil {
			m.log.Error().Err(err).Str("path", dirIndexPath).Msg("error migrating index")
		}
	}

	// migrate project indexes
	for _, spaceType := range []string{"personal", "project", "share"} {
		indexPath := filepath.Join(root, "indexes", "by-type", spaceType+".mpk")
		dirIndexPath := filepath.Join(root, "indexes", "by-type", spaceType)

		_, err := os.Stat(indexPath)
		if err != nil || os.IsNotExist(err) {
			continue
		}

		m.log.Info().Str("root", m.lu.InternalRoot()).Msg("Migrating " + indexPath + " to symlinks index format...")
		err = downSpaceIndex(indexPath, dirIndexPath)
		if err != nil {
			m.log.Error().Err(err).Str("path", dirIndexPath).Msg("error migrating index")
		}
	}
	return stateDown, nil
}

func downSpaceIndex(indexPath, dirIndexPath string) error {
	d, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	links := map[string][]byte{}
	err = msgpack.Unmarshal(d, &links)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dirIndexPath, 0700)
	if err != nil {
		return err
	}
	for link, target := range links {
		err = os.Symlink(string(target), filepath.Join(dirIndexPath, link))
		if err != nil {
			return err
		}
	}

	return os.Remove(indexPath)
}
