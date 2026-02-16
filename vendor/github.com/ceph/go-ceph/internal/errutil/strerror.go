/*
Package errutil provides common functions for dealing with error conditions for
all ceph api wrappers.
*/
package errutil

/* force XSI-complaint strerror_r() */

// #define _POSIX_C_SOURCE 200112L
// #undef _GNU_SOURCE
// #include <stdlib.h>
// #include <errno.h>
// #include <string.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// FormatErrno returns the absolute value of the errno as well as a string
// describing the errno. The string will be empty is the errno is not known.
func FormatErrno(errno int) (int, string) {
	buf := make([]byte, 1024)
	// strerror expects errno >= 0
	if errno < 0 {
		errno = -errno
	}

	ret := C.strerror_r(
		C.int(errno),
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.size_t(len(buf)))
	if ret != 0 {
		return errno, ""
	}

	return errno, C.GoString((*C.char)(unsafe.Pointer(&buf[0])))
}

// FormatErrorCode returns a string that describes the supplied error source
// and error code as a string. Suitable to use in Error() methods.  If the
// error code maps to an errno the string will contain a description of the
// error. Otherwise the string will only indicate the source and value if the
// value does not map to a known errno.
func FormatErrorCode(source string, errValue int) string {
	_, s := FormatErrno(errValue)
	if s == "" {
		return fmt.Sprintf("%s: ret=%d", source, errValue)
	}
	return fmt.Sprintf("%s: ret=%d, %s", source, errValue, s)
}
