// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/rogpeppe/go-internal/lockedfile"
	"github.com/rs/zerolog"
)

var allMigrations = []string{"0001", "0002", "0003"}

const (
	resultFailed            = "failed"
	resultSucceeded         = "succeeded"
	resultSucceededRunAgain = "runagain"
)

type migrationState struct {
	State   string
	Message string
}
type migrationStates map[string]migrationState

// Result represents the result of a migration run
type Result string

// Migrator runs migrations on an existing decomposedfs
type Migrator struct {
	lu     *lookup.Lookup
	states migrationStates
	log    *zerolog.Logger
}

// New returns a new Migrator instance
func New(lu *lookup.Lookup, log *zerolog.Logger) Migrator {
	return Migrator{
		lu:  lu,
		log: log,
	}
}

// RunMigrations runs all migrations in sequence. Note this sequence must not be changed or it might
// damage existing decomposed fs.
func (m *Migrator) RunMigrations() error {
	lock, err := lockedfile.OpenFile(filepath.Join(m.lu.InternalRoot(), ".migrations.lock"), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer lock.Close()

	err = m.readStates()
	if err != nil {
		return err
	}

	for _, migration := range allMigrations {
		s := m.states[migration]
		if s.State == "succeeded" {
			continue
		}

		migrateMethod := reflect.ValueOf(m).MethodByName("Migration" + migration)
		v := migrateMethod.Call(nil)
		s.State = string(v[0].Interface().(Result))
		if v[1].Interface() != nil {
			err := v[1].Interface().(error)
			m.log.Error().Err(err).Msg("migration " + migration + " failed")
			s.Message = err.Error()
		}

		m.states[migration] = s
		err = m.writeStates()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrator) readStates() error {
	m.states = migrationStates{}

	d, err := os.ReadFile(filepath.Join(m.lu.InternalRoot(), ".migrations"))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if len(d) > 0 {
		err = json.Unmarshal(d, &m.states)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) writeStates() error {
	d, err := json.Marshal(m.states)
	if err != nil {
		m.log.Error().Err(err).Msg("could not marshal migration states")
		return nil
	}
	return os.WriteFile(filepath.Join(m.lu.InternalRoot(), ".migrations"), d, 0600)
}
