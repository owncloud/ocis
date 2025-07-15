// Copyright 2018-2022 CERN
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

package tags

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

var (
	// character used to separate tags in lists
	_tagsep = ","
	// maximum number of tags
	_maxtags = 100
)

// Tags is a helper struct for merging, deleting and deduplicating the tags while preserving the order
type Tags struct {
	sep     string
	maxtags int

	t       []string
	exists  map[string]bool
	numtags int
}

// New creates a Tag struct from a slice of tags, e.g. ["tag1", "tag2"] or a list of tags, e.g. "tag1,tag2"
func New(ts ...string) *Tags {
	t := &Tags{sep: _tagsep, maxtags: _maxtags, exists: make(map[string]bool), t: make([]string, 0)}
	t.addTags(t.normalize(ts))
	return t
}

// Add appends a list of new tags and returns true if at least one was appended
func (t *Tags) Add(ts ...string) bool {
	return len(t.addTags(t.normalize(ts))) > 0
}

// AddValidated appends a list of new tags and validates them using the provided validator function
// It returns true if at least one was appended and an error if the validation failed
func (t *Tags) AddValidated(validator func([]string) error, ts ...string) (bool, error) {
	newTags := t.normalize(ts)
	err := validator(newTags)
	if err != nil {
		return false, err
	}
	return len(t.addTags(newTags)) > 0, nil
}

// Remove removes a list of tags and returns true if at least one was removed
func (t *Tags) Remove(s ...string) bool {
	var removed bool
	tags := t.normalize(s)

	for _, tag := range tags {
		if !t.exists[tag] {
			continue
		}

		for i, tt := range t.t {
			if tt == tag {
				t.t = append(t.t[:i], t.t[i+1:]...)
				break
			}
		}

		delete(t.exists, tag)
		removed = true
	}
	return removed
}

// AsList returns the tags converted to a list
func (t *Tags) AsList() string {
	return strings.Join(t.t, t.sep)
}

// AsSlice returns the tags as slice of strings
func (t *Tags) AsSlice() []string {
	return t.t
}

// adds the tags and returns a list of added tags
// the function receiving a normalized slice of tags, e.g.["tag1", "tag2"]
func (t *Tags) addTags(s []string) []string {
	added := make([]string, 0, len(t.t)+len(s))
	for _, tag := range s {
		if tag == "" {
			// ignore empty tags
			continue
		}

		if t.exists[tag] {
			// tag is already existing
			continue
		}

		if t.numtags >= t.maxtags {
			// max number of tags reached. We return silently without warning anyone
			break
		}

		added = append(added, tag)
		t.exists[tag] = true
		t.numtags++
	}

	t.t = append(added, t.t...)
	return added
}

// normalize splits the tags and removes empty tags
// the function receiving a slice of tags, e.g.["tag1", "tag2"] or a list of tags, e.g."tag1,tag2" or mixed.
func (t *Tags) normalize(s []string) []string {
	res := make([]string, 0, t.maxtags/2)
	for _, tt := range s {
		for _, tag := range strings.Split(tt, t.sep) {
			ttr := strings.TrimSpace(tag)
			if ttr == "" {
				// ignore empty tags
				continue
			}
			res = append(res, ttr)
		}
	}
	return res
}

// MaxLengthValidator returns a function that validates the length of each tag in a slice
func MaxLengthValidator(maxTagLength int) func([]string) error {
	if maxTagLength <= 0 {
		return func(tags []string) error { return nil }
	}
	return func(tags []string) error {
		t := make([]string, 0, maxTagLength)
		for _, tag := range tags {
			if !utf8.ValidString(tag) {
				return fmt.Errorf("tag [%s] contains invalid characters", tag)
			}
			if utf8.RuneCount([]byte(tag)) > maxTagLength {
				t = append(t, tag)
				continue
			}
		}
		if len(t) > 0 {
			return fmt.Errorf("tag [%s] too long, max length is %d", strings.Join(t, ", "), maxTagLength)
		}
		return nil
	}
}
