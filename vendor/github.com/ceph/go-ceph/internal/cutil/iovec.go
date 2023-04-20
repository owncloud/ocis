package cutil

/*
#include <stdlib.h>
#include <sys/uio.h>
*/
import "C"
import (
	"unsafe"
)

// Iovec is a slice of iovec structs. Might have allocated C memory, so it must
// be freed with the Free() method.
type Iovec struct {
	iovec []C.struct_iovec
	sbs   []*SyncBuffer
}

const iovecSize = C.sizeof_struct_iovec

// ByteSlicesToIovec creates an Iovec and links it to Go buffers in data.
func ByteSlicesToIovec(data [][]byte) (v Iovec) {
	n := len(data)
	iovecMem := C.malloc(iovecSize * C.size_t(n))
	v.iovec = (*[MaxIdx]C.struct_iovec)(iovecMem)[:n:n]
	for i, b := range data {
		sb := NewSyncBuffer(CPtr(&v.iovec[i].iov_base), b)
		v.sbs = append(v.sbs, sb)
		v.iovec[i].iov_len = C.size_t(len(b))
	}
	return
}

// Sync makes sure the slices contain the same as the C buffers
func (v *Iovec) Sync() {
	for _, sb := range v.sbs {
		sb.Sync()
	}
}

// Pointer returns a pointer to the iovec
func (v *Iovec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&v.iovec[0])
}

// Len returns a pointer to the iovec
func (v *Iovec) Len() int {
	return len(v.iovec)
}

// Free the C memory in the Iovec.
func (v *Iovec) Free() {
	for _, sb := range v.sbs {
		sb.Release()
	}
	if len(v.iovec) != 0 {
		C.free(unsafe.Pointer(&v.iovec[0]))
	}
	v.iovec = nil
}
