// Copyright 2014 The MOAC-core Authors
// This file is part of the MOAC-core library.
//
// The MOAC-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The MOAC-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the MOAC-core library. If not, see <http://www.gnu.org/licenses/>.

package rlp

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"Chain3Go/lib/common"
)

func TestStreamKind(t *testing.T) {
	tests := []struct {
		input    string
		wantKind Kind
		wantLen  uint64
	}{
		// {"00", Byte, 0},
		// {"01", Byte, 0},
		// {"7F", Byte, 0},
		// {"80", String, 0},
		// {"B7", String, 55},
		// {"B90400", String, 1024},
		// {"BFFFFFFFFFFFFFFFFF", String, ^uint64(0)},
		// {"C0", List, 0},
		// {"C8", List, 8},
		// {"F7", List, 55},
		// {"F90400", List, 1024},
		// {"FFFFFFFFFFFFFFFFFF", List, ^uint64(0)},
		{"f87f01808a3039353032663930303086323030623230947312f4b8a4457a36827f185325fd6b66a3f8bb8b9030313633343537383564386130303030008025a0ea4e8768cd1d60ed59dc4e3fd7a371d0edb166fa057bc83c0585ed30e1293c41a04c349910f262b5075de9258f3978b119f14c93f5b5b8a8d1a2af22f0c14ed2f7", List, 127},
	}

	for i, test := range tests {
		// using plainReader to inhibit input limit errors.
		s := NewStream(newPlainReader(unhex(test.input)), 0)
		kind, len, err := s.Kind()
		if err != nil {
			t.Errorf("test %d: Kind returned error: %v", i, err)
			continue
		}
		if kind != test.wantKind {
			t.Errorf("test %d: kind mismatch: got %d, want %d", i, kind, test.wantKind)
		}
		if len != test.wantLen {
			t.Errorf("test %d: len mismatch: got %d, want %d", i, len, test.wantLen)
		}
	}
}

/*
 * This is to test an input with end
 * { from: '0xa8863fc8Ce3816411378685223C03DAae9770ebB',
        //  nonce: '0x16',
        //  gasPrice: '0x09502f9000',
        //  gasLimit: '0x5208',
        //  to: '0x7312F4B8A4457a36827f185325Fd6B66a3f8BB8B',
        //  value: '0x016345785d8a0000',
        //  data: '0x00' }

*/
func TestInputCommandStream(t *testing.T) {
	incmd := "f86c168509502f9000825208947312f4b8a4457a36827f185325fd6b66a3f8bb8b88016345785d8a0000001ba0c515be40cdaed29997219eb5fb1002ac03e35b8666ce84c3b9e181ad13760ea4a01eb2ded06283fb4663652120f031b5b0cef06b9cf0a831d3c2b5a158460790c1"
	ls := NewListStream(bytes.NewReader(unhex(incmd)), 108)
	if k, size, err := ls.Kind(); k != List || size != 108 || err != nil {
		t.Errorf("Kind() returned (%v, %d, %v), expected (List, 3, nil)", k, size, err)
	}

	if size, err := ls.List(); size != 108 || err != nil {
		t.Errorf("List() returned (%d, %v), expected (108, nil)", size, err)
	}

	for i := 0; i < 3; i++ {
		if val, err := ls.Uint(); val != 1 || err != nil {
			t.Errorf("Uint() returned (%d, %v), expected (1, nil)", val, err)
		}
	}
	if err := ls.ListEnd(); err != nil {
		t.Errorf("ListEnd() returned %v, expected (3, nil)", err)
	}
}

func TestNewListStream(t *testing.T) {
	ls := NewListStream(bytes.NewReader(unhex("0101010101")), 3)
	if k, size, err := ls.Kind(); k != List || size != 3 || err != nil {
		t.Errorf("Kind() returned (%v, %d, %v), expected (List, 3, nil)", k, size, err)
	}
	if size, err := ls.List(); size != 3 || err != nil {
		t.Errorf("List() returned (%d, %v), expected (3, nil)", size, err)
	}
	for i := 0; i < 3; i++ {
		if val, err := ls.Uint(); val != 1 || err != nil {
			t.Errorf("Uint() returned (%d, %v), expected (1, nil)", val, err)
		}
	}
	if err := ls.ListEnd(); err != nil {
		t.Errorf("ListEnd() returned %v, expected (3, nil)", err)
	}
}

