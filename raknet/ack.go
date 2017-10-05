package raknet

import (
	"bytes"
	"github.com/cr0sh/encore/util/binary"
	"github.com/sirupsen/logrus"
	"io"
	"sort"
)

// ACKMap is a set type for saving ACK/NACK packet IDs.
type ACKMap = map[uint32]struct{}

// EncodeACK encodes given ACKMap to Writer.
func EncodeACK(ack ACKMap, wr io.Writer) error {
	var warned bool

	b := make([]byte, 7)
	keys := make([]int, len(ack))
	i := 0
	for k, _ := range ack {
		keys[i] = int(k)
		i++
	}
	sort.Ints(keys)

	buf := new(bytes.Buffer)
	keycnt := len(keys)
	records := 0

	if keycnt > 0 {
		ptr := 1
		start := keys[0]
		end := keys[0]

		for ptr < keycnt {
			current := keys[ptr]
			ptr++
			diff := current - end
			if diff == 1 {
				end = current
			} else if diff > 1 {
				if start == end {
					b[0] = 0x01
					binary.LittleEndian.PutTriad(b[1:4], uint32(start))
					end = current
					start = end
					if _, err := buf.Write(b[:4]); err != nil {
						return err
					}
				} else {
					b[0] = 0x00
					binary.LittleEndian.PutTriad(b[1:4], uint32(start))
					binary.LittleEndian.PutTriad(b[4:7], uint32(end))
					end = current
					start = end
					if _, err := buf.Write(b); err != nil {
						return err
					}
				}
				records++
			} else {
				if !warned {
					logrus.Warn("Duplicate ACK Key while encoding(maybe by a bad ACK/NACK Queue?)")
					warned = true
				}
			}
		}

		if start == end {
			b[0] = 0x01
			binary.LittleEndian.PutTriad(b[1:4], uint32(start))
			buf.Write(b[:4])
		} else {
			b[0] = 0x00
			binary.LittleEndian.PutTriad(b[1:4], uint32(start))
			binary.LittleEndian.PutTriad(b[4:7], uint32(end))
			buf.Write(b)
		}
	}

	binary.BigEndian.PutUint16(b[:2], uint16(records))
	if _, err := wr.Write(b[:2]); err != nil {
		return err
	}

	_, err := io.Copy(wr, buf)
	return err
}

// DecodeACK returns decoded list from reader.
func DecodeACK(rd io.Reader) ([]uint32, error) {
	b := make([]byte, 7)

	if _, err := rd.Read(b[:2]); err != nil {
		return nil, err
	}

	keyscnt := int(binary.BigEndian.Uint16(b))
	keys := make([]uint32, 0, keyscnt)

	for i := 0; i < keyscnt; i++ {
		if _, err := rd.Read(b[:1]); err != nil {
			return keys, err
		}
		if b[0] == 0 {
			if _, err := rd.Read(b[:6]); err != nil {
				return keys, err
			}
			start := binary.LittleEndian.Triad(b[:3])
			end := binary.LittleEndian.Triad(b[3:6])
			/*
				if end - start > 512 {
					end = start + 512
				}
			*/
			for j := start; j <= end; j++ {
				keys = append(keys, uint32(j))
			}
		} else {
			if _, err := rd.Read(b[:3]); err != nil {
				return keys, err
			}
			keys = append(keys, uint32(binary.LittleEndian.Triad(b[:3])))
		}
	}

	return keys, nil
}
