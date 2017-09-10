package packet

import (
	"bytes"
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
