package cutil

/*
#include <stdlib.h>
#include <string.h>
typedef void* voidptr;
*/
import "C"

import (
	"math"
	"unsafe"
)

const (
	// MaxIdx is the maximum index on 32 bit systems
	MaxIdx = math.MaxInt32 // 2GB, max int32 value, should be safe

	// PtrSize is the size of a pointer
	PtrSize = C.sizeof_voidptr

	// SizeTSize is the size of C.size_t
	SizeTSize = C.sizeof_size_t
)

// Compile-time assertion ensuring that Go's `int` is at least as large as C's.
const _ = unsafe.Sizeof(int(0)) - C.sizeof_int

// SizeT wraps size_t from C.
type SizeT C.size_t

// This section contains a bunch of types that are basically just
// unsafe.Pointer but have specific types to help "self document" what the
// underlying pointer is really meant to represent.

// CPtr is an unsafe.Pointer to C allocated memory
type CPtr unsafe.Pointer

// CharPtrPtr is an unsafe pointer wrapping C's `char**`.
type CharPtrPtr unsafe.Pointer

// CharPtr is an unsafe pointer wrapping C's `char*`.
type CharPtr unsafe.Pointer

// SizeTPtr is an unsafe pointer wrapping C's `size_t*`.
type SizeTPtr unsafe.Pointer

// FreeFunc is a wrapper around calls to, or act like, C's free function.
type FreeFunc func(unsafe.Pointer)

// Malloc is C.malloc
func Malloc(s SizeT) CPtr { return CPtr(C.malloc(C.size_t(s))) }

// Free is C.free
func Free(p CPtr) { C.free(unsafe.Pointer(p)) }

// CString is C.CString
func CString(s string) CharPtr { return CharPtr((C.CString(s))) }

// CBytes is C.CBytes
func CBytes(b []byte) CPtr { return CPtr(C.CBytes(b)) }

// Memcpy is C.memcpy
func Memcpy(dst, src CPtr, n SizeT) {
	C.memcpy(unsafe.Pointer(dst), unsafe.Pointer(src), C.size_t(n))
}
