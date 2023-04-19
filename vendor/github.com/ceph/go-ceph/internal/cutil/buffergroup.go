package cutil

// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

// BufferGroup is a helper structure that holds Go-allocated slices of
// C-allocated strings and their respective lengths. Useful for C functions
// that consume byte buffers with explicit length instead of null-terminated
// strings. When used as input arguments in C functions, caller must make sure
// the C code will not hold any pointers to either of the struct's attributes
// after that C function returns.
type BufferGroup struct {
	// C-allocated buffers.
	Buffers []CharPtr
	// Lengths of C buffers, where Lengths[i] = length(Buffers[i]).
	Lengths []SizeT
}

// TODO: should BufferGroup implementation change and the slices would contain
//       nested Go pointers, they must be pinned with PtrGuard.

// NewBufferGroupStrings returns new BufferGroup constructed from strings.
func NewBufferGroupStrings(strs []string) *BufferGroup {
	s := &BufferGroup{
		Buffers: make([]CharPtr, len(strs)),
		Lengths: make([]SizeT, len(strs)),
	}

	for i, str := range strs {
		bs := []byte(str)
		s.Buffers[i] = CharPtr(C.CBytes(bs))
		s.Lengths[i] = SizeT(len(bs))
	}

	return s
}

// NewBufferGroupBytes returns new BufferGroup constructed
// from slice of byte slices.
func NewBufferGroupBytes(bss [][]byte) *BufferGroup {
	s := &BufferGroup{
		Buffers: make([]CharPtr, len(bss)),
		Lengths: make([]SizeT, len(bss)),
	}

	for i, bs := range bss {
		s.Buffers[i] = CharPtr(C.CBytes(bs))
		s.Lengths[i] = SizeT(len(bs))
	}

	return s
}

// Free free()s the C-allocated memory.
func (s *BufferGroup) Free() {
	for _, ptr := range s.Buffers {
		C.free(unsafe.Pointer(ptr))
	}

	s.Buffers = nil
	s.Lengths = nil
}

// BuffersPtr returns a pointer to the beginning of the Buffers slice.
func (s *BufferGroup) BuffersPtr() CharPtrPtr {
	if len(s.Buffers) == 0 {
		return nil
	}

	return CharPtrPtr(&s.Buffers[0])
}

// LengthsPtr returns a pointer to the beginning of the Lengths slice.
func (s *BufferGroup) LengthsPtr() SizeTPtr {
	if len(s.Lengths) == 0 {
		return nil
	}

	return SizeTPtr(&s.Lengths[0])
}

func testBufferGroupGet(s *BufferGroup, index int) (str string, length int) {
	bs := C.GoBytes(unsafe.Pointer(s.Buffers[index]), C.int(s.Lengths[index]))
	return string(bs), int(s.Lengths[index])
}
