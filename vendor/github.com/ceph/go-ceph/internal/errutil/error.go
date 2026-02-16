package errutil

type cephErrno int

// Error returns the error string for the errno.
func (e cephErrno) Error() string {
	_, strerror := FormatErrno(int(e))
	return strerror
}

// cephError combines the source/component that generated the error and its
// related errno.
type cephError struct {
	source string
	errno  cephErrno
}

// Error returns the error string with the source and errno.
func (e cephError) Error() string {
	return FormatErrorCode(e.source, int(e.errno))
}

// Unwrap returns an error without the source.
func (e cephError) Unwrap() error {
	if e.errno == 0 {
		return nil
	}

	return e.errno
}

// Is checks if both errors have the same errno.
func (e cephError) Is(err error) bool {
	ce, ok := err.(cephError)
	if !ok {
		return false
	}

	return e.errno == ce.errno
}

// ErrorCode returns the errno of the error.
func (e cephError) ErrorCode() int {
	return int(e.errno)
}

// GetError returns a new error that can be compared with errors.Is(),
// independently of the source/component of the error.
func GetError(source string, e int) error {
	if e == 0 {
		return nil
	}
	return cephError{
		source: source,
		errno:  cephErrno(int(e)),
	}
}
