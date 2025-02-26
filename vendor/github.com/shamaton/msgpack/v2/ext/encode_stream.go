package ext

import (
	"io"
	"reflect"

	"github.com/shamaton/msgpack/v2/internal/common"
)

// StreamEncoder is interface that extended encoder should implement
type StreamEncoder interface {
	Code() int8
	Type() reflect.Type
	Write(w StreamWriter, value reflect.Value) error
}

// StreamWriter is provided some writing functions for extended format by user
type StreamWriter struct {
	w   io.Writer
	buf *common.Buffer
}

func CreateStreamWriter(w io.Writer, buf *common.Buffer) StreamWriter {
	return StreamWriter{w, buf}
}

func (w *StreamWriter) WriteByte1Int64(value int64) error {
	return w.buf.Write(w.w,
		byte(value),
	)
}

func (w *StreamWriter) WriteByte2Int64(value int64) error {
	return w.buf.Write(w.w,
		byte(value>>8),
		byte(value),
	)
}

func (w *StreamWriter) WriteByte4Int64(value int64) error {
	return w.buf.Write(w.w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
}

func (w *StreamWriter) WriteByte8Int64(value int64) error {
	return w.buf.Write(w.w,
		byte(value>>56),
		byte(value>>48),
		byte(value>>40),
		byte(value>>32),
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
}

func (w *StreamWriter) WriteByte1Uint64(value uint64) error {
	return w.buf.Write(w.w,
		byte(value),
	)
}

func (w *StreamWriter) WriteByte2Uint64(value uint64) error {
	return w.buf.Write(w.w,
		byte(value>>8),
		byte(value),
	)
}

func (w *StreamWriter) WriteByte4Uint64(value uint64) error {
	return w.buf.Write(w.w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
}

func (w *StreamWriter) WriteByte8Uint64(value uint64) error {
	return w.buf.Write(w.w,
		byte(value>>56),
		byte(value>>48),
		byte(value>>40),
		byte(value>>32),
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
}

func (w *StreamWriter) WriteByte1Int(value int) error {
	return w.buf.Write(w.w,
		byte(value),
	)
}

func (w *StreamWriter) WriteByte2Int(value int) error {
	return w.buf.Write(w.w,
		byte(value>>8),
		byte(value),
	)
}

func (w *StreamWriter) WriteByte4Int(value int) error {
	return w.buf.Write(w.w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
}

func (w *StreamWriter) WriteByte4Uint32(value uint32) error {
	return w.buf.Write(w.w,
		byte(value>>24),
		byte(value>>16),
		byte(value>>8),
		byte(value),
	)
}

func (w *StreamWriter) WriteBytes(bs []byte) error {
	return w.buf.Write(w.w, bs...)
}
