package retry

// WithTries repeatedly calls a SizeFunc with increasing sizes until either it
// returns nil, or a maximum number of tries has been reached. If the returned hint is
// DoubleSize or indicating a size not greater than the current size, the size
// is doubled.
// WithTries should be used in conjunction with functions that provide a size
// hint other than DoubleSize. If the underlying function provides a size hint
// then the number of iterations can be kept low (ideally 2) and the system
// will make fewer memory allocations.
func WithTries(size int, maxTries int, f SizeFunc) {
	attempt := 0
	for {
		attempt++
		if attempt > maxTries {
			break
		}
		hint := f(size)
		if hint == nil {
			break
		}
		if hint.size() > size {
			size = hint.size()
		} else {
			size *= 2
		}
	}
}
