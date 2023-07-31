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

	"github.com/shamaton/msgpack/v2"
)

func init() {
	registerMigration("0005", Migration0005{})
}

type Migration0005 struct{}

// Migrate fixes the messagepack space index data structure
func (Migration0005) Migrate(migrator *Migrator) (Result, error) {
	root := migrator.lu.InternalRoot()

	indexes, err := filepath.Glob(filepath.Join(root, "indexes", "**", "*.mpk"))
	if err != nil {
		return stateFailed, err
	}
	for _, i := range indexes {
		migrator.log.Info().Str("root", migrator.lu.InternalRoot()).Msg("Fixing index format of " + i)

		// Read old-format index
		oldData, err := os.ReadFile(i)
		if err != nil {
			return stateFailed, err
		}
		oldIndex := map[string][]byte{}
		err = msgpack.Unmarshal(oldData, &oldIndex)
		if err != nil {
			// likely already migrated -> skip
			migrator.log.Warn().Str("root", migrator.lu.InternalRoot()).Msg("Invalid index format found in " + i)
			continue
		}

		// Write new-format index
		newIndex := map[string]string{}
		for k, v := range oldIndex {
			newIndex[k] = string(v)
		}
		newData, err := msgpack.Marshal(newIndex)
		if err != nil {
			return stateFailed, err
		}
		err = os.WriteFile(i, newData, 0600)
		if err != nil {
			return stateFailed, err
		}
	}
	migrator.log.Info().Msg("done.")
	return stateSucceeded, nil
}

// Rollback rolls back the migration
func (Migration0005) Rollback(migrator *Migrator) (Result, error) {
	root := migrator.lu.InternalRoot()

	indexes, err := filepath.Glob(filepath.Join(root, "indexes", "**", "*.mpk"))
	if err != nil {
		return stateFailed, err
	}
	for _, i := range indexes {
		migrator.log.Info().Str("root", migrator.lu.InternalRoot()).Msg("Fixing index format of " + i)

		oldData, err := os.ReadFile(i)
		if err != nil {
			return stateFailed, err
		}
		oldIndex := map[string]string{}
		err = msgpack.Unmarshal(oldData, &oldIndex)
		if err != nil {
			migrator.log.Warn().Str("root", migrator.lu.InternalRoot()).Msg("Invalid index format found in " + i)
			continue
		}

		// Write new-format index
		newIndex := map[string][]byte{}
		for k, v := range oldIndex {
			newIndex[k] = []byte(v)
		}
		newData, err := msgpack.Marshal(newIndex)
		if err != nil {
			return stateFailed, err
		}
		err = os.WriteFile(i, newData, 0600)
		if err != nil {
			return stateFailed, err
		}
	}
	migrator.log.Info().Msg("done.")
	return stateDown, nil
}
