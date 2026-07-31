//go:build ceph_preview

package rados

/*
#cgo LDFLAGS: -lrados
#include <stdlib.h>
#include <rados/librados.h>
*/
import "C"

import "unsafe"

// Checksum calculates the checksum of the given object data, using one of the supported checksum algorithms.
//
// Implements:
//
//	int rados_checksum(rados_ioctx_t io,
//	                   const char *oid,
//	                   rados_checksum_type_t type,
//	                   const char *init_value,
//	                   size_t init_value_len,
//	                   size_t len,
//	                   uint64_t off,
//	                   size_t chunk_size,
//	                   char *pchecksum,
//	                   size_t checksum_len);
func (ioctx *IOContext) Checksum(oid string, checksumType ChecksumType, dst []byte, opts *ChecksumOptions) error {
	// apply defaults
	if opts == nil {
		opts = &ChecksumOptions{}
	}
	if opts.InitValue == nil {
		initLen := 4
		if checksumType == ChecksumTypeXXHash64 {
			initLen = 8
		}
		opts.InitValue = make([]byte, initLen)
	}

	// call library
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	return getError(C.rados_checksum(
		ioctx.ioctx,
		coid,
		C.rados_checksum_type_t(checksumType),
		(*C.char)(unsafe.Pointer(&opts.InitValue[0])),
		C.size_t(len(opts.InitValue)),
		C.size_t(opts.Len),
		C.uint64_t(opts.Off),
		C.size_t(opts.ChunkSize),
		(*C.char)(unsafe.Pointer(&dst[0])),
		C.size_t(len(dst)),
	))
}

// ChecksumType indicates checksum algorithm types supported by the IOContext.Checksum method.
// Equivalent to the rados_checksum_type_t enum.
type ChecksumType uint32

const (
	// ChecksumTypeXXHash32 produces an encoded le32 checksum of the given object.
	ChecksumTypeXXHash32 = ChecksumType(C.LIBRADOS_CHECKSUM_TYPE_XXHASH32)
	// ChecksumTypeXXHash64 produces an encoded le64 checksum of the given object.
	ChecksumTypeXXHash64 = ChecksumType(C.LIBRADOS_CHECKSUM_TYPE_XXHASH64)
	// ChecksumTypeCRC32C produces an encoded le32 checksum of the given object.
	ChecksumTypeCRC32C = ChecksumType(C.LIBRADOS_CHECKSUM_TYPE_CRC32C)
)

// ChecksumOptions exposes non-required parameters for the Checksum method.
type ChecksumOptions struct {
	// Off sets the object offset to start checksumming in the object.
	// By default, the entire object will be checksummed.
	Off uint64
	// Len sets the the number of bytes to checksum in the object.
	// By default, the entire object will be checksummed.
	Len uint64
	// ChunkSize sets the length-aligned chunk size for the checksum calculation.
	// By default, the entire object will be checksummed as a single chunk.
	ChunkSize uint64
	// InitValue sets the initial value for the checksum calculation.
	// By default, the initial value will be a zeroed-out byte slice.
	InitValue []byte
}
