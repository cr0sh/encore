package raknet

import (
	"github.com/cr0sh/encore/util/binary"
	log "github.com/sirupsen/logrus"
	"io"
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
func (ep EncapsulatedPacket) MarshalStream(wr io.Writer) (err error) {
	flag := ep.Reliability << 5
	if ep.IsSplit {
		flag |= (1 << 4)
	}

	buf := make([]byte, 3, 20)
	buf[0] = flag
	binary.BigEndian.PutUint16(buf[1:3], uint16(len(ep.Payload)<<3))

	if ep.Reliability > 0 {
		if ep.Reliability >= 2 && ep.Reliability != 5 {
			b := make([]byte, 3)
			binary.LittleEndian.PutTriad(b, ep.MessageIndex)
			buf = append(buf, b...)
		}
		if ep.Reliability <= 4 && ep.Reliability != 2 {
			b := make([]byte, 4)
			binary.LittleEndian.PutTriad(b, ep.OrderIndex)
			b[3] = ep.OrderChannel
			buf = append(buf, b...)
		}
	}

	if ep.IsSplit {
		b := make([]byte, 10)
		binary.BigEndian.PutUint32(b, ep.SplitCount)
		binary.BigEndian.PutUint16(b[4:], ep.SplitID)
		binary.BigEndian.PutUint32(b[6:], ep.SplitIndex)
		buf = append(buf, b...)
	}

	if ep.headLen() != len(buf) {
		log.WithFields(log.Fields{ // TODO: Remove this
			"expected": ep.headLen(),
			"real":     len(buf),
		}).Warn("Incorrect EncapsulatedPacket header length")
	}

	if _, err = wr.Write(buf); err != nil {
		return
	}

	_, err = wr.Write(ep.Payload)
	return
}

// UnmarshalStream implements Stream Unmarshaler interface.
func (ep *EncapsulatedPacket) UnmarshalStream(rd io.Reader) (err error) {
	b := make([]byte, 10)
	if _, err = rd.Read(b[:3]); err != nil {
		return
	}
	ep.Reliability = b[0] >> 5
	ep.IsSplit = b[0]&(1<<4) > 0

	payloadLen := binary.BigEndian.Uint16(b[1:3]) >> 3
	if b[2]&7 != 0 {
		payloadLen++
	}

	if ep.Reliability > 0 {
		if ep.Reliability >= 2 && ep.Reliability != 5 {
			if _, err = rd.Read(b[:3]); err != nil {
				return
			}
			ep.MessageIndex = binary.LittleEndian.Triad(b[:3])
		}
		if ep.Reliability <= 4 && ep.Reliability != 2 {
			if _, err = rd.Read(b[:4]); err != nil {
				return
			}
			ep.OrderIndex = binary.LittleEndian.Triad(b[:3])
			ep.OrderChannel = b[3]
		}
	}

	if ep.IsSplit {
		if _, err = rd.Read(b[:10]); err != nil {
			return
		}
		ep.SplitCount = binary.BigEndian.Uint32(b[:4])
		ep.SplitID = binary.BigEndian.Uint16(b[4:6])
		ep.SplitIndex = binary.BigEndian.Uint32(b[6:10])
	}

	ep.Payload = make([]byte, payloadLen)
	_, err = rd.Read(ep.Payload)
	return
}

// DataPacket is a set of EncapsulatedPackets with sequence number used in MCPE protocols.
type DataPacket struct {
	Seq     binary.LTriad
	Packets []EncapsulatedPacket
}

// MarshalStream implements Stream Marshaler interface.
func (dp DataPacket) MarshalStream(wr io.Writer) (err error) {
	b := make([]byte, 3)
	binary.LittleEndian.PutTriad(b, uint32(dp.Seq))
	if _, err = wr.Write(b); err != nil {
		return
	}
	for i := range dp.Packets {
		if err = dp.Packets[i].MarshalStream(wr); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalStream implements Stream Unmarshaler interface.
func (dp *DataPacket) UnmarshalStream(rd io.Reader) (err error) {
	b := make([]byte, 3)
	if _, err = rd.Read(b); err != nil {
		return
	}
	dp.Seq = binary.LTriad(binary.LittleEndian.Triad(b))

	for {
		ep := new(EncapsulatedPacket)
		if err = ep.UnmarshalStream(rd); err == io.EOF {
			return nil
		} else if err != nil {
			return
		}

		if len(ep.Payload) == 0 {
			// I don't know why PM/Nukkit implemented like this
			// ref https://github.com/Nukkit/Nukkit/blob/master/src/main/java/cn/nukkit/raknet/protocol/DataPacket.java#L45
			break
		}
		dp.Packets = append(dp.Packets, *ep)
	}

	return nil
}
