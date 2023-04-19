package rados

// #cgo LDFLAGS: -lrados
// #include <errno.h>
// #include <stdlib.h>
// #include <rados/librados.h>
//
import "C"

// OperationFlags control the behavior of read and write operations.
type OperationFlags int

const (
	// OperationNoFlag indicates no special behavior is requested.
	OperationNoFlag = OperationFlags(C.LIBRADOS_OPERATION_NOFLAG)
	// OperationBalanceReads TODO
	OperationBalanceReads = OperationFlags(C.LIBRADOS_OPERATION_BALANCE_READS)
	// OperationLocalizeReads TODO
	OperationLocalizeReads = OperationFlags(C.LIBRADOS_OPERATION_LOCALIZE_READS)
	// OperationOrderReadsWrites TODO
	OperationOrderReadsWrites = OperationFlags(C.LIBRADOS_OPERATION_ORDER_READS_WRITES)
	// OperationIgnoreCache TODO
	OperationIgnoreCache = OperationFlags(C.LIBRADOS_OPERATION_IGNORE_CACHE)
	// OperationSkipRWLocks TODO
	OperationSkipRWLocks = OperationFlags(C.LIBRADOS_OPERATION_SKIPRWLOCKS)
	// OperationIgnoreOverlay TODO
	OperationIgnoreOverlay = OperationFlags(C.LIBRADOS_OPERATION_IGNORE_OVERLAY)
	// OperationFullTry send request to a full cluster or pool, ops such as delete
	// can succeed while other ops will return out-of-space errors.
	OperationFullTry = OperationFlags(C.LIBRADOS_OPERATION_FULL_TRY)
	// OperationFullForce TODO
	OperationFullForce = OperationFlags(C.LIBRADOS_OPERATION_FULL_FORCE)
	// OperationIgnoreRedirect TODO
	OperationIgnoreRedirect = OperationFlags(C.LIBRADOS_OPERATION_IGNORE_REDIRECT)
	// OperationOrderSnap TODO
	OperationOrderSnap = OperationFlags(C.LIBRADOS_OPERATION_ORDERSNAP)
)
