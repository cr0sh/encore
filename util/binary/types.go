package binary

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Triad is a 3-byte unsigned, big-endian integer.
type Triad uint32

// MarshalStream implements Marshaler interface.
func (x Triad) MarshalStream(wr io.Writer) error {
	b := make([]byte, 3)
	BigEndian.PutTriad(b, uint32(x))
	_, err := wr.Write(b)
	return err
}

// UnmarshalStream implements Unmarshaler interface.
func (x *Triad) UnmarshalStream(rd *bytes.Buffer) error {
	b := make([]byte, 3)
	if _, err := rd.Read(b); err != nil {
		return err
	}
	*x = Triad(BigEndian.Triad(b))
	return nil
}

// LTriad is a 3-byte unsigned, little-endian integer.
type LTriad uint32

// MarshalStream implements Marshaler interface.
func (x LTriad) MarshalStream(wr io.Writer) error {
	b := make([]byte, 3)
	LittleEndian.PutTriad(b, uint32(x))
	_, err := wr.Write(b)
	return err
}

// UnmarshalStream implements Unmarshaler interface.
func (x *LTriad) UnmarshalStream(rd io.Reader) error {
	b := make([]byte, 3)
	if _, err := rd.Read(b); err != nil {
		return err
	}
	*x = LTriad(LittleEndian.Triad(b))
	return nil
}

type byteReaderHelper struct {
	io.Reader
	b []byte
}

func (h byteReaderHelper) ReadByte() (_ byte, err error) {
	if _, err = h.Reader.Read(h.b); err != nil {
		return
	}
	return h.b[0], nil
}

// MCString is a common string type used in minecraft protocol.
type MCString string

// MarshalStream implements Marshaler interface.
func (s MCString) MarshalStream(wr io.Writer) error {
	b := make([]byte, len(s)+32)
	var n int
	if n = PutString(b, string(s)); n < 0 {
		return errors.New("PutString failed: PutString returned " + string(n))
	}
	if _, err := wr.Write(b[:n]); err != nil {
		return err
	}
	return nil
}

// UnmarshalStream implements Unmarshaler interface.
func (s *MCString) UnmarshalStream(rd io.Reader) error {
	var length uint64
	var err error
	if length, err = binary.ReadUvarint(byteReaderHelper{rd, make([]byte, 1)}); err != nil {
		return err
	}
	b := make([]byte, length)
	if _, err := rd.Read(b); err != nil {
		return err
	}
	*s = MCString(b)
	return nil
}

// FixedMCString is a old-style MCString type with fixed-size length field.
type FixedMCString string

// MarshalStream implements Marshaler interface.
func (s FixedMCString) MarshalStream(wr io.Writer) (err error) {
	b := make([]byte, 2)
	BigEndian.PutInt16(b, int16(len(s)))
	if _, err = wr.Write(b); err != nil {
		return
	}
	_, err = wr.Write([]byte(s))
	return
}

// UnmarshalStream implements Unmarshaler interface.
func (s *FixedMCString) UnmarshalStream(rd io.Reader) (err error) {
	b := make([]byte, 2)
	if _, err = rd.Read(b); err != nil {
		return
	}
	b_ := make([]byte, BigEndian.Int16(b))
	if _, err = rd.Read(b_); err != nil {
		return
	}
	*s = FixedMCString(b_)
	return nil
}
