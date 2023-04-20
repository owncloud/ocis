/*
 * Copyright 2017-2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package utils

import (
	"fmt"
)

// ErrorWithDescription is an interface binding the standard error
// inteface with a description.
type ErrorWithDescription interface {
	error
	Description() string
}

// DescribeError returns a wrapped version for errors which contain additional
// fields. The wrapped version contains all fields as a string value. Use this
// for general purpose logging of rich errors.
func DescribeError(err error) error {
	switch err.(type) {
	case ErrorWithDescription:
		err = fmt.Errorf("%s - %s", err, err.(ErrorWithDescription).Description())
	}

	return err
}

// ErrorAsFields returns a mapping of all fields of the provided error.
func ErrorAsFields(err error) map[string]interface{} {
	if err == nil {
		return nil
	}

	fields := make(map[string]interface{})
	fields["error"] = err.Error()
	switch err.(type) {
	case ErrorWithDescription:
		fields["desc"] = err.(ErrorWithDescription).Description()
	}

	return fields
}
