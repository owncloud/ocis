package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
//
import "C"

// AllocHintFlags control the behavior of read and write operations.
type AllocHintFlags uint32

const (
	// AllocHintNoHint indicates no predefined behavior
	AllocHintNoHint = AllocHintFlags(0)
	// AllocHintSequentialWrite TODO
	AllocHintSequentialWrite = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_SEQUENTIAL_WRITE)
	// AllocHintRandomWrite TODO
	AllocHintRandomWrite = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_RANDOM_WRITE)
	// AllocHintSequentialRead TODO
	AllocHintSequentialRead = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_SEQUENTIAL_READ)
	// AllocHintRandomRead TODO
	AllocHintRandomRead = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_RANDOM_READ)
	// AllocHintAppendOnly TODO
	AllocHintAppendOnly = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_APPEND_ONLY)
	// AllocHintImmutable TODO
	AllocHintImmutable = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_IMMUTABLE)
	// AllocHintShortlived TODO
	AllocHintShortlived = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_SHORTLIVED)
	// AllocHintLonglived TODO
	AllocHintLonglived = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_LONGLIVED)
	// AllocHintCompressible TODO
	AllocHintCompressible = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_COMPRESSIBLE)
	// AllocHintIncompressible TODO
	AllocHintIncompressible = AllocHintFlags(C.LIBRADOS_ALLOC_HINT_FLAG_INCOMPRESSIBLE)
)
