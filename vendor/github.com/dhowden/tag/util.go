// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"bytes"
	"encoding/binary"
	"io"
)

func getBit(b byte, n uint) bool {
	x := byte(1 << n)
	return (b & x) == x
}

func get7BitChunkedInt(b []byte) int {
	var n int
	for _, x := range b {
		n = n << 7
		n |= int(x)
	}
	return n
}

func getInt(b []byte) int {
	var n int
	for _, x := range b {
		n = n << 8
		n |= int(x)
	}
	return n
}

func readUint64LittleEndian(r io.Reader) (uint64, error) {
	b, err := readBytes(r, 8)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b), nil
}

// readBytesMaxUpfront is the max up-front allocation allowed
const readBytesMaxUpfront = 10 << 20 // 10MB

func readBytes(r io.Reader, n uint) ([]byte, error) {
	if n > readBytesMaxUpfront {
		b := &bytes.Buffer{}
		if _, err := io.CopyN(b, r, int64(n)); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}

	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func readString(r io.Reader, n uint) (string, error) {
	b, err := readBytes(r, n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func readUint(r io.Reader, n uint) (uint, error) {
	x, err := readInt(r, n)
	if err != nil {
		return 0, err
	}
	return uint(x), nil
}

func readInt(r io.Reader, n uint) (int, error) {
	b, err := readBytes(r, n)
	if err != nil {
		return 0, err
	}
	return getInt(b), nil
}

func read7BitChunkedUint(r io.Reader, n uint) (uint, error) {
	b, err := readBytes(r, n)
	if err != nil {
		return 0, err
	}
	return uint(get7BitChunkedInt(b)), nil
}

func readUint32LittleEndian(r io.Reader) (uint32, error) {
	b, err := readBytes(r, 4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}
