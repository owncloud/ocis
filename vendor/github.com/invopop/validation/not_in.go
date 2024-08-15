// Copyright 2018 Qiang Xue, Google LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import "reflect"

// ErrNotInInvalid is the error that returns when a value is in a list.
var ErrNotInInvalid = NewError("validation_not_in_invalid", "must not be in list")

// NotIn returns a validation rule that checks if a value is absent from the given list of values.
// Like with In(), reflect.DeepEqual() will be used to determine if two values are equal.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func NotIn[T any](values ...T) NotInRule[T] {
	return NotInRule[T]{
		elements: values,
		err:      ErrNotInInvalid,
	}
}

// NotInRule is a validation rule that checks if a value is absent from the given list of values.
type NotInRule[T any] struct {
	elements []T
	err      Error
}

// Validate checks if the given value is valid or not.
func (r NotInRule[T]) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return nil
	}

	for _, e := range r.elements {
		if reflect.DeepEqual(e, value) {
			return r.err
		}
	}
	return nil
}

// Error sets the error message for the rule.
func (r NotInRule[T]) Error(message string) NotInRule[T] {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r NotInRule[T]) ErrorObject(err Error) NotInRule[T] {
	r.err = err
	return r
}
