//go:build !no_ptrguard
// +build !no_ptrguard

package cutil

import (
	"unsafe"
)

// SyncBuffer is a C buffer connected to a data slice
type SyncBuffer struct {
	pg *PtrGuard
}

// NewSyncBuffer creates a C buffer from a data slice and stores it at CPtr
func NewSyncBuffer(cPtr CPtr, data []byte) *SyncBuffer {
	var v SyncBuffer
	v.pg = NewPtrGuard(cPtr, unsafe.Pointer(&data[0]))
	return &v
}

// Release releases the C buffer and nulls its stored pointer
func (v *SyncBuffer) Release() {
	v.pg.Release()
}

// Sync asserts that changes in the C buffer are available in the data
// slice
func (*SyncBuffer) Sync() {}
