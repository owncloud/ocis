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
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
)

func init() {
	registerMigration("0001", Migration0001{})
}

type Migration0001 struct{}

// Migrate creates the spaces directory structure
func (m Migration0001) Migrate(migrator *Migrator) (Result, error) {
	migrator.log.Info().Msg("Migrating spaces directory structure...")

	// create spaces folder and iterate over existing nodes to populate it
	nodesPath := filepath.Join(migrator.lu.InternalRoot(), "nodes")
	fi, err := os.Stat(nodesPath)
	if err == nil && fi.IsDir() {
		f, err := os.Open(nodesPath)
		if err != nil {
			return stateFailed, err
		}
		nodes, err := f.Readdir(0)
		if err != nil {
			return stateFailed, err
		}

		for _, n := range nodes {
			nodePath := filepath.Join(nodesPath, n.Name())

			attr, err := migrator.lu.MetadataBackend().Get(context.Background(), nodePath, prefixes.ParentidAttr)
			if err == nil && string(attr) == node.RootID {
				if err := m.moveNode(migrator, n.Name(), n.Name()); err != nil {
					migrator.log.Error().Err(err).
						Str("space", n.Name()).
						Msg("could not move space")
					continue
				}
				m.linkSpaceNode(migrator, "personal", n.Name())
			}
		}
		// TODO delete nodesPath if empty
	}
	return stateSucceeded, nil
}

// Rollback is not implemented
func (Migration0001) Rollback(_ *Migrator) (Result, error) {
	return stateFailed, errors.New("rollback not implemented")
}

func (m Migration0001) moveNode(migrator *Migrator, spaceID, nodeID string) error {
	dirPath := filepath.Join(migrator.lu.InternalRoot(), "nodes", nodeID)
	f, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	children, err := f.Readdir(0)
	if err != nil {
		return err
	}
	for _, child := range children {
		old := filepath.Join(migrator.lu.InternalRoot(), "nodes", child.Name())
		new := filepath.Join(migrator.lu.InternalRoot(), "spaces", lookup.Pathify(spaceID, 1, 2), "nodes", lookup.Pathify(child.Name(), 4, 2))
		if err := os.Rename(old, new); err != nil {
			migrator.log.Error().Err(err).
				Str("space", spaceID).
				Str("nodes", child.Name()).
				Str("oldpath", old).
				Str("newpath", new).
				Msg("could not rename node")
		}
		if child.IsDir() {
			if err := m.moveNode(migrator, spaceID, child.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

// linkSpace creates a new symbolic link for a space with the given type st, and node id
func (m Migration0001) linkSpaceNode(migrator *Migrator, spaceType, spaceID string) {
	spaceTypesPath := filepath.Join(migrator.lu.InternalRoot(), "spacetypes", spaceType, spaceID)
	expectedTarget := "../../spaces/" + lookup.Pathify(spaceID, 1, 2) + "/nodes/" + lookup.Pathify(spaceID, 4, 2)
	linkTarget, err := os.Readlink(spaceTypesPath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Symlink(expectedTarget, spaceTypesPath)
		if err != nil {
			migrator.log.Error().Err(err).
				Str("space_type", spaceType).
				Str("space", spaceID).
				Msg("could not create symlink")
		}
	} else {
		if err != nil {
			migrator.log.Error().Err(err).
				Str("space_type", spaceType).
				Str("space", spaceID).
				Msg("could not read symlink")
		}
		if linkTarget != expectedTarget {
			migrator.log.Warn().
				Str("space_type", spaceType).
				Str("space", spaceID).
				Str("expected", expectedTarget).
				Str("actual", linkTarget).
				Msg("expected a different link target")
		}
	}
}
