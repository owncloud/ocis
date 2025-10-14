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

package exporters

import (
	"github.com/owncloud/reva/v2/pkg/mentix/exchangers"
	"github.com/owncloud/reva/v2/pkg/mentix/meshdata"
)

// Exporter is the interface that all exporters must implement.
type Exporter interface {
	exchangers.Exchanger
}

// BaseExporter implements basic exporter functionality common to all exporters.
type BaseExporter struct {
	exchangers.BaseExchanger
}

// Start starts the exporter.
func (exporter *BaseExporter) Start() error {
	// Initialize the exporter with empty data
	_ = exporter.Update(meshdata.Map{})
	return nil
}
