// +build ignore

// Package binary provides wrapper functions to
// read/write common data types from/to bytes.Buffer object.
package binary

import (
	"bytes"
	"unsafe"
)

// ReadTriad reads 3-byte big-endian number from buf.
func ReadTriad(buf *bytes.Buffer) (n uint32, err error) {
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n = uint32(b) << 16
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b)
	return
}

// WriteTriad writes 3-byte big-endian number to buf.
func WriteTriad(buf *bytes.Buffer, n uint32) (err error) {
	if err = buf.WriteByte(byte(n >> 16)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 8)); err != nil {
		return
	}
	err = buf.WriteByte(byte(n))
	return
}

// ReadLTriad reads 3-byte little-endian number from buf.
func ReadLTriad(buf *bytes.Buffer) (n uint32, err error) {
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n = uint32(b)
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b) << 16
	return
}

// WriteLTriad writes 3-byte little-endian number to buf.
func WriteLTriad(buf *bytes.Buffer, n uint32) (err error) {
	if err = buf.WriteByte(byte(n)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 8)); err != nil {
		return
	}
	err = buf.WriteByte(byte(n >> 16))
	return
}

// ReadBool reads boolean from buf.
func ReadBool(buf *bytes.Buffer) (b bool, err error) {
	var t byte
	if t, err = buf.ReadByte(); err != nil {
		return
	}
	return t != 0, err
}

// WriteBool writes boolean to buf.
func WriteBool(buf *bytes.Buffer, b bool) (err error) {
	if b {
		return buf.WriteByte(1)
	}
	return buf.WriteByte(0)
}

// ReadByte reads unsigned byte from buf.
func ReadByte(buf *bytes.Buffer) (b byte, err error) {
	b, err = buf.ReadByte()
	return
}

// WriteByte writes unsigned byte to buf.
func WriteByte(buf *bytes.Buffer, b byte) (err error) {
	return buf.WriteByte(b)
}

// ReadShort reads 16-bit signed big-endian number.
func ReadShort(buf *bytes.Buffer) (n int16, err error) {
	var u uint16
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u = uint16(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint16(b)
	return *(*int16)(unsafe.Pointer(&u)), nil
}

// WriteShort writes 16-bit signed big-endian number.
func WriteShort(buf *bytes.Buffer, n int16) (err error) {
	var u uint16 = *(*uint16)(unsafe.Pointer(&n))
	if err = buf.WriteByte(byte(u >> 8)); err != nil {
		return
	}
	return buf.WriteByte(byte(u))
}

// ReadShortUnsigned reads 16-bit unsigned big-endian number.
func ReadShortUnsigned(buf *bytes.Buffer) (n uint16, err error) {
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n = uint16(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint16(b)
	return
}

// WriteShortUnsigned writes 16-bit unsigned big-endian number.
func WriteShortUnsigned(buf *bytes.Buffer, n uint16) (err error) {
	if err = buf.WriteByte(byte(n >> 8)); err != nil {
		return
	}
	return buf.WriteByte(byte(n))
}

// ReadLShort reads 16-bit signed little-endian number.
func ReadLShort(buf *bytes.Buffer) (n int16, err error) {
	var u uint16
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u = uint16(b)
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint16(b) << 8
	return *(*int16)(unsafe.Pointer(&u)), nil
}

// WriteLShort writes 16-bit signed little-endian number.
func WriteLShort(buf *bytes.Buffer, n int16) (err error) {
	var u uint16 = *(*uint16)(unsafe.Pointer(&n))
	if err = buf.WriteByte(byte(u)); err != nil {
		return
	}
	return buf.WriteByte(byte(u >> 8))
}

// ReadLShortUnsigned reads 16-bit unsigned little-endian number.
func ReadLShortUnsigned(buf *bytes.Buffer) (n uint16, err error) {
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n = uint16(b)
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint16(b) << 8
	return
}

// WriteLShortUnsigned writes 16-bit unsigned little-endian number.
func WriteLShortUnsigned(buf *bytes.Buffer, n uint16) (err error) {
	if err = buf.WriteByte(byte(n)); err != nil {
		return
	}
	return buf.WriteByte(byte(n >> 8))
}

// ReadInt reads 32-bit signed big-endian number.
func ReadInt(buf *bytes.Buffer) (n int32, err error) {
	var u uint32
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u = uint32(b) << 24
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint32(b) << 16
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint32(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint32(b)
	return *(*int32)(unsafe.Pointer(&u)), nil
}

// WriteInt writes 32-bit signed big-endian number.
func WriteInt(buf *bytes.Buffer, n int32) (err error) {
	var u uint32 = *(*uint32)(unsafe.Pointer(&n))
	if err = buf.WriteByte(byte(u >> 24)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 16)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 8)); err != nil {
		return
	}
	return buf.WriteByte(byte(u))
}

// ReadLInt reads 32-bit signed little-endian number.
func ReadLInt(buf *bytes.Buffer) (n int32, err error) {
	var u uint32
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u = uint32(b)
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint32(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint32(b) << 16
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint32(b) << 24
	return *(*int32)(unsafe.Pointer(&u)), nil
}

// WriteLInt writes 32-bit signed little-endian number.
func WriteLInt(buf *bytes.Buffer, n int32) (err error) {
	var u uint32 = *(*uint32)(unsafe.Pointer(&n))
	if err = buf.WriteByte(byte(u)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 8)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 16)); err != nil {
		return
	}
	return buf.WriteByte(byte(u >> 24))
}

