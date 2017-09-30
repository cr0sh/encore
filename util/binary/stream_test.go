package binary

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	cases := []struct {
		pack   interface{}
		expect []byte
	}{
		{
			struct {
				A byte
			}{
				'a',
			},
			[]byte{'a'},
		},
		{
			struct {
				A byte
				B uint32
				c int
				D uint
				e uint
				F Triad
				G int8
			}{
				'h',
				1,
				3,
				4,
				5,
				0xfefd09,
				-3,
			},
			[]byte{'h', 0, 0, 0, 1, 0xfe, 0xfd, 0x09, 0xfd},
		},
		{
			struct {
				A uint32 `stream:"pass"`
				B LTriad
				C float64
				D uint16 `stream:"little"`
			}{
				3,
				0xdeadff,
				math.Pi,
				0xfdaf,
			},
			[]byte{
				0xff, 0xad, 0xde, 0x40,
				0x09, 0x21, 0xfb, 0x54,
				0x44, 0x2d, 0x18, 0xaf,
				0xfd,
			},
		},
		{
			struct {
				A [3]struct {
					A LTriad
					b int
					C int32 `stream:"pass"`
					D int
				}
			}{

				[3]struct {
					A LTriad
					b int
					C int32 `stream:"pass"`
					D int
				}{
					{
						0xfefd09,
						1,
						2,
						3,
					},
					{
						0x010203,
						1,
						4,
						0xffff,
					},
					{
						0x123456,
						0xdead,
						0xbeef,
						0x1234,
					},
				},
			},
			[]byte{
				0x09, 0xfd, 0xfe,
				0x03, 0x02, 0x01,
				0x56, 0x34, 0x12,
			},
		},
		{
			struct {
				A []struct {
					A LTriad
					b int
					C int32 `stream:"pass"`
					D int
				}
			}{

				[]struct {
					A LTriad
					b int
					C int32 `stream:"pass"`
					D int
				}{
					{
						0xfefd09,
						1,
						2,
						3,
					},
					{
						0x010203,
						1,
						4,
						0xffff,
					},
				},
			},
			[]byte{
				0x09, 0xfd, 0xfe,
				0x03, 0x02, 0x01,
			},
		},
		{
			struct {
				A []byte
			}{
				[]byte("Hello guys!가나다라마바사아."),
			},
			[]byte("Hello guys!가나다라마바사아."),
		},
		{
			struct {
				A MCString
			}{
				"Hello!",
			},
			[]byte("\x06Hello!"),
		},
	}

	for i, c := range cases {
		var buf bytes.Buffer
		if err := Marshal(c.pack, &buf); err != nil {
			t.Errorf("Test #%d: Marshal returned error %v", i, err)
		}
		if !reflect.DeepEqual(c.expect, buf.Bytes()) {
			t.Errorf("Test #%d: Expected %v\nGot %v", i, c.expect, buf.Bytes())
		}
	}
}

func TestUnmarshal(t *testing.T) {
	cases := []struct {
		b      []byte
		expect interface{}
	}{
		{
			[]byte{0, 0, 1, 2},
			struct {
				A uint32
			}{
				0x00000102,
			},
		},
		{
			[]byte{0, 0, 1, 2, 0xfa},
			struct {
				A uint32
				B byte
			}{
				0x00000102,
				0xfa,
			},
		},
		{
			[]byte{0, 0, 0, 3},
			struct {
				A struct {
					A int32
				}
			}{
				struct {
					A int32
				}{
					3,
				},
			},
		},
		{
			[]byte("hello! 안녕하세요!"),
			struct {
				A []byte
			}{
				[]byte("hello! 안녕하세요!"),
			},
		},
		{
			[]byte{1, 0, 0},
			struct {
				A uint
				B int
				c uint32
				D LTriad
			}{
				D: 1,
			},
		},
		{
			[]byte("\x06Hello!"),
			struct {
				A MCString
			}{
				"Hello!",
			},
		},
		{
			[]byte("\x01\x02\x03\x04"),
			struct {
				A [3]byte
			}{
				[3]byte{1, 2, 3},
			},
		},
	}

	for i, c := range cases {
		buf := bytes.NewBuffer(c.b)
		v := reflect.New(reflect.TypeOf(c.expect))
		switch i { // case-specific initializations
		case 3:
			v.Elem().Field(0).Set(reflect.ValueOf(make([]byte, len(c.expect.(struct {
				A []byte
			}).A))))
		}
		if err := Unmarshal(v.Interface(), buf); err != nil {
			t.Errorf("Test #%d: Unmarshal returned error %v", i, err)
		}
		if !reflect.DeepEqual(c.expect, v.Elem().Interface()) {
			t.Errorf("Test #%d: Expected %v\nGot %v", i, c.expect, v.Elem().Interface())
		}
	}
}
