package faiss

/*
#include <faiss/c_api/impl/AuxIndexStructures_c.h>
*/
import "C"

// IDSelector represents a set of IDs to remove.
type IDSelector struct {
	sel *C.FaissIDSelector
}

// Delete frees the memory associated with s.
func (s *IDSelector) Delete() {
	if s == nil || s.sel == nil {
		return
	}

	C.faiss_IDSelector_free(s.sel)
}

type IDSelectorBatch struct {
	sel      *C.FaissIDSelector
	batchSel *C.FaissIDSelector
}

// Delete frees the memory associated with s.
func (s *IDSelectorBatch) Delete() {
	if s == nil {
		return
	}

	if s.sel != nil {
		C.faiss_IDSelector_free(s.sel)
	}
	if s.batchSel != nil {
		C.faiss_IDSelector_free(s.batchSel)
	}
}

// NewIDSelectorRange creates a selector that removes IDs on [imin, imax).
func NewIDSelectorRange(imin, imax int64) (*IDSelector, error) {
	var sel *C.FaissIDSelectorRange
	c := C.faiss_IDSelectorRange_new(&sel, C.idx_t(imin), C.idx_t(imax))
	if c != 0 {
		return nil, getLastError()
	}
	return &IDSelector{(*C.FaissIDSelector)(sel)}, nil
}

// NewIDSelectorBatch creates a new batch selector.
func NewIDSelectorBatch(indices []int64) (*IDSelector, error) {
	var sel *C.FaissIDSelectorBatch
	if c := C.faiss_IDSelectorBatch_new(
		&sel,
		C.size_t(len(indices)),
		(*C.idx_t)(&indices[0]),
	); c != 0 {
		return nil, getLastError()
	}
	return &IDSelector{(*C.FaissIDSelector)(sel)}, nil
}

// NewIDSelectorNot creates a new Not selector, wrapped arround a
// batch selector, with the IDs in 'exclude'.
func NewIDSelectorNot(exclude []int64) (*IDSelectorBatch, error) {
	batchSelector, err := NewIDSelectorBatch(exclude)
	if err != nil {
		return nil, err
	}

	var sel *C.FaissIDSelectorNot
	if c := C.faiss_IDSelectorNot_new(
		&sel,
		batchSelector.sel,
	); c != 0 {
		batchSelector.Delete()
		return nil, getLastError()
	}
	return &IDSelectorBatch{sel: (*C.FaissIDSelector)(sel), batchSel: batchSelector.sel}, nil
}
