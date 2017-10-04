package raknet

import (
	"bytes"
	"github.com/cr0sh/encore/util/binary"
	"reflect"
	"testing"
)

func TestMarshalEncapsulatedPacket(t *testing.T) {
	cases := []struct {
		ep     EncapsulatedPacket
		expect []byte
	}{
		{
			ep: EncapsulatedPacket{
				Reliability: 0,
				Payload:     []byte("\x00\x01\x02\x03"),
			},
			expect: []byte("\x00\x00\x20\x00\x01\x02\x03"),
		},
		{
			ep: EncapsulatedPacket{
				Reliability:  2,
				MessageIndex: 10,
				Payload:      []byte("\x00\x12\x45"),
			},
			expect: []byte("\x40\x00\x18\x0a\x00\x00\x00\x12\x45"),
		},
		{
			ep: EncapsulatedPacket{
				Reliability:  2,
				MessageIndex: 16,

				IsSplit:    true,
				SplitCount: 10,
				SplitID:    3,
				SplitIndex: 1,

				Payload: []byte("\x00\x01\x02\x03"),
			},
			expect: []byte("\x50\x00\x20\x10\x00\x00\x00\x00\x00\x0a\x00\x03\x00\x00\x00\x01\x00\x01\x02\x03"),
		},
	}

	for i, c := range cases {
		buf := new(bytes.Buffer)
		if err := binary.Marshal(c.ep, buf); err != nil {
			t.Errorf("Test $%d: Marshal returned error %v", i, err)
			return
		}
		if !reflect.DeepEqual(c.expect, buf.Bytes()) {
			t.Errorf("Test #%d: Expected %v,\nGot %v", i, c.expect, buf.Bytes())
			return
		}
	}
}

func TestUnmarshalEncapsulatedPacket(t *testing.T) {
	cases := []struct {
		payload []byte
		expect  EncapsulatedPacket
	}{
		{
			payload: []byte("\x00\x00\x20\x00\x01\x02\x03"),
			expect: EncapsulatedPacket{
				Reliability: 0,
				Payload:     []byte("\x00\x01\x02\x03"),
			},
		},
		{
			payload: []byte("\x40\x00\x18\x0a\x00\x00\x00\x12\x45"),
			expect: EncapsulatedPacket{
				Reliability:  2,
				MessageIndex: 10,
				Payload:      []byte("\x00\x12\x45"),
			},
		},
		{
			payload: []byte("\x50\x00\x20\x10\x00\x00\x00\x00\x00\x0a\x00\x03\x00\x00\x00\x01\x00\x01\x02\x03"),
			expect: EncapsulatedPacket{
				Reliability:  2,
				MessageIndex: 16,

				IsSplit:    true,
				SplitCount: 10,
				SplitID:    3,
				SplitIndex: 1,

				Payload: []byte("\x00\x01\x02\x03"),
			},
		},
		{
			payload: []byte("\x90\x00\x20" +
				"\x03\x00\x00" +
				"\x02\x00\x00\x0a" +
				"\x00\x00\x00\x18\x00\x12\x00\x00\x00\x13" +
				"\x02\x03\x12\x11"),
			expect: EncapsulatedPacket{
				Reliability:  4,
				MessageIndex: 3,
				OrderIndex:   2,
				OrderChannel: 10,

				IsSplit:    true,
				SplitCount: 24,
				SplitID:    18,
				SplitIndex: 19,

				Payload: []byte("\x02\x03\x12\x11"),
			},
		},
	}

	for i, c := range cases {
		buf := bytes.NewBuffer(c.payload)
		ep := EncapsulatedPacket{}
		if err := binary.Unmarshal(&ep, buf); err != nil {
			t.Errorf("Test $%d: Unmarshal returned error %v", i, err)
			return
		}
		if !reflect.DeepEqual(c.expect, ep) {
			t.Errorf("Test #%d: Expected %v,\nGot %v", i, c.expect, ep)
			return
		}
	}

}

func TestMarshalDataPacket(t *testing.T) {
	cases := []struct {
		dp     DataPacket
		expect []byte
	}{
		{
			dp: DataPacket{
				Seq: 1,
			},
			expect: []byte("\x01\x00\x00"),
		},
	}

	for i, c := range cases {
		buf := new(bytes.Buffer)
		if err := binary.Marshal(c.dp, buf); err != nil {
			t.Errorf("Test $%d: Marshal returned error %v", i, err)
			return
		}
		if !reflect.DeepEqual(c.expect, buf.Bytes()) {
			t.Errorf("Test #%d: Expected %v,\nGot %v", i, c.expect, buf.Bytes())
			return
		}
	}
}
