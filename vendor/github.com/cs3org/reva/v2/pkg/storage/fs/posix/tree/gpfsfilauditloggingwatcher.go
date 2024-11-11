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
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"time"
)

type GpfsFileAuditLoggingWatcher struct {
	tree *Tree
}

type lwe struct {
	Event        string
	Path         string
	BytesWritten string
}

func NewGpfsFileAuditLoggingWatcher(tree *Tree, auditLogFile string) (*GpfsFileAuditLoggingWatcher, error) {
	w := &GpfsFileAuditLoggingWatcher{
		tree: tree,
	}

	_, err := os.Stat(auditLogFile)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *GpfsFileAuditLoggingWatcher) Watch(path string) {
start:
	file, err := os.Open(path)
	if err != nil {
		// try again later
		time.Sleep(5 * time.Second)
		goto start
	}
	defer file.Close()

	// Seek to the end of the file
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		time.Sleep(5 * time.Second)
		goto start
	}

	reader := bufio.NewReader(file)
	ev := &lwe{}
	for {
		line, err := reader.ReadString('\n')
		switch err {
		case nil:
			err := json.Unmarshal([]byte(line), ev)
			if err != nil {
				continue
			}
			if isLockFile(ev.Path) || isTrash(ev.Path) || w.tree.isUpload(ev.Path) {
				continue
			}
			switch ev.Event {
			case "CREATE":
				go func() { _ = w.tree.Scan(ev.Path, ActionCreate, false) }()
			case "CLOSE":
				bytesWritten, err := strconv.Atoi(ev.BytesWritten)
				if err == nil && bytesWritten > 0 {
					go func() { _ = w.tree.Scan(ev.Path, ActionUpdate, false) }()
				}
			case "RENAME":
				go func() {
					_ = w.tree.Scan(ev.Path, ActionMove, false)
					_ = w.tree.WarmupIDCache(ev.Path, false, false)
				}()
			}
		case io.EOF:
			time.Sleep(1 * time.Second)
		default:
			time.Sleep(5 * time.Second)
			goto start
		}
	}
}
