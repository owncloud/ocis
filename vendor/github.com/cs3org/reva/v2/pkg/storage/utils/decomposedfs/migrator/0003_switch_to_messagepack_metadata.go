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
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cs3org/reva/v2/pkg/storage/cache"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
)

func init() {
	registerMigration("0003", Migration0003{})
}

type Migration0003 struct{}

// Migrate migrates the file metadata to the current backend.
// Only the xattrs -> messagepack path is supported.
func (m Migration0003) Migrate(migrator *Migrator) (Result, error) {
	bod := lookup.DetectBackendOnDisk(migrator.lu.InternalRoot())
	if bod == "" {
		return stateFailed, errors.New("could not detect metadata backend on disk")
	}

	if bod != "xattrs" || migrator.lu.MetadataBackend().Name() != "messagepack" {
		return stateSucceededRunAgain, nil
	}

	migrator.log.Info().Str("root", migrator.lu.InternalRoot()).Msg("Migrating to messagepack metadata backend...")
	xattrs := metadata.XattrsBackend{}
	mpk := metadata.NewMessagePackBackend(migrator.lu.InternalRoot(), cache.Config{})

	spaces, _ := filepath.Glob(filepath.Join(migrator.lu.InternalRoot(), "spaces", "*", "*"))
	for _, space := range spaces {
		err := filepath.WalkDir(filepath.Join(space, "nodes"), func(path string, _ fs.DirEntry, err error) error {
			// Do not continue on error
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, ".mpk") || strings.HasSuffix(path, ".flock") {
				// None of our business
				return nil
			}

			fi, err := os.Lstat(path)
			if err != nil {
				return err
			}

			if !fi.IsDir() && !fi.Mode().IsRegular() {
				return nil
			}

			mpkPath := mpk.MetadataPath(path)
			_, err = os.Stat(mpkPath)
			if err == nil {
				return nil
			}

			attribs, err := xattrs.All(context.Background(), path)
			if err != nil {
				migrator.log.Error().Err(err).Str("path", path).Msg("error converting file")
				return err
			}
			if len(attribs) == 0 {
				return nil
			}

			err = mpk.SetMultiple(context.Background(), path, attribs, false)
			if err != nil {
				migrator.log.Error().Err(err).Str("path", path).Msg("error setting attributes")
				return err
			}

			for k := range attribs {
				err = xattrs.Remove(context.Background(), path, k)
				if err != nil {
					migrator.log.Debug().Err(err).Str("path", path).Msg("error removing xattr")
				}
			}

			return nil
		})
		if err != nil {
			migrator.log.Error().Err(err).Msg("error migrating nodes to messagepack metadata backend")
		}
	}

	migrator.log.Info().Msg("done.")
	return stateSucceeded, nil
}

// Rollback is not implemented
func (Migration0003) Rollback(_ *Migrator) (Result, error) {
	return stateFailed, errors.New("rollback not implemented")
}
