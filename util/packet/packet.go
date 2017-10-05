// Package packet provides tools for handling packets.
package packet

import (
	"bytes"
	"github.com/cr0sh/encore/util/binary"
)

// Packet is a interface required for each packet struct to qualify.
type Packet interface {
	ID() byte
}

// Marshal marshals packet with appropriate ID.
func Marshal(pk Packet, buf *bytes.Buffer) (err error) {
	if err = buf.WriteByte(pk.ID()); err != nil {
		return
	}
	return binary.Marshal(pk, buf)
}