// ReadIntUnsigned reads 32-bit unsigned big-endian number.
func ReadIntUnsigned(buf *bytes.Buffer) (n uint32, err error) {
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n = uint32(b) << 24
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b) << 16
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b)
	return
}

// WriteIntUnsigned writes 32-bit unsigned big-endian number.
func WriteIntUnsigned(buf *bytes.Buffer, n uint32) (err error) {
	if err = buf.WriteByte(byte(n >> 24)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 16)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 8)); err != nil {
		return
	}
	return buf.WriteByte(byte(n))
}

// ReadLIntUnsigned reads 32-bit unsigned little-endian number.
func ReadLIntUnsigned(buf *bytes.Buffer) (n uint32, err error) {
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n = uint32(b)
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b) << 16
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint32(b) << 24
	return
}

// WriteLIntUnsigned writes 32-bit unsigned little-endian number.
func WriteLIntUnsigned(buf *bytes.Buffer, n uint32) (err error) {
	if err = buf.WriteByte(byte(n)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 8)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 16)); err != nil {
		return
	}
	return buf.WriteByte(byte(n >> 24))
}

// ReadLong reads 64-bit signed big-endian number.
func ReadLong(buf *bytes.Buffer) (n int64, err error) {
	var u uint64
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u = uint64(b) << 56
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 48
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 40
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 32
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 24
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 16
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b)
	return *(*int64)(unsafe.Pointer(&u)), nil
}

// WriteLong writes 64-bit signed big-endian number.
func WriteLong(buf *bytes.Buffer, n int64) (err error) {
	var u uint64 = *(*uint64)(unsafe.Pointer(&n))
	if err = buf.WriteByte(byte(u >> 56)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 48)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 40)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 32)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 24)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 16)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 8)); err != nil {
		return
	}
	return buf.WriteByte(byte(u))
}

// ReadLLong reads 64-bit signed little-endian number.
func ReadLLong(buf *bytes.Buffer) (n int64, err error) {
	var u uint64
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u = uint64(b)
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 16
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 24
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 32
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 40
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 48
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	u += uint64(b) << 56
	return *(*int64)(unsafe.Pointer(&u)), nil
}

// WriteLLong writes 64-bit signed little-endian number.
func WriteLLong(buf *bytes.Buffer, n int64) (err error) {
	var u uint64 = *(*uint64)(unsafe.Pointer(&n))
	if err = buf.WriteByte(byte(u)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 8)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 16)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 24)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 32)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 40)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(u >> 48)); err != nil {
		return
	}
	return buf.WriteByte(byte(u >> 56))
}

// ReadLongUnsigned reads 64-bit unsigned big-endian number.
func ReadLongUnsigned(buf *bytes.Buffer) (n uint64, err error) {
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n = uint64(b) << 56
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 48
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 40
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 32
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 24
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 16
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b)
	return
}

// WriteLongUnsigned writes 64-bit unsigned big-endian number.
func WriteLongUnsigned(buf *bytes.Buffer, n uint64) (err error) {
	if err = buf.WriteByte(byte(n >> 56)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 48)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 40)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 32)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 24)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 16)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 8)); err != nil {
		return
	}
	return buf.WriteByte(byte(n))
}

// ReadLLongUnsigned reads 64-bit unsigned little-endian number.
func ReadLLongUnsigned(buf *bytes.Buffer) (n uint64, err error) {
	var b byte
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n = uint64(b)
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 8
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 16
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 24
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 32
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 40
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 48
	if b, err = buf.ReadByte(); err != nil {
		return
	}
	n += uint64(b) << 56
	return n, nil
}

// WriteLLongUnsigned writes 64-bit unsigned little-endian number.
func WriteLLongUnsigned(buf *bytes.Buffer, n uint64) (err error) {
	if err = buf.WriteByte(byte(n)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 8)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 16)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 24)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 32)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 40)); err != nil {
		return
	}
	if err = buf.WriteByte(byte(n >> 48)); err != nil {
		return
	}
	return buf.WriteByte(byte(n >> 56))
}
