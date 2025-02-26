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
	"errors"
	"fmt"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/cs3org/reva/v2/pkg/storage/utils/indexer/option"
)

func getType(v interface{}) (reflect.Value, error) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return reflect.Value{}, errors.New("failed to read value via reflection")
	}

	return rv, nil
}

func getTypeFQN(t interface{}) string {
	typ, _ := getType(t)
	typeName := path.Join(typ.Type().PkgPath(), typ.Type().Name())
	typeName = strings.ReplaceAll(typeName, "/", ".")
	return typeName
}

func valueOf(v interface{}, indexBy option.IndexBy) (string, error) {
	switch idxBy := indexBy.(type) {
	case option.IndexByField:
		return valueOfField(v, string(idxBy))
	case option.IndexByFunc:
		return idxBy.Func(v)
	default:
		return "", fmt.Errorf("unknown indexBy type")
	}
}

func valueOfField(v interface{}, field string) (string, error) {
	parts := strings.Split(field, ".")
	for i, part := range parts {
		r := reflect.ValueOf(v)
		if r.Kind() == reflect.Ptr {
			r = r.Elem()
		}
		f := reflect.Indirect(r).FieldByName(part)
		if f.Kind() == reflect.Ptr {
			f = f.Elem()
		}

		switch {
		case f.Kind() == reflect.Struct && i != len(parts)-1:
			v = f.Interface()
		case f.Kind() == reflect.String:
			return f.String(), nil
		case f.IsZero():
			return "", nil
		default:
			return strconv.Itoa(int(f.Int())), nil
		}
	}
	return "", nil
}
