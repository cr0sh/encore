package raknet

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestPacketWindow(t *testing.T) {
	type input struct {
		order uint64
		ptr   unsafe.Pointer
	}

	ns := []int{1, 2, 3, 4, 5}
	cases := []struct {
		put    input
		expect []unsafe.Pointer
	}{
		{
			put: input{
				order: 0,
				ptr:   unsafe.Pointer(&ns[0]),
			},
			expect: []unsafe.Pointer{unsafe.Pointer(&ns[0])},
		},
		{
			put: input{
				order: 99999,
				ptr:   unsafe.Pointer(&ns[1]),
			},
			expect: nil,
		},
		{
			put: input{
				order: 0,
				ptr:   unsafe.Pointer(&ns[0]),
			},
			expect: nil,
		},
		{
			put: input{
				order: 4,
				ptr:   unsafe.Pointer(&ns[4]),
			},
			expect: []unsafe.Pointer{},
		},

		{
			put: input{
				order: 2,
				ptr:   unsafe.Pointer(&ns[2]),
			},
			expect: []unsafe.Pointer{},
		},
		{
			put: input{
				order: 3,
				ptr:   unsafe.Pointer(&ns[3]),
			},
			expect: []unsafe.Pointer{},
		},
		{
			put: input{
				order: 1,
				ptr:   unsafe.Pointer(&ns[1]),
			},
			expect: []unsafe.Pointer{
				unsafe.Pointer(&ns[1]),
				unsafe.Pointer(&ns[2]),
				unsafe.Pointer(&ns[3]),
				unsafe.Pointer(&ns[4]),
			},
		},
	}

	window := new(PacketWindow).Init(true)
	for i, c := range cases {
		if ret := window.Put(c.put.order, c.put.ptr); !reflect.DeepEqual(ret, c.expect) {
			t.Errorf("Test #%d: expected %v,\ngot %v", i, c.expect, ret)
			return
		}
		switch i {
		case 3:
			if !reflect.DeepEqual(window.missing, map[uint64]struct{}{
				1: struct{}{},
				2: struct{}{},
				3: struct{}{},
			}) {
				t.Errorf("Test 3: GetMissing() mismatch: got %v", window.missing)
				return
			}
		case 4:
			if !reflect.DeepEqual(window.missing, map[uint64]struct{}{
				1: struct{}{},
				3: struct{}{},
			}) {
				t.Errorf("Test 4: GetMissing() mismatch: got %v", window.missing)
				return
			}
		case 5:
			if !reflect.DeepEqual(window.missing, map[uint64]struct{}{
				1: struct{}{},
			}) {
				t.Errorf("Test 5: GetMissing() mismatch: got %v", window.missing)
				return
			}
		}
	}
}
