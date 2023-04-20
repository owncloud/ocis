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
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
)

// Migration0003 migrates the file metadata to the current backend.
// Only the xattrs -> messagepack path is supported.
func (m *Migrator) Migration0003() (Result, error) {
	bod := lookup.DetectBackendOnDisk(m.lu.InternalRoot())
	if bod == "" {
		return resultFailed, errors.New("could not detect metadata backend on disk")
	}

	if bod != "xattrs" || m.lu.MetadataBackend().Name() != "messagepack" {
		return resultSucceededRunAgain, nil
	}

	m.log.Info().Str("root", m.lu.InternalRoot()).Msg("Migrating to messagepack metadata backend...")
	xattrs := metadata.XattrsBackend{}
	mpk := metadata.NewMessagePackBackend(m.lu.InternalRoot(), options.CacheOptions{})

	spaces, _ := filepath.Glob(filepath.Join(m.lu.InternalRoot(), "spaces", "*", "*"))
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

			attribs, err := xattrs.All(path)
			if err != nil {
				m.log.Error().Err(err).Str("path", path).Msg("error converting file")
				return err
			}
			if len(attribs) == 0 {
				return nil
			}

			err = mpk.SetMultiple(path, attribs, false)
			if err != nil {
				m.log.Error().Err(err).Str("path", path).Msg("error setting attributes")
				return err
			}

			for k := range attribs {
				err = xattrs.Remove(path, k)
				if err != nil {
					m.log.Debug().Err(err).Str("path", path).Msg("error removing xattr")
				}
			}

			return nil
		})
		if err != nil {
			m.log.Error().Err(err).Msg("error migrating nodes to messagepack metadata backend")
		}
	}

	m.log.Info().Msg("done.")
	return resultSucceeded, nil
}
