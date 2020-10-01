package errors

import (
	"fmt"
)

type AlreadyExistsErr struct {
	TypeName, Key, Value string
}

func (e *AlreadyExistsErr) Error() string {
	return fmt.Sprintf("%s with %s=%s does already exist", e.TypeName, e.Key, e.Value)
}

func IsAlreadyExistsErr(e error) bool {
	_, ok := e.(*AlreadyExistsErr)
	return ok
}

type NotFoundErr struct {
	TypeName, Key, Value string
}

func (e *NotFoundErr) Error() string {
	return fmt.Sprintf("%s with %s=%s not found", e.TypeName, e.Key, e.Value)
}

func IsNotFoundErr(e error) bool {
	_, ok := e.(*NotFoundErr)
	return ok
}
