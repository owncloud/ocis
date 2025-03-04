// Copyright 2018-2024 CERN
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

package tree

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	kafka "github.com/segmentio/kafka-go"
)

type GpfsWatchFolderWatcher struct {
	tree    *Tree
	brokers []string
	log     *zerolog.Logger
}

func NewGpfsWatchFolderWatcher(tree *Tree, kafkaBrokers []string, log *zerolog.Logger) (*GpfsWatchFolderWatcher, error) {
	return &GpfsWatchFolderWatcher{
		tree:    tree,
		brokers: kafkaBrokers,
		log:     log,
	}, nil
}

func (w *GpfsWatchFolderWatcher) Watch(topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: w.brokers,
		GroupID: "ocis-posixfs",
		Topic:   topic,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		lwev := &lwe{}
		err = json.Unmarshal(m.Value, lwev)
		if err != nil {
			continue
		}

		if isLockFile(lwev.Path) || isTrash(lwev.Path) || w.tree.isUpload(lwev.Path) {
			continue
		}

		go func() {
			isDir := strings.Contains(lwev.Event, "IN_ISDIR")

			var err error
			switch {
			case strings.Contains(lwev.Event, "IN_DELETE"):
				err = w.tree.Scan(lwev.Path, ActionDelete, isDir)

			case strings.Contains(lwev.Event, "IN_MOVE_FROM"):
				err = w.tree.Scan(lwev.Path, ActionMoveFrom, isDir)

			case strings.Contains(lwev.Event, "IN_CREATE"):
				err = w.tree.Scan(lwev.Path, ActionCreate, isDir)

			case strings.Contains(lwev.Event, "IN_CLOSE_WRITE"):
				bytesWritten, convErr := strconv.Atoi(lwev.BytesWritten)
				if convErr == nil && bytesWritten > 0 {
					err = w.tree.Scan(lwev.Path, ActionUpdate, isDir)
				}
			case strings.Contains(lwev.Event, "IN_MOVED_TO"):
				err = w.tree.Scan(lwev.Path, ActionMove, isDir)
			}
			if err != nil {
				w.log.Error().Err(err).Str("path", lwev.Path).Msg("error scanning path")
			}
		}()
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
