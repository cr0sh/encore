package binary

import (
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	cases := []struct {
		buf []byte
		s   string
		n   int
	}{
		{
			[]byte("\x01a"),
			"a",
			2,
		},
		{
			[]byte("\x02ab"),
			"ab",
			3,
		},
		{
			[]byte("\x03a"),
			"a",
			-2,
		},
		{
			[]byte(""),
			"",
			0,
		},
	}

	for i, c := range cases {
		ts, tn := String(c.buf)
		if ts != c.s || tn != c.n {
			t.Errorf("Test #%d: Expected s, n %v, %d\nGot s, n %v, %d", i, c.s, c.n, ts, tn)
		}
	}
}

func TestPutString(t *testing.T) {
	cases := []struct {
		buflen int
		s      string
		n      int
		buf    []byte
	}{
		{
			3,
			"a",
			2,
			[]byte("\x01a\x00"),
		},
		{
			12,
			"hello world",
			12,
			[]byte("\x0bhello world"),
		},
	}
	for i, c := range cases {
		buf := make([]byte, c.buflen)
		tn := PutString(buf, c.s)
		if tn != c.n {
			t.Errorf("Test #%d: Expected n %d\nGot n %d", i, c.n, tn)
		}
		if !reflect.DeepEqual(buf, c.buf) {
			t.Errorf("Test #%d: Expected buf %v\nGot buf %v", i, c.buf, buf)
		}
	}
}

func TestPutStringPanic(t *testing.T) {
	cases := []struct {
		buflen int
		s      string
	}{
		{
			0,
			"abcd",
		},
		{
			1,
			string(make([]byte, 128)),
		},
		{
			2,
			string(make([]byte, 1<<15)),
		},
		{
			3,
			string(make([]byte, 1<<23+2)),
		},
	}

	for i, c := range cases {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Test #%d did not panic", i)
				}
			}()
			buf := make([]byte, c.buflen)
			PutString(buf, c.s)
		}()
	}
}
