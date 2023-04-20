package timespec

/*
#include <time.h>
*/
import "C"

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// Timespec behaves similarly to C's struct timespec.
// Timespec is used to retain fidelity to the C based file systems
// apis that could be lossy with the use of Go time types.
type Timespec unix.Timespec

// CTimespecPtr is an unsafe pointer wrapping C's `struct timespec`.
type CTimespecPtr unsafe.Pointer

// CStructToTimespec creates a new Timespec for the C 'struct timespec'.
func CStructToTimespec(cts CTimespecPtr) Timespec {
	t := (*C.struct_timespec)(cts)

	return Timespec{
		Sec:  int64(t.tv_sec),
		Nsec: int64(t.tv_nsec),
	}
}

// CopyToCStruct copies the time values from a Timespec to a previously
// allocated C `struct timespec`. Due to restrictions on Cgo the C pointer
// must be passed via the CTimespecPtr wrapper.
func CopyToCStruct(ts Timespec, cts CTimespecPtr) {
	t := (*C.struct_timespec)(cts)
	t.tv_sec = C.time_t(ts.Sec)
	t.tv_nsec = C.long(ts.Nsec)
}
