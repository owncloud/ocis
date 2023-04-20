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

package mime

import (
	"path"
	"sync"

	gomime "github.com/cubewise-code/go-mime"
)

const defaultMimeDir = "httpd/unix-directory"

var mimes sync.Map

func init() {
	mimes = sync.Map{}
}

// RegisterMime is a package level function that registers
// a mime type with the given extension.
// TODO(labkode): check that we do not override mime type mappings?
func RegisterMime(ext, mime string) {
	mimes.Store(ext, mime)
}

// Detect returns the mimetype associated with the given filename.
func Detect(isDir bool, fn string) string {
	if isDir {
		return defaultMimeDir
	}

	ext := path.Ext(fn)

	mimeType := getCustomMime(ext)

	if mimeType == "" {
		mimeType = gomime.TypeByExtension(ext)
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return mimeType
}

func getCustomMime(ext string) string {
	if m, ok := mimes.Load(ext); ok {
		return m.(string)
	}
	return ""
}
