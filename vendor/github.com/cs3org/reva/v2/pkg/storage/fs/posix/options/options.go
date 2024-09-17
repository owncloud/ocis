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

package options

import (
	"time"

	decomposedoptions "github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Options struct {
	decomposedoptions.Options

	UseSpaceGroups bool `mapstructure:"use_space_groups"`

	ScanDebounceDelay time.Duration `mapstructure:"scan_debounce_delay"`

	WatchFS                 bool   `mapstructure:"watch_fs"`
	WatchType               string `mapstructure:"watch_type"`
	WatchPath               string `mapstructure:"watch_path"`
	WatchFolderKafkaBrokers string `mapstructure:"watch_folder_kafka_brokers"`
}

// New returns a new Options instance for the given configuration
func New(m map[string]interface{}) (*Options, error) {
	o := &Options{}
	if err := mapstructure.Decode(m, o); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}

	do, err := decomposedoptions.New(m)
	if err != nil {
		return nil, err
	}
	o.Options = *do

	return o, nil
}
