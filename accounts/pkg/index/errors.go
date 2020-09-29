package index

import (
	"fmt"
)

type alreadyExistsErr struct {
	typeName, key, val string
}

func (e *alreadyExistsErr) Error() string {
	return fmt.Sprintf("%s with %s=%s does already exist", e.typeName, e.key, e.val)
}

func IsAlreadyExistsErr(e error) bool {
	_, ok := e.(*alreadyExistsErr)
	return ok
}

type notFoundErr struct {
	typeName, key, val string
}

func (e *notFoundErr) Error() string {
	return fmt.Sprintf("%s with %s=%s not found", e.typeName, e.key, e.val)
}

func IsNotFoundErr(e error) bool {
	_, ok := e.(*notFoundErr)
	return ok
}
