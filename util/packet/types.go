package packet

import (
	"bytes"
	gobinary "encoding/binary"
	"errors"
	"github.com/cr0sh/encore/util/binary"
)

// Triad is a 3-byte unsigned, big-endian integer.
type Triad uint32

// MarshalPacket implements Marshaler interface.
func (x Triad) MarshalPacket(buf *bytes.Buffer) error {
	b := make([]byte, 3)
	binary.BigEndian.PutTriad(b, uint32(x))
	_, err := buf.Write(b)
	return err
}

// UnmarshalPacket implements Unmarshaler interface.
func (x *Triad) UnmarshalPacket(buf *bytes.Buffer) error {
	b := make([]byte, 3)
	if _, err := buf.Read(b); err != nil {
		return err
	}
	*x = Triad(binary.BigEndian.Triad(b))
	return nil
}

// LTriad is a 3-byte unsigned, little-endian integer.
type LTriad uint32

// MarshalPacket implements Marshaler interface.
func (x LTriad) MarshalPacket(buf *bytes.Buffer) error {
	b := make([]byte, 3)
	binary.LittleEndian.PutTriad(b, uint32(x))
	_, err := buf.Write(b)
	return err
}

// UnmarshalPacket implements Unmarshaler interface.
func (x *LTriad) UnmarshalPacket(buf *bytes.Buffer) error {
	b := make([]byte, 3)
	if _, err := buf.Read(b); err != nil {
		return err
	}
	*x = LTriad(binary.LittleEndian.Triad(b))
	return nil
}

// MCString is a common string type used in minecraft protocol.
type MCString string

// MarshalPacket implements Marshaler interface.
func (s MCString) MarshalPacket(buf *bytes.Buffer) error {
	b := make([]byte, len(s)+32)
	var n int
	if n = binary.PutString(b, string(s)); n < 0 {
		return errors.New("PutString failed: PutString returned " + string(n))
	}
	if _, err := buf.Write(b[:n]); err != nil {
		return err
	}
	return nil
}

// UnmarshalPacket implements Unmarshaler interface.
func (s *MCString) UnmarshalPacket(buf *bytes.Buffer) error {
	var length uint64
	var err error
	if length, err = gobinary.ReadUvarint(buf); err != nil {
		return err
	}
	b := make([]byte, length)
	if _, err := buf.Read(b); err != nil {
		return err
	}
	*s = MCString(b)
	return nil
}
