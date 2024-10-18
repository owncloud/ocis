package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

import (
	"unsafe"
)

// SetLocator sets the key for mapping objects to pgs within an io context.
// Until a different locator key is set, all objects in this io context will be placed in the same pg.
// To reset the locator, an empty string must be set.
//
// Implements:
//
//	void rados_ioctx_locator_set_key(rados_ioctx_t io, const char *key);
func (ioctx *IOContext) SetLocator(locator string) {
	if locator == "" {
		C.rados_ioctx_locator_set_key(ioctx.ioctx, nil)
	} else {
		var cLoc *C.char = C.CString(locator)
		defer C.free(unsafe.Pointer(cLoc))
		C.rados_ioctx_locator_set_key(ioctx.ioctx, cLoc)
	}
}