func TestStreamErrors(t *testing.T) {
	withoutInputLimit := func(b []byte) *Stream {
		return NewStream(newPlainReader(b), 0)
	}
	withCustomInputLimit := func(limit uint64) func([]byte) *Stream {
		return func(b []byte) *Stream {
			return NewStream(bytes.NewReader(b), limit)
		}
	}

	type calls []string
	tests := []struct {
		string
		calls
		newStream func([]byte) *Stream // uses bytes.Reader if nil
		error     error
	}{
		{"C0", calls{"Bytes"}, nil, ErrExpectedString},
		{"C0", calls{"Uint"}, nil, ErrExpectedString},
		{"89000000000000000001", calls{"Uint"}, nil, errUintOverflow},
		{"00", calls{"List"}, nil, ErrExpectedList},
		{"80", calls{"List"}, nil, ErrExpectedList},
		{"C0", calls{"List", "Uint"}, nil, EOL},
		{"C8C9010101010101010101", calls{"List", "Kind"}, nil, ErrElemTooLarge},
		{"C3C2010201", calls{"List", "List", "Uint", "Uint", "ListEnd", "Uint"}, nil, EOL},
		{"00", calls{"ListEnd"}, nil, errNotInList},
		{"C401020304", calls{"List", "Uint", "ListEnd"}, nil, errNotAtEOL},

		// Non-canonical integers (e.g. leading zero bytes).
		{"00", calls{"Uint"}, nil, ErrCanonInt},
		{"820002", calls{"Uint"}, nil, ErrCanonInt},
		{"8133", calls{"Uint"}, nil, ErrCanonSize},
		{"817F", calls{"Uint"}, nil, ErrCanonSize},
		{"8180", calls{"Uint"}, nil, nil},

		// Non-valid boolean
		{"02", calls{"Bool"}, nil, errors.New("rlp: invalid boolean value: 2")},

		// Size tags must use the smallest possible encoding.
		// Leading zero bytes in the size tag are also rejected.
		{"8100", calls{"Uint"}, nil, ErrCanonSize},
		{"8100", calls{"Bytes"}, nil, ErrCanonSize},
		{"8101", calls{"Bytes"}, nil, ErrCanonSize},
		{"817F", calls{"Bytes"}, nil, ErrCanonSize},
		{"8180", calls{"Bytes"}, nil, nil},
		{"B800", calls{"Kind"}, withoutInputLimit, ErrCanonSize},
		{"B90000", calls{"Kind"}, withoutInputLimit, ErrCanonSize},
		{"B90055", calls{"Kind"}, withoutInputLimit, ErrCanonSize},
		{"BA0002FFFF", calls{"Bytes"}, withoutInputLimit, ErrCanonSize},
		{"F800", calls{"Kind"}, withoutInputLimit, ErrCanonSize},
		{"F90000", calls{"Kind"}, withoutInputLimit, ErrCanonSize},
		{"F90055", calls{"Kind"}, withoutInputLimit, ErrCanonSize},
		{"FA0002FFFF", calls{"List"}, withoutInputLimit, ErrCanonSize},

		// Expected EOF
		{"", calls{"Kind"}, nil, io.EOF},
		{"", calls{"Uint"}, nil, io.EOF},
		{"", calls{"List"}, nil, io.EOF},
		{"8180", calls{"Uint", "Uint"}, nil, io.EOF},
		{"C0", calls{"List", "ListEnd", "List"}, nil, io.EOF},

		{"", calls{"List"}, withoutInputLimit, io.EOF},
		{"8180", calls{"Uint", "Uint"}, withoutInputLimit, io.EOF},
		{"C0", calls{"List", "ListEnd", "List"}, withoutInputLimit, io.EOF},

		// Input limit errors.
		{"81", calls{"Bytes"}, nil, ErrValueTooLarge},
		{"81", calls{"Uint"}, nil, ErrValueTooLarge},
		{"81", calls{"Raw"}, nil, ErrValueTooLarge},
		{"BFFFFFFFFFFFFFFFFFFF", calls{"Bytes"}, nil, ErrValueTooLarge},
		{"C801", calls{"List"}, nil, ErrValueTooLarge},

		// Test for list element size check overflow.
		{"CD04040404FFFFFFFFFFFFFFFFFF0303", calls{"List", "Uint", "Uint", "Uint", "Uint", "List"}, nil, ErrElemTooLarge},

		// Test for input limit overflow. Since we are counting the limit
		// down toward zero in Stream.remaining, reading too far can overflow
		// remaining to a large value, effectively disabling the limit.
		{"C40102030401", calls{"Raw", "Uint"}, withCustomInputLimit(5), io.EOF},
		{"C4010203048180", calls{"Raw", "Uint"}, withCustomInputLimit(6), ErrValueTooLarge},

		// Check that the same calls are fine without a limit.
		{"C40102030401", calls{"Raw", "Uint"}, withoutInputLimit, nil},
		{"C4010203048180", calls{"Raw", "Uint"}, withoutInputLimit, nil},

		// Unexpected EOF. This only happens when there is
		// no input limit, so the reader needs to be 'dumbed down'.
		{"81", calls{"Bytes"}, withoutInputLimit, io.ErrUnexpectedEOF},
		{"81", calls{"Uint"}, withoutInputLimit, io.ErrUnexpectedEOF},
		{"BFFFFFFFFFFFFFFF", calls{"Bytes"}, withoutInputLimit, io.ErrUnexpectedEOF},
		{"C801", calls{"List", "Uint", "Uint"}, withoutInputLimit, io.ErrUnexpectedEOF},

		// This test verifies that the input position is advanced
		// correctly when calling Bytes for empty strings. Kind can be called
		// any number of times in between and doesn't advance.
		{"C3808080", calls{
			"List",  // enter the list
			"Bytes", // past first element

			"Kind", "Kind", "Kind", // this shouldn't advance

			"Bytes", // past second element

			"Kind", "Kind", // can't hurt to try

			"Bytes", // past final element
			"Bytes", // this one should fail
		}, nil, EOL},
	}

testfor:
	for i, test := range tests {
		if test.newStream == nil {
			test.newStream = func(b []byte) *Stream { return NewStream(bytes.NewReader(b), 0) }
		}
		s := test.newStream(unhex(test.string))
		rs := reflect.ValueOf(s)
		for j, call := range test.calls {
			fval := rs.MethodByName(call)
			ret := fval.Call(nil)
			err := "<nil>"
			if lastret := ret[len(ret)-1].Interface(); lastret != nil {
				err = lastret.(error).Error()
			}
			if j == len(test.calls)-1 {
				want := "<nil>"
				if test.error != nil {
					want = test.error.Error()
				}
				if err != want {
					t.Log(test)
					t.Errorf("test %d: last call (%s) error mismatch\ngot:  %s\nwant: %s",
						i, call, err, test.error)
				}
			} else if err != "<nil>" {
				t.Log(test)
				t.Errorf("test %d: call %d (%s) unexpected error: %q", i, j, call, err)
				continue testfor
			}
		}
	}
}

