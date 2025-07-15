package rados

// #include <stdint.h>
import "C"

import (
	"unsafe"
)

type readStep struct {
	withoutUpdate
	withoutFree
	// the c pointer utilizes the Go byteslice data and no free is needed

	// inputs:
	b []byte

	// arguments:
	cBuffer  *C.char
	cReadLen C.size_t
	cOffset  C.uint64_t
}

func newReadStep(b []byte, offset uint64) *readStep {
	return &readStep{
		b:        b,
		cBuffer:  (*C.char)(unsafe.Pointer(&b[0])), // TODO: must be pinned
		cReadLen: C.size_t(len(b)),
		cOffset:  C.uint64_t(offset),
	}
}
