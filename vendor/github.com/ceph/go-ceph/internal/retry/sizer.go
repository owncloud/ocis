package retry

// Hint is a type for retry hints
type Hint interface {
	If(bool) Hint
	size() int
}

type hintInt int

func (hint hintInt) size() int {
	return int(hint)
}

// If is a convenience function, that returns a given hint only if a certain
// condition is met (for example a test for a "buffer too small" error).
// Otherwise it returns a nil which stops the retries.
func (hint hintInt) If(cond bool) Hint {
	if cond {
		return hint
	}
	return nil
}

// DoubleSize is a hint to retry with double the size
const DoubleSize = hintInt(0)

// Size returns a hint for a specific size
func Size(s int) Hint {
	return hintInt(s)
}

// SizeFunc is used to implement 'resize loops' that hides the complexity of the
// sizing away from most of the application. It's a function that takes a size
// argument and returns nil, if no retry is necessary, or a hint indicating the
// size for the next retry. If errors or other results are required from the
// function, the function can write them to function closures of the surrounding
// scope. See tests for examples.
type SizeFunc func(size int) (hint Hint)

// WithSizes repeatingly calls a SizeFunc with increasing sizes until either it
// returns nil, or the max size has been reached. If the returned hint is
// DoubleSize or indicating a size not greater than the current size, the size
// is doubled. If the hint or next size is greater than the max size, the max
// size is used for a last retry.
func WithSizes(size int, max int, f SizeFunc) {
	if size > max {
		return
	}
	for {
		hint := f(size)
		if hint == nil || size == max {
			break
		}
		if hint.size() > size {
			size = hint.size()
		} else {
			size *= 2
		}
		if size > max {
			size = max
		}
	}
}