func TestStreamList(t *testing.T) {
	s := NewStream(bytes.NewReader(unhex("C80102030405060708")), 0)

	len, err := s.List()
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len != 8 {
		t.Fatalf("List returned invalid length, got %d, want 8", len)
	}

	for i := uint64(1); i <= 8; i++ {
		v, err := s.Uint()
		if err != nil {
			t.Fatalf("Uint error: %v", err)
		}
		if i != v {
			t.Errorf("Uint returned wrong value, got %d, want %d", v, i)
		}
	}

	if _, err := s.Uint(); err != EOL {
		t.Errorf("Uint error mismatch, got %v, want %v", err, EOL)
	}
	if err = s.ListEnd(); err != nil {
		t.Fatalf("ListEnd error: %v", err)
	}
}

func TestStreamRaw(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			"C58401010101",
			"8401010101",
		},
		{
			"F842B84001010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101",
			"B84001010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101",
		},
	}
	for i, tt := range tests {
		s := NewStream(bytes.NewReader(unhex(tt.input)), 0)
		s.List()

		want := unhex(tt.output)
		raw, err := s.Raw()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(want, raw) {
			t.Errorf("test %d: raw mismatch: got %x, want %x", i, raw, want)
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	r := bytes.NewReader(nil)

	if err := Decode(r, nil); err != errDecodeIntoNil {
		t.Errorf("Decode(r, nil) error mismatch, got %q, want %q", err, errDecodeIntoNil)
	}

	var nilptr *struct{}
	if err := Decode(r, nilptr); err != errDecodeIntoNil {
		t.Errorf("Decode(r, nilptr) error mismatch, got %q, want %q", err, errDecodeIntoNil)
	}

	if err := Decode(r, struct{}{}); err != errNoPointer {
		t.Errorf("Decode(r, struct{}{}) error mismatch, got %q, want %q", err, errNoPointer)
	}

	expectErr := "rlp: type chan bool is not RLP-serializable"
	if err := Decode(r, new(chan bool)); err == nil || err.Error() != expectErr {
		t.Errorf("Decode(r, new(chan bool)) error mismatch, got %q, want %q", err, expectErr)
	}

	if err := Decode(r, new(uint)); err != io.EOF {
		t.Errorf("Decode(r, new(int)) error mismatch, got %q, want %q", err, io.EOF)
	}
}

type decodeTest struct {
	input string
	ptr   interface{}
	value interface{}
	error string
}

type simplestruct struct {
	A uint
	B string
}

type recstruct struct {
	I     uint
	Child *recstruct `rlp:"nil"`
}

type invalidTail1 struct {
	A uint `rlp:"tail"`
	B string
}

type invalidTail2 struct {
	A uint
	B string `rlp:"tail"`
}

type tailRaw struct {
	A    uint
	Tail []RawValue `rlp:"tail"`
}

type tailUint struct {
	A    uint
	Tail []uint `rlp:"tail"`
}

var (
	veryBigInt = big.NewInt(0).Add(
		big.NewInt(0).Lsh(big.NewInt(0xFFFFFFFFFFFFFF), 16),
		big.NewInt(0xFFFF),
	)
)

type hasIgnoredField struct {
	A uint
	B uint `rlp:"-"`
	C uint
}

var decodeTests = []decodeTest{
	// booleans
	{input: "01", ptr: new(bool), value: true},
	{input: "80", ptr: new(bool), value: false},
	{input: "02", ptr: new(bool), error: "rlp: invalid boolean value: 2"},

	// integers
	{input: "05", ptr: new(uint32), value: uint32(5)},
	{input: "80", ptr: new(uint32), value: uint32(0)},
	{input: "820505", ptr: new(uint32), value: uint32(0x0505)},
	{input: "83050505", ptr: new(uint32), value: uint32(0x050505)},
	{input: "8405050505", ptr: new(uint32), value: uint32(0x05050505)},
	{input: "850505050505", ptr: new(uint32), error: "rlp: input string too long for uint32"},
	{input: "C0", ptr: new(uint32), error: "rlp: expected input string or byte for uint32"},
	{input: "00", ptr: new(uint32), error: "rlp: non-canonical integer (leading zero bytes) for uint32"},
	{input: "8105", ptr: new(uint32), error: "rlp: non-canonical size information for uint32"},
	{input: "820004", ptr: new(uint32), error: "rlp: non-canonical integer (leading zero bytes) for uint32"},
	{input: "B8020004", ptr: new(uint32), error: "rlp: non-canonical size information for uint32"},

	// slices
	{input: "C0", ptr: new([]uint), value: []uint{}},
	{input: "C80102030405060708", ptr: new([]uint), value: []uint{1, 2, 3, 4, 5, 6, 7, 8}},
	{input: "F8020004", ptr: new([]uint), error: "rlp: non-canonical size information for []uint"},

	// arrays
	{input: "C50102030405", ptr: new([5]uint), value: [5]uint{1, 2, 3, 4, 5}},
	{input: "C0", ptr: new([5]uint), error: "rlp: input list has too few elements for [5]uint"},
	{input: "C102", ptr: new([5]uint), error: "rlp: input list has too few elements for [5]uint"},
	{input: "C6010203040506", ptr: new([5]uint), error: "rlp: input list has too many elements for [5]uint"},
	{input: "F8020004", ptr: new([5]uint), error: "rlp: non-canonical size information for [5]uint"},

	// zero sized arrays
	{input: "C0", ptr: new([0]uint), value: [0]uint{}},
	{input: "C101", ptr: new([0]uint), error: "rlp: input list has too many elements for [0]uint"},

	// byte slices
	{input: "01", ptr: new([]byte), value: []byte{1}},
	{input: "80", ptr: new([]byte), value: []byte{}},
	{input: "8D6162636465666768696A6B6C6D", ptr: new([]byte), value: []byte("abcdefghijklm")},
	{input: "C0", ptr: new([]byte), error: "rlp: expected input string or byte for []uint8"},
	{input: "8105", ptr: new([]byte), error: "rlp: non-canonical size information for []uint8"},

	// byte arrays
	{input: "02", ptr: new([1]byte), value: [1]byte{2}},
	{input: "8180", ptr: new([1]byte), value: [1]byte{128}},
	{input: "850102030405", ptr: new([5]byte), value: [5]byte{1, 2, 3, 4, 5}},

	// byte array errors
	{input: "02", ptr: new([5]byte), error: "rlp: input string too short for [5]uint8"},
	{input: "80", ptr: new([5]byte), error: "rlp: input string too short for [5]uint8"},
	{input: "820000", ptr: new([5]byte), error: "rlp: input string too short for [5]uint8"},
	{input: "C0", ptr: new([5]byte), error: "rlp: expected input string or byte for [5]uint8"},
	{input: "C3010203", ptr: new([5]byte), error: "rlp: expected input string or byte for [5]uint8"},
	{input: "86010203040506", ptr: new([5]byte), error: "rlp: input string too long for [5]uint8"},
	{input: "8105", ptr: new([1]byte), error: "rlp: non-canonical size information for [1]uint8"},
	{input: "817F", ptr: new([1]byte), error: "rlp: non-canonical size information for [1]uint8"},

	// zero sized byte arrays
	{input: "80", ptr: new([0]byte), value: [0]byte{}},
	{input: "01", ptr: new([0]byte), error: "rlp: input string too long for [0]uint8"},
	{input: "8101", ptr: new([0]byte), error: "rlp: input string too long for [0]uint8"},

	// strings
	{input: "00", ptr: new(string), value: "\000"},
	{input: "8D6162636465666768696A6B6C6D", ptr: new(string), value: "abcdefghijklm"},
	{input: "C0", ptr: new(string), error: "rlp: expected input string or byte for string"},

	// big ints
	{input: "01", ptr: new(*big.Int), value: big.NewInt(1)},
	{input: "89FFFFFFFFFFFFFFFFFF", ptr: new(*big.Int), value: veryBigInt},
	{input: "10", ptr: new(big.Int), value: *big.NewInt(16)}, // non-pointer also works
	{input: "C0", ptr: new(*big.Int), error: "rlp: expected input string or byte for *big.Int"},
	{input: "820001", ptr: new(big.Int), error: "rlp: non-canonical integer (leading zero bytes) for *big.Int"},
	{input: "8105", ptr: new(big.Int), error: "rlp: non-canonical size information for *big.Int"},

	// structs
	{
		input: "C50583343434",
		ptr:   new(simplestruct),
		value: simplestruct{5, "444"},
	},
	{
		input: "C601C402C203C0",
		ptr:   new(recstruct),
		value: recstruct{1, &recstruct{2, &recstruct{3, nil}}},
	},

	// struct errors
	{
		input: "C0",
		ptr:   new(simplestruct),
		error: "rlp: too few elements for rlp.simplestruct",
	},
	{
		input: "C105",
		ptr:   new(simplestruct),
		error: "rlp: too few elements for rlp.simplestruct",
	},
	{
		input: "C7C50583343434C0",
		ptr:   new([]*simplestruct),
		error: "rlp: too few elements for rlp.simplestruct, decoding into ([]*rlp.simplestruct)[1]",
	},
	{
		input: "83222222",
		ptr:   new(simplestruct),
		error: "rlp: expected input list for rlp.simplestruct",
	},
	{
		input: "C3010101",
		ptr:   new(simplestruct),
		error: "rlp: input list has too many elements for rlp.simplestruct",
	},
	{
		input: "C501C3C00000",
		ptr:   new(recstruct),
		error: "rlp: expected input string or byte for uint, decoding into (rlp.recstruct).Child.I",
	},
	{
		input: "C0",
		ptr:   new(invalidTail1),
		error: "rlp: invalid struct tag \"tail\" for rlp.invalidTail1.A (must be on last field)",
	},
	{
		input: "C0",
		ptr:   new(invalidTail2),
		error: "rlp: invalid struct tag \"tail\" for rlp.invalidTail2.B (field type is not slice)",
	},
	{
		input: "C50102C20102",
		ptr:   new(tailUint),
		error: "rlp: expected input string or byte for uint, decoding into (rlp.tailUint).Tail[1]",
	},

	// struct tag "tail"
	{
		input: "C3010203",
		ptr:   new(tailRaw),
		value: tailRaw{A: 1, Tail: []RawValue{unhex("02"), unhex("03")}},
	},
	{
		input: "C20102",
		ptr:   new(tailRaw),
		value: tailRaw{A: 1, Tail: []RawValue{unhex("02")}},
	},
	{
		input: "C101",
		ptr:   new(tailRaw),
		value: tailRaw{A: 1, Tail: []RawValue{}},
	},

	// struct tag "-"
	{
		input: "C20102",
		ptr:   new(hasIgnoredField),
		value: hasIgnoredField{A: 1, C: 2},
	},

	// RawValue
	{input: "01", ptr: new(RawValue), value: RawValue(unhex("01"))},
	{input: "82FFFF", ptr: new(RawValue), value: RawValue(unhex("82FFFF"))},
	{input: "C20102", ptr: new([]RawValue), value: []RawValue{unhex("01"), unhex("02")}},

	// pointers
	{input: "00", ptr: new(*[]byte), value: &[]byte{0}},
	{input: "80", ptr: new(*uint), value: uintp(0)},
	{input: "C0", ptr: new(*uint), error: "rlp: expected input string or byte for uint"},
	{input: "07", ptr: new(*uint), value: uintp(7)},
	{input: "817F", ptr: new(*uint), error: "rlp: non-canonical size information for uint"},
	{input: "8180", ptr: new(*uint), value: uintp(0x80)},
	{input: "C109", ptr: new(*[]uint), value: &[]uint{9}},
	{input: "C58403030303", ptr: new(*[][]byte), value: &[][]byte{{3, 3, 3, 3}}},

	// check that input position is advanced also for empty values.
	{input: "C3808005", ptr: new([]*uint), value: []*uint{uintp(0), uintp(0), uintp(5)}},

	// interface{}
	{input: "00", ptr: new(interface{}), value: []byte{0}},
	{input: "01", ptr: new(interface{}), value: []byte{1}},
	{input: "80", ptr: new(interface{}), value: []byte{}},
	{input: "850505050505", ptr: new(interface{}), value: []byte{5, 5, 5, 5, 5}},
	{input: "C0", ptr: new(interface{}), value: []interface{}{}},
	{input: "C50183040404", ptr: new(interface{}), value: []interface{}{[]byte{1}, []byte{4, 4, 4}}},
	{
		input: "C3010203",
		ptr:   new([]io.Reader),
		error: "rlp: type io.Reader is not RLP-serializable",
	},

	// fuzzer crashes
	{
		input: "c330f9c030f93030ce3030303030303030bd303030303030",
		ptr:   new(interface{}),
		error: "rlp: element is larger than containing list",
	},
}

func uintp(i uint) *uint { return &i }

func runTests(t *testing.T, decode func([]byte, interface{}) error) {
	for i, test := range decodeTests {
		input, err := hex.DecodeString(test.input)
		if err != nil {
			t.Errorf("test %d: invalid hex input %q", i, test.input)
			continue
		}
		err = decode(input, test.ptr)
		if err != nil && test.error == "" {
			t.Errorf("test %d: unexpected Decode error: %v\ndecoding into %T\ninput %q",
				i, err, test.ptr, test.input)
			continue
		}
		if test.error != "" && fmt.Sprint(err) != test.error {

			t.Errorf("test %d: Decode error mismatch\ngot  %v\nwant %v\ndecoding into %T\ninput %q",
				i, err, test.error, test.ptr, test.input)
			continue
		}
		deref := reflect.ValueOf(test.ptr).Elem().Interface()
		if err == nil && !reflect.DeepEqual(deref, test.value) {
			t.Errorf("test %d: value mismatch\ngot  %#v\nwant %#v\ndecoding into %T\ninput %q",
				i, deref, test.value, test.ptr, test.input)
		}
	}
}

func TestDecodeWithByteReader(t *testing.T) {
	runTests(t, func(input []byte, into interface{}) error {
		return Decode(bytes.NewReader(input), into)
	})
}

// plainReader reads from a byte slice but does not
// implement ReadByte. It is also not recognized by the
// size validation. This is useful to test how the decoder
// behaves on a non-buffered input stream.
type plainReader []byte

func newPlainReader(b []byte) io.Reader {
	return (*plainReader)(&b)
}

func (r *plainReader) Read(buf []byte) (n int, err error) {
	if len(*r) == 0 {
		return 0, io.EOF
	}
	n = copy(buf, *r)
	*r = (*r)[n:]
	return n, nil
}

func TestDecodeWithNonByteReader(t *testing.T) {
	runTests(t, func(input []byte, into interface{}) error {
		return Decode(newPlainReader(input), into)
	})
}

func TestDecodeStreamReset(t *testing.T) {
	s := NewStream(nil, 0)
	runTests(t, func(input []byte, into interface{}) error {
		s.Reset(bytes.NewReader(input), 0)
		return s.Decode(into)
	})
}

type ethtxdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     *big.Int        `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

//2018/03/15
//Revised data format, only use SystemContract and Sharding
type txdata struct {
	AccountNonce   uint64          `json:"nonce"    gencodec:"required"`
	SystemContract uint64          `json:"syscnt" gencodec:"required"`
	Price          *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit       *big.Int        `json:"gas"      gencodec:"required"`
	Recipient      *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount         *big.Int        `json:"value"    gencodec:"required"`
	Payload        []byte          `json:"input"    gencodec:"required"`
	ShardingFlag   uint64          `json:"shardingFlag" gencodec:"required"`
	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

/*
 * Original data decoder
 */
func TestMoacDataDecoder(t *testing.T) {

	mctxformat := txdata{
		AccountNonce:   0,
		Recipient:      nil,
		Payload:        nil,
		Amount:         new(big.Int),
		GasLimit:       new(big.Int),
		Price:          new(big.Int),
		SystemContract: 0,
		ShardingFlag:   0,
		V:              new(big.Int),
		R:              new(big.Int),
		S:              new(big.Int),
	}
	// cmd := "f87f01808a3039353032663930303086323030623230947312f4b8a4457a36827f185325fd6b66a3f8bb8b9030313633343537383564386130303030008025a0ea4e8768cd1d60ed59dc4e3fd7a371d0edb166fa057bc83c0585ed30e1293c41a04c349910f262b5075de9258f3978b119f14c93f5b5b8a8d1a2af22f0c14ed2f7"
	//   cmd := "f87f01808a3039353032663930303086323030623230947312f4b8a4457a36827f185325fd6b66a3f8bb8b9030313633343537383564386130303030008025a0ea4e8768cd1d60ed59dc4e3fd7a371d0edb166fa057bc83c0585ed30e1293c41a04c349910f262b5075de9258f3978b119f14c93f5b5b8a8d1a2af22f0c14ed2f7"
	//
	// { from: '0xa8863fc8Ce3816411378685223C03DAae9770ebB',
	// 	nonce: 4,
	// 	gasPrice: 4000000,
	// 	gasLimit: 2000,
	// 	to: '0x7312F4B8A4457a36827f185325Fd6B66a3f8BB8B',
	// 	value: '01100000000000000000',
	// 	data: '0x12',
	// 	shardingFlag: 0 }
	// [ '0x04',
	// '0x',
	// '0x3d0900',
	// '0x07d0',
	// '0x7312f4b8a4457a36827f185325fd6b66a3f8bb8b',
	// '0x3031313030303030303030303030303030303030',
	// '0x12',
	// '0x',
	// '0x25',
	// '0x4ae57512a74c26b0a7c282d1241093b774a24e1c9f1a425d756860da23a7249d',
	// '0x16795a6164178063d3915a5f0e124648139ba11f714b19fc811277fc1f674c25' ]
	cmd := "f86f01808509502f900083200b20947312f4b8a4457a36827f185325fd6b66a3f8bb8b880f43fc2c04ee0000008025a0dfa37f170535100ff2c6081c326406fec05d7f9b314dfe256af0bfef6c3ba421a01ff2ce9f8fe5167da27d3de416ddc4844572ca88b9735c39307b430d1e13926a"
	addr := "0xa8863fc8Ce3816411378685223C03DAae9770ebB"
	if err := Decode(bytes.NewReader(unhex(cmd)), &mctxformat); err != nil {
		t.Fatalf("Decode error: %v", err)
	} else {
		fmt.Printf("Recipient:%v\n", mctxformat.Recipient.Hex())
		fmt.Printf("TX nonce: %v\n", mctxformat.AccountNonce)
		if addr != mctxformat.Recipient.Hex() {
			fmt.Printf("Unmatch address!%v vs %v\n", addr, mctxformat.Recipient)
		}

		if mctxformat.AccountNonce != 1 {
			fmt.Printf("TX:%v\n", mctxformat.Recipient)
		}
		fmt.Printf("Recipient:%v\n", mctxformat.Price)
	}

}

/*
 * This test can be used for changing transaction data format
 * as defined in Moac-lib/transaction/transaction.go
 * for correctly decoding the data
 */
func TestTransactionDataDecoder(t *testing.T) {
	ethformat := ethtxdata{
		AccountNonce: 0,
		Recipient:    nil,
		Payload:      nil,
		Amount:       new(big.Int),
		GasLimit:     new(big.Int),
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	cmd := "f86601843b9aca00835d161494a8863fc8ce3816411378685223c03daae9770ebb821388001ba02ca8516dfa28afdf0931b7993100bbf2497e8dc0e29bf4b9847e0708f96c98a2a05a0ce782188b3f42e5ad8cdb7b1f701a04f2509039fecc1c2c350c33b389d7b2"
	addr := "0xa8863fc8Ce3816411378685223C03DAae9770ebB"
	if err := Decode(bytes.NewReader(unhex(cmd)), &ethformat); err != nil {
		t.Fatalf("Decode error: %v", err)
	} else {
		fmt.Printf("Recipient:%v\n", ethformat.Recipient.Hex())
		if addr != ethformat.Recipient.Hex() {
			fmt.Println("Unmatch address!")
		}
		fmt.Printf("Recipient:%v\n", ethformat.Price)
	}

}

type testDecoder struct{ called bool }

func (t *testDecoder) DecodeRLP(s *Stream) error {
	if _, err := s.Uint(); err != nil {
		return err
	}
	t.called = true
	return nil
}

func TestDecodeDecoder(t *testing.T) {
	var s struct {
		T1 testDecoder
		T2 *testDecoder
		T3 **testDecoder
	}
	if err := Decode(bytes.NewReader(unhex("C3010203")), &s); err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	if !s.T1.called {
		t.Errorf("DecodeRLP was not called for (non-pointer) testDecoder")
	}

	if s.T2 == nil {
		t.Errorf("*testDecoder has not been allocated")
	} else if !s.T2.called {
		t.Errorf("DecodeRLP was not called for *testDecoder")
	}

	if s.T3 == nil || *s.T3 == nil {
		t.Errorf("**testDecoder has not been allocated")
	} else if !(*s.T3).called {
		t.Errorf("DecodeRLP was not called for **testDecoder")
	}
}

type byteDecoder byte

func (bd *byteDecoder) DecodeRLP(s *Stream) error {
	_, err := s.Uint()
	*bd = 255
	return err
}

func (bd byteDecoder) called() bool {
	return bd == 255
}

// This test verifies that the byte slice/byte array logic
// does not kick in for element types implementing Decoder.
func TestDecoderInByteSlice(t *testing.T) {
	var slice []byteDecoder
	if err := Decode(bytes.NewReader(unhex("C101")), &slice); err != nil {
		t.Errorf("unexpected Decode error %v", err)
	} else if !slice[0].called() {
		t.Errorf("DecodeRLP not called for slice element")
	}

	var array [1]byteDecoder
	if err := Decode(bytes.NewReader(unhex("C101")), &array); err != nil {
		t.Errorf("unexpected Decode error %v", err)
	} else if !array[0].called() {
		t.Errorf("DecodeRLP not called for array element")
	}
}

func ExampleDecode() {
	input, _ := hex.DecodeString("C90A1486666F6F626172")

	type example struct {
		A, B    uint
		private uint // private fields are ignored
		String  string
	}

	var s example
	err := Decode(bytes.NewReader(input), &s)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Decoded value: %#v\n", s)
	}
	// Output:
	// Decoded value: rlp.example{A:0xa, B:0x14, private:0x0, String:"foobar"}
}

func ExampleDecode_structTagNil() {
	// In this example, we'll use the "nil" struct tag to change
	// how a pointer-typed field is decoded. The input contains an RLP
	// list of one element, an empty string.
	input := []byte{0xC1, 0x80}

	// This type uses the normal rules.
	// The empty input string is decoded as a pointer to an empty Go string.
	var normalRules struct {
		String *string
	}
	Decode(bytes.NewReader(input), &normalRules)
	fmt.Printf("normal: String = %q\n", *normalRules.String)

	// This type uses the struct tag.
	// The empty input string is decoded as a nil pointer.
	var withEmptyOK struct {
		String *string `rlp:"nil"`
	}
	Decode(bytes.NewReader(input), &withEmptyOK)
	fmt.Printf("with nil tag: String = %v\n", withEmptyOK.String)

	// Output:
	// normal: String = ""
	// with nil tag: String = <nil>
}

func ExampleStream() {
	input, _ := hex.DecodeString("C90A1486666F6F626172")
	s := NewStream(bytes.NewReader(input), 0)

	// Check what kind of value lies ahead
	kind, size, _ := s.Kind()
	fmt.Printf("Kind: %v size:%d\n", kind, size)

	// Enter the list
	if _, err := s.List(); err != nil {
		fmt.Printf("List error: %v\n", err)
		return
	}

	// Decode elements
	fmt.Println(s.Uint())
	fmt.Println(s.Uint())
	fmt.Println(s.Bytes())

	// Acknowledge end of list
	if err := s.ListEnd(); err != nil {
		fmt.Printf("ListEnd error: %v\n", err)
	}
	// Output:
	// Kind: List size:9
	// 10 <nil>
	// 20 <nil>
	// [102 111 111 98 97 114] <nil>
}

func BenchmarkDecode(b *testing.B) {
	enc := encodeTestSlice(90000)
	b.SetBytes(int64(len(enc)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var s []uint
		r := bytes.NewReader(enc)
		if err := Decode(r, &s); err != nil {
			b.Fatalf("Decode error: %v", err)
		}
	}
}

func BenchmarkDecodeIntSliceReuse(b *testing.B) {
	enc := encodeTestSlice(100000)
	b.SetBytes(int64(len(enc)))
	b.ReportAllocs()
	b.ResetTimer()

	var s []uint
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(enc)
		if err := Decode(r, &s); err != nil {
			b.Fatalf("Decode error: %v", err)
		}
	}
}

func encodeTestSlice(n uint) []byte {
	s := make([]uint, n)
	for i := uint(0); i < n; i++ {
		s[i] = i
	}
	b, err := EncodeToBytes(s)
	if err != nil {
		panic(fmt.Sprintf("encode error: %v", err))
	}
	return b
}

func unhex(str string) []byte {
	b, err := hex.DecodeString(strings.Replace(str, " ", "", -1))
	if err != nil {
		panic(fmt.Sprintf("invalid hex string: %q", str))
	}
	return b
}
