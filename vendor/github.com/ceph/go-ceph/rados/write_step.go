package rados

// #include <stdint.h>
import "C"

import (
	"unsafe"
)

type writeStep struct {
	withoutUpdate
	withoutFree
	// the c pointer utilizes the Go byteslice data and no free is needed

	// inputs:
	b []byte

	// arguments:
	cBuffer   *C.char
	cDataLen  C.size_t
	cWriteLen C.size_t
	cOffset   C.uint64_t
}

func newWriteStep(b []byte, writeLen, offset uint64) *writeStep {
	return &writeStep{
		b:         b,
		cBuffer:   (*C.char)(unsafe.Pointer(&b[0])), // TODO: must be pinned
		cDataLen:  C.size_t(len(b)),
		cWriteLen: C.size_t(writeLen),
		cOffset:   C.uint64_t(offset),
	}
}
