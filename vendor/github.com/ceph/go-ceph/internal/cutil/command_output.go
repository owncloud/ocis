package cutil

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

// CommandOutput can be used to manage the outputs of ceph's *_command
// functions.
type CommandOutput struct {
	free      FreeFunc
	outBuf    *C.char
	outBufLen C.size_t
	outs      *C.char
	outsLen   C.size_t
}

// NewCommandOutput returns an empty CommandOutput. The pointers that
// a CommandOutput provides can be used to get the results of ceph's
// *_command functions.
func NewCommandOutput() *CommandOutput {
	return &CommandOutput{
		free: free,
	}
}

// SetFreeFunc sets the function used to free memory held by CommandOutput.
// Not all uses of CommandOutput expect to use the basic C.free function
// and either require or prefer the use of a custom deallocation function.
// Use SetFreeFunc to change the free function and return the modified
// CommandOutput object.
func (co *CommandOutput) SetFreeFunc(f FreeFunc) *CommandOutput {
	co.free = f
	return co
}

// Free any C memory tracked by this object.
func (co *CommandOutput) Free() {
	if co.outBuf != nil {
		co.free(unsafe.Pointer(co.outBuf))
	}
	if co.outs != nil {
		co.free(unsafe.Pointer(co.outs))
	}
}

// OutBuf returns an unsafe wrapper around a pointer to a `char*`.
func (co *CommandOutput) OutBuf() CharPtrPtr {
	return CharPtrPtr(&co.outBuf)
}

// OutBufLen returns an unsafe wrapper around a pointer to a size_t.
func (co *CommandOutput) OutBufLen() SizeTPtr {
	return SizeTPtr(&co.outBufLen)
}

// Outs returns an unsafe wrapper around a pointer to a `char*`.
func (co *CommandOutput) Outs() CharPtrPtr {
	return CharPtrPtr(&co.outs)
}

// OutsLen returns an unsafe wrapper around a pointer to a size_t.
func (co *CommandOutput) OutsLen() SizeTPtr {
	return SizeTPtr(&co.outsLen)
}

// GoValues returns native go values converted from the internal C-language
// values tracked by this object.
func (co *CommandOutput) GoValues() (buf []byte, status string) {
	if co.outBufLen > 0 {
		buf = C.GoBytes(unsafe.Pointer(co.outBuf), C.int(co.outBufLen))
	}
	if co.outsLen > 0 {
		status = C.GoStringN(co.outs, C.int(co.outsLen))
	}
	return buf, status
}

// testSetString is only used by the unit tests for this file.
// It is located here due to the restriction on having import "C" in
// go test files. :-(
// It mimics a C function that takes a pointer to a
// string and length and allocates memory and sets the pointers
// to the new string and its length.
func testSetString(strp CharPtrPtr, lenp SizeTPtr, s string) {
	sp := (**C.char)(strp)
	lp := (*C.size_t)(lenp)
	*sp = C.CString(s)
	*lp = C.size_t(len(s))
}

// free wraps C.free.
// Required for unit tests that may not use cgo directly.
func free(p unsafe.Pointer) {
	C.free(p)
}
