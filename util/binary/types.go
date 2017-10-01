package binary

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// Triad is a 3-byte unsigned, big-endian integer.
type Triad uint32

// MarshalStream implements Marshaler interface.
func (x Triad) MarshalStream(buf *bytes.Buffer) error {
	b := make([]byte, 3)
	BigEndian.PutTriad(b, uint32(x))
	_, err := buf.Write(b)
	return err
}

// UnmarshalStream implements Unmarshaler interface.
func (x *Triad) UnmarshalStream(buf *bytes.Buffer) error {
	b := make([]byte, 3)
	if _, err := buf.Read(b); err != nil {
		return err
	}
	*x = Triad(BigEndian.Triad(b))
	return nil
}

// LTriad is a 3-byte unsigned, little-endian integer.
type LTriad uint32

// MarshalStream implements Marshaler interface.
func (x LTriad) MarshalStream(buf *bytes.Buffer) error {
	b := make([]byte, 3)
	LittleEndian.PutTriad(b, uint32(x))
	_, err := buf.Write(b)
	return err
}

// UnmarshalStream implements Unmarshaler interface.
func (x *LTriad) UnmarshalStream(buf *bytes.Buffer) error {
	b := make([]byte, 3)
	if _, err := buf.Read(b); err != nil {
		return err
	}
	*x = LTriad(LittleEndian.Triad(b))
	return nil
}

// MCString is a common string type used in minecraft protocol.
type MCString string

// MarshalStream implements Marshaler interface.
func (s MCString) MarshalStream(buf *bytes.Buffer) error {
	b := make([]byte, len(s)+32)
	var n int
	if n = PutString(b, string(s)); n < 0 {
		return errors.New("PutString failed: PutString returned " + string(n))
	}
	if _, err := buf.Write(b[:n]); err != nil {
		return err
	}
	return nil
}

// UnmarshalStream implements Unmarshaler interface.
func (s *MCString) UnmarshalStream(buf *bytes.Buffer) error {
	var length uint64
	var err error
	if length, err = binary.ReadUvarint(buf); err != nil {
		return err
	}
	b := make([]byte, length)
	if _, err := buf.Read(b); err != nil {
		return err
	}
	*s = MCString(b)
	return nil
}

// FixedMCString is a old-style MCString type with fixed-size length field.
type FixedMCString string

// MarshalStream implements Marshaler interface.
func (s FixedMCString) MarshalStream(buf *bytes.Buffer) (err error) {
	b := make([]byte, 2)
	BigEndian.PutInt16(b, int16(len(s)))
	if _, err = buf.Write(b); err != nil {
		return
	}
	_, err = buf.WriteString(string(s))
	return
}

// UnmarshalStream implements Unmarshaler interface.
func (s *FixedMCString) UnmarshalStream(buf *bytes.Buffer) (err error) {
	b := buf.Next(int(BigEndian.Int16(buf.Next(2))))
	*s = FixedMCString(b)
	return nil
}
