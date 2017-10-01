// Package binary evtends standard encoding/binary package,
// with additional functionality for other types
// used in raknet/minecraft protocols. (boolean, triad, etc.)
package binary

import (
	"encoding/binary"
	"unsafe"
)

const maxuint = ^uint(0)
const maxint = int(maxuint >> 1)

// A ByteOrder evtends encoding/binary.ByteOrder's functionality.
type ByteOrder interface {
	binary.ByteOrder

	Int16([]byte) int16
	Int32([]byte) int32
	Int64([]byte) int64

	PutInt16([]byte, int16)
	PutInt32([]byte, int32)
	PutInt64([]byte, int64)

	Triad([]byte) uint32
	PutTriad([]byte, uint32)

	Float32([]byte) float32
	Float64([]byte) float64
	PutFloat32([]byte, float32)
	PutFloat64([]byte, float64)
}

// LittleEndian wraps encoding/binary.LittleEndian
var LittleEndian ByteOrder = littleEndian{binary.LittleEndian}

// BigEndian wraps encoding/binary.BigEndian

type littleEndian struct{ binary.ByteOrder }

func (s littleEndian) Int16(b []byte) int16 {
	v := s.ByteOrder.Uint16(b)
	return *(*int16)(unsafe.Pointer(&v))
}

func (s littleEndian) PutInt16(b []byte, v int16) {
	s.ByteOrder.PutUint16(b, *(*uint16)(unsafe.Pointer(&v)))
}

func (s littleEndian) Int32(b []byte) int32 {
	v := s.ByteOrder.Uint32(b)
	return *(*int32)(unsafe.Pointer(&v))
}

func (s littleEndian) PutInt32(b []byte, v int32) {
	s.ByteOrder.PutUint32(b, *(*uint32)(unsafe.Pointer(&v)))
}

func (s littleEndian) Int64(b []byte) int64 {
	v := s.ByteOrder.Uint64(b)
	return *(*int64)(unsafe.Pointer(&v))
}

func (s littleEndian) PutInt64(b []byte, v int64) {
	s.ByteOrder.PutUint64(b, *(*uint64)(unsafe.Pointer(&v)))
}

func (littleEndian) Triad(b []byte) uint32 {
	_ = b[2] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

func (littleEndian) PutTriad(b []byte, v uint32) {
	_ = b[2] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}

func (s littleEndian) Float32(b []byte) float32 {
	n := s.ByteOrder.Uint32(b)
	return *(*float32)(unsafe.Pointer(&n))
}

func (s littleEndian) Float64(b []byte) float64 {
	n := s.ByteOrder.Uint64(b)
	return *(*float64)(unsafe.Pointer(&n))
}

func (s littleEndian) PutFloat32(b []byte, v float32) {
	s.PutUint32(b, *(*uint32)(unsafe.Pointer(&v)))
}

func (s littleEndian) PutFloat64(b []byte, v float64) {
	s.PutUint64(b, *(*uint64)(unsafe.Pointer(&v)))
}

type bigEndian struct{ binary.ByteOrder }

func (s bigEndian) Int16(b []byte) int16 {
	v := s.ByteOrder.Uint16(b)
	return *(*int16)(unsafe.Pointer(&v))
}

func (s bigEndian) PutInt16(b []byte, v int16) {
	s.ByteOrder.PutUint16(b, *(*uint16)(unsafe.Pointer(&v)))
}

func (s bigEndian) Int32(b []byte) int32 {
	v := s.ByteOrder.Uint32(b)
	return *(*int32)(unsafe.Pointer(&v))
}

func (s bigEndian) PutInt32(b []byte, v int32) {
	s.ByteOrder.PutUint32(b, *(*uint32)(unsafe.Pointer(&v)))
}

func (s bigEndian) Int64(b []byte) int64 {
	v := s.ByteOrder.Uint64(b)
	return *(*int64)(unsafe.Pointer(&v))
}

func (s bigEndian) PutInt64(b []byte, v int64) {
	s.ByteOrder.PutUint64(b, *(*uint64)(unsafe.Pointer(&v)))
}

func (bigEndian) Triad(b []byte) uint32 {
	_ = b[2] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

func (bigEndian) PutTriad(b []byte, v uint32) {
	_ = b[2] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

func (s bigEndian) Float32(b []byte) float32 {
	n := s.ByteOrder.Uint32(b)
	return *(*float32)(unsafe.Pointer(&n))
}

func (s bigEndian) Float64(b []byte) float64 {
	n := s.ByteOrder.Uint64(b)
	return *(*float64)(unsafe.Pointer(&n))
}

func (s bigEndian) PutFloat32(b []byte, v float32) {
	s.PutUint32(b, *(*uint32)(unsafe.Pointer(&v)))
}

func (s bigEndian) PutFloat64(b []byte, v float64) {
	s.PutUint64(b, *(*uint64)(unsafe.Pointer(&v)))
}

var BigEndian ByteOrder = bigEndian{binary.BigEndian}

// String decodes a string from buf and returns that and the number of the bytes read (> 0).
// If an error occured, the value is an empty string and the number of bytes n is <= meaning:
//	n == 0: buf too small
//	n < 0: invalid string length, or string length overflowed,
//	       and -n is the number of bytes read
// If buf contains too few string than expected, String reads string from buf to the end.
func String(buf []byte) (string, int) {
	l, n := binary.Uvarint(buf)
	if n <= 0 || l > uint64(maxint) {
		return "", n
	}
	if len(buf[n:]) < int(l) {
		return string(buf[n:]), -(n + len(buf[n:]))
	}
	return string(buf[n : n+int(l)]), n + int(l)
}

// PutString encodes a string into buf and returns the number of bytes written.
// If the buffer is too small to contain string length, PutString will panic.
// Else PutString will put string into buf as much as possible.
func PutString(buf []byte, s string) int {
	n := binary.PutUvarint(buf, uint64(len([]byte(s))))
	if n < 0 {
		return n
	}
	return n + copy(buf[n:], s)
}
