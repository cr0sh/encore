package raknet

import (
	"bytes"
	"github.com/cr0sh/encore/util/binary"
	log "github.com/sirupsen/logrus"
)

// EncapsulatedPacket is an advanced raknet packet format
// with splits, reliability improvements(e.g. ACK), etc.
type EncapsulatedPacket struct {
	Reliability byte

	IsSplit    bool
	SplitCount uint32
	SplitID    uint16
	SplitIndex uint32

	MessageIndex uint32 // Little-Endian Triad

	OrderIndex   uint32 // Little-Endian Triad
	OrderChannel byte

	Payload []byte
}

func (ep EncapsulatedPacket) headLen() (length int) {
	switch ep.Reliability {
	case 0, 5:
		length = 3
	case 1:
		length = 7
	case 3, 4:
		length = 10
	default:
		length = 6
	}
	if ep.IsSplit {
		return length + 10
	}
	return
}

// Len returns estimated size of Marshaled EncapsulatedPacket in bytes.
func (ep EncapsulatedPacket) Len() int {
	return ep.headLen() + len(ep.Payload)
}

// MarshalStream implements Stream Marshaler interface.
func (ep EncapsulatedPacket) MarshalStream(buf *bytes.Buffer) (err error) {
	flag := ep.Reliability << 5
	if ep.IsSplit {
		flag |= (1 << 4)
	}

	buf_ := make([]byte, 3, 20)
	buf_[0] = flag
	binary.BigEndian.PutUint16(buf_[1:3], uint16(len(ep.Payload)<<3))

	if ep.Reliability > 0 {
		if ep.Reliability >= 2 && ep.Reliability != 5 {
			b := make([]byte, 3)
			binary.LittleEndian.PutTriad(b, ep.MessageIndex)
			buf_ = append(buf_, b...)
		}
		if ep.Reliability <= 4 && ep.Reliability != 2 {
			b := make([]byte, 4)
			binary.LittleEndian.PutTriad(b, ep.OrderIndex)
			b[3] = ep.OrderChannel
			buf_ = append(buf_, b...)
		}
	}

	if ep.IsSplit {
		b := make([]byte, 10)
		binary.BigEndian.PutUint32(b, ep.SplitCount)
		binary.BigEndian.PutUint16(b[4:], ep.SplitID)
		binary.BigEndian.PutUint32(b[6:], ep.SplitIndex)
		buf_ = append(buf_, b...)
	}

	if ep.headLen() != len(buf_) {
		log.WithFields(log.Fields{
			"expected": ep.headLen(),
			"real":     len(buf_),
		}).Warn("Incorrect EncapsulatedPacket header length")
	}

	if _, err = buf.Write(buf_); err != nil {
		return
	}

	_, err = buf.Write(ep.Payload)
	return
}

// UnmarshalStream implements Stream Unmarshaler interface.
func (ep *EncapsulatedPacket) UnmarshalStream(buf *bytes.Buffer) (err error) {
	b := make([]byte, 10)
	if _, err = buf.Read(b[:3]); err != nil {
		return
	}
	ep.Reliability = b[0] >> 5
	ep.IsSplit = b[0]&(1<<4) > 0

	payloadLen := binary.BigEndian.Uint16(b[:2]) >> 3
	if b[1]&7 != 0 {
		payloadLen++
	}

	if ep.Reliability > 0 {
		if ep.Reliability >= 2 && ep.Reliability != 5 {
			if _, err = buf.Read(b[:3]); err != nil {
				return
			}
			ep.MessageIndex = binary.LittleEndian.Triad(b[:3])
		}
		if ep.Reliability <= 4 && ep.Reliability != 2 {
			if _, err = buf.Read(b[:4]); err != nil {
				return
			}
			ep.OrderIndex = binary.LittleEndian.Triad(b[:3])
			ep.OrderChannel = b[3]
		}
	}

	if ep.IsSplit {
		if _, err = buf.Read(b[:10]); err != nil {
			return
		}
		ep.SplitCount = binary.BigEndian.Uint32(b[:4])
		ep.SplitID = binary.BigEndian.Uint16(b[4:6])
		ep.SplitIndex = binary.BigEndian.Uint32(b[6:10])
	}

	ep.Payload = make([]byte, payloadLen)
	_, err = buf.Read(ep.Payload)
	return
}

// DataPacket is a set of EncapsulatedPackets with sequence number used in MCPE protocols.
type DataPacket struct {
	Seq     uint32
	Packets []EncapsulatedPacket
}

// MarshalStream implements Stream Marshaler interface.
func (dp DataPacket) MarshalStream(buf *bytes.Buffer) (err error) {
	b := make([]byte, 3)
	binary.LittleEndian.PutTriad(b, dp.Seq)
	if _, err = buf.Write(b); err != nil {
		return
	}
	for i := range dp.Packets {
		if err = dp.Packets[i].MarshalStream(buf); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalStream implements Stream Unmarshaler interface.
func (dp *DataPacket) UnmarshalStream(buf *bytes.Buffer) (err error) {
	b := make([]byte, 3)
	if _, err = buf.Read(b); err != nil {
		return
	}
	dp.Seq = binary.LittleEndian.Triad(b)

	for !(buf.Len() > 0) {
		ep := new(EncapsulatedPacket)
		if err = ep.UnmarshalStream(buf); err != nil {
			return
		}

		if len(ep.Payload) == 0 {
			// I don't know Why
			// ref https://github.com/Nukkit/Nukkit/blob/master/src/main/java/cn/nukkit/raknet/protocol/DataPacket.java#L45
			break
		}
		dp.Packets = append(dp.Packets, *ep)
	}

	return nil
}
