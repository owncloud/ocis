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
	"fmt"

	"github.com/pablodz/inotifywaitgo/inotifywaitgo"
)

type InotifyWatcher struct {
	tree *Tree
}

func NewInotifyWatcher(tree *Tree) *InotifyWatcher {
	return &InotifyWatcher{
		tree: tree,
	}
}

func (iw *InotifyWatcher) Watch(path string) {
	events := make(chan inotifywaitgo.FileEvent)
	errors := make(chan error)

	go inotifywaitgo.WatchPath(&inotifywaitgo.Settings{
		Dir:        path,
		FileEvents: events,
		ErrorChan:  errors,
		KillOthers: true,
		Options: &inotifywaitgo.Options{
			Recursive: true,
			Events: []inotifywaitgo.EVENT{
				inotifywaitgo.CREATE,
				inotifywaitgo.MOVED_TO,
				inotifywaitgo.MOVED_FROM,
				inotifywaitgo.CLOSE_WRITE,
				inotifywaitgo.DELETE,
			},
			Monitor: true,
		},
		Verbose: false,
	})

	for {
		select {
		case event := <-events:
			for _, e := range event.Events {
				if isLockFile(event.Filename) || isTrash(event.Filename) || iw.tree.isUpload(event.Filename) {
					continue
				}
				go func() {
					switch e {
					case inotifywaitgo.DELETE:
						_ = iw.tree.Scan(event.Filename, ActionDelete, event.IsDir)
					case inotifywaitgo.MOVED_FROM:
						_ = iw.tree.Scan(event.Filename, ActionMoveFrom, event.IsDir)
					case inotifywaitgo.CREATE, inotifywaitgo.MOVED_TO:
						_ = iw.tree.Scan(event.Filename, ActionCreate, event.IsDir)
					case inotifywaitgo.CLOSE_WRITE:
						_ = iw.tree.Scan(event.Filename, ActionUpdate, event.IsDir)
					}
				}()
			}

		case err := <-errors:
			switch err.Error() {
			case inotifywaitgo.NOT_INSTALLED:
				panic("Error: inotifywait is not installed")
			case inotifywaitgo.INVALID_EVENT:
				// ignore
			default:
				fmt.Printf("Error: %s\n", err)
			}
		}
	}
}
