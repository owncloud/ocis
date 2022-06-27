package cutil

// The following code needs some explanation:
// This creates slices on top of the C memory buffers allocated before in
// order to safely and comfortably use them as arrays. First the void pointer
// is cast to a pointer to an array of the type that will be stored in the
// array. Because the size of an array is a constant, but the real array size
// is dynamic, we just use the biggest possible index value MaxIdx, to make
// sure it's always big enough. (Nothing is allocated by casting, so the size
// can be arbitrarily big.) So, if the array should store items of myType, the
// cast would be (*[MaxIdx]myItem)(myCMemPtr).
// From that array pointer a slice is created with the [start:end:capacity]
// syntax. The capacity must be set explicitly here, because by default it
// would be set to the size of the original array, which is MaxIdx, which
// doesn't reflect reality in this case. This results in definitions like:
// cSlice := (*[MaxIdx]myItem)(myCMemPtr)[:numOfItems:numOfItems]

////////// CPtr //////////

// CPtrCSlice is a C allocated slice of C pointers.
type CPtrCSlice []CPtr

// NewCPtrCSlice returns a CPtrSlice.
// Similar to CString it must be freed with slice.Free()
func NewCPtrCSlice(size int) CPtrCSlice {
	if size == 0 {
		return nil
	}
	cMem := Malloc(SizeT(size) * PtrSize)
	cSlice := (*[MaxIdx]CPtr)(cMem)[:size:size]
	return cSlice
}

// Ptr returns a pointer to CPtrSlice
func (v *CPtrCSlice) Ptr() CPtr {
	if len(*v) == 0 {
		return nil
	}
	return CPtr(&(*v)[0])
}

// Free frees a CPtrSlice
func (v *CPtrCSlice) Free() {
	Free(v.Ptr())
	*v = nil
}

////////// SizeT //////////

// SizeTCSlice is a C allocated slice of C.size_t.
type SizeTCSlice []SizeT

// NewSizeTCSlice returns a SizeTCSlice.
// Similar to CString it must be freed with slice.Free()
func NewSizeTCSlice(size int) SizeTCSlice {
	if size == 0 {
		return nil
	}
	cMem := Malloc(SizeT(size) * SizeTSize)
	cSlice := (*[MaxIdx]SizeT)(cMem)[:size:size]
	return cSlice
}

// Ptr returns a pointer to SizeTCSlice
func (v *SizeTCSlice) Ptr() CPtr {
	if len(*v) == 0 {
		return nil
	}
	return CPtr(&(*v)[0])
}

// Free frees a SizeTCSlice
func (v *SizeTCSlice) Free() {
	Free(v.Ptr())
	*v = nil
}
