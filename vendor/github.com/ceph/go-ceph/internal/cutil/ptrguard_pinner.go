//go:build go1.21
// +build go1.21

package cutil

import (
	"runtime"
	"unsafe"
)

// PtrGuard respresents a guarded Go pointer (pointing to memory allocated by Go
// runtime) stored in C memory (allocated by C)
type PtrGuard struct {
	cPtr   CPtr
	pinner runtime.Pinner
}

// NewPtrGuard writes the goPtr (pointing to Go memory) into C memory at the
// position cPtr, and returns a PtrGuard object.
func NewPtrGuard(cPtr CPtr, goPtr unsafe.Pointer) *PtrGuard {
	var v PtrGuard
	v.pinner.Pin(goPtr)
	v.cPtr = cPtr
	p := (*unsafe.Pointer)(unsafe.Pointer(cPtr))
	*p = goPtr
	return &v
}

// Release removes the guarded Go pointer from the C memory by overwriting it
// with NULL.
func (v *PtrGuard) Release() {
	p := (*unsafe.Pointer)(unsafe.Pointer(v.cPtr))
	*p = nil
	v.pinner.Unpin()
}
