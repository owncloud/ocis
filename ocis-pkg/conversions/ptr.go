package conversions

// ToPointer converts a value to a pointer
func ToPointer[T any](val T) *T {
	return &val
}

// ToValue converts a pointer to a value
func ToValue[T any](ptr *T) T {
	if ptr == nil {
		var t T
		return t
	}

	return *ptr
}

// ToPointerSlice converts a slice of values to a slice of pointers
func ToPointerSlice[E any](s []E) []*E {
	rs := make([]*E, len(s))

	for i, v := range s {
		rs[i] = ToPointer(v)
	}

	return rs
}

// ToValueSlice converts a slice of pointers to a slice of values
func ToValueSlice[E any](s []*E) []E {
	rs := make([]E, len(s))

	for i, v := range s {
		rs[i] = ToValue(v)
	}

	return rs
}
