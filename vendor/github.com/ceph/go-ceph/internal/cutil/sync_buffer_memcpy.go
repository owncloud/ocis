//go:build no_ptrguard
// +build no_ptrguard

package cutil

// SyncBuffer is a C buffer connected to a data slice
type SyncBuffer struct {
	data []byte
	cPtr *CPtr
}

// NewSyncBuffer creates a C buffer from a data slice and stores it at CPtr
func NewSyncBuffer(cPtr CPtr, data []byte) *SyncBuffer {
	var v SyncBuffer
	v.data = data
	v.cPtr = (*CPtr)(cPtr)
	*v.cPtr = CBytes(data)
	return &v
}

// Release releases the C buffer and nulls its stored pointer
func (v *SyncBuffer) Release() {
	if v.cPtr != nil {
		Free(*v.cPtr)
		*v.cPtr = nil
		v.cPtr = nil
	}
	v.data = nil
}

// Sync asserts that changes in the C buffer are available in the data
// slice
func (v *SyncBuffer) Sync() {
	if v.cPtr == nil {
		return
	}
	Memcpy(CPtr(&v.data[0]), CPtr(*v.cPtr), SizeT(len(v.data)))
}
