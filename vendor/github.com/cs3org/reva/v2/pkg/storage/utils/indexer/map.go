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

package indexer

import (
	"github.com/cs3org/reva/v2/pkg/storage/utils/indexer/index"
	"github.com/iancoleman/strcase"
)

// typeMap stores the indexer layout at runtime.

type typeMap map[tName]typeMapping
type tName = string
type fieldName = string

type typeMapping struct {
	PKFieldName    string
	IndicesByField map[fieldName][]index.Index
}

func (m typeMap) addIndex(typeName string, pkName string, idx index.Index) {
	if val, ok := m[typeName]; ok {
		val.IndicesByField[strcase.ToCamel(idx.IndexBy().String())] = append(val.IndicesByField[strcase.ToCamel(idx.IndexBy().String())], idx)
		return
	}
	m[typeName] = typeMapping{
		PKFieldName: pkName,
		IndicesByField: map[string][]index.Index{
			strcase.ToCamel(idx.IndexBy().String()): {idx},
		},
	}
}
