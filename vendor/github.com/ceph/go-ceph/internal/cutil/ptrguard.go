package cutil

import (
	"sync"
	"unsafe"
)

// PtrGuard respresents a guarded Go pointer (pointing to memory allocated by Go
// runtime) stored in C memory (allocated by C)
type PtrGuard struct {
	// These mutexes will be used as binary semaphores for signalling events from
	// one thread to another, which - in contrast to other languages like C++ - is
	// possible in Go, that is a Mutex can be locked in one thread and unlocked in
	// another.
	stored, release sync.Mutex
	released        bool
}

// WARNING: using binary semaphores (mutexes) for signalling like this is quite
// a delicate task in order to avoid deadlocks or panics. Whenever changing the
// code logic, please review at least three times that there is no unexpected
// state possible. Usually the natural choice would be to use channels instead,
// but these can not easily passed to C code because of the pointer-to-pointer
// cgo rule, and would require the use of a Go object registry.

// NewPtrGuard writes the goPtr (pointing to Go memory) into C memory at the
// position cPtr, and returns a PtrGuard object.
func NewPtrGuard(cPtr CPtr, goPtr unsafe.Pointer) *PtrGuard {
	var v PtrGuard
	// Since the mutexes are used for signalling, they have to be initialized to
	// locked state, so that following lock attempts will block.
	v.release.Lock()
	v.stored.Lock()
	// Start a background go routine that lives until Release is called. This
	// calls a special function that makes sure the garbage collector doesn't touch
	// goPtr, stores it into C memory at position cPtr and then waits until it
	// reveices the "release" signal, after which it nulls out the C memory at
	// cPtr and then exits.
	go func() {
		storeUntilRelease(&v, (*CPtr)(cPtr), uintptr(goPtr))
	}()
	// Wait for the "stored" signal from the go routine when the Go pointer has
	// been stored to the C memory. <--(1)
	v.stored.Lock()
	return &v
}

// Release removes the guarded Go pointer from the C memory by overwriting it
// with NULL.
func (v *PtrGuard) Release() {
	if !v.released {
		v.released = true
		v.release.Unlock() // Send the "release" signal to the go routine. -->(2)
		v.stored.Lock()    // Wait for the second "stored" signal when the C memory
		//                    has been nulled out. <--(3)

	}
}

// The uintptrPtr() helper function below assumes that uintptr has the same size
// as a pointer, although in theory it could be larger.  Therefore we use this
// constant expression to assert size equality as a safeguard at compile time.
// How it works: if sizes are different, either the inner or outer expression is
// negative, which always fails with "constant ... overflows uintptr", because
// unsafe.Sizeof() is a uintptr typed constant.
const _ = -(unsafe.Sizeof(uintptr(0)) - PtrSize) // size assert
func uintptrPtr(p *CPtr) *uintptr {
	return (*uintptr)(unsafe.Pointer(p))
}

//go:uintptrescapes

// From https://golang.org/src/cmd/compile/internal/gc/lex.go:
// For the next function declared in the file any uintptr arguments may be
// pointer values converted to uintptr. This directive ensures that the
// referenced allocated object, if any, is retained and not moved until the call
// completes, even though from the types alone it would appear that the object
// is no longer needed during the call. The conversion to uintptr must appear in
// the argument list.
// Also see https://golang.org/cmd/compile/#hdr-Compiler_Directives

func storeUntilRelease(v *PtrGuard, cPtr *CPtr, goPtr uintptr) {
	uip := uintptrPtr(cPtr)
	*uip = goPtr      // store Go pointer in C memory at c_ptr
	v.stored.Unlock() // send "stored" signal to main thread -->(1)
	v.release.Lock()  // wait for "release" signal from main thread when
	//                   Release() has been called. <--(2)
	*uip = 0          // reset C memory to NULL
	v.stored.Unlock() // send second "stored" signal to main thread -->(3)
}
