package preprocessor

import (
	"io"

	"github.com/pkg/errors"
)

type LazyReadSeeker struct {
	reader   io.Reader
	buffer   []byte
	position int64
}

func NewLazyReadSeeker(r io.Reader) *LazyReadSeeker {
	return &LazyReadSeeker{
		reader: r,
		buffer: make([]byte, 0),
	}
}

func (l *LazyReadSeeker) Read(p []byte) (int, error) {
	// Fill buffer if necessary
	if l.position >= int64(len(l.buffer)) {
		temp := make([]byte, len(p))
		n, err := l.reader.Read(temp)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			return 0, io.EOF
		}
		l.buffer = append(l.buffer, temp[:n]...)
	}

	// Read from buffer
	n := copy(p, l.buffer[l.position:])
	l.position += int64(n)

	return n, nil
}

func (l *LazyReadSeeker) Seek(offset int64, whence int) (int64, error) {
	var newPos int64
	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = l.position + offset
	default:
		return 0, errors.New("seekEnd is not supported")
	}

	if newPos < 0 {
		return 0, errors.New("negative position is invalid")
	}

	l.position = newPos
	return l.position, nil
}
