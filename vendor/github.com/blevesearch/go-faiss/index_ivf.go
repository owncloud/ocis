package faiss

/*
#include <faiss/c_api/IndexIVFFlat_c.h>
#include <faiss/c_api/MetaIndexes_c.h>
#include <faiss/c_api/Index_c.h>
#include <faiss/c_api/IndexIVF_c.h>
#include <faiss/c_api/IndexIVF_c_ex.h>
#include <faiss/c_api/IndexScalarQuantizer_c.h>
*/
import "C"
import (
	"fmt"
)

func (idx *faissIndex) SetDirectMap(mapType int) (err error) {

	ivfPtr := C.faiss_IndexIVF_cast(idx.cPtr())
	if ivfPtr == nil {
		return errNotIVFIndex
	}
	if c := C.faiss_IndexIVF_set_direct_map(
		ivfPtr,
		C.int(mapType),
	); c != 0 {
		err = getLastError()
	}
	return err
}

func (idx *faissIndex) GetSubIndex() (Index, error) {

	ptr := C.faiss_IndexIDMap2_cast(idx.cPtr())
	if ptr == nil {
		return nil, fmt.Errorf("index is not a id map")
	}

	subIdx := C.faiss_IndexIDMap2_sub_index(ptr)
	if subIdx == nil {
		return nil, fmt.Errorf("couldn't retrieve the sub index")
	}

	return &IndexImpl{&faissIndex{subIdx}}, nil
}

// pass nprobe to be set as index time option for IVF indexes only.
// varying nprobe impacts recall but with an increase in latency.
func (idx *faissIndex) SetNProbe(nprobe int32) {
	ivfPtr := C.faiss_IndexIVF_cast(idx.cPtr())
	if ivfPtr == nil {
		return
	}
	C.faiss_IndexIVF_set_nprobe(ivfPtr, C.size_t(nprobe))
}

func (idx *faissIndex) IVFParams() (nprobe, nlist int) {
	ivfPtr := C.faiss_IndexIVF_cast(idx.cPtr())
	if ivfPtr == nil {
		return 0, 0
	}
	return int(C.faiss_IndexIVF_nprobe(ivfPtr)),
		int(C.faiss_IndexIVF_nlist(ivfPtr))
}

func (idx *faissIndex) IsSQIndex() bool {
	sqPtr := C.faiss_IndexScalarQuantizer_cast(idx.cPtr())
	return sqPtr != nil
}

func (idx *faissIndex) SetQuantizers(srcIndex Index) error {
	if !(idx.IsIVFIndex() && srcIndex.IsIVFIndex()) &&
		!(idx.IsSQIndex() && srcIndex.IsSQIndex()) {
		return fmt.Errorf("faissIndex SetQuantizers: %w, index type not supported", errFailedToSetQuantizers)
	}

	srcIndexPtr := srcIndex.cPtr()
	if srcIndexPtr == nil {
		return fmt.Errorf("coarse quantizer is not valid")
	}

	err := C.faiss_Set_quantizers(idx.idx, srcIndexPtr)
	if err != 0 {
		return fmt.Errorf("faissIndex SetQuantizers: %w", errFailedToSetQuantizers)
	}

	return nil
}
