//go:build !mimic
// +build !mimic

package rados

// #include <rados/librados.h>
import "C"

const (
	// OpFlagFAdviseFUA optionally support FUA (force unit access) on write
	// requests.
	OpFlagFAdviseFUA = OpFlags(C.LIBRADOS_OP_FLAG_FADVISE_FUA)
)
